// Package pkg defines the fundamental interfaces and abstractions for the scoring system.
// These contracts allow for decoupled development, making it easy to add new 
// calculation rules or swap engine implementations without breaking the system.
package pkg

import (
	"context"
	"techsupport/core/pkg/models"
)

// ScoreCalculator defines the standard behavior for an individual scoring rule.
// Any structure that implements this interface can be injected into the 
// main scoring pipeline.
type ScoreCalculator interface {
	// Calculate performs a specific validation logic (e.g., IP check, Age check)
	// and returns a detailed result based on the provided user and database records.
	Calculate(ctx context.Context, user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult
}

// Engine acts as the high-level orchestrator for the scoring process.
// It is responsible for gathering data, executing calculators, and 
// synthesizing the final output.
type Engine interface {
	// Run takes the raw input data and executes the full scoring lifecycle,
	// returning a finalized OutputData structure ready for the caller.
	Run(input models.InputData) models.OutputData
}