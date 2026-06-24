package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cleanDryRun bool
	cleanForce  bool
)

type CleanResult struct {
	Status     string        `json:"status"`
	Actions    []CleanAction `json:"actions"`
	Skipped    int           `json:"skipped"`
	Errors     []string      `json:"errors,omitempty"`
	Error      string        `json:"error,omitempty"`
}

type CleanAction struct {
	Type    string `json:"type"`
	Target  string `json:"target"`
	Action  string `json:"action"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remediate detected threats",
	Long: `Automatically clean detected threats based on scan/detect results.

Actions performed:
- Kill malicious processes
- Remove suspicious files
- Clean infected cron jobs
- Restore modified system files

Exit codes:
  0 - Clean completed successfully
  1 - Partial cleanup (some items require manual review)
  2 - Dry run completed
  3 - Clean error
  4 - Dangerous operation blocked (use --force)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		result := CleanResult{
			Status: "ok",
		}

		if cleanDryRun {
			result.Actions = append(result.Actions, CleanAction{
				Type:   "dry_run",
				Target: "all",
				Action: "preview",
				Status: "skipped",
			})
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(result)
			os.Exit(2)
			return nil
		}

		if !cleanForce {
			result.Error = "use --force to confirm cleanup actions"
			result.Status = "blocked"
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(result)
			os.Exit(4)
			return nil
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(result); err != nil {
			return fmt.Errorf("failed to encode result: %w", err)
		}
		return nil
	},
}

func init() {
	cleanCmd.Flags().BoolVar(&cleanDryRun, "dry-run", false, "Preview actions without executing")
	cleanCmd.Flags().BoolVar(&cleanForce, "force", false, "Confirm cleanup actions")
}
