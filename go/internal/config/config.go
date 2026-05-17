package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jdonohoo/legal-bot/go/internal/embedded"
)

// Config represents the legal-bot configuration.
type Config struct {
	Version        string                        `json:"version"`
	TimeoutSeconds int                           `json:"timeout_seconds"`
	MaxRetries     int                           `json:"max_retries"`
	PipelineMode   string                        `json:"pipeline_mode"`
	Pipelines      map[string][]PipelineStep     `json:"discovery_pipelines"`
	LegalPipelines map[string][]PipelineStep     `json:"legal_pipelines"`
	ModelProfiles  map[string]ModelProfileConfig `json:"model_profiles"`
	LLMs           map[string]bool               `json:"llms"`
	LLMMode        string                        `json:"llm_mode"`
	LLMModes       map[string]LLMModeConfig      `json:"llm_modes"`
	VernHole       VernHoleConfig                `json:"vernhole"`
	Timeouts       TimeoutConfig                 `json:"timeouts"`

	// User preferences (persisted across sessions)
	DefaultDiscoveryPath string `json:"default_discovery_path,omitempty"`

	// Backward compat: old config format
	LegacyPipeline []PipelineStep `json:"discovery_pipeline"`

	// SourcePath is the file this config was loaded from (not serialized).
	SourcePath string `json:"-"`
}

// LLMModeConfig defines fallback behavior for a given LLM mode.
type LLMModeConfig struct {
	Description  string            `json:"description"`
	Fallback     map[string]string `json:"fallback"`
	SynthesisLLM string            `json:"synthesis_llm"`
	OverrideLLM  string            `json:"override_llm,omitempty"`
}

// ModelProfileConfig maps a provider-neutral model profile to concrete models.
type ModelProfileConfig struct {
	Preferred ModelTargetSpec   `json:"preferred,omitempty"`
	Fallback  ModelTargetSpec   `json:"fallback,omitempty"`
	Primary   ModelTargetSpec   `json:"primary,omitempty"`
	Fallbacks []ModelTargetSpec `json:"fallbacks,omitempty"`
}

// PipelineStep defines a single step in a discovery pipeline.
type PipelineStep struct {
	Step         int    `json:"step"`
	Name         string `json:"name"`
	Persona      string `json:"persona"`
	LLM          string `json:"llm"`
	ModelProfile string `json:"model_profile,omitempty"`
	ContextMode  string `json:"context_mode"`
	PromptPrefix string `json:"prompt_prefix"`
}

// TimeoutConfig holds granular timeout settings (all in seconds).
type TimeoutConfig struct {
	PipelineStep int `json:"pipeline_step"` // per discovery pipeline step
	Historian    int `json:"historian"`     // historian indexing
	Oracle       int `json:"oracle"`        // oracle consult
	OracleApply  int `json:"oracle_apply"`  // architect applying oracle vision
}

// VernHoleConfig holds VernHole-specific settings.
type VernHoleConfig struct {
	DefaultCouncil string `json:"default_council"`
	Min            int    `json:"min"`
}

// Load reads configuration using the 4-tier chain:
//  1. ~/.claude/legal-bot-config.json (user config, Claude Code plugin)
//  2. ~/.config/legal-bot/config.json (standalone user config)
//  3. {projectRoot}/config.default.json (project defaults)
//  4. Hardcoded defaults
func Load(projectRoot string) *Config {
	// Tier 1: Claude Code plugin user config
	userConfig := filepath.Join(os.Getenv("HOME"), ".claude", "legal-bot-config.json")
	if cfg, err := loadFile(userConfig); err == nil {
		cfg.SourcePath = userConfig
		return cfg
	}

	// Tier 2: standalone user config
	standaloneConfig := filepath.Join(os.Getenv("HOME"), ".config", "legal-bot", "config.json")
	if cfg, err := loadFile(standaloneConfig); err == nil {
		cfg.SourcePath = standaloneConfig
		return cfg
	}

	// Tier 3: project default config (on disk)
	defaultConfig := filepath.Join(projectRoot, "config.default.json")
	if cfg, err := loadFile(defaultConfig); err == nil {
		cfg.SourcePath = defaultConfig
		return cfg
	}

	// Tier 4: embedded default config (compiled into binary)
	if cfg, err := loadEmbeddedConfig(); err == nil {
		return cfg
	}

	// Tier 5: hardcoded defaults
	return hardcodedDefaults()
}

