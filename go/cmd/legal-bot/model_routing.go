package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jdonohoo/legal-bot/go/internal/config"
	"github.com/jdonohoo/legal-bot/go/internal/llm"
)

type routeSelection struct {
	Target   config.ResolvedModelTarget
	Warnings []string
}

func selectStepRoute(cfg *config.Config, profileName string, fallbackLLM string) (*routeSelection, error) {
	if override, warnings, err := cfg.GetOverrideTarget(); err != nil {
		return nil, err
	} else if override != nil {
		return chooseRoute([]config.ResolvedModelTarget{*override}, warnings)
	}

	targets, warnings := cfg.ResolveModelProfileTargets(profileName, fallbackLLM)
	return chooseRoute(targets, warnings)
}

func chooseRoute(targets []config.ResolvedModelTarget, warnings []string) (*routeSelection, error) {
	target, selectWarnings, err := config.SelectRunnableTarget(targets)
	allWarnings := append([]string{}, warnings...)
	allWarnings = append(allWarnings, selectWarnings...)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &routeSelection{Target: target, Warnings: allWarnings}, nil
}

func describeRoute(profileName string, selection *routeSelection) string {
	source := "step"
	if profileName != "" {
		source = profileName
	} else if selection.Target.SourceProfile != "" {
		source = selection.Target.SourceProfile
	}
	return fmt.Sprintf("%s -> %s", source, selection.Target.DebugString())
}

func printRouteWarnings(warnings []string) {
	for _, warning := range warnings {
		if strings.TrimSpace(warning) == "" {
			continue
		}
		fmt.Printf("    Route note: %s\n", warning)
	}
}

func runOptionsForRoute(selection *routeSelection, prompt string, outputFile string, persona string, timeoutSeconds int, agentsDir string) llm.RunOptions {
	return llm.RunOptions{
		LLM:        selection.Target.Engine,
		Model:      selection.Target.Model,
		Effort:     selection.Target.Effort,
		Prompt:     prompt,
		OutputFile: outputFile,
		Persona:    persona,
		Timeout:    llmDuration(timeoutSeconds),
		AgentsDir:  agentsDir,
	}
}

func llmDuration(timeoutSeconds int) time.Duration {
	return time.Duration(timeoutSeconds) * time.Second
}
