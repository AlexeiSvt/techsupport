package repo

import (
	"context"
	"database/sql"
	"fmt"
	"techsupport/storage/internal/db/queries"
    stModels "techsupport/storage/pkg/models"

)


func (s *Storage) GetTicketsForAgent(ctx context.Context, accTag string) ([]stModels.TicketAgentView, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    rows, err := s.Pool.Query(ctx, queries.SelectTicketDecisionView, accTag)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch agent view: %w", err)
    }
    defer rows.Close()

    var views []stModels.TicketAgentView
    for rows.Next() {
        var v stModels.TicketAgentView
        var mID sql.NullString 

        err := rows.Scan(
            &v.TicketID,
            &v.SubmittedAt,
            &v.UpdatedAt,
            &mID,
            &v.Decision,
        )
        if err != nil {
            return nil, fmt.Errorf("scan agent view failed: %w", err)
        }
        
        v.MachineID = mID.String 
        views = append(views, v)
    }

    return views, nil
}