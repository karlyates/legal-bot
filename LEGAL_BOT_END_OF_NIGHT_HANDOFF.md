# Legal-Bot End-of-Night Handoff

## Executive Summary

- Current repo state: heavily dirty working tree with broad uncommitted changes across legacy wrappers, agent files, docs, embedded assets, and Go control-plane code. This is not a clean checkout.
- What changed since v3: workflow output contracts are now present and tested, the final synthesizer is explicitly contract-driven, and the active surface has been cleaned of obvious stale Vern-era branding.
- Workflow output contracts implemented: Verified.
- Repo identity cleanup implemented: Partially verified. Active surfaces are cleaned; legacy compatibility/history references remain intentionally.
- Tests passing: Verified.
- Safe to continue tomorrow: Yes, but the next step should be practical validation of real workflow output quality rather than more naming cleanup.
- Recommended next step: run a controlled dry-run or non-sensitive sample matter through `guide`/`workflow` and inspect one real final packet.

## Current Baseline

- Active command surface: `guide`, `routes`, `intake`, `review`, `workflow`, `run`, `setup`, `historian`, `tui`.
- Route catalog: centralized in `go/cmd/legal-bot/routes_catalog.go`; it still owns review modes, workflow types, conditional specialists, and legacy surface inventory.
- Model routing: still split into engine / model / effort; `openai:*` normalization remains backward compatible; the build/test pass confirms the routing layer is healthy.
- Guide: deterministic front door remains in place with `--list`, `--dry-run`, `--choice`, and `--intent`.
- Workflow contracts: implemented in `go/cmd/legal-bot/workflow_contracts.go`, injected through `workflow.go`, and documented in `docs/workflow_output_contracts.md`.
- Legacy surfaces: still intentionally present in the route catalog, `docs/legacy_surfaces.md`, `docs/legacy_discovery_migration.md`, and the legacy packages under `go/internal/{pipeline,council,generate,vts,tobeads,tui}`.
- Repo identity cleanup: active user-facing strings are now Legal-Bot-branded; compatibility names and legacy docs remain by design.
- Tests/build status: full Go suite passes, the CLI builds, and the non-LLM smoke checks passed.

## Changes Since Pipeline Architecture Report v3

| Area | v3 Baseline | Current State | Change | Confidence |
|---|---|---|---|---|
| Command surface | Active `guide`, `routes`, `intake`, `review`, `workflow`, utility commands | Same active surface remains | No major change | Verified |
| Route catalog | Central route metadata and legacy inventory | Still centralized; diagnostics output is live | No major change | Verified |
| Model routing | Engine/model/effort split with backward-compatible `openai:*` normalization | Same, and tests still pass | No major change | Verified |
| Guide | Deterministic guided front door exists | Still present; `guide --list` works | No major change | Verified |
| Workflow output contracts | Not part of the v3 baseline summary | Fully implemented and documented | Major addition | Verified |
| Managing partner synthesizer | Final synthesis existed, but not yet explicitly contract-driven | Now a decision-maker that follows workflow contracts, source discipline, and Do Not Chase guidance | Strengthened | Verified |
| Shared source / strategy rules | Not highlighted in the v3 summary | Present in prompt corpus and aligned to contract-driven final synthesis | Strengthened | Inferred from code/docs |
| Docs | Active docs described workflow and legacy surfaces | `docs/workflow_output_contracts.md` now anchors the final-packet contract layer; active docs reflect it | Expanded | Verified |
| Tests | General suite passed locally | General suite passes; new contract and repo-identity tests are present | Stronger | Verified |
| Embedded assets | Baseline contained generated prompt/config assets | `go/internal/embedded/agents_generated.go` is modified, indicating embedded prompt assets moved with the agent corpus | Likely regenerated | Inferred from diff |
| Repo identity cleanup | Not yet captured in v3 baseline | Active user-facing stale Vern branding is cleaned; audit/test artifacts exist | New completion artifact | Partially verified |
| Legacy surfaces | Kept as compatibility/history | Still intact and explicitly documented | No major change | Verified |
| Module path / imports | Module path left alone | Still `github.com/jdonohoo/legal-bot/go`; imports remain consistent | No change | Verified |
| Scripts / bin | Legacy wrappers existed | Legacy wrappers are still present and documented; active wrappers remain compatibility-oriented | No major change | Verified |
| Build / test status | Last known state was green | Full test suite and build passed during this inspection | Verified and current | Verified |

## Git State

