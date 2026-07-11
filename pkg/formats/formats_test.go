package formats

import (
	"testing"
)

// --- Round-trip tests ---

func TestJSONRoundTrip(t *testing.T) {
	h := jsonHandler{}
	testRoundTrip(t, h, `{"server":{"port":8080,"host":"localhost"}}`, "server.port", "9090", "server.host", "example.com")
}

func TestYAMLRoundTrip(t *testing.T) {
	h := yamlHandler{}
	testRoundTrip(t, h, "server:\n  port: 8080\n  host: localhost\n", "server.port", "9090", "server.host", "example.com")
}

func TestTOMLRoundTrip(t *testing.T) {
	h := tomlHandler{}
	testRoundTrip(t, h, "[server]\nport = 8080\nhost = \"localhost\"\n", "server.port", "9090", "server.host", "example.com")
}

func TestINIRoundTrip(t *testing.T) {
	h := iniHandler{}
	input := "[server]\nport = 8080\nhost = localhost\n"
	data := []byte(input)

	val, err := h.Read(data, "server.port")
	if err != nil {
		t.Fatalf("INI read failed: %v", err)
	}
	if val != "8080" {
		t.Fatalf("expected 8080, got %s", val)
	}

	data, err = h.Write(data, "server.port", "9090")
	if err != nil {
		t.Fatalf("INI write failed: %v", err)
	}

	val, err = h.Read(data, "server.port")
	if err != nil {
		t.Fatalf("INI re-read failed: %v", err)
	}
	if val != "9090" {
		t.Fatalf("expected 9090, got %s", val)
	}
}

func TestEnvRoundTrip(t *testing.T) {
	h := envHandler{}
	input := "DATABASE_HOST=localhost\nDATABASE_PORT=5432\n"
	data := []byte(input)

	val, err := h.Read(data, "DATABASE_HOST")
	if err != nil {
		t.Fatalf("env read failed: %v", err)
	}
	if val != "localhost" {
		t.Fatalf("expected localhost, got %s", val)
	}

	data, err = h.Write(data, "DATABASE_HOST", "prod.example.com")
	if err != nil {
		t.Fatalf("env write failed: %v", err)
	}

	val, err = h.Read(data, "DATABASE_HOST")
	if err != nil {
		t.Fatalf("env re-read failed: %v", err)
	}
	if val != "prod.example.com" {
		t.Fatalf("expected prod.example.com, got %s", val)
	}
}

func TestXMLRoundTrip(t *testing.T) {
	h := xmlHandler{}
	input := `<config><entry key="host">localhost</entry></config>`
	data := []byte(input)

	val, err := h.Read(data, "host")
	if err != nil {
		t.Fatalf("xml read failed: %v", err)
	}
	if val != "localhost" {
		t.Fatalf("expected localhost, got %s", val)
	}

	data, err = h.Write(data, "host", "example.com")
	if err != nil {
		t.Fatalf("xml write failed: %v", err)
	}

	val, err = h.Read(data, "host")
	if err != nil {
		t.Fatalf("xml re-read failed: %v", err)
	}
	if val != "example.com" {
		t.Fatalf("expected example.com, got %s", val)
	}
}

// --- Edge case tests ---

func TestJSONEmptyFile(t *testing.T) {
	h := jsonHandler{}
	data, err := h.Write([]byte{}, "server.port", "8080")
	if err != nil {
		t.Fatalf("JSON write to empty file failed: %v", err)
	}

	val, err := h.Read(data, "server.port")
	if err != nil {
		t.Fatalf("JSON read from new file failed: %v", err)
	}
	if val != "8080" {
		t.Fatalf("expected 8080, got %s", val)
	}
}

