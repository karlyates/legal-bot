package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jdonohoo/legal-bot/go/internal/persona"
)

func TestLoadFile(t *testing.T) {
	dir := t.TempDir()
	configJSON := `{
		"version": "1.6.0",
		"timeout_seconds": 600,
		"max_retries": 2,
		"pipeline_mode": "default",
		"discovery_pipelines": {
			"default": [
				{
					"step": 1,
					"name": "Test Step",
					"persona": "mighty",
					"llm": "codex",
					"context_mode": "prompt_only",
					"prompt_prefix": "Test prompt"
				}
			]
		},
		"vernhole": {
			"default_council": "hammers",
			"min": 3
		}
	}`
	path := filepath.Join(dir, "config.json")
	os.WriteFile(path, []byte(configJSON), 0644)

	cfg, err := loadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.TimeoutSeconds != 600 {
		t.Errorf("timeout: got %d, want 600", cfg.TimeoutSeconds)
	}
	if cfg.MaxRetries != 2 {
		t.Errorf("max_retries: got %d, want 2", cfg.MaxRetries)
	}
	if cfg.VernHole.DefaultCouncil != "hammers" {
		t.Errorf("council: got %q, want %q", cfg.VernHole.DefaultCouncil, "hammers")
	}

	steps := cfg.GetPipeline("default")
	if len(steps) != 1 {
		t.Fatalf("steps: got %d, want 1", len(steps))
	}
	if steps[0].Name != "Test Step" {
		t.Errorf("step name: got %q, want %q", steps[0].Name, "Test Step")
	}
}

func TestLegacyConfig(t *testing.T) {
	dir := t.TempDir()
	configJSON := `{
		"version": "1.0.0",
		"discovery_pipeline": [
			{
				"step": 1,
				"name": "Legacy Step",
				"persona": "mighty",
				"llm": "claude",
				"context_mode": "prompt_only",
				"prompt_prefix": "Legacy prompt"
			}
		]
	}`
	path := filepath.Join(dir, "config.json")
	os.WriteFile(path, []byte(configJSON), 0644)

	cfg, err := loadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	steps := cfg.GetPipeline("default")
	if len(steps) != 1 {
		t.Fatalf("steps: got %d, want 1", len(steps))
	}
	if steps[0].Name != "Legacy Step" {
		t.Errorf("step name: got %q, want %q", steps[0].Name, "Legacy Step")
	}
}

func TestHardcodedDefaults(t *testing.T) {
	cfg := hardcodedDefaults()
	if cfg.TimeoutSeconds != 1200 {
		t.Errorf("timeout: got %d, want 1200", cfg.TimeoutSeconds)
	}

	defaultSteps := cfg.GetPipeline("default")
	if len(defaultSteps) != 5 {
		t.Errorf("default steps: got %d, want 5", len(defaultSteps))
	}

	expandedSteps := cfg.GetPipeline("expanded")
	if len(expandedSteps) != 7 {
		t.Errorf("expanded steps: got %d, want 7", len(expandedSteps))
	}
}

func TestGetPipelineFallback(t *testing.T) {
	cfg := &Config{
		Pipelines: map[string][]PipelineStep{
			"default": {{Step: 1, Name: "Only Default"}},
		},
	}

	// Request non-existent pipeline, should fall back to default
	steps := cfg.GetPipeline("nonexistent")
	if len(steps) != 1 || steps[0].Name != "Only Default" {
		t.Errorf("fallback failed: got %v", steps)
	}
}

func TestLoadWithProjectRoot(t *testing.T) {
	// Load should work with the actual project config.default.json
	cfg := Load(filepath.Join("..", "..", ".."))
	if cfg == nil {
		t.Fatal("Load returned nil")
	}
	defaultSteps := cfg.GetPipeline("default")
	if len(defaultSteps) < 3 {
		t.Errorf("expected at least 3 default steps, got %d", len(defaultSteps))
	}
}

func TestDefaultLegalPipelinesResolveAgentsAndProfiles(t *testing.T) {
	repoRoot := filepath.Join("..", "..", "..")
	cfg, err := loadFile(filepath.Join(repoRoot, "config.default.json"))
	if err != nil {
		t.Fatal(err)
	}

	agentsDir := filepath.Join(repoRoot, "agents")
	for pipelineName, steps := range cfg.LegalPipelines {
		if len(steps) == 0 {
			t.Errorf("legal pipeline %q has no steps", pipelineName)
			continue
		}
		for _, step := range steps {
			personaName := persona.ResolveName(step.Persona)
			if _, err := os.Stat(filepath.Join(agentsDir, personaName+".md")); err != nil {
				t.Errorf("pipeline %q step %d references missing agent %q", pipelineName, step.Step, step.Persona)
			}
			if step.ModelProfile != "" {
				if _, ok := cfg.ModelProfiles[step.ModelProfile]; !ok {
					t.Errorf("pipeline %q step %d references missing model profile %q", pipelineName, step.Step, step.ModelProfile)
				}
			}
		}
	}
}

