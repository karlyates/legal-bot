package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type guidedChoice struct {
	ID              int
	Label           string
	Kind            string
	WorkflowType    string
	DefaultAudience string
	DefaultUrgency  string
	DefaultGoal     string
}

type guidedRunPlan struct {
	Matter         string
	Kind           string
	WorkflowType   string
	ReviewMode     string
	DocumentType   string
	Level          string
	Audience       string
	Urgency        string
	Situation      string
	Goal           string
	Questions      string
	DraftPath      string
	ChildRelated   bool
	Financial      bool
	DryRun         bool
	OutputPath     string
	GuidedFilePath string
	Choice         guidedChoice
}

type matterStatus struct {
	Matter          string
	MatterDir       string
	InputDir        string
	DraftsDir       string
	KnowledgeDir    string
	WorkingDir      string
	OutputDir       string
	InputFiles      []string
	DraftFiles      []string
	NewestDraft     string
	HasKnowledge    bool
	MatterExists    bool
	InputDirExists  bool
	DraftsDirExists bool
}

var (
	guideChoiceFlag     int
	guideIntentFlag     string
	guideListFlag       bool
	guideDryRunFlag     bool
	guideSituationFlag  string
	guideGoalFlag       string
	guideUrgencyFlag    string
	guideAudienceFlag   string
	guideQuestionsFlag  string
	guideDraftFlag      string
	guideReviewModeFlag string
	guideDocumentType   string
	guideLevelFlag      string
)

var guideCmd = &cobra.Command{
	Use:     "guide [matter]",
	Aliases: []string{"start"},
	Short:   "Guided front door for intake, workflow, and review routes",
	Args:    cobra.MaximumNArgs(1),
	RunE:    runGuide,
}

func init() {
	guideCmd.Flags().IntVar(&guideChoiceFlag, "choice", 0, "Guided menu choice number")
	guideCmd.Flags().StringVar(&guideIntentFlag, "intent", "", "Intent alias such as triage, counsel-brief, review, or intake")
	guideCmd.Flags().BoolVar(&guideListFlag, "list", false, "List guided choices and exit")
	guideCmd.Flags().BoolVar(&guideDryRunFlag, "dry-run", false, "Print the planned route without running it")
	guideCmd.Flags().StringVar(&guideSituationFlag, "situation", "", "Guided situation or topic text")
	guideCmd.Flags().StringVar(&guideGoalFlag, "goal", "", "Guided goal text")
	guideCmd.Flags().StringVar(&guideUrgencyFlag, "urgency", "", "Urgency: low, normal, high, emergency")
	guideCmd.Flags().StringVar(&guideAudienceFlag, "audience", "", "Audience hint for guided routes")
	guideCmd.Flags().StringVar(&guideQuestionsFlag, "questions", "", "Specific questions to answer")
	guideCmd.Flags().StringVar(&guideDraftFlag, "draft", "", "Draft path for review or draft-aware workflows")
	guideCmd.Flags().StringVar(&guideReviewModeFlag, "mode", "", "Review mode for guided review")
	guideCmd.Flags().StringVar(&guideDocumentType, "document-type", "", "Review document type for guided review")
	guideCmd.Flags().StringVar(&guideLevelFlag, "level", string(analysisLevelScout), "Analysis level for routed intake, review, or workflow commands")
	rootCmd.AddCommand(guideCmd)
}

