package hardware

import (
	"os"
	"runtime"
	"strconv"
	"strings"
)

type Info struct {
	Platform    string        `json:"platform"`
	Arch        string        `json:"arch"`
	CPUCount    int           `json:"cpu_count"`
	CPUModel    string        `json:"cpu_model"`
	MemoryTotal uint64        `json:"memory_total_bytes"`
	MemoryHuman string        `json:"memory_human"`
	Hostname    string        `json:"hostname"`
	Kernel      string        `json:"kernel"`
	Disks       []DiskInfo    `json:"disks,omitempty"`
	Network     []NetworkInfo `json:"network,omitempty"`
}

type DiskInfo struct {
	Mount   string `json:"mount"`
	FsType  string `json:"fs_type"`
	Total   uint64 `json:"total_bytes"`
	Free    uint64 `json:"free_bytes"`
	Used    uint64 `json:"used_bytes"`
	Percent  float64 `json:"used_percent"`
}

type NetworkInfo struct {
	Name   string   `json:"name"`
	IPs    []string `json:"ips"`
	MAC    string   `json:"mac,omitempty"`
}

func Detect() *Info {
	info := &Info{
		Platform: runtime.GOOS,
		Arch:     runtime.GOARCH,
		CPUCount: runtime.NumCPU(),
	}

	hostname, _ := os.Hostname()
	info.Hostname = hostname

	info.CPUModel = detectCPUModel()
	info.MemoryTotal, info.MemoryHuman = detectMemory()
	info.Kernel = detectKernel()
	info.Disks = detectDisks()
	info.Network = detectNetwork()

	return info
}

func detectCPUModel() string {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "unknown"
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "model name") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "unknown"
}

func detectMemory() (uint64, string) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, "unknown"
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				valStr := strings.TrimSpace(parts[1])
				valStr = strings.ReplaceAll(valStr, " kB", "")
				val, _ := strconv.ParseUint(valStr, 10, 64)
			 bytes := val * 1024
				return bytes, formatBytes(bytes)
			}
		}
	}
	return 0, "unknown"
}

func detectKernel() string {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return "unknown"
	}
	parts := strings.SplitN(string(data), " ", 3)
	if len(parts) >= 1 {
		return parts[0]
	}
	return "unknown"
}

func detectDisks() []DiskInfo {
	return nil
}

func detectNetwork() []NetworkInfo {
	return nil
}

func formatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)
	switch {
	case bytes >= TB:
		return strconv.FormatFloat(float64(bytes)/float64(TB), 'f', 1, 64) + " TB"
	case bytes >= GB:
		return strconv.FormatFloat(float64(bytes)/float64(GB), 'f', 1, 64) + " GB"
	case bytes >= MB:
		return strconv.FormatFloat(float64(bytes)/float64(MB), 'f', 1, 64) + " MB"
	case bytes >= KB:
		return strconv.FormatFloat(float64(bytes)/float64(KB), 'f', 1, 64) + " KB"
	default:
		return strconv.FormatUint(bytes, 10) + " B"
	}
}
