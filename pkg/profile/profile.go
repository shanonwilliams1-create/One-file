// Package profile manages the OmniConfig user profile.
// The profile is stored in ~/.config/omniconfig/profile.json.
package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Profile holds the user's information.
type Profile struct {
	Name     string `json:"name"`
	Business string `json:"business"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

// configDir returns the OmniConfig config directory (~/.config/omniconfig).
func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot find home directory: %w", err)
	}
	dir := filepath.Join(home, ".config", "omniconfig")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("cannot create config directory: %w", err)
	}
	return dir, nil
}

// ProfilePath returns the path to the profile file.
func ProfilePath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "profile.json"), nil
}

// Load reads the profile from disk.
// Returns an empty profile if the file doesn't exist or can't be read.
func Load() (*Profile, error) {
	path, err := ProfilePath()
	if err != nil {
		return &Profile{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Profile{}, nil
		}
		return &Profile{}, nil // silently return empty
	}

	var p Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return &Profile{}, nil
	}

	return &p, nil
}

// Save writes the profile to disk.
func Save(p *Profile) error {
	path, err := ProfilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profile: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write profile: %w", err)
	}

	return nil
}

// IsComplete returns true if the profile has at least a name.
func (p *Profile) IsComplete() bool {
	return p.Name != ""
}

// Display returns a formatted string of the profile.
func (p *Profile) Display() string {
	s := fmt.Sprintf("Name:     %s\n", p.Name)
	s += fmt.Sprintf("Business: %s\n", p.Business)
	if p.Email != "" {
		s += fmt.Sprintf("Email:    %s\n", p.Email)
	}
	if p.Phone != "" {
		s += fmt.Sprintf("Phone:    %s\n", p.Phone)
	}
	return s
}