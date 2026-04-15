package models

import (
	coreModels "techsupport/core/pkg/models"
	sysModels "techsupport/sysinfo/pkg/models"
	"time"
)

type TicketRecord struct {
	TicketID        int64   `db:"ticket_id"`
	AccTag          string  `db:"acc_tag"`
	ClaimantTag     string  `db:"claimant_tag"`
	DeviceID        string  `db:"device_id"`
	FinalPercentage string  `db:"final_percentage"`
	Knowledge       float64 `db:"knowledge_score"`
	Penalty         float64 `db:"penalty_score"`

	UserData    coreModels.UserData     `db:"user_data"`
	UserHistory coreModels.UserHistory  `db:"user_history"`
	SysInfo     sysModels.SystemInfo    `db:"sys_info"`
	Details     []coreModels.CalcResult `db:"details"`
}

type TicketAgentView struct {
    TicketID    int64     `json:"ticket_id"`
    SubmittedAt time.Time `json:"submitted_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    MachineID   string    `json:"machine_id"`
    Decision    string    `json:"decision"`
}
