package main

import (
	"bufio"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/jdonohoo/legal-bot/go/internal/config"
)

type analysisLevel string

const (
	analysisLevelScout    analysisLevel = "scout"
	analysisLevelStandard analysisLevel = "standard"
	analysisLevelDeep     analysisLevel = "deep"
	analysisLevelMax      analysisLevel = "max"
)

type planStep struct {
	Step         int
	Name         string
	Persona      string
	ModelProfile string
	Route        string
	Warnings     []string
}

type runPlanSummary struct {
	CommandType string
	Matter      string
	Level       analysisLevel
	PipelineKey string
	Steps       []planStep
}

func parseAnalysisLevel(raw string) (analysisLevel, error) {
	switch normalizeRouteName(raw) {
	case "", string(analysisLevelScout):
		return analysisLevelScout, nil
	case string(analysisLevelStandard):
		return analysisLevelStandard, nil
	case string(analysisLevelDeep):
		return analysisLevelDeep, nil
	case string(analysisLevelMax):
		return analysisLevelMax, nil
	default:
		return "", fmt.Errorf("invalid --level %q (supported: scout, standard, deep, max)", raw)
	}
}

func (l analysisLevel) requiresConfirmation() bool {
	return l == analysisLevelDeep || l == analysisLevelMax
}

func applyIntakeLevel(steps []config.PipelineStep, level analysisLevel) []config.PipelineStep {
	selected := clonePipelineSteps(steps)
	switch level {
	case analysisLevelScout:
		selected = filterStepsByOriginalNumber(selected, map[int]bool{1: true})
		overrideProfilesByOriginalStep(selected, map[int]string{1: "cheap_fast"})
	case analysisLevelStandard:
		selected = filterStepsByOriginalNumber(selected, map[int]bool{1: true, 3: true, 4: true, 5: true})
		overrideProfilesByOriginalStep(selected, map[int]string{
			1: "cheap_fast",
			3: "structured_extraction",
			4: "structured_extraction",
			5: "cheap_fast",
		})
	case analysisLevelDeep:
		overrideProfilesByOriginalStep(selected, map[int]string{5: "cheap_fast"})
	case analysisLevelMax:
	}
	resequenceSteps(selected)
	return selected
}

func applyWorkflowLevel(steps []config.PipelineStep, level analysisLevel) []config.PipelineStep {
	selected := clonePipelineSteps(steps)

	switch level {
	case analysisLevelScout:
		selected = selectWorkflowScoutSteps(selected)
		for i := range selected {
			if selected[i].Persona == "managing-partner-final-synthesizer" {
				selected[i].ModelProfile = "final_synthesis_scout"
			} else {
				selected[i].ModelProfile = "cheap_fast"
			}
		}
	case analysisLevelStandard:
		priorities := []string{
			"chronology-clerk",
			"atomic-fact-extractor",
			"family-law-attorney-reviewer",
			personaChildBestInterests,
			personaFinancialSupport,
			"opposing-counsel",
			"practical-resolution-reviewer",
			personaLegalAuthority,
		}
		selected = selectPrioritySteps(selected, priorities, 5, true, "managing-partner-final-synthesizer")
		for i := range selected {
			selected[i].ModelProfile = downgradedWorkflowProfile(selected[i], level)
		}
	case analysisLevelDeep:
		for i := range selected {
			selected[i].ModelProfile = downgradedWorkflowProfile(selected[i], level)
		}
	case analysisLevelMax:
	}

	resequenceSteps(selected)
	return selected
}

func applyReviewLevel(steps []config.PipelineStep, level analysisLevel) []config.PipelineStep {
	selected := clonePipelineSteps(steps)

	switch level {
	case analysisLevelScout:
		selected = selectReviewScoutSteps(selected)
		for i := range selected {
			selected[i].ModelProfile = downgradedReviewProfile(selected[i], level)
		}
	case analysisLevelStandard:
		priorities := []string{
			"litigation-paralegal",
			"trial-fact-checker",
			personaChildBestInterests,
			personaFinancialSupport,
			"family-law-attorney-reviewer",
			"opposing-counsel",
			"neutral-court-reader",
			"practical-resolution-reviewer",
			personaLegalAuthority,
			"legal-writing-preservation-editor",
		}
		selected = selectPrioritySteps(selected, priorities, 5, false, "managing-partner-final-synthesizer")
		for i := range selected {
			selected[i].ModelProfile = downgradedReviewProfile(selected[i], level)
		}
	case analysisLevelDeep:
		for i := range selected {
			selected[i].ModelProfile = downgradedReviewProfile(selected[i], level)
		}
	case analysisLevelMax:
	}

	resequenceSteps(selected)
	return selected
}

