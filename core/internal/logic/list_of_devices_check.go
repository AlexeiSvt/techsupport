// Package logic provides calculation engines for scoring user data.
// It includes thread-safe implementations for comparing various user attributes 
// in high-concurrency environments.
package logic

import (
    "context"
    "fmt"
    "strings"
    "sync"
    "sync/atomic"
    "time"

    "techsupport/core/internal/constants"
    "techsupport/core/internal/errors"
    "techsupport/core/pkg"
    "techsupport/core/pkg/models"
    logPkg "techsupport/log/pkg"
)

// Ensure DevicesCalculator implements the ScoreCalculator interface at compile time.
var _ pkg.ScoreCalculator = (*DevicesCalculator)(nil)

// DevicesCalculator analyzes the history of devices associated with a use.
// It applies matching ratios and penalties based on account age and device count.
// It is fully thread-safe and optimized for concurrent execution.
type DevicesCalculator struct {
    mu  sync.RWMutex // Protects the logger instance during runtime updates.
    log logPkg.Logger

    totalCalculations uint64 // Atomic counter for the total number of calculation attempts.
    matchCount        uint64 // Atomic counter for the total number of successful matches.
}

// NewDevicesCalculator initializes and returns a new DevicesCalculator with the given logger.
func NewDevicesCalculator(logger logPkg.Logger) *DevicesCalculator {
    return &DevicesCalculator{
        log: logger,
    }
}

// Calculate performs the device history scoring by comparing user claims against database records.
// It respects context cancellation and updates internal metrics atomically.
func (c *DevicesCalculator) Calculate(ctx context.Context, claim models.UserClaim, support models.SupportContext, db models.DBRecord, weights models.Weights) models.CalcResult {
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

    // Capture the logger instance safely.
    c.mu.RLock()
    logger := c.log
    c.mu.RUnlock()

    // Note: In the new flat architecture, device lists should be pre-fetched 
    // or passed via a specialized context/field. Assuming they are now 
    // part of the extended claim/db logic.
    userDevices := []string{claim.FirstDeviceName} // Simplified for flat model context.
    dbDevices := []string{db.FirstDevice}          // In a real graph scenario, these come from relationships.

    if logger != nil {
        logger.Debugw("calculating all devices match",
            "ub_ticket_id", claim.UBTicketID,
            "user_devices_count", len(userDevices),
            "db_devices_count", len(dbDevices),
        )
    }

    // Execute core comparison logic.
    res := c.checkAllDevices(ctx, userDevices, dbDevices, db.RegDate)

    // If the calculation result indicates a successful match, increment match counter.
    if res.Status == errors.StatusMatch {
        atomic.AddUint64(&c.matchCount, 1)
    }

    calcRes := models.CalcResult{
        Name:    "Device History Match",
        Code:    "devices_list",
        Value:   res.Value,
        Weight:  weights.Devices,
        Result:  res.Value * weights.Devices,
        Status:  res.Status,
        Comment: res.Comment,
    }

    if logger != nil {
        logger.Infow("device history calculation finished",
            "status", calcRes.Status,
            "result", calcRes.Result,
            "ub_ticket_id", claim.UBTicketID,
        )
    }

    return calcRes
}

// checkAllDevices calculates the match ratio between two lists of devices.
// It builds a lookup map to perform the intersection in O(n) time.
func (c *DevicesCalculator) checkAllDevices(ctx context.Context, userDevices []string, dbDevices []string, regDate time.Time) rawCheckResult {
    // Input validation for empty lists.
    if len(userDevices) == 0 || len(dbDevices) == 0 {
        return rawCheckResult{
            Value:   constants.NoMatch,
            Status:  errors.StatusNoData,
            Comment: "Device list is empty",
        }
    }

    // Build a lookup map from database devices for efficient matching.
    deviceMap := make(map[string]struct{}, len(dbDevices))
    for _, d := range dbDevices {
        deviceMap[strings.ToLower(d)] = struct{}{}
    }

    matches := 0.0
    for _, u := range userDevices {
        lowU := strings.ToLower(u)
        if _, exists := deviceMap[lowU]; exists {
            matches++
            // Delete to avoid double-counting if the user list has duplicates.
            delete(deviceMap, lowU)
        }
    }

    if matches == 0 {
        return rawCheckResult{
            Value:   constants.NoMatch,
            Status:  errors.StatusNoMatch,
            Comment: "No common devices found",
        }
    }

    // Calculate the ratio based on the total number of devices in the database.
    matchRatio := matches / float64(len(dbDevices))
    if matchRatio > constants.IdealMatch {
        matchRatio = constants.IdealMatch
    }

    // Apply penalty based on the age of the account and the number of devices.
    penalty := c.getDevicePenalty(regDate, len(dbDevices))
    finalValue := matchRatio * penalty

    comment := fmt.Sprintf("Matches: %.0f/%d. Ratio: %.2f. Penalty: %.2f",
        matches, len(dbDevices), matchRatio, penalty)

    status := errors.StatusPartial
    if finalValue >= constants.IdealMatch {
        status = errors.StatusMatch
    }

    return rawCheckResult{
        Value:   finalValue,
        Status:  status,
        Comment: comment,
    }
}

// getDevicePenalty determines a multiplier based on device "churn."
// If an account has too many devices relative to its age, a penalty is returned.
func (c *DevicesCalculator) getDevicePenalty(regDate time.Time, devicesCount int) float64 {
    durationSinceReg := time.Since(regDate)

    if durationSinceReg < 0 {
        return constants.IdealMatch
    }

    yearsActive := durationSinceReg.Hours() / constants.OneYearInHours
    if yearsActive < 1 {
        yearsActive = 1
    }

    // Basic anti-fraud rule: allow 4 new devices per year, with a base of 10.
    allowedDevices := 10
    calculatedAllowed := int(yearsActive * 4)
    if calculatedAllowed > allowedDevices {
        allowedDevices = calculatedAllowed
    }

    if devicesCount > allowedDevices {
        return constants.MostlyMatch
    }

    return constants.IdealMatch
}

// GetStats returns the total number of calculations and matches processed by this instance.
func (c *DevicesCalculator) GetStats() (uint64, uint64) {
    return atomic.LoadUint64(&c.totalCalculations), atomic.LoadUint64(&c.matchCount)
}

// SetLogger allows for safe runtime updates of the logger instance.
func (c *DevicesCalculator) SetLogger(newLogger logPkg.Logger) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.log = newLogger
}