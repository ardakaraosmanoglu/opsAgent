package usecase

import (
	"context"

	"github.com/opsagent/opsagent/internal/infrastructure/storage/sqlite"
	"github.com/opsagent/opsagent/internal/infrastructure/system"
)

func RunDeepScan(ctx context.Context, store *sqlite.Store) error {
	scan, err := system.RunDeepScan()
	if err != nil {
		return err
	}

	if len(scan.TopCPUProcesses) > 0 {
		if err := store.InsertProcessSnapshot(ctx, scan.TopCPUProcesses); err != nil {
			return err
		}
	}

	for _, port := range scan.OpenPorts {
		if err := store.UpsertPort(ctx, port); err != nil {
			return err
		}
	}

	for _, svc := range scan.Services {
		if err := store.UpsertService(ctx, svc); err != nil {
			return err
		}
	}

	return nil
}

func CleanupOldMetrics(ctx context.Context, store *sqlite.Store) error {
	return store.CleanupOldMetrics(ctx)
}