- `git status --short`: the working tree is heavily dirty. There are broad deletions of old Vern-era assets, many new legal agent files, modifications in `go/internal/*`, docs, wrappers, and generated prompt assets.
- `git log --oneline -10`: latest commit is `dc6889d Capture LLM subprocess stderr, harden tobeads CLI, bump to v2.9.1`.
- `git diff --stat`: `326 files changed, 4112 insertions(+), 30962 deletions(-)`.
- `git diff --name-only`: broad changes touch active command/control-plane files, internals, docs, agents, bin wrappers, scripts, and generated assets.

Changed-file categories:

- Go command/control-plane: `go/cmd/legal-bot/*`, especially `workflow.go`, `workflow_contracts.go`, `routes_catalog.go`, and the CLI command files.
- Go internals: `go/internal/config`, `go/internal/llm`, `go/internal/persona`, `go/internal/embedded`, and the legacy helper packages.
- Agents/prompts: `agents/*.md` plus `agents/shared/*`.
- Docs: `README.md`, `docs/*.md`, including `docs/workflow_output_contracts.md`.
- Tests: `*_test.go` across command, config, llm, persona, embedded, and legacy packages.
- Scripts/bin: `bin/*`, `scripts/*`, `install-hooks.sh`.
- Generated assets: `go/internal/embedded/agents_generated.go`.
- Audit/report files: `REPO_IDENTITY_CLEANUP_AUDIT.md`, `LEGAL_BOT_END_OF_NIGHT_HANDOFF.md`, `ARCHITECTURE_DELTA_REPORT-v4.md`.

## Test and Build Results

- `cd go; & "C:\Program Files\Go\bin\go.exe" test ./...` -> Verified, exit 0, all packages passed.
- `cd go; & "C:\Program Files\Go\bin\go.exe" build -o bin/legal-bot.exe ./cmd/legal-bot` -> Verified, exit 0.

Notes:

- The first sandboxed test/build attempts hit `Access is denied` when Go tried to create repo-local cache/temp directories. I reran them outside the sandbox with approval and they passed.
- No real LLM calls were made.

## Non-LLM Smoke Check Results

- `cd go; & "C:\Program Files\Go\bin\go.exe" run ./cmd/legal-bot routes` -> Verified. It printed the active commands, review modes, review document types, workflow types, conditional specialists, and legacy surfaces.
- `cd go; & "C:\Program Files\Go\bin\go.exe" run ./cmd/legal-bot routes --json` -> Verified. It emitted the JSON route snapshot.
- `cd go; & "C:\Program Files\Go\bin\go.exe" run ./cmd/legal-bot guide --list` -> Verified. It printed the guided choices list.

## Workflow Output Contracts

- Location: `go/cmd/legal-bot/workflow_contracts.go`.
- Documentation: `docs/workflow_output_contracts.md`.
- Prompt injection: `workflow.go` builds a concise contract summary for every workflow step and injects the full contract into the final `managing-partner-final-synthesizer` step.
- Final synthesizer: receives the full contract and is instructed to choose, prioritize, recommend, and name what not to chase.
- Earlier steps: receive the workflow metadata and concise contract summary, which keeps the whole panel aligned.
- Test coverage: `go/cmd/legal-bot/workflow_contracts_test.go` covers contract existence, headings, global rules, final synthesizer prompt mentions, and doc coverage.
- Gaps: no obvious structural gap in coverage; the next likely improvement is practical output calibration using a real controlled matter.

| Workflow | Contract Exists? | Required Sections Present? | Do Not Chase? | Source Strength? | External-Safe Rules? | Notes |
|---|---:|---:|---:|---:|---:|---|
| triage | Yes | Yes | Yes | Yes | Yes | Includes Bottom Line, Source Strength, Options Matrix, Do Not Chase Check |
| counsel-brief | Yes | Yes | Yes | Yes | Yes | Focuses on counsel-facing synthesis and attorney questions |
| sm-triage | Yes | Yes | Yes | Yes | Yes | Includes SM scope, ripeness, requested directive, and Do Not Chase |
| draft-message | Yes | Yes | Yes | Yes | Yes | Explicitly separates safe facts from risky claims |
| evidence-packet | Yes | Yes | Yes | Yes | Yes | Focuses on evidence inventory, missing evidence, and packet structure |
| pattern-review | Yes | Yes | Yes | Yes | Yes | Includes pattern theory, strongest examples, and Do Not Chase |
| fact-lock | Yes | Yes | Yes | Yes | Yes | Strong separation between locked facts, inference, and authority gaps |
| authority-check | Yes | Yes | Yes | Yes | Yes | Explicitly prevents overclaiming currentness or authority status |
| prep | Yes | Yes | Yes | Yes | Yes | Centers decisions needed, pushback, and what not to say |
| journal-entry | Yes | Yes | Yes | Yes | Yes | Neutral record-oriented packet with caveats and evidence anchors |

