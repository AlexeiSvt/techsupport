// Package transactions provides validation logic for financial activities.
package transactions

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"techsupport/core/internal/constants"
	"techsupport/core/internal/errors"
	"techsupport/core/internal/logic"
	"techsupport/core/pkg"
	"techsupport/core/pkg/models"
	logPkg "techsupport/log/pkg"
)

// Ensure FirstTransactionScoreCalculator implements the ScoreCalculator interface.
var _ pkg.ScoreCalculator = (*FirstTransactionScoreCalculator)(nil)

// FirstTransactionScoreCalculator evaluates the validity of the user's first transaction.
// It integrates multiple validation checks to detect anomalies and fraud patterns.
type FirstTransactionScoreCalculator struct {
	mu        sync.RWMutex // Protects the logger and validator during runtime updates.
	log       logPkg.Logger
	validator *TxValidator

	totalCalculations uint64 // Atomic counter for total analysis attempts.
	blockedAnomalies  uint64 // Atomic counter for transactions blocked due to anomalies.
}

// NewFirstTransactionScoreCalculator initializes a new calculator with a validator.
func NewFirstTransactionScoreCalculator(logger logPkg.Logger) *FirstTransactionScoreCalculator {
	if logger == nil {
		return nil
	}
	
	return &FirstTransactionScoreCalculator{
		log:       logger,
		validator: NewTxValidator(logger),
	}
}

// Calculate orchestrates the transaction analysis process.
func (c *FirstTransactionScoreCalculator) Calculate(ctx context.Context, claim models.UserClaim, support models.SupportContext, db models.DBRecord, weights models.Weights) models.CalcResult {
	atomic.AddUint64(&c.totalCalculations, 1)

	// Early exit if context is cancelled.
	select {
	case <-ctx.Done():
		return models.CalcResult{Status: "context_cancelled", Comment: ctx.Err().Error()}
	default:
	}

	c.mu.RLock()
	logger := c.log
	validator := c.validator
	c.mu.RUnlock()

	// Lazy initialization of validator if it's missing.
	if validator == nil {
		validator = NewTxValidator(logger)
	}

	// Skip calculation if the policy weight is zero.
	if weights.FirstTransaction <= 0 {
		if logger != nil {
			logger.Debugw("skipping transaction check", "reason", "weight is zero", "ub_ticket_id", claim.UBTicketID)
		}
		return models.CalcResult{
			Name:    "First Transaction Check",
			Code:    "first_tx",
			Status:  errors.StatusSkipped,
			Comment: "Weight is zero",
		}
	}

	if logger != nil {
		logger.Debugw("starting first transaction analysis", 
			"tx_id", claim.FirstTransaction.TransactionID,
			"ub_ticket_id", claim.UBTicketID,
		)
	}

	// EXECUTION: Passing 'support' instead of 'db' to access History context.
	res := c.checkFirstTransaction(ctx, claim, support, validator)

	// If blocked by anomalies, track it in metrics.
	if res.Status == errors.StatusAnomalyBlock {
		atomic.AddUint64(&c.blockedAnomalies, 1)
	}

	calcRes := models.CalcResult{
		Name:    "First Transaction Analysis",
		Code:    "first_tx",
		Value:   res.Value,
		Weight:  weights.FirstTransaction,
		Result:  res.Value * weights.FirstTransaction,
		Status:  res.Status,
		Comment: res.Comment,
	}

	if logger != nil {
		logger.Infow("transaction analysis finished",
			"status", calcRes.Status,
			"final_value", calcRes.Value,
			"ub_ticket_id", claim.UBTicketID,
		)
	}

	return calcRes
}

// checkFirstTransaction executes specific validation rules and applies penalties for anomalies.
func (c *FirstTransactionScoreCalculator) checkFirstTransaction(
	ctx context.Context, 
	claim models.UserClaim, 
	support models.SupportContext, 
	v *TxValidator,
) logic.RawTxResult {
	tx := claim.FirstTransaction
	var baseScore float64
	anomalyCount := 0
	comment := ""

	// Pulling slices directly from the flattened SupportContext.History.
	firstWindow := support.History.FirstWindow
	lastWindow := support.History.LastWindow
	allPayments := support.History.AllPayments

	// 1. Geography & Device Match Scoring.
	scoreFirst := v.CalculateWindowScore(ctx, tx, firstWindow)
	scoreLast := v.CalculateWindowScore(ctx, tx, lastWindow)

	baseScore = mathMax(scoreFirst, scoreLast)

	// Fallback to boolean check if historical weight scoring is zero.
	if baseScore == 0 {
		if v.IsRegionAndDeviceKnown(ctx, tx, firstWindow, lastWindow, claim.UBTicketID) {
			baseScore = constants.PartialMatch
			comment = "Known region and device found in history"
		} else {
			return logic.RawTxResult{
				Value:   0,
				Status:  errors.StatusNoMatch,
				Comment: "Transaction environment is completely unknown",
			}
		}
	} else {
		comment = fmt.Sprintf("Matched via session history (base: %.2f)", baseScore)
	}

	// 2. Anomaly Detection Phase.

	// Anomaly: High Frequency (Velocity check).
	if v.IsHighFrequencyTransaction(ctx, allPayments, tx, claim.UBTicketID) {
		anomalyCount++
		baseScore *= constants.MostlyMatch
		comment += " | Anomaly: High frequency"
	}

	// Anomaly: Sudden High Donation (Amount check).
	if v.IsSuddenHighDonation(ctx, tx, allPayments, claim.UBTicketID) {
		anomalyCount++
		baseScore *= constants.MostlyMatch
		comment += " | Anomaly: Sudden high donation"
	}

	// Anomaly: Environmental Mismatch.
	if !v.IsRegionAndDeviceKnown(ctx, tx, firstWindow, lastWindow, claim.UBTicketID) {
		anomalyCount++
		baseScore *= constants.MostlyMatch
		comment += " | Anomaly: Unknown region/device combination"
	}

	// 3. Final Decision Logic.

	// Anti-Fraud Policy: Cumulative anomalies result in a hard block.
	if anomalyCount >= 2 {
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

// GetStats returns metrics for total processed transactions and total anomaly blocks.
func (c *FirstTransactionScoreCalculator) GetStats() (uint64, uint64) {
	return atomic.LoadUint64(&c.totalCalculations), atomic.LoadUint64(&c.blockedAnomalies)
}

// SetLogger safely updates the logger and internal validator at runtime.
func (c *FirstTransactionScoreCalculator) SetLogger(newLogger logPkg.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.log = newLogger
	if c.validator != nil {
		c.validator.SetLogger(newLogger)
	}
}

// mathMax is a helper for float64 comparison.
func mathMax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}