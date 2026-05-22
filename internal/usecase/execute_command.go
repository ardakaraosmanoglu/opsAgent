package usecase

import (
    "context"

    "github.com/opsagent/opsagent/internal/config"
    "github.com/opsagent/opsagent/internal/infrastructure/executor"
    "github.com/opsagent/opsagent/internal/infrastructure/storage/sqlite"
    "github.com/opsagent/opsagent/internal/policy"
)

func ExecuteApprovedCommand(ctx context.Context, cfg *config.Config, store *sqlite.Store, commandID int64) error {
    cmd, err := store.GetCommandByID(ctx, commandID)
    if err != nil {
        return err
    }

    policyResult := policy.Classify(cmd.Command)
    if !policyResult.Allowed {
        _ = store.MarkCommandBlocked(ctx, commandID, policyResult.Reason)
        _ = store.InsertAuditLog(ctx, "command_blocked", "system", policyResult.Reason, nil)
        return err
    }

    if policyResult.RequiresApproval && cmd.Status != "approved" {
        return err
    }

    runner := executor.NewRunner(
        cfg.Security.CommandTimeoutSeconds,
        cfg.Security.MaxOutputSizeKB,
    )

    _ = store.MarkCommandRunning(ctx, commandID)

    result, err := runner.Run(ctx, cmd.Command)
    if err != nil {
        _ = store.SaveCommandExecution(ctx, commandID, "", err.Error(), -1, "failed")
        _ = store.MarkCommandFailed(ctx, commandID)
        return err
    }

    status := "success"
    if result.ExitCode != 0 {
        status = "failed"
    }

    if err := store.SaveCommandExecution(ctx, commandID, result.Stdout, result.Stderr, result.ExitCode, status); err != nil {
        return err
    }

    if status == "success" {
        _ = store.MarkCommandSuccess(ctx, commandID)
    } else {
        _ = store.MarkCommandFailed(ctx, commandID)
    }

    _ = store.InsertAuditLog(ctx, "command_executed", "admin", "Command executed", map[string]any{
        "command_id": commandID,
        "exit_code":  result.ExitCode,
        "status":     status,
    })

    return nil
}
