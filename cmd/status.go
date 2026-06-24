package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/neko233-com/linuxsafe/internal/distro"
	"github.com/neko233-com/linuxsafe/internal/hardware"
	"github.com/spf13/cobra"
)

type StatusResult struct {
	Status      string         `json:"status"`
	Version     string         `json:"version"`
	Distro      *distro.Info   `json:"distro"`
	Hardware    *hardware.Info `json:"hardware,omitempty"`
	Signatures  string         `json:"signatures"`
	LastScan    string         `json:"last_scan"`
	LastUpdate  string         `json:"last_update"`
	EngineReady bool           `json:"engine_ready"`
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show linuxsafe status",
	Long: `Display current status of the linuxsafe agent:
- Engine readiness
- Signature database version
- Detected Linux distribution
- Hardware information
- Last scan/update time`,
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