func runGuide(cmd *cobra.Command, args []string) error {
	if guideListFlag {
		printGuidedChoices(cmd.OutOrStdout())
		return nil
	}

	projectRoot := resolveProjectRoot()
	reader := bufio.NewReader(cmd.InOrStdin())
	writer := cmd.OutOrStdout()
	interactive := guideChoiceFlag == 0 && strings.TrimSpace(guideIntentFlag) == "" && len(args) == 0

	matter := ""
	if len(args) > 0 {
		matter = strings.TrimSpace(args[0])
	}
	if matter == "" {
		var err error
		matter, err = promptForMatter(reader, writer, projectRoot)
		if err != nil {
			return err
		}
	}

	status := detectMatterStatus(projectRoot, matter)
	if !status.MatterExists {
		if !confirmWithDefault(reader, writer, fmt.Sprintf("Matter %q does not exist. Create the standard matter layout? [Y/n] ", matter), true) {
			return fmt.Errorf("matter %q does not exist", matter)
		}
		if err := createMatterLayout(status); err != nil {
			return err
		}
		status = detectMatterStatus(projectRoot, matter)
	}

	choice, err := resolveGuidedChoice(guideChoiceFlag, guideIntentFlag)
	if err != nil && !interactive {
		return err
	}
	if interactive || choice.ID == 0 {
		choice, err = promptForChoice(reader, writer)
		if err != nil {
			return err
		}
	}

	plan, err := buildGuidedPlan(reader, writer, status, choice)
	if err != nil {
		return err
	}
	plan.DryRun = guideDryRunFlag

	preRunAction, err := handleMatterReadiness(reader, writer, status, &plan)
	if err != nil {
		return err
	}

	if !guideDryRunFlag {
		if err := maybeWriteGuidedRequest(status, &plan); err != nil {
			return err
		}
	}

	printPlan(writer, plan)
	if guideDryRunFlag {
		return nil
	}

	runNow := true
	if interactive {
		runNow = confirmWithDefault(reader, writer, "Run this now? [Y/n] ", true)
	}
	if !runNow {
		fmt.Fprintf(writer, "\nEquivalent command:\n%s\n", guidedEquivalentCommand(plan))
		return nil
	}

	if preRunAction == "intake" {
		if err := executeIntake(intakeOptions{Matter: plan.Matter, Level: plan.Level}); err != nil {
			return err
		}
	}

	switch plan.Kind {
	case "intake":
		return executeIntake(intakeOptions{Matter: plan.Matter, Level: plan.Level})
	case "review":
		return executeReview(reviewOptions{
			Matter:       plan.Matter,
			Draft:        plan.DraftPath,
			Mode:         plan.ReviewMode,
			DocumentType: plan.DocumentType,
			ChildRelated: plan.ChildRelated,
			Financial:    plan.Financial,
			Level:        plan.Level,
		})
	case "workflow":
		return executeWorkflow(workflowRunOptions{
			Matter:       plan.Matter,
			Type:         plan.WorkflowType,
			Audience:     plan.Audience,
			Urgency:      plan.Urgency,
			SituationArg: inlineOrPath(plan.Situation),
			GoalArg:      inlineOrPath(plan.Goal),
			QuestionsArg: inlineOrPath(plan.Questions),
			Draft:        plan.DraftPath,
			ChildRelated: plan.ChildRelated,
			Financial:    plan.Financial,
			Level:        plan.Level,
		})
	default:
		return fmt.Errorf("unknown guided route kind %q", plan.Kind)
	}
}

