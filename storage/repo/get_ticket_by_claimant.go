package repo

import (
	"context"
	"fmt"
	"techsupport/storage/internal/db/queries"
	stModels "techsupport/storage/pkg/models"
	coreModels "techsupport/core/pkg/models"
)

func (s *Storage) GetTicketsByClaimant(ctx context.Context, claimantTag string) ([]stModels.TicketRecord, error) {
    if claimantTag == "" {
        return nil, fmt.Errorf("validation error: claimant_tag is empty")
    }

    s.mu.RLock()
    defer s.mu.RUnlock()

    rows, err := s.Pool.Query(ctx, queries.SelectTicketsByClaimant, claimantTag)
    if err != nil {
        return nil, fmt.Errorf("failed to select tickets by claimant: %w", err)
    }
    defer rows.Close()

    var tickets []stModels.TicketRecord

    for rows.Next() {
        var t stModels.TicketRecord
        if err := rows.Scan(&t.TicketID, &t.AccTag, &t.FinalPercentage); err != nil {
            return nil, fmt.Errorf("scan ticket failed: %w", err)
        }
        t.ClaimantTag = claimantTag
        tickets = append(tickets, t)
    }

    for i := range tickets {
        err := s.Pool.QueryRow(ctx, queries.SelectSysInfoByTicket, tickets[i].TicketID).Scan(
            &tickets[i].SysInfo.OS,
            &tickets[i].SysInfo.CPUModel,
            &tickets[i].SysInfo.MachineID,
        )
        if err != nil {
            continue
        }

        dRows, err := s.Pool.Query(ctx, queries.SelectDetailsByTicket, tickets[i].TicketID)
        if err != nil {
            continue
        }
        
        for dRows.Next() {
            var d coreModels.CalcResult
            if err := dRows.Scan(&d.Name, &d.Code, &d.Result, &d.Status); err == nil {
                tickets[i].Details = append(tickets[i].Details, d)
            }
        }
        dRows.Close()
    }

    return tickets, nil
}