package formats

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	Register(envHandler{})
}

type envHandler struct{}

func (e envHandler) Name() string         { return "env" }
func (e envHandler) Extensions() []string { return []string{".env"} }

func (e envHandler) Read(data []byte, key string) (string, error) {
	envMap, err := godotenv.Unmarshal(string(data))
	if err != nil {
		return "", fmt.Errorf("env: failed to parse: %w", err)
	}

	val, ok := envMap[key]
	if !ok {
		return "", fmt.Errorf("env: key %q not found", key)
	}

	return val, nil
}

func (e envHandler) Write(data []byte, key, value string) ([]byte, error) {
	envMap, err := godotenv.Unmarshal(string(data))
	if err != nil {
		// Start fresh if parsing fails
		envMap = make(map[string]string)
	}

	envMap[key] = value

	// Serialize back to .env format
	var buf strings.Builder
	for k, v := range envMap {
		// Quote value if it contains spaces or special chars
		if strings.ContainsAny(v, " #\"'") {
			v = fmt.Sprintf("%q", v)
		}
		buf.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}

	return []byte(buf.String()), nil
}