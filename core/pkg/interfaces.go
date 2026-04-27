// Package pkg defines the fundamental interfaces and abstractions for the scoring system.
package pkg

import (
    "context"
    "techsupport/core/pkg/models"
)

// ScoreCalculator defines the standard behavior for an individual scoring rule.
// It now uses the updated flat models designed for the graph architecture.
type ScoreCalculator interface {
    // Calculate performs validation logic using the user's claim, support context, 
    // and historical database records.
    Calculate(
        ctx context.Context, 
        claim models.UserClaim, 
        support models.SupportContext, 
        db models.DBRecord, 
        weights models.Weights,
    ) models.CalcResult
}

// Engine acts as the high-level orchestrator for the scoring process.
type Engine interface {
    // Run executes the scoring lifecycle. 
    // Note: If you renamed InputData to something else, update the parameter here.
    // Usually, it takes the combined context of the claim and the support interaction.
    Run(ctx context.Context, claim models.UserClaim, support models.SupportContext) models.OutputData
}