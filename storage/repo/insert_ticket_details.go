package repo

import (
	"context"
	"fmt"
	coreModels "techsupport/core/pkg/models"
	"techsupport/storage/internal/db/queries"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) InsertTicketDetails(ctx context.Context, ticketID int64, details []coreModels.CalcResult) error {
	if len(details) == 0 {
		return nil 
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	batch := &pgx.Batch{}

	for _, d := range details {
		if d.Name == "" {
			continue
		}

		batch.Queue(queries.InsertTicketDetails,
			ticketID,
			d.Name,
			d.Code,
			d.Value,
			d.Weight,
			d.Result,
			d.Comment,
			d.Status,
		)
	}


	br := s.Pool.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("batch insert detail failed at index %d: %w", i, err)
		}
	}

	return nil
}