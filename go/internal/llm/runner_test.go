package llm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestResolveLLM(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"claude", "claude"},
		{"c", "claude"},
		{"Claude", "claude"},
	}
	for _, tt := range tests {
		got := resolveLLM(tt.input)
		if got != tt.want {
			t.Errorf("resolveLLM(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestBuildEngineArgs(t *testing.T) {
	codexArgs := buildCodexArgs("out.md", "c:\\repo", RunOptions{Model: "gpt-5.4", Effort: "high"})
	assertContainsSequence(t, codexArgs, []string{"--model", "gpt-5.4"})
	assertContainsSequence(t, codexArgs, []string{"-c", "model_reasoning_effort=high"})
	assertContainsSequence(t, codexArgs, []string{"-o", "out.md"})
	if containsArg(codexArgs, "--dangerously-bypass-approvals-and-sandbox") {
		t.Fatal("codex args should not include bypass flag by default")
	}
	for _, arg := range codexArgs {
		if arg == "prompt" {
			t.Fatal("codex prompt should be sent via stdin, not argv")
		}
	}

	claudeArgs := buildClaudeArgs("prompt", RunOptions{Model: "sonnet", Effort: "xhigh"})
	assertContainsSequence(t, claudeArgs, []string{"--model", "sonnet"})
	assertContainsSequence(t, claudeArgs, []string{"--effort", "xhigh"})

	geminiArgs := buildGeminiArgs("prompt", RunOptions{Model: "flash", Effort: "high"})
	assertContainsSequence(t, geminiArgs, []string{"--model", "flash"})
	if containsArg(geminiArgs, "--effort") {
		t.Fatal("gemini args should not include effort")
	}

	copilotArgs := buildCopilotArgs("prompt", RunOptions{Model: "auto", Effort: "high"})
	assertContainsSequence(t, copilotArgs, []string{"--model", "auto"})
	if containsArg(copilotArgs, "--effort") {
		t.Fatal("copilot args should not include effort")
	}
}

func TestFormatRunDiagnosticHidesPromptAndIncludesExitContext(t *testing.T) {
	diag := formatRunDiagnostic(
		"claude",
		RunOptions{Model: "sonnet", Effort: "high", Prompt: "secret prompt"},
		"claude",
		[]string{"--model", "sonnet", "--effort", "high", "-p", "secret prompt"},
		17,
		&Result{ExitCode: 1, Stderr: "fatal: backend unavailable"},
		fmt.Errorf("exit status 1"),
	)

	if !strings.Contains(diag, "engine=claude") {
		t.Fatalf("diagnostic missing engine: %s", diag)
	}
	if !strings.Contains(diag, "model=sonnet") || !strings.Contains(diag, "effort=high") {
		t.Fatalf("diagnostic missing model/effort: %s", diag)
	}
	if !strings.Contains(diag, "exit_code=1") || !strings.Contains(diag, "stdout_bytes=17") {
		t.Fatalf("diagnostic missing exit/stdout context: %s", diag)
	}
	if strings.Contains(diag, "secret prompt") {
		t.Fatalf("diagnostic leaked prompt: %s", diag)
	}
	if !strings.Contains(diag, "<prompt>") {
		t.Fatalf("diagnostic should include sanitized prompt marker: %s", diag)
	}
	if !strings.Contains(diag, "stderr_tail=") {
		t.Fatalf("diagnostic missing stderr tail: %s", diag)
	}
}

func TestResolveCodexExecutablePrefersCmdOnWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-specific PATH resolution")
	}

	pathDir := t.TempDir()
	writeFakeCLI(t, pathDir, "codex.cmd")
	writeFakeCLI(t, pathDir, "codex.exe")
	t.Setenv("PATH", pathDir)

	got, err := resolveCodexExecutable()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(strings.ToLower(got), "codex.cmd") {
		t.Fatalf("resolveCodexExecutable() = %q, want codex.cmd", got)
	}
}

func TestLoadPersonaContext(t *testing.T) {
	dir := t.TempDir()
	content := `---
name: mighty
description: Lead Analyst
model: opus
---

You are the Lead Analyst. You wield the power of Codex.
`
	os.WriteFile(filepath.Join(dir, "mighty.md"), []byte(content), 0644)

	ctx := loadPersonaContext(dir, "mighty")
	if ctx == "" {
		t.Error("persona context should not be empty")
	}
	if !containsStr(ctx, "=== PERSONA ===") {
		t.Error("should contain PERSONA markers")
	}
	if !containsStr(ctx, "You are the Lead Analyst") {
		t.Error("should contain persona body")
	}
	// Should NOT contain frontmatter
	if containsStr(ctx, "model: opus") {
		t.Error("should not contain frontmatter")
	}
}

func TestLoadPersonaContextMissing(t *testing.T) {
	ctx := loadPersonaContext("/nonexistent", "missing")
	if ctx != "" {
		t.Error("missing persona should return empty string")
	}
}

func TestExitCodeFromErr(t *testing.T) {
	if exitCodeFromErr(nil) != 0 {
		t.Error("nil error should return 0")
	}
}