## Final Synthesizer / Agent Prompt Changes

- Managing partner: `agents/managing-partner-final-synthesizer.md` now says the agent is a decision-maker, not a summarizer, and it explicitly requires workflow contracts, source-supported facts, external-safe separation, and Do Not Chase guidance.
- Shared source/strategy rules: the prompt corpus includes shared safety and source-hierarchy rules, and the final synthesizer explicitly references them.
- Legal authority scholar: present in the route/conditional system and still focused on legal authority gaps and currentness discipline.
- Practical resolution reviewer: still emphasizes administrable relief, proportionality, and Do Not Chase-style restraint where appropriate.
- Embedded assets: `go/internal/embedded/agents_generated.go` is modified, so the embedded prompt bundle reflects the current prompt corpus.

## Repo Identity Cleanup

- `REPO_IDENTITY_CLEANUP_AUDIT.md` exists and documents the earlier active-surface cleanup pass.
- Active references cleaned: `Historian Vern`, `HISTORIAN VERN`, `VernHole` in active TUI help text, stale `go/cmd/legal-bot/vern` comment text, stale Vern-branded fallback prompt prefixes, and stale active test/example strings.
- Legacy references intentionally left: route catalog legacy inventory, `docs/legacy_surfaces.md`, `docs/legacy_discovery_migration.md`, legacy packages, and compatibility wrappers.
- `go/go.mod`: module path was left unchanged.
- Imports: still consistent with `github.com/jdonohoo/legal-bot/go/...`.
- Stale active reference test: `go/cmd/legal-bot/repo_identity_test.go` exists and is aimed only at active surfaces.
- Remaining suspicious references: mostly compatibility env vars and comments (`VERN_LOG`, `VERN_TIMEOUT`, `VERN_ROOT`, `VERN_WORKING_DIR`) plus legacy docs/wrappers, which are intentional.

| Reference Category | Count / Files | Status | Follow-Up Needed |
|---|---|---|---|
| Active user-facing cleaned | Several command/help strings and config/test fixtures | Cleaned | Only if new active branding leaks reappear |
| Active code remaining | `go/internal/config/config.go`, `go/internal/llm/log.go`, `go/cmd/legal-bot/run.go` | Compatibility references remain intentionally | Maybe, in a separate compatibility migration |
| Legacy intentional | `go/internal/{pipeline,council,generate,vts,tobeads,tui}`, `docs/legacy_*`, `bin/*` | Preserved | No immediate action |
| Module path / import identity | `go/go.mod`, local imports | Unchanged | Defer module-path cleanup |
| Generated assets | `go/internal/embedded/agents_generated.go` | Present / current | Regenerate again only if prompt sources change |
| Docs historical | `README.md`, `docs/legacy_*`, legacy sections of active docs | Intentional | No immediate action |
| Suspicious remaining | `VERN_*` env fallbacks, legacy comments, some wrapper labels | Intentional compatibility | Only if you decide to retire compatibility support |

## Current Active Command Contract

- `guide`: active front door and deterministic router.
- `routes`: active route/legacy inventory diagnostic.
- `intake`: active matter ingestion and knowledge-building.
- `review`: active draft review and authority-aware review.
- `workflow`: active draftless legal workflow engine.
- `run`: utility surface for one-shot LLM subprocesses.
- `setup`: utility surface for first-run configuration.
- `historian`: utility surface for directory indexing.
- `tui`: utility/legacy-oriented interactive surface.

Explicitly non-current legacy commands:

- `discovery`
- `council`
- `oracle`
- `generate`

## Current Known Risks

### Technical Risks

- The working tree is very dirty and spans many files that were not changed in this inspection pass.
- Compatibility layer code still carries `VERN_*`, `discovery_*`, and `vernhole` names.

### UX Risks

- The active workflow is much more polished now, but the utility surfaces still expose legacy concepts in their compatibility wrappers and docs.

### Prompt / Output Quality Risks

- The contracts exist, but they still need real-world calibration on controlled sample matters to verify output quality and reduce verbosity or over-hedging.

