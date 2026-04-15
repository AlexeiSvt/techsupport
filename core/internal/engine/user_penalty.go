// Package engine orchestrates the execution of multiple scoring rules.
package engine

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"techsupport/core/internal/constants"
	"techsupport/core/pkg"
	"techsupport/core/pkg/models"
	logPkg "techsupport/log/pkg"
)

// Ensure DeviceBruteforcePenaltyCalculator implements the ScoreCalculator interface.
var _ pkg.ScoreCalculator = (*DeviceBruteforcePenaltyCalculator)(nil)

// DeviceBruteforcePenaltyCalculator penalizes users if their device history shows
// abnormal growth or if their IP address has a low trust score.
type DeviceBruteforcePenaltyCalculator struct {
	mu  sync.RWMutex
	log logPkg.Logger

	totalCalculations uint64 // Atomic counter for tracking usage.
	penaltiesApplied  uint64 // Atomic counter for tracking detected threats.
}

// NewDeviceBruteforcePenaltyCalculator creates a new instance with the provided logger.
func NewDeviceBruteforcePenaltyCalculator(logger logPkg.Logger) *DeviceBruteforcePenaltyCalculator {
	return &DeviceBruteforcePenaltyCalculator{
		log: logger,
	}
}

// Calculate assesses security risks and applies penalties. 
// It caps the total penalty to constants.FullPenalty to maintain scoring balance.
func (c *DeviceBruteforcePenaltyCalculator) Calculate(ctx context.Context, user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	atomic.AddUint64(&c.totalCalculations, 1)

	// Respect context cancellation.
	select {
	case <-ctx.Done():
		return models.CalcResult{Status: "context_cancelled", Comment: ctx.Err().Error()}
	default:
	}

	c.mu.RLock()
	logger := c.log
	c.mu.RUnlock()

	penalty := 0.0
	comment := "No penalties"

	if logger != nil {
		logger.Debugw("starting security penalty calculation",
			"user_devices_count", len(user.UserClaim.Devices),
			"db_devices_count", len(db.Devices),
		)
	}

	// 1. Device Bruteforce Detection: Check if the number of claimed devices 
	// exceeds the known history plus a predefined threshold.
	if len(user.UserClaim.Devices) > len(db.Devices)+constants.NewDevicesThreshold {
		penalty += constants.BruteforcePenalty
		comment = "New devices threshold exceeded"
		
		if logger != nil {
			logger.Warnw("device bruteforce penalty applied",
				"u_count", len(user.UserClaim.Devices),
				"db_count", len(db.Devices),
				"penalty", constants.BruteforcePenalty,
			)
		}
	}

	// 2. IP Reputation Check: Integrate penalties from external IP intelligence.
	if user.UserClaim.IPInfo != nil {
		ipPenalty := user.UserClaim.IPInfo.GetPenaltyScore(logger)
		
		if ipPenalty > 0 {
			penalty += ipPenalty
			comment += fmt.Sprintf(" | IP Penalty: %.1f", ipPenalty)
		}
	} else {
		if logger != nil {
			logger.Warnw("IPInfo is missing during penalty calculation - potential data collection skip")
		}
	}

	// 3. Cap the penalty: Ensure the score doesn't become mathematically invalid.
	if penalty > float64(constants.FullPenalty) {
		if logger != nil {
			logger.Debugw("capping penalty to max value", "raw_penalty", penalty, "max", constants.FullPenalty)
		}
		penalty = float64(constants.FullPenalty)
	}

	status := "clear"
	if penalty > 0 {
		status = "penalty"
		atomic.AddUint64(&c.penaltiesApplied, 1)
		if logger != nil {
			logger.Infow("security penalty final result", "total_penalty", penalty, "status", status)
		}
	}

	return models.CalcResult{
		Name:    "Security Penalty",
		Code:    "penalty_score",
		Value:   penalty,
		Weight:  1.0, // Penalties usually have a fixed impact.
		Result:  penalty,
		Status:  status,
		Comment: comment,
	}
}

// GetStats returns metrics for total checks and penalties triggered.
func (c *DeviceBruteforcePenaltyCalculator) GetStats() (uint64, uint64) {
	return atomic.LoadUint64(&c.totalCalculations), atomic.LoadUint64(&c.penaltiesApplied)
}

// SetLogger safely updates the logger at runtime.
func (c *DeviceBruteforcePenaltyCalculator) SetLogger(newLogger logPkg.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.log = newLogger
}