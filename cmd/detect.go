package cmd

import (
	"fmt"
	"strings"

	"github.com/omnicofig/cli/pkg/detect"
	"github.com/spf13/cobra"
)

func NewDetectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "detect",
		Short: "Detect operating system and config file format",
		Long: `Detect the current operating system and common configuration file
locations and formats, including mobile platforms (iOS, Android).

OmniConfig detects:
  - Operating system: Linux, macOS, Windows, iOS, Android
  - Common config file locations for the detected OS
  - Config format by file extension or content sniffing`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			osInfo := detect.DetectOS()

			fmt.Printf("Operating System: %s (%s)\n", osInfo.GOOS, osInfo.PrettyName)
			fmt.Printf("Native Binary:    %s\n", osInfo.NativeBinary)
			fmt.Println()

			if len(osInfo.ConfigDirs) > 0 {
				fmt.Println("Config Directories:")
				for _, d := range osInfo.ConfigDirs {
					fmt.Printf("  - %s\n", d)
				}
				fmt.Println()
			}

			if len(osInfo.ConfigFiles) > 0 {
				fmt.Println("Common Config Files:")
				for _, f := range osInfo.ConfigFiles {
					fmt.Printf("  - %s\n", f)
				}
				fmt.Println()
			}

			fmt.Println("Supported Formats:")
			for _, f := range detect.AllFormats {
				exts := strings.Join(f.Extensions, ", ")
				if exts == "" {
					exts = "(detected by content)"
				}
				fmt.Printf("  - %s (%s)\n", strings.ToUpper(f.Name), exts)
			}
			return nil
		},
	}
}