// LoadProjectDefault loads the repo's config.default.json directly, falling
// back to the embedded config and then hardcoded defaults.
func LoadProjectDefault(projectRoot string) (*Config, error) {
	defaultConfig := filepath.Join(projectRoot, "config.default.json")
	if cfg, err := loadFile(defaultConfig); err == nil {
		cfg.SourcePath = defaultConfig
		return cfg, nil
	}
	if cfg, err := loadEmbeddedConfig(); err == nil {
		return cfg, nil
	}
	return hardcodedDefaults(), nil
}

func loadEmbeddedConfig() (*Config, error) {
	data := embedded.GetDefaultConfig()
	if data == "" {
		return nil, fmt.Errorf("no embedded config")
	}

	cfg := &Config{}
	if err := json.Unmarshal([]byte(data), cfg); err != nil {
		return nil, fmt.Errorf("parse embedded config: %w", err)
	}

	// Apply same defaults as loadFile
	if cfg.TimeoutSeconds == 0 {
		cfg.TimeoutSeconds = 1200
	}
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 2
	}
	if cfg.PipelineMode == "" {
		cfg.PipelineMode = "default"
	}
	if cfg.VernHole.Min == 0 {
		cfg.VernHole.Min = 3
	}
	if cfg.VernHole.DefaultCouncil == "" {
		cfg.VernHole.DefaultCouncil = "random"
	}
	if cfg.LLMMode == "" {
		cfg.LLMMode = "mixed_claude_fallback"
	}
	if cfg.LLMModes == nil {
		cfg.LLMModes = defaultLLMModes()
	}
	if cfg.ModelProfiles == nil {
		cfg.ModelProfiles = defaultModelProfiles()
	}
	applyTimeoutDefaults(cfg)

	return cfg, nil
}

func loadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	// Backward compat: migrate legacy discovery_pipeline to discovery_pipelines
	if cfg.Pipelines == nil && cfg.LegacyPipeline != nil {
		cfg.Pipelines = map[string][]PipelineStep{
			"default": cfg.LegacyPipeline,
		}
		cfg.PipelineMode = "default"
	}

	// Apply defaults for missing fields
	if cfg.TimeoutSeconds == 0 {
		cfg.TimeoutSeconds = 1200
	}
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 2
	}
	if cfg.PipelineMode == "" {
		cfg.PipelineMode = "default"
	}
	if cfg.VernHole.Min == 0 {
		cfg.VernHole.Min = 3
	}
	if cfg.VernHole.DefaultCouncil == "" {
		cfg.VernHole.DefaultCouncil = "random"
	}
	if cfg.LLMMode == "" {
		cfg.LLMMode = "mixed_claude_fallback"
	}
	if cfg.LLMModes == nil {
		cfg.LLMModes = defaultLLMModes()
	}
	if cfg.ModelProfiles == nil {
		cfg.ModelProfiles = defaultModelProfiles()
	}
	applyTimeoutDefaults(cfg)

	return cfg, nil
}

// GetPipeline returns the pipeline steps for the given mode.
// Falls back to "default" if the mode doesn't exist.
func (c *Config) GetPipeline(mode string) []PipelineStep {
	if mode == "" {
		mode = c.PipelineMode
	}
	if steps, ok := c.Pipelines[mode]; ok {
		return steps
	}
	if steps, ok := c.Pipelines["default"]; ok {
		return steps
	}
	return hardcodedDefaults().Pipelines["default"]
}

// GetLegalPipeline returns the legal pipeline steps for the given legal route.
func (c *Config) GetLegalPipeline(pipelineType string) []PipelineStep {
	if c.LegalPipelines != nil {
		if steps, ok := c.LegalPipelines[pipelineType]; ok {
			return steps
		}
		if pipelineType == "review" {
			if steps, ok := c.LegalPipelines["review_standard"]; ok {
				return steps
			}
		}
	}
	return nil
}

