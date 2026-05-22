package web

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/opsagent/opsagent/internal/infrastructure/auth"
	"github.com/opsagent/opsagent/internal/config"
	"github.com/opsagent/opsagent/internal/infrastructure/executor"
	"github.com/opsagent/opsagent/internal/ui"
	"github.com/opsagent/opsagent/internal/infrastructure/storage/sqlite"
	"github.com/opsagent/opsagent/internal/usecase"
	"golang.org/x/crypto/bcrypt"
)

func NewRouter(cfg *config.Config, store *sqlite.Store) http.Handler {
	mux := http.NewServeMux()

	// Create JWT authenticator
	jwtAuth, err := auth.NewJWTAuth(cfg, store)
	if err != nil {
		panic("failed to create JWT auth: " + err.Error())
	}

	mux.HandleFunc("/api/system/service-status", methodGuard("GET", handleServiceStatus(cfg)))
	mux.HandleFunc("/health", methodGuard("GET", handleHealth()))
	mux.HandleFunc("/", handleRoot())
	mux.HandleFunc("/assets/", handleAssets())

	// Auth routes (JWT) - no middleware needed
	mux.HandleFunc("/api/setup/required", methodGuard("GET", handleSetupRequired(store)))
	mux.HandleFunc("/api/setup/admin", methodGuard("POST", handleCreateAdmin(store)))
	mux.HandleFunc("/api/auth/login", jwtAuth.LoginHandler())
	mux.HandleFunc("/api/auth/logout", jwtAuth.LogoutHandler())
	mux.HandleFunc("/api/auth/me", jwtAuth.MeHandler())
	mux.HandleFunc("/api/auth/refresh", jwtAuth.RefreshHandler())

	// Protected routes with JWT middleware
	jwtMw := jwtAuth.Middleware()

	// Alert routes
	mux.Handle("/api/alerts", jwtMw(methodGuard("GET", handleGetAlerts(store))))
	mux.Handle("/api/alerts/", jwtMw(handleAlertByID(store)))

	// Dashboard
	mux.Handle("/api/dashboard/summary", jwtMw(methodGuard("GET", handleDashboardSummary(store))))

	// Assistant and task routes
	mux.Handle("/api/assistant/message", jwtMw(methodGuard("POST", handleAssistantMessage(cfg, store))))
	mux.Handle("/api/tasks", jwtMw(methodGuard("GET", handleGetTasks(store))))
	mux.Handle("/api/tasks/", jwtMw(handleTasksByID(cfg, store)))

	// Settings and system routes
	mux.Handle("/api/settings", jwtMw(methodGuard("GET", handleGetSettings(store, cfg))))
	mux.Handle("/api/settings/ai", jwtMw(methodGuard("POST", handleUpdateAISettings(store, cfg))))
	mux.Handle("/api/settings/password", jwtMw(methodGuard("POST", handleUpdatePassword(store))))
	mux.Handle("/api/settings/access-mode", jwtMw(methodGuard("PATCH", handleUpdateAccessMode(store, cfg))))
	mux.Handle("/api/system/info", jwtMw(methodGuard("GET", handleSystemInfo(store))))
	mux.Handle("/api/system/update-check", jwtMw(methodGuard("GET", handleUpdateCheck(store))))
	mux.Handle("/api/system/update", jwtMw(methodGuard("POST", handleUpdateAgent(store))))

	return withSecurityHeaders(mux)
}

func methodGuard(expectedMethod string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != expectedMethod {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fn(w, r)
	}
}

func handleServiceStatus(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var output string
		var err error

		// Run systemctl status opsagent --no-pager
		runner := executor.NewRunner(10, 64)
		result, execErr := runner.Run(ctx, "systemctl status opsagent --no-pager")

		if execErr != nil {
			output = "Service status unavailable: " + execErr.Error()
		} else {
			output = result.Stdout
			if result.Stderr != "" {
				output = output + "\n[stderr]\n" + result.Stderr
			}
		}

		_ = err

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":  "ok",
			"version": "0.1.0",
			"output":  output,
		})
	}
}

func handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Serve React app for all non-API routes
		if r.URL.Path == "/" || !strings.HasPrefix(r.URL.Path, "/api") {
			index, err := ui.Files.ReadFile("dist/index.html")
			if err != nil {
				http.Error(w, "Dashboard not found. Run: cd web/dashboard && npm run build", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write(index)
			return
		}
		http.NotFound(w, r)
	}
}

func handleAssets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/assets/")
		if path == "" {
			http.NotFound(w, r)
			return
		}

		content, err := ui.Files.ReadFile("dist/assets/" + path)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		// Set content type based on file extension
		contentType := "application/octet-stream"
		if strings.HasSuffix(path, ".js") {
			contentType = "application/javascript"
		} else if strings.HasSuffix(path, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(path, ".html") {
			contentType = "text/html"
		} else if strings.HasSuffix(path, ".svg") {
			contentType = "image/svg+xml"
		} else if strings.HasSuffix(path, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") {
			contentType = "image/jpeg"
		} else if strings.HasSuffix(path, ".woff2") {
			contentType = "font/woff2"
		}

		w.Header().Set("Content-Type", contentType)
		w.Write(content)
	}
}

func withSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}

// Auth handlers

func handleSetupRequired(store *sqlite.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := store.GetUserByUsername(context.Background(), "admin")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"required": user == nil})
	}
}

func handleCreateAdmin(store *sqlite.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if len(req.Username) < 3 {
			http.Error(w, "username must be at least 3 characters", http.StatusBadRequest)
			return
		}

		if len(req.Password) < 8 {
			http.Error(w, "password must be at least 8 characters", http.StatusBadRequest)
			return
		}

		// Check if users exist
		user, _ := store.GetUserByUsername(context.Background(), req.Username)
		if user != nil {
			http.Error(w, "username already exists", http.StatusConflict)
			return
		}

		// Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "failed to hash password", http.StatusInternalServerError)
			return
		}

		// Create user
		err = store.CreateUser(context.Background(), req.Username, string(hash), "admin")
		if err != nil {
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}

		// Mark setup as completed
		_ = store.SetSetting(context.Background(), "setup_completed", "true")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

func handleLogin(store *sqlite.Store, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		user, err := store.GetUserByUsername(context.Background(), req.Username)
		if err != nil || user == nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		// Create session token
		tokenBytes := make([]byte, 32)
		_, _ = rand.Read(tokenBytes)
		token := base64.URLEncoding.EncodeToString(tokenBytes)

		// Store session
		_ = store.SetSetting(context.Background(), "session_"+token, user.Username)
		_ = store.SetSetting(context.Background(), "session_"+token+"_exp", time.Now().Add(24*time.Hour).Format(time.RFC3339))

		_ = store.InsertAuditLog(context.Background(), "login", user.Username, "User logged in", nil)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"token": token,
			"user": map[string]string{
				"username": user.Username,
				"role":     user.Role,
			},
		})
	}
}

func handleLogout(store *sqlite.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		token := r.Header.Get("Authorization")
		if token != "" {
			token = strings.TrimPrefix(token, "Bearer ")
			_ = store.DeleteSetting(context.Background(), "session_"+token)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

func handleMe(store *sqlite.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		username, _ := store.GetSetting(context.Background(), "session_"+token)
		expStr, _ := store.GetSetting(context.Background(), "session_"+token+"_exp")

		if username == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		exp, _ := time.Parse(time.RFC3339, expStr)
		if time.Now().After(exp) {
			_ = store.DeleteSetting(context.Background(), "session_"+token)
			http.Error(w, "session expired", http.StatusUnauthorized)
			return
		}

		user, _ := store.GetUserByUsername(context.Background(), username)
		if user == nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"username": user.Username,
			"role":     user.Role,
		})
	}
}

