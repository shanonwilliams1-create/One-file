// Package audit provides audit logging for OmniConfig config changes.
// Audit logs are stored in JSON Lines format at ~/.config/omniconfig/audit.log.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp  string `json:"timestamp"`
	User       string `json:"user"`
	Business   string `json:"business"`
	File       string `json:"file"`
	Key        string `json:"key"`
	OldValue   string `json:"old_value"`
	NewValue   string `json:"new_value"`
	BackupPath string `json:"backup_path,omitempty"`
}

// LogDir returns the OmniConfig log directory.
func LogDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot find home directory: %w", err)
	}
	dir := filepath.Join(home, ".config", "omniconfig")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("cannot create log directory: %w", err)
	}
	return dir, nil
}

// LogPath returns the path to the audit log file.
func LogPath() (string, error) {
	dir, err := LogDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "audit.log"), nil
}

// Log appends an audit entry to the log file.
func Log(entry Entry) error {
	entry.Timestamp = time.Now().UTC().Format(time.RFC3339)

	path, err := LogPath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal audit entry: %w", err)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open audit log: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write audit log: %w", err)
	}

	return nil
}

// ReadAll reads all audit log entries, most recent first.
// Supports optional filtering by file path and date.
func ReadAll(filterFile, filterSince string) ([]Entry, error) {
	path, err := LogPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Entry{}, nil
		}
		return nil, fmt.Errorf("failed to read audit log: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	var entries []Entry

	for _, line := range lines {
		if line == "" {
			continue
		}
		var e Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			continue // skip corrupt lines
		}

		// Apply filters
		if filterFile != "" && !strings.Contains(e.File, filterFile) {
			continue
		}
		if filterSince != "" {
			if e.Timestamp < filterSince {
				continue
			}
		}

		entries = append(entries, e)
	}

	// Most recent first
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp > entries[j].Timestamp
	})

	return entries, nil
}

// Display formats an entry for terminal output.
func (e Entry) Display() string {
	s := fmt.Sprintf("[%s]\n", e.Timestamp)
	if e.User != "" {
		s += fmt.Sprintf("  User:     %s\n", e.User)
	}
	if e.Business != "" {
		s += fmt.Sprintf("  Business: %s\n", e.Business)
	}
	s += fmt.Sprintf("  File:     %s\n", e.File)
	s += fmt.Sprintf("  Key:      %s\n", e.Key)
	s += fmt.Sprintf("  Old:      %s\n", e.OldValue)
	s += fmt.Sprintf("  New:      %s\n", e.NewValue)
	if e.BackupPath != "" {
		s += fmt.Sprintf("  Backup:   %s\n", e.BackupPath)
	}
	return s
}