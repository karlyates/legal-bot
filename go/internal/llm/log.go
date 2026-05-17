package llm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type logEntry struct {
	Time            string `json:"time"`
	LLMRequested    string `json:"llm_requested"`
	ModelRequested  string `json:"model_requested,omitempty"`
	EffortRequested string `json:"effort_requested,omitempty"`
	LLMUsed         string `json:"llm_used"`
	ModelUsed       string `json:"model_used,omitempty"`
	EffortUsed      string `json:"effort_used,omitempty"`
	ExitCode        int    `json:"exit_code"`
	TimedOut        bool   `json:"timed_out"`
	DurationMs      int64  `json:"duration_ms"`
	Error           string `json:"error,omitempty"`
	Stderr          string `json:"stderr,omitempty"`
	OutputFile      string `json:"output_file,omitempty"`
	OutputBytes     int    `json:"output_bytes"`
	PromptPreview   string `json:"prompt_preview"`
}

// logRun appends a JSONL entry to ~/.config/legal-bot/logs/legal-bot.log.
// Disabled if LEGAL_BOT_LOG=0 or legacy VERN_LOG=0.
func logRun(opts RunOptions, llmRequested string, result *Result, runErr error, writeErr error) {
	if os.Getenv("LEGAL_BOT_LOG") == "0" || os.Getenv("VERN_LOG") == "0" {
		return
	}

	logDir := filepath.Join(configDir(), "logs")
	if mkErr := os.MkdirAll(logDir, 0755); mkErr != nil {
		return
	}

	entry := logEntry{
		Time:            time.Now().UTC().Format(time.RFC3339),
		LLMRequested:    llmRequested,
		ModelRequested:  opts.Model,
		EffortRequested: opts.Effort,
		PromptPreview:   truncatePrompt(opts.Prompt, 200),
	}

	if opts.OutputFile != "" {
		entry.OutputFile = opts.OutputFile
	}

	if result != nil {
		entry.LLMUsed = result.LLMUsed
		entry.ModelUsed = result.ModelUsed
		entry.EffortUsed = result.EffortUsed
		entry.ExitCode = result.ExitCode
		entry.TimedOut = result.TimedOut
		entry.DurationMs = result.Duration.Milliseconds()
		entry.OutputBytes = len(result.Output)
		if result.Stderr != "" {
			entry.Stderr = truncatePrompt(result.Stderr, 500)
		}
	}

	if runErr != nil {
		entry.Error = runErr.Error()
	} else if writeErr != nil {
		entry.Error = writeErr.Error()
	}

	data, jsonErr := json.Marshal(entry)
	if jsonErr != nil {
		return
	}
	data = append(data, '\n')

	logPath := filepath.Join(logDir, "legal-bot.log")
	f, fErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fErr != nil {
		return
	}
	defer f.Close()

	f.Write(data)
}

func truncatePrompt(prompt string, max int) string {
	if len(prompt) <= max {
		return prompt
	}
	return prompt[:max] + "..."
}

func configDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[legal-bot-log] Warning: cannot determine home directory: %v\n", err)
		return "/tmp/legal-bot"
	}
	return filepath.Join(home, ".config", "legal-bot")
}
