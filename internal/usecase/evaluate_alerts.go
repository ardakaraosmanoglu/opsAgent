package usecase

import (
    "context"
    "fmt"

    "github.com/opsagent/opsagent/internal/config"
    "github.com/opsagent/opsagent/internal/infrastructure/storage/sqlite"
)

func EvaluateAlerts(ctx context.Context, cfg *config.Config, store *sqlite.Store) error {
    metric, err := store.GetLatestMetric(ctx)
    if err != nil {
        return err
    }

    if metric == nil {
        return nil
    }

    // Disk critical
    if metric.DiskUsage >= cfg.Monitoring.DiskCriticalThreshold {
        err = store.UpsertOpenAlert(ctx, sqlite.AlertInput{
            Type:     "disk_usage",
            Severity: "critical",
            Title:    "Disk usage is critical",
            Message:  fmt.Sprintf("Disk usage reached %.1f%% — action required immediately", metric.DiskUsage),
        })
        if err != nil {
            return err
        }
        return nil
    }

    // Disk warning
    if metric.DiskUsage >= cfg.Monitoring.DiskWarningThreshold {
        err = store.UpsertOpenAlert(ctx, sqlite.AlertInput{
            Type:     "disk_usage",
            Severity: "warning",
            Title:    "Disk usage is high",
            Message:  fmt.Sprintf("Disk usage reached %.1f%%", metric.DiskUsage),
        })
        if err != nil {
            return err
        }
    }

    // Memory warning
    if metric.MemoryUsage >= cfg.Monitoring.MemoryWarningThreshold {
        err = store.UpsertOpenAlert(ctx, sqlite.AlertInput{
            Type:     "memory_usage",
            Severity: "warning",
            Title:    "Memory usage is high",
            Message:  fmt.Sprintf("Memory usage reached %.1f%%", metric.MemoryUsage),
        })
        if err != nil {
            return err
        }
    }

    return nil
}
