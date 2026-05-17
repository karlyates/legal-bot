package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jdonohoo/legal-bot/go/internal/config"
	"github.com/jdonohoo/legal-bot/go/internal/llm"
	"github.com/spf13/cobra"
)

var reviewCmd = &cobra.Command{
	Use:     "review <matter>",
	Aliases: []string{"review-draft"},
	Short:   "Run the draft-review workflow on a matter",
	Long: `Review a legal draft against the matter's knowledge layer.
The most recent .txt or .md file in matters/<matter>/drafts/ is reviewed.

Example:
  legal-bot review TRO-MTE-GAL-SM`,
	Args: cobra.ExactArgs(1),
	RunE: runReview,
}

var (
	reviewSingleLLM    string
	reviewLLMMode      string
	reviewDraft        string
	reviewMode         string
	reviewDocumentType string
	reviewChildRelated bool
	reviewFinancial    bool
	reviewLevel        string
	reviewPlan         bool
	reviewYes          bool
)

func init() {
	reviewCmd.Flags().StringVar(&reviewSingleLLM, "single-llm", "", "Use a single LLM for all steps")
	reviewCmd.Flags().StringVar(&reviewLLMMode, "llm-mode", "", "LLM fallback mode")
	reviewCmd.Flags().StringVar(&reviewDraft, "draft", "", "Specific draft file to review (default: latest in drafts/)")
	reviewCmd.Flags().StringVar(&reviewMode, "mode", "standard", "Review mode: quick, standard, deep, legal-research, or strategy")
	reviewCmd.Flags().StringVar(&reviewDocumentType, "document-type", "", "Optional document type route: motion, opposition, response, reply, declaration, proposed-order, co-parenting-communication")
	reviewCmd.Flags().BoolVar(&reviewChildRelated, "child-related", false, "Add the child best-interests/family-dynamics reviewer")
	reviewCmd.Flags().BoolVar(&reviewFinancial, "financial", false, "Add the financial support reviewer")
	reviewCmd.Flags().StringVar(&reviewLevel, "level", string(analysisLevelScout), "Analysis level: scout, standard, deep, or max")
	reviewCmd.Flags().BoolVar(&reviewPlan, "plan", false, "Print the planned review route without calling any LLMs")
	reviewCmd.Flags().BoolVar(&reviewYes, "yes", false, "Skip confirmation for deep/max review runs")
	rootCmd.AddCommand(reviewCmd)
}

func runReview(cmd *cobra.Command, args []string) error {
	return executeReview(reviewOptions{
		Matter:       args[0],
		SingleLLM:    reviewSingleLLM,
		LLMMode:      reviewLLMMode,
		Draft:        reviewDraft,
		Mode:         reviewMode,
		DocumentType: reviewDocumentType,
		ChildRelated: reviewChildRelated,
		Financial:    reviewFinancial,
		Level:        reviewLevel,
		Plan:         reviewPlan,
		Yes:          reviewYes,
	})
}

type reviewOptions struct {
	Matter       string
	SingleLLM    string
	LLMMode      string
	Draft        string
	Mode         string
	DocumentType string
	ChildRelated bool
	Financial    bool
	Level        string
	Plan         bool
	Yes          bool
}

