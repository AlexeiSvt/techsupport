// Package models contains the data structures used throughout the scoring engine.
package models

import (
    "time"
)

// OutputData represents the final verdict and summary of the scoring process.
// In a graph database, detailed metrics and individual calculations are 
// moved to separate nodes linked to this report via relationships.
type OutputData struct {
    // TicketID is the unique identifier mapping this result to a support ticket.
    TicketID int64 `json:"ticket_id" cypher:"ticket_id"`

    // UBTicketID is the internal Bot-Session identifier used for tracing.
    UBTicketID string `json:"ub_ticket_id" cypher:"ub_ticket_id"`

    // CreatedAt marks the initiation of the scoring engine.
    CreatedAt time.Time `json:"created_at" cypher:"created_at"`
    
    // UpdatedAt marks the completion of the scoring process.
    UpdatedAt time.Time `json:"updated_at" cypher:"updated_at"`
    
    // FinalPercentage is the formatted reliability score (e.g., "95.00%").
    FinalPercentage string `json:"final_percentage" cypher:"final_percentage"`

    // FinalScore is the raw numerical sum of all weighted results.
    FinalScore float64 `json:"final_score" cypher:"final_score"`

    // KnowledgeSum is the flat cumulative trust score.
    KnowledgeSum float64 `json:"knowledge_sum" cypher:"knowledge_sum"`

    // PenaltySum is the flat cumulative risk score.
    PenaltySum float64 `json:"penalty_sum" cypher:"penalty_sum"`

    // Results holds the detailed breakdown of each individual calculator's output.
    // This is used for auditing and debugging before being persisted to the graph.
    Results []CalcResult `json:"results" cypher:"-"`
}