// ResolveModelProfile returns the runnable LLM engine for a model profile.
// This is kept for backward compatibility with older engine-only call sites.
func (c *Config) ResolveModelProfile(profileName string, fallbackLLM string) string {
	targets, _ := c.ResolveModelProfileTargets(profileName, fallbackLLM)
	for _, target := range targets {
		if target.Engine != "" {
			return target.Engine
		}
	}
	return fallbackLLM
}

// GetFallbackLLM returns the fallback LLM for a given original LLM based on the active LLM mode.
// Returns empty string if no fallback is configured.
func (c *Config) GetFallbackLLM(originalLLM string) string {
	mode := c.getActiveMode()
	if mode == nil {
		return "claude" // backward compat default
	}
	if fb, ok := mode.Fallback[originalLLM]; ok && fb != originalLLM {
		return fb
	}
	return ""
}

// GetSynthesisLLM returns the LLM to use for VernHole synthesis.
func (c *Config) GetSynthesisLLM() string {
	mode := c.getActiveMode()
	if mode != nil && mode.SynthesisLLM != "" {
		return mode.SynthesisLLM
	}
	return "claude" // backward compat default
}

// GetOverrideLLM returns the override LLM for single_llm mode, or empty string.
func (c *Config) GetOverrideLLM() string {
	mode := c.getActiveMode()
	if mode != nil && mode.OverrideLLM != "" {
		target, err := ParseModelTargetSpec(mode.OverrideLLM)
		if err == nil && target.Engine != "" {
			return target.Engine
		}
		return mode.OverrideLLM
	}
	return ""
}

// GetPipelineStepTimeout returns the per-step timeout for discovery pipeline steps (seconds).
func (c *Config) GetPipelineStepTimeout() int {
	if c.Timeouts.PipelineStep > 0 {
		return c.Timeouts.PipelineStep
	}
	return 1200
}

// GetHistorianTimeout returns the timeout for historian indexing (seconds).
func (c *Config) GetHistorianTimeout() int {
	if c.Timeouts.Historian > 0 {
		return c.Timeouts.Historian
	}
	return 1200
}

// GetOracleTimeout returns the timeout for oracle consult (seconds).
func (c *Config) GetOracleTimeout() int {
	if c.Timeouts.Oracle > 0 {
		return c.Timeouts.Oracle
	}
	return 1200
}

// GetOracleApplyTimeout returns the timeout for oracle apply (seconds).
func (c *Config) GetOracleApplyTimeout() int {
	if c.Timeouts.OracleApply > 0 {
		return c.Timeouts.OracleApply
	}
	return 1200
}

func (c *Config) getActiveMode() *LLMModeConfig {
	if c.LLMMode == "" || c.LLMModes == nil {
		return nil
	}
	mode, ok := c.LLMModes[c.LLMMode]
	if !ok {
		return nil
	}
	return &mode
}

// applyTimeoutDefaults populates the granular Timeouts struct from the legacy
// TimeoutSeconds field if needed, ensuring backward compatibility.
func applyTimeoutDefaults(cfg *Config) {
	zero := TimeoutConfig{}
	if cfg.Timeouts == zero {
		base := cfg.TimeoutSeconds
		if base <= 0 {
			base = 1200
		}
		cfg.Timeouts = TimeoutConfig{
			PipelineStep: base,
			Historian:    base,
			Oracle:       base,
			OracleApply:  base,
		}
	}
}

