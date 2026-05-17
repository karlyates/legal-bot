package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jdonohoo/legal-bot/go/internal/config"
	"github.com/jdonohoo/legal-bot/go/internal/embedded"
	personapkg "github.com/jdonohoo/legal-bot/go/internal/persona"
)

const (
	routeKindIntake   = "intake"
	routeKindReview   = "review"
	routeKindWorkflow = "workflow"

	reviewModeQuick         = "quick"
	reviewModeStandard      = "standard"
	reviewModeDeep          = "deep"
	reviewModeLegalResearch = "legal-research"
	reviewModeStrategy      = "strategy"

	reviewDocumentMotion                   = "motion"
	reviewDocumentOpposition               = "opposition"
	reviewDocumentResponse                 = "response"
	reviewDocumentReply                    = "reply"
	reviewDocumentDeclaration              = "declaration"
	reviewDocumentProposedOrder            = "proposed-order"
	reviewDocumentCoParentingCommunication = "co-parenting-communication"

	workflowTypeTriage         = "triage"
	workflowTypeCounselBrief   = "counsel-brief"
	workflowTypeSMTriage       = "sm-triage"
	workflowTypeDraftMessage   = "draft-message"
	workflowTypeEvidencePacket = "evidence-packet"
	workflowTypePatternReview  = "pattern-review"
	workflowTypeFactLock       = "fact-lock"
	workflowTypeAuthorityCheck = "authority-check"
	workflowTypePrep           = "prep"
	workflowTypeJournalEntry   = "journal-entry"

	personaChildBestInterests = "child-best-interests-family-dynamics-reviewer"
	personaFinancialSupport   = "financial-support-reviewer"
	personaLegalAuthority     = "legal-authority-scholar"
	personaStrategicOptions   = "strategic-options-architect"
)

type reviewModeSpec struct {
	Name        string   `json:"name"`
	PipelineKey string   `json:"pipeline_key"`
	Aliases     []string `json:"aliases,omitempty"`
	Description string   `json:"description"`
}

type reviewDocumentTypeSpec struct {
	Name        string   `json:"name"`
	PipelineKey string   `json:"pipeline_key"`
	MapsTo      string   `json:"maps_to"`
	Aliases     []string `json:"aliases,omitempty"`
	Description string   `json:"description"`
}

type workflowRouteCatalogSpec struct {
	Type                 string   `json:"type"`
	DisplayName          string   `json:"display_name"`
	PipelineKey          string   `json:"pipeline_key"`
	OutputFile           string   `json:"output_file"`
	Description          string   `json:"description"`
	Aliases              []string `json:"aliases,omitempty"`
	RequireSituation     bool     `json:"require_situation"`
	RequireAudience      bool     `json:"require_audience"`
	AllowDraft           bool     `json:"allow_draft"`
	AutoDraftFallback    bool     `json:"auto_draft_fallback"`
	RequiresIntake       bool     `json:"requires_intake"`
	DefaultAudience      string   `json:"default_audience,omitempty"`
	DefaultGoal          string   `json:"default_goal,omitempty"`
	ChildConditional     bool     `json:"child_conditional"`
	FinancialConditional bool     `json:"financial_conditional"`
}

type conditionalRuleSpec struct {
	Name        string   `json:"name"`
	Persona     string   `json:"persona"`
	AppliesTo   []string `json:"applies_to"`
	Trigger     string   `json:"trigger"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords,omitempty"`
}

type legacySurfaceSpec struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Status      string `json:"status"`
	Replacement string `json:"replacement,omitempty"`
	Notes       string `json:"notes,omitempty"`
}

type routeSnapshot struct {
	ActiveCommands      []string                   `json:"active_commands"`
	ReviewModes         []reviewModeSpec           `json:"review_modes"`
	ReviewDocumentTypes []reviewDocumentTypeSpec   `json:"review_document_types"`
	WorkflowTypes       []workflowRouteCatalogSpec `json:"workflow_types"`
	ConditionalRules    []conditionalRuleSpec      `json:"conditional_rules"`
	LegacySurfaces      []legacySurfaceSpec        `json:"legacy_surfaces"`
}

