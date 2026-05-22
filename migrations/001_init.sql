-- OpsAgent Initial Schema

CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT 'admin',
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS settings (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL,
  updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS metrics (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  cpu_usage REAL,
  memory_usage REAL,
  disk_usage REAL,
  load_average_1 REAL,
  load_average_5 REAL,
  load_average_15 REAL,
  uptime_seconds INTEGER,
  created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS process_snapshots (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  pid INTEGER,
  name TEXT,
  user TEXT,
  cpu_usage REAL,
  memory_usage REAL,
  command TEXT,
  created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS ports (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  protocol TEXT,
  local_address TEXT,
  port INTEGER,
  process_name TEXT,
  pid INTEGER,
  first_seen_at DATETIME NOT NULL,
  last_seen_at DATETIME NOT NULL,
  status TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS services (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  status TEXT,
  active_state TEXT,
  sub_state TEXT,
  first_seen_at DATETIME NOT NULL,
  last_seen_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS alerts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  type TEXT NOT NULL,
  severity TEXT NOT NULL,
  title TEXT NOT NULL,
  message TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'open',
  metadata_json TEXT,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  source TEXT NOT NULL,
  prompt TEXT,
  status TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS task_plans (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  task_id INTEGER NOT NULL,
  summary TEXT NOT NULL,
  risk_level TEXT NOT NULL,
  requires_approval INTEGER NOT NULL DEFAULT 1,
  metadata_json TEXT,
  created_at DATETIME NOT NULL,
  FOREIGN KEY (task_id) REFERENCES tasks(id)
);

CREATE TABLE IF NOT EXISTS commands (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  task_plan_id INTEGER NOT NULL,
  command TEXT NOT NULL,
  command_type TEXT NOT NULL,
  risk_level TEXT NOT NULL,
  requires_approval INTEGER NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending',
  created_at DATETIME NOT NULL,
  FOREIGN KEY (task_plan_id) REFERENCES task_plans(id)
);

CREATE TABLE IF NOT EXISTS command_executions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  command_id INTEGER NOT NULL,
  output TEXT,
  error_output TEXT,
  exit_code INTEGER,
  started_at DATETIME NOT NULL,
  finished_at DATETIME,
  status TEXT NOT NULL,
  FOREIGN KEY (command_id) REFERENCES commands(id)
);

CREATE TABLE IF NOT EXISTS audit_logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  event_type TEXT NOT NULL,
  actor TEXT,
  message TEXT NOT NULL,
  metadata_json TEXT,
  created_at DATETIME NOT NULL
);
