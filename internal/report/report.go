package report

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/neko233-com/linuxsafe/internal/distro"
)

type ScanReport struct {
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	Distro    *distro.Info           `json:"distro"`
	Summary   ReportSummary          `json:"summary"`
	Findings  []FindingReport        `json:"findings,omitempty"`
	Actions   []ActionReport         `json:"actions,omitempty"`
	System    *SystemInfo            `json:"system,omitempty"`
}

type ReportSummary struct {
	TotalFiles   int    `json:"total_files"`
	ThreatsFound int    `json:"threats_found"`
	CleanFiles   int    `json:"clean_files"`
	Duration     string `json:"duration"`
}

type FindingReport struct {
	File       string  `json:"file"`
	Threat     string  `json:"threat"`
	Severity   string  `json:"severity"`
	Confidence float64 `json:"confidence"`
	Md5        string  `json:"md5,omitempty"`
	Sha256     string  `json:"sha256,omitempty"`
}

type ActionReport struct {
	Type    string `json:"type"`
	Target  string `json:"target"`
	Action  string `json:"action"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type SystemInfo struct {
	Hostname string `json:"hostname"`
	Kernel   string `json:"kernel"`
	Uptime   string `json:"uptime"`
}

func NewScanReport(findings []FindingReport, duration time.Duration, totalFiles int) *ScanReport {
	return &ScanReport{
		Type:      "scan",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Distro:    distro.Detect(),
		Summary: ReportSummary{
			TotalFiles:   totalFiles,
			ThreatsFound: len(findings),
			CleanFiles:   totalFiles - len(findings),
			Duration:     duration.String(),
		},
		Findings: findings,
	}
}

func NewDetectReport(findings []FindingReport) *ScanReport {
	return &ScanReport{
		Type:      "detect",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Distro:    distro.Detect(),
		Summary: ReportSummary{
			ThreatsFound: len(findings),
		},
		Findings: findings,
	}
}

func NewCleanReport(actions []ActionReport) *ScanReport {
	return &ScanReport{
		Type:      "clean",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Distro:    distro.Detect(),
		Actions:   actions,
	}
}

func (r *ScanReport) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func (r *ScanReport) ToHTML() string {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Linuxsafe Report - %s</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, sans-serif; max-width: 900px; margin: 0 auto; padding: 20px; background: #f5f5f5; }
        .header { background: #2563eb; color: white; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
        .summary { display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px; margin-bottom: 20px; }
        .card { background: white; padding: 15px; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
        .card h3 { margin: 0; font-size: 24px; }
        .card p { margin: 5px 0 0; color: #666; }
        table { width: 100%%; border-collapse: collapse; background: white; border-radius: 8px; overflow: hidden; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #eee; }
        th { background: #f8f9fa; }
        .critical { color: #dc2626; }
        .high { color: #ea580c; }
        .medium { color: #ca8a04; }
        .low { color: #16a34a; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Linuxsafe Security Report</h1>
        <p>%s | %s</p>
    </div>
    <div class="summary">
        <div class="card">
            <h3>%d</h3>
            <p>Total Files</p>
        </div>
        <div class="card">
            <h3 class="critical">%d</h3>
            <p>Threats Found</p>
        </div>
        <div class="card">
            <h3 class="low">%d</h3>
            <p>Clean Files</p>
        </div>
        <div class="card">
            <h3>%s</h3>
            <p>Duration</p>
        </div>
    </div>
`, r.Type, r.Timestamp, r.Distro.PrettyName, r.Summary.TotalFiles, r.Summary.ThreatsFound, r.Summary.CleanFiles, r.Summary.Duration)

	if len(r.Findings) > 0 {
		html += `<table>
        <tr><th>File</th><th>Threat</th><th>Severity</th><th>Confidence</th></tr>`
		for _, f := range r.Findings {
			html += fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td class="%s">%s</td><td>%.0f%%</td></tr>`,
				f.File, f.Threat, f.Severity, f.Severity, f.Confidence*100)
		}
		html += `</table>`
	}

	html += `</body></html>`
	return html
}

func (r *ScanReport) WriteJSON(path string) error {
	data, err := r.ToJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (r *ScanReport) WriteHTML(path string) error {
	return os.WriteFile(path, []byte(r.ToHTML()), 0644)
}
