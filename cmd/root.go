package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile       string
	listPlatforms bool
)

var rootCmd = &cobra.Command{
	Use:   "omniconfig",
	Short: "One CLI tool to read, write, and manage config files across any OS",
	Long: `OmniConfig — one CLI tool that reads, writes, and manages configuration
files across any operating system (Linux, macOS, Windows, iOS, Android)
in any format (JSON, YAML, TOML, INI, XML, .env, and more).

Auto-detects the OS and config format so you never have to think about it.
Just run one command, anywhere.`,
	Run: func(cmd *cobra.Command, args []string) {
		if listPlatforms {
			printPlatforms()
			return
		}
		cmd.Help()
	},
}

// printPlatforms prints all supported target platforms
func printPlatforms() {
	fmt.Println("OmniConfig Supported Platforms")
	fmt.Println("==============================")
	fmt.Println()
	fmt.Println("Desktop:")
	fmt.Println("  linux/amd64     — Linux x86_64")
	fmt.Println("  linux/arm64     — Linux ARM64 (Raspberry Pi, AWS Graviton)")
	fmt.Println("  darwin/amd64    — macOS Intel (x86_64)")
	fmt.Println("  darwin/arm64    — macOS Apple Silicon (M1/M2/M3/M4)")
	fmt.Println("  windows/amd64   — Windows x86_64")
	fmt.Println("  windows/arm64   — Windows ARM64 (Surface Pro X, etc.)")
	fmt.Println()
	fmt.Println("Mobile:")
	fmt.Println("  ios/arm64       — iOS/iPadOS devices (requires Xcode)")
	fmt.Println("  android/arm64   — Android ARM64 devices (requires NDK)")
	fmt.Println("  android/amd64   — Android x86_64 (emulator)")
	fmt.Println()
	fmt.Println("Binary naming:  omniconfig-{os}-{arch}[.exe]")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to config file (default: auto-detected per OS)")
	rootCmd.PersistentFlags().BoolVar(&listPlatforms, "list-platforms", false, "list all supported target platforms")
	rootCmd.AddCommand(NewDetectCmd())
	rootCmd.AddCommand(NewGetCmd())
	rootCmd.AddCommand(NewSetCmd())
	rootCmd.AddCommand(NewVersionCmd())
}