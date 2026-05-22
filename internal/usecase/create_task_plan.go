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

	// Classify commands and separate read-only from write operations
	var readOnlyCommands []string
	var requiresApproval bool
	var hasBlocked bool

	for i := range plan.Commands {
		result := policy.Classify(plan.Commands[i].Command)

		if !result.Allowed {
			plan.Commands[i].Command = "BLOCKED:" + plan.Commands[i].Command
			plan.Commands[i].RiskLevel = "critical"
			hasBlocked = true
			continue
		}

		if result.RequiresApproval {
			requiresApproval = true
		} else {
			readOnlyCommands = append(readOnlyCommands, plan.Commands[i].Command)
		}
	}

	// If has blocked commands, don't allow execution
	if hasBlocked {
		requiresApproval = true
	}

	riskLevel := plan.RiskLevel
	if riskLevel == "" {
		riskLevel = "low"
	}

	planID, err := store.CreateTaskPlan(ctx, taskID, plan.Summary, riskLevel, requiresApproval)
	if err != nil {
		return nil, err
	}

	// Store commands in DB
	for _, cmd := range plan.Commands {
		cmdType := "read"
		risk := "low"
		requiresApprovalCmd := false

		if cmd.Command != "" && strings.HasPrefix(cmd.Command, "BLOCKED:") {
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

	// For read-only queries (no approval needed), execute commands and return answer
	if !requiresApproval && len(readOnlyCommands) > 0 {
		// Execute read-only commands
		output, execErr := ExecuteReadOnlyCommands(ctx, cfg, readOnlyCommands)
		if execErr != nil {
			output = "Command execution error: " + execErr.Error()
		}

		// Generate natural language answer
		answer, _ := GenerateAnswer(ctx, cfg, prompt, output)

		status := "completed"
		if err := store.UpdateTaskStatus(ctx, taskID, status); err != nil {
			return nil, err
		}

		_ = store.InsertAuditLog(ctx, "prompt_submitted", "assistant", prompt, map[string]any{
			"task_id":    taskID,
			"plan_id":    planID,
			"type":       "read_only",
			"answer":     answer,
		})

		return map[string]any{
			"type":              "answer",
			"task_id":           taskID,
			"summary":           answer,
			"requires_approval": false,
		}, nil
	}

	// For write operations requiring approval, return plan
	status := "waiting_approval"
	if err := store.UpdateTaskStatus(ctx, taskID, status); err != nil {
		return nil, err
	}

	_ = store.InsertAuditLog(ctx, "plan_generated", "assistant", prompt, map[string]any{
		"task_id": taskID,
		"plan_id": planID,
		"status":  status,
	})

	return map[string]any{
		"type":              "plan",
		"task_id":           taskID,
		"summary":           plan.Summary,
		"requires_approval": true,
	}, nil
}

type fallbackPlan struct {
	Summary   string
	RiskLevel string
	Commands  []sqlite.FallbackCommand
}

func buildFallbackPlan(prompt string) fallbackPlan {
	lower := strings.ToLower(prompt)

	// Disk analysis
	if strings.Contains(lower, "disk") || strings.Contains(lower, "alan") || strings.Contains(lower, "doluyor") {
		return fallbackPlan{
			Summary:   "Analyzing disk usage",
			RiskLevel: "low",
			Commands: []sqlite.FallbackCommand{
				{Command: "df -h", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
				{Command: "du -sh /var/log/* 2>/dev/null | sort -rh | head -10", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
				{Command: "df -i 2>/dev/null | sort -k5 -rh | head -10", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
			},
		}
	}

	// Memory analysis
	if strings.Contains(lower, "memory") || strings.Contains(lower, "ram") || strings.Contains(lower, "bellek") {
		return fallbackPlan{
			Summary:   "Analyzing memory usage",
			RiskLevel: "low",
			Commands: []sqlite.FallbackCommand{
				{Command: "free -m", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
				{Command: "ps aux --sort=-%mem | head -10", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
			},
		}
	}

	// CPU analysis
	if strings.Contains(lower, "cpu") || strings.Contains(lower, "işlemci") || strings.Contains(lower, "yüksek") {
		return fallbackPlan{
			Summary:   "Analyzing CPU usage",
			RiskLevel: "low",
			Commands: []sqlite.FallbackCommand{
				{Command: "top -bn1 | head -20", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
				{Command: "uptime", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
			},
		}
	}

	// Port analysis
	if strings.Contains(lower, "port") || strings.Contains(lower, "açık") || strings.Contains(lower, "listening") {
		return fallbackPlan{
			Summary:   "Analyzing open ports",
			RiskLevel: "low",
			Commands: []sqlite.FallbackCommand{
				{Command: "ss -tulpn 2>/dev/null || netstat -tulpn 2>/dev/null", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
			},
		}
	}

	// Service analysis
	if strings.Contains(lower, "service") || strings.Contains(lower, "servis") || strings.Contains(lower, "nginx") || strings.Contains(lower, "docker") {
		return fallbackPlan{
			Summary:   "Analyzing services",
			RiskLevel: "low",
			Commands: []sqlite.FallbackCommand{
				{Command: "systemctl list-units --type=service --state=running | head -30", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
				{Command: "docker ps 2>/dev/null || echo 'Docker not available'", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
			},
		}
	}

	// Log analysis
	if strings.Contains(lower, "log") || strings.Contains(lower, "journal") || strings.Contains(lower, "日志") {
		return fallbackPlan{
			Summary:   "Analyzing system logs",
			RiskLevel: "low",
			Commands: []sqlite.FallbackCommand{
				{Command: "journalctl --disk-usage", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
				{Command: "du -sh /var/log/* 2>/dev/null | sort -rh | head -10", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
				{Command: "journalctl -n 20 --no-pager 2>/dev/null", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
			},
		}
	}

	// Package update check
	if strings.Contains(lower, "update") || strings.Contains(lower, "upgrade") || strings.Contains(lower, "güncelle") || strings.Contains(lower, "paket") {
		return fallbackPlan{
			Summary:   "Checking available updates",
			RiskLevel: "medium",
			Commands: []sqlite.FallbackCommand{
				{Command: "apt list --upgradable 2>/dev/null || yum list updates 2>/dev/null || echo 'No package manager available'", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
			},
		}
	}

	// Default: general system info
	return fallbackPlan{
		Summary:   "System information",
		RiskLevel: "low",
		Commands: []sqlite.FallbackCommand{
			{Command: "uptime", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
			{Command: "free -m", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
			{Command: "df -h", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
			{Command: "ss -tulpn 2>/dev/null | wc -l", CommandType: "read", RiskLevel: "low", RequiresApproval: false},
		},
	}
}