func guidedChoices() []guidedChoice {
	triageRoute, _ := workflowRouteCatalogByType(workflowTypeTriage)
	evidenceRoute, _ := workflowRouteCatalogByType(workflowTypeEvidencePacket)
	counselRoute, _ := workflowRouteCatalogByType(workflowTypeCounselBrief)
	smRoute, _ := workflowRouteCatalogByType(workflowTypeSMTriage)
	messageRoute, _ := workflowRouteCatalogByType(workflowTypeDraftMessage)
	patternRoute, _ := workflowRouteCatalogByType(workflowTypePatternReview)
	factLockRoute, _ := workflowRouteCatalogByType(workflowTypeFactLock)
	authorityRoute, _ := workflowRouteCatalogByType(workflowTypeAuthorityCheck)
	prepRoute, _ := workflowRouteCatalogByType(workflowTypePrep)
	journalRoute, _ := workflowRouteCatalogByType(workflowTypeJournalEntry)

	return []guidedChoice{
		{ID: 1, Label: "Understand a new situation", Kind: routeKindWorkflow, WorkflowType: workflowTypeTriage, DefaultAudience: triageRoute.DefaultAudience, DefaultUrgency: "normal", DefaultGoal: triageRoute.DefaultGoal},
		{ID: 2, Label: "Decide whether to escalate, preserve, or let something go", Kind: routeKindWorkflow, WorkflowType: workflowTypeTriage, DefaultAudience: triageRoute.DefaultAudience, DefaultUrgency: "normal", DefaultGoal: "Determine whether this issue should be pursued now, preserved for pattern evidence, handled with a limited record-building message, raised with counsel, raised with the Special Master, or let go."},
		{ID: 3, Label: "Build an evidence packet", Kind: routeKindWorkflow, WorkflowType: workflowTypeEvidencePacket, DefaultAudience: evidenceRoute.DefaultAudience, DefaultUrgency: "normal", DefaultGoal: evidenceRoute.DefaultGoal},
		{ID: 4, Label: "Write to Kaitlyn/Jessika", Kind: routeKindWorkflow, WorkflowType: workflowTypeCounselBrief, DefaultAudience: counselRoute.DefaultAudience, DefaultUrgency: "normal", DefaultGoal: counselRoute.DefaultGoal},
		{ID: 5, Label: "Review a draft from counsel", Kind: routeKindReview},
		{ID: 6, Label: "Decide whether to bring something to the Special Master", Kind: routeKindWorkflow, WorkflowType: workflowTypeSMTriage, DefaultAudience: smRoute.DefaultAudience, DefaultUrgency: "normal", DefaultGoal: smRoute.DefaultGoal},
		{ID: 7, Label: "Write an OFW / SM / provider / school / counsel message", Kind: routeKindWorkflow, WorkflowType: workflowTypeDraftMessage, DefaultAudience: messageRoute.DefaultAudience, DefaultUrgency: "normal", DefaultGoal: messageRoute.DefaultGoal},
		{ID: 8, Label: "Analyze whether recurring events add up to a pattern", Kind: routeKindWorkflow, WorkflowType: workflowTypePatternReview, DefaultAudience: patternRoute.DefaultAudience, DefaultUrgency: "normal", DefaultGoal: patternRoute.DefaultGoal},
		{ID: 9, Label: "Lock down facts before external use", Kind: routeKindWorkflow, WorkflowType: workflowTypeFactLock, DefaultAudience: factLockRoute.DefaultAudience, DefaultUrgency: "normal", DefaultGoal: factLockRoute.DefaultGoal},
		{ID: 10, Label: "Check legal authority gaps", Kind: routeKindWorkflow, WorkflowType: workflowTypeAuthorityCheck, DefaultAudience: authorityRoute.DefaultAudience, DefaultUrgency: "normal", DefaultGoal: authorityRoute.DefaultGoal},
		{ID: 11, Label: "Prepare for a call/hearing/SM conference", Kind: routeKindWorkflow, WorkflowType: workflowTypePrep, DefaultAudience: prepRoute.DefaultAudience, DefaultUrgency: "normal", DefaultGoal: prepRoute.DefaultGoal},
		{ID: 12, Label: "Create a journal entry", Kind: routeKindWorkflow, WorkflowType: workflowTypeJournalEntry, DefaultAudience: journalRoute.DefaultAudience, DefaultUrgency: "normal", DefaultGoal: journalRoute.DefaultGoal},
		{ID: 13, Label: "Run intake on new source material", Kind: routeKindIntake},
	}
}

func guidedIntentAliases() map[string]int {
	return map[string]int{
		"new-situation": 1, "triage": 1,
		"pursue-preserve-let-go": 2,
		"evidence":               3, "evidence-packet": 3,
		"counsel": 4, "counsel-brief": 4,
		"draft-review": 5, "review": 5,
		"special-master": 6, "sm": 6, "sm-triage": 6,
		"message": 7, "draft-message": 7,
		"pattern": 8, "pattern-review": 8,
		"fact-lock": 9,
		"authority": 10, "authority-check": 10,
		"prep":    11,
		"journal": 12, "journal-entry": 12,
		"intake": 13,
	}
}

func resolveGuidedChoice(choiceID int, intent string) (guidedChoice, error) {
	if strings.TrimSpace(intent) != "" {
		id, ok := guidedIntentAliases()[normalizeRouteName(intent)]
		if !ok {
			return guidedChoice{}, fmt.Errorf("unknown intent %q", intent)
		}
		choiceID = id
	}
	for _, choice := range guidedChoices() {
		if choice.ID == choiceID {
			return choice, nil
		}
	}
	return guidedChoice{}, fmt.Errorf("unknown choice %d", choiceID)
}

func printGuidedChoices(w io.Writer) {
	fmt.Fprintln(w, "Guided choices:")
	for _, choice := range guidedChoices() {
		target := choice.Kind
		if choice.WorkflowType != "" {
			target = "workflow " + choice.WorkflowType
		}
		if choice.ID == 5 {
			target = "review"
		}
		if choice.ID == 13 {
			target = "intake"
		}
		fmt.Fprintf(w, "%d. %s -> %s\n", choice.ID, choice.Label, target)
	}
}

func promptForMatter(reader *bufio.Reader, w io.Writer, projectRoot string) (string, error) {
	matterNames := listMatterNames(filepath.Join(projectRoot, "matters"))
	if len(matterNames) > 0 {
		fmt.Fprintln(w, "Known matters:")
		for _, name := range matterNames {
			fmt.Fprintf(w, "- %s\n", name)
		}
	}
	return promptText(reader, w, "Matter name? ", false)
}