func TestDefaultReviewRoutesExist(t *testing.T) {
	repoRoot := filepath.Join("..", "..", "..")
	cfg, err := loadFile(filepath.Join(repoRoot, "config.default.json"))
	if err != nil {
		t.Fatal(err)
	}

	for _, route := range []string{
		"review_quick",
		"review_standard",
		"review_deep",
		"review_legal_research",
		"review_strategy",
		"review_document_reply",
		"review_document_declaration",
		"review_document_proposed_order",
		"review_document_coparenting_communication",
	} {
		if steps := cfg.GetLegalPipeline(route); len(steps) == 0 {
			t.Errorf("expected legal review route %q to resolve steps", route)
		}
	}
}

func TestDefaultWorkflowRoutesExist(t *testing.T) {
	repoRoot := filepath.Join("..", "..", "..")
	cfg, err := loadFile(filepath.Join(repoRoot, "config.default.json"))
	if err != nil {
		t.Fatal(err)
	}

	for _, route := range []string{
		"triage_new_situation",
		"brief_counsel",
		"triage_special_master",
		"draft_external_message",
		"build_evidence_packet",
		"review_pattern",
		"fact_lock",
		"authority_check",
		"prep_hearing_or_call",
		"journal_entry",
	} {
		if steps := cfg.GetLegalPipeline(route); len(steps) == 0 {
			t.Errorf("expected legal workflow route %q to resolve steps", route)
		}
	}
}

func TestActiveLegalAgentsExistAndMatchFrontmatter(t *testing.T) {
	repoRoot := filepath.Join("..", "..", "..")
	agentsDir := filepath.Join(repoRoot, "agents")

	expectedAgents := []string{
		"atomic-fact-extractor",
		"attorney-prep-questioner",
		"child-best-interests-family-dynamics-reviewer",
		"chronology-clerk",
		"family-law-attorney-reviewer",
		"financial-support-reviewer",
		"legal-authority-scholar",
		"legal-writing-preservation-editor",
		"litigation-paralegal",
		"managing-partner-final-synthesizer",
		"neutral-court-reader",
		"opposing-counsel",
		"practical-resolution-reviewer",
		"relief-and-order-alignment-counsel",
		"source-document-analyst",
		"strategic-options-architect",
		"trial-fact-checker",
	}

	for _, name := range expectedAgents {
		path := filepath.Join(agentsDir, name+".md")
		p, err := persona.LoadFile(path)
		if err != nil {
			t.Fatalf("failed to load %s: %v", name, err)
		}
		if p.Name != name {
			t.Errorf("frontmatter name for %s = %q", name, p.Name)
		}
		if p.ModelProfile == "" {
			t.Errorf("%s missing model_profile", name)
		}
	}

	if _, err := os.Stat(filepath.Join(agentsDir, "bulldog-advocate.md")); err == nil {
		t.Error("bulldog-advocate.md should not exist as an active agent file")
	}
}

func TestParseModelTargetSpecCompactForms(t *testing.T) {
	tests := []struct {
		input  string
		engine string
		model  string
		effort string
	}{
		{"claude", "claude", "", ""},
		{"codex:gpt-5.4:medium", "codex", "gpt-5.4", "medium"},
		{"codex:gpt-5.5:high", "codex", "gpt-5.5", "high"},
		{"openai:gpt-5.5", "codex", "gpt-5.5", ""},
		{"claude:sonnet:high", "claude", "sonnet", "high"},
		{"claude:opus:xhigh", "claude", "opus", "xhigh"},
		{"gemini:pro", "gemini", "pro", ""},
		{"gemini:flash", "gemini", "flash", ""},
		{"copilot:auto", "copilot", "auto", ""},
	}

	for _, tt := range tests {
		got, err := ParseModelTargetSpec(tt.input)
		if err != nil {
			t.Fatalf("ParseModelTargetSpec(%q): %v", tt.input, err)
		}
		if got.Engine != tt.engine || got.Model != tt.model || got.Effort != tt.effort {
			t.Fatalf("ParseModelTargetSpec(%q) = %#v", tt.input, got)
		}
	}
}

