package config

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

var validEngines = map[string]struct{}{
	"claude":  {},
	"codex":   {},
	"gemini":  {},
	"copilot": {},
}

var validEfforts = map[string]struct{}{
	"minimal": {},
	"low":     {},
	"medium":  {},
	"high":    {},
	"xhigh":   {},
	"max":     {},
}

type ModelTargetSpec struct {
	Engine       string `json:"engine,omitempty"`
	Model        string `json:"model,omitempty"`
	Effort       string `json:"effort,omitempty"`
	Raw          string `json:"-"`
	Set          bool   `json:"-"`
	LegacyOpenAI bool   `json:"-"`
}

type ResolvedModelTarget struct {
	Engine                string
	Model                 string
	Effort                string
	SourceProfile         string
	SourceSpec            string
	FallbackIndex         int
	LegacyAliasNormalized bool
}

func (m *ModelTargetSpec) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*m = ModelTargetSpec{}
		return nil
	}

	if len(data) > 0 && data[0] == '"' {
		var raw string
		if err := json.Unmarshal(data, &raw); err != nil {
			return err
		}
		spec, err := ParseModelTargetSpec(raw)
		if err != nil {
			return err
		}
		*m = spec
		return nil
	}

	type alias ModelTargetSpec
	var decoded alias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	decoded.Engine = normalizeEngine(decoded.Engine)
	decoded.Effort = normalizeEffort(decoded.Effort)
	decoded.Set = decoded.Engine != "" || decoded.Model != "" || decoded.Effort != ""
	*m = ModelTargetSpec(decoded)
	return nil
}

func ParseModelTargetSpec(raw string) (ModelTargetSpec, error) {
	spec := ModelTargetSpec{Raw: strings.TrimSpace(raw), Set: true}
	if spec.Raw == "" {
		return spec, fmt.Errorf("empty model target")
	}

	parts := strings.Split(spec.Raw, ":")
	switch len(parts) {
	case 1:
		spec.Engine = normalizeEngine(parts[0])
		if spec.Engine == "" {
			spec.Engine = strings.ToLower(strings.TrimSpace(parts[0]))
		}
	case 2, 3:
		first := strings.ToLower(strings.TrimSpace(parts[0]))
		if first == "openai" {
			spec.Engine = "codex"
			spec.LegacyOpenAI = true
		} else {
			spec.Engine = normalizeEngine(first)
			if spec.Engine == "" {
				spec.Engine = first
			}
		}
		spec.Model = strings.TrimSpace(parts[1])
		if len(parts) == 3 {
			spec.Effort = normalizeEffort(parts[2])
		}
	default:
		return spec, fmt.Errorf("invalid compact model target %q", raw)
	}

	return spec, nil
}

func (m ModelTargetSpec) IsZero() bool {
	return !m.Set && m.Engine == "" && m.Model == "" && m.Effort == ""
}

func (m ModelTargetSpec) String() string {
	if m.Raw != "" {
		return m.Raw
	}
	parts := []string{}
	if m.Engine != "" {
		parts = append(parts, m.Engine)
	}
	if m.Model != "" {
		parts = append(parts, m.Model)
	}
	if m.Effort != "" {
		parts = append(parts, m.Effort)
	}
	return strings.Join(parts, ":")
}

func (c *Config) ResolveModelProfileTargets(profileName string, fallbackLLM string) ([]ResolvedModelTarget, []string) {
	warnings := []string{}
	if profileName == "" || c.ModelProfiles == nil {
		return appendFallbackTarget(nil, profileName, fallbackLLM), warnings
	}

	profile, ok := c.ModelProfiles[profileName]
	if !ok {
		warnings = append(warnings, fmt.Sprintf("model profile %q not found; using step fallback %q", profileName, fallbackLLM))
		return appendFallbackTarget(nil, profileName, fallbackLLM), warnings
	}

	targets := []ResolvedModelTarget{}
	for idx, spec := range profile.targets() {
		if spec.IsZero() {
			continue
		}
		target, targetWarnings := resolveTargetSpec(profileName, idx, spec)
		warnings = append(warnings, targetWarnings...)
		targets = append(targets, target)
	}

	targets = appendFallbackTarget(targets, profileName, fallbackLLM)
	return targets, dedupeStrings(warnings)
}

func (c *Config) GetOverrideTarget() (*ResolvedModelTarget, []string, error) {
	mode := c.getActiveMode()
	if mode == nil || strings.TrimSpace(mode.OverrideLLM) == "" {
		return nil, nil, nil
	}

	spec, err := ParseModelTargetSpec(mode.OverrideLLM)
	if err != nil {
		return nil, nil, fmt.Errorf("parse single_llm override %q: %w", mode.OverrideLLM, err)
	}

	target, warnings := resolveTargetSpec("single_llm", 0, spec)
	return &target, warnings, nil
}

