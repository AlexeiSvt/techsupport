// Package logic provides calculation engines for scoring user data.
package logic

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"

	"techsupport/core/internal/constants"
	"techsupport/core/internal/errors"
	"techsupport/core/pkg"
	"techsupport/core/pkg/models"
	logPkg "techsupport/log/pkg"
)

// Ensure FirstPhoneCalculator implements the ScoreCalculator interface at compile time.
var _ pkg.ScoreCalculator = (*FirstPhoneCalculator)(nil)

// FirstPhoneCalculator handles comparison and scoring for user phone numbers.
// It supports full matches and partial suffix matches to account for international formatting.
// It is fully thread-safe and optimized for concurrent execution.
type FirstPhoneCalculator struct {
	mu  sync.RWMutex // Protects the logger instance during runtime updates.
	log logPkg.Logger

	totalCalculations uint64 // Atomic counter for the total number of calculation attempts.
	matchCount        uint64 // Atomic counter for the total number of successful matches.
}

// NewFirstPhoneCalculator initializes and returns a new FirstPhoneCalculator with the given logger.
func NewFirstPhoneCalculator(logger logPkg.Logger) *FirstPhoneCalculator {
	return &FirstPhoneCalculator{
		log: logger,
	}
}

// Calculate compares user and database phone records to determine a match score.
// It respects context cancellation and updates internal metrics atomically.
func (c *FirstPhoneCalculator) Calculate(ctx context.Context, user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
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
		logger.Debugw("calculating phone match score",
			"user_phone_raw", user.UserClaim.Phone,
			"db_phone_raw", db.Phone,
		)
	}

	// Execute core comparison logic.
	res := c.checkPhone(ctx, user.UserClaim.Phone, db.Phone)

	// If the calculation result indicates a successful match, increment match counter.
	if res.Status == errors.StatusMatch {
		atomic.AddUint64(&c.matchCount, 1)
	}

	calcRes := models.CalcResult{
		Name:    "First Phone Match",
		Code:    "first_phone",
		Value:   res.Value,
		Weight:  weights.Phone,
		Result:  res.Value * weights.Phone,
		Status:  res.Status,
		Comment: res.Comment,
	}

	if logger != nil {
		logger.Infow("phone score calculation finished",
			"status", calcRes.Status,
			"result", calcRes.Result,
		)
	}

	return calcRes
}

// checkPhone cleans and compares two phone number strings.
// It performs a full match first, followed by a suffix match if full comparison fails.
func (c *FirstPhoneCalculator) checkPhone(ctx context.Context, userPhone, dbPhone string) rawCheckResult {
	// Read-Locking for safe access to the logger.
	c.mu.RLock()
	defer c.mu.RUnlock()

	if userPhone == "" || dbPhone == "" {
		if c.log != nil {
			c.log.Warnw(errors.ErrEmptyPhoneData.Error(), "u_phone_empty", userPhone == "", "db_phone_empty", dbPhone == "")
		}
		return rawCheckResult{
			Value:   constants.NoMatch,
			Status:  errors.StatusNoData,
			Comment: "Missing phone data",
		}
	}

	// Clean numbers to contain only digits.
	uPhone := CleanPhoneNumber(userPhone)
	dPhone := CleanPhoneNumber(dbPhone)

	// Attempt exact match.
	if uPhone == dPhone {
		return rawCheckResult{
			Value:   constants.IdealMatch,
			Status:  errors.StatusMatch,
			Comment: "Full phone number match",
		}
	}

	// Check minimum length requirement before partial matching.
	if len(uPhone) < constants.MinLen || len(dPhone) < constants.MinLen {
		if c.log != nil {
			c.log.Debugw(errors.ErrPhoneTooShort.Error(), "u_len", len(uPhone), "d_len", len(dPhone))
		}
		return rawCheckResult{
			Value:   constants.NoMatch,
			Status:  errors.StatusNoMatch,
			Comment: "Phone length below minimum",
		}
	}

	// Calculate suffix match (e.g., last 10 digits) to handle country code variations.
	checkLen := 10
	if len(dPhone) < checkLen {
		checkLen = len(dPhone)
	}
	if len(uPhone) < checkLen {
		checkLen = len(uPhone)
	}

	if uPhone[len(uPhone)-checkLen:] == dPhone[len(dPhone)-checkLen:] {
		if c.log != nil {
			c.log.Debugw("partial phone match detected", "check_len", checkLen)
		}
		return rawCheckResult{
			Value:   constants.PartialMatch,
			Status:  errors.StatusPartial,
			Comment: "Matched by last digits",
		}
	}

	return rawCheckResult{
		Value:   constants.NoMatch,
		Status:  errors.StatusNoMatch,
		Comment: "Numbers do not match",
	}
}

// CleanPhoneNumber removes all non-numeric characters from the input string.
func CleanPhoneNumber(phone string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)
}

// GetStats returns the total number of calculations and matches processed by this instance.
func (c *FirstPhoneCalculator) GetStats() (uint64, uint64) {
	return atomic.LoadUint64(&c.totalCalculations), atomic.LoadUint64(&c.matchCount)
}

// SetLogger allows for safe runtime updates of the logger instance.
func (c *FirstPhoneCalculator) SetLogger(newLogger logPkg.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.log = newLogger
}