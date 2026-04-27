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
    
    // CreatedAt marks the initiation of the scoring engine.
    CreatedAt time.Time `json:"created_at" cypher:"created_at"`
    
    // UpdatedAt marks the completion of the scoring process.
    UpdatedAt time.Time `json:"updated_at" cypher:"updated_at"`
    
    // FinalPercentage is the formatted reliability score (e.g., "95.00%").
    FinalPercentage string `json:"final_percentage" cypher:"final_percentage"`

    // KnowledgeSum is the flat cumulative trust score.
    KnowledgeSum float64 `json:"knowledge_sum" cypher:"knowledge_sum"`

    // PenaltySum is the flat cumulative risk score.
    PenaltySum float64 `json:"penalty_sum" cypher:"penalty_sum"`
}