func executeReview(opts reviewOptions) error {
	matter := opts.Matter
	level, err := parseAnalysisLevel(opts.Level)
	if err != nil {
		return err
	}

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
	draftsDir := filepath.Join(matterDir, "drafts")
	knowledgeDir := filepath.Join(matterDir, "knowledge")
	workingDir := filepath.Join(matterDir, "working")
	outputDir := filepath.Join(matterDir, "output")
	caseDir := filepath.Join(projectRoot, "case")

	registerPath := filepath.Join(knowledgeDir, "document_register.csv")
	if _, err := os.Stat(registerPath); os.IsNotExist(err) {
		return fmt.Errorf("no knowledge layer found for matter %s\nRun intake first: legal-bot intake %s", matter, matter)
	}

	draftPath, err := findDraft(draftsDir, opts.Draft)
	if err != nil {
		return err
	}

	draftData, err := os.ReadFile(draftPath)
	if err != nil {
		return fmt.Errorf("read draft: %w", err)
	}
	draftContent := string(draftData)
	draftName := filepath.Base(draftPath)

	timeout := cfg.GetPipelineStepTimeout()

	caseContext := loadCaseContext(caseDir)
	knowledgeContext := buildKnowledgeContext(knowledgeDir)

	pipelineKey := reviewPipelineKey(opts.Mode, opts.DocumentType)
	steps := cfg.GetLegalPipeline(pipelineKey)
	if len(steps) == 0 {
		return fmt.Errorf("no review pipeline steps found in config for %q", pipelineKey)
	}
	steps = applyConditionalReviewRouting(steps, draftContent, opts.Mode, opts.DocumentType, opts.ChildRelated, opts.Financial)
	steps = applyReviewLevel(steps, level)
	if opts.Plan {
		fmt.Println()
		printRunPlan(os.Stdout, buildRunPlanSummary(cfg, routeKindReview, matter, level, pipelineKey, steps))
		return nil
	}
	maybeWarnLargeRun(os.Stdout, steps)
	if err := maybeConfirmLargeRun(routeKindReview, level, steps, opts.Yes); err != nil {
		return err
	}

	os.MkdirAll(workingDir, 0755)
	os.MkdirAll(outputDir, 0755)

	fmt.Println()
	fmt.Println("  ============================================")
	fmt.Println("   Legal-Bot Draft Review Pipeline")
	fmt.Printf("   Matter: %s\n", matter)
	fmt.Printf("   Draft: %s\n", draftName)
	fmt.Printf("   Level: %s\n", level)
	fmt.Println("  ============================================")
	fmt.Println()

	var reviewOutputs strings.Builder

	for _, step := range steps {
		stepNum := step.Step
		fmt.Printf("\n>>> Step %d/%d: %s (%s)\n", stepNum, len(steps), step.Name, step.LLM)

		var outputFile string
		if stepNum == len(steps) {
			outputFile = filepath.Join(outputDir, "final_review_packet.md")
		} else {
			slug := strings.ToLower(strings.ReplaceAll(step.Name, " ", "_"))
			outputFile = filepath.Join(workingDir, fmt.Sprintf("%02d-%s.md", stepNum, slug))
		}

		var prompt string
		switch step.ContextMode {
		case "draft_plus_knowledge":
			prompt = step.PromptPrefix +
				"\n\n=== CASE CONTEXT ===\n" + caseContext + "\n=== END CASE CONTEXT ===" +
				"\n\n=== DRAFT UNDER REVIEW ===\n" + draftContent + "\n=== END DRAFT ===" +
				"\n\n=== KNOWLEDGE LAYER ===\n" + knowledgeContext + "\n=== END KNOWLEDGE LAYER ==="
		case "draft_plus_reviews":
			prompt = step.PromptPrefix +
				"\n\n=== CASE CONTEXT ===\n" + caseContext + "\n=== END CASE CONTEXT ===" +
				"\n\n=== DRAFT UNDER REVIEW ===\n" + draftContent + "\n=== END DRAFT ===" +
				"\n\n=== PREVIOUS REVIEW FINDINGS ===\n" + reviewOutputs.String() + "\n=== END REVIEW FINDINGS ==="
		case "all_reviews":
			prompt = step.PromptPrefix +
				"\n\n=== CASE CONTEXT ===\n" + caseContext + "\n=== END CASE CONTEXT ===" +
				"\n\n=== DRAFT FILENAME ===\n" + draftName + "\n=== END DRAFT FILENAME ===" +
				"\n\n=== ALL REVIEW FINDINGS ===\n" + reviewOutputs.String() + "\n=== END REVIEW FINDINGS ==="
		default:
			prompt = step.PromptPrefix +
				"\n\n=== DRAFT ===\n" + draftContent + "\n=== END DRAFT ==="
		}

		selection, routeErr := selectStepRoute(cfg, step.ModelProfile, step.LLM)
		if routeErr != nil {
			fmt.Printf("    FAILED: %s %s\n", step.Name, routeErr.Error())
			os.WriteFile(outputFile, []byte(fmt.Sprintf("# STEP FAILED\n\nStep %d (%s) failed.\n\nRoute error: %s\n", stepNum, step.Name, routeErr.Error())), 0644)
			continue
		}
		fmt.Printf("    Route: %s\n", describeRoute(step.ModelProfile, selection))
		printRouteWarnings(selection.Warnings)

		result, runErr := llm.Run(runOptionsForRoute(selection, prompt, outputFile, step.Persona, timeout, agentsDir))
		if runErr != nil || result.ExitCode != 0 {
			detail := ""
			if result != nil && result.Stderr != "" {
				detail = llm.FirstLine(result.Stderr)
			}
			fmt.Printf("    FAILED: %s %s\n", step.Name, detail)
			os.WriteFile(outputFile, []byte(fmt.Sprintf("# STEP FAILED\n\nStep %d (%s) failed.\n", stepNum, step.Name)), 0644)
			continue
		}

		outBytes := fileSize(outputFile)
		fmt.Printf("    OK (%s, %d bytes, %s)\n", result.LLMUsed, outBytes, result.Duration.Truncate(time.Second))

		data, _ := os.ReadFile(outputFile)
		reviewOutputs.WriteString(fmt.Sprintf("\n\n--- %s ---\n%s\n", step.Name, string(data)))
	}

	fmt.Println()
	fmt.Println("  ============================================")
	fmt.Println("   Review Complete")
	fmt.Printf("   Output: matters/%s/output/final_review_packet.md\n", matter)
	fmt.Println("  ============================================")
	fmt.Println()
	fmt.Println("  Next steps:")
	fmt.Println("  1. Read the review packet")
	fmt.Println("  2. Check Bottom Line, Must Fix, and Do Not Chase")
	fmt.Println("  3. Verify anything you plan to file, send, or rely on externally")
	fmt.Println()
	fmt.Println("  > AI-generated legal analysis and strategy support. Verify before external use.")

	return nil
}