func activeCommandNames() []string {
	return []string{"guide", "intake", "review", "workflow", "run", "setup", "historian", "tui"}
}

func legalReviewModes() []reviewModeSpec {
	return []reviewModeSpec{
		{Name: reviewModeQuick, PipelineKey: "review_quick", Description: "Fast attorney-readiness screen."},
		{Name: reviewModeStandard, PipelineKey: "review_standard", Description: "Default multi-reviewer draft review."},
		{Name: reviewModeDeep, PipelineKey: "review_deep", Description: "Expanded high-stakes review with strategy depth."},
		{Name: reviewModeLegalResearch, PipelineKey: "review_legal_research", Aliases: []string{"legal_research", "research"}, Description: "Authority-focused review path."},
		{Name: reviewModeStrategy, PipelineKey: "review_strategy", Description: "Strategy-focused review path."},
	}
}

func legalReviewDocumentTypes() []reviewDocumentTypeSpec {
	return []reviewDocumentTypeSpec{
		{Name: reviewDocumentMotion, PipelineKey: "review_standard", MapsTo: reviewModeStandard, Description: "Motion folds into standard review."},
		{Name: reviewDocumentOpposition, PipelineKey: "review_standard", MapsTo: reviewModeStandard, Description: "Opposition folds into standard review."},
		{Name: reviewDocumentResponse, PipelineKey: "review_standard", MapsTo: reviewModeStandard, Description: "Response folds into standard review."},
		{Name: reviewDocumentReply, PipelineKey: "review_document_reply", MapsTo: "document-specific", Description: "Reply-specific review route."},
		{Name: reviewDocumentDeclaration, PipelineKey: "review_document_declaration", MapsTo: "document-specific", Description: "Declaration-specific review route."},
		{Name: reviewDocumentProposedOrder, PipelineKey: "review_document_proposed_order", MapsTo: "document-specific", Aliases: []string{"proposed_order", "order"}, Description: "Proposed-order-specific review route."},
		{Name: reviewDocumentCoParentingCommunication, PipelineKey: "review_document_coparenting_communication", MapsTo: "document-specific", Aliases: []string{"coparenting-communication", "co_parenting_communication", "communication"}, Description: "Co-parenting communication route."},
	}
}

