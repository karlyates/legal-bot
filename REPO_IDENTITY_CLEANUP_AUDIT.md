# Repo Identity Cleanup Audit

## Summary

- Scope: targeted active-surface repo identity cleanup only; no broad legacy-package rewrite.
- Active code fixed now: 5 files.
- Active user-facing/help text fixed now: 2 command files.
- Tests updated or added: 4 files.
- Generated assets regenerated: no.
- Legacy references intentionally preserved: legacy packages, legacy wrapper inventory, legacy migration docs, and compatibility env/config keys.

## What Changed

- Updated active historian help/output branding in `go/cmd/legal-bot/historian.go` from `Historian Vern` to `Legal-Bot Historian`.
- Updated active TUI help text in `go/cmd/legal-bot/tui.go` to describe the legacy utility surface without exposing `VernHole` by name.
- Updated an active binary-path comment in `go/cmd/legal-bot/intake.go` from the old `vern` path to `legal-bot`.
- Replaced stale Vern-branded hardcoded fallback prompt prefixes in `go/internal/config/config.go` with neutral role-based wording.
- Replaced stale example strings/comments in `go/internal/persona/loader.go`.
- Updated stale active test fixtures in `go/internal/config/config_test.go`, `go/internal/persona/loader_test.go`, and `go/internal/llm/runner_test.go`.
- Added `go/cmd/legal-bot/repo_identity_test.go` to guard active surfaces against a small allowlist-based set of stale repo-identity leaks.

## What Did Not Change

- `go/go.mod` module path was not changed.
- Current imports using `github.com/jdonohoo/legal-bot/go/...` were left unchanged.
- Legacy packages under `go/internal/pipeline`, `go/internal/council`, `go/internal/generate`, `go/internal/vts`, `go/internal/tobeads`, and `go/internal/tui` were not broadly refactored.
- Legacy wrapper inventory in `go/cmd/legal-bot/routes_catalog.go` and the legacy docs were preserved.
- `config.default.json` was not changed because the targeted stale active branding from this pass was not present there.
- Embedded assets were not regenerated because no agent files or embedded config sources changed.

## Active References Cleaned

- `Historian Vern`
- `HISTORIAN VERN`
- `VernHole` in active TUI help text
- `go/cmd/legal-bot/vern` comment reference
- Hardcoded fallback prompt strings:
  - `MightyVern`
  - `Vernile the Great`
  - `YOLO Vern`
  - `Architect Vern`
  - `Startup Vern`
  - `Vern the Mediocre`
- Test/example strings:
  - `MightyVern / Codex Vern`
  - `YOLO Vern`
  - `/Users/justin/projects/jdonohoo/legal-bot`

## Legacy References Intentionally Preserved

- Legacy surface inventory and explanations in:
  - `docs/legacy_surfaces.md`
  - `docs/legacy_discovery_migration.md`
  - `go/cmd/legal-bot/routes_catalog.go`
- Legacy packages and their APIs:
  - `go/internal/pipeline`
  - `go/internal/council`
  - `go/internal/generate`
  - `go/internal/vts`
  - `go/internal/tobeads`
  - `go/internal/tui`
- Compatibility wrappers in `bin/`.
- Compatibility env/config identities still in use:
  - `VERN_TIMEOUT`
  - `VERN_WORKING_DIR`
  - `VERN_ROOT`
  - `VERN_LOG`
  - `discovery_pipelines`
  - `discovery_pipeline`
  - `vernhole`
  - `default_discovery_path`

## Module Path Decision

- Left unchanged.
- Current module/import identity remains `github.com/jdonohoo/legal-bot/go`.
- No concrete build or import bug required a module-path rewrite in this pass.

## Tests / Validation

- Ran `gofmt` on all changed Go files.
- Ran full Go test suite from `C:\GIT\legal-bot\go` with repo-local Go cache env vars.
- Ran non-LLM CLI smoke checks:
  - `go run ./cmd/legal-bot routes`
  - `go run ./cmd/legal-bot routes --json`
  - `go run ./cmd/legal-bot guide --list`
- Result: all validations passed.

## Remaining Known References

- Legacy route/catalog text still mentions `Vern`, `VernHole`, `oracle`, `council`, `discovery`, and `VTS` where the CLI is explicitly documenting legacy surfaces.
- Compatibility env/config names still contain `VERN`/`discovery`/`vernhole` because changing those would be a separate compatibility migration.
- Current Go imports still use `github.com/jdonohoo/legal-bot/go`, which is intentionally unchanged in this pass.

## Follow-Up Recommendations

- Do a separate compatibility-focused config/env cleanup pass if you want to migrate legacy keys like `vernhole`, `discovery_pipeline`, and `VERN_*` without breaking existing setups.
- Do a separate module-path decision pass if you want to rename `github.com/jdonohoo/legal-bot/go`; treat that as an import/build migration, not a branding cleanup.
- If the legacy TUI is kept long-term, consider a later pass that re-frames its internal labels more consistently while preserving compatibility behavior.
