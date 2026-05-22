package sqlite

import (
    "context"
    "database/sql"
    "encoding/json"
    "errors"
    "fmt"

    "github.com/opsagent/opsagent/internal/domain"
    "github.com/opsagent/opsagent/internal/infrastructure/system"
    _ "modernc.org/sqlite"
)

type Store struct {
    DB *sql.DB
}

func Open(path string) (*Store, error) {
    db, err := sql.Open("sqlite", path)
    if err != nil {
        return nil, fmt.Errorf("opening sqlite: %w", err)
    }

    db.SetMaxOpenConns(1)

    if _, err := db.Exec(`PRAGMA journal_mode=WAL`); err != nil {
        return nil, fmt.Errorf("setting WAL mode: %w", err)
    }

    if _, err := db.Exec(`PRAGMA foreign_keys=ON`); err != nil {
        return nil, fmt.Errorf("enabling foreign keys: %w", err)
    }

    return &Store{DB: db}, nil
}

func (s *Store) Migrate() error {
    _, err := s.DB.Exec(initSQL)
    return err
}

// AlertInput is used for creating/updating alerts
type AlertInput struct {
    Type     string
    Severity string
    Title    string
    Message  string
}

// Repository methods needed by usecases

func (s *Store) InsertMetric(ctx context.Context, m *system.Metrics) error {
    _, err := s.DB.ExecContext(ctx, `
        INSERT INTO metrics (hostname, os_info, kernel_version, cpu_usage, memory_usage, disk_usage, load_average_1, load_average_5, load_average_15, uptime_seconds, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))`,
        m.Hostname, m.OSInfo, m.KernelVersion, m.CPUUsage, m.MemoryUsage, m.DiskUsage, m.LoadAverage1, m.LoadAverage5, m.LoadAverage15, m.UptimeSeconds)
    return err
}

func (s *Store) GetLatestMetric(ctx context.Context) (*domain.Metric, error) {
    row := s.DB.QueryRowContext(ctx, `
        SELECT id, hostname, os_info, kernel_version, cpu_usage, memory_usage, disk_usage, load_average_1, load_average_5, load_average_15, uptime_seconds, created_at
        FROM metrics ORDER BY created_at DESC LIMIT 1`)

    var m domain.Metric
    err := row.Scan(&m.ID, &m.Hostname, &m.OSInfo, &m.KernelVersion, &m.CPUUsage, &m.MemoryUsage, &m.DiskUsage, &m.LoadAverage1, &m.LoadAverage5, &m.LoadAverage15, &m.UptimeSeconds, &m.CreatedAt)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &m, nil
}

func (s *Store) UpsertOpenAlert(ctx context.Context, input AlertInput) error {
    // Check if there's already an open alert of this type
    row := s.DB.QueryRowContext(ctx, `
        SELECT id FROM alerts WHERE type = ? AND status = 'open' LIMIT 1`, input.Type)
    var existingID int64
    err := row.Scan(&existingID)

    if errors.Is(err, sql.ErrNoRows) {
        // Create new alert
        _, err = s.DB.ExecContext(ctx, `
            INSERT INTO alerts (type, severity, title, message, status, created_at, updated_at)
            VALUES (?, ?, ?, ?, 'open', datetime('now'), datetime('now'))`,
            input.Type, input.Severity, input.Title, input.Message)
        return err
    }

    if err != nil {
        return err
    }

    // Update existing alert
    _, err = s.DB.ExecContext(ctx, `
        UPDATE alerts SET severity = ?, title = ?, message = ?, updated_at = datetime('now')
        WHERE id = ?`, input.Severity, input.Title, input.Message, existingID)
    return err
}