func downgradedWorkflowProfile(step config.PipelineStep, level analysisLevel) string {
	if level == analysisLevelMax {
		return step.ModelProfile
	}
	switch step.Persona {
	case "managing-partner-final-synthesizer":
		switch level {
		case analysisLevelScout:
			return "final_synthesis_scout"
		case analysisLevelStandard:
			return "final_synthesis_standard"
		default:
			return step.ModelProfile
		}
	case "chronology-clerk", "atomic-fact-extractor", "source-document-analyst":
		if level == analysisLevelScout {
			return "cheap_fast"
		}
		return "structured_extraction"
	case personaStrategicOptions:
		if level == analysisLevelStandard {
			return "strategy_reasoning_standard"
		}
		return step.ModelProfile
	case personaLegalAuthority:
		if level == analysisLevelScout {
			return "cheap_fast"
		}
		if level == analysisLevelStandard {
			return "review_reasoning_standard"
		}
		return step.ModelProfile
	default:
		if level == analysisLevelScout {
			return "cheap_fast"
		}
		if level == analysisLevelStandard {
			return "review_reasoning_standard"
		}
		return step.ModelProfile
	}
}

func downgradedReviewProfile(step config.PipelineStep, level analysisLevel) string {
	if level == analysisLevelMax {
		return step.ModelProfile
	}
	switch step.Persona {
	case "managing-partner-final-synthesizer":
		switch level {
		case analysisLevelScout:
			return "final_synthesis_scout"
		case analysisLevelStandard:
			return "final_synthesis_standard"
		default:
			return step.ModelProfile
		}
	case "chronology-clerk", "atomic-fact-extractor", "source-document-analyst":
		if level == analysisLevelScout {
			return "cheap_fast"
		}
		return "structured_extraction"
	case personaStrategicOptions:
		if level == analysisLevelStandard {
			return "strategy_reasoning_standard"
		}
		return step.ModelProfile
	case personaLegalAuthority:
		if level == analysisLevelScout {
			return "cheap_fast"
		}
		if level == analysisLevelStandard {
			return "review_reasoning_standard"
		}
		return step.ModelProfile
	default:
		if level == analysisLevelScout {
			return "review_reasoning_scout"
		}
		if level == analysisLevelStandard {
			return "review_reasoning_standard"
		}
		return step.ModelProfile
	}
}

func buildRunPlanSummary(cfg *config.Config, commandType string, matter string, level analysisLevel, pipelineKey string, steps []config.PipelineStep) runPlanSummary {
	summary := runPlanSummary{
		CommandType: commandType,
		Matter:      matter,
		Level:       level,
		PipelineKey: pipelineKey,
	}

	for _, step := range steps {
		planned := planStep{
			Step:         step.Step,
			Name:         step.Name,
			Persona:      step.Persona,
			ModelProfile: step.ModelProfile,
		}

		if selection, err := selectStepRoute(cfg, step.ModelProfile, step.LLM); err == nil {
			planned.Route = describeRoute(step.ModelProfile, selection)
			planned.Warnings = append(planned.Warnings, selection.Warnings...)
		} else {
			planned.Route = "unresolved"
			planned.Warnings = append(planned.Warnings, err.Error())
		}

		summary.Steps = append(summary.Steps, planned)
	}

	return summary
}

func printRunPlan(w io.Writer, plan runPlanSummary) {
	fmt.Fprintf(w, "Command type: %s\n", plan.CommandType)
	fmt.Fprintf(w, "Matter: %s\n", plan.Matter)
	fmt.Fprintf(w, "Selected level: %s\n", plan.Level)
	fmt.Fprintf(w, "Pipeline key: %s\n", plan.PipelineKey)
	fmt.Fprintf(w, "Planned LLM steps: %d\n", len(plan.Steps))
	fmt.Fprintln(w, "Steps:")
	for _, step := range plan.Steps {
		fmt.Fprintf(w, "%d. %s\n", step.Step, step.Name)
		fmt.Fprintf(w, "   persona: %s\n", step.Persona)
		fmt.Fprintf(w, "   model profile: %s\n", step.ModelProfile)
		fmt.Fprintf(w, "   resolved route: %s\n", step.Route)
		for _, warning := range step.Warnings {
			if strings.TrimSpace(warning) == "" {
				continue
			}
			fmt.Fprintf(w, "   warning: %s\n", warning)
		}
	}
	if len(plan.Steps) > 5 {
		fmt.Fprintf(w, "Warning: this run will invoke %d LLM steps. Consider --level scout or --plan first.\n", len(plan.Steps))
	}
	fmt.Fprintln(w, "No LLMs were called.")
}

func maybeWarnLargeRun(w io.Writer, steps []config.PipelineStep) {
	if len(steps) > 5 {
		fmt.Fprintf(w, "Warning: this run will invoke %d LLM steps. Consider --level scout or --plan first.\n", len(steps))
	}
}

