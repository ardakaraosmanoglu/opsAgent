package app

import (
    "context"
    "log"
    "time"

    "github.com/opsagent/opsagent/internal/config"
    "github.com/opsagent/opsagent/internal/infrastructure/storage/sqlite"
    "github.com/opsagent/opsagent/internal/usecase"
)

func StartScheduler(ctx context.Context, cfg *config.Config, store *sqlite.Store) error {
    metricsTicker := time.NewTicker(time.Duration(cfg.Monitoring.IntervalSeconds) * time.Second)
    deepScanTicker := time.NewTicker(5 * time.Minute)
    cleanupTicker := time.NewTicker(1 * time.Hour)

    defer metricsTicker.Stop()
    defer deepScanTicker.Stop()
    defer cleanupTicker.Stop()

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-metricsTicker.C:
            if err := usecase.CollectMetrics(ctx, store); err != nil {
                log.Printf("collect metrics failed: %v", err)
            }
            if err := usecase.EvaluateAlerts(ctx, cfg, store); err != nil {
                log.Printf("evaluate alerts failed: %v", err)
            }
        case <-deepScanTicker.C:
            if err := usecase.RunDeepScan(ctx, store); err != nil {
                log.Printf("deep scan failed: %v", err)
            }
        case <-cleanupTicker.C:
            if err := usecase.CleanupOldMetrics(ctx, store); err != nil {
                log.Printf("cleanup old metrics failed: %v", err)
            }
        }
    }
}
