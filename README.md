# linuxsafe

Agent-first Linux security CLI. Automated threat detection, investigation, and remediation.

## Quick Install

```bash
curl -sSL https://raw.githubusercontent.com/neko233-com/linuxsafe/main/install.sh | bash
```

## Build

```bash
go build -o linuxsafe .
```

## Usage

```bash
# Scan for threats
linuxsafe scan /

# Detect system threats
linuxsafe detect

# View hardware info
linuxsafe hw

# Generate security report
linuxsafe report

# Check status (includes distro detection)
linuxsafe status

# Investigate a specific target
linuxsafe investigate <pid|path|user>

# Clean threats (dry-run first)
linuxsafe clean --dry-run
linuxsafe clean --force
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success / Clean |
| 1 | Threats found / Warnings |
| 2 | Scan error / Critical threats |
| 3 | Configuration error |
| 4 | Dangerous operation blocked |

## Features

- **Threat Detection**: Hash-based malware detection, pattern matching
- **Rootkit Detection**: LKM, /dev/mem, /dev/kmem checks
- **Persistence Checks**: Cron jobs, SSH keys, SUID anomalies
- **Auto Cleanup**: Automated remediation with safety controls
- **Reports**: JSON and HTML security reports
- **Hardware Info**: CPU, memory, disk, network detection
- **Distro Auto-Detect**: Ubuntu, Debian, CentOS, Alpine support

## Docker Testing

Test across multiple Linux distributions:

```bash
# Run all Docker tests
make test-docker

# Test specific distro
make test-docker-ubuntu
make test-docker-debian
make test-docker-centos
make test-docker-alpine
```

### Supported Distributions

| Distribution | Status | Package Manager |
|-------------|--------|-----------------|
| Ubuntu 22.04 | ✅ Tested | apt-get |
| Debian 12 | ✅ Tested | apt-get |
| CentOS Stream 9 | ✅ Tested | yum |
| Alpine 3.19 | ✅ Tested | apk |

## CI/CD

- **CI**: Runs on every push/PR (lint, test, build)
- **Release**: Auto-builds on tag push (v*)
- **Docs**: Auto-deploys to GitHub Pages

## Architecture

```
linuxsafe/
├── cmd/                    # CLI commands (cobra)
│   ├── root.go            # Root command
│   ├── scan.go            # File scanning
│   ├── detect.go          # Threat detection
│   ├── clean.go           # Remediation
│   ├── investigate.go     # Forensic analysis
│   ├── status.go          # Status (with distro)
│   ├── hw.go              # Hardware info
│   ├── report.go          # Report generation
│   └── update.go          # Signature updates
├── internal/
│   ├── scanner/           # File scanning engine
│   ├── distro/            # Distro detection
│   ├── hardware/          # Hardware detection
│   └── report/            # Report generation
├── test/
│   ├── docker/            # Docker test files
│   └── fixtures/          # Test data
└── docs/                  # Documentation
```

## Agent Integration

All commands output structured JSON by default. Parse stdout for results, exit code for decision.

```bash
# Agent workflow
result=$(linuxsafe scan /home)
if [ $? -eq 1 ]; then
  linuxsafe clean --dry-run
  linuxsafe clean --force
fi
```

## License

MIT
