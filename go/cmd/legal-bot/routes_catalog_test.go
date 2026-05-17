package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jdonohoo/legal-bot/go/internal/config"
)

func repoRootForTests(t *testing.T) string {
	t.Helper()
	root := filepath.Clean(filepath.Join("..", "..", ".."))
	if _, err := os.Stat(filepath.Join(root, "config.default.json")); err != nil {
		t.Fatalf("repo root not found from test cwd: %v", err)
	}
	return root
}

func TestLegalReviewModesPresent(t *testing.T) {
	var got []string
	for _, mode := range legalReviewModes() {
		got = append(got, mode.Name)
	}
	want := []string{reviewModeQuick, reviewModeStandard, reviewModeDeep, reviewModeLegalResearch, reviewModeStrategy}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("review modes = %v, want %v", got, want)
	}
}

func TestLegalReviewDocumentTypesPresent(t *testing.T) {
	var got []string
	for _, document := range legalReviewDocumentTypes() {
		got = append(got, document.Name)
	}
	want := []string{
		reviewDocumentMotion,
		reviewDocumentOpposition,
		reviewDocumentResponse,
		reviewDocumentReply,
		reviewDocumentDeclaration,
		reviewDocumentProposedOrder,
		reviewDocumentCoParentingCommunication,
	}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("review document types = %v, want %v", got, want)
	}
}

func TestLegalWorkflowTypesPresent(t *testing.T) {
	specs := workflowSpecs()
	for _, workflowType := range []string{
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
		if _, ok := specs[workflowType]; !ok {
			t.Fatalf("workflow type %q missing", workflowType)
		}
	}
}

