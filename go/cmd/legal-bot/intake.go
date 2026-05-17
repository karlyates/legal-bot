package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jdonohoo/legal-bot/go/internal/config"
	"github.com/jdonohoo/legal-bot/go/internal/llm"
	"github.com/spf13/cobra"
)

var intakeCmd = &cobra.Command{
	Use:   "intake <matter>",
	Short: "Run the intake pipeline on a matter's source documents",
	Long: `Scan source documents in matters/<matter>/input/, build a knowledge layer
with document register, document cards, timeline, facts, and open questions.

Example:
  legal-bot intake TRO-MTE-GAL-SM`,
	Args: cobra.ExactArgs(1),
	RunE: runIntake,
}

var (
	intakeSingleLLM string
	intakeLLMMode   string
)

func init() {
	intakeCmd.Flags().StringVar(&intakeSingleLLM, "single-llm", "", "Use a single LLM for all steps")
	intakeCmd.Flags().StringVar(&intakeLLMMode, "llm-mode", "", "LLM fallback mode")
	rootCmd.AddCommand(intakeCmd)
}

func runIntake(cmd *cobra.Command, args []string) error {
	return executeIntake(intakeOptions{
		Matter:    args[0],
		SingleLLM: intakeSingleLLM,
		LLMMode:   intakeLLMMode,
	})
}

type intakeOptions struct {
	Matter    string
	SingleLLM string
	LLMMode   string
}

