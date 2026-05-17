package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestActiveSurfaceHasNoStaleRepoIdentityLeaks(t *testing.T) {
	root := repoRootForTests(t)

	targets := []string{
		"README.md",
		"config.default.json",
		"docs/architecture.md",
		"docs/agent_persona_system.md",
		"docs/token_strategy.md",
		"docs/user_workflow.md",
		"docs/workflow_output_contracts.md",
		"go/cmd/legal-bot",
		"go/internal/config",
		"go/internal/llm",
		"go/internal/persona",
		"scripts",
	}

	disallowed := []string{
		"Historian Vern",
		"HISTORIAN VERN",
		"MightyVern",
		"Codex Vern",
		"YOLO Vern",
		"Vernile the Great",
		"Architect Vern",
		"Startup Vern",
		"Vern the Mediocre",
		"/Users/justin/projects/jdonohoo/legal-bot",
		"go/cmd/legal-bot/vern",
	}

	for _, rel := range targets {
		full := filepath.Join(root, filepath.FromSlash(rel))
		info, err := os.Stat(full)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			t.Fatalf("stat %s: %v", rel, err)
		}

		if info.IsDir() {
			err = filepath.WalkDir(full, func(path string, d os.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}
				if d.IsDir() {
					return nil
				}
				if filepath.Base(path) == "agents_generated.go" || filepath.Base(path) == "repo_identity_test.go" {
					return nil
				}
				checkFileForDisallowedIdentity(t, path, root, disallowed)
				return nil
			})
			if err != nil {
				t.Fatalf("walk %s: %v", rel, err)
			}
			continue
		}

		checkFileForDisallowedIdentity(t, full, root, disallowed)
	}
}

func checkFileForDisallowedIdentity(t *testing.T, path string, root string, disallowed []string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	content := string(data)
	for _, needle := range disallowed {
		if strings.Contains(content, needle) {
			rel, relErr := filepath.Rel(root, path)
			if relErr != nil {
				rel = path
			}
			t.Fatalf("%s still contains stale identity reference %q", filepath.ToSlash(rel), needle)
		}
	}
}
