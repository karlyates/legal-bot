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

type workflowSpec struct {
	Type              string
	DisplayName       string
	PipelineKey       string
	OutputFile        string
	RequireSituation  bool
	RequireAudience   bool
	AllowDraft        bool
	AutoDraftFallback bool
}

type workflowInputs struct {
	SituationPath string
	Situation     string
	GoalPath      string
	Goal          string
	QuestionsPath string
	Questions     string
	DraftPath     string
	DraftName     string
	Draft         string
}

var workflowCmd = &cobra.Command{
	Use:   "workflow <matter>",
	Short: "Run a user-job workflow such as triage, counsel brief, or authority check",
	Long: `Run a workflow-oriented legal analysis path without requiring the user to think in raw pipeline names.

Examples:
  legal-bot workflow TRO-MTE-GAL-SM --type triage
  legal-bot workflow TRO-MTE-GAL-SM --type counsel-brief
  legal-bot workflow TRO-MTE-GAL-SM --type draft-message --audience ofw`,
	Args: cobra.ExactArgs(1),
	RunE: runWorkflow,
}

var (
	workflowType         string
	workflowAudience     string
	workflowUrgency      string
	workflowSingleLLM    string
	workflowLLMMode      string
	workflowSituation    string
	workflowGoal         string
	workflowQuestions    string
	workflowDraft        string
	workflowChildRelated bool
	workflowFinancial    bool
)

func init() {
	workflowCmd.Flags().StringVar(&workflowType, "type", "", "Workflow type: triage, counsel-brief, sm-triage, draft-message, evidence-packet, pattern-review, fact-lock, authority-check, prep, or journal-entry")
	workflowCmd.Flags().StringVar(&workflowAudience, "audience", "", "Audience hint for workflows like draft-message: ofw, special-master, provider, school, counsel, court-cover, or private-journal")
	workflowCmd.Flags().StringVar(&workflowUrgency, "urgency", "", "Urgency or mode hint: urgent, normal, call-agenda, evidence-update, draft-review-feedback, or strategic-question")
	workflowCmd.Flags().StringVar(&workflowSingleLLM, "single-llm", "", "Use a single LLM for all steps")
	workflowCmd.Flags().StringVar(&workflowLLMMode, "llm-mode", "", "LLM fallback mode")
	workflowCmd.Flags().StringVar(&workflowSituation, "situation", "", "Path to a situation file (default: matter-level situation.md conventions)")
	workflowCmd.Flags().StringVar(&workflowGoal, "goal", "", "Path to a goal file (default: matter-level goal.md conventions)")
	workflowCmd.Flags().StringVar(&workflowQuestions, "questions", "", "Path to a questions file (default: matter-level questions.md conventions)")
	workflowCmd.Flags().StringVar(&workflowDraft, "draft", "", "Optional draft file to include in workflow context")
	workflowCmd.Flags().BoolVar(&workflowChildRelated, "child-related", false, "Force child/family-dynamics routing when relevant")
	workflowCmd.Flags().BoolVar(&workflowFinancial, "financial", false, "Force financial/support routing when relevant")
	_ = workflowCmd.MarkFlagRequired("type")
	rootCmd.AddCommand(workflowCmd)
}

func runWorkflow(cmd *cobra.Command, args []string) error {
	return executeWorkflow(workflowRunOptions{
		Matter:       args[0],
		Type:         workflowType,
		Audience:     workflowAudience,
		Urgency:      workflowUrgency,
		SingleLLM:    workflowSingleLLM,
		LLMMode:      workflowLLMMode,
		SituationArg: workflowSituation,
		GoalArg:      workflowGoal,
		QuestionsArg: workflowQuestions,
		Draft:        workflowDraft,
		ChildRelated: workflowChildRelated,
		Financial:    workflowFinancial,
	})
}

type workflowRunOptions struct {
	Matter       string
	Type         string
	Audience     string
	Urgency      string
	SingleLLM    string
	LLMMode      string
	SituationArg string
	GoalArg      string
	QuestionsArg string
	Draft        string
	ChildRelated bool
	Financial    bool
}

