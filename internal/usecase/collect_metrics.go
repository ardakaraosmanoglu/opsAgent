package usecase

import (
    "context"

    "github.com/opsagent/opsagent/internal/infrastructure/storage/sqlite"
    "github.com/opsagent/opsagent/internal/infrastructure/system"
)

func CollectMetrics(ctx context.Context, store *sqlite.Store) error {
    m, err := system.CollectMetrics()
    if err != nil {
        return err
    }

    return store.InsertMetric(ctx, m)
}
