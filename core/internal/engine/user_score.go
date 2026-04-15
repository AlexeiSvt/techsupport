// Package engine orchestrates the creation and execution of scoring rules.
package engine

import (
	"context"

	"techsupport/core/internal/logic"
	"techsupport/core/internal/logic/transactions"
	"techsupport/core/pkg"
	"techsupport/core/pkg/models"
	logPkg "techsupport/log/pkg"
)

// NewScoringEngine acts as a Factory, initializing all specialized calculators.
// It sets up dependencies and ensures all rules follow a unified logging policy.
func NewScoringEngine(log logPkg.Logger) []pkg.ScoreCalculator {
	// Initialize calculators with private fields using SetLogger.
	regCountry := &logic.RegCountryCalculator{}
	regCountry.SetLogger(log)

	regCity := &logic.RegCityCalculator{}
	regCity.SetLogger(log)

	firstEmail := &logic.FirstEmailCalculator{}
	firstEmail.SetLogger(log)

	firstDevice := &logic.FirstDeviceCalculator{}
	firstDevice.SetLogger(log)

	devicesCalc := &logic.DevicesCalculator{}
	devicesCalc.SetLogger(log)

	// This constructor only wants the logger, as per your package definition.
	firstTx := transactions.NewFirstTransactionScoreCalculator(log)

	// Initializing security-focused calculators.
	bruteForce := NewDeviceBruteforcePenaltyCalculator(log)

	return []pkg.ScoreCalculator{
		logic.NewRegDateCalculator(log),
		regCountry,
		regCity,
		firstEmail,
		logic.NewFirstPhoneCalculator(log),
		firstDevice,
		devicesCalc,
		firstTx,
		bruteForce,
	}
}

// CalculateAll executes the full suite of scoring rules and returns granular results.
// It is designed for "Heavy Scoring" where a full audit trail (details) is required.
func CalculateAll(ctx context.Context, log logPkg.Logger, user models.UserData, db models.DBRecord, weights models.Weights, calculators []pkg.ScoreCalculator) ([]models.CalcResult, float64) {
	var details []models.CalcResult
	var total float64

	if log != nil {
		log.Infow("starting full scoring session", 
			"is_donator", weights.FirstTransaction > 0,
			"calculators_loaded", len(calculators),
		)
	}

	// Iterate through all registered rules.
	for _, calc := range calculators {
		// Respect context cancellation between calculation steps.
		select {
		case <-ctx.Done():
			if log != nil {
				log.Warnw("scoring session aborted by context", "error", ctx.Err())
			}
			return details, total
		default:
		}

		// Execute the specific rule logic.
		res := calc.Calculate(ctx, user, db, weights)
		
		// Accumulate results and collect metadata for the audit report.
		total += res.Result 
		details = append(details, res)
	}

	if log != nil {
		log.Infow("scoring session completed", 
			"total_score", total, 
			"results_count", len(details),
		)
	}

	return details, total
}