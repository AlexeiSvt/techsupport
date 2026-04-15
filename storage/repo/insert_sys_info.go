package repo

import (
	"context"
	"fmt"
	"techsupport/storage/internal/db/queries"
	sysModels "techsupport/sysinfo/pkg/models"
)

func (s *Storage) InsertSystemInfo(ctx context.Context, ticketID int64, sysInfo sysModels.SystemInfo) error {
	if sysInfo.MachineID == "" {
		return fmt.Errorf("machine_id is required for system info")
	}

	s.mu.Lock()
	defer s.mu.RUnlock()

	_, err := s.Pool.Exec(ctx, queries.InsertSysInfo,
		ticketID,
		sysInfo.OS,
		sysInfo.Platform,
		sysInfo.Arch,
		sysInfo.Kernel,
		sysInfo.CPUModel,
		sysInfo.CPUCores,
		sysInfo.TotalRAM,
		sysInfo.Hostname,
		sysInfo.MachineID,
		sysInfo.Username,
	)

	if err != nil {
		return fmt.Errorf("failed to insert system info for ticket %d: %w", ticketID, err)
	}

	return nil
}