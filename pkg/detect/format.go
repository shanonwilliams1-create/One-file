package detect

import (
    "path/filepath"
    "strings"
)

// Format represents a supported config file format.
type Format struct {
    Name       string
    Extensions []string
    Priority   int
}

// Supported format definitions
var (
    FormatJSON  = Format{Name: "json", Extensions: []string{".json"}, Priority: 10}
    FormatYAML  = Format{Name: "yaml", Extensions: []string{".yaml", ".yml"}, Priority: 10}
    FormatTOML  = Format{Name: "toml", Extensions: []string{".toml"}, Priority: 10}
    FormatINI   = Format{Name: "ini", Extensions: []string{".ini", ".cfg", ".conf"}, Priority: 5}
    FormatEnv   = Format{Name: "env", Extensions: []string{".env"}, Priority: 5}
    FormatXML   = Format{Name: "xml", Extensions: []string{".xml"}, Priority: 5}
    FormatPlist = Format{Name: "plist", Extensions: []string{".plist"}, Priority: 4} // macOS/iOS
    FormatPlain = Format{Name: "plain", Extensions: nil, Priority: 0}

    // AllFormats is the complete list of supported formats.
    AllFormats = []Format{FormatJSON, FormatYAML, FormatTOML, FormatINI, FormatEnv, FormatXML, FormatPlist}
)

// DetectFormat detects the config file format from a file path.
// First tries by file extension, then falls back to known filenames.
func DetectFormat(path string) Format {
    ext := strings.ToLower(filepath.Ext(path))
    for _, f := range AllFormats {
        for _, e := range f.Extensions {
            if e == ext {
                return f
            }
        }
    }

    // Try by base filename (e.g., ".env" has no extension)
    base := strings.ToLower(filepath.Base(path))
    if base == ".env" {
        return FormatEnv
    }

    return FormatPlain
}

// DetectFormatFromContent detects format by sniffing file content.
// Uses simple heuristics on the first line/content prefix.
func DetectFormatFromContent(content string) Format {
    content = strings.TrimSpace(content)
    if len(content) == 0 {
        return FormatPlain
    }

    firstLine := content
    if idx := strings.Index(content, "\n"); idx >= 0 {
        firstLine = content[:idx]
    }
    firstLine = strings.TrimSpace(firstLine)

    // XML: starts with <
    if strings.HasPrefix(firstLine, "<") {
        return FormatXML
    }

    // TOML / INI: [section] — must check before JSON array [
    if strings.HasPrefix(firstLine, "[") && strings.Contains(firstLine, "]") {
        // Check if it looks like a JSON array (content after [ is a quote or digit)
        contentAfter := firstLine[strings.Index(firstLine, "[")+1:]
        trimmed := strings.TrimSpace(contentAfter)
        if len(trimmed) > 0 && (trimmed[0] == '"' || (trimmed[0] >= '0' && trimmed[0] <= '9') || trimmed[0] == '{' || trimmed[0] == '[') {
            // JSON array
            return FormatJSON
        }
        // TOML section header
        return FormatTOML
    }

    // JSON: starts with { or [ (but not a TOML section)
    if strings.HasPrefix(firstLine, "{") || strings.HasPrefix(firstLine, "[") {
        return FormatJSON
    }

    // YAML: starts with ---
    if strings.HasPrefix(firstLine, "---") {
        return FormatYAML
    }

    // .env / INI: key=value
    if strings.Contains(firstLine, "=") {
        return FormatEnv
    }

    // YAML: key: value
    if strings.Contains(firstLine, ": ") {
        return FormatYAML
    }

    return FormatPlain
}