func promptForChoice(reader *bufio.Reader, w io.Writer) (guidedChoice, error) {
	fmt.Fprintln(w, "What are you trying to do?")
	for _, choice := range guidedChoices() {
		fmt.Fprintf(w, "%d. %s\n", choice.ID, choice.Label)
	}
	for {
		value, err := promptText(reader, w, "Choice number? ", false)
		if err != nil {
			return guidedChoice{}, err
		}
		num, convErr := strconv.Atoi(value)
		if convErr == nil {
			if choice, err := resolveGuidedChoice(num, ""); err == nil {
				return choice, nil
			}
		}
		fmt.Fprintln(w, "Please enter a number from 1 to 13.")
	}
}

func buildGuidedPlan(reader *bufio.Reader, w io.Writer, status matterStatus, choice guidedChoice) (guidedRunPlan, error) {
	plan := guidedRunPlan{
		Matter:       status.Matter,
		Kind:         choice.Kind,
		WorkflowType: choice.WorkflowType,
		Level:        guideLevelFlag,
		Audience:     choice.DefaultAudience,
		Urgency:      choice.DefaultUrgency,
		Goal:         choice.DefaultGoal,
		Choice:       choice,
	}
	if plan.Urgency == "" {
		plan.Urgency = "normal"
	}

	switch choice.ID {
	case 1:
		plan.Situation = valueOrPrompt(reader, w, guideSituationFlag, "What happened or what are you trying to analyze? ")
		plan.Goal = valueOrPromptDefault(reader, w, guideGoalFlag, "What outcome do you want from this run? ", plan.Goal)
		plan.Urgency = valueOrPromptDefault(reader, w, guideUrgencyFlag, "Urgency? [low/normal/high/emergency] ", plan.Urgency)
		plan.Audience = valueOrPromptDefault(reader, w, guideAudienceFlag, "Audience? [internal/counsel/OFW/Special Master/provider/school/court] ", plan.Audience)
		plan.ChildRelated = promptBool(reader, w, "Is this child-related? [Y/n] ", true)
		plan.Financial = promptBool(reader, w, "Is this financial/support-related? [y/N] ", false)
	case 2:
		plan.Situation = valueOrPrompt(reader, w, guideSituationFlag, "What happened or what are you trying to analyze? ")
		plan.Urgency = valueOrPromptDefault(reader, w, guideUrgencyFlag, "Urgency? [low/normal/high/emergency] ", plan.Urgency)
		plan.ChildRelated = promptBool(reader, w, "Is this child-related? [Y/n] ", true)
		plan.Financial = promptBool(reader, w, "Is this financial/support-related? [y/N] ", false)
	case 3:
		plan.Situation = valueOrPrompt(reader, w, guideSituationFlag, "Situation or topic? ")
		plan.Goal = valueOrPromptDefault(reader, w, guideGoalFlag, "What outcome do you want from this run? ", plan.Goal)
		plan.Urgency = valueOrPromptDefault(reader, w, guideUrgencyFlag, "Urgency? [low/normal/high/emergency] ", plan.Urgency)
		inputSummary := "Current input/ files: none detected."
		if len(status.InputFiles) > 0 {
			inputSummary = "Current input/ files: " + strings.Join(status.InputFiles, ", ")
		}
		packetType := valueOrPrompt(reader, w, "", "What output is needed? [attorney packet / SM packet / court packet / journal support / internal review] ")
		plan.Questions = strings.TrimSpace(inputSummary + "\nRequested output: " + packetType)
	case 4:
		target := valueOrPrompt(reader, w, "", "Is this for Kaitlyn, Jessika, or both? ")
		plan.Situation = valueOrPrompt(reader, w, guideSituationFlag, "What happened or what are you trying to analyze? ")
		plan.Goal = valueOrPromptDefault(reader, w, guideGoalFlag, "What outcome do you want from this run? ", plan.Goal)
		plan.Urgency = valueOrPromptDefault(reader, w, guideUrgencyFlag, "Urgency? [low/normal/high/emergency] ", plan.Urgency)
		format := valueOrPrompt(reader, w, "", "What format? [concise email / call agenda / issue list / evidence request list] ")
		plan.Questions = strings.TrimSpace("Counsel target: " + target + "\nRequested format: " + format)
	case 5:
		plan.ReviewMode = "standard"
		plan.DraftPath = resolveGuideDraft(status, guideDraftFlag)
		if plan.DraftPath == "" {
			useLatest := status.NewestDraft != ""
			if useLatest && confirmWithDefault(reader, w, fmt.Sprintf("Use latest draft in drafts/? [%s] ", filepath.Base(status.NewestDraft)), true) {
				plan.DraftPath = status.NewestDraft
			} else {
				plan.DraftPath = valueOrPrompt(reader, w, "", "Draft path? ")
			}
		}
		plan.DocumentType = normalizeRouteName(valueOrPromptDefault(reader, w, guideDocumentType, "Document type? [motion/opposition/response/reply/declaration/proposed-order/co-parenting-communication/unknown-general] ", "unknown-general"))
		if plan.DocumentType == "unknown-general" || plan.DocumentType == "unknown" || plan.DocumentType == "general" {
			plan.DocumentType = ""
			plan.ReviewMode = normalizeRouteName(valueOrPromptDefault(reader, w, guideReviewModeFlag, "Mode? [quick/standard/deep/legal-research/strategy] ", "standard"))
		}
		plan.ChildRelated = promptBool(reader, w, "Is this child-related? [Y/n] ", true)
		plan.Financial = promptBool(reader, w, "Is this financial/support-related? [y/N] ", false)
	case 6:
		plan.Situation = valueOrPrompt(reader, w, guideSituationFlag, "What happened or what are you trying to analyze? ")
		directive := valueOrPrompt(reader, w, "", "What directive do you want, if any? ")
		plan.Goal = strings.TrimSpace(plan.Goal + "\nDesired directive: " + directive)
		plan.Urgency = valueOrPromptDefault(reader, w, guideUrgencyFlag, "Urgency? [low/normal/high/emergency] ", plan.Urgency)
		prior := valueOrPrompt(reader, w, "", "Has there already been an OFW attempt or attorney discussion? ")
		plan.Questions = strings.TrimSpace("Prior OFW/attorney step: " + prior)
		plan.ChildRelated = promptBool(reader, w, "Is this child-related? [Y/n] ", true)
		plan.Financial = promptBool(reader, w, "Is this financial/support-related? [y/N] ", false)
	case 7:
		plan.Audience = normalizeRouteName(valueOrPromptDefault(reader, w, guideAudienceFlag, "Who is the audience? [OFW/Emily, Special Master, provider, school, counsel, court-facing, internal] ", "ofw"))
		plan.Situation = valueOrPrompt(reader, w, guideSituationFlag, "What happened or what are you trying to analyze? ")
		plan.Goal = valueOrPromptDefault(reader, w, guideGoalFlag, "What outcome do you want from this run? ", plan.Goal)
		plan.Urgency = valueOrPromptDefault(reader, w, guideUrgencyFlag, "Urgency? [low/normal/high/emergency] ", plan.Urgency)
		plan.Questions = "Requested message type: " + valueOrPrompt(reader, w, "", "What kind of message? [brief record-building / direct request / response to accusation / scheduling-logistics / provider access-request / Special Master directive request] ")
		plan.ChildRelated = promptBool(reader, w, "Is this child-related? [Y/n] ", true)
		plan.Financial = promptBool(reader, w, "Is this financial/support-related? [y/N] ", false)
	case 8:
		plan.Situation = valueOrPrompt(reader, w, guideSituationFlag, "Pattern or topic? ")
		dateRange := valueOrPrompt(reader, w, "", "Approximate date range? ")
		related := valueOrPrompt(reader, w, "", "Related issues or events? ")
		goalUse := valueOrPrompt(reader, w, "", "Goal? [MTE / SM / counsel briefing / journal cleanup / internal triage] ")
		plan.Goal = "Primary use: " + goalUse
		plan.Questions = strings.TrimSpace("Date range: " + dateRange + "\nRelated issues/events: " + related)
		plan.ChildRelated = promptBool(reader, w, "Is this child-related? [Y/n] ", true)
		plan.Financial = promptBool(reader, w, "Is this financial/support-related? [y/N] ", false)
	case 9:
		plan.Situation = valueOrPrompt(reader, w, guideSituationFlag, "What facts or claims need to be verified? ")
		plan.Audience = valueOrPromptDefault(reader, w, guideAudienceFlag, "Intended external use? [counsel email / OFW / Special Master / court filing / provider-school / internal only] ", "internal")
		authority := valueOrPrompt(reader, w, "", "Are authority or legal claims involved? ")
		plan.Questions = "Authority or legal claims involved: " + authority
		plan.ChildRelated = promptBool(reader, w, "Is this child-related? [Y/n] ", true)
		plan.Financial = promptBool(reader, w, "Is this financial/support-related? [y/N] ", false)
	case 10:
		plan.Situation = valueOrPrompt(reader, w, guideSituationFlag, "Legal issue or question? ")
		jurisdiction := valueOrPrompt(reader, w, "", "Jurisdiction if relevant? ")
		use := valueOrPrompt(reader, w, "", "Intended use? [attorney question / motion review / SM argument / internal triage] ")
		plan.Goal = valueOrPromptDefault(reader, w, guideGoalFlag, "What outcome do you want from this run? ", plan.Goal)
		plan.Questions = strings.TrimSpace("Jurisdiction: " + jurisdiction + "\nIntended use: " + use + "\nDo not overclaim live currentness or citator status.")
	case 11:
		prepType := valueOrPrompt(reader, w, "", "Prep type? [attorney call / hearing / SM conference / mediation / provider-school meeting] ")
		plan.Situation = valueOrPrompt(reader, w, guideSituationFlag, "Situation or topic? ")
		plan.Goal = valueOrPromptDefault(reader, w, guideGoalFlag, "Goal or desired outcome? ", plan.Goal)
		when := valueOrPrompt(reader, w, "", "Urgency or date if known? ")
		plan.Urgency = valueOrPromptDefault(reader, w, guideUrgencyFlag, "Urgency? [low/normal/high/emergency] ", plan.Urgency)
		plan.Questions = strings.TrimSpace("Prep type: " + prepType + "\nDate or timing: " + when)
		plan.ChildRelated = promptBool(reader, w, "Is this child-related? [Y/n] ", true)
		plan.Financial = promptBool(reader, w, "Is this financial/support-related? [y/N] ", false)
	case 12:
		plan.Situation = valueOrPrompt(reader, w, guideSituationFlag, "What happened? ")
		dateRange := valueOrPrompt(reader, w, "", "Date or date range? ")
		anchors := valueOrPrompt(reader, w, "", "Source anchors or evidence if known? ")
		category := valueOrPrompt(reader, w, "", "Category or related issue if known? ")
		priv := valueOrPrompt(reader, w, "", "Should this be internal-only or privileged? ")
		plan.Questions = strings.TrimSpace("Date/date range: " + dateRange + "\nSource anchors/evidence: " + anchors + "\nCategory/issue: " + category + "\nInternal-only/privileged: " + priv)
		plan.ChildRelated = promptBool(reader, w, "Is this child-related? [Y/n] ", true)
		plan.Financial = promptBool(reader, w, "Is this financial/support-related? [y/N] ", false)
	case 13:
		plan.OutputPath = filepath.Join(status.OutputDir, "ingestion_report.md")
	default:
		return plan, fmt.Errorf("unsupported guided choice %d", choice.ID)
	}

	if plan.Kind == "workflow" {
		spec := workflowSpecs()[plan.WorkflowType]
		plan.OutputPath = filepath.Join(status.OutputDir, spec.OutputFile)
	}
	if plan.Kind == "review" {
		plan.OutputPath = filepath.Join(status.OutputDir, "final_review_packet.md")
	}
	return plan, nil
}

