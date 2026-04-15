// Package engine orchestrates the execution of multiple scoring rules.
package engine

import (
	"context"
	"techsupport/core/pkg"
	"techsupport/core/pkg/models"
	logPkg "techsupport/log/pkg"
)

// CalculateScore executes a set of scoring calculators against user and database records.
// It aggregates the results into a single numerical score.
// This function is designed to be the primary entry point for the scoring logic.
func CalculateScore(ctx context.Context, log logPkg.Logger, user models.UserData, db models.DBRecord, weights models.Weights, calculators []pkg.ScoreCalculator) float64 {
	total := 0.0

	if log != nil {
		log.Debugw("starting fast score calculation",
			"calculators_count", len(calculators),
			"user_id", user.UserClaim.Phone, // Using phone as a trace identifier
		)
	}

	// Iterate through each provided calculator.
	// Since calculators are standardized via the pkg.ScoreCalculator interface,
	// we can process them uniformly.
	for _, calc := range calculators {
		// Check if the context was cancelled before starting the next calculation.
		select {
		case <-ctx.Done():
			if log != nil {
				log.Warnw("score calculation interrupted by context", "error", ctx.Err())
			}
			return total
		default:
		}

		// Perform the specific calculation rule.
		res := calc.Calculate(ctx, user, db, weights)
		
		// Accumulate the weighted result.
		total += res.Result
	}

	if log != nil {
		log.Infow("fast score calculation finished",
			"total_score", total,
		)
	}

	return total
}