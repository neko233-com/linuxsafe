# linuxsafe

Agent-first Linux security CLI. Automated threat detection, investigation, and remediation.

## Build

```bash
go build -o linuxsafe .
```

## Usage

```bash
# Scan for threats (JSON output by default)
./linuxsafe scan /

# Detect system threats (rootkit, backdoors, persistence)
./linuxsafe detect

# Investigate a specific target
./linuxsafe investigate <pid|path|user>

# Clean threats (dry-run first)
./linuxsafe clean --dry-run
./linuxsafe clean --force

# Check status
./linuxsafe status

# Update signatures
./linuxsafe update
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success / Clean |
| 1 | Threats found / Warnings |
| 2 | Scan error / Critical threats |
| 3 | Configuration error |
| 4 | Dangerous operation blocked |

## Agent Integration

All commands output structured JSON by default. Parse stdout for results, exit code for decision.

```bash
# Agent workflow
result=$(./linuxsafe scan /home)
if [ $? -eq 1 ]; then
  ./linuxsafe clean --dry-run
  ./linuxsafe clean --force
fi
```

## Architecture

```
cmd/           CLI commands (cobra)
internal/
  scanner/     File scanning engine
  detector/    Threat detection (planned)
  signature/   Signature database (planned)
  agent/       Agent protocol (planned)
```
