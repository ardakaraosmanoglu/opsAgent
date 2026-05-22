# OpsAgent

Local-first Linux AI ops agent. Single-binary, SQLite backend, web dashboard.

## Prerequisites

- Go 1.25+
- Linux (x86_64)
- SQLite3 (usually pre-installed)

## Build

```bash
go build -o opsagent ./cmd/opsagent
```

## Run

```bash
./opsagent
```

Agent starts on `0.0.0.0:8080` by default. Dashboard at `http://localhost:8080`.

## Install as Systemd Service

```bash
sudo cp opsagent /usr/local/bin/
sudo cp scripts/install.sh /usr/local/bin/opsagent-install
sudo opsagent-install
```

## Configuration

Edit `config.yaml` before running:

```yaml
server:
  host: "0.0.0.0"
  port: 8080

database:
  path: "./opsagent.db"

ai:
  # Optional: set API key via environment variable
  # OPENAI_API_KEY=sk-... ./opsagent
```

## First Setup

1. Open dashboard at `http://localhost:8080`
2. Create admin account
3. Optionally configure AI API key for AI-assisted analysis
4. Agent begins monitoring automatically

## Uninstall

```bash
sudo scripts/uninstall.sh
```