func legalWorkflowRouteCatalog() []workflowRouteCatalogSpec {
	return []workflowRouteCatalogSpec{
		{Type: workflowTypeTriage, DisplayName: "New Situation Triage", PipelineKey: "triage_new_situation", OutputFile: "triage.md", Description: "Understand a new situation and decide pursue / preserve / let go.", RequireSituation: true, AllowDraft: true, DefaultAudience: "internal", DefaultGoal: "Analyze the situation, identify source strength and risks, and recommend whether to pursue, preserve, or let go.", ChildConditional: true, FinancialConditional: true},
		{Type: workflowTypeCounselBrief, DisplayName: "Counsel Brief", PipelineKey: "brief_counsel", OutputFile: "counsel-brief.md", Description: "Prepare an attorney-ready issue brief.", RequireSituation: true, AllowDraft: true, RequiresIntake: false, DefaultAudience: "counsel", DefaultGoal: "Prepare a concise attorney-ready brief with key facts, why it matters, evidence available, evidence missing, options, and specific questions for counsel.", ChildConditional: true, FinancialConditional: true},
		{Type: workflowTypeSMTriage, DisplayName: "Special Master Triage", PipelineKey: "triage_special_master", OutputFile: "sm-triage.md", Description: "Decide whether an issue is ready for the Special Master.", RequireSituation: true, AllowDraft: true, DefaultAudience: "special-master", DefaultGoal: "Determine whether this issue is ready and appropriate for the Special Master, what directive should be requested, what evidence is needed, and whether the issue is too small or premature.", ChildConditional: true, FinancialConditional: true},
		{Type: workflowTypeDraftMessage, DisplayName: "External Message Drafting", PipelineKey: "draft_external_message", OutputFile: "draft-message.md", Description: "Draft an audience-appropriate external message.", RequireSituation: true, RequireAudience: true, AllowDraft: true, DefaultGoal: "Draft or evaluate a court-safe, audience-appropriate message that preserves the record without over-escalating.", ChildConditional: true, FinancialConditional: true},
		{Type: workflowTypeEvidencePacket, DisplayName: "Evidence Packet Builder", PipelineKey: "build_evidence_packet", OutputFile: "evidence-packet.md", Description: "Build the most useful source packet for an issue.", RequireSituation: true, DefaultAudience: "internal", DefaultGoal: "Identify the source documents, evidence anchors, chronology, fact gaps, and packet structure needed for this issue.", ChildConditional: true, FinancialConditional: true},
		{Type: workflowTypePatternReview, DisplayName: "Pattern Review", PipelineKey: "review_pattern", OutputFile: "pattern-review.md", Description: "Analyze whether events form a meaningful pattern.", RequireSituation: true, DefaultAudience: "internal", DefaultGoal: "Analyze whether recurring events add up to a meaningful pattern, what examples are strongest, and whether the pattern supports counsel, SM, court, or journal use.", ChildConditional: true, FinancialConditional: true},
		{Type: workflowTypeFactLock, DisplayName: "Fact Lock / Source Strength Review", PipelineKey: "fact_lock", OutputFile: "fact-lock.md", Description: "Separate locked facts from recollection, inference, and unsupported claims.", RequireSituation: true, AllowDraft: true, DefaultGoal: "Separate locked/source-supported facts from recollection, inference, disputed claims, unsupported claims, and legal propositions needing authority.", ChildConditional: true, FinancialConditional: true},
		{Type: workflowTypeAuthorityCheck, DisplayName: "Authority Check", PipelineKey: "authority_check", OutputFile: "authority-check.md", Description: "Check legal authority gaps without overclaiming currentness.", AllowDraft: true, AutoDraftFallback: true, DefaultAudience: "internal", DefaultGoal: "Identify legal authority questions, authority gaps, and issues needing attorney or live legal research verification."},
		{Type: workflowTypePrep, DisplayName: "Hearing / Call Prep", PipelineKey: "prep_hearing_or_call", OutputFile: "prep.md", Description: "Prepare for a call, hearing, or conference.", RequireSituation: true, AllowDraft: true, DefaultAudience: "internal", DefaultGoal: "Prepare a focused call/hearing/conference plan with key facts, likely attacks, questions, decisions needed, and recommended next steps.", ChildConditional: true, FinancialConditional: true},
		{Type: workflowTypeJournalEntry, DisplayName: "Post-Event Journal Entry", PipelineKey: "journal_entry", OutputFile: "journal-entry.md", Description: "Create a neutral, source-aware journal entry.", RequireSituation: true, DefaultAudience: "internal", DefaultGoal: "Create a neutral, source-aware journal entry with summary, what happened, legal/co-parenting significance, evaluator notes, evidence anchors, caveats, and related events.", ChildConditional: true, FinancialConditional: true},
	}
}

func workflowRouteCatalogByType(workflowType string) (workflowRouteCatalogSpec, bool) {
	for _, route := range legalWorkflowRouteCatalog() {
		if route.Type == normalizeRouteName(workflowType) {
			return route, true
		}
	}
	return workflowRouteCatalogSpec{}, false
}

