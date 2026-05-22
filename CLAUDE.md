# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
# Build Go binary
go build -o opsagent ./cmd/opsagent

# Run server
./opsagent serve --config /etc/opsagent/config.yaml

# Frontend (web/dashboard/)
cd web/dashboard && pnpm install && pnpm dev

# Production frontend build (outputs to web/dashboard/dist/, then copied to internal/ui/dist for embed)
cd web/dashboard && pnpm build
```

## Architecture

### Backend (Go)
- **Entry**: `cmd/opsagent/main.go`
- **Server**: `internal/web/server.go` - HTTP router, API handlers, serves embedded React
- **Use cases**: `internal/usecase/` - task planning, command execution
- **Domain**: `internal/domain/` - User, Task, Command, Metric, Alert, Audit models
- **Infrastructure**: `internal/infrastructure/`
  - `auth/` - JWT authentication
  - `storage/sqlite/` - SQLite store with migrations
  - `executor/` - Command runner with timeout/output limits
  - `ai/` - AI client for plan generation
  - `system/` - Metrics collection

### Frontend (React + TypeScript + Vite)
- **Location**: `web/dashboard/`
- **Routing**: React Router DOM 6
- **Styling**: Tailwind CSS 3
- **Components**: Radix UI primitives
- **API client**: `web/dashboard/src/lib/api.ts` - fetch wrapper with JWT handling

### UI Embedding
React build output embedded into Go binary via `//go:embed` directive pointing to `internal/ui/dist/`. After `pnpm build`, copy dist contents to `internal/ui/dist/` before Go build.

### API
- RESTful JSON under `/api/*`
- JWT Bearer token auth
- Setup wizard for initial admin creation
- Config: `/etc/opsagent/config.yaml`, DB: `/var/lib/opsagent/opsagent.db`, Port: 8787

## Security Model
- Commands classified: `read` (allowed), `write` (approval required), `blocked` (rejected)
- Blocked: `rm -rf /`, `mkfs`, `dd if=`, `chmod -R 777 /`, etc.
- All operations audit logged to SQLite

## Key Files
| Path | Purpose |
|------|---------|
| `cmd/opsagent/main.go` | Entry point |
| `internal/web/server.go` | HTTP server + all route handlers |
| `internal/config/config.go` | Config loading |
| `internal/usecase/evaluate_alerts.go` | Alert evaluation logic |
| `internal/usecase/create_task_plan.go` | AI task planning |
| `internal/infrastructure/executor/executor.go` | Command execution |
| `internal/policy/policy.go` | Security policy classifier |
| `web/dashboard/src/App.tsx` | Main React app |
| `web/dashboard/src/pages/SettingsPage.tsx` | Settings UI |