func maybeConfirmLargeRun(kind string, level analysisLevel, steps []config.PipelineStep, yes bool) error {
	if yes || !level.requiresConfirmation() {
		return nil
	}
	reader := bufio.NewReader(os.Stdin)
	label := fmt.Sprintf("This %s run at level %s will invoke %d LLM steps. Continue? [y/N] ", kind, level, len(steps))
	if !confirmWithDefault(reader, os.Stdout, label, false) {
		return fmt.Errorf("%s run cancelled", kind)
	}
	return nil
}

func clonePipelineSteps(steps []config.PipelineStep) []config.PipelineStep {
	cloned := make([]config.PipelineStep, len(steps))
	copy(cloned, steps)
	return cloned
}

func resequenceSteps(steps []config.PipelineStep) {
	for i := range steps {
		steps[i].Step = i + 1
	}
}

func filterStepsByOriginalNumber(steps []config.PipelineStep, keep map[int]bool) []config.PipelineStep {
	filtered := make([]config.PipelineStep, 0, len(steps))
	for _, step := range steps {
		if keep[step.Step] {
			filtered = append(filtered, step)
		}
	}
	return filtered
}

func overrideProfilesByOriginalStep(steps []config.PipelineStep, overrides map[int]string) {
	for i := range steps {
		if profile, ok := overrides[steps[i].Step]; ok {
			steps[i].ModelProfile = profile
		}
	}
}

func selectWorkflowScoutSteps(steps []config.PipelineStep) []config.PipelineStep {
	selected := map[int]bool{}
	if len(steps) == 0 {
		return nil
	}
	selected[0] = true

	substantiveOrder := []string{
		"family-law-attorney-reviewer",
		"practical-resolution-reviewer",
		"neutral-court-reader",
		"opposing-counsel",
	}
	if idx := firstMatchingPersonaIndex(steps, substantiveOrder, selected); idx >= 0 {
		selected[idx] = true
	} else if idx := firstNonFinalIndex(steps, selected); idx >= 0 {
		selected[idx] = true
	}

	if idx := personaIndex(steps, "managing-partner-final-synthesizer"); idx >= 0 {
		selected[idx] = true
	}

	return stepsFromSelection(steps, selected, 3)
}

func selectReviewScoutSteps(steps []config.PipelineStep) []config.PipelineStep {
	selected := map[int]bool{}
	if idx := firstMatchingPersonaIndex(steps, []string{"trial-fact-checker", "litigation-paralegal"}, selected); idx >= 0 {
		selected[idx] = true
	}
	if idx := firstMatchingPersonaIndex(steps, []string{"opposing-counsel", "neutral-court-reader"}, selected); idx >= 0 {
		selected[idx] = true
	}
	if idx := personaIndex(steps, "managing-partner-final-synthesizer"); idx >= 0 {
		selected[idx] = true
	}
	if len(selected) == 0 && len(steps) > 0 {
		selected[0] = true
	}
	return stepsFromSelection(steps, selected, 3)
}

func selectPrioritySteps(steps []config.PipelineStep, priorities []string, cap int, includeFirst bool, finalPersona string) []config.PipelineStep {
	selected := map[int]bool{}
	if includeFirst && len(steps) > 0 {
		selected[0] = true
	}

	finalIdx := personaIndex(steps, finalPersona)
	reserved := 0
	if finalIdx >= 0 && !selected[finalIdx] {
		reserved = 1
	}

	for _, persona := range priorities {
		if len(selected) >= cap-reserved {
			break
		}
		if idx := personaIndex(steps, persona); idx >= 0 && !selected[idx] {
			selected[idx] = true
		}
	}

	if finalIdx >= 0 && len(selected) < cap {
		selected[finalIdx] = true
	}

	if len(selected) < cap {
		for i := range steps {
			if len(selected) >= cap {
				break
			}
			selected[i] = true
		}
	}

	return stepsFromSelection(steps, selected, cap)
}

func stepsFromSelection(steps []config.PipelineStep, selected map[int]bool, cap int) []config.PipelineStep {
	indexes := slices.Collect(maps.Keys(selected))
	slices.Sort(indexes)
	if cap > 0 && len(indexes) > cap {
		indexes = indexes[:cap]
	}
	result := make([]config.PipelineStep, 0, len(indexes))
	for _, idx := range indexes {
		if idx >= 0 && idx < len(steps) {
			result = append(result, steps[idx])
		}
	}
	return result
}

func personaIndex(steps []config.PipelineStep, persona string) int {
	for i, step := range steps {
		if step.Persona == persona {
			return i
		}
	}
	return -1
}

func firstMatchingPersonaIndex(steps []config.PipelineStep, personas []string, selected map[int]bool) int {
	for _, persona := range personas {
		for i, step := range steps {
			if selected[i] {
				continue
			}
			if step.Persona == persona {
				return i
			}
		}
	}
	return -1
}

func firstNonFinalIndex(steps []config.PipelineStep, selected map[int]bool) int {
	for i, step := range steps {
		if selected[i] || step.Persona == "managing-partner-final-synthesizer" {
			continue
		}
		return i
	}
	return -1
}