func detectMatterStatus(projectRoot string, matter string) matterStatus {
	matterDir := filepath.Join(projectRoot, "matters", matter)
	inputDir := filepath.Join(matterDir, "input")
	draftsDir := filepath.Join(matterDir, "drafts")
	knowledgeDir := filepath.Join(matterDir, "knowledge")
	workingDir := filepath.Join(matterDir, "working")
	outputDir := filepath.Join(matterDir, "output")

	status := matterStatus{
		Matter:          matter,
		MatterDir:       matterDir,
		InputDir:        inputDir,
		DraftsDir:       draftsDir,
		KnowledgeDir:    knowledgeDir,
		WorkingDir:      workingDir,
		OutputDir:       outputDir,
		MatterExists:    pathExists(matterDir),
		InputDirExists:  pathExists(inputDir),
		DraftsDirExists: pathExists(draftsDir),
		HasKnowledge:    pathExists(filepath.Join(knowledgeDir, "document_register.csv")),
		InputFiles:      listTextFiles(inputDir),
		DraftFiles:      listTextFiles(draftsDir),
	}
	if len(status.DraftFiles) > 0 {
		status.NewestDraft = newestTextFile(draftsDir)
	}
	return status
}

func handleMatterReadiness(reader *bufio.Reader, w io.Writer, status matterStatus, plan *guidedRunPlan) (string, error) {
	switch plan.Kind {
	case "intake":
		if len(status.InputFiles) == 0 {
			return "", fmt.Errorf("no .md or .txt files found in matters/%s/input/\nPut source files there and run intake again", status.Matter)
		}
	case "review":
		if !status.HasKnowledge {
			if len(status.InputFiles) == 0 {
				return "", fmt.Errorf("review expects intake-generated knowledge, but matters/%s/knowledge/document_register.csv is missing.\nPut source files in matters/%s/input/ and run intake first", status.Matter, status.Matter)
			}
			if confirmWithDefault(reader, w, "No knowledge layer found. Run intake first? [Y/n] ", true) {
				return "intake", nil
			}
			return "", fmt.Errorf("review cancelled because the knowledge layer is missing")
		}
		if strings.TrimSpace(plan.DraftPath) == "" {
			return "", fmt.Errorf("no draft selected")
		}
		if !pathExists(plan.DraftPath) {
			return "", fmt.Errorf("draft path not found: %s", plan.DraftPath)
		}
	case "workflow":
		if !status.HasKnowledge {
			fmt.Fprintln(w, "Warning: no intake-generated knowledge layer found. This workflow can continue, but source coverage will be weaker.")
			if len(status.InputFiles) > 0 && confirmWithDefault(reader, w, "Run intake first? [y/N] ", false) {
				return "intake", nil
			}
		}
		if plan.Urgency == "emergency" {
			fmt.Fprintln(w, "Emergency urgency selected. Legal-Bot can prepare analysis, but court/attorney timing should be handled outside the tool.")
		}
	}
	return "", nil
}