func TestLogRunWritesJSONL(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir)
	// Clear any LEGAL_BOT_LOG setting
	t.Setenv("LEGAL_BOT_LOG", "")

	// Create the config dir structure so configDir() resolves here
	logDir := filepath.Join(tmpDir, ".config", "legal-bot", "logs")

	opts := RunOptions{
		LLM:        "gemini",
		Prompt:     "Analyze this idea about building a distributed system",
		OutputFile: "output.md",
	}
	result := &Result{
		Output:   "Here is my analysis...",
		ExitCode: 0,
		TimedOut: false,
		LLMUsed:  "claude",
		Duration: 4500 * time.Millisecond,
	}

	logRun(opts, "gemini", result, nil, nil)

	logPath := filepath.Join(logDir, "legal-bot.log")
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 log line, got %d", len(lines))
	}

	var entry logEntry
	if err := json.Unmarshal([]byte(lines[0]), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if entry.LLMRequested != "gemini" {
		t.Errorf("llm_requested = %q, want %q", entry.LLMRequested, "gemini")
	}
	if entry.LLMUsed != "claude" {
		t.Errorf("llm_used = %q, want %q", entry.LLMUsed, "claude")
	}
	if entry.ExitCode != 0 {
		t.Errorf("exit_code = %d, want 0", entry.ExitCode)
	}
	if entry.TimedOut {
		t.Error("timed_out should be false")
	}
	if entry.DurationMs != 4500 {
		t.Errorf("duration_ms = %d, want 4500", entry.DurationMs)
	}
	if entry.Error != "" {
		t.Errorf("error should be empty, got %q", entry.Error)
	}
	if entry.OutputFile != "output.md" {
		t.Errorf("output_file = %q, want %q", entry.OutputFile, "output.md")
	}
	if entry.OutputBytes != len("Here is my analysis...") {
		t.Errorf("output_bytes = %d, want %d", entry.OutputBytes, len("Here is my analysis..."))
	}
	if entry.PromptPreview != "Analyze this idea about building a distributed system" {
		t.Errorf("prompt_preview = %q, want full prompt", entry.PromptPreview)
	}
	if entry.Time == "" {
		t.Error("time should not be empty")
	}
}

func TestLogRunWithError(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir)
	t.Setenv("LEGAL_BOT_LOG", "")

	opts := RunOptions{
		LLM:    "codex",
		Prompt: "test prompt",
	}
	result := &Result{
		ExitCode: 1,
		LLMUsed:  "codex",
		Duration: 200 * time.Millisecond,
	}

	logRun(opts, "codex", result, fmt.Errorf("signal: killed"), nil)

	logPath := filepath.Join(tmpDir, ".config", "legal-bot", "logs", "legal-bot.log")
	data, _ := os.ReadFile(logPath)

	var entry logEntry
	json.Unmarshal([]byte(strings.TrimSpace(string(data))), &entry)

	if entry.Error != "signal: killed" {
		t.Errorf("error = %q, want %q", entry.Error, "signal: killed")
	}
	if entry.ExitCode != 1 {
		t.Errorf("exit_code = %d, want 1", entry.ExitCode)
	}
}

func TestLogRunDisabledByEnv(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir)
	t.Setenv("LEGAL_BOT_LOG", "0")

	opts := RunOptions{LLM: "claude", Prompt: "test"}
	result := &Result{LLMUsed: "claude", Duration: time.Second}

	logRun(opts, "claude", result, nil, nil)

	logPath := filepath.Join(tmpDir, ".config", "legal-bot", "logs", "legal-bot.log")
	if _, err := os.Stat(logPath); err == nil {
		t.Error("log file should not exist when LEGAL_BOT_LOG=0")
	}
}

func TestTruncatePrompt(t *testing.T) {
	short := "short prompt"
	if got := truncatePrompt(short, 200); got != short {
		t.Errorf("short prompt should not be truncated, got %q", got)
	}

	long := strings.Repeat("a", 300)
	got := truncatePrompt(long, 200)
	if len(got) != 203 { // 200 + "..."
		t.Errorf("truncated length = %d, want 203", len(got))
	}
	if !strings.HasSuffix(got, "...") {
		t.Error("truncated prompt should end with ...")
	}
}

func TestLogRunAppendsMultipleEntries(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir)
	t.Setenv("LEGAL_BOT_LOG", "")

	opts := RunOptions{LLM: "claude", Prompt: "first"}
	result := &Result{LLMUsed: "claude", Duration: time.Second}

	logRun(opts, "claude", result, nil, nil)
	logRun(opts, "claude", result, nil, nil)

	logPath := filepath.Join(tmpDir, ".config", "legal-bot", "logs", "legal-bot.log")
	data, _ := os.ReadFile(logPath)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 log lines, got %d", len(lines))
	}
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func containsArg(args []string, needle string) bool {
	for _, arg := range args {
		if arg == needle {
			return true
		}
	}
	return false
}

func assertContainsSequence(t *testing.T, args []string, sequence []string) {
	t.Helper()
	for i := 0; i <= len(args)-len(sequence); i++ {
		match := true
		for j := range sequence {
			if args[i+j] != sequence[j] {
				match = false
				break
			}
		}
		if match {
			return
		}
	}
	t.Fatalf("args %v do not contain sequence %v", args, sequence)
}

func writeFakeCLI(t *testing.T, dir string, name string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte("@echo off\r\n"), 0644); err != nil {
		t.Fatal(err)
	}
}
