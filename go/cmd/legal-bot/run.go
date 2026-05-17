package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jdonohoo/legal-bot/go/internal/llm"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <llm> <prompt>",
	Short: "Run a single LLM subprocess",
	Long: `Spawn an LLM subprocess with timeout and optional persona context.

Supported LLMs: claude (c), codex (x), gemini (g), copilot (p)

Exit codes:
  0    Success
  1    Usage error / unknown LLM
  124  Timeout (GNU timeout convention)`,
	Args: cobra.ExactArgs(2),
	RunE: runRun,
}

var (
	runOutputFile string
	runPersona    string
	runTimeout    int
)

func init() {
	runCmd.Flags().StringVarP(&runOutputFile, "output", "o", "", "File to save output to")
	runCmd.Flags().StringVarP(&runPersona, "persona", "p", "", "Persona name (loads agents/{persona}.md)")
	runCmd.Flags().IntVarP(&runTimeout, "timeout", "t", 0, "Timeout in seconds (default: LEGAL_BOT_TIMEOUT or legacy VERN_TIMEOUT, else 1200)")
	rootCmd.AddCommand(runCmd)
}

func runRun(cmd *cobra.Command, args []string) error {
	llmName := args[0]
	prompt := args[1]

	// Resolve timeout: flag > env > default
	timeout := 1200
	if envTimeout := os.Getenv("LEGAL_BOT_TIMEOUT"); envTimeout != "" {
		fmt.Sscanf(envTimeout, "%d", &timeout)
	} else if envTimeout := os.Getenv("VERN_TIMEOUT"); envTimeout != "" {
		fmt.Sscanf(envTimeout, "%d", &timeout)
	}
	if runTimeout > 0 {
		timeout = runTimeout
	}

	// Find agents dir relative to binary
	agentsDir := resolveAgentsDir()

	opts := llm.RunOptions{
		LLM:        llmName,
		Prompt:     prompt,
		OutputFile: runOutputFile,
		Persona:    runPersona,
		Timeout:    time.Duration(timeout) * time.Second,
		WorkingDir: envOrFallback("LEGAL_BOT_WORKING_DIR", "VERN_WORKING_DIR"),
		AgentsDir:  agentsDir,
	}

	result, err := llm.Run(opts)
	if err != nil {
		return err
	}

	// Print output to stdout (tee behavior when output file is set)
	if result.Output != "" {
		fmt.Print(result.Output)
	}

	if result.ExitCode != 0 {
		os.Exit(result.ExitCode)
	}

	return nil
}

// resolveAgentsDir finds the agents/ directory relative to the binary or project root.
func resolveAgentsDir() string {
	// Try relative to binary: binary is at go/bin/legal-bot or go/cmd/legal-bot/legal-bot.
	exe, err := os.Executable()
	if err == nil {
		exe, _ = filepath.EvalSymlinks(exe)

		// Binary at {root}/go/bin/legal-bot -> agents at {root}/agents/.
		root := filepath.Dir(filepath.Dir(filepath.Dir(exe)))
		candidate := filepath.Join(root, "agents")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}

		// Binary at {root}/go/cmd/legal-bot/legal-bot -> agents at {root}/agents/.
		root = filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(exe))))
		candidate = filepath.Join(root, "agents")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
	}

	// Fallback: look for LEGAL_BOT_ROOT first, then legacy VERN_ROOT.
	if root := envOrFallback("LEGAL_BOT_ROOT", "VERN_ROOT"); root != "" {
		return filepath.Join(root, "agents")
	}

	// Last resort: relative to cwd
	return "agents"
}

func envOrFallback(primary string, legacy string) string {
	if value := os.Getenv(primary); value != "" {
		return value
	}
	return os.Getenv(legacy)
}