func defaultLLMModes() map[string]LLMModeConfig {
	return map[string]LLMModeConfig{
		"mixed_claude_fallback": {
			Description:  "Mixed LLMs, claude as safety net",
			Fallback:     map[string]string{"codex": "claude", "gemini": "claude", "copilot": "claude"},
			SynthesisLLM: "claude",
		},
		"mixed_codex_fallback": {
			Description:  "Mixed LLMs, codex as safety net",
			Fallback:     map[string]string{"claude": "codex", "gemini": "codex", "copilot": "codex"},
			SynthesisLLM: "codex",
		},
		"mixed_gemini_fallback": {
			Description:  "Mixed LLMs, gemini as safety net",
			Fallback:     map[string]string{"claude": "gemini", "codex": "gemini", "copilot": "gemini"},
			SynthesisLLM: "gemini",
		},
		"mixed_copilot_fallback": {
			Description:  "Mixed LLMs, copilot as safety net",
			Fallback:     map[string]string{"claude": "copilot", "codex": "copilot", "gemini": "copilot"},
			SynthesisLLM: "copilot",
		},
		"single_llm": {
			Description:  "Single LLM for everything",
			OverrideLLM:  "",
			Fallback:     map[string]string{},
			SynthesisLLM: "",
		},
	}
}

func defaultModelProfiles() map[string]ModelProfileConfig {
	return map[string]ModelProfileConfig{
		"intake_long_context": {
			Primary: ModelTargetSpec{Engine: "gemini", Model: "pro", Set: true},
			Fallbacks: []ModelTargetSpec{
				{Engine: "claude", Model: "sonnet", Effort: "medium", Set: true},
				{Engine: "codex", Model: "gpt-5.4", Effort: "medium", Set: true},
			},
		},
		"structured_extraction": {
			Primary: ModelTargetSpec{Engine: "codex", Model: "gpt-5.4-mini", Effort: "medium", Set: true},
			Fallbacks: []ModelTargetSpec{
				{Engine: "gemini", Model: "flash", Set: true},
				{Engine: "claude", Model: "haiku", Effort: "medium", Set: true},
			},
		},
		"review_reasoning": {
			Primary: ModelTargetSpec{Engine: "codex", Model: "gpt-5.4", Effort: "medium", Set: true},
			Fallbacks: []ModelTargetSpec{
				{Engine: "claude", Model: "sonnet", Effort: "medium", Set: true},
				{Engine: "gemini", Model: "pro", Set: true},
			},
		},
		"final_synthesis": {
			Primary: ModelTargetSpec{Engine: "codex", Model: "gpt-5.4", Effort: "high", Set: true},
			Fallbacks: []ModelTargetSpec{
				{Engine: "claude", Model: "opus", Effort: "high", Set: true},
				{Engine: "gemini", Model: "pro", Set: true},
			},
		},
		"cheap_fast": {
			Primary: ModelTargetSpec{Engine: "codex", Model: "gpt-5.4-mini", Effort: "low", Set: true},
			Fallbacks: []ModelTargetSpec{
				{Engine: "gemini", Model: "flash-lite", Set: true},
				{Engine: "claude", Model: "haiku", Effort: "medium", Set: true},
			},
		},
		"strategy_reasoning": {
			Primary: ModelTargetSpec{Engine: "codex", Model: "gpt-5.4", Effort: "high", Set: true},
			Fallbacks: []ModelTargetSpec{
				{Engine: "claude", Model: "opus", Effort: "high", Set: true},
				{Engine: "gemini", Model: "pro", Set: true},
			},
		},
		"legal_research": {
			Primary: ModelTargetSpec{Engine: "codex", Model: "gpt-5.4", Effort: "high", Set: true},
			Fallbacks: []ModelTargetSpec{
				{Engine: "claude", Model: "opus", Effort: "high", Set: true},
				{Engine: "gemini", Model: "pro", Set: true},
			},
		},
	}
}

