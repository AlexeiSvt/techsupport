package repo

import (
	"context"
	"database/sql"
	"fmt"
	coreModels "techsupport/core/pkg/models"
	"techsupport/storage/internal/db/queries"
	"time"
)

func (s *Storage) GetSessionHistory(ctx context.Context, accTag string) ([]coreModels.Session, error) {
    if accTag == "" {
        return nil, fmt.Errorf("acc_tag is required")
    }

    s.mu.RLock()
    defer s.mu.RUnlock()

    rows, err := s.Pool.Query(ctx, queries.SelectUserSessionHistory, accTag)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch sessions: %w", err)
    }
    defer rows.Close()

    var history []coreModels.Session

    for rows.Next() {
        var sess coreModels.Session
        var startTime time.Time 
        var endTime sql.NullTime

        err := rows.Scan(
            &sess.SessionID,
            &sess.SessionIP,
            &sess.DeviceID,
            &sess.ASN,
            &sess.Country,
            &sess.City,
            &startTime, 
            &endTime,   
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan session row: %w", err)
        }

        sess.StartTime = startTime.Format(time.RFC3339)

        if endTime.Valid {
            sess.EndTime = endTime.Time.Format(time.RFC3339)
        } else {
            sess.EndTime = "active"
        }

        history = append(history, sess)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error during rows iteration: %w", err)
    }

    return history, nil
}