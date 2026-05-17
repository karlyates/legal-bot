package persona

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFile(t *testing.T) {
	// Create a temp agent file
	dir := t.TempDir()
	content := `---
name: mighty
description: Lead Analyst / Codex Analyst - Raw computational power.
model: opus
model_profile: review_reasoning
color: blue
---

You are the Lead Analyst. You wield the power of Codex.

PERSONALITY:
- Powerful and thorough
`
	path := filepath.Join(dir, "mighty.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if p.Name != "mighty" {
		t.Errorf("name: got %q, want %q", p.Name, "mighty")
	}
	if p.Model != "opus" {
		t.Errorf("model: got %q, want %q", p.Model, "opus")
	}
	if p.ModelProfile != "review_reasoning" {
		t.Errorf("model_profile: got %q, want %q", p.ModelProfile, "review_reasoning")
	}
	if p.Color != "blue" {
		t.Errorf("color: got %q, want %q", p.Color, "blue")
	}
	if p.Body == "" {
		t.Error("body should not be empty")
	}
}

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	content := `---
name: yolo
description: Stress Tester - No guardrails.
model: sonnet
---

Full send.
`
	if err := os.WriteFile(filepath.Join(dir, "yolo.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p, err := Load(dir, "yolo")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "yolo" {
		t.Errorf("name: got %q, want %q", p.Name, "yolo")
	}
}

func TestModelToLLM(t *testing.T) {
	tests := []struct {
		model string
		want  string
	}{
		{"opus", "claude"},
		{"sonnet", "claude"},
		{"haiku", "claude"},
		{"unknown", "claude"},
	}
	for _, tt := range tests {
		got := ModelToLLM(tt.model)
		if got != tt.want {
			t.Errorf("ModelToLLM(%q) = %q, want %q", tt.model, got, tt.want)
		}
	}
}

func TestResolveNameLegacyLegalAgents(t *testing.T) {
	tests := map[string]string{
		"intake-mapper":              "litigation-paralegal",
		"document-card-generator":    "source-document-analyst",
		"fact-extractor":             "atomic-fact-extractor",
		"timeline-builder":           "chronology-clerk",
		"open-questions-generator":   "attorney-prep-questioner",
		"fact-auditor":               "trial-fact-checker",
		"legal-sufficiency-reviewer": "family-law-attorney-reviewer",
		"court-reader":               "neutral-court-reader",
		"attack-surface-reviewer":    "opposing-counsel",
		"relief-alignment-reviewer":  "relief-and-order-alignment-counsel",
		"preservation-editor":        "legal-writing-preservation-editor",
		"bulldog-advocate":           "strategic-options-architect",
		"final-synthesizer":          "managing-partner-final-synthesizer",
	}
	for legacy, want := range tests {
		if got := ResolveName(legacy); got != want {
			t.Errorf("ResolveName(%q) = %q, want %q", legacy, got, want)
		}
	}
}

func TestShortDescription(t *testing.T) {
	tests := []struct {
		desc string
		want string
	}{
		{"Lead Analyst / Codex Analyst - Raw computational power. Comprehensive solutions.", "Raw computational power"},
		{"Stress Tester - No guardrails.", "No guardrails"},
		{"Simple description", "Simple description"},
	}
	for _, tt := range tests {
		got := ShortDescription(tt.desc)
		if got != tt.want {
			t.Errorf("ShortDescription(%q) = %q, want %q", tt.desc, got, tt.want)
		}
	}
}
