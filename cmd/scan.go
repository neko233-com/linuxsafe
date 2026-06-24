package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/neko233-com/linuxsafe/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	scanPaths   []string
	scanExclude []string
	scanDeep    bool
)

type ScanResult struct {
	Status    string              `json:"status"`
	Scanned   int                 `json:"scanned"`
	Threats   int                 `json:"threats"`
	Clean     int                 `json:"clean"`
	Duration  string              `json:"duration"`
	Findings  []scanner.Finding   `json:"findings,omitempty"`
	Error     string              `json:"error,omitempty"`
}

var scanCmd = &cobra.Command{
	Use:   "scan [paths...]",
	Short: "Scan filesystem for threats",
	Long: `Scan specified paths (default: /) for malware, rootkits, and suspicious patterns.

Exit codes:
  0 - Clean, no threats found
  1 - Threats found
  2 - Scan error (permission denied, I/O error, etc.)
  3 - Configuration error`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		paths := args
		if len(paths) == 0 {
			paths = []string{"/"}
		}

		s := scanner.New(scanner.Config{
			Paths:    paths,
			Exclude:  scanExclude,
			DeepScan: scanDeep,
			Verbose:  verbose,
		})

		result, err := s.Run()
		if err != nil {
			return outputResult(ScanResult{
				Status: "error",
				Error:  err.Error(),
			}, 2)
		}

		sr := ScanResult{
			Status:   "ok",
			Scanned:  result.Scanned,
			Threats:  result.Threats,
			Clean:    result.Scanned - result.Threats,
			Duration: result.Duration.String(),
			Findings: result.Findings,
		}

		exitCode := 0
		if result.Threats > 0 {
			exitCode = 1
		}

		return outputResult(sr, exitCode)
	},
}

func init() {
	scanCmd.Flags().StringSliceVarP(&scanExclude, "exclude", "e", nil, "Exclude paths (glob patterns)")
	scanCmd.Flags().BoolVarP(&scanDeep, "deep", "d", false, "Deep scan (includes archive inspection)")
}

func outputResult(result interface{}, exitCode int) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		return fmt.Errorf("failed to encode result: %w", err)
	}
	if exitCode != 0 {
		os.Exit(exitCode)
	}
	return nil
}