func legalConditionalRules() []conditionalRuleSpec {
	return []conditionalRuleSpec{
		{Name: "child-related specialist", Persona: personaChildBestInterests, AppliesTo: []string{routeKindReview, routeKindWorkflow}, Trigger: "explicit --child-related flag, route/document-type hints, or child-related keywords", Description: "Adds child best-interests and family-dynamics review.", Keywords: []string{"child", "children", "custody", "parent-time", "parenting time", "school", "therapy", "medical", "doctor", "gal", "special master"}},
		{Name: "financial specialist", Persona: personaFinancialSupport, AppliesTo: []string{routeKindReview, routeKindWorkflow}, Trigger: "explicit --financial flag, route/document-type hints, or financial keywords", Description: "Adds financial/support review.", Keywords: []string{"child support", "support", "alimony", "income", "expense", "expenses", "reimbursement", "attorney fees", "medical expenses", "school expenses", "imputed income", "billing", "financial"}},
		{Name: "authority specialist", Persona: personaLegalAuthority, AppliesTo: []string{routeKindReview}, Trigger: "legal-research mode, legal-sensitive document types, or legal authority keywords", Description: "Adds legal authority scholar review.", Keywords: []string{"under utah law", "the law requires", "rule ", "utah rule", "utah code", "statute", "case law", "joint legal custody", "best interests", "contempt", "preliminary injunction"}},
		{Name: "strategy specialist", Persona: personaStrategicOptions, AppliesTo: []string{routeKindReview}, Trigger: "strategy mode, strategy-sensitive document types, or strategy keywords", Description: "Adds strategic options review.", Keywords: []string{"status quo", "temporary restraining order", "tro", "injunction", "emergency", "expedited", "special master", "gal", "custody evaluation", "psychological evaluation", "motion to enforce"}},
	}
}

func legacySurfacesCatalog() []legacySurfaceSpec {
	return []legacySurfaceSpec{
		{Name: "legal-bot-discovery", Path: "bin/legal-bot-discovery", Status: "retired", Replacement: "legal-bot guide | intake | review | workflow", Notes: "Legacy Vern discovery wrapper. Current Go CLI does not expose a discovery subcommand."},
		{Name: "legal-bot-discovery.cmd", Path: "bin/legal-bot-discovery.cmd", Status: "retired", Replacement: "legal-bot guide | intake | review | workflow", Notes: "Windows legacy Vern discovery wrapper."},
		{Name: "legal-bot-council", Path: "bin/legal-bot-council", Status: "retired", Replacement: "legal-bot tui", Notes: "Legacy VernHole council wrapper. Current Go CLI does not expose hole."},
		{Name: "legal-bot-council.cmd", Path: "bin/legal-bot-council.cmd", Status: "retired", Replacement: "legal-bot tui", Notes: "Windows legacy VernHole council wrapper."},
		{Name: "legal-bot-oracle", Path: "bin/legal-bot-oracle", Status: "retired", Replacement: "docs/legacy_surfaces.md", Notes: "Legacy oracle wrapper. Current legal CLI does not expose oracle."},
		{Name: "legal-bot-generate", Path: "bin/legal-bot-generate", Status: "retired", Replacement: "docs/legacy_surfaces.md", Notes: "Legacy persona-generation wrapper. Current legal CLI does not expose generate."},
		{Name: "legal-bot-run", Path: "bin/legal-bot-run", Status: "partial", Replacement: "legal-bot run", Notes: "Compatibility wrapper for the active low-level run command."},
		{Name: "legal-bot-run.cmd", Path: "bin/legal-bot-run.cmd", Status: "partial", Replacement: "legal-bot run", Notes: "Windows compatibility wrapper for the active run command."},
		{Name: "legal-bot-historian", Path: "bin/legal-bot-historian", Status: "partial", Replacement: "legal-bot historian", Notes: "Compatibility wrapper for the active historian utility."},
		{Name: "install-legal-bot-cli", Path: "bin/install-legal-bot-cli", Status: "partial", Replacement: "cd go && go build -o bin/legal-bot ./cmd/legal-bot", Notes: "Binary install helper, not part of the current legal workflow."},
		{Name: "install-legal-bot-cli.cmd", Path: "bin/install-legal-bot-cli.cmd", Status: "partial", Replacement: "cd go && go build -o bin/legal-bot.exe ./cmd/legal-bot", Notes: "Windows binary install helper."},
		{Name: "go/internal/pipeline", Path: "go/internal/pipeline", Status: "partial", Replacement: "Active for historian and TUI compatibility", Notes: "Legacy/mixed orchestration package still used by historian and TUI."},
		{Name: "go/internal/council", Path: "go/internal/council", Status: "legacy", Replacement: "docs/legacy_surfaces.md", Notes: "VernHole council support, not part of the active legal workflow contract."},
		{Name: "go/internal/generate", Path: "go/internal/generate", Status: "legacy", Replacement: "docs/legacy_surfaces.md", Notes: "Legacy persona-generation helpers."},
		{Name: "go/internal/vts", Path: "go/internal/vts", Status: "legacy", Replacement: "docs/legacy_surfaces.md", Notes: "Legacy Vern Task Spec helpers."},
		{Name: "go/internal/tobeads", Path: "go/internal/tobeads", Status: "legacy", Replacement: "docs/legacy_surfaces.md", Notes: "Legacy VTS-to-Beads import tooling."},
		{Name: "go/internal/tui", Path: "go/internal/tui", Status: "partial", Replacement: "legal-bot tui", Notes: "Still exposed, but oriented around legacy discovery/VernHole utilities rather than the core legal workflow."},
	}
}

