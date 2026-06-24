package scanner

import (
	"time"
)

type Config struct {
	Paths    []string
	Exclude  []string
	DeepScan bool
	Verbose  bool
}

type Finding struct {
	File       string   `json:"file"`
	Threat     string   `json:"threat"`
	Severity   string   `json:"severity"`
	Md5        string   `json:"md5,omitempty"`
	Sha256     string   `json:"sha256,omitempty"`
	MimeType   string   `json:"mime_type,omitempty"`
	Size       int64    `json:"size"`
	Confidence float64  `json:"confidence"`
	Rules      []string `json:"rules,omitempty"`
}

type ScanOutput struct {
	Scanned  int
	Threats  int
	Duration time.Duration
	Findings []Finding
}

type Scanner struct {
	cfg Config
}

func New(cfg Config) *Scanner {
	return &Scanner{cfg: cfg}
}

func (s *Scanner) Run() (*ScanOutput, error) {
	start := time.Now()

	output := &ScanOutput{
		Findings: make([]Finding, 0),
	}

	for _, path := range s.cfg.Paths {
		findings, err := s.scanDir(path)
		if err != nil {
			continue
		}
		output.Findings = append(output.Findings, findings...)
		output.Scanned++
	}

	output.Threats = len(output.Findings)
	output.Duration = time.Since(start)

	return output, nil
}
