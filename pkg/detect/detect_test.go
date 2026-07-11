package detect

import (
	"testing"
)

func TestDetectOS_Linux(t *testing.T) {
	// Can't change runtime.GOOS at runtime, but we can verify the structure
	info := DetectOS()
	if info.GOOS == "" {
		t.Fatal("expected non-empty GOOS")
	}
	if info.PrettyName == "" {
		t.Fatal("expected non-empty PrettyName")
	}
	if info.NativeBinary == "" {
		t.Fatal("expected non-empty NativeBinary")
	}
	if len(info.ConfigDirs) == 0 {
		t.Fatal("expected at least one ConfigDir")
	}
	if len(info.ConfigFiles) == 0 {
		t.Fatal("expected at least one ConfigFile")
	}
}

func TestIsMobileOS(t *testing.T) {
	// On the current platform (linux), should return false
	result := IsMobileOS()
	t.Logf("IsMobileOS() = %v (current OS)", result)
}

func TestDetectFormat_ByExtension(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"config.json", "json"},
		{"config.yaml", "yaml"},
		{"config.yml", "yaml"},
		{"config.toml", "toml"},
		{"config.ini", "ini"},
		{"config.cfg", "ini"},
		{"config.conf", "ini"},
		{".env", "env"},
		{"config.env", "env"},
		{"config.xml", "xml"},
		{"config.plist", "plist"},
		{"config.unknown", "plain"},
		{"config", "plain"},
	}

	for _, tc := range tests {
		f := DetectFormat(tc.path)
		if f.Name != tc.expected {
			t.Errorf("DetectFormat(%q) = %s, want %s", tc.path, f.Name, tc.expected)
		}
	}
}

func TestDetectFormatFromContent(t *testing.T) {
	tests := []struct {
		content  string
		expected string
	}{
		{`{"key": "value"}`, "json"},
		{`[{"key": 1}]`, "json"},
		{`key: value`, "yaml"},
		{`---`, "yaml"},
		{`<root><item>val</item></root>`, "xml"},
		{`<config>`, "xml"},
		{`[server]`, "toml"},
		{`[server]port = 8080`, "toml"},
		{`KEY=value`, "env"},
		{`KEY="value with spaces"`, "env"},
		{`random text`, "plain"},
		{"", "plain"},
		{"   ", "plain"},
	}

	for _, tc := range tests {
		f := DetectFormatFromContent(tc.content)
		if f.Name != tc.expected {
			t.Errorf("DetectFormatFromContent(%q) = %s, want %s", tc.content, f.Name, tc.expected)
		}
	}
}

func TestFormatList(t *testing.T) {
	if len(AllFormats) == 0 {
		t.Fatal("AllFormats should not be empty")
	}

	// Verify all formats have names
	for _, f := range AllFormats {
		if f.Name == "" {
			t.Error("Format has empty name")
		}
	}
}

func TestOSInfoFields(t *testing.T) {
	info := DetectOS()

	// Verify that the struct fields are accessible
	_ = info.GOOS
	_ = info.PrettyName
	_ = info.NativeBinary
	_ = info.IsMobile
	_ = info.ConfigDirs
	_ = info.ConfigFiles

	t.Logf("OS: %s (%s), Mobile: %v", info.GOOS, info.PrettyName, info.IsMobile)
	t.Logf("Config dirs: %v", info.ConfigDirs)
	t.Logf("Config files: %v", info.ConfigFiles)
}