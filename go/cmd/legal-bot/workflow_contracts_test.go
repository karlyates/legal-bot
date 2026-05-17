package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEveryWorkflowTypeHasOutputContract(t *testing.T) {
	contracts := map[string]workflowOutputContract{}
	for _, contract := range allWorkflowOutputContracts() {
		contracts[contract.Type] = contract
	}

	for _, route := range legalWorkflowRouteCatalog() {
		contract, ok := contracts[route.Type]
		if !ok {
			t.Fatalf("workflow type %q missing output contract", route.Type)
		}
		if contract.Purpose == "" {
			t.Fatalf("workflow type %q has empty purpose", route.Type)
		}
		if len(contract.RequiredSections) == 0 {
			t.Fatalf("workflow type %q has no required sections", route.Type)
		}
	}
}

func TestNoWorkflowContractExistsForUnknownRoute(t *testing.T) {
	known := map[string]bool{}
	for _, route := range legalWorkflowRouteCatalog() {
		known[route.Type] = true
	}
	for _, contract := range allWorkflowOutputContracts() {
		if !known[contract.Type] {
			t.Fatalf("contract exists for unknown workflow type %q", contract.Type)
		}
	}
}

func TestRenderedWorkflowContractIncludesGlobalRules(t *testing.T) {
	contract, ok := workflowOutputContractForType(workflowTypeTriage)
	if !ok {
		t.Fatal("triage contract missing")
	}
	rendered := renderWorkflowContract(contract)
	for _, fragment := range []string{
		"Recommended Track:",
		"Source-strength taxonomy:",
		"Locked / source-supported fact",
		"Attorney-only / internal strategy",
		"Do not convert internal strategy into accusations",
		"Do Not Chase rules:",
	} {
		if !strings.Contains(rendered, fragment) {
			t.Fatalf("rendered contract missing %q", fragment)
		}
	}
}

func TestContractsHaveNonGenericDecisionAndRuleContent(t *testing.T) {
	for _, contract := range allWorkflowOutputContracts() {
		if strings.TrimSpace(contract.DecisionType) == "" {
			t.Fatalf("%s contract missing decision type", contract.Type)
		}
		if len(contract.RecommendedTracks) == 0 {
			t.Fatalf("%s contract missing recommended tracks", contract.Type)
		}
		if len(contract.Rules) == 0 {
			t.Fatalf("%s contract missing rules", contract.Type)
		}
	}
}

func TestContractsCoverAttorneyQuestionsExternalLanguageAndDoNotChase(t *testing.T) {
	for _, contract := range allWorkflowOutputContracts() {
		if len(contract.AttorneyQuestionRules) == 0 {
			t.Fatalf("%s contract missing attorney-question rules", contract.Type)
		}
		if len(contract.ExternalLanguageRules) == 0 {
			t.Fatalf("%s contract missing external-language rules", contract.Type)
		}
		if len(contract.DoNotChaseRules) == 0 {
			t.Fatalf("%s contract missing Do Not Chase rules", contract.Type)
		}
	}
}

func TestWorkflowSpecificRequiredHeadings(t *testing.T) {
	tests := map[string][]string{
		workflowTypeTriage:         {"Bottom Line", "Source Strength", "Options Matrix", "Do Not Chase Check"},
		workflowTypeCounselBrief:   {"Short Version", "Legal / Strategy Questions for Counsel", "Attorney-Ready Draft"},
		workflowTypeSMTriage:       {"Scope / Authority Fit", "Ripeness", "Requested Directive"},
		workflowTypeDraftMessage:   {"Facts Safe to Say", "Facts or Claims to Avoid", "Recommended Message"},
		workflowTypeEvidencePacket: {"Evidence Inventory", "Missing Evidence", "Fact Map"},
		workflowTypePatternReview:  {"Pattern Theory", "Strongest Examples", "Weak or Risky Examples"},
		workflowTypeFactLock:       {"Locked Facts", "Partially Supported Facts", "Legal Propositions Needing Authority"},
		workflowTypeAuthorityCheck: {"Authority Gaps", "Currentness / Citator Status", "Do Not Overclaim"},
		workflowTypePrep:           {"Decisions Needed", "Likely Pushback / Attacks", "What Not to Say"},
		workflowTypeJournalEntry:   {"What Happened", "Neutral Evaluator Notes", "Evidence / Source Anchors", "Caveats / Verification Needed"},
	}

	for workflowType, expected := range tests {
		contract, ok := workflowOutputContractForType(workflowType)
		if !ok {
			t.Fatalf("missing contract for %s", workflowType)
		}
		rendered := renderWorkflowContract(contract)
		for _, heading := range expected {
			if !strings.Contains(rendered, "## "+heading) {
				t.Fatalf("%s contract missing heading %q", workflowType, heading)
			}
		}
	}
}

func TestManagingPartnerPersonaMentionsWorkflowContractsAndDoNotChase(t *testing.T) {
	root := repoRootForTests(t)
	path := filepath.Join(root, "agents", "managing-partner-final-synthesizer.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	for _, fragment := range []string{
		"decision-maker",
		"workflow output contract",
		"Do Not Chase",
		"not a summarizer",
		"source-supported facts",
		"external-safe",
	} {
		if !strings.Contains(strings.ToLower(content), strings.ToLower(fragment)) {
			t.Fatalf("managing partner persona missing %q", fragment)
		}
	}
}

func TestWorkflowContractDocsExistAndMentionAllWorkflowTypes(t *testing.T) {
	root := repoRootForTests(t)
	path := filepath.Join(root, "docs", "workflow_output_contracts.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	for _, fragment := range []string{
		"source-strength taxonomy",
		"internal strategy",
		"external-safe language",
		"Do Not Chase",
		workflowTypeTriage,
		workflowTypeCounselBrief,
		workflowTypeSMTriage,
		workflowTypeDraftMessage,
		workflowTypeEvidencePacket,
		workflowTypePatternReview,
		workflowTypeFactLock,
		workflowTypeAuthorityCheck,
		workflowTypePrep,
		workflowTypeJournalEntry,
	} {
		if !strings.Contains(content, fragment) {
			t.Fatalf("workflow contract doc missing %q", fragment)
		}
	}
}