func maybeWriteGuidedRequest(status matterStatus, plan *guidedRunPlan) error {
	if plan.DryRun {
		return nil
	}
	if plan.Kind != "workflow" {
		return nil
	}
	if strings.TrimSpace(plan.Situation) == "" && strings.TrimSpace(plan.Goal) == "" && strings.TrimSpace(plan.Questions) == "" {
		return nil
	}
	if err := os.MkdirAll(status.WorkingDir, 0755); err != nil {
		return err
	}
	ts := time.Now().Format("20060102_150405")
	path := filepath.Join(status.WorkingDir, fmt.Sprintf("guided_%s_%s.md", ts, plan.WorkflowType))
	body := fmt.Sprintf(`# Guided Legal-Bot Request

Matter: %s
Workflow / Route: %s
Audience: %s
Urgency: %s
Child-related: %t
Financial: %t
Situation:
%s

Goal:
%s

Questions:
%s

Draft path: %s
Created by: legal-bot guide
Created at: %s
`, plan.Matter, plan.WorkflowType, plan.Audience, plan.Urgency, plan.ChildRelated, plan.Financial, plan.Situation, plan.Goal, plan.Questions, plan.DraftPath, time.Now().Format(time.RFC3339))
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		return err
	}
	plan.GuidedFilePath = path
	return nil
}

