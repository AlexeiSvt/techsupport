package repo

import (
	"context"
	"database/sql"
	"fmt"
	"techsupport/storage/internal/db/queries"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) GetFullUserProfile(ctx context.Context, accTag string) (map[string]any, error) {
    if accTag == "" {
        return nil, fmt.Errorf("acc_tag is required")
    }

    s.mu.RLock()
    defer s.mu.RUnlock()

    var (
        tag, email     string
        isDonator      bool
        regDate        time.Time
        os, cpu, mID   sql.NullString 
    )

    err := s.Pool.QueryRow(ctx, queries.SelectFullUserProfileByTag, accTag).Scan(
        &tag,
        &email,
        &isDonator,
        &regDate,
        &os,
        &cpu,
        &mID,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, fmt.Errorf("user not found: %s", accTag)
        }
        return nil, fmt.Errorf("db error: %w", err)
    }

    return map[string]any{
        "acc_tag":     tag,
        "first_email": email,
        "is_donator":  isDonator,
        "reg_date":    regDate,
        "os":          os.String,
        "cpu_model":   cpu.String,
        "machine_id":  mID.String,
    }, nil
}