func TestReviewPipelineMappings(t *testing.T) {
	tests := map[string]string{
		reviewDocumentMotion:                   "review_standard",
		reviewDocumentOpposition:               "review_standard",
		reviewDocumentResponse:                 "review_standard",
		reviewDocumentReply:                    "review_document_reply",
		reviewDocumentDeclaration:              "review_document_declaration",
		reviewDocumentProposedOrder:            "review_document_proposed_order",
		reviewDocumentCoParentingCommunication: "review_document_coparenting_communication",
	}
	for input, want := range tests {
		if got := reviewPipelineKey("", input); got != want {
			t.Fatalf("reviewPipelineKey(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestGuidedChoicesReferenceKnownWorkflowTypes(t *testing.T) {
	valid := workflowSpecs()
	for _, choice := range guidedChoices() {
		if choice.Kind != routeKindWorkflow {
			continue
		}
		if _, ok := valid[choice.WorkflowType]; !ok {
			t.Fatalf("guided choice %d references unknown workflow type %q", choice.ID, choice.WorkflowType)
		}
	}
}

func TestConditionalReviewRoutingAddsSpecialistsAndAvoidsDuplicates(t *testing.T) {
	projectRoot := repoRootForTests(t)
	cfg, err := config.LoadProjectDefault(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	steps := cfg.GetLegalPipeline("review_standard")
	if len(steps) == 0 {
		t.Fatal("review_standard pipeline missing")
	}

	routed := applyConditionalReviewRouting(steps, "This child support draft cites Utah law and asks for emergency relief.", reviewModeDeep, reviewDocumentMotion, true, true)
	personaCounts := map[string]int{}
	for _, step := range routed {
		personaCounts[step.Persona]++
	}
	for _, persona := range []string{personaChildBestInterests, personaFinancialSupport, personaLegalAuthority, personaStrategicOptions} {
		if personaCounts[persona] != 1 {
			t.Fatalf("expected %s exactly once, got %d", persona, personaCounts[persona])
		}
	}

	routedAgain := applyConditionalReviewRouting(routed, "This child support draft cites Utah law and asks for emergency relief.", reviewModeDeep, reviewDocumentMotion, true, true)
	personaCounts = map[string]int{}
	for _, step := range routedAgain {
		personaCounts[step.Persona]++
	}
	for _, persona := range []string{personaChildBestInterests, personaFinancialSupport, personaLegalAuthority, personaStrategicOptions} {
		if personaCounts[persona] != 1 {
			t.Fatalf("expected %s exactly once after reroute, got %d", persona, personaCounts[persona])
		}
	}
}

func TestConditionalWorkflowRoutingAddsSpecialistsAndAvoidsDuplicates(t *testing.T) {
	projectRoot := repoRootForTests(t)
	cfg, err := config.LoadProjectDefault(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	steps := cfg.GetLegalPipeline("triage_new_situation")
	if len(steps) == 0 {
		t.Fatal("triage_new_situation pipeline missing")
	}
	spec := workflowSpecs()[workflowTypeTriage]
	inputs := workflowInputs{Situation: "Child support issue at school and medical billing.", Goal: "Decide next steps."}
	routed := applyConditionalWorkflowRouting(steps, spec, inputs, "internal", "high", true, true)

	personaCounts := map[string]int{}
	for _, step := range routed {
		personaCounts[step.Persona]++
	}
	for _, persona := range []string{personaChildBestInterests, personaFinancialSupport} {
		if personaCounts[persona] != 1 {
			t.Fatalf("expected %s exactly once, got %d", persona, personaCounts[persona])
		}
	}

	routedAgain := applyConditionalWorkflowRouting(routed, spec, inputs, "internal", "high", true, true)
	personaCounts = map[string]int{}
	for _, step := range routedAgain {
		personaCounts[step.Persona]++
	}
	for _, persona := range []string{personaChildBestInterests, personaFinancialSupport} {
		if personaCounts[persona] != 1 {
			t.Fatalf("expected %s exactly once after reroute, got %d", persona, personaCounts[persona])
		}
	}
}

func TestRouteSnapshotJSON(t *testing.T) {
	data, err := marshalRouteSnapshotJSON()
	if err != nil {
		t.Fatal(err)
	}
	var snapshot routeSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		t.Fatal(err)
	}
	if len(snapshot.ReviewModes) == 0 || len(snapshot.WorkflowTypes) == 0 || len(snapshot.LegacySurfaces) == 0 {
		t.Fatalf("unexpected empty route snapshot: %+v", snapshot)
	}
}

func TestRoutesCommandOutputIncludesKeySections(t *testing.T) {
	var buf bytes.Buffer
	routesCmd.SetOut(&buf)
	routesCmd.SetErr(&buf)
	routesJSON = false
	if err := runRoutes(routesCmd, nil); err != nil {
		t.Fatal(err)
	}
	output := buf.String()
	for _, fragment := range []string{"Active commands:", "Review modes:", "Workflow types:", "Conditional specialists:", "Legacy surfaces:"} {
		if !strings.Contains(output, fragment) {
			t.Fatalf("routes output missing %q", fragment)
		}
	}
}

func TestValidateRouteCatalog(t *testing.T) {
	errs := validateRouteCatalog(repoRootForTests(t))
	if len(errs) > 0 {
		var lines []string
		for _, err := range errs {
			lines = append(lines, err.Error())
		}
		t.Fatalf("validateRouteCatalog returned errors:\n%s", strings.Join(lines, "\n"))
	}
}

func TestLegacyWrapperInventoryIsDocumented(t *testing.T) {
	root := repoRootForTests(t)
	docPath := filepath.Join(root, "docs", "legacy_surfaces.md")
	data, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("read legacy surfaces doc: %v", err)
	}
	content := string(data)
	for _, surface := range []string{"bin/legal-bot-discovery", "bin/legal-bot-council", "bin/legal-bot-oracle", "bin/legal-bot-generate"} {
		if !strings.Contains(content, surface) {
			t.Fatalf("legacy surfaces doc missing %s", surface)
		}
	}
}

func TestRootCommandsPresent(t *testing.T) {
	expected := []string{"guide", "intake", "review", "workflow", "routes", "run", "setup", "historian", "tui"}
	names := map[string]bool{}
	for _, cmd := range rootCmd.Commands() {
		names[cmd.Name()] = true
	}
	for _, name := range expected {
		if !names[name] {
			t.Fatalf("root command %q missing", name)
		}
	}
}

func TestRetiredLegacyWrappersExplainThemselves(t *testing.T) {
	root := repoRootForTests(t)
	for _, rel := range []string{
		"bin/legal-bot-discovery",
		"bin/legal-bot-discovery.cmd",
		"bin/legal-bot-council",
		"bin/legal-bot-council.cmd",
		"bin/legal-bot-oracle",
		"bin/legal-bot-generate",
	} {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatalf("read %s: %v", rel, err)
		}
		content := strings.ToLower(string(data))
		if !strings.Contains(content, "legacy") || !strings.Contains(content, "current legal-bot") {
			t.Fatalf("%s is not clearly marked legacy", rel)
		}
	}
}
