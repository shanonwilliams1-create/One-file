// Package formats provides a pluggable interface for reading and writing
// configuration files in various formats (JSON, YAML, TOML, INI, XML, .env, etc.).
package formats

import "strings"

// Handler defines the interface for reading and writing config files.
type Handler interface {
	// Name returns the format name (e.g., "json", "yaml").
	Name() string

	// Extensions returns the file extensions for this format.
	Extensions() []string

	// Read reads a value by key from the config data.
	// Keys support dot-separated paths (e.g., "database.host").
	Read(data []byte, key string) (string, error)

	// Write sets a value by key in the config data.
	// Returns the modified config bytes.
	Write(data []byte, key, value string) ([]byte, error)
}

// Registry maps format names to their handlers.
var Registry = make(map[string]Handler)

// Register adds a format handler to the global registry.
func Register(h Handler) {
	Registry[h.Name()] = h
}

// Get returns the handler for a given format name.
// Returns nil if the format is not registered.
func Get(name string) Handler {
	return Registry[name]
}

// GetByExtension returns the handler for a given file extension.
// Returns nil if no handler supports that extension.
func GetByExtension(ext string) Handler {
	ext = strings.ToLower(ext)
	for _, h := range Registry {
		for _, e := range h.Extensions() {
			if e == ext {
				return h
			}
		}
	}
	return nil
}

// splitKey splits a dot-notation key into path segments.
// e.g., "database.host" -> ["database", "host"]
func splitKey(key string) []string {
	return strings.Split(key, ".")
}

// joinKey joins path segments into a dot-notation key.
func joinKey(parts []string) string {
	return strings.Join(parts, ".")
}