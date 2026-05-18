package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGuidedChoiceMapping(t *testing.T) {
	tests := []struct {
		id   int
		kind string
		typ  string
	}{
		{1, "workflow", "triage"},
		{2, "workflow", "triage"},
		{3, "workflow", "evidence-packet"},
		{4, "workflow", "counsel-brief"},
		{5, "review", ""},
		{6, "workflow", "sm-triage"},
		{7, "workflow", "draft-message"},
		{8, "workflow", "pattern-review"},
		{9, "workflow", "fact-lock"},
		{10, "workflow", "authority-check"},
		{11, "workflow", "prep"},
		{12, "workflow", "journal-entry"},
		{13, "intake", ""},
	}
	for _, tt := range tests {
		choice, err := resolveGuidedChoice(tt.id, "")
		if err != nil {
			t.Fatalf("choice %d: %v", tt.id, err)
		}
		if choice.Kind != tt.kind || choice.WorkflowType != tt.typ {
			t.Fatalf("choice %d mapped to kind=%q type=%q", tt.id, choice.Kind, choice.WorkflowType)
		}
	}
}

func TestGuidedIntentAliases(t *testing.T) {
	choice, err := resolveGuidedChoice(0, "counsel-brief")
	if err != nil {
		t.Fatal(err)
	}
	if choice.ID != 4 {
		t.Fatalf("expected counsel-brief -> 4, got %d", choice.ID)
	}
}

func TestGuidedEquivalentCommand(t *testing.T) {
	plan := guidedRunPlan{
		Matter:       "Revere-042926",
		Kind:         "workflow",
		WorkflowType: "sm-triage",
		Audience:     "special-master",
		Urgency:      "high",
		ChildRelated: true,
	}
	cmd := guidedEquivalentCommand(plan)
	for _, fragment := range []string{"workflow Revere-042926", "--type sm-triage", "--audience special-master", "--urgency high", "--child-related"} {
		if !strings.Contains(cmd, fragment) {
			t.Fatalf("command missing %q: %s", fragment, cmd)
		}
	}
}

func TestGuidedChoiceTwoUsesPursuePreserveLetGoGoal(t *testing.T) {
	choice, err := resolveGuidedChoice(2, "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(choice.DefaultGoal, "preserved for pattern evidence") {
		t.Fatalf("unexpected goal: %s", choice.DefaultGoal)
	}
}

func TestDetectMatterStatusFindsNewestDraft(t *testing.T) {
	root := t.TempDir()
	drafts := filepath.Join(root, "matters", "Matter-A", "drafts")
	if err := os.MkdirAll(drafts, 0755); err != nil {
		t.Fatal(err)
	}
	oldFile := filepath.Join(drafts, "old.md")
	newFile := filepath.Join(drafts, "new.md")
	if err := os.WriteFile(oldFile, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(newFile, []byte("new"), 0644); err != nil {
		t.Fatal(err)
	}
	status := detectMatterStatus(root, "Matter-A")
	if filepath.Base(status.NewestDraft) != "new.md" {
		t.Fatalf("expected newest draft new.md, got %s", status.NewestDraft)
	}
}

func TestMaybeWriteGuidedRequestSkipsDryRun(t *testing.T) {
	root := t.TempDir()
	workingDir := filepath.Join(root, "working")
	status := matterStatus{WorkingDir: workingDir}
	plan := guidedRunPlan{
		Matter:       "Matter-A",
		Kind:         "workflow",
		WorkflowType: "triage",
		Situation:    "Test situation",
		DryRun:       true,
	}

	if err := maybeWriteGuidedRequest(status, &plan); err != nil {
		t.Fatal(err)
	}
	if plan.GuidedFilePath != "" {
		t.Fatalf("expected no guided file path in dry-run, got %q", plan.GuidedFilePath)
	}
	if _, err := os.Stat(workingDir); !os.IsNotExist(err) {
		t.Fatalf("expected working dir to remain absent, stat err = %v", err)
	}
}