func authMiddleware(store *sqlite.Store, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		username, _ := store.GetSetting(context.Background(), "session_"+token)
		expStr, _ := store.GetSetting(context.Background(), "session_"+token+"_exp")

		if username == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		exp, _ := time.Parse(time.RFC3339, expStr)
		if time.Now().After(exp) {
			_ = store.DeleteSetting(context.Background(), "session_"+token)
			http.Error(w, "session expired", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func handleGetAlerts(store *sqlite.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 50
		offset := 0
		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
				limit = parsed
			}
		}
		if o := r.URL.Query().Get("offset"); o != "" {
			if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
				offset = parsed
			}
		}

		alerts, err := store.GetAlerts(r.Context(), limit, offset)
		if err != nil {
			http.Error(w, "failed to fetch alerts", http.StatusInternalServerError)
			return
		}

		count, _ := store.GetOpenAlertCount(r.Context())

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"alerts": alerts,
			"total":  len(alerts),
			"open":   count,
		})
	}
}

func handleAlertByID(store *sqlite.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract ID from path: /api/alerts/:id
		path := strings.TrimPrefix(r.URL.Path, "/api/alerts/")
		path = strings.TrimSuffix(path, "/")

		if path == "" {
			http.Error(w, "alert id required", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(path, 10, 64)
		if err != nil {
			http.Error(w, "invalid alert id", http.StatusBadRequest)
			return
		}

		// Handle sub-paths: acknowledge, resolve, ignore
		if r.URL.Query().Get("acknowledge") != "" || strings.HasSuffix(r.URL.Path, "/acknowledge") {
			if r.Method != "POST" {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			err := store.UpdateAlertStatus(r.Context(), id, "acknowledged")
			if err != nil {
				http.Error(w, "failed to update alert", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]bool{"success": true})
			return
		}

		if r.URL.Query().Get("resolve") != "" || strings.HasSuffix(r.URL.Path, "/resolve") {
			if r.Method != "POST" {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			err := store.UpdateAlertStatus(r.Context(), id, "resolved")
			if err != nil {
				http.Error(w, "failed to update alert", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]bool{"success": true})
			return
		}

		if r.URL.Query().Get("ignore") != "" || strings.HasSuffix(r.URL.Path, "/ignore") {
			if r.Method != "POST" {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			err := store.UpdateAlertStatus(r.Context(), id, "ignored")
			if err != nil {
				http.Error(w, "failed to update alert", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]bool{"success": true})
			return
		}

		// Get single alert
		alert, err := store.GetAlertByID(r.Context(), id)
		if err != nil {
			http.Error(w, "failed to fetch alert", http.StatusInternalServerError)
			return
		}
		if alert == nil {
			http.Error(w, "alert not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(alert)
	}
}

func handleDashboardSummary(store *sqlite.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		metric, err := store.GetLatestMetric(ctx)
		openAlerts, _ := store.GetOpenAlertCount(ctx)
		ports, _ := store.GetPorts(ctx)

		summary := map[string]any{
			"open_alerts":  openAlerts,
			"cpu_usage":    "-",
			"memory_usage": "-",
			"disk_usage":   "-",
			"uptime":       "-",
			"load":         "-",
			"open_ports":   "-",
		}

		if err == nil && metric != nil {
			summary["cpu_usage"] = fmt.Sprintf("%.0f", metric.CPUUsage)
			summary["memory_usage"] = fmt.Sprintf("%.0f", metric.MemoryUsage)
			summary["disk_usage"] = fmt.Sprintf("%.0f", metric.DiskUsage)
			summary["load"] = fmt.Sprintf("%.2f", metric.LoadAverage1)
			if metric.UptimeSeconds > 0 {
				days := metric.UptimeSeconds / 86400
				hours := (metric.UptimeSeconds % 86400) / 3600
				mins := (metric.UptimeSeconds % 3600) / 60
				if days > 0 {
					summary["uptime"] = fmt.Sprintf("%dd %dh %dm", days, hours, mins)
				} else if hours > 0 {
					summary["uptime"] = fmt.Sprintf("%dh %dm", hours, mins)
				} else {
					summary["uptime"] = fmt.Sprintf("%dm", mins)
				}
			}
		}

		if len(ports) > 0 {
			summary["open_ports"] = len(ports)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(summary)
	}
}

func handleAssistantMessage(cfg *config.Config, store *sqlite.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Message string `json:"message"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		req.Message = strings.TrimSpace(req.Message)
		if req.Message == "" {
			http.Error(w, "message is required", http.StatusBadRequest)
			return
		}

		result, err := usecase.CreateTaskPlanFromPrompt(r.Context(), cfg, store, req.Message)
		if err != nil {
			http.Error(w, "failed to create plan", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func handleGetTasks(store *sqlite.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 50
		offset := 0
		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
				limit = parsed
			}
		}
		if o := r.URL.Query().Get("offset"); o != "" {
			if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
				offset = parsed
			}
		}

		tasks, err := store.GetAllTasks(r.Context(), limit, offset)
		if err != nil {
			http.Error(w, "failed to fetch tasks", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"tasks": tasks,
			"total": len(tasks),
		})
	}
}

func handleTasksByID(cfg *config.Config, store *sqlite.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract task ID from path
		path := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
		path = strings.TrimSuffix(path, "/")

		if path == "" {
			http.Error(w, "task id required", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(path, 10, 64)
		if err != nil {
			http.Error(w, "invalid task id", http.StatusBadRequest)
			return
		}

		// Handle sub-paths
		if strings.HasSuffix(r.URL.Path, "/approve") {
			if r.Method != "POST" {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			err := store.UpdateTaskStatus(r.Context(), id, "approved")
			if err != nil {
				http.Error(w, "failed to approve task", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]bool{"success": true})
			return
		}

		if strings.HasSuffix(r.URL.Path, "/reject") {
			if r.Method != "POST" {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			err := store.UpdateTaskStatus(r.Context(), id, "rejected")
			if err != nil {
				http.Error(w, "failed to reject task", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]bool{"success": true})
			return
		}

		if strings.HasSuffix(r.URL.Path, "/run") {
			if r.Method != "POST" {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			runner := executor.NewRunner(cfg.Security.CommandTimeoutSeconds, cfg.Security.MaxOutputSizeKB)

			// Execute approved commands
			task, _ := store.GetTask(r.Context(), id)
			if task == nil {
				http.Error(w, "task not found", http.StatusNotFound)
				return
			}

			plans, _ := store.GetTaskPlan(r.Context(), id)
			if len(plans) == 0 {
				http.Error(w, "no plan found", http.StatusNotFound)
				return
			}

			plan := plans[0]
			cmds, _ := store.GetCommandsByPlanID(r.Context(), plan.ID)

			results := []map[string]any{}
			for _, cmd := range cmds {
				if cmd.Type == "blocked" {
					_ = store.UpdateCommandStatus(r.Context(), cmd.ID, "blocked")
					results = append(results, map[string]any{
						"command_id": cmd.ID,
						"command":    cmd.Command,
						"status":     "blocked",
					})
					continue
				}

				if cmd.RequiresApproval && task.Status != "approved" {
					results = append(results, map[string]any{
						"command_id": cmd.ID,
						"command":    cmd.Command,
						"status":     "skipped",
						"reason":     "requires approval",
					})
					continue
				}

				_ = store.UpdateCommandStatus(r.Context(), cmd.ID, "running")

				result, err := runner.Run(r.Context(), cmd.Command)

				if err != nil {
					_ = store.SaveCommandExecution(r.Context(), cmd.ID, "", err.Error(), -1, "failed")
					_ = store.UpdateCommandStatus(r.Context(), cmd.ID, "failed")
					_ = store.InsertAuditLog(r.Context(), "command_executed", "admin", fmt.Sprintf("Command failed: %s", cmd.Command), map[string]any{
						"command_id": cmd.ID,
						"error":      err.Error(),
					})
					results = append(results, map[string]any{
						"command_id": cmd.ID,
						"command":    cmd.Command,
						"status":     "failed",
						"error":      err.Error(),
					})
					continue
				}

				status := "success"
				if result.ExitCode != 0 {
					status = "failed"
				}

				_ = store.SaveCommandExecution(r.Context(), cmd.ID, result.Stdout, result.Stderr, result.ExitCode, status)
				_ = store.UpdateCommandStatus(r.Context(), cmd.ID, status)
				_ = store.InsertAuditLog(r.Context(), "command_executed", "admin", fmt.Sprintf("Command executed: %s", cmd.Command), map[string]any{
					"command_id":  cmd.ID,
					"exit_code":  result.ExitCode,
					"status":     status,
					"duration_ms": result.Duration.Milliseconds(),
				})

				output := result.Stdout
				if result.Stderr != "" {
					output = output + "\n[stderr]\n" + result.Stderr
				}

				results = append(results, map[string]any{
					"command_id":  cmd.ID,
					"command":     cmd.Command,
					"status":      status,
					"exit_code":   result.ExitCode,
					"output":      output,
					"duration_ms": result.Duration.Milliseconds(),
				})
			}

			_ = store.UpdateTaskStatus(r.Context(), id, "completed")

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"task_id": id,
				"results": results,
			})
			return
		}

		// Check for executions endpoint
		if strings.HasSuffix(r.URL.Path, "/executions") {
			idStr := strings.TrimSuffix(path, "/executions")
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				http.Error(w, "invalid task id", http.StatusBadRequest)
				return
			}

			plans, _ := store.GetTaskPlan(r.Context(), id)
			if len(plans) == 0 {
				http.Error(w, "no plan found", http.StatusNotFound)
				return
			}

			allExecutions := []map[string]any{}
			for _, plan := range plans {
				cmds, _ := store.GetCommandsByPlanID(r.Context(), plan.ID)
				for _, cmd := range cmds {
					execs, _ := store.GetCommandExecutions(r.Context(), cmd.ID)
					for _, e := range execs {
						allExecutions = append(allExecutions, map[string]any{
							"command_id":  e.CommandID,
							"command":     cmd.Command,
							"output":      e.Output,
							"error":       e.ErrorOutput,
							"exit_code":   e.ExitCode,
							"started_at":  e.StartedAt,
							"finished_at": e.FinishedAt,
							"status":      e.Status,
						})
					}
				}
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"executions": allExecutions})
			return
		}

		// Get task with plan and commands
		task, err := store.GetTask(r.Context(), id)
		if err != nil || task == nil {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}

		plans, _ := store.GetTaskPlan(r.Context(), id)

		taskData := map[string]any{
			"task":  task,
			"plans": []map[string]any{},
		}

		for _, plan := range plans {
			cmds, _ := store.GetCommandsByPlanID(r.Context(), plan.ID)
			planMap := map[string]any{
				"id":                plan.ID,
				"summary":           plan.Summary,
				"risk_level":        plan.RiskLevel,
				"requires_approval": plan.RequiresApproval,
				"commands":          cmds,
			}
			taskData["plans"] = append(taskData["plans"].([]map[string]any), planMap)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(taskData)
	}
}

func handleGetSettings(store *sqlite.Store, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		setupCompleted, _ := store.GetSetting(ctx, "setup_completed")

		// Read AI settings from database (user may have updated them)
		aiEnabled, _ := store.GetSetting(ctx, "ai_enabled")
		aiProvider, _ := store.GetSetting(ctx, "ai_provider")
		aiModel, _ := store.GetSetting(ctx, "ai_model")
		aiAPIKey, _ := store.GetSetting(ctx, "ai_api_key")

		settings := map[string]any{
			"setup_completed": setupCompleted == "true",
			"server": map[string]any{
				"bind_address": cfg.Server.BindAddress,
				"port":         cfg.Server.Port,
			},
			"ai": map[string]any{
				"enabled":  aiEnabled == "true",
				"provider": aiProvider,
				"model":   aiModel,
				"has_key":  aiAPIKey != "",
				"api_key":  "", // Never expose API key
			},
			"monitoring": map[string]any{
				"interval_seconds":           cfg.Monitoring.IntervalSeconds,
				"disk_warning_threshold":     cfg.Monitoring.DiskWarningThreshold,
				"disk_critical_threshold":   cfg.Monitoring.DiskCriticalThreshold,
				"memory_warning_threshold":   cfg.Monitoring.MemoryWarningThreshold,
				"cpu_warning_threshold":      cfg.Monitoring.CPUWarningThreshold,
				"cpu_warning_duration_secs":  cfg.Monitoring.CPUWarningDurationSecond,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	}
}

func handleUpdateAISettings(store *sqlite.Store, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Enabled  bool   `json:"enabled"`
			Provider string `json:"provider"`
			APIKey  string `json:"api_key"`
			Model   string `json:"model"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		// Store AI settings in database AND update in-memory config
		_ = store.SetSetting(r.Context(), "ai_provider", req.Provider)
		_ = store.SetSetting(r.Context(), "ai_model", req.Model)
		_ = store.SetSetting(r.Context(), "ai_enabled", strconv.FormatBool(req.Enabled))
		// Only update API key if provided (don't clear existing key)
		if req.APIKey != "" {
			_ = store.SetSetting(r.Context(), "ai_api_key", req.APIKey)
			cfg.AI.APIKey = req.APIKey
		}
		cfg.AI.Provider = req.Provider
		cfg.AI.Model = req.Model
		cfg.AI.Enabled = req.Enabled
		_ = store.InsertAuditLog(r.Context(), "settings_updated", "admin", "AI settings updated", map[string]any{
			"provider": req.Provider,
			"model":   req.Model,
			"enabled":  req.Enabled,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

func handleSystemInfo(store *sqlite.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		metric, _ := store.GetLatestMetric(ctx)

		info := map[string]any{
			"hostname":        "unknown",
			"os_info":        "unknown",
			"kernel_version": "unknown",
			"uptime_seconds": 0,
			"load_average_1": 0,
			"load_average_5": 0,
			"load_average_15": 0,
		}

		if metric != nil {
			info["hostname"] = metric.Hostname
			info["os_info"] = metric.OSInfo
			info["kernel_version"] = metric.KernelVersion
			info["uptime_seconds"] = metric.UptimeSeconds
			info["load_average_1"] = metric.LoadAverage1
			info["load_average_5"] = metric.LoadAverage5
			info["load_average_15"] = metric.LoadAverage15
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	}
}

func handleUpdateAccessMode(store *sqlite.Store, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			BindAddress string `json:"bind_address"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		// Validate bind_address value
		validBindAddresses := map[string]bool{"127.0.0.1": true, "0.0.0.0": true}
		if !validBindAddresses[req.BindAddress] {
			http.Error(w, "invalid bind_address: must be 127.0.0.1 or 0.0.0.0", http.StatusBadRequest)
			return
		}

		// Store in database for persistence across restarts
		err := store.SetSetting(r.Context(), "server_bind_address", req.BindAddress)
		if err != nil {
			http.Error(w, "failed to save setting", http.StatusInternalServerError)
			return
		}

		_ = store.InsertAuditLog(r.Context(), "settings_updated", "admin", "Access mode updated", map[string]any{
			"bind_address": req.BindAddress,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"message": "bind_address updated. Restart required for changes to take effect.",
		})
	}
}

func handleUpdateCheck(store *sqlite.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		runner := executor.NewRunner(30, 128)
		result, err := runner.Run(ctx, "curl -fsSL https://api.github.com/repos/ardakaraosmanoglu/opsAgent/releases/latest")

		var latestVersion string
		var downloadURL string

		if err == nil {
			// Try to parse tag_name from JSON response
			var release map[string]any
			if json.Unmarshal([]byte(result.Stdout), &release) == nil {
				if tag, ok := release["tag_name"].(string); ok {
					latestVersion = tag
				}
				if assets, ok := release["assets"].([]any); ok && len(assets) > 0 {
					for _, a := range assets {
						if asset, ok := a.(map[string]any); ok {
							if name, _ := asset["name"].(string); name == "opsagent-linux-amd64" {
								if url, _ := asset["browser_download_url"].(string); url != "" {
									downloadURL = url
								}
							}
						}
					}
				}
			}
		}

		currentVersion := "0.1.0"
		needsUpdate := false
		latestClean := strings.TrimPrefix(latestVersion, "v")
		currentClean := strings.TrimPrefix(currentVersion, "v")
		if latestVersion != "" && latestClean != currentClean {
			needsUpdate = true
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"current_version": currentVersion,
			"latest_version":  latestVersion,
			"needs_update":    needsUpdate,
			"download_url":    downloadURL,
		})
	}
}

func handleUpdateAgent(store *sqlite.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()

		// Get latest release info first
		runner := executor.NewRunner(30, 128)
		result, err := runner.Run(ctx, "curl -fsSL https://api.github.com/repos/ardakaraosmanoglu/opsAgent/releases/latest")

		var downloadURL string
		if err == nil {
			var release map[string]any
			if json.Unmarshal([]byte(result.Stdout), &release) == nil {
				if assets, ok := release["assets"].([]any); ok {
					for _, a := range assets {
						if asset, ok := a.(map[string]any); ok {
							if name, _ := asset["name"].(string); name == "opsagent-linux-amd64" {
								if url, _ := asset["browser_download_url"].(string); url != "" {
									downloadURL = url
								}
							}
						}
					}
				}
			}
		}

		if downloadURL == "" {
			http.Error(w, "failed to find download URL", http.StatusInternalServerError)
			return
		}

		// Download new binary to temp location
		_, err = runner.Run(ctx, fmt.Sprintf("curl -fsSL -L %s -o /tmp/opsagent-new && chmod +x /tmp/opsagent-new", downloadURL))
		if err != nil {
			http.Error(w, "failed to download update: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Stop service, replace binary, restart
		commands := []string{
			"systemctl stop opsagent",
			"mv /tmp/opsagent-new /usr/local/bin/opsagent",
			"chmod 755 /usr/local/bin/opsagent",
			"systemctl start opsagent",
		}

		for _, cmd := range commands {
			_, err = runner.Run(ctx, cmd)
			if err != nil {
				http.Error(w, "update failed at: "+cmd+" - "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		_ = store.InsertAuditLog(ctx, "agent_updated", "admin", "Agent updated to latest version", map[string]any{
			"download_url": downloadURL,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"message": "Agent updated successfully. Service restarted.",
		})
	}
}

func handleUpdatePassword(store *sqlite.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			CurrentPassword string `json:"current_password"`
			NewPassword     string `json:"new_password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if req.CurrentPassword == "" || req.NewPassword == "" {
			http.Error(w, "current_password and new_password are required", http.StatusBadRequest)
			return
		}

		if len(req.NewPassword) < 8 {
			http.Error(w, "new password must be at least 8 characters", http.StatusBadRequest)
			return
		}

		// Get current user from auth middleware context
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		username, _ := store.GetSetting(r.Context(), "session_"+token)
		if username == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := store.GetUserByUsername(r.Context(), username)
		if err != nil || user == nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		// Verify current password
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
			http.Error(w, "current password is incorrect", http.StatusUnauthorized)
			return
		}

		// Hash new password
		hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "failed to hash password", http.StatusInternalServerError)
			return
		}

		// Update password
		if err := store.UpdateUserPassword(r.Context(), user.ID, string(hash)); err != nil {
			http.Error(w, "failed to update password", http.StatusInternalServerError)
			return
		}

		_ = store.InsertAuditLog(r.Context(), "password_changed", user.Username, "User changed their password", nil)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}
