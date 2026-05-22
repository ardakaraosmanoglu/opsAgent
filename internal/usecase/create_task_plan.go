package usecase

import (
	"context"
	"strings"

	"github.com/opsagent/opsagent/internal/config"
	"github.com/opsagent/opsagent/internal/infrastructure/ai"
	"github.com/opsagent/opsagent/internal/infrastructure/storage/sqlite"
	"github.com/opsagent/opsagent/internal/policy"
)

func CreateTaskPlanFromPrompt(ctx context.Context, cfg *config.Config, store *sqlite.Store, prompt string) (map[string]any, error) {
	taskID, err := store.CreateTask(ctx, "assistant", prompt, "draft")
	if err != nil {
		return nil, err
	}

	var plan *ai.Plan
	var requiresApproval bool

	// Try AI if configured
	if cfg.AI.Enabled && cfg.AI.APIKey != "" {
		aiClient := &ai.Client{
			APIKey:  cfg.AI.APIKey,
			BaseURL: cfg.AI.Provider,
			Model:   cfg.AI.Model,
		}

		// Get latest metrics for context
		metric, _ := store.GetLatestMetric(ctx)

		systemContext := map[string]any{}
		if metric != nil {
			systemContext["cpu_usage"] = metric.CPUUsage
			systemContext["memory_usage"] = metric.MemoryUsage
			systemContext["disk_usage"] = metric.DiskUsage
		}

		plan, err = aiClient.CreatePlan(ctx, prompt, systemContext)
		if err != nil {
			// Fall back to template
			plan = nil
		}
	}

	// Use template plan if AI not available
	if plan == nil {
		template := buildFallbackPlan(prompt)
		plan = &ai.Plan{
			Summary:          template.Summary,
			RiskLevel:        template.RiskLevel,
			RequiresApproval: false,
			Commands:         nil,
		}
		for _, cmd := range template.Commands {
			plan.Commands = append(plan.Commands, ai.PlanCommand{
				Command:        cmd.Command,
				Purpose:        "",
				ExpectedEffect: "",
				RiskLevel:      cmd.RiskLevel,
			})
		}
	}

	requiresApproval = false
	for i := range plan.Commands {
		result := policy.Classify(plan.Commands[i].Command)

		if !result.Allowed {
			plan.Commands[i].Command = "BLOCKED:" + plan.Commands[i].Command
			plan.Commands[i].RiskLevel = "critical"
			requiresApproval = false
			continue
		}

		if result.RequiresApproval {
			requiresApproval = true
		}
	}

	status := "completed"
	if requiresApproval {
		status = "waiting_approval"
	}

	riskLevel := plan.RiskLevel
	if riskLevel == "" {
		riskLevel = "low"
	}

	planID, err := store.CreateTaskPlan(ctx, taskID, plan.Summary, riskLevel, requiresApproval)
	if err != nil {
		return nil, err
	}

	for _, cmd := range plan.Commands {
		cmdType := "read"
		risk := "low"
		requiresApprovalCmd := false

		if cmd.Command != "" && len(cmd.Command) > 8 && cmd.Command[:8] == "BLOCKED:" {
			cmdType = "blocked"
			risk = "critical"
		} else {
			result := policy.Classify(cmd.Command)
			cmdType = string(result.Type)
			risk = string(result.RiskLevel)
			requiresApprovalCmd = result.RequiresApproval
		}

		err := store.CreateCommand(ctx, planID, sqlite.FallbackCommand{
			Command:          cmd.Command,
			CommandType:      cmdType,
			RiskLevel:        risk,
			RequiresApproval: requiresApprovalCmd,
		})
		if err != nil {
			return nil, err
		}
	}

	if err := store.UpdateTaskStatus(ctx, taskID, status); err != nil {
		return nil, err
	}

	_ = store.InsertAuditLog(ctx, "plan_generated", "assistant", prompt, map[string]any{
		"task_id": taskID,
		"plan_id": planID,
		"status":  status,
	})

	responseType := "answer"
	if status == "waiting_approval" {
		responseType = "plan"
	}

	return map[string]any{
		"type":              responseType,
		"task_id":           taskID,
		"summary":           plan.Summary,
		"requires_approval": requiresApproval,
	}, nil
}

type fallbackPlan struct {
	Summary   string
	RiskLevel string
	Commands  []sqlite.FallbackCommand
}

func buildFallbackPlan(prompt string) fallbackPlan {
	if contains(prompt, "disk") {
		return fallbackPlan{
			Summary:   "Analyzing disk usage",
			RiskLevel: "low",
			Commands: []sqlite.FallbackCommand{
				{Command: "df -h", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
				{Command: "du -sh /var/log/* 2>/dev/null | sort -rh | head -10", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
			},
		}
	}

	return fallbackPlan{
		Summary:   "System information",
		RiskLevel: "low",
		Commands: []sqlite.FallbackCommand{
			{Command: "uptime", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
			{Command: "free -m", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
		},
	}
}

func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
