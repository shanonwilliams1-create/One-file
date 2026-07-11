package cmd

import (
	"fmt"

	"github.com/omnicofig/cli/pkg/audit"
	"github.com/spf13/cobra"
)

var (
	logFile  string
	logSince string
)

func NewLogCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "View the audit log of config changes",
		Long: `View the audit log showing all config file changes made with
'omniconfig set'. Each entry shows who made the change, what
file and key was changed, and the old/new values.

Filters:
  --file <path>   Show only entries for a specific config file
  --since <date>  Show only entries after a date (RFC3339 format)

Examples:
  omniconfig log
  omniconfig log --file /tmp/app.yaml
  omniconfig log --since 2026-07-01`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := audit.ReadAll(logFile, logSince)
			if err != nil {
				return fmt.Errorf("failed to read audit log: %w", err)
			}

			if len(entries) == 0 {
				fmt.Println("No audit log entries found.")
				return nil
			}

			fmt.Printf("Audit log (%d entries)", len(entries))
			if logFile != "" {
				fmt.Printf(" for file: %s", logFile)
			}
			if logSince != "" {
				fmt.Printf(" since: %s", logSince)
			}
			fmt.Println()
			fmt.Println()

			for _, e := range entries {
				fmt.Println(e.Display())
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&logFile, "file", "", "Filter by config file path")
	cmd.Flags().StringVar(&logSince, "since", "", "Filter by date (RFC3339)")

	return cmd
}