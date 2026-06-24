package distro

import (
	"os"
	"path/filepath"
	"strings"
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