func (s *Store) GetAlerts(ctx context.Context, limit, offset int) ([]domain.Alert, error) {
    rows, err := s.DB.QueryContext(ctx, `
        SELECT id, type, severity, title, message, status, metadata_json, created_at, updated_at
        FROM alerts ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var alerts []domain.Alert
    for rows.Next() {
        var a domain.Alert
        if err := rows.Scan(&a.ID, &a.Type, &a.Severity, &a.Title, &a.Message, &a.Status, &a.MetadataJSON, &a.CreatedAt, &a.UpdatedAt); err != nil {
            return nil, err
        }
        alerts = append(alerts, a)
    }
    return alerts, rows.Err()
}

func (s *Store) GetAlertByID(ctx context.Context, id int64) (*domain.Alert, error) {
    row := s.DB.QueryRowContext(ctx, `
        SELECT id, type, severity, title, message, status, metadata_json, created_at, updated_at
        FROM alerts WHERE id = ?`, id)
    var a domain.Alert
    err := row.Scan(&a.ID, &a.Type, &a.Severity, &a.Title, &a.Message, &a.Status, &a.MetadataJSON, &a.CreatedAt, &a.UpdatedAt)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &a, nil
}

func (s *Store) UpdateAlertStatus(ctx context.Context, id int64, status string) error {
    _, err := s.DB.ExecContext(ctx, `
        UPDATE alerts SET status = ?, updated_at = datetime('now') WHERE id = ?`, status, id)
    return err
}

func (s *Store) GetOpenAlertCount(ctx context.Context) (int, error) {
    row := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM alerts WHERE status = 'open'`)
    var count int
    err := row.Scan(&count)
    return count, err
}

func (s *Store) GetAlertCountByType(ctx context.Context, alertType string) (int, error) {
    row := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM alerts WHERE type = ? AND status = 'open'`, alertType)
    var count int
    err := row.Scan(&count)
    return count, err
}

func (s *Store) CreateTask(ctx context.Context, source, prompt, status string) (int64, error) {
    result, err := s.DB.ExecContext(ctx, `
        INSERT INTO tasks (source, prompt, status, created_at, updated_at)
        VALUES (?, ?, ?, datetime('now'), datetime('now'))`,
        source, prompt, status)
    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}

func (s *Store) CreateTaskPlan(ctx context.Context, taskID int64, summary, riskLevel string, requiresApproval bool) (int64, error) {
    result, err := s.DB.ExecContext(ctx, `
        INSERT INTO task_plans (task_id, summary, risk_level, requires_approval, created_at)
        VALUES (?, ?, ?, ?, datetime('now'))`,
        taskID, summary, riskLevel, requiresApproval)
    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}

// FallbackCommand is exported for use by usecase layer
type FallbackCommand struct {
    Command          string
    CommandType      string
    RiskLevel        string
    RequiresApproval bool
}

func (s *Store) CreateCommand(ctx context.Context, planID int64, cmd FallbackCommand) error {
    _, err := s.DB.ExecContext(ctx, `
        INSERT INTO commands (task_plan_id, command, command_type, risk_level, requires_approval, status, created_at)
        VALUES (?, ?, ?, ?, ?, 'pending', datetime('now'))`,
        planID, cmd.Command, cmd.CommandType, cmd.RiskLevel, cmd.RequiresApproval)
    return err
}

func (s *Store) UpdateTaskStatus(ctx context.Context, taskID int64, status string) error {
    _, err := s.DB.ExecContext(ctx, `
        UPDATE tasks SET status = ?, updated_at = datetime('now') WHERE id = ?`,
        status, taskID)
    return err
}

func (s *Store) InsertAuditLog(ctx context.Context, eventType, actor, message string, meta map[string]any) error {
    var metaJSON []byte
    if meta != nil {
        metaJSON, _ = json.Marshal(meta)
    }
    _, err := s.DB.ExecContext(ctx, `
        INSERT INTO audit_logs (event_type, actor, message, metadata_json, created_at)
        VALUES (?, ?, ?, ?, datetime('now'))`,
        eventType, actor, message, string(metaJSON))
    return err
}

func (s *Store) GetCommandByID(ctx context.Context, commandID int64) (*domain.Command, error) {
    row := s.DB.QueryRowContext(ctx, `
        SELECT id, task_plan_id, command, command_type, risk_level, requires_approval, status
        FROM commands WHERE id = ?`, commandID)

    var cmd domain.Command
    err := row.Scan(&cmd.ID, &cmd.TaskPlanID, &cmd.Command, &cmd.Type, &cmd.RiskLevel, &cmd.RequiresApproval, &cmd.Status)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, errors.New("command not found")
    }
    if err != nil {
        return nil, err
    }
    return &cmd, nil
}

func (s *Store) MarkCommandRunning(ctx context.Context, commandID int64) error {
    _, err := s.DB.ExecContext(ctx, `UPDATE commands SET status = 'running' WHERE id = ?`, commandID)
    return err
}

func (s *Store) MarkCommandBlocked(ctx context.Context, commandID int64, reason string) error {
    _, err := s.DB.ExecContext(ctx, `UPDATE commands SET status = 'blocked' WHERE id = ?`, commandID)
    return err
}

func (s *Store) MarkCommandSuccess(ctx context.Context, commandID int64) error {
    _, err := s.DB.ExecContext(ctx, `UPDATE commands SET status = 'success' WHERE id = ?`, commandID)
    return err
}

func (s *Store) MarkCommandFailed(ctx context.Context, commandID int64) error {
    _, err := s.DB.ExecContext(ctx, `UPDATE commands SET status = 'failed' WHERE id = ?`, commandID)
    return err
}

func (s *Store) SaveCommandExecution(ctx context.Context, commandID int64, stdout, stderr string, exitCode int, status string) error {
    _, err := s.DB.ExecContext(ctx, `
        INSERT INTO command_executions (command_id, output, error_output, exit_code, started_at, finished_at, status)
        VALUES (?, ?, ?, ?, datetime('now'), datetime('now'), ?)`,
        commandID, stdout, stderr, exitCode, status)
    return err
}

func (s *Store) InsertProcessSnapshot(ctx context.Context, procs []system.ProcessInfo) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, p := range procs {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO process_snapshots (pid, name, user, cpu_usage, memory_usage, command, created_at)
			VALUES (?, ?, ?, ?, ?, ?, datetime('now'))`,
			p.PID, p.Name, p.User, p.CPUUsage, p.MemUsage, p.Command)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) UpsertPort(ctx context.Context, port system.PortInfo) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO ports (protocol, local_address, port, process_name, pid, first_seen_at, last_seen_at, status)
		VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'), ?)
		ON CONFLICT(port, protocol) DO UPDATE SET
			last_seen_at = datetime('now'),
			process_name = excluded.process_name,
			pid = excluded.pid`,
		port.Protocol, port.LocalAddr, port.Port, port.ProcessName, port.PID, port.State)
	return err
}

func (s *Store) UpsertService(ctx context.Context, svc system.ServiceInfo) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO services (name, status, active_state, sub_state, first_seen_at, last_seen_at)
		VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
		ON CONFLICT(name) DO UPDATE SET
			last_seen_at = datetime('now'),
			status = excluded.status,
			active_state = excluded.active_state,
			sub_state = excluded.sub_state`,
		svc.Name, svc.Status, svc.ActiveState, svc.SubState)
	return err
}

