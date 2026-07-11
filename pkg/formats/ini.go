package formats

import (
	"fmt"
	"strings"

	"gopkg.in/ini.v1"
)

func init() {
	Register(iniHandler{})
}

type iniHandler struct{}

func (i iniHandler) Name() string         { return "ini" }
func (i iniHandler) Extensions() []string { return []string{".ini", ".cfg", ".conf"} }

func (i iniHandler) Read(data []byte, key string) (string, error) {
	cfg, err := ini.Load(data)
	if err != nil {
		return "", fmt.Errorf("ini: failed to parse: %w", err)
	}

	// Parse key as [section]key or key
	section, k := parseINIKey(key)

	iniSection := cfg.Section(section)
	if !iniSection.HasKey(k) {
		return "", fmt.Errorf("ini: key %q not found", key)
	}

	return iniSection.Key(k).String(), nil
}

func (i iniHandler) Write(data []byte, key, value string) ([]byte, error) {
	cfg, err := ini.Load(data)
	if err != nil {
		// Create new config if file is empty/corrupt
		cfg = ini.Empty()
		// If data is non-empty but failed to parse, it's a real error
		if len(strings.TrimSpace(string(data))) > 0 {
			// Try to load as empty
			cfg, _ = ini.Load([]byte{})
		}
	}

	// Parse key as [section]key or key
	section, k := parseINIKey(key)

	iniSection := cfg.Section(section)
	iniSection.Key(k).SetValue(value)

	var buf strings.Builder
	if _, err := cfg.WriteTo(&buf); err != nil {
		return nil, fmt.Errorf("ini: failed to write: %w", err)
	}

	return []byte(buf.String()), nil
}

// parseINIKey parses a key like "section.key" into ("section", "key")
// or just "key" into (ini.DEFAULT_SECTION, "key")
func parseINIKey(key string) (string, string) {
	parts := strings.SplitN(key, ".", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return ini.DEFAULT_SECTION, parts[0]
}