func executeIntake(opts intakeOptions) error {
	matter := opts.Matter

	projectRoot := resolveProjectRoot()
	agentsDir := resolveAgentsDir()
	cfg := config.Load(projectRoot)

	if opts.SingleLLM != "" {
		cfg.LLMMode = "single_llm"
		if mode, ok := cfg.LLMModes["single_llm"]; ok {
			mode.OverrideLLM = opts.SingleLLM
			mode.SynthesisLLM = opts.SingleLLM
			cfg.LLMModes["single_llm"] = mode
		}
	} else if opts.LLMMode != "" {
		cfg.LLMMode = opts.LLMMode
	}

	matterDir := filepath.Join(projectRoot, "matters", matter)
	inputDir := filepath.Join(matterDir, "input")
	knowledgeDir := filepath.Join(matterDir, "knowledge")
	cardDir := filepath.Join(knowledgeDir, "document_cards")
	workingDir := filepath.Join(matterDir, "working")
	caseDir := filepath.Join(projectRoot, "case")

	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		return fmt.Errorf("input directory not found: %s\nCreate it and add source documents first", inputDir)
	}

	os.MkdirAll(knowledgeDir, 0755)
	os.MkdirAll(cardDir, 0755)
	os.MkdirAll(workingDir, 0755)

	inputFiles := countTextFiles(inputDir)
	if inputFiles == 0 {
		return fmt.Errorf("no .txt or .md files found in %s\nAdd source documents and try again", inputDir)
	}

	timeout := cfg.GetPipelineStepTimeout()

	fmt.Println()
	fmt.Println("  ============================================")
	fmt.Println("   Legal-Bot Intake Pipeline")
	fmt.Printf("   Matter: %s\n", matter)
	fmt.Printf("   Input files: %d\n", inputFiles)
	fmt.Println("  ============================================")
	fmt.Println()

	caseContext := loadCaseContext(caseDir)
	inputContext := buildInputContext(inputDir)

	steps := cfg.GetLegalPipeline("intake")
	if len(steps) == 0 {
		return fmt.Errorf("no intake pipeline steps found in config")
	}

	var prevOutput string
	var knowledgeContext string
	var failedSteps []string

	for _, step := range steps {
		stepNum := step.Step
		fmt.Printf("\n>>> Step %d/%d: %s (%s)\n", stepNum, len(steps), step.Name, step.LLM)

		var outputFile string
		switch stepNum {
		case 1:
			outputFile = filepath.Join(workingDir, "01-intake-mapping.md")
		case 2:
			outputFile = filepath.Join(workingDir, "02-document-cards.md")
		case 3:
			outputFile = filepath.Join(knowledgeDir, "case_timeline.md")
		case 4:
			outputFile = filepath.Join(knowledgeDir, "fact_candidates.md")
		case 5:
			outputFile = filepath.Join(knowledgeDir, "open_questions_for_user.md")
		default:
			outputFile = filepath.Join(workingDir, fmt.Sprintf("%02d-%s.md", stepNum, step.Persona))
		}

		var prompt string
		switch step.ContextMode {
		case "prompt_only":
			prompt = step.PromptPrefix + "\n\n=== CASE CONTEXT ===\n" + caseContext + "\n=== END CASE CONTEXT ===\n\n=== INPUT DOCUMENTS ===\n" + inputContext + "\n=== END INPUT DOCUMENTS ==="
		case "previous":
			prompt = step.PromptPrefix + "\n\n=== CASE CONTEXT ===\n" + caseContext + "\n=== END CASE CONTEXT ===\n\n=== PREVIOUS STEP OUTPUT ===\n" + prevOutput + "\n=== END PREVIOUS STEP OUTPUT ===\n\n=== INPUT DOCUMENTS ===\n" + inputContext + "\n=== END INPUT DOCUMENTS ==="
		case "knowledge_layer":
			if knowledgeContext == "" {
				knowledgeContext = buildKnowledgeContext(knowledgeDir)
			}
			prompt = step.PromptPrefix + "\n\n=== CASE CONTEXT ===\n" + caseContext + "\n=== END CASE CONTEXT ===\n\n=== KNOWLEDGE LAYER ===\n" + knowledgeContext + "\n=== END KNOWLEDGE LAYER ==="
		default:
			prompt = step.PromptPrefix + "\n\n=== CASE CONTEXT ===\n" + caseContext + "\n=== END CASE CONTEXT ===\n\n" + inputContext
		}

		selection, routeErr := selectStepRoute(cfg, step.ModelProfile, step.LLM)
		if routeErr != nil {
			fmt.Printf("    FAILED: %s %s\n", step.Name, routeErr.Error())
			failedSteps = append(failedSteps, step.Name)
			os.WriteFile(outputFile, []byte(fmt.Sprintf("# STEP FAILED\n\nStep %d (%s) failed.\n\nRoute error: %s\n", stepNum, step.Name, routeErr.Error())), 0644)
			continue
		}
		fmt.Printf("    Route: %s\n", describeRoute(step.ModelProfile, selection))
		printRouteWarnings(selection.Warnings)

		result, err := llm.Run(runOptionsForRoute(selection, prompt, outputFile, step.Persona, timeout, agentsDir))
		if err != nil || result.ExitCode != 0 {
			detail := ""
			if result != nil && result.Stderr != "" {
				detail = llm.FirstLine(result.Stderr)
			}
			fmt.Printf("    FAILED: %s %s\n", step.Name, detail)
			failedSteps = append(failedSteps, step.Name)
			os.WriteFile(outputFile, []byte(fmt.Sprintf("# STEP FAILED\n\nStep %d (%s) failed.\n", stepNum, step.Name)), 0644)
			continue
		}

		outputBytes := fileSize(outputFile)
		fmt.Printf("    OK (%s, %d bytes, %s)\n", result.LLMUsed, outputBytes, result.Duration.Truncate(time.Second))

		data, _ := os.ReadFile(outputFile)
		prevOutput = string(data)

		if stepNum == 2 {
			splitDocumentCards(string(data), cardDir)
			knowledgeContext = buildKnowledgeContext(knowledgeDir)
		}
		if stepNum == 1 {
			extractRegister(string(data), filepath.Join(knowledgeDir, "document_register.csv"))
			knowledgeContext = buildKnowledgeContext(knowledgeDir)
		}
	}

	writeIngestionReport(matterDir, matter, inputFiles)

	statusTitle, statusNote, summaryErr := intakeSummary(len(steps), failedSteps)

	fmt.Println()
	fmt.Println("  ============================================")
	fmt.Printf("   %s\n", statusTitle)
	fmt.Printf("   Knowledge: matters/%s/knowledge/\n", matter)
	if statusNote != "" {
		fmt.Printf("   %s\n", statusNote)
	}
	fmt.Println("  ============================================")
	fmt.Println()
	if len(failedSteps) > 0 {
		fmt.Println("  Failed steps:")
		for _, name := range failedSteps {
			fmt.Printf("  - %s\n", name)
		}
		fmt.Println()
	}
	fmt.Println("  Next steps:")
	fmt.Printf("  1. Review: matters/%s/knowledge/open_questions_for_user.md\n", matter)
	fmt.Printf("  2. Add or update: matters/%s/situation.md, goal.md, and drafts/ as needed\n", matter)
	fmt.Println("  3. Run a workflow or draft review")

	return summaryErr
}

// countTextFiles counts .txt and .md files in a directory
func countTextFiles(dir string) int {
	count := 0
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext == ".txt" || ext == ".md" {
			count++
		}
	}
	return count
}

// loadCaseContext reads all files from the case/ directory
func loadCaseContext(caseDir string) string {
	var ctx strings.Builder
	entries, err := os.ReadDir(caseDir)
	if err != nil {
		return ""
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext != ".md" && ext != ".txt" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(caseDir, e.Name()))
		if err != nil {
			continue
		}
		ctx.WriteString(fmt.Sprintf("\n--- %s ---\n%s\n", e.Name(), string(data)))
	}
	return ctx.String()
}

// buildInputContext reads all text files from the input directory
func buildInputContext(inputDir string) string {
	var ctx strings.Builder
	entries, _ := os.ReadDir(inputDir)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext != ".txt" && ext != ".md" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(inputDir, e.Name()))
		if err != nil {
			continue
		}
		ctx.WriteString(fmt.Sprintf("\n=== FILE: %s ===\n%s\n=== END FILE ===\n", e.Name(), string(data)))
	}
	return ctx.String()
}