func reviewModePipelineKey(mode string) string {
	normalized := normalizeRouteName(mode)
	for _, spec := range legalReviewModes() {
		if normalized == spec.Name || slices.Contains(spec.Aliases, normalized) {
			return spec.PipelineKey
		}
	}
	return "review_standard"
}

func reviewDocumentTypePipelineKey(documentType string) (string, bool) {
	normalized := normalizeRouteName(documentType)
	for _, spec := range legalReviewDocumentTypes() {
		if normalized == spec.Name || slices.Contains(spec.Aliases, normalized) {
			return spec.PipelineKey, true
		}
	}
	return "", false
}

func reviewDocumentTriggersChildIssues(route string) bool {
	return route == "custody" || route == "school" || route == "therapy" || route == "medical"
}

func reviewDocumentTriggersFinancialIssues(route string) bool {
	return route == "financial" || route == "support"
}

func reviewDocumentTriggersStrategy(route string) bool {
	switch route {
	case reviewDocumentMotion, reviewDocumentDeclaration, "custody", "school", "therapy", "medical":
		return true
	default:
		return false
	}
}

func reviewDocumentTriggersLegalAuthority(route string) bool {
	switch route {
	case reviewDocumentMotion, reviewDocumentOpposition, reviewDocumentResponse, reviewDocumentReply, reviewDocumentDeclaration, reviewDocumentProposedOrder, "proposed_order", "order", "communication", "custody", "school", "therapy", "medical", "financial", "support":
		return true
	default:
		return false
	}
}

func buildRouteSnapshot() routeSnapshot {
	return routeSnapshot{
		ActiveCommands:      activeCommandNames(),
		ReviewModes:         legalReviewModes(),
		ReviewDocumentTypes: legalReviewDocumentTypes(),
		WorkflowTypes:       legalWorkflowRouteCatalog(),
		ConditionalRules:    legalConditionalRules(),
		LegacySurfaces:      legacySurfacesCatalog(),
	}
}

func marshalRouteSnapshotJSON() ([]byte, error) {
	return json.MarshalIndent(buildRouteSnapshot(), "", "  ")
}

func collectAvailableAgentNames(projectRoot string) map[string]bool {
	available := map[string]bool{}
	for _, name := range embedded.ListAgents() {
		available[name] = true
	}
	agentsDir := filepath.Join(projectRoot, "agents")
	entries, err := os.ReadDir(agentsDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
				continue
			}
			available[strings.TrimSuffix(entry.Name(), ".md")] = true
		}
	}
	return available
}

