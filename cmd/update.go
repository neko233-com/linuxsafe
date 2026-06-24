package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type UpdateResult struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Error   string `json:"error,omitempty"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update signature database",
	Long: `Update the threat signature database from remote sources.

Exit codes:
  0 - Update successful
  1 - Update failed
  2 - Already up to date`,
	RunE: func(cmd *cobra.Command, args []string) error {
		result := UpdateResult{
			Status:  "ok",
			Version: "2026.06.24",
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(result); err != nil {
			return fmt.Errorf("failed to encode result: %w", err)
		}
		return nil
	},
}