// buildKnowledgeContext reads all knowledge artifacts
func buildKnowledgeContext(knowledgeDir string) string {
	var ctx strings.Builder

	// Read top-level files
	entries, _ := os.ReadDir(knowledgeDir)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(knowledgeDir, e.Name()))
		if err != nil {
			continue
		}
		ctx.WriteString(fmt.Sprintf("\n--- %s ---\n%s\n", e.Name(), string(data)))
	}

	// Read document cards
	cardDir := filepath.Join(knowledgeDir, "document_cards")
	cards, _ := os.ReadDir(cardDir)
	for _, c := range cards {
		if c.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(cardDir, c.Name()))
		if err != nil {
			continue
		}
		ctx.WriteString(fmt.Sprintf("\n--- document_cards/%s ---\n%s\n", c.Name(), string(data)))
	}

	return ctx.String()
}

// splitDocumentCards splits a combined document cards output into individual files
func splitDocumentCards(content string, cardDir string) {
	heading := "# Source Document Card:"
	if !strings.Contains(content, heading) {
		heading = "# Document Card:"
	}
	sections := strings.Split(content, heading)
	for i, section := range sections {
		if i == 0 || strings.TrimSpace(section) == "" {
			continue
		}
		// Extract doc ID from first line
		lines := strings.SplitN(section, "\n", 2)
		docID := strings.TrimSpace(lines[0])
		docID = strings.ReplaceAll(docID, " ", "-")
		if docID == "" {
			docID = fmt.Sprintf("card-%d", i)
		}
		filename := fmt.Sprintf("%s.md", docID)
		cardContent := heading + section
		os.WriteFile(filepath.Join(cardDir, filename), []byte(cardContent), 0644)
	}
}

// extractRegister extracts CSV content from intake mapper output
func extractRegister(content string, outputPath string) {
	// Look for CSV block in the output
	lines := strings.Split(content, "\n")
	var csvLines []string
	inCSV := false
	for _, line := range lines {
		if strings.HasPrefix(line, "doc_id,") || strings.HasPrefix(line, "DOC-") {
			inCSV = true
		}
		if inCSV {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || trimmed == "```" {
				if len(csvLines) > 0 {
					break
				}
				continue
			}
			csvLines = append(csvLines, trimmed)
		}
	}
	if len(csvLines) > 0 {
		os.WriteFile(outputPath, []byte(strings.Join(csvLines, "\n")+"\n"), 0644)
	}
}

// writeIngestionReport creates a summary report
func writeIngestionReport(matterDir string, matter string, inputFiles int) {
	report := fmt.Sprintf(`# Ingestion Report

**Matter:** %s
**Date:** %s
**Input Files:** %d

## Generated Artifacts
- document_register.csv
- document_cards/ (per-document summaries)
- case_timeline.md
- fact_candidates.md
- open_questions_for_user.md

> This report was generated by Legal-Bot. It is AI-generated legal analysis and knowledge-layer support.
> Verify important points against source documents before filing, serving, sending, or relying on them externally.
`, matter, time.Now().Format("2006-01-02 15:04:05"), inputFiles)

	os.WriteFile(filepath.Join(matterDir, "output", "ingestion_report.md"), []byte(report), 0644)
}

// fileSize returns file size in bytes
func fileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

func intakeSummary(totalSteps int, failedSteps []string) (string, string, error) {
	failed := len(failedSteps)
	switch {
	case failed == 0:
		return "Intake Complete", "", nil
	case failed >= totalSteps && totalSteps > 0:
		return "Intake Failed", fmt.Sprintf("%d/%d steps failed.", failed, totalSteps), fmt.Errorf("intake failed: %d/%d steps failed", failed, totalSteps)
	default:
		return "Intake Completed with Failures", fmt.Sprintf("%d/%d steps failed.", failed, totalSteps), fmt.Errorf("intake completed with failures: %d/%d steps failed", failed, totalSteps)
	}
}

// resolveProjectRoot finds the project root directory
func resolveProjectRoot() string {
	// Try relative to binary
	exe, err := os.Executable()
	if err == nil {
		exe, _ = filepath.EvalSymlinks(exe)
		// Binary at {root}/go/bin/legal-bot
		root := filepath.Dir(filepath.Dir(filepath.Dir(exe)))
		if _, err := os.Stat(filepath.Join(root, "case")); err == nil {
			return root
		}
		// Binary at {root}/go/cmd/legal-bot/legal-bot
		root = filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(exe))))
		if _, err := os.Stat(filepath.Join(root, "case")); err == nil {
			return root
		}
	}

	// Try LEGAL_BOT_ROOT env
	if root := os.Getenv("LEGAL_BOT_ROOT"); root != "" {
		return root
	}

	// Fallback: current directory
	cwd, _ := os.Getwd()
	return cwd
}