func TestJSONMissingKey(t *testing.T) {
	h := jsonHandler{}
	_, err := h.Read([]byte(`{"a":1}`), "b")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestJSONNestedCreate(t *testing.T) {
	h := jsonHandler{}
	data, err := h.Write([]byte(`{"a":1}`), "b.c.d", "deep")
	if err != nil {
		t.Fatalf("JSON nested write failed: %v", err)
	}

	val, err := h.Read(data, "b.c.d")
	if err != nil {
		t.Fatalf("JSON nested read failed: %v", err)
	}
	if val != "deep" {
		t.Fatalf("expected deep, got %s", val)
	}
}

func TestYAMLEmptyFile(t *testing.T) {
	h := yamlHandler{}
	data, err := h.Write([]byte{}, "server.port", "8080")
	if err != nil {
		t.Fatalf("YAML write to empty file failed: %v", err)
	}

	val, err := h.Read(data, "server.port")
	if err != nil {
		t.Fatalf("YAML read from new file failed: %v", err)
	}
	if val != "8080" {
		t.Fatalf("expected 8080, got %s", val)
	}
}

func TestYAMLMissingKey(t *testing.T) {
	h := yamlHandler{}
	_, err := h.Read([]byte("a: 1\n"), "b")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestTOMLMissingKey(t *testing.T) {
	h := tomlHandler{}
	_, err := h.Read([]byte("[server]\nport=8080\n"), "server.host")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestINIDefaultSection(t *testing.T) {
	h := iniHandler{}
	input := "port = 8080\nhost = localhost\n"
	data := []byte(input)

	val, err := h.Read(data, "port")
	if err != nil {
		t.Fatalf("INI default section read failed: %v", err)
	}
	if val != "8080" {
		t.Fatalf("expected 8080, got %s", val)
	}
}

func TestINIMissingKey(t *testing.T) {
	h := iniHandler{}
	_, err := h.Read([]byte("[server]\nport=8080\n"), "server.host")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestEnvMissingKey(t *testing.T) {
	h := envHandler{}
	_, err := h.Read([]byte("A=1\n"), "B")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestEnvEmptyFile(t *testing.T) {
	h := envHandler{}
	data, err := h.Write([]byte{}, "NEW_KEY", "new_value")
	if err != nil {
		t.Fatalf("env write to empty file failed: %v", err)
	}

	val, err := h.Read(data, "NEW_KEY")
	if err != nil {
		t.Fatalf("env read from new file failed: %v", err)
	}
	if val != "new_value" {
		t.Fatalf("expected new_value, got %s", val)
	}
}

func TestXMLMissingKey(t *testing.T) {
	h := xmlHandler{}
	_, err := h.Read([]byte(`<config><entry key="host">localhost</entry></config>`), "missing")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestXMLSectionKey(t *testing.T) {
	h := xmlHandler{}
	input := `<config><section name="server"><entry key="port">8080</entry></section></config>`
	data := []byte(input)

	val, err := h.Read(data, "server.port")
	if err != nil {
		t.Fatalf("XML section read failed: %v", err)
	}
	if val != "8080" {
		t.Fatalf("expected 8080, got %s", val)
	}

	data, err = h.Write(data, "server.port", "9090")
	if err != nil {
		t.Fatalf("XML section write failed: %v", err)
	}

	val, err = h.Read(data, "server.port")
	if err != nil {
		t.Fatalf("XML section re-read failed: %v", err)
	}
	if val != "9090" {
		t.Fatalf("expected 9090, got %s", val)
	}
}

// --- GetByExtension tests ---

func TestGetByExtension(t *testing.T) {
	tests := []struct {
		ext      string
		expected string
	}{
		{".json", "json"},
		{".yaml", "yaml"},
		{".yml", "yaml"},
		{".toml", "toml"},
		{".ini", "ini"},
		{".cfg", "ini"},
		{".conf", "ini"},
		{".env", "env"},
		{".xml", "xml"},
		{".unknown", ""},
		{"", ""},
	}

	for _, tc := range tests {
		h := GetByExtension(tc.ext)
		if tc.expected == "" {
			if h != nil {
				t.Errorf("GetByExtension(%q) expected nil, got %s", tc.ext, h.Name())
			}
		} else {
			if h == nil {
				t.Errorf("GetByExtension(%q) expected %s, got nil", tc.ext, tc.expected)
			} else if h.Name() != tc.expected {
				t.Errorf("GetByExtension(%q) expected %s, got %s", tc.ext, tc.expected, h.Name())
			}
		}
	}
}

// --- Helper: generic round-trip test ---

func testRoundTrip(t *testing.T, h Handler, input, key1, val1, key2, val2 string) {
	t.Helper()
	data := []byte(input)

	// Read initial value
	val, err := h.Read(data, key1)
	if err != nil {
		t.Fatalf("%s read failed: %v", h.Name(), err)
	}
	t.Logf("%s read %s = %s", h.Name(), key1, val)

	// Write new value
	data, err = h.Write(data, key1, val1)
	if err != nil {
		t.Fatalf("%s write failed: %v", h.Name(), err)
	}

	// Read back
	val, err = h.Read(data, key1)
	if err != nil {
		t.Fatalf("%s re-read failed: %v", h.Name(), err)
	}
	if val != val1 {
		t.Fatalf("%s expected %s, got %s", h.Name(), val1, val)
	}

	// Write second key
	data, err = h.Write(data, key2, val2)
	if err != nil {
		t.Fatalf("%s write key2 failed: %v", h.Name(), err)
	}

	// Verify both keys
	val, err = h.Read(data, key1)
	if err != nil {
		t.Fatalf("%s re-read key1 failed: %v", h.Name(), err)
	}
	if val != val1 {
		t.Fatalf("%s expected %s, got %s", h.Name(), val1, val)
	}

	val, err = h.Read(data, key2)
	if err != nil {
		t.Fatalf("%s re-read key2 failed: %v", h.Name(), err)
	}
	if val != val2 {
		t.Fatalf("%s expected %s, got %s", h.Name(), val2, val)
	}
}