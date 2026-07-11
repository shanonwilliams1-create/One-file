package cmd

import (
	"fmt"

	"github.com/omnicofig/cli/pkg/profile"
	"github.com/spf13/cobra"
)

func NewWhoamiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show your profile name and business",
		Long:  `Quickly display your name and business from the profile.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := profile.Load()
			if err != nil {
				return fmt.Errorf("failed to load profile: %w", err)
			}

			if !p.IsComplete() {
				fmt.Println("No profile set. Run 'omniconfig profile init' to create one.")
				return nil
			}

			if p.Name != "" {
				fmt.Printf("%s", p.Name)
				if p.Business != "" {
					fmt.Printf(" — %s", p.Business)
				}
				fmt.Println()
			}

			return nil
		},
	}
}