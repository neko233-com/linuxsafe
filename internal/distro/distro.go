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
		strings.Contains(i.IDLike, "debian") || strings.Contains(i.IDLike, "ubuntu")
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
