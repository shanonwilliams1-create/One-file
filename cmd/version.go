package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.1.0-dev"

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of OmniConfig",
		Long:  `Print the current version of OmniConfig.`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("OmniConfig v%s\n", version)
		},
	}
}