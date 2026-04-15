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

// Ensure calculators implement the ScoreCalculator interface at compile time.
var _ pkg.ScoreCalculator = (*RegCountryCalculator)(nil)
var _ pkg.ScoreCalculator = (*RegCityCalculator)(nil)

// RegCountryCalculator handles scoring based on the user's registration country.
type RegCountryCalculator struct {
	mu  sync.RWMutex
	log logPkg.Logger

	totalCalculations uint64
	matchCount        uint64
}

// Calculate compares the registration country from user claims and database records.
func (c *RegCountryCalculator) Calculate(ctx context.Context, user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	atomic.AddUint64(&c.totalCalculations, 1)

	select {
	case <-ctx.Done():
		return models.CalcResult{Status: "context_cancelled", Comment: ctx.Err().Error()}
	default:
	}

	c.mu.RLock()
	logger := c.log
	c.mu.RUnlock()

	if logger != nil {
		logger.Debugw("calculating country match", "u_country", user.UserClaim.RegCountry, "db_country", db.RegCountry)
	}

	res := checkLocationGeneric(ctx, logger, user.UserClaim.RegCountry, db.RegCountry, "Country")

	if res.Status == errors.StatusMatch {
		atomic.AddUint64(&c.matchCount, 1)
	}

	return models.CalcResult{
		Name:    "Registration Country Match",
		Code:    "reg_country",
		Value:   res.Value,
		Weight:  weights.RegCountry,
		Result:  res.Value * weights.RegCountry,
		Status:  res.Status,
		Comment: res.Comment,
	}
}

// RegCityCalculator handles scoring based on the user's registration city.
type RegCityCalculator struct {
	mu  sync.RWMutex
	log logPkg.Logger

	totalCalculations uint64
	matchCount        uint64
}

// Calculate compares the registration city from user claims and database records.
func (c *RegCityCalculator) Calculate(ctx context.Context, user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	atomic.AddUint64(&c.totalCalculations, 1)

	select {
	case <-ctx.Done():
		return models.CalcResult{Status: "context_cancelled", Comment: ctx.Err().Error()}
	default:
	}

	c.mu.RLock()
	logger := c.log
	c.mu.RUnlock()

	if logger != nil {
		logger.Debugw("calculating city match", "u_city", user.UserClaim.RegCity, "db_city", db.RegCity)
	}

	res := checkLocationGeneric(ctx, logger, user.UserClaim.RegCity, db.RegCity, "City")

	if res.Status == errors.StatusMatch {
		atomic.AddUint64(&c.matchCount, 1)
	}

	return models.CalcResult{
		Name:    "Registration City Match",
		Code:    "reg_city",
		Value:   res.Value,
		Weight:  weights.RegCity,
		Result:  res.Value * weights.RegCity,
		Status:  res.Status,
		Comment: res.Comment,
	}
}

// checkLocationGeneric is a thread-safe helper function to compare string-based location data.
// It is used by both Country and City calculators to maintain logic consistency.
func checkLocationGeneric(ctx context.Context, log logPkg.Logger, userVal, dbVal, label string) rawCheckResult {
	u := strings.TrimSpace(userVal)
	d := strings.TrimSpace(dbVal)

	if u == "" || d == "" {
		if log != nil {
			log.Warnw(errors.ErrEmptyLocationData.Error(), "label", label, "is_user_empty", u == "", "is_db_empty", d == "")
		}
		return rawCheckResult{
			Value:   constants.NoMatch,
			Status:  errors.StatusNoData,
			Comment: label + " data is missing",
		}
	}

	if strings.EqualFold(u, d) {
		return rawCheckResult{
			Value:   constants.IdealMatch,
			Status:  errors.StatusMatch,
			Comment: label + " matches exactly",
		}
	}

	return rawCheckResult{
		Value:   constants.NoMatch,
		Status:  errors.StatusNoMatch,
		Comment: label + " mismatch",
	}
}

// GetStats returns the calculation metrics for the Country calculator.
func (c *RegCountryCalculator) GetStats() (uint64, uint64) {
	return atomic.LoadUint64(&c.totalCalculations), atomic.LoadUint64(&c.matchCount)
}

// GetStats returns the calculation metrics for the City calculator.
func (c *RegCityCalculator) GetStats() (uint64, uint64) {
	return atomic.LoadUint64(&c.totalCalculations), atomic.LoadUint64(&c.matchCount)
}

// SetLogger safely updates the logger for the Country calculator.
func (c *RegCountryCalculator) SetLogger(newLogger logPkg.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.log = newLogger
}

// SetLogger safely updates the logger for the City calculator.
func (c *RegCityCalculator) SetLogger(newLogger logPkg.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.log = newLogger
}