func executeWorkflow(opts workflowRunOptions) error {
	matter := opts.Matter

	spec, ok := workflowSpecs()[normalizeRouteName(opts.Type)]
	if !ok {
		return fmt.Errorf("unknown workflow type %q\nSupported types: %s", opts.Type, strings.Join(sortedWorkflowTypes(), ", "))
	}
	contract, ok := workflowOutputContractForType(spec.Type)
	if !ok {
		return fmt.Errorf("missing workflow output contract for type %q", spec.Type)
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
	inputDir := filepath.Join(matterDir, "input")
	knowledgeDir := filepath.Join(matterDir, "knowledge")
	draftsDir := filepath.Join(matterDir, "drafts")
	workingDir := filepath.Join(matterDir, "working")
	outputDir := filepath.Join(matterDir, "output")
	caseDir := filepath.Join(projectRoot, "case")

	if spec.RequireAudience && strings.TrimSpace(opts.Audience) == "" {
		return fmt.Errorf("%s requires --audience\nExamples: ofw, special-master, provider, school, counsel, court-cover", spec.DisplayName)
	}

	inputs, err := loadWorkflowInputs(spec, matter, matterDir, draftsDir, opts)
	if err != nil {
		return err
	}

	steps := cfg.GetLegalPipeline(spec.PipelineKey)
	if len(steps) == 0 {
		return fmt.Errorf("no workflow pipeline steps found in config for %q", spec.PipelineKey)
	}

	caseContext := loadCaseContext(caseDir)
	knowledgeContext := buildKnowledgeContext(knowledgeDir)
	inputManifest := buildInputManifest(inputDir)
	steps = applyConditionalWorkflowRouting(steps, spec, inputs, opts.Audience, opts.Urgency, opts.ChildRelated, opts.Financial)

	os.MkdirAll(workingDir, 0755)
	os.MkdirAll(outputDir, 0755)

	timeout := cfg.GetPipelineStepTimeout()
	finalOutput := filepath.Join(outputDir, spec.OutputFile)

	fmt.Println()
	fmt.Println("  ============================================")
	fmt.Printf("   Legal-Bot %s\n", spec.DisplayName)
	fmt.Printf("   Matter: %s\n", matter)
	if inputs.DraftName != "" {
		fmt.Printf("   Draft context: %s\n", inputs.DraftName)
	}
	if opts.Audience != "" {
		fmt.Printf("   Audience: %s\n", opts.Audience)
	}
	fmt.Println("  ============================================")
	fmt.Println()

	var workflowOutputs strings.Builder
	workflowMeta := buildWorkflowMeta(spec, matter, opts.Audience, opts.Urgency, knowledgeContext != "", inputManifest != "")
	contractSummary := renderWorkflowContractSummary(contract)

	for _, step := range steps {
		stepNum := step.Step
		fmt.Printf("\n>>> Step %d/%d: %s (%s)\n", stepNum, len(steps), step.Name, step.LLM)
		isFinalStep := stepNum == len(steps)

		var outputFile string
		if isFinalStep {
			outputFile = finalOutput
		} else {
			slug := strings.ToLower(strings.ReplaceAll(step.Name, " ", "_"))
			outputFile = filepath.Join(workingDir, fmt.Sprintf("%s-%02d-%s.md", spec.Type, stepNum, slug))
		}

		contractDetail := ""
		if isFinalStep || step.Persona == "managing-partner-final-synthesizer" {
			contractDetail = renderWorkflowContract(contract)
		}

		prompt := buildWorkflowPrompt(step.ContextMode, step.PromptPrefix, workflowMeta, contractSummary, contractDetail, caseContext, inputs, knowledgeContext, inputManifest, workflowOutputs.String())

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
		workflowOutputs.WriteString(fmt.Sprintf("\n\n--- %s ---\n%s\n", step.Name, string(data)))
	}

	fmt.Println()
	fmt.Println("  ============================================")
	fmt.Printf("   %s Complete\n", spec.DisplayName)
	fmt.Printf("   Output: matters/%s/output/%s\n", matter, spec.OutputFile)
	fmt.Println("  ============================================")
	fmt.Println()
	fmt.Println("  Next steps:")
	fmt.Println("  1. Read the Bottom Line and Recommendation sections first")
	fmt.Println("  2. Verify any important factual or legal points against sources")
	fmt.Println("  3. Use any draft message only after confirming the audience and record posture")
	fmt.Println()
	fmt.Println("  > AI-generated legal analysis and strategy support. Verify before external use.")

	return nil
}

func workflowSpecs() map[string]workflowSpec {
	specs := map[string]workflowSpec{}
	for _, route := range legalWorkflowRouteCatalog() {
		specs[route.Type] = workflowSpec{
			Type:              route.Type,
			DisplayName:       route.DisplayName,
			PipelineKey:       route.PipelineKey,
			OutputFile:        route.OutputFile,
			RequireSituation:  route.RequireSituation,
			RequireAudience:   route.RequireAudience,
			AllowDraft:        route.AllowDraft,
			AutoDraftFallback: route.AutoDraftFallback,
		}
	}
	return specs
}

func sortedWorkflowTypes() []string {
	types := make([]string, 0, len(legalWorkflowRouteCatalog()))
	for _, route := range legalWorkflowRouteCatalog() {
		types = append(types, route.Type)
	}
	sort.Strings(types)
	return types
}

func loadWorkflowInputs(spec workflowSpec, matter string, matterDir string, draftsDir string, opts workflowRunOptions) (workflowInputs, error) {
	inputs := workflowInputs{}

	situationPath, situationText, err := resolveOptionalMatterFile(opts.SituationArg, matterDir, []string{
		"situation.md",
		"triage.md",
		filepath.Join("input", "situation.md"),
		filepath.Join("working", "user_situation.md"),
	})
	if err != nil {
		return inputs, err
	}
	inputs.SituationPath = situationPath
	inputs.Situation = situationText

	goalPath, goalText, err := resolveOptionalMatterFile(opts.GoalArg, matterDir, []string{
		"goal.md",
		filepath.Join("input", "goal.md"),
	})
	if err != nil {
		return inputs, err
	}
	inputs.GoalPath = goalPath
	inputs.Goal = goalText

	questionsPath, questionsText, err := resolveOptionalMatterFile(opts.QuestionsArg, matterDir, []string{
		"questions.md",
		filepath.Join("input", "questions.md"),
	})
	if err != nil {
		return inputs, err
	}
	inputs.QuestionsPath = questionsPath
	inputs.Questions = questionsText

	if spec.AllowDraft {
		draftPath, draftName, draftText, err := resolveWorkflowDraft(opts.Draft, draftsDir, spec.AutoDraftFallback)
		if err != nil {
			return inputs, err
		}
		inputs.DraftPath = draftPath
		inputs.DraftName = draftName
		inputs.Draft = draftText
	}

	if spec.RequireSituation && strings.TrimSpace(inputs.Situation) == "" {
		return inputs, fmt.Errorf("%s requires a situation file.\nCreate one of:\n  matters/%s/situation.md\n  matters/%s/triage.md\n  matters/%s/input/situation.md\nOptional supporting files:\n  matters/%s/goal.md\n  matters/%s/questions.md", spec.DisplayName, matter, matter, matter, matter, matter)
	}

	if normalizeRouteName(spec.Type) == "authority-check" && strings.TrimSpace(inputs.Situation) == "" && strings.TrimSpace(inputs.Draft) == "" {
		return inputs, fmt.Errorf("Authority Check requires either a situation file or a draft.\nCreate matters/%s/situation.md or pass --draft path\\to\\draft.md", matter)
	}

	return inputs, nil
}

func resolveOptionalMatterFile(flagPath string, matterDir string, candidates []string) (string, string, error) {
	if strings.TrimSpace(flagPath) != "" {
		absPath, _ := filepath.Abs(flagPath)
		if data, err := os.ReadFile(absPath); err == nil {
			return absPath, string(data), nil
		}
		if pathExists(absPath) {
			return "", "", fmt.Errorf("read %s: %w", flagPath, os.ErrNotExist)
		}
		return "", flagPath, nil
	}

	for _, candidate := range candidates {
		path := filepath.Join(matterDir, candidate)
		data, err := os.ReadFile(path)
		if err == nil {
			return path, string(data), nil
		}
	}

	return "", "", nil
}

func resolveWorkflowDraft(specified string, draftsDir string, autoFallback bool) (string, string, string, error) {
	if strings.TrimSpace(specified) != "" {
		absPath, _ := filepath.Abs(specified)
		data, err := os.ReadFile(absPath)
		if err != nil {
			return "", "", "", fmt.Errorf("read draft %s: %w", specified, err)
		}
		return absPath, filepath.Base(absPath), string(data), nil
	}

	if !autoFallback {
		return "", "", "", nil
	}

	entries, err := os.ReadDir(draftsDir)
	if err != nil {
		return "", "", "", nil
	}

	type draftFile struct {
		path    string
		name    string
		modTime time.Time
	}
	var drafts []draftFile

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".txt" && ext != ".md" {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		drafts = append(drafts, draftFile{
			path:    filepath.Join(draftsDir, entry.Name()),
			name:    entry.Name(),
			modTime: info.ModTime(),
		})
	}

	if len(drafts) == 0 {
		return "", "", "", nil
	}

	sort.Slice(drafts, func(i, j int) bool {
		return drafts[i].modTime.After(drafts[j].modTime)
	})

	data, err := os.ReadFile(drafts[0].path)
	if err != nil {
		return "", "", "", fmt.Errorf("read draft %s: %w", drafts[0].path, err)
	}

	return drafts[0].path, drafts[0].name, string(data), nil
}

func buildInputManifest(inputDir string) string {
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return ""
	}

	var lines []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			lines = append(lines, fmt.Sprintf("- %s", entry.Name()))
			continue
		}
		lines = append(lines, fmt.Sprintf("- %s (%d bytes, modified %s)", entry.Name(), info.Size(), info.ModTime().Format("2006-01-02 15:04")))
	}
	return strings.Join(lines, "\n")
}