func (s *Store) GetPorts(ctx context.Context) ([]system.PortInfo, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT protocol, local_address, port, process_name, pid, status FROM ports ORDER BY port`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ports []system.PortInfo
	for rows.Next() {
		var p system.PortInfo
		if err := rows.Scan(&p.Protocol, &p.LocalAddr, &p.Port, &p.ProcessName, &p.PID, &p.State); err != nil {
			return nil, err
		}
		ports = append(ports, p)
	}
	return ports, rows.Err()
}

func (s *Store) GetServices(ctx context.Context) ([]system.ServiceInfo, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT name, status, active_state, sub_state FROM services ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []system.ServiceInfo
	for rows.Next() {
		var svc system.ServiceInfo
		if err := rows.Scan(&svc.Name, &svc.Status, &svc.ActiveState, &svc.SubState); err != nil {
			return nil, err
		}
		services = append(services, svc)
	}
	return services, rows.Err()
}

func (s *Store) GetTopCPUProcesses(ctx context.Context, limit int) ([]system.ProcessInfo, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT pid, name, user, cpu_usage, memory_usage, command
		FROM process_snapshots
		ORDER BY created_at DESC, cpu_usage DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var procs []system.ProcessInfo
	for rows.Next() {
		var p system.ProcessInfo
		if err := rows.Scan(&p.PID, &p.Name, &p.User, &p.CPUUsage, &p.MemUsage, &p.Command); err != nil {
			return nil, err
		}
		procs = append(procs, p)
	}
	return procs, rows.Err()
}

func (s *Store) GetTopMemoryProcesses(ctx context.Context, limit int) ([]system.ProcessInfo, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT pid, name, user, cpu_usage, memory_usage, command
		FROM process_snapshots
		ORDER BY created_at DESC, memory_usage DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var procs []system.ProcessInfo
	for rows.Next() {
		var p system.ProcessInfo
		if err := rows.Scan(&p.PID, &p.Name, &p.User, &p.CPUUsage, &p.MemUsage, &p.Command); err != nil {
			return nil, err
		}
		procs = append(procs, p)
	}
	return procs, rows.Err()
}

func (s *Store) CleanupOldMetrics(ctx context.Context) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM metrics WHERE created_at < datetime('now', '-7 days')`)
	return err
}

