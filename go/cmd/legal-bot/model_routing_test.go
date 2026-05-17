package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jdonohoo/legal-bot/go/internal/config"
)

func TestSelectStepRouteOverrideBeatsProfile(t *testing.T) {
	pathDir := t.TempDir()
	writeFakeCLI(t, pathDir, "codex.cmd")
	writeFakeCLI(t, pathDir, "gemini.cmd")
	t.Setenv("PATH", pathDir)

	cfg := &config.Config{
		LLMMode: "single_llm",
		LLMModes: map[string]config.LLMModeConfig{
			"single_llm": {
				OverrideLLM: "codex:gpt-5.4:high",
			},
		},
		ModelProfiles: map[string]config.ModelProfileConfig{
			"review_reasoning": {
				Primary: config.ModelTargetSpec{Engine: "gemini", Model: "pro", Set: true},
			},
		},
	}

	selection, err := selectStepRoute(cfg, "review_reasoning", "claude")
	if err != nil {
		t.Fatal(err)
	}
	if selection.Target.Engine != "codex" || selection.Target.Model != "gpt-5.4" || selection.Target.Effort != "high" {
		t.Fatalf("unexpected route selection: %#v", selection.Target)
	}
}

func writeFakeCLI(t *testing.T, dir string, name string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte("@echo off\r\n"), 0644); err != nil {
		t.Fatal(err)
	}
}
