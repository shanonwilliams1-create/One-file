package formats

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func init() {
	Register(yamlHandler{})
}

type yamlHandler struct{}

func (y yamlHandler) Name() string         { return "yaml" }
func (y yamlHandler) Extensions() []string { return []string{".yaml", ".yml"} }

func (y yamlHandler) Read(data []byte, key string) (string, error) {
	var v interface{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		return "", fmt.Errorf("yaml: failed to parse: %w", err)
	}

	parts := splitKey(key)
	val, err := getYAMLNested(v, parts)
	if err != nil {
		return "", fmt.Errorf("yaml: key %q not found: %w", key, err)
	}

	return fmt.Sprintf("%v", val), nil
}

func (y yamlHandler) Write(data []byte, key, value string) ([]byte, error) {
	var root interface{}
	if len(strings.TrimSpace(string(data))) > 0 {
		if err := yaml.Unmarshal(data, &root); err != nil {
			return nil, fmt.Errorf("yaml: failed to parse: %w", err)
		}
	}

	m, ok := root.(map[string]interface{})
	if !ok || root == nil {
		m = make(map[string]interface{})
	}

	parts := splitKey(key)
	if err := setYAMLNested(m, parts, parseYAMLValue(value)); err != nil {
		return nil, fmt.Errorf("yaml: failed to set key %q: %w", key, err)
	}

	result, err := yaml.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("yaml: failed to marshal: %w", err)
	}
	return result, nil
}

func getYAMLNested(v interface{}, parts []string) (interface{}, error) {
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
		return getYAMLNested(child, remaining)
	case map[interface{}]interface{}:
		child, ok := m[current]
		if !ok {
			return nil, fmt.Errorf("key %q not found", current)
		}
		return getYAMLNested(child, remaining)
	default:
		if idx, err := strconv.Atoi(current); err == nil {
			if arr, ok := v.([]interface{}); ok && idx >= 0 && idx < len(arr) {
				return getYAMLNested(arr[idx], remaining)
			}
		}
		return nil, fmt.Errorf("key %q not found", current)
	}
}

func setYAMLNested(m map[string]interface{}, parts []string, val interface{}) error {
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

	return setYAMLNested(childMap, remaining, val)
}

func parseYAMLValue(s string) interface{} {
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}
	if s == "~" || s == "null" {
		return nil
	}

	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}

	return s
}