func buildWorkflowMeta(spec workflowSpec, matter string, audience string, urgency string, knowledgeAvailable bool, inputManifestAvailable bool) string {
	var b strings.Builder
	b.WriteString("=== WORKFLOW METADATA ===\n")
	b.WriteString("workflow_type: " + spec.Type + "\n")
	b.WriteString("workflow_label: " + spec.DisplayName + "\n")
	b.WriteString("matter: " + matter + "\n")
	if audience != "" {
		b.WriteString("audience: " + audience + "\n")
	}
	if urgency != "" {
		b.WriteString("urgency: " + urgency + "\n")
	}
	if knowledgeAvailable {
		b.WriteString("knowledge_layer: available\n")
	} else {
		b.WriteString("knowledge_layer: not_available\n")
	}
	if inputManifestAvailable {
		b.WriteString("input_manifest: available\n")
	} else {
		b.WriteString("input_manifest: not_available\n")
	}
	b.WriteString("decision_style: practical, opinionated, source-disciplined, honest about uncertainty\n")
	b.WriteString("=== END WORKFLOW METADATA ===")
	return b.String()
}

func buildWorkflowPrompt(contextMode string, prefix string, meta string, contractSummary string, contractDetail string, caseContext string, inputs workflowInputs, knowledgeContext string, inputManifest string, priorOutputs string) string {
	situation := valueOrFallback(inputs.Situation, "No situation file was provided.")
	goal := valueOrFallback(inputs.Goal, "No goal file was provided.")
	questions := valueOrFallback(inputs.Questions, "No questions file was provided.")
	draft := valueOrFallback(inputs.Draft, "No draft file was provided.")
	knowledge := valueOrFallback(knowledgeContext, "No matter knowledge layer is available. Flag evidence gaps clearly and avoid overclaiming.")
	manifest := valueOrFallback(inputManifest, "No matter input manifest is available.")
	contractBlock := ""
	if strings.TrimSpace(contractSummary) != "" {
		contractBlock += "\n\n" + contractSummary
	}
	if strings.TrimSpace(contractDetail) != "" {
		contractBlock += "\n\n" + contractDetail
	}

	switch contextMode {
	case "situation_plus_reviews":
		return prefix +
			"\n\n" + meta +
			contractBlock +
			"\n\n=== CASE CONTEXT ===\n" + caseContext + "\n=== END CASE CONTEXT ===" +
			"\n\n=== USER SITUATION ===\n" + situation + "\n=== END USER SITUATION ===" +
			"\n\n=== USER GOAL ===\n" + goal + "\n=== END USER GOAL ===" +
			"\n\n=== USER QUESTIONS ===\n" + questions + "\n=== END USER QUESTIONS ===" +
			"\n\n=== OPTIONAL DRAFT CONTEXT ===\n" + draft + "\n=== END OPTIONAL DRAFT CONTEXT ===" +
			"\n\n=== KNOWLEDGE LAYER ===\n" + knowledge + "\n=== END KNOWLEDGE LAYER ===" +
			"\n\n=== MATTER INPUT MANIFEST ===\n" + manifest + "\n=== END MATTER INPUT MANIFEST ===" +
			"\n\n=== PREVIOUS WORKFLOW FINDINGS ===\n" + priorOutputs + "\n=== END PREVIOUS WORKFLOW FINDINGS ==="
	case "all_reviews":
		return prefix +
			"\n\n" + meta +
			contractBlock +
			"\n\n=== CASE CONTEXT ===\n" + caseContext + "\n=== END CASE CONTEXT ===" +
			"\n\n=== USER SITUATION ===\n" + situation + "\n=== END USER SITUATION ===" +
			"\n\n=== USER GOAL ===\n" + goal + "\n=== END USER GOAL ===" +
			"\n\n=== USER QUESTIONS ===\n" + questions + "\n=== END USER QUESTIONS ===" +
			"\n\n=== OPTIONAL DRAFT CONTEXT ===\n" + draft + "\n=== END OPTIONAL DRAFT CONTEXT ===" +
			"\n\n=== ALL WORKFLOW FINDINGS ===\n" + priorOutputs + "\n=== END WORKFLOW FINDINGS ==="
	default:
		return prefix +
			"\n\n" + meta +
			contractBlock +
			"\n\n=== CASE CONTEXT ===\n" + caseContext + "\n=== END CASE CONTEXT ===" +
			"\n\n=== USER SITUATION ===\n" + situation + "\n=== END USER SITUATION ===" +
			"\n\n=== USER GOAL ===\n" + goal + "\n=== END USER GOAL ===" +
			"\n\n=== USER QUESTIONS ===\n" + questions + "\n=== END USER QUESTIONS ===" +
			"\n\n=== OPTIONAL DRAFT CONTEXT ===\n" + draft + "\n=== END OPTIONAL DRAFT CONTEXT ===" +
			"\n\n=== KNOWLEDGE LAYER ===\n" + knowledge + "\n=== END KNOWLEDGE LAYER ===" +
			"\n\n=== MATTER INPUT MANIFEST ===\n" + manifest + "\n=== END MATTER INPUT MANIFEST ==="
	}
}

