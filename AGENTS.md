# Agent Instructions for linuxsafe

## Build & Test
- Build: `go build -o linuxsafe .`
- Run: `./linuxsafe <command> [args]`
- All output is JSON by default (agent-first design)
- Exit codes: 0=ok, 1=threats/warnings, 2=error/critical, 3=config, 4=blocked

## Code Style
- Go standard library preferred
- cobra for CLI framework
- Structured JSON output for all commands
- Deterministic exit codes for agent decision-making
- Idempotent operations safe for automation

## Architecture
- `cmd/` - CLI entry points
- `internal/scanner/` - File scanning, hash analysis, pattern matching
- `internal/detector/` - System threat detection (rootkit, persistence)
- `internal/signature/` - Malware signature database
- `internal/agent/` - Agent protocol, structured I/O

## Key Principles
1. Agent-first: Every output machine-parseable
2. Idempotent: Safe to re-run any command
3. Deterministic: Same input = same output
4. Autonomous: Can run full scan→detect→clean cycle without human input
