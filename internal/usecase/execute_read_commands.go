package usecase

import (
	"context"
	"strings"

	"github.com/opsagent/opsagent/internal/config"
	"github.com/opsagent/opsagent/internal/infrastructure/ai"
	"github.com/opsagent/opsagent/internal/infrastructure/executor"
)

// ExecuteReadOnlyCommands runs read-only commands and returns aggregated output
func ExecuteReadOnlyCommands(ctx context.Context, cfg *config.Config, commands []string) (string, error) {
	if len(commands) == 0 {
		return "", nil
	}

	runner := executor.NewRunner(cfg.Security.CommandTimeoutSeconds, cfg.Security.MaxOutputSizeKB)

	var output strings.Builder

	for _, cmd := range commands {
		if cmd == "" {
			continue
		}

		result, err := runner.Run(ctx, cmd)
		if err != nil {
			output.WriteString("Command: " + cmd + "\nError: " + err.Error() + "\n\n")
			continue
		}

		output.WriteString("=== " + cmd + " ===\n")
		if result.Stdout != "" {
			output.WriteString(result.Stdout)
		}
		if result.Stderr != "" {
			output.WriteString("\n[stderr] " + result.Stderr)
		}
		output.WriteString("\n\n")
	}

	return strings.TrimSpace(output.String()), nil
}

// GenerateAnswer creates a natural language response from command outputs
func GenerateAnswer(ctx context.Context, cfg *config.Config, prompt string, commandOutputs string) (string, error) {
	// If AI is enabled, use it to generate analysis
	if cfg.AI.Enabled && cfg.AI.APIKey != "" {
		aiClient := &ai.Client{
			APIKey:  cfg.AI.APIKey,
			BaseURL: cfg.AI.Provider,
			Model:   cfg.AI.Model,
		}

		analysisPrompt := "You are OpsAgent assistant. Analyze the command outputs and provide a clear, concise answer in Turkish if the user asked in Turkish, otherwise in English.\n\n"
		analysisPrompt += "User question: " + prompt + "\n\n"
		analysisPrompt += "Command outputs:\n" + commandOutputs + "\n\n"
		analysisPrompt += "Provide a helpful response that:\n"
		analysisPrompt += "1. Directly answers the user's question\n"
		analysisPrompt += "2. Highlights any concerning values (high disk, low memory, etc.)\n"
		analysisPrompt += "3. Suggests actions if there are problems\n"
		analysisPrompt += "4. Is concise but informative"

		plan, err := aiClient.CreatePlan(ctx, analysisPrompt, nil)
		if err == nil && plan != nil && plan.Summary != "" {
			return plan.Summary, nil
		}
	}

	// Fallback: parse output directly based on prompt
	return parseCommandOutputForAnswer(prompt, commandOutputs), nil
}

// parseCommandOutputForAnswer generates a response from command outputs without AI
func parseCommandOutputForAnswer(prompt string, output string) string {
	lower := strings.ToLower(prompt)
	var answer strings.Builder

	// Disk analysis
	if strings.Contains(lower, "disk") || strings.Contains(lower, "alan") || strings.Contains(lower, "doluyor") {
		answer.WriteString("Disk Durumu:\n")
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Use%") || strings.Contains(line, " использование") ||
				(strings.Count(line, "%") > 0 && strings.Contains(line, "G")) {
				answer.WriteString(line + "\n")
			}
		}
		answer.WriteString("\n")
	}

	// Memory analysis
	if strings.Contains(lower, "memory") || strings.Contains(lower, "ram") || strings.Contains(lower, "bellek") {
		answer.WriteString("Bellek Durumu:\n")
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Mem:") || strings.Contains(line, "total") {
				answer.WriteString(line + "\n")
			}
		}
		answer.WriteString("\n")
	}

	// Default: return truncated output
	if answer.Len() == 0 {
		lines := strings.Split(output, "\n")
		maxLines := 30
		if len(lines) < maxLines {
			maxLines = len(lines)
		}
		for i := 0; i < maxLines; i++ {
			answer.WriteString(lines[i] + "\n")
		}
		if len(lines) > maxLines {
			answer.WriteString("\n... (" + string(rune('0'+len(lines)-maxLines)) + " more lines)")
		}
	}

	return strings.TrimSpace(answer.String())
}
