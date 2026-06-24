package cmd

import (
	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "linuxsafe",
	Short: "Linux security agent — automated threat detection, investigation, and remediation",
	Long: `linuxsafe is an agent-first CLI for Linux server security.

Designed for both human operators and AI agents:
- Structured JSON output by default (--json, always on)
- Human-readable via --pretty flag
- Deterministic exit codes for agent decision-making
- Idempotent operations safe for automation`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", true, "Output JSON (default, agent-first)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().MarkHidden("json")

	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(detectCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(investigateCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(updateCmd)
}
