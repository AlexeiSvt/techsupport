// Package logic provides calculation engines for scoring user data.
package logic

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"techsupport/core/internal/constants"
	"techsupport/core/internal/errors"
	"techsupport/core/pkg"
	"techsupport/core/pkg/models"
	logPkg "techsupport/log/pkg"
)

// Ensure RegDateCalculator implements the ScoreCalculator interface at compile time.
var _ pkg.ScoreCalculator = (*RegDateCalculator)(nil)

// RegDateCalculator compares the user's claimed registration date against the database record.
// It applies scoring based on the temporal distance (difference) between the two dates.
// It is fully thread-safe and respects context-based execution.
type RegDateCalculator struct {
	mu  sync.RWMutex // Protects the logger instance during runtime updates.
	log logPkg.Logger

	totalCalculations uint64 // Atomic counter for the total number of calculation attempts.
	matchCount        uint64 // Atomic counter for the total number of successful matches.
}

// NewRegDateCalculator initializes and returns a new RegDateCalculator with the given logger.
func NewRegDateCalculator(logger logPkg.Logger) *RegDateCalculator {
	return &RegDateCalculator{
		log: logger,
	}
}

// Calculate determines the match score based on the registration date difference.
// It handles anomaly detection for high discrepancies and respects context cancellation.
func (c *RegDateCalculator) Calculate(ctx context.Context, user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	// Increment total calculations atomically.
	atomic.AddUint64(&c.totalCalculations, 1)

	// Early exit if the context is cancelled or timed out.
	select {
	case <-ctx.Done():
		return models.CalcResult{
			Status:  "context_cancelled",
			Comment: ctx.Err().Error(),
		}
	default:
	}

	// Capture the logger instance safely under a Read-Lock.
	c.mu.RLock()
	logger := c.log
	c.mu.RUnlock()

	if logger != nil {
		logger.Debugw("calculating registration date match",
			"user_reg_date", user.UserClaim.RegDate,
			"db_reg_date", db.RegDate,
		)
	}

	// Execute core comparison logic.
	res := c.checkCreationAge(ctx, user.UserClaim.RegDate, db.RegDate)

	// If the calculation result indicates a successful match, increment match counter.
	if res.Status == errors.StatusMatch {
		atomic.AddUint64(&c.matchCount, 1)
	}

	calcRes := models.CalcResult{
		Name:    "Registration Date Match",
		Code:    "reg_date",
		Value:   res.Value,
		Weight:  weights.RegDate,
		Result:  res.Value * weights.RegDate,
		Status:  res.Status,
		Comment: res.Comment,
	}

	if logger != nil {
		logger.Infow("reg date score calculated",
			"status", calcRes.Status,
			"diff_value", calcRes.Value,
		)
	}

	return calcRes
}

// checkCreationAge calculates the distance between two time.Time objects in months.
// It applies a tolerance range and identifies anomalies if the discrepancy is too large.
func (c *RegDateCalculator) checkCreationAge(ctx context.Context, userClaim time.Time, dbRecord time.Time) rawCheckResult {
	// Acquire Read-Lock for safe logging access.
	c.mu.RLock()
	defer c.mu.RUnlock()

	if userClaim.IsZero() || dbRecord.IsZero() {
		if c.log != nil {
			c.log.Warnw(errors.ErrEmptyRegDate.Error(),
				"user_date_zero", userClaim.IsZero(),
				"db_date_zero", dbRecord.IsZero(),
			)
		}
		return rawCheckResult{
			Value:   constants.NoMatch,
			Status:  errors.StatusNoData,
			Comment: "Missing registration date in claim or database",
		}
	}

	// Calculate absolute difference in hours and convert to months.
	diffHours := math.Abs(userClaim.Sub(dbRecord).Hours())
	diffMonths := diffHours / constants.AvgAmountOfHoursInMonth
	toleranceMonths := constants.ToleranceHours / constants.AvgAmountOfHoursInMonth

	commentBase := fmt.Sprintf("Diff: %.1f months (tolerance: %.1f)", diffMonths, toleranceMonths)

	if c.log != nil {
		c.log.Debugw("date difference calculated", "diff_months", diffMonths, "tolerance", toleranceMonths)
	}

	// 1. Check for Ideal Match range.
	if diffMonths <= constants.IdealMatchofMonths+toleranceMonths {
		return rawCheckResult{
			Value:   constants.IdealMatch,
			Status:  errors.StatusMatch,
			Comment: commentBase + " - Within ideal range",
		}
	}

	// 2. Check for Partial Match range.
	if diffMonths <= constants.PartialMatchofMonths+toleranceMonths {
		return rawCheckResult{
			Value:   constants.PartialMatch,
			Status:  errors.StatusPartial,
			Comment: commentBase + " - Within partial range",
		}
	}

	// 3. Detect Anomaly (e.g., difference > 1 year).
	// This usually results in a negative score multiplier in many scoring systems.
	if diffMonths-toleranceMonths >= constants.OneYearofMonths {
		if c.log != nil {
			c.log.Warnw(errors.ErrRegDateAnomaly.Error(), "diff_months", diffMonths)
		}

		return rawCheckResult{
			Value:   -constants.IdealMatch,
			Status:  errors.StatusAnomaly,
			Comment: commentBase + " - High discrepancy (over 1 year)",
		}
	}

	// Default case: no match within acceptable parameters.
	return rawCheckResult{
		Value:   constants.NoMatch,
		Status:  errors.StatusNoMatch,
		Comment: commentBase + " - Outside acceptable ranges",
	}
}

// GetStats returns the total number of calculations and matches processed by this instance.
func (c *RegDateCalculator) GetStats() (uint64, uint64) {
	return atomic.LoadUint64(&c.totalCalculations), atomic.LoadUint64(&c.matchCount)
}

// SetLogger allows for safe runtime updates of the logger instance.
func (c *RegDateCalculator) SetLogger(newLogger logPkg.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.log = newLogger
}