func hardcodedDefaults() *Config {
	return &Config{
		Version:        "2.9.1",
		TimeoutSeconds: 1200,
		MaxRetries:     2,
		PipelineMode:   "default",
		LLMMode:        "mixed_claude_fallback",
		LLMModes:       defaultLLMModes(),
		ModelProfiles:  defaultModelProfiles(),
		VernHole: VernHoleConfig{
			DefaultCouncil: "random",
			Min:            3,
		},
		Timeouts: TimeoutConfig{
			PipelineStep: 1200,
			Historian:    1200,
			Oracle:       1200,
			OracleApply:  1200,
		},
		LLMs: map[string]bool{
			"claude":  true,
			"codex":   true,
			"gemini":  true,
			"copilot": true,
		},
		Pipelines: map[string][]PipelineStep{
			"default": {
				{Step: 1, Name: "Initial Analysis", Persona: "mighty", LLM: "codex", ContextMode: "prompt_only",
					PromptPrefix: "You are the lead analyst. Analyze this idea and provide comprehensive initial analysis including: problem space, technical requirements, proposed architecture, unknowns and risks."},
				{Step: 2, Name: "Refinement", Persona: "great", LLM: "claude", ContextMode: "previous",
					PromptPrefix: "You are the refinement reviewer. Review and refine this analysis. Identify gaps, add architectural considerations, consider maintainability and elegance."},
				{Step: 3, Name: "Chaos Check", Persona: "yolo", LLM: "gemini", ContextMode: "previous",
					PromptPrefix: "You are the stress tester. Challenge and stress-test this plan. What could go wrong? What unconventional approaches exist? No sacred cows."},
				{Step: 4, Name: "Consolidation", Persona: "mighty", LLM: "codex", ContextMode: "all_previous",
					PromptPrefix: "You are the lead synthesizer. Synthesize all inputs into a master plan. Merge insights, resolve contradictions, create unified vision, prioritize features."},
				{Step: 5, Name: "Architect Breakdown", Persona: "architect", LLM: "claude", ContextMode: "consolidation",
					PromptPrefix: "You are the task architect. Break down this master plan into actionable implementation tasks. Format each task with an h3 header exactly like this: ### TASK 1: Title Here. Include for each task: **Description:** what needs to be done, **Acceptance Criteria:** bullet list, **Complexity:** S|M|L|XL, **Dependencies:** Task N references or None, **Files:** list of files likely touched. Think in systems. Consider failure modes. Make it maintainable."},
			},
			"expanded": {
				{Step: 1, Name: "Initial Analysis", Persona: "mighty", LLM: "codex", ContextMode: "prompt_only",
					PromptPrefix: "You are the lead analyst. Analyze this idea and provide comprehensive initial analysis including: problem space, technical requirements, proposed architecture, unknowns and risks."},
				{Step: 2, Name: "Refinement", Persona: "great", LLM: "claude", ContextMode: "previous",
					PromptPrefix: "You are the refinement reviewer. Review and refine this analysis. Identify gaps, add architectural considerations, consider maintainability and elegance."},
				{Step: 3, Name: "Reality Check", Persona: "mediocre", LLM: "claude", ContextMode: "previous",
					PromptPrefix: "You are the pragmatic reviewer. Reality-check this plan. What's over-engineered? What can be simplified? Where is cleverness hiding complexity? Cut the fluff, keep what ships."},
				{Step: 4, Name: "Chaos Check", Persona: "yolo", LLM: "gemini", ContextMode: "previous",
					PromptPrefix: "You are the stress tester. Challenge and stress-test this plan. What could go wrong? What unconventional approaches exist? No sacred cows."},
				{Step: 5, Name: "MVP Lens", Persona: "startup", LLM: "claude", ContextMode: "previous",
					PromptPrefix: "You are the scope cutter. What's the MVP here? Cut scope ruthlessly. What can ship in week one? What's a nice-to-have disguised as a must-have? If you're not embarrassed by v1, you shipped too late."},
				{Step: 6, Name: "Consolidation", Persona: "mighty", LLM: "codex", ContextMode: "all_previous",
					PromptPrefix: "You are the lead synthesizer. Synthesize all inputs into a master plan. Merge insights, resolve contradictions, create unified vision, prioritize features."},
				{Step: 7, Name: "Architect Breakdown", Persona: "architect", LLM: "claude", ContextMode: "consolidation",
					PromptPrefix: "You are the task architect. Break down this master plan into actionable implementation tasks. Format each task with an h3 header exactly like this: ### TASK 1: Title Here. Include for each task: **Description:** what needs to be done, **Acceptance Criteria:** bullet list, **Complexity:** S|M|L|XL, **Dependencies:** Task N references or None, **Files:** list of files likely touched. Think in systems. Consider failure modes. Make it maintainable."},
			},
		},
	}
}
