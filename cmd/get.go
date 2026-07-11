package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/omnicofig/cli/pkg/detect"
	"github.com/omnicofig/cli/pkg/formats"
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Read a configuration value",
		Long: `Read a configuration value by key from the auto-detected config file.

OmniConfig will automatically detect your operating system and the
config file format based on the file extension or content. If a
specific config file is not provided via --config, it will search
common locations for your OS.

Supports dot-notation keys for nested formats (e.g., "database.host").

Examples:
  omniconfig get database.host
  omniconfig get --config ./app.yaml server.port`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]

			configPath, err := resolveConfigPath(cfgFile)
			if err != nil {
				return fmt.Errorf("unable to resolve config file: %w", err)
			}

			// Read the file
			data, err := os.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("cannot read config file %s: %w", configPath, err)
			}

			// Detect format
			ext := filepath.Ext(configPath)
			handler := formats.GetByExtension(ext)
			if handler == nil {
				// Try content-based detection
				f := detect.DetectFormatFromContent(string(data))
				handler = formats.Get(f.Name)
			}
			if handler == nil {
				return fmt.Errorf("unsupported config format: %s", ext)
			}

			// Read the value
			val, err := handler.Read(data, key)
			if err != nil {
				return fmt.Errorf("failed to read %q: %w", key, err)
			}

			fmt.Println(val)
			return nil
		},
	}
}

// resolveConfigPath returns the config file path to use.
// If cfgFile is set, use it. Otherwise, search common locations.
func resolveConfigPath(cfgFile string) (string, error) {
	if cfgFile != "" {
		// User specified a path
		abs, err := filepath.Abs(cfgFile)
		if err != nil {
			return "", err
		}
		if _, err := os.Stat(abs); err != nil {
			return "", fmt.Errorf("config file not found: %s", abs)
		}
		return abs, nil
	}

	// Auto-detect: search common config files
	osInfo := detect.DetectOS()
	for _, dir := range osInfo.ConfigDirs {
		expanded := expandPath(dir)
		entries, err := os.ReadDir(expanded)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			ext := filepath.Ext(entry.Name())
			if formats.GetByExtension(ext) != nil {
				return filepath.Join(expanded, entry.Name()), nil
			}
		}
	}

	// Check common individual config files
	for _, f := range osInfo.ConfigFiles {
		expanded := expandPath(f)
		if _, err := os.Stat(expanded); err == nil {
			ext := filepath.Ext(expanded)
			if formats.GetByExtension(ext) != nil {
				return expanded, nil
			}
		}
	}

	return "", fmt.Errorf("no config file found. Use --config to specify one")
}

// expandPath expands ~ to the home directory.
func expandPath(path string) string {
	if len(path) > 1 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err == nil {
			return home + path[1:]
		}
	}
	return path
}