func applyConditionalWorkflowRouting(steps []config.PipelineStep, spec workflowSpec, inputs workflowInputs, audience string, urgency string, childFlag bool, financialFlag bool) []config.PipelineStep {
	combined := strings.ToLower(strings.Join([]string{
		spec.Type,
		audience,
		urgency,
		inputs.Situation,
		inputs.Goal,
		inputs.Questions,
		inputs.Draft,
	}, "\n"))

	rules := legalConditionalRules()
	childDetected := childFlag || draftMentionsAny(combined, rules[0].Keywords)
	financialDetected := financialFlag || draftMentionsAny(combined, rules[1].Keywords)

	if childDetected {
		steps = insertBeforeFirstAvailable(steps, []string{
			"family-law-attorney-reviewer",
			"neutral-court-reader",
			"practical-resolution-reviewer",
			"managing-partner-final-synthesizer",
		}, config.PipelineStep{
			Name:         "Child Best Interests and Family Dynamics Review",
			Persona:      "child-best-interests-family-dynamics-reviewer",
			LLM:          "claude",
			ModelProfile: "review_reasoning",
			ContextMode:  "situation_plus_knowledge",
			PromptPrefix: "You are the Child Best Interests and Family Dynamics Reviewer. Distinguish child welfare significance, loyalty-conflict risk, source strength, escalation value, and pattern-preservation value. Follow your agent instructions exactly.",
		})
	}

	if financialDetected {
		steps = insertBeforeFirstAvailable(steps, []string{
			"family-law-attorney-reviewer",
			"neutral-court-reader",
			"practical-resolution-reviewer",
			"managing-partner-final-synthesizer",
		}, config.PipelineStep{
			Name:         "Financial Support Review",
			Persona:      "financial-support-reviewer",
			LLM:          "claude",
			ModelProfile: "review_reasoning",
			ContextMode:  "situation_plus_knowledge",
			PromptPrefix: "You are the Financial Support Reviewer. Check billing, support, reimbursement, income, expense, and documentation issues for strength, gaps, and escalation value. Follow your agent instructions exactly.",
		})
	}

	for i := range steps {
		steps[i].Step = i + 1
	}

	return steps
}

func insertBeforeFirstAvailable(steps []config.PipelineStep, targetPersonas []string, step config.PipelineStep) []config.PipelineStep {
	for _, existing := range steps {
		if existing.Persona == step.Persona {
			return steps
		}
	}
	for _, target := range targetPersonas {
		for i, existing := range steps {
			if existing.Persona == target {
				steps = append(steps, config.PipelineStep{})
				copy(steps[i+1:], steps[i:])
				steps[i] = step
				return steps
			}
		}
	}
	return insertBeforeFinal(steps, step)
}

func valueOrFallback(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