func validateRouteCatalog(projectRoot string) []error {
	cfg, err := config.LoadProjectDefault(projectRoot)
	if err != nil {
		return []error{err}
	}

	availableAgents := collectAvailableAgentNames(projectRoot)
	var errs []error

	checkAgent := func(agent string, context string) {
		resolved := personapkg.ResolveName(agent)
		if !availableAgents[resolved] {
			errs = append(errs, fmt.Errorf("%s references missing agent %q (resolved %q)", context, agent, resolved))
		}
	}

	for pipelineKey, steps := range cfg.LegalPipelines {
		for _, step := range steps {
			checkAgent(step.Persona, "config pipeline "+pipelineKey)
		}
	}

	for _, mode := range legalReviewModes() {
		if _, ok := cfg.LegalPipelines[mode.PipelineKey]; !ok {
			errs = append(errs, fmt.Errorf("review mode %q maps to missing pipeline key %q", mode.Name, mode.PipelineKey))
		}
	}
	for _, document := range legalReviewDocumentTypes() {
		if _, ok := cfg.LegalPipelines[document.PipelineKey]; !ok {
			errs = append(errs, fmt.Errorf("review document type %q maps to missing pipeline key %q", document.Name, document.PipelineKey))
		}
	}
	for _, workflow := range legalWorkflowRouteCatalog() {
		if _, ok := cfg.LegalPipelines[workflow.PipelineKey]; !ok {
			errs = append(errs, fmt.Errorf("workflow type %q maps to missing pipeline key %q", workflow.Type, workflow.PipelineKey))
		}
	}
	for _, conditional := range legalConditionalRules() {
		checkAgent(conditional.Persona, "conditional rule "+conditional.Name)
	}
	for legacy, active := range personapkg.LegacyAliases() {
		if personapkg.ResolveName(legacy) != active {
			errs = append(errs, fmt.Errorf("legacy persona alias %q does not resolve to %q", legacy, active))
		}
		checkAgent(active, "legacy persona alias "+legacy)
	}

	reviewAliases := map[string]string{"review-draft": "review"}
	for alias, target := range reviewAliases {
		if target != "review" {
			errs = append(errs, fmt.Errorf("unexpected review alias target for %q: %q", alias, target))
		}
	}

	intentSeen := map[string]string{}
	for alias, choice := range guidedIntentAliases() {
		if prior, ok := intentSeen[alias]; ok {
			errs = append(errs, fmt.Errorf("duplicate guided alias %q reused for %s and %d", alias, prior, choice))
			continue
		}
		intentSeen[alias] = fmt.Sprintf("%d", choice)
	}

	for _, path := range []string{"scripts/run-intake.bat", "scripts/run-review.bat", "scripts/open-latest-output.bat", "scripts/run-guide.bat"} {
		full := filepath.Join(projectRoot, filepath.FromSlash(path))
		if _, statErr := os.Stat(full); statErr != nil {
			errs = append(errs, fmt.Errorf("expected active script missing: %s", path))
		}
	}

	for _, surface := range legacySurfacesCatalog() {
		full := filepath.Join(projectRoot, filepath.FromSlash(surface.Path))
		if _, statErr := os.Stat(full); statErr != nil {
			errs = append(errs, fmt.Errorf("legacy surface inventory references missing path %s", surface.Path))
			continue
		}
		if strings.HasPrefix(surface.Path, "bin/") {
			data, readErr := os.ReadFile(full)
			if readErr != nil {
				errs = append(errs, fmt.Errorf("read %s: %w", surface.Path, readErr))
				continue
			}
			content := strings.ToLower(string(data))
			if surface.Status == "retired" && !strings.Contains(content, "legacy") {
				errs = append(errs, fmt.Errorf("retired wrapper %s should be clearly marked legacy", surface.Path))
			}
		}
	}

	return errs
}
