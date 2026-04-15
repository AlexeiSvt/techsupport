// Package logic provides calculation engines for scoring user data.
// It includes thread-safe implementations for comparing various user attributes 
// in high-concurrency environments.
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

// Ensure FirstDeviceCalculator implements the ScoreCalculator interface at compile time.
var _ pkg.ScoreCalculator = (*FirstDeviceCalculator)(nil)

// FirstDeviceCalculator handles the comparison and scoring logic for user device identifiers.
// It is fully thread-safe and optimized for concurrent read access to its logger.
type FirstDeviceCalculator struct {
	mu  sync.RWMutex // Protects the logger instance during runtime updates.
	log logPkg.Logger

	totalCalculations uint64 // Atomic counter for the total number of calculation attempts.
	matchCount        uint64 // Atomic counter for the total number of successful matches.
}

// NewFirstDeviceCalculator initializes and returns a new FirstDeviceCalculator with the given logger.
func NewFirstDeviceCalculator(logger logPkg.Logger) *FirstDeviceCalculator {
	return &FirstDeviceCalculator{
		log: logger,
	}
}

// Calculate compares user and database device records to produce a scoring result.
// This method is designed to be executed within a goroutine and respects context cancellation.
func (c *FirstDeviceCalculator) Calculate(ctx context.Context, user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	// Increment total calculations atomically to avoid race conditions.
	atomic.AddUint64(&c.totalCalculations, 1)

	// Early exit if the parent context has been cancelled or timed out.
	select {
	case <-ctx.Done():
		return models.CalcResult{
			Status:  "context_cancelled",
			Comment: ctx.Err().Error(),
		}
	default:
	}

	// Capture the logger instance under a Read-Lock.
	c.mu.RLock()
	logger := c.log
	c.mu.RUnlock()

	userDev := user.UserClaim.FirstDevice
	dbDev := db.FirstDevice

	if logger != nil {
		logger.Debugw("checking first device match", "user_device", userDev, "db_device", dbDev)
	}

	// Execute core comparison logic.
	res := c.checkFirstDevice(ctx, userDev, dbDev)

	// Increment match counter atomically if the comparison was successful.
	if res.Status == errors.StatusMatch {
		atomic.AddUint64(&c.matchCount, 1)
	}

	// Construct the final calculation result.
	calcRes := models.CalcResult{
		Name:    "First Device Match",
		Code:    "first_device",
		Value:   res.Value,
		Weight:  weights.FirstDevice,
		Result:  res.Value * weights.FirstDevice,
		Status:  res.Status,
		Comment: res.Comment,
	}

	if logger != nil {
		logger.Infow("device score calculated",
			"status", calcRes.Status,
			"result", calcRes.Result,
		)
	}

	return calcRes
}

// checkFirstDevice is an internal helper that performs case-insensitive comparison for device IDs.
// It trims whitespace from inputs and handles empty data scenarios.
func (c *FirstDeviceCalculator) checkFirstDevice(ctx context.Context, userDev, dbDev string) rawCheckResult {
	u := strings.TrimSpace(userDev)
	d := strings.TrimSpace(dbDev)

	// Read-Locking to safely access the logger for warning messages.
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Validate inputs to ensure both device IDs are present.
	if u == "" || d == "" {
		if c.log != nil {
			c.log.Warnw(errors.ErrEmptyDeviceData.Error(), "u_dev_empty", u == "", "db_dev_empty", d == "")
		}

		return rawCheckResult{
			Value:   constants.NoMatch,
			Status:  errors.StatusNoData,
			Comment: "Missing first device data",
		}
	}

	// Perform case-insensitive match using strings.EqualFold.
	if strings.EqualFold(u, d) {
		return rawCheckResult{
			Value:   constants.IdealMatch,
			Status:  errors.StatusMatch,
			Comment: "First device matches exactly",
		}
	}

	return rawCheckResult{
		Value:   constants.NoMatch,
		Status:  errors.StatusNoMatch,
		Comment: "First device mismatch",
	}
}

// GetStats returns current metrics: total calculations performed and total matches found.
// This data is retrieved atomically without blocking the main execution flow.
func (c *FirstDeviceCalculator) GetStats() (uint64, uint64) {
	return atomic.LoadUint64(&c.totalCalculations), atomic.LoadUint64(&c.matchCount)
}

// SetLogger allows for safe runtime updates of the logger instance.
// It uses a Write-Lock to ensure exclusive access during the pointer swap.
func (c *FirstDeviceCalculator) SetLogger(newLogger logPkg.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.log = newLogger
}