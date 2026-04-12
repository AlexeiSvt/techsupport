package models

import (
	"time"
)

type InputData struct {
	UserData UserData `json:"user_data"`
	DBRecord DBRecord `json:"db_record"`
}

type OutputData struct {
	TicketID       int64    `json:"ticket_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	FinaPercentage string    `json:"final_percentage"`
	Metrics        Metrics   `json:"metrics"`
	Details []CalcResult `json:"details,omitempty"`
}

type Metrics struct {
	Knowledge    float64 `json:"knowledge"`
	PenaltyScore float64 `json:"penalty_score"`
}
