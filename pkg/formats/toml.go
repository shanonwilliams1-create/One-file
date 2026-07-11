package formats

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

func init() {
	Register(tomlHandler{})
}

type tomlHandler struct{}

func (t tomlHandler) Name() string         { return "toml" }
func (t tomlHandler) Extensions() []string { return []string{".toml"} }

func (t tomlHandler) Read(data []byte, key string) (string, error) {
	var v map[string]interface{}
	if err := toml.Unmarshal(data, &v); err != nil {
		return "", fmt.Errorf("toml: failed to parse: %w", err)
	}

	parts := splitKey(key)
	val, err := getTOMLNested(v, parts)
	if err != nil {
		return "", fmt.Errorf("toml: key %q not found: %w", key, err)
	}

	return fmt.Sprintf("%v", val), nil
}

func (t tomlHandler) Write(data []byte, key, value string) ([]byte, error) {
	var v map[string]interface{}
	if len(strings.TrimSpace(string(data))) > 0 {
		if _, err := toml.Decode(string(data), &v); err != nil {
			return nil, fmt.Errorf("toml: failed to parse: %w", err)
		}
	}
	if v == nil {
		v = make(map[string]interface{})
	}

	parts := splitKey(key)
	if err := setTOMLNested(v, parts, parseTOMLValue(value)); err != nil {
		return nil, fmt.Errorf("toml: failed to set key %q: %w", key, err)
	}

	// TOML doesn't have a standard marshal-to-string, so we use the Go library
	var buf strings.Builder
	enc := toml.NewEncoder(&buf)
	enc.Indent = ""
	if err := enc.Encode(v); err != nil {
		return nil, fmt.Errorf("toml: failed to encode: %w", err)
	}

	return []byte(buf.String()), nil
}

func getTOMLNested(v map[string]interface{}, parts []string) (interface{}, error) {
	if len(parts) == 0 {
		return v, nil
	}

	current := parts[0]
	remaining := parts[1:]

	child, ok := v[current]
	if !ok {
		return nil, fmt.Errorf("key %q not found", current)
	}

	if len(remaining) == 0 {
		return child, nil
	}

	childMap, ok := child.(map[string]interface{})
	if !ok {
		// Try to convert from toml's internal map type
		if m, ok := child.(map[interface{}]interface{}); ok {
			childMap = make(map[string]interface{})
			for k, val := range m {
				childMap[fmt.Sprintf("%v", k)] = val
			}
		} else {
			return nil, fmt.Errorf("key %q is not a table/section", current)
		}
	}

	return getTOMLNested(childMap, remaining)
}

func setTOMLNested(m map[string]interface{}, parts []string, val interface{}) error {
	if len(parts) == 0 {
		return nil
	}

	if len(parts) == 1 {
		m[parts[0]] = val
		return nil
	}

	current := parts[0]
	remaining := parts[1:]

	child, exists := m[current]
	if !exists {
		child = make(map[string]interface{})
		m[current] = child
	}

	childMap, ok := child.(map[string]interface{})
	if !ok {
		childMap = make(map[string]interface{})
		m[current] = childMap
	}

	return setTOMLNested(childMap, remaining, val)
}

func parseTOMLValue(s string) interface{} {
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}

	// Simple integer/float parsing
	var i int
	if _, err := fmt.Sscanf(s, "%d", &i); err == nil {
		// Check if it's actually a float
		if !strings.Contains(s, ".") {
			return i
		}
	}
	var f float64
	if _, err := fmt.Sscanf(s, "%f", &f); err == nil {
		return f
	}

	return s
}