// User management methods

func (s *Store) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
    row := s.DB.QueryRowContext(ctx, `SELECT id, username, password_hash, role, created_at, updated_at FROM users WHERE username = ?`, username)
    var u domain.User
    err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &u, nil
}

func (s *Store) CreateUser(ctx context.Context, username, passwordHash, role string) error {
    _, err := s.DB.ExecContext(ctx, `
        INSERT INTO users (username, password_hash, role, created_at, updated_at)
        VALUES (?, ?, ?, datetime('now'), datetime('now'))`,
        username, passwordHash, role)
    return err
}

func (s *Store) GetSetting(ctx context.Context, key string) (string, error) {
    row := s.DB.QueryRowContext(ctx, `SELECT value FROM settings WHERE key = ?`, key)
    var val string
    err := row.Scan(&val)
    if errors.Is(err, sql.ErrNoRows) {
        return "", nil
    }
    return val, err
}

func (s *Store) SetSetting(ctx context.Context, key, value string) error {
    _, err := s.DB.ExecContext(ctx, `
        INSERT INTO settings (key, value, updated_at) VALUES (?, ?, datetime('now'))
        ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at`,
        key, value)
    return err
}

func (s *Store) DeleteSetting(ctx context.Context, key string) error {
    _, err := s.DB.ExecContext(ctx, `DELETE FROM settings WHERE key = ?`, key)
    return err
}

func (s *Store) GetTask(ctx context.Context, taskID int64) (*domain.Task, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, source, prompt, status, created_at, updated_at FROM tasks WHERE id = ?`, taskID)
	var t domain.Task
	err := row.Scan(&t.ID, &t.Source, &t.Prompt, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *Store) GetAllTasks(ctx context.Context, limit, offset int) ([]domain.Task, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, source, prompt, status, created_at, updated_at
		FROM tasks ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(&t.ID, &t.Source, &t.Prompt, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (s *Store) GetTaskPlan(ctx context.Context, taskID int64) ([]domain.TaskPlan, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, task_id, summary, risk_level, requires_approval, metadata_json, created_at
		FROM task_plans WHERE task_id = ?`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []domain.TaskPlan
	for rows.Next() {
		var p domain.TaskPlan
		if err := rows.Scan(&p.ID, &p.TaskID, &p.Summary, &p.RiskLevel, &p.RequiresApproval, &p.MetadataJSON, &p.CreatedAt); err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, rows.Err()
}

func (s *Store) GetCommandsByPlanID(ctx context.Context, planID int64) ([]domain.Command, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, task_plan_id, command, command_type, risk_level, requires_approval, status
		FROM commands WHERE task_plan_id = ?`, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cmds []domain.Command
	for rows.Next() {
		var c domain.Command
		if err := rows.Scan(&c.ID, &c.TaskPlanID, &c.Command, &c.Type, &c.RiskLevel, &c.RequiresApproval, &c.Status); err != nil {
			return nil, err
		}
		cmds = append(cmds, c)
	}
	return cmds, rows.Err()
}

func (s *Store) UpdateCommandStatus(ctx context.Context, commandID int64, status string) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE commands SET status = ? WHERE id = ?`, status, commandID)
	return err
}

func (s *Store) GetCommandExecutions(ctx context.Context, commandID int64) ([]CommandExecution, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, command_id, output, error_output, exit_code, started_at, finished_at, status
		FROM command_executions WHERE command_id = ? ORDER BY started_at DESC`, commandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var execs []CommandExecution
	for rows.Next() {
		var e CommandExecution
		if err := rows.Scan(&e.ID, &e.CommandID, &e.Output, &e.ErrorOutput, &e.ExitCode, &e.StartedAt, &e.FinishedAt, &e.Status); err != nil {
			return nil, err
		}
		execs = append(execs, e)
	}
	return execs, rows.Err()
}

type CommandExecution struct {
	ID          int64
	CommandID   int64
	Output      string
	ErrorOutput string
	ExitCode    int
	StartedAt   string
	FinishedAt  string
	Status      string
}