func printPlan(w io.Writer, plan guidedRunPlan) {
	fmt.Fprintln(w, "\nPlanned Legal-Bot run:")
	fmt.Fprintf(w, "- Matter: %s\n", plan.Matter)
	fmt.Fprintf(w, "- Route: %s\n", plan.Kind)
	if plan.Level != "" {
		fmt.Fprintf(w, "- Level: %s\n", plan.Level)
	}
	if plan.WorkflowType != "" {
		fmt.Fprintf(w, "- Type: %s\n", plan.WorkflowType)
	}
	if plan.DocumentType != "" {
		fmt.Fprintf(w, "- Document type: %s\n", plan.DocumentType)
	}
	if plan.ReviewMode != "" {
		fmt.Fprintf(w, "- Review mode: %s\n", plan.ReviewMode)
	}
	if plan.Audience != "" {
		fmt.Fprintf(w, "- Audience: %s\n", plan.Audience)
	}
	if plan.Urgency != "" && plan.Kind == "workflow" {
		fmt.Fprintf(w, "- Urgency: %s\n", plan.Urgency)
	}
	fmt.Fprintf(w, "- Child-related: %s\n", yesNo(plan.ChildRelated))
	fmt.Fprintf(w, "- Financial: %s\n", yesNo(plan.Financial))
	if plan.DraftPath != "" {
		fmt.Fprintf(w, "- Draft path: %s\n", plan.DraftPath)
	}
	if plan.Situation != "" {
		fmt.Fprintln(w, "- Situation source: interactive input")
	}
	if plan.OutputPath != "" {
		fmt.Fprintf(w, "- Output: %s\n", filepath.ToSlash(plan.OutputPath))
	}
	if plan.GuidedFilePath != "" {
		fmt.Fprintf(w, "- Guided input saved: %s\n", filepath.ToSlash(plan.GuidedFilePath))
	}
	fmt.Fprintf(w, "\nEquivalent command:\n%s\n", guidedEquivalentCommand(plan))
}

