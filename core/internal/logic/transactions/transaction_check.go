package transactions

import (
	"fmt"
	"techsupport/core/internal/constants"
	"techsupport/core/internal/models"
	"techsupport/core/pkg"
)

var _ pkg.ScoreCalculator = (*FirstTransactionScoreCalculator)(nil)

type FirstTransactionScoreCalculator struct{}

type rawTxResult struct {
	Value   float64
	Status  string
	Comment string
}

func (c FirstTransactionScoreCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	if weights.FirstTransaction <= 0 {
		return models.CalcResult{
			Name: "First Transaction Check",
			Code: "first_tx",
			Status: "skipped",
			Comment: "Weight is zero",
		}
	}

	res := checkFirstTransaction(db, user.UserClaim)

	return models.CalcResult{
		Name:    "First Transaction Analysis",
		Code:    "first_tx",
		Value:   res.Value,
		Weight:  weights.FirstTransaction,
		Result:  res.Value * weights.FirstTransaction, 
		Status:  res.Status,
		Comment: res.Comment,
	}
}

func checkFirstTransaction(dbRecord models.DBRecord, userClaim models.UserClaim) rawTxResult {
	tx := userClaim.FirstTransaction
	var baseScore float64
	anomalyCount := 0
	comment := ""

	scoreFirst := calculateWindowScore(tx, dbRecord.UserHistory.FirstWindow)
	scoreLast := calculateWindowScore(tx, dbRecord.UserHistory.LastWindow)
	
	baseScore = max(scoreFirst, scoreLast)

	if baseScore == 0 {
		if isRegionAndDeviceKnown(tx, dbRecord.UserHistory) {
			baseScore = constants.PartialMatch
			comment = "Known region and device found in history"
		} else {
			return rawTxResult{Value: 0, Status: "no_match", Comment: "Transaction environment is completely unknown"}
		}
	} else {
		comment = fmt.Sprintf("Matched via session history (base: %.2f)", baseScore)
	}

	if isHighFrequencyTransaction(dbRecord.UserHistory.AllPayments, tx) {
		anomalyCount++
		baseScore *= constants.MostlyMatch
		comment += " | Anomaly: High frequency"
	}

	if isSuddenHighDonation(tx, dbRecord.UserHistory) {
		anomalyCount++
		baseScore *= constants.MostlyMatch
		comment += " | Anomaly: Sudden high donation"
	}

	if !isRegionAndDeviceKnown(tx, dbRecord.UserHistory) {
		anomalyCount++
		baseScore *= constants.MostlyMatch
		comment += " | Anomaly: Unknown region/device combination"
	}

	if anomalyCount >= 2 {
		return rawTxResult{
			Value:   0,
			Status:  "anomaly_block",
			Comment: fmt.Sprintf("Blocked: %d anomalies detected | %s", anomalyCount, comment),
		}
	}

	status := "match"
	if baseScore < 1.0 {
		status = "partial"
	}

	return rawTxResult{
		Value:   baseScore,
		Status:  status,
		Comment: comment,
	}
}