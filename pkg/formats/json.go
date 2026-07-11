package formats

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func init() {
	Register(jsonHandler{})
}

type jsonHandler struct{}

func (j jsonHandler) Name() string         { return "json" }
func (j jsonHandler) Extensions() []string { return []string{".json"} }

func (j jsonHandler) Read(data []byte, key string) (string, error) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return "", fmt.Errorf("json: failed to parse: %w", err)
	}

	parts := splitKey(key)
	val, err := getNestedValue(v, parts)
	if err != nil {
		return "", fmt.Errorf("json: key %q not found: %w", key, err)
	}

	return fmt.Sprintf("%v", val), nil
}

func (j jsonHandler) Write(data []byte, key, value string) ([]byte, error) {
	var v interface{}
	if len(strings.TrimSpace(string(data))) > 0 {
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, fmt.Errorf("json: failed to parse: %w", err)
		}
	} else {
		// Empty file — start with empty object
		v = make(map[string]interface{})
	}

	parsed, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("json: root must be an object")
	}

	parts := splitKey(key)
	if err := setNestedValue(parsed, parts, parseValue(value)); err != nil {
		return nil, fmt.Errorf("json: failed to set key %q: %w", key, err)
	}

	result, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("json: failed to marshal: %w", err)
	}
	result = append(result, '\n')
	return result, nil
}

// getNestedValue navigates a nested structure using path segments.
func getNestedValue(v interface{}, parts []string) (interface{}, error) {
	if len(parts) == 0 {
		return v, nil
	}

	current := parts[0]
	remaining := parts[1:]

	switch m := v.(type) {
	case map[string]interface{}:
		child, ok := m[current]
		if !ok {
			return nil, fmt.Errorf("key %q not found", current)
		}
		return getNestedValue(child, remaining)
	default:
		// Try array index
		if idx, err := strconv.Atoi(current); err == nil {
			if arr, ok := v.([]interface{}); ok && idx >= 0 && idx < len(arr) {
				return getNestedValue(arr[idx], remaining)
			}
		}
		return nil, fmt.Errorf("key %q not found", current)
	}
}

// setNestedValue sets a value at a nested path, creating intermediate maps as needed.
func setNestedValue(m map[string]interface{}, parts []string, val interface{}) error {
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

	return setNestedValue(childMap, remaining, val)
}

// parseValue attempts to parse a string value into the appropriate Go type.
func parseValue(s string) interface{} {
	// Try bool
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}

	// Try number (int first, then float)
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}

	// Try JSON array/object
	if strings.HasPrefix(s, "{") || strings.HasPrefix(s, "[") {
		var parsed interface{}
		if err := json.Unmarshal([]byte(s), &parsed); err == nil {
			return parsed
		}
	}

	// Default: string
	return s
}