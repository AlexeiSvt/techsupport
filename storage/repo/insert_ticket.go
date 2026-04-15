package repo

import (
	"context"
	"fmt"
	"techsupport/storage/internal/db/queries"
	stModels "techsupport/storage/pkg/models"
)

func (s *Storage) InsertTicket(ctx context.Context, ticket stModels.TicketRecord) (int64, error) {
	if ticket.AccTag == "" || ticket.ClaimantTag == "" {
		return 0, fmt.Errorf("validation error: acc_tag and claimant_tag are required")
	}
	if ticket.SysInfo.MachineID == "" {
		return 0, fmt.Errorf("validation error: machine_id is missing")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var ticketID int64

	err = tx.QueryRow(ctx, queries.InsertTicketAndReturnID,
		ticket.AccTag,
		ticket.ClaimantTag,
		ticket.DeviceID,
		ticket.FinalPercentage,
		ticket.Knowledge,
		ticket.Penalty,
		ticket.UserData,
	).Scan(&ticketID)

	if err != nil {
		return 0, fmt.Errorf("insert ticket failed: %w", err)
	}

	for _, d := range ticket.Details {
		if d.Weight <= 0 {
			continue
		}
		_, err = tx.Exec(ctx, queries.InsertTicketDetails,
			ticketID, d.Name, d.Code, d.Value, d.Weight, d.Result, d.Comment, d.Status,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to insert detail %s: %w", d.Name, err)
		}
	}

	_, err = tx.Exec(ctx, queries.InsertSysInfo,
		ticketID,
		ticket.SysInfo.OS,
		ticket.SysInfo.Platform,
		ticket.SysInfo.Arch,
		ticket.SysInfo.Kernel,
		ticket.SysInfo.CPUModel,
		ticket.SysInfo.CPUCores,
		ticket.SysInfo.TotalRAM,
		ticket.SysInfo.Hostname,
		ticket.SysInfo.MachineID,
		ticket.SysInfo.Username,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert system info: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit: %w", err)
	}

	return ticketID, nil
}
