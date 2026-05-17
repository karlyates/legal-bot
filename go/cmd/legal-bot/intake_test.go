package main

import "testing"

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