func reviewPipelineKey(mode string, documentType string) string {
	if key, ok := reviewDocumentTypePipelineKey(documentType); ok {
		return key
	}
	return reviewModePipelineKey(mode)
}

func applyConditionalReviewRouting(steps []config.PipelineStep, draft string, mode string, documentType string, childFlag bool, financialFlag bool) []config.PipelineStep {
	route := normalizeRouteName(documentType)
	modeName := normalizeRouteName(mode)
	deep := modeName == "deep"
	strategyMode := modeName == "strategy"
	legalResearchMode := modeName == "legal-research" || modeName == "legal_research" || modeName == "research"
	rules := legalConditionalRules()
	childDetected := childFlag || reviewDocumentTriggersChildIssues(route) || (deep && draftMentionsAny(draft, rules[0].Keywords))
	financialDetected := financialFlag || reviewDocumentTriggersFinancialIssues(route) || (deep && draftMentionsAny(draft, rules[1].Keywords))
	strategyDetected := strategyMode || reviewDocumentTriggersStrategy(route) || draftMentionsAny(draft, rules[3].Keywords)
	legalDetected := legalResearchMode || reviewDocumentTriggersLegalAuthority(route) || draftMentionsAny(draft, rules[2].Keywords)

	if childDetected {
		steps = insertBeforePersona(steps, "legal-writing-preservation-editor", config.PipelineStep{
			Name:         "Child Best Interests and Family Dynamics Review",
			Persona:      "child-best-interests-family-dynamics-reviewer",
			LLM:          "claude",
			ModelProfile: "review_reasoning",
			ContextMode:  "draft_plus_knowledge",
			PromptPrefix: "You are the Child Best Interests and Family Dynamics Reviewer. Review child-related, custody, school, therapy, medical, parent-time, and co-parenting content carefully. Follow your agent instructions exactly.",
		})
	}
	if financialDetected {
		steps = insertBeforePersona(steps, "legal-writing-preservation-editor", config.PipelineStep{
			Name:         "Financial Support Review",
			Persona:      "financial-support-reviewer",
			LLM:          "claude",
			ModelProfile: "review_reasoning",
			ContextMode:  "draft_plus_knowledge",
			PromptPrefix: "You are the Financial Support Reviewer. Review financial/support issues, assumptions, documentation, and attorney questions. Follow your agent instructions exactly.",
		})
	}
	if legalDetected {
		steps = insertAfterPersona(steps, "trial-fact-checker", config.PipelineStep{
			Name:         "Legal Authority Scholar Review",
			Persona:      "legal-authority-scholar",
			LLM:          "claude",
			ModelProfile: "legal_research",
			ContextMode:  "draft_plus_knowledge",
			PromptPrefix: "You are the Legal Authority Scholar. Verify legal standards, statutes, rules, cases, procedural vehicles, and remedy authority using a citation-first approach. Follow your agent instructions exactly.",
		})
	}
	if strategyDetected {
		steps = insertBeforePersona(steps, "practical-resolution-reviewer", config.PipelineStep{
			Name:         "Strategic Options Review",
			Persona:      "strategic-options-architect",
			LLM:          "claude",
			ModelProfile: "strategy_reasoning",
			ContextMode:  "draft_plus_knowledge",
			PromptPrefix: "You are the Strategic Options Architect. Identify lawful strategic options, status-quo leverage, procedural paths, and expert or evaluation ideas for attorney review. Follow your agent instructions exactly.",
		})
	}

	for i := range steps {
		steps[i].Step = i + 1
	}
	return steps
}

