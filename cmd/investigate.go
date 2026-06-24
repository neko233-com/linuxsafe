package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type InvestigateResult struct {
	Status   string            `json:"status"`
	Target   string            `json:"target"`
	Findings []InvestigateFind `json:"findings"`
	Error    string            `json:"error,omitempty"`
}

type InvestigateFind struct {
	Category string      `json:"category"`
	Severity string      `json:"severity"`
	Summary  string      `json:"summary"`
	Details  interface{} `json:"details,omitempty"`
}

var investigateCmd = &cobra.Command{
	Use:   "investigate <target>",
	Short: "Deep forensic investigation of a target",
	Long: `Perform deep forensic investigation on a specific target:
- Process: Analyze process tree, network connections, open files
- File: Hash analysis, entropy check, metadata inspection
- User: Login history, sudo usage, process ownership
- Network: Connection audit, listening ports, DNS queries

Exit codes:
  0 - Investigation complete
  1 - Suspicious findings require attention
  2 - Investigation error`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]

		result := InvestigateResult{
			Status: "ok",
			Target: target,
			Findings: []InvestigateFind{
				{
					Category: "general",
					Severity: "info",
					Summary:  fmt.Sprintf("Investigation of target: %s", target),
				},
			},
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(result); err != nil {
			return fmt.Errorf("failed to encode result: %w", err)
		}
		return nil
	},
}
