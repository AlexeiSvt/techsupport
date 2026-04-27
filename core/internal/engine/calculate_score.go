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
func CalculateScore(
	ctx context.Context, 
	log logPkg.Logger, 
	claim models.UserClaim, 
	support models.SupportContext, 
	db models.DBRecord, 
	weights models.Weights, 
	calculators []pkg.ScoreCalculator,
) models.OutputData {
	// Initialize the output structure with metadata.
	output := models.OutputData{
		UBTicketID: claim.UBTicketID,
		Results:    make([]models.CalcResult, 0, len(calculators)),
	}

	if log != nil {
		log.Debugw("starting fast score calculation",
			"calculators_count", len(calculators),
			"ub_ticket_id", claim.UBTicketID, // Using UB nuance as a trace identifier
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
			return output
		default:
		}

		if calc == nil {
			if log != nil {
				log.Fatalw("CRITICAL: Calculator is nil. Check your engine initialization!")
			}
			panic("nil calculator")
		}

		// Perform the specific calculation rule using the new flat interface.
		res := calc.Calculate(ctx, claim, support, db, weights)

		// Aggregate Knowledge and Penalty sums based on the calculation result.
		// Positive values contribute to Knowledge, negative to Penalties (anomalies).
		if res.Value < 0 {
			output.PenaltySum += res.Result
		} else {
			output.KnowledgeSum += res.Result
		}

		// Store individual result for the audit trail.
		output.Results = append(output.Results, res)
	}

	// Calculate the final combined score.
	output.FinalScore = output.KnowledgeSum + output.PenaltySum

	if log != nil {
		log.Infow("fast score calculation finished",
			"ub_ticket_id", output.UBTicketID,
			"knowledge_sum", output.KnowledgeSum,
			"penalty_sum", output.PenaltySum,
			"final_score", output.FinalScore,
		)
	}

	return output
}