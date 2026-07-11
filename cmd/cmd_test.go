package cmd

import (
	"testing"
)

func TestCommandStructure(t *testing.T) {
	// Verify commands are registered
	cmds := rootCmd.Commands()
	cmdNames := make(map[string]bool)
	for _, c := range cmds {
		cmdNames[c.Name()] = true
	}

	expected := []string{"detect", "get", "set", "version"}
	for _, name := range expected {
		if !cmdNames[name] {
			t.Errorf("expected command %q to be registered", name)
		}
	}
}

func TestRootHasFlags(t *testing.T) {
	// Verify root flags exist
	configFlag := rootCmd.PersistentFlags().Lookup("config")
	if configFlag == nil {
		t.Error("expected --config flag")
	}

	listFlag := rootCmd.PersistentFlags().Lookup("list-platforms")
	if listFlag == nil {
		t.Error("expected --list-platforms flag")
	}
}

func TestGetCmdArgs(t *testing.T) {
	getCmd := NewGetCmd()
	if getCmd == nil {
		t.Fatal("NewGetCmd() returned nil")
	}
	if getCmd.Args == nil {
		t.Error("get command should have arg validation")
	}
}

func TestSetCmdArgs(t *testing.T) {
	setCmd := NewSetCmd()
	if setCmd == nil {
		t.Fatal("NewSetCmd() returned nil")
	}
	if setCmd.Args == nil {
		t.Error("set command should have arg validation")
	}
}

func TestDetectCmd(t *testing.T) {
	detectCmd := NewDetectCmd()
	if detectCmd == nil {
		t.Fatal("NewDetectCmd() returned nil")
	}
	if detectCmd.Use != "detect" {
		t.Errorf("expected use 'detect', got %q", detectCmd.Use)
	}
}

func TestVersionCmd(t *testing.T) {
	versionCmd := NewVersionCmd()
	if versionCmd == nil {
		t.Fatal("NewVersionCmd() returned nil")
	}
	if versionCmd.Use != "version" {
		t.Errorf("expected use 'version', got %q", versionCmd.Use)
	}
}

func TestExpandPath(t *testing.T) {
	result := expandPath("~/.config")
	if result == "" {
		t.Error("expandPath returned empty string")
	}
	if result == "~/.config" {
		t.Error("expandPath should have expanded ~")
	}
}

func TestDetectHelp(t *testing.T) {
	detectCmd := NewDetectCmd()
	if detectCmd.Long == "" {
		t.Error("detect command should have a long description")
	}
	if detectCmd.Short == "" {
		t.Error("detect command should have a short description")
	}
}

func TestGetHelp(t *testing.T) {
	getCmd := NewGetCmd()
	if getCmd.Long == "" {
		t.Error("get command should have a long description")
	}
	if getCmd.Short == "" {
		t.Error("get command should have a short description")
	}
}

func TestSetHelp(t *testing.T) {
	setCmd := NewSetCmd()
	if setCmd.Long == "" {
		t.Error("set command should have a long description")
	}
	if setCmd.Short == "" {
		t.Error("set command should have a short description")
	}
}