func guidedEquivalentCommand(plan guidedRunPlan) string {
	var parts []string
	parts = append(parts, "legal-bot")
	switch plan.Kind {
	case "intake":
		parts = append(parts, "intake", plan.Matter)
	case "review":
		parts = append(parts, "review", plan.Matter)
		if plan.Level != "" {
			parts = append(parts, "--level", plan.Level)
		}
		if plan.ReviewMode != "" {
			parts = append(parts, "--mode", plan.ReviewMode)
		}
		if plan.DocumentType != "" {
			parts = append(parts, "--document-type", plan.DocumentType)
		}
		if plan.DraftPath != "" {
			parts = append(parts, "--draft", quoteIfNeeded(plan.DraftPath))
		}
	case "workflow":
		parts = append(parts, "workflow", plan.Matter, "--type", plan.WorkflowType)
		if plan.Level != "" {
			parts = append(parts, "--level", plan.Level)
		}
		if plan.Audience != "" {
			parts = append(parts, "--audience", quoteIfNeeded(plan.Audience))
		}
		if plan.Urgency != "" {
			parts = append(parts, "--urgency", plan.Urgency)
		}
		if plan.DraftPath != "" {
			parts = append(parts, "--draft", quoteIfNeeded(plan.DraftPath))
		}
	}
	if plan.Kind == "intake" && plan.Level != "" {
		parts = append(parts, "--level", plan.Level)
	}
	if plan.ChildRelated {
		parts = append(parts, "--child-related")
	}
	if plan.Financial {
		parts = append(parts, "--financial")
	}
	return strings.Join(parts, " ")
}

func listMatterNames(mattersDir string) []string {
	entries, err := os.ReadDir(mattersDir)
	if err != nil {
		return nil
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	slices.Sort(names)
	return names
}

func createMatterLayout(status matterStatus) error {
	for _, dir := range []string{status.InputDir, status.DraftsDir, status.KnowledgeDir, status.WorkingDir, status.OutputDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func listTextFiles(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".md" || ext == ".txt" {
			files = append(files, entry.Name())
		}
	}
	slices.Sort(files)
	return files
}

func newestTextFile(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	var newest string
	var newestTime time.Time
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".md" && ext != ".txt" {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if newest == "" || info.ModTime().After(newestTime) {
			newest = filepath.Join(dir, entry.Name())
			newestTime = info.ModTime()
		}
	}
	return newest
}

func resolveGuideDraft(status matterStatus, value string) string {
	if strings.TrimSpace(value) != "" {
		return value
	}
	return ""
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func promptText(reader *bufio.Reader, w io.Writer, label string, allowEmpty bool) (string, error) {
	for {
		fmt.Fprint(w, label)
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		value := strings.TrimSpace(line)
		if value != "" || allowEmpty || err == io.EOF {
			return value, nil
		}
	}
}

func valueOrPrompt(reader *bufio.Reader, w io.Writer, current string, label string) string {
	if strings.TrimSpace(current) != "" {
		return strings.TrimSpace(current)
	}
	value, _ := promptText(reader, w, label, false)
	return value
}

func valueOrPromptDefault(reader *bufio.Reader, w io.Writer, current string, label string, defaultValue string) string {
	if strings.TrimSpace(current) != "" {
		return strings.TrimSpace(current)
	}
	value, _ := promptText(reader, w, label, true)
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}

func promptBool(reader *bufio.Reader, w io.Writer, label string, defaultValue bool) bool {
	value, _ := promptText(reader, w, label, true)
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return defaultValue
	}
	return value == "y" || value == "yes"
}

func confirmWithDefault(reader *bufio.Reader, w io.Writer, label string, defaultValue bool) bool {
	return promptBool(reader, w, label, defaultValue)
}

func yesNo(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}

func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t") {
		return strconv.Quote(v)
	}
	return v
}

func inlineOrPath(value string) string {
	return value
}
