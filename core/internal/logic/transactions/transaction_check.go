package transactions

import (
	"fmt"
	"techsupport/core/internal/constants"
	"techsupport/core/internal/errors"
	"techsupport/core/internal/models"
	"techsupport/core/pkg"
	logPkg "techsupport/log/pkg"
    "techsupport/core/internal/logic"
)

var _ pkg.ScoreCalculator = (*FirstTransactionScoreCalculator)(nil)

type FirstTransactionScoreCalculator struct {
	Log       logPkg.Logger
	Validator *TxValidator 
}

func (c *FirstTransactionScoreCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	if c.Validator == nil {
		c.Validator = &TxValidator{Log: c.Log}
	}

	if weights.FirstTransaction <= 0 {
		if c.Log != nil {
			c.Log.Debugw("skipping transaction check", "reason", "weight is zero")
		}
		return models.CalcResult{
			Name:    "First Transaction Check",
			Code:    "first_tx",
			Status:  errors.StatusSkipped,
			Comment: "Weight is zero",
		}
	}

	if c.Log != nil {
		c.Log.Debugw("starting first transaction analysis", "tx_id", user.UserClaim.FirstTransaction.TransactionID)
	}

	res := c.checkFirstTransaction(db, user.UserClaim)

	calcRes := models.CalcResult{
		Name:    "First Transaction Analysis",
		Code:    "first_tx",
		Value:   res.Value,
		Weight:  weights.FirstTransaction,
		Result:  res.Value * weights.FirstTransaction,
		Status:  res.Status,
		Comment: res.Comment,
	}

	if c.Log != nil {
		c.Log.Infow("transaction analysis finished", 
			"status", calcRes.Status, 
			"final_value", calcRes.Value,
		)
	}

	return calcRes
}

func (c *FirstTransactionScoreCalculator) checkFirstTransaction(dbRecord models.DBRecord, userClaim models.UserClaim) logic.RawTxResult {
	tx := userClaim.FirstTransaction
	var baseScore float64
	anomalyCount := 0
	comment := ""

	scoreFirst := c.Validator.CalculateWindowScore(tx, dbRecord.UserHistory.FirstWindow)
	scoreLast := c.Validator.CalculateWindowScore(tx, dbRecord.UserHistory.LastWindow)
	
	baseScore = max(scoreFirst, scoreLast)

	if baseScore == 0 {
		if c.Validator.IsRegionAndDeviceKnown(tx, dbRecord.UserHistory) {
			baseScore = constants.PartialMatch
			comment = "Known region and device found in history"
		} else {
			if c.Log != nil {
				c.Log.Warnw("transaction environment unknown", "tx_id", tx.TransactionID)
			}
			return logic.RawTxResult{Value: 0, Status: errors.StatusNoMatch, Comment: "Transaction environment is completely unknown"}
		}
	} else {
		comment = fmt.Sprintf("Matched via session history (base: %.2f)", baseScore)
	}

	if c.Validator.IsHighFrequencyTransaction(dbRecord.UserHistory.AllPayments, tx) {
		anomalyCount++
		baseScore *= constants.MostlyMatch
		comment += " | Anomaly: High frequency"
	}

	if c.Validator.IsSuddenHighDonation(tx, dbRecord.UserHistory) {
		anomalyCount++
		baseScore *= constants.MostlyMatch
		comment += " | Anomaly: Sudden high donation"
	}

	if !c.Validator.IsRegionAndDeviceKnown(tx, dbRecord.UserHistory) {
		anomalyCount++
		baseScore *= constants.MostlyMatch
		comment += " | Anomaly: Unknown region/device combination"
	}

	if anomalyCount >= 2 {
		if c.Log != nil {
			c.Log.Errorw("transaction blocked due to multiple anomalies", 
				"tx_id", tx.TransactionID, 
				"anomaly_count", anomalyCount,
			)
		}
		return logic.RawTxResult{
			Value:   0,
			Status:  errors.StatusAnomalyBlock,
			Comment: fmt.Sprintf("Blocked: %d anomalies detected | %s", anomalyCount, comment),
		}
	}

	status := errors.StatusMatch
	if baseScore < 1.0 {
		status = errors.StatusPartial
	}

	return logic.RawTxResult{
		Value:   baseScore,
		Status:  status,
		Comment: comment,
	}
}