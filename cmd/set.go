package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/omnicofig/cli/pkg/detect"
	"github.com/omnicofig/cli/pkg/formats"
	"github.com/spf13/cobra"
)

func NewSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Write a configuration value",
		Long: `Write a configuration value by key to the auto-detected config file.

OmniConfig will automatically detect your operating system and the
config file format based on the file extension or content. If a
specific config file is not provided via --config, it will search
common locations for your OS.

Creates a backup of the original file before modifying (.bak timestamp).

Supports dot-notation keys for nested formats (e.g., "database.host").

Examples:
  omniconfig set database.host localhost
  omniconfig set --config ./app.yaml server.port 8080`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			configPath, err := resolveConfigPath(cfgFile)
			if err != nil {
				// If file doesn't exist and --config was given, we can create it
				if cfgFile != "" {
					configPath = cfgFile
				} else {
					return fmt.Errorf("unable to resolve config file: %w", err)
				}
			}

			// Read existing file (if it exists)
			var data []byte
			if _, statErr := os.Stat(configPath); statErr == nil {
				data, err = os.ReadFile(configPath)
				if err != nil {
					return fmt.Errorf("cannot read config file %s: %w", configPath, err)
				}

				// Create backup before modifying
				backupPath := configPath + ".bak." + time.Now().Format("20060102-150405")
				if err := os.WriteFile(backupPath, data, 0644); err != nil {
					return fmt.Errorf("failed to create backup: %w", err)
				}
				fmt.Fprintf(os.Stderr, "Backup saved to: %s\n", backupPath)
			}

			// Detect format
			ext := filepath.Ext(configPath)
			handler := formats.GetByExtension(ext)
			if handler == nil && len(data) > 0 {
				f := detect.DetectFormatFromContent(string(data))
				handler = formats.Get(f.Name)
			}
			if handler == nil {
				// Default to JSON if we can't detect
				handler = formats.Get("json")
				if handler == nil {
					return fmt.Errorf("no format handler available")
				}
				// If creating a new file, ensure it has .json extension
				if ext == "" {
					configPath += ".json"
				}
			}

			// Write the value
			newData, err := handler.Write(data, key, value)
			if err != nil {
				return fmt.Errorf("failed to set %q: %w", key, err)
			}

			// Write back to file
			if err := os.WriteFile(configPath, newData, 0644); err != nil {
				return fmt.Errorf("failed to write config file: %w", err)
			}

			fmt.Printf("Set %s = %s in %s\n", key, value, configPath)
			return nil
		},
	}
}