func TestObjectStyleModelProfileParsing(t *testing.T) {
	dir := t.TempDir()
	configJSON := `{
		"model_profiles": {
			"review_reasoning": {
				"primary": {
					"engine": "codex",
					"model": "gpt-5.4",
					"effort": "medium"
				},
				"fallbacks": [
					{
						"engine": "claude",
						"model": "sonnet",
						"effort": "high"
					}
				],
				"fallback": "gemini:pro"
			}
		}
	}`
	path := filepath.Join(dir, "config.json")
	os.WriteFile(path, []byte(configJSON), 0644)

	cfg, err := loadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	targets, warnings := cfg.ResolveModelProfileTargets("review_reasoning", "claude")
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %v", warnings)
	}
	if len(targets) != 4 {
		t.Fatalf("expected 4 targets, got %d", len(targets))
	}
	if targets[0].Engine != "codex" || targets[0].Model != "gpt-5.4" || targets[0].Effort != "medium" {
		t.Fatalf("unexpected primary target: %#v", targets[0])
	}
	if targets[1].Engine != "claude" || targets[1].Model != "sonnet" || targets[1].Effort != "high" {
		t.Fatalf("unexpected fallback[0]: %#v", targets[1])
	}
	if targets[2].Engine != "gemini" || targets[2].Model != "pro" {
		t.Fatalf("unexpected fallback[1]: %#v", targets[2])
	}
	if targets[3].Engine != "claude" || targets[3].Model != "" {
		t.Fatalf("unexpected step fallback target: %#v", targets[3])
	}
}

func TestResolveModelProfileTargetsBackwardCompatible(t *testing.T) {
	cfg := &Config{
		ModelProfiles: map[string]ModelProfileConfig{
			"review_reasoning": {
				Preferred: ModelTargetSpec{Engine: "claude", Set: true},
				Fallback:  ModelTargetSpec{Engine: "gemini", Set: true},
			},
			"legacy_openai": {
				Preferred: ModelTargetSpec{Engine: "codex", Model: "gpt-5.5", LegacyOpenAI: true, Raw: "openai:gpt-5.5", Set: true},
				Fallback:  ModelTargetSpec{Engine: "claude", Set: true},
			},
			"bad_engine": {
				Primary: ModelTargetSpec{Engine: "wat", Model: "mystery", Set: true},
			},
		},
	}

	engine := cfg.ResolveModelProfile("review_reasoning", "copilot")
	if engine != "claude" {
		t.Fatalf("ResolveModelProfile returned %q", engine)
	}

	targets, _ := cfg.ResolveModelProfileTargets("legacy_openai", "gemini")
	if len(targets) < 1 {
		t.Fatal("expected at least one target")
	}
	if targets[0].Engine != "codex" || targets[0].Model != "gpt-5.5" || targets[0].Effort != "medium" {
		t.Fatalf("legacy openai target resolved incorrectly: %#v", targets[0])
	}

	_, warnings := cfg.ResolveModelProfileTargets("bad_engine", "claude")
	joined := strings.Join(warnings, " ")
	if !strings.Contains(joined, "unknown engine") {
		t.Fatalf("expected unknown engine warning, got %v", warnings)
	}
}

func TestGetOverrideTargetParsesCompactSpec(t *testing.T) {
	cfg := &Config{
		LLMMode: "single_llm",
		LLMModes: map[string]LLMModeConfig{
			"single_llm": {
				OverrideLLM: "codex:gpt-5.4:high",
			},
		},
	}

	target, warnings, err := cfg.GetOverrideTarget()
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if target == nil || target.Engine != "codex" || target.Model != "gpt-5.4" || target.Effort != "high" {
		t.Fatalf("unexpected override target: %#v", target)
	}
	if got := cfg.GetOverrideLLM(); got != "codex" {
		t.Fatalf("GetOverrideLLM() = %q", got)
	}
}

func TestSelectRunnableTargetPrefersFirstInstalledEngine(t *testing.T) {
	pathDir := t.TempDir()
	writeFakeCommand(t, pathDir, "gemini.cmd")
	writeFakeCommand(t, pathDir, "codex.cmd")
	t.Setenv("PATH", pathDir)

	targets := []ResolvedModelTarget{
		{Engine: "gemini", Model: "pro"},
		{Engine: "codex", Model: "gpt-5.4", Effort: "medium"},
	}

	selected, warnings, err := SelectRunnableTarget(targets)
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if selected.Engine != "gemini" || selected.Model != "pro" {
		t.Fatalf("unexpected selected target: %#v", selected)
	}
}

func TestSelectRunnableTargetSkipsUnknownAndMissingEngines(t *testing.T) {
	pathDir := t.TempDir()
	writeFakeCommand(t, pathDir, "codex.cmd")
	t.Setenv("PATH", pathDir)

	targets := []ResolvedModelTarget{
		{Engine: "wat", Model: "mystery"},
		{Engine: "claude", Model: "sonnet"},
		{Engine: "codex", Model: "gpt-5.4"},
	}

	selected, warnings, err := SelectRunnableTarget(targets)
	if err != nil {
		t.Fatal(err)
	}
	if selected.Engine != "codex" {
		t.Fatalf("expected codex, got %#v", selected)
	}
	joined := strings.Join(warnings, " ")
	if !strings.Contains(joined, "unknown engine") || !strings.Contains(joined, "not found on PATH") {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
}

func writeFakeCommand(t *testing.T, dir string, name string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("@echo off\r\n"), 0644); err != nil {
		t.Fatal(err)
	}
}