func insertBeforeFinal(steps []config.PipelineStep, step config.PipelineStep) []config.PipelineStep {
	for _, existing := range steps {
		if existing.Persona == step.Persona {
			return steps
		}
	}
	finalIdx := len(steps)
	for i, existing := range steps {
		if existing.Persona == "managing-partner-final-synthesizer" {
			finalIdx = i
			break
		}
	}
	steps = append(steps, config.PipelineStep{})
	copy(steps[finalIdx+1:], steps[finalIdx:])
	steps[finalIdx] = step
	return steps
}

func insertBeforePersona(steps []config.PipelineStep, targetPersona string, step config.PipelineStep) []config.PipelineStep {
	for _, existing := range steps {
		if existing.Persona == step.Persona {
			return steps
		}
	}
	targetIdx := len(steps)
	for i, existing := range steps {
		if existing.Persona == targetPersona {
			targetIdx = i
			break
		}
	}
	if targetIdx == len(steps) {
		return insertBeforeFinal(steps, step)
	}
	steps = append(steps, config.PipelineStep{})
	copy(steps[targetIdx+1:], steps[targetIdx:])
	steps[targetIdx] = step
	return steps
}

func insertAfterPersona(steps []config.PipelineStep, targetPersona string, step config.PipelineStep) []config.PipelineStep {
	for _, existing := range steps {
		if existing.Persona == step.Persona {
			return steps
		}
	}
	for i, existing := range steps {
		if existing.Persona == targetPersona {
			insertIdx := i + 1
			steps = append(steps, config.PipelineStep{})
			copy(steps[insertIdx+1:], steps[insertIdx:])
			steps[insertIdx] = step
			return steps
		}
	}
	return insertBeforeFinal(steps, step)
}

func normalizeRouteName(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func draftMentionsAny(draft string, terms []string) bool {
	lower := strings.ToLower(draft)
	for _, term := range terms {
		if strings.Contains(lower, term) {
			return true
		}
	}
	return false
}

// findDraft locates the draft file to review
func findDraft(draftsDir string, specified string) (string, error) {
	if specified != "" {
		absPath, _ := filepath.Abs(specified)
		if _, err := os.Stat(absPath); err == nil {
			return absPath, nil
		}
		return "", fmt.Errorf("specified draft not found: %s", specified)
	}

	entries, err := os.ReadDir(draftsDir)
	if err != nil {
		return "", fmt.Errorf("drafts directory not found: %s\nCreate it and add a draft file", draftsDir)
	}

	// Find most recently modified .txt or .md file
	type draftFile struct {
		path    string
		modTime time.Time
	}
	var drafts []draftFile

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext != ".txt" && ext != ".md" {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		drafts = append(drafts, draftFile{
			path:    filepath.Join(draftsDir, e.Name()),
			modTime: info.ModTime(),
		})
	}

	if len(drafts) == 0 {
		return "", fmt.Errorf("no .txt or .md files found in %s\nAdd your attorney's draft and try again", draftsDir)
	}

	sort.Slice(drafts, func(i, j int) bool {
		return drafts[i].modTime.After(drafts[j].modTime)
	})

	fmt.Printf("  Reviewing: %s\n", filepath.Base(drafts[0].path))
	return drafts[0].path, nil
}