### Repo Maintenance Risks

- Module-path cleanup is still deferred.
- Legacy packages are intentionally preserved, so the repo will continue to look mixed until a separate compatibility-migration pass happens.

### Legal-Use / Source-Discipline Risks

- The contracts and final synthesizer reduce risk, but the tool still depends on careful source discipline and human verification before anything external is used.

## Recommended Next Steps Tomorrow

1. Run one controlled sample matter through `guide --dry-run` and then a real non-sensitive workflow.
   - Why it matters: validates the new contracts against a real packet.
   - Risk: you may discover the packet wants more or less structure than the contract currently demands.
   - Suggested model: use the current default routing; if you force one engine, start with `claude` for final synthesis.
   - Type: manual test, then light iteration.
2. Inspect one actual `workflow` output packet end to end.
   - Why it matters: tells us whether the contract headings are useful or just theoretically correct.
   - Risk: output may be too verbose or too generic.
   - Suggested model: current default routing.
   - Type: manual test.
3. Calibrate the contract/doc text after seeing real output.
   - Why it matters: the best contract language is the one that matches actual packets.
   - Risk: overfitting to one matter.
   - Suggested model: not applicable unless you re-run a packet.
   - Type: implementation if needed, otherwise audit.
4. Defer module-path cleanup.
   - Why it matters: it is a separate migration with import and build consequences.
   - Risk: touching it too early can create noisy breakage for little operator value.
   - Suggested model: not applicable.
   - Type: defer.

## Do Not Chase Tomorrow

- Do not rename the Go module path unless a concrete import/build bug appears.
- Do not delete all legacy packages in one sweep.
- Do not rewrite the TUI just to erase old names.
- Do not add a direct OpenAI API path.
- Do not expand persona proliferation without a concrete workflow need.
- Do not change workflow routing again until you have seen a real packet.

## Commands Karl Should Run Tomorrow

```powershell
cd C:\GIT\legal-bot\go
$env:GOCACHE='C:\GIT\legal-bot\.go-cache\gocache'
$env:GOMODCACHE='C:\GIT\legal-bot\.go-cache\gomodcache'
$env:GOTMPDIR='C:\GIT\legal-bot\.go-cache\tmp'
New-Item -ItemType Directory -Force $env:GOCACHE | Out-Null
New-Item -ItemType Directory -Force $env:GOMODCACHE | Out-Null
New-Item -ItemType Directory -Force $env:GOTMPDIR | Out-Null
& "C:\Program Files\Go\bin\go.exe" test ./...
& "C:\Program Files\Go\bin\go.exe" build -o bin/legal-bot.exe ./cmd/legal-bot
& "C:\Program Files\Go\bin\go.exe" run ./cmd/legal-bot routes
& "C:\Program Files\Go\bin\go.exe" run ./cmd/legal-bot routes --json
& "C:\Program Files\Go\bin\go.exe" run ./cmd/legal-bot guide --list
& "C:\Program Files\Go\bin\go.exe" run ./cmd/legal-bot guide Revere-042926 --dry-run
```

Optional safe setup for a tiny sample matter:

```powershell
cd C:\GIT\legal-bot
New-Item -ItemType Directory -Force matters\sample-matter\input, matters\sample-matter\drafts | Out-Null
Set-Content matters\sample-matter\input\source-1.md "# Sample source"
```

## ChatGPT Continuation Context

Paste this to resume tomorrow if needed:

The repo is now in the v3 architecture plus two additional passes state. Workflow output contracts are implemented in `go/cmd/legal-bot/workflow_contracts.go`, documented in `docs/workflow_output_contracts.md`, injected in `go/cmd/legal-bot/workflow.go`, and tested in `go/cmd/legal-bot/workflow_contracts_test.go`. The managing-partner synthesizer is explicitly decision-maker oriented and contract-driven. Repo-identity cleanup already removed obvious active Vern-era branding from user-facing surfaces and added `go/cmd/legal-bot/repo_identity_test.go`, while legacy compatibility surfaces remain intentionally documented. Full Go tests, build, and non-LLM smoke checks all passed during the current inspection. The biggest open question is not code correctness but whether the new workflow contracts produce a useful real packet on a controlled sample matter. Next file to inspect for implementation follow-up: `go/cmd/legal-bot/workflow_contracts.go`, then `go/cmd/legal-bot/workflow.go`, then `agents/managing-partner-final-synthesizer.md`.