func SelectRunnableTarget(targets []ResolvedModelTarget) (ResolvedModelTarget, []string, error) {
	warnings := []string{}
	tried := []string{}

	for _, target := range targets {
		targetLabel := target.ShortString()
		tried = append(tried, targetLabel)

		if !isKnownEngine(target.Engine) {
			warnings = append(warnings, fmt.Sprintf("unknown engine %q in %s", target.Engine, targetLabel))
			continue
		}

		if _, err := exec.LookPath(target.Engine); err == nil {
			return target, warnings, nil
		}

		warnings = append(warnings, fmt.Sprintf("%s not found on PATH for %s", target.Engine, targetLabel))
	}

	if len(tried) == 0 {
		return ResolvedModelTarget{}, warnings, fmt.Errorf("no model targets were configured")
	}

	return ResolvedModelTarget{}, warnings, fmt.Errorf("no configured model engine is installed; tried %s", strings.Join(tried, ", "))
}

func (r ResolvedModelTarget) ShortString() string {
	parts := []string{r.Engine}
	if r.Model != "" {
		parts = append(parts, r.Model)
	}
	if r.Effort != "" {
		parts = append(parts, r.Effort)
	}
	return strings.Join(parts, " / ")
}

func (r ResolvedModelTarget) DebugString() string {
	label := r.ShortString()
	if r.FallbackIndex > 0 {
		label += fmt.Sprintf(" (fallback[%d])", r.FallbackIndex-1)
	}
	if r.LegacyAliasNormalized {
		label += " [legacy openai:* -> codex]"
	}
	return label
}

func (m ModelProfileConfig) targets() []ModelTargetSpec {
	targets := []ModelTargetSpec{}
	if !m.Primary.IsZero() {
		targets = append(targets, m.Primary)
	}
	if !m.Preferred.IsZero() {
		targets = append(targets, m.Preferred)
	}
	if len(m.Fallbacks) > 0 {
		targets = append(targets, m.Fallbacks...)
	}
	if !m.Fallback.IsZero() {
		targets = append(targets, m.Fallback)
	}
	return targets
}

func resolveTargetSpec(profileName string, idx int, spec ModelTargetSpec) (ResolvedModelTarget, []string) {
	warnings := []string{}
	engine := normalizeEngine(spec.Engine)
	if engine == "" {
		engine = strings.ToLower(strings.TrimSpace(spec.Engine))
	}

	effort := normalizeEffort(spec.Effort)
	if effort == "" && spec.Effort != "" {
		warnings = append(warnings, fmt.Sprintf("profile %q target %q uses unsupported effort %q", profileName, spec.String(), spec.Effort))
	}
	if effort == "" && spec.LegacyOpenAI {
		effort = defaultLegacyOpenAIEffort(profileName)
	}

	target := ResolvedModelTarget{
		Engine:                engine,
		Model:                 strings.TrimSpace(spec.Model),
		Effort:                effort,
		SourceProfile:         profileName,
		SourceSpec:            spec.String(),
		FallbackIndex:         idx,
		LegacyAliasNormalized: spec.LegacyOpenAI,
	}

	if !isKnownEngine(target.Engine) && target.Engine != "" {
		warnings = append(warnings, fmt.Sprintf("profile %q target %q uses unknown engine %q", profileName, spec.String(), target.Engine))
	}
	if target.Effort != "" && (target.Engine == "gemini" || target.Engine == "copilot") {
		warnings = append(warnings, fmt.Sprintf("%s ignores effort %q for target %q", target.Engine, target.Effort, spec.String()))
	}
	if spec.LegacyOpenAI && target.Model == "" {
		warnings = append(warnings, fmt.Sprintf("legacy target %q is missing a model", spec.String()))
	}

	return target, warnings
}

func appendFallbackTarget(targets []ResolvedModelTarget, profileName string, fallbackLLM string) []ResolvedModelTarget {
	fallbackLLM = strings.TrimSpace(fallbackLLM)
	if fallbackLLM == "" {
		return targets
	}

	spec, err := ParseModelTargetSpec(fallbackLLM)
	if err != nil {
		return targets
	}
	target, _ := resolveTargetSpec(profileName, len(targets), spec)
	for _, existing := range targets {
		if existing.Engine == target.Engine && existing.Model == target.Model && existing.Effort == target.Effort {
			return targets
		}
	}
	return append(targets, target)
}

func normalizeEngine(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "claude", "c":
		return "claude"
	case "codex", "x":
		return "codex"
	case "gemini", "g":
		return "gemini"
	case "copilot", "p":
		return "copilot"
	default:
		return ""
	}
}

func normalizeEffort(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if _, ok := validEfforts[value]; ok {
		return value
	}
	return ""
}

func isKnownEngine(value string) bool {
	_, ok := validEngines[strings.ToLower(strings.TrimSpace(value))]
	return ok
}

func defaultLegacyOpenAIEffort(profileName string) string {
	switch profileName {
	case "cheap_fast":
		return "low"
	case "structured_extraction":
		return "medium"
	default:
		return "medium"
	}
}

func dedupeStrings(values []string) []string {
	if len(values) == 0 {
		return values
	}
	seen := map[string]struct{}{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
