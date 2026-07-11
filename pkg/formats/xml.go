package formats

import (
	"encoding/xml"
	"fmt"
	"strings"
)

func init() {
	Register(xmlHandler{})
}

type xmlHandler struct{}

func (x xmlHandler) Name() string         { return "xml" }
func (x xmlHandler) Extensions() []string { return []string{".xml"} }

// xmlDocument is a simple wrapper for XML config files
type xmlDocument struct {
	XMLName xml.Name     `xml:"config"`
	Entries []xmlEntry   `xml:"entry,omitempty"`
	Sections []xmlSection `xml:"section,omitempty"`
}

type xmlEntry struct {
	Key   string `xml:"key,attr"`
	Value string `xml:",chardata"`
}

type xmlSection struct {
	Name    string    `xml:"name,attr"`
	Entries []xmlEntry `xml:"entry,omitempty"`
}

func (x xmlHandler) Read(data []byte, key string) (string, error) {
	// Try parsing as simple key-value format first
	var doc xmlDocument
	if err := xml.Unmarshal(data, &doc); err != nil {
		return "", fmt.Errorf("xml: failed to parse: %w", err)
	}

	parts := splitKey(key)

	// If single part, look in root entries
	if len(parts) == 1 {
		for _, e := range doc.Entries {
			if e.Key == parts[0] {
				return e.Value, nil
			}
		}
		return "", fmt.Errorf("xml: key %q not found", key)
	}

	// If two parts, look in section.key
	if len(parts) == 2 {
		sectionName := parts[0]
		entryKey := parts[1]
		for _, s := range doc.Sections {
			if s.Name == sectionName {
				for _, e := range s.Entries {
					if e.Key == entryKey {
						return e.Value, nil
					}
				}
			}
		}
		return "", fmt.Errorf("xml: key %q not found", key)
	}

	return "", fmt.Errorf("xml: key %q not found (max depth: section.key)", key)
}

func (x xmlHandler) Write(data []byte, key, value string) ([]byte, error) {
	var doc xmlDocument
	if len(strings.TrimSpace(string(data))) > 0 {
		if err := xml.Unmarshal(data, &doc); err != nil {
			// Start fresh
			doc = xmlDocument{}
		}
	}

	parts := splitKey(key)

	if len(parts) == 1 {
		// Update or add root entry
		found := false
		for i, e := range doc.Entries {
			if e.Key == parts[0] {
				doc.Entries[i].Value = value
				found = true
				break
			}
		}
		if !found {
			doc.Entries = append(doc.Entries, xmlEntry{Key: parts[0], Value: value})
		}
	} else if len(parts) == 2 {
		sectionName := parts[0]
		entryKey := parts[1]

		// Find or create section
		var section *xmlSection
		for i := range doc.Sections {
			if doc.Sections[i].Name == sectionName {
				section = &doc.Sections[i]
				break
			}
		}
		if section == nil {
			doc.Sections = append(doc.Sections, xmlSection{Name: sectionName})
			section = &doc.Sections[len(doc.Sections)-1]
		}

		// Update or add entry in section
		found := false
		for i, e := range section.Entries {
			if e.Key == entryKey {
				section.Entries[i].Value = value
				found = true
				break
			}
		}
		if !found {
			section.Entries = append(section.Entries, xmlEntry{Key: entryKey, Value: value})
		}
	} else {
		return nil, fmt.Errorf("xml: key %q not supported (max depth: section.key)", key)
	}

	result, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("xml: failed to marshal: %w", err)
	}
	result = append([]byte(xml.Header), result...)
	result = append(result, '\n')
	return result, nil
}