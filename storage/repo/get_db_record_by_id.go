package repo

import (
	"context"
	"fmt"
	"techsupport/storage/internal/db/queries"
	coreModels "techsupport/core/pkg/models"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) GetDBRecordByTag(ctx context.Context, accTag string) (coreModels.DBRecord, error) {
    if accTag == "" {
        return coreModels.DBRecord{}, fmt.Errorf("acc_tag is empty")
    }

    s.mu.RLock()
    defer s.mu.RUnlock()

    var rec coreModels.DBRecord

    err := s.Pool.QueryRow(ctx, queries.SelectDBRecordByTag, accTag).Scan(
        &rec.AccTag,
        &rec.RegCountry,
        &rec.RegCity,
        &rec.FirstEmail,
        &rec.Phone,
        &rec.FirstDevice,
        &rec.IsDonator,
        &rec.RegDate,
        &rec.Devices,          
        &rec.FirstTransaction, 
        &rec.UserHistory,  
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return coreModels.DBRecord{}, fmt.Errorf("record not found for tag: %s", accTag)
        }
        return coreModels.DBRecord{}, fmt.Errorf("failed to fetch db_record: %w", err)
    }

    return rec, nil
}