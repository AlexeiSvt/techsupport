// Package pkg defines the fundamental interfaces and abstractions for the scoring system.
package pkg

import (
    "context"
    "techsupport/core/pkg/models"
)

// ScoreCalculator defines the standard behavior for an individual scoring rule.
// It now uses the updated flat models designed for the graph architecture.
type ScoreCalculator interface {
// Calculate performs the specific scoring logic for a given user and database record.
// It takes a context for cancellation, a UserNode representing the user's claims and interactions,
// a DBRecordInfoNode representing the verified historical state of the user's account, and Weights for scoring adjustments.
// The method returns a CalcResultInfoNode, which includes the outcome of the calculation and any relevant comments or metadata.
    Calculate(ctx context.Context, user *models.UserNode, db models.DBRecordInfoNode, weights models.Weights) models.CalcResultInfoNode
}


// Engine acts as the high-level orchestrator for the scoring process.
// It coordinates the execution of multiple ScoreCalculators, aggregates their results,
// and produces the final OutputData. The Engine is designed to be flexible and extensible,
// allowing for easy integration of new scoring rules and logic without modifying the core orchestration.
type Engine interface {

// Run executes the scoring process for a given user and support context, returning the final OutputData.
// It takes a context for cancellation, a UserNode representing the user's claims and interactions,
// and a SupportContextNodeInfo providing additional information about the support ticket and environment. 
// The method is designed to be the primary entry point for the scoring logic, 
// coordinating multiple ScoreCalculators and aggregating their results into 
// a comprehensive OutputData structure.
    Run(ctx context.Context, user *models.UserNode, support models.SupportContextNodeInfo) models.OutputData
}