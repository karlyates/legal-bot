package main

import (
	"path/filepath"
	"testing"

	"github.com/jdonohoo/legal-bot/go/internal/config"
)

func TestIntakeSummary(t *testing.T) {
	tests := []struct {
		name       string
		totalSteps int
		failed     []string
		wantTitle  string
		wantErr    bool
	}{
		{
			name:       "all succeed",
			totalSteps: 5,
			failed:     nil,
			wantTitle:  "Intake Complete",
			wantErr:    false,
		},
		{
			name:       "partial failure",
			totalSteps: 5,
			failed:     []string{"Source Document Cards"},
			wantTitle:  "Intake Completed with Failures",
			wantErr:    true,
		},
		{
			name:       "all failed",
			totalSteps: 5,
			failed:     []string{"Document Intake and Registration", "Source Document Cards", "Chronology Construction", "Atomic Fact Extraction", "Attorney Prep Questions"},
			wantTitle:  "Intake Failed",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, note, err := intakeSummary(tt.totalSteps, tt.failed)
			if title != tt.wantTitle {
				t.Fatalf("intakeSummary() title = %q, want %q", title, tt.wantTitle)
			}
			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if tt.wantErr && note == "" {
				t.Fatal("expected failure note, got empty string")
			}
		})
	}
}

func TestIntakeOutputFileForStepUsesSemanticIdentity(t *testing.T) {
	workingDir := filepath.Join("matter", "working")
	knowledgeDir := filepath.Join("matter", "knowledge")

	tests := []struct {
		name string
		step config.PipelineStep
		want string
	}{
		{
			name: "chronology keeps timeline file after resequencing",
			step: config.PipelineStep{Step: 2, Persona: "chronology-clerk", Name: "Chronology Construction"},
			want: filepath.Join(knowledgeDir, "case_timeline.md"),
		},
		{
			name: "atomic facts keep fact candidate file after resequencing",
			step: config.PipelineStep{Step: 3, Persona: "atomic-fact-extractor", Name: "Atomic Fact Extraction"},
			want: filepath.Join(knowledgeDir, "fact_candidates.md"),
		},
		{
			name: "fallback still uses step number and persona",
			step: config.PipelineStep{Step: 7, Persona: "custom-persona", Name: "Custom Step"},
			want: filepath.Join(workingDir, "07-custom-persona.md"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := intakeOutputFileForStep(tt.step, workingDir, knowledgeDir)
			if got != tt.want {
				t.Fatalf("intakeOutputFileForStep() = %q, want %q", got, tt.want)
			}
		})
	}
}
