package main

import (
	"strings"
	"testing"
)

func TestWorkflowSpecsIncludeCoreTypes(t *testing.T) {
	specs := workflowSpecs()
	for _, workflowType := range []string{
		"triage",
		"counsel-brief",
		"sm-triage",
		"draft-message",
		"evidence-packet",
		"pattern-review",
		"fact-lock",
		"authority-check",
		"prep",
		"journal-entry",
	} {
		spec, ok := specs[workflowType]
		if !ok {
			t.Fatalf("workflow type %q missing", workflowType)
		}
		if spec.PipelineKey == "" {
			t.Fatalf("workflow type %q missing pipeline key", workflowType)
		}
		if spec.OutputFile == "" {
			t.Fatalf("workflow type %q missing output file", workflowType)
		}
	}
}

func TestBuildWorkflowPromptIncludesCoreContext(t *testing.T) {
	inputs := workflowInputs{
		Situation: "A new provider called about an unexpected referral.",
		Goal:      "Decide whether to raise this with counsel.",
		Questions: "Should this be preserved or escalated?",
	}

	prompt := buildWorkflowPrompt(
		"situation_plus_knowledge",
		"You are testing.",
		"=== WORKFLOW METADATA ===\nworkflow_type: triage\n=== END WORKFLOW METADATA ===",
		"=== WORKFLOW PURPOSE AND OUTPUT CONTRACT SUMMARY ===\npurpose: test\n=== END WORKFLOW PURPOSE AND OUTPUT CONTRACT SUMMARY ===",
		"",
		"case context",
		inputs,
		"knowledge",
		"input manifest",
		"prior outputs",
	)

	for _, fragment := range []string{
		"workflow_type: triage",
		"=== USER SITUATION ===",
		"unexpected referral",
		"=== USER GOAL ===",
		"=== USER QUESTIONS ===",
		"=== KNOWLEDGE LAYER ===",
		"input manifest",
	} {
		if !strings.Contains(prompt, fragment) {
			t.Fatalf("prompt missing fragment %q", fragment)
		}
	}
}

func TestBuildWorkflowPromptIncludesContractSummary(t *testing.T) {
	inputs := workflowInputs{
		Situation: "School access issue.",
		Goal:      "Decide whether to preserve or escalate.",
	}

	prompt := buildWorkflowPrompt(
		"situation_plus_knowledge",
		"You are testing.",
		"=== WORKFLOW METADATA ===\nworkflow_type: triage\n=== END WORKFLOW METADATA ===",
		"=== WORKFLOW PURPOSE AND OUTPUT CONTRACT SUMMARY ===\npurpose: Analyze a new or messy situation.\nrequired_sections:\n- Bottom Line\n=== END WORKFLOW PURPOSE AND OUTPUT CONTRACT SUMMARY ===",
		"",
		"case context",
		inputs,
		"knowledge",
		"input manifest",
		"",
	)

	for _, fragment := range []string{
		"=== WORKFLOW PURPOSE AND OUTPUT CONTRACT SUMMARY ===",
		"purpose: Analyze a new or messy situation.",
		"- Bottom Line",
	} {
		if !strings.Contains(prompt, fragment) {
			t.Fatalf("prompt missing fragment %q", fragment)
		}
	}
}

func TestBuildWorkflowPromptIncludesFullContractForFinalSynthesis(t *testing.T) {
	inputs := workflowInputs{
		Situation: "Provider access problem.",
		Goal:      "Decide whether to involve the Special Master.",
	}
	contract, ok := workflowOutputContractForType(workflowTypeSMTriage)
	if !ok {
		t.Fatal("sm-triage contract missing")
	}

	prompt := buildWorkflowPrompt(
		"all_reviews",
		"You are the managing partner final synthesizer.",
		"=== WORKFLOW METADATA ===\nworkflow_type: sm-triage\n=== END WORKFLOW METADATA ===",
		"=== WORKFLOW PURPOSE AND OUTPUT CONTRACT SUMMARY ===\nworkflow_type: sm-triage\n=== END WORKFLOW PURPOSE AND OUTPUT CONTRACT SUMMARY ===",
		renderWorkflowContract(contract),
		"case context",
		inputs,
		"knowledge",
		"input manifest",
		"prior outputs",
	)

	for _, fragment := range []string{
		"=== FULL WORKFLOW OUTPUT CONTRACT ===",
		"Follow this contract exactly",
		"## Scope / Authority Fit",
		"## Ripeness",
		"## Requested Directive",
	} {
		if !strings.Contains(prompt, fragment) {
			t.Fatalf("final prompt missing fragment %q", fragment)
		}
	}
}
