# Pipeline Architecture Report v3

## Executive Summary

- The current active Legal-Bot architecture is now a Go CLI control plane around `intake`, `review`, `workflow`, `guide`, and `routes`, with `run`, `setup`, `historian`, and `tui` still exposed as utility/compatibility commands.
- The biggest v2 drift has been reduced: route names, review modes, workflow types, conditional specialists, and legacy wrappers are now centralized in `go/cmd/legal-bot/routes_catalog.go` and exposed through `legal-bot routes`.
- Model routing is materially improved. `config.default.json` now uses explicit `model_profiles` plus `llm_modes`, `go/internal/config/model_routing.go` parses compact specs and normalizes legacy `openai:*` targets to Codex, and `go/internal/llm/runner.go` supports `claude`, `codex`, `gemini`, and `copilot` with model/effort propagation where supported.
- `legal-bot guide` exists and is wired to the existing `intake`, `review`, and `workflow` execution paths. It is deterministic and does not use an LLM to choose a route.
- The guide is not perfectly side-effect free: workflow plans still write a guided trace file under `matters/<matter>/working/` even when `--dry-run` is used.
- Legacy Vern-era surfaces are de-risked much better than v2. The old discovery/council/oracle/generate wrappers are now explicit legacy stubs or documented compatibility surfaces instead of looking like normal operator paths.
- The biggest remaining operator risks are configuration drift in `config.default.json`, the fact that route catalog centralization does not yet replace code-driven prompt assembly, and the environment sensitivity of LLM CLI availability/cache behavior.
- Karl should first treat `legal-bot guide`, `legal-bot routes`, `legal-bot intake`, `legal-bot review`, and `legal-bot workflow` as the current legal workflow contract, and treat Vern-era discovery/council/oracle/generate artifacts as historical or compatibility-only.

## v2 to v3 Change Summary

| Area | v2 Status | v3 Status | Net Change | Confidence |
|---|---|---|---|---|
| Active command surface | `intake`, `review`, `workflow`, `run`, `setup`, `historian`, `tui`; no guided front door or route diagnostics | Adds `guide` and `routes` to the active contract; `review-draft` remains an alias | More operator-facing clarity | High |
| Intake | Active, declarative pipeline in config | Still active, unchanged in behavior | Stable | High |
| Review | Active, code/config mixed route selection | Still active, but route names and document types are now cataloged centrally | Less drift risk | High |
| Workflow | Active, code/config mixed route selection | Still active, but workflow types now live in a shared catalog | Less drift risk | High |
| Guide | Absent | Active deterministic guided front door with `--list`, `--dry-run`, `--choice`, `--intent` | New operator UX layer | High |
| Route catalog / control plane | Implicit and spread across files | Explicit route catalog in `routes_catalog.go` plus `legal-bot routes` diagnostic output | Major improvement | High |
| Conditional routing | Heuristic, hard to inspect | Named conditional rules and keywords are centralized and tested | Much clearer | High |
| Model profile resolution | `openai:*` mismatch collapsed into Claude-like behavior; engine/model/effort conflated | Engine/model/effort separation exists; `openai:*` normalized to Codex; compact specs supported | Core mismatch fixed | High |
| LLM runner support | Engine-only in practice, with legacy fallback behavior | `claude`, `codex`, `gemini`, `copilot` all supported with model/effort propagation where relevant | Broader and cleaner | High |
| Agent/persona loading | Legacy aliases existed but were easy to lose track of | Active agent inventory plus alias resolution and tests | Safer, clearer | High |
| Embedded assets | Generated assets existed but were not clearly tied to validation | Generator and embedded tests are explicit; active agent/config inventory is validated | Better guardrails | Medium |
| Case guidance loading | Global prompt context, but operator docs were thinner | Behavior is documented more clearly in docs | Better operator understanding | High |
| Matter layout | Ad hoc matter folder expectations | Guide now detects/creates standard matter layout | Better onboarding | High |
| Output behavior | Static output paths by pipeline type | Same output behavior remains | No behavioral drift | High |
| Windows scripts | Intake/review wrappers only, guide absent | `run-guide.bat` added; intake/review wrappers remain | More direct UX | High |
| `bin/` wrappers | Vern-era wrappers looked runnable | Discovery/council/oracle/generate are now legacy stubs or clearly documented legacy paths | Major risk reduction | High |
| Legacy Vern surfaces | Easy to confuse with current workflow | Clearly separated into docs and wrapper stubs | Major risk reduction | High |
| Docs | Mixed, with legacy and active surfaces less clearly separated | README/user workflow/architecture/legacy inventory are aligned to the current contract | Better operator truth | High |
| Tests | Some route/config tests existed | Route catalog, guide, model routing, aliases, and legacy inventory are tested more directly | Stronger regression coverage | High |

## Active Command Surface

| Command | Active? | Purpose | Key Flags | Calls LLM? | Notes |
|---|---:|---|---|---:|---|
| `guide` | Yes | Guided front door to `intake`, `review`, or `workflow` | `--list`, `--dry-run`, `--choice`, `--intent`, `--audience`, `--urgency`, `--situation`, `--goal`, `--questions`, `--draft`, `--mode`, `--document-type`, `--child-related`, `--financial` | Indirect | Alias `start` exists; deterministic route selection |
| `intake` | Yes | Build the matter knowledge layer from source files | Matter arg only plus existing hidden/config-driven behavior | Yes | Requires source files in `matters/<matter>/input/` |
| `review` | Yes | Review a draft against the knowledge layer | `--mode`, `--document-type`, `--draft`, `--single-llm`, `--llm-mode`, `--child-related`, `--financial` | Yes | Alias `review-draft` remains |
| `workflow` | Yes | Run a user-job workflow such as triage, counsel brief, or authority check | `--type`, `--audience`, `--urgency`, `--situation`, `--goal`, `--questions`, `--draft`, `--single-llm`, `--llm-mode`, `--child-related`, `--financial` | Yes | Shared route catalog |
| `routes` | Yes | Show the active route catalog and legacy surface inventory | `--json` | No | Diagnostic/control-plane command |
| `run` | Yes | Run a single LLM subprocess | `-o/--output`, `-p/--persona`, `-t/--timeout` | Yes | Low-level utility command |
| `setup` | Yes | Detect available LLMs and save config | No route flags | No | Utility command; updates `LLMs`/`LLMMode` |
| `historian` | Yes | Index a directory into `input-history.md` | `--llm`, `--timeout` | Yes | Uses legacy pipeline package for utility history work |
| `tui` | Yes | Launch the Bubble Tea utility UI | None on the root command | Usually yes | Legacy-oriented utility surface, not the main legal front door |
| `doctor` | No | N/A | N/A | No | No active root command with this name |
| `doctor-routes` | No | N/A | N/A | No | No active root command with this name |
| `models` | No | N/A | N/A | No | No active root command with this name |
| `discovery` | No | N/A | N/A | No | Legacy Vern surface, not a current root command |
| `council` | No | N/A | N/A | No | Legacy Vern surface, not a current root command |
| `oracle` | No | N/A | N/A | No | Legacy Vern surface, not a current root command |
| `generate` | No | N/A | N/A | No | Legacy Vern surface, not a current root command |

## Active Pipeline Definition Locations

| Area | File/Path | Role | Active? | Changed Since v2? | Notes |
|---|---|---|---:|---:|---|
| Root CLI registration | `go/cmd/legal-bot/main.go` | Cobra root command | Yes | Yes | Root surface is still the active legal control plane |
| Intake command | `go/cmd/legal-bot/intake.go` | Intake execution | Yes | No | Still reads case context and matter source files |
| Review command | `go/cmd/legal-bot/review.go` | Draft review execution | Yes | Yes | Uses shared review route helpers and conditional routing |
| Workflow command | `go/cmd/legal-bot/workflow.go` | Workflow execution | Yes | Yes | Uses shared workflow catalog and conditional routing |
| Guide command | `go/cmd/legal-bot/guide.go` | Guided UX front door | Yes | Yes | Deterministic menu/plan/execution wrapper around existing commands |
| Routes diagnostic | `go/cmd/legal-bot/routes_command.go` | Prints route inventory | Yes | Yes | Supports `--json` |
| Route catalog | `go/cmd/legal-bot/routes_catalog.go` | Central route/alias/legacy inventory | Yes | Yes | Main control-plane catalog |
| Config loader | `go/internal/config/config.go` | Config loading/defaults | Yes | Yes | Adds `LoadProjectDefault` and supports model-profile resolution |
| Model routing | `go/internal/config/model_routing.go` | Model target parsing/selection | Yes | Yes | Fixes engine/model/effort separation and legacy `openai:*` aliasing |
| LLM runner | `go/internal/llm/runner.go` | Subprocess execution | Yes | Yes | Dispatches to local CLIs and logs run metadata |
| Persona loader | `go/internal/persona/loader.go` | Persona parsing and legacy alias resolution | Yes | Yes | Keeps legacy alias compatibility while naming active personas explicitly |
| Embedded generator | `go/internal/embedded/gen.go` | Generates embedded agent/config assets | Yes | No | Regenerate with `cd go && go generate ./internal/embedded/` when agents/config change |
| Embedded output | `go/internal/embedded/agents_generated.go` | Compiled-in personas/config | Yes | Yes | Source of truth for binary-only fallback |
| Active agents | `agents/*.md` | Persona bodies | Yes | Yes | 17 active legal personas plus shared rule docs in `agents/shared/` |
| Case context | `case/*.md` | Global prompt context | Yes | No | Loaded as prompt context, not parsed as rules |
| Scripts | `scripts/*.bat` | Windows wrappers | Yes | Yes | Guide wrapper added; no workflow wrapper yet |
| Legacy wrappers | `bin/*` | Compatibility / legacy stubs | Partial/legacy | Yes | Discovery/council/oracle/generate are now visibly legacy |
| Legacy packages | `go/internal/pipeline`, `go/internal/council`, `go/internal/generate`, `go/internal/vts`, `go/internal/tobeads`, `go/internal/tui` | Compatibility or historical code | Partial/legacy | Mostly no | Still present; not the main legal workflow contract |

## Route Catalog / Control Plane

The route catalog now lives in `go/cmd/legal-bot/routes_catalog.go` and is reused by `review.go`, `workflow.go`, `guide.go`, and `routes_command.go`.

| Route / Concept | Defined In | Used By | Centralized? | Drift Risk |
|---|---|---|---:|---|
| Review modes (`quick`, `standard`, `deep`, `legal-research`, `strategy`) | `routes_catalog.go` | `review.go`, `guide.go`, `routes` | Yes | Low |
| Review document types (`motion`, `opposition`, `response`, `reply`, `declaration`, `proposed-order`, `co-parenting-communication`) | `routes_catalog.go` | `review.go`, `guide.go`, `routes` | Yes | Low |
| Workflow types (`triage`, `counsel-brief`, `sm-triage`, `draft-message`, `evidence-packet`, `pattern-review`, `fact-lock`, `authority-check`, `prep`, `journal-entry`) | `routes_catalog.go` | `workflow.go`, `guide.go`, `routes` | Yes | Low |
| Guided choices / intents | `guide.go` backed by `routes_catalog.go` | `guide`, `guide --list`, `guide --choice`, `guide --intent` | Mostly | Low |
| Conditional specialists | `routes_catalog.go` | `review.go`, `workflow.go`, `routes` | Yes | Medium |
| Legacy persona aliases | `go/internal/persona/loader.go` | Persona loading, validation, tests | Yes | Low |
| Model profiles | `config.default.json`, `go/internal/config/model_routing.go` | `setup`, `review`, `workflow`, `run` | Partially | Medium |
| `review-draft` alias | `review.go` | Root command parsing | Yes | Low |
| `start` alias for guide | `guide.go` | Root command parsing | Yes | Low |

The important distinction is that the catalog centralizes names and metadata, but not the entire prompt assembly or pipeline order. Pipeline order still comes from `config.default.json`.

## Actual Pipeline Sequences

### Intake Pipeline

Source of truth: `config.default.json` `legal_pipelines.intake`

1. `litigation-paralegal` - Document Intake and Registration
2. `source-document-analyst` - Source Document Cards
3. `chronology-clerk` - Chronology Construction
4. `atomic-fact-extractor` - Atomic Fact Extraction
5. `attorney-prep-questioner` - Attorney Prep Questions

### Quick Review

Source of truth: `config.default.json` `legal_pipelines.review_quick`

1. `litigation-paralegal` - Draft Mechanics Review
2. `family-law-attorney-reviewer` - Family Law Attorney Review
3. `opposing-counsel` - Opposing Counsel Review
4. `managing-partner-final-synthesizer` - Managing Partner Final Synthesis

### Standard Review

Source of truth: `config.default.json` `legal_pipelines.review_standard`

1. `litigation-paralegal`
2. `trial-fact-checker`
3. `family-law-attorney-reviewer`
4. `opposing-counsel`
5. `neutral-court-reader`
6. `relief-and-order-alignment-counsel`
7. `legal-writing-preservation-editor`
8. `managing-partner-final-synthesizer`

### Deep Review

Source of truth: `config.default.json` `legal_pipelines.review_deep`

1. `litigation-paralegal`
2. `trial-fact-checker`
3. `family-law-attorney-reviewer`
4. `opposing-counsel`
5. `neutral-court-reader`
6. `relief-and-order-alignment-counsel`
7. `strategic-options-architect`
8. `practical-resolution-reviewer`
9. `legal-writing-preservation-editor`
10. `managing-partner-final-synthesizer`

### Legal Research Review

Source of truth: `config.default.json` `legal_pipelines.review_legal_research`

1. `legal-authority-scholar`
2. `family-law-attorney-reviewer`
3. `managing-partner-final-synthesizer`

### Strategy Review

Source of truth: `config.default.json` `legal_pipelines.review_strategy`

1. `chronology-clerk`
2. `trial-fact-checker`
3. `family-law-attorney-reviewer`
4. `opposing-counsel`
5. `neutral-court-reader`
6. `relief-and-order-alignment-counsel`
7. `strategic-options-architect`
8. `practical-resolution-reviewer`
9. `managing-partner-final-synthesizer`

### Document-Type Routes

| Document type | Maps to | Actual route sequence |
|---|---|---|
| `motion` | Standard review | Uses `review_standard` |
| `opposition` | Standard review | Uses `review_standard` |
| `response` | Standard review | Uses `review_standard` |
| `reply` | Distinct route | `trial-fact-checker` -> `opposing-counsel` -> `neutral-court-reader` -> `legal-writing-preservation-editor` -> `managing-partner-final-synthesizer` |
| `declaration` | Distinct route | `litigation-paralegal` -> `trial-fact-checker` -> `opposing-counsel` -> `legal-writing-preservation-editor` -> `managing-partner-final-synthesizer` |
| `proposed-order` | Distinct route | `litigation-paralegal` -> `relief-and-order-alignment-counsel` -> `neutral-court-reader` -> `practical-resolution-reviewer` -> `managing-partner-final-synthesizer` |
| `co-parenting-communication` | Distinct route | `opposing-counsel` -> `neutral-court-reader` -> `legal-writing-preservation-editor` -> `practical-resolution-reviewer` -> `managing-partner-final-synthesizer` |

### Workflow Routes

Source of truth: `config.default.json` `legal_pipelines.*` plus `workflow.go`

| Workflow type | Actual sequence |
|---|---|
| `triage` | `litigation-paralegal` -> `chronology-clerk` -> `atomic-fact-extractor` -> `family-law-attorney-reviewer` -> `opposing-counsel` -> `neutral-court-reader` -> `relief-and-order-alignment-counsel` -> `strategic-options-architect` -> `practical-resolution-reviewer` -> `attorney-prep-questioner` -> `managing-partner-final-synthesizer` |
| `counsel-brief` | `litigation-paralegal` -> `chronology-clerk` -> `family-law-attorney-reviewer` -> `strategic-options-architect` -> `opposing-counsel` -> `attorney-prep-questioner` -> `legal-writing-preservation-editor` -> `managing-partner-final-synthesizer` |
| `sm-triage` | `chronology-clerk` -> `neutral-court-reader` -> `practical-resolution-reviewer` -> `opposing-counsel` -> `relief-and-order-alignment-counsel` -> `legal-writing-preservation-editor` -> `managing-partner-final-synthesizer` |
| `draft-message` | `opposing-counsel` -> `neutral-court-reader` -> `practical-resolution-reviewer` -> `legal-writing-preservation-editor` -> `managing-partner-final-synthesizer` |
| `evidence-packet` | `source-document-analyst` -> `litigation-paralegal` -> `chronology-clerk` -> `trial-fact-checker` -> `relief-and-order-alignment-counsel` -> `attorney-prep-questioner` -> `managing-partner-final-synthesizer` |
| `pattern-review` | `chronology-clerk` -> `atomic-fact-extractor` -> `trial-fact-checker` -> `family-law-attorney-reviewer` -> `opposing-counsel` -> `neutral-court-reader` -> `strategic-options-architect` -> `managing-partner-final-synthesizer` |
| `fact-lock` | `trial-fact-checker` -> `source-document-analyst` -> `legal-authority-scholar` -> `opposing-counsel` -> `managing-partner-final-synthesizer` |
| `authority-check` | `legal-authority-scholar` -> `family-law-attorney-reviewer` -> `managing-partner-final-synthesizer` |
| `prep` | `litigation-paralegal` -> `chronology-clerk` -> `family-law-attorney-reviewer` -> `opposing-counsel` -> `neutral-court-reader` -> `strategic-options-architect` -> `attorney-prep-questioner` -> `managing-partner-final-synthesizer` |
| `journal-entry` | `litigation-paralegal` -> `chronology-clerk` -> `trial-fact-checker` -> `neutral-court-reader` -> `managing-partner-final-synthesizer` |

## Guided UX

Implemented files:

- `go/cmd/legal-bot/guide.go`
- `go/cmd/legal-bot/guide_test.go`
- `scripts/run-guide.bat`

Behavior verified from source:

- Interactive and deterministic.
- Supports `guide`, `guide <matter>`, `--list`, `--dry-run`, `--choice`, and `--intent`.
- Uses the shared route catalog rather than duplicating route strings.
- Asks only the follow-up questions needed for the chosen route.
- Builds a plan, prints the equivalent command, and confirms before running in interactive mode.
- Routes to the existing `executeIntake`, `executeReview`, and `executeWorkflow` helpers instead of shelling out to a new binary.
- Detects missing matter/input/draft/knowledge conditions and offers sane prompts rather than crashing.
- Can create the standard matter layout when the matter folder does not exist.

Important nuance:

- `--dry-run` suppresses the pipeline run, but workflow plans still write the guided request trace file when enough data exists. So the current implementation is LLM-free, but not fully side-effect-free.

| Guided Choice / Intent | Maps To | Required Inputs | Optional Inputs | Notes |
|---|---|---|---|---|
| 1 / `triage`, `new-situation` | `workflow --type triage` | matter, situation | goal, urgency, audience, child-related, financial | Default goal is the new-situation triage goal |
| 2 / `pursue-preserve-let-go` | `workflow --type triage` | matter, situation | urgency, child-related, financial | Goal emphasizes pursue/preserve/let-go decision-making |
| 3 / `evidence`, `evidence-packet` | `workflow --type evidence-packet` | matter, situation, output need | goal, urgency | Question text captures input files and desired packet type |
| 4 / `counsel`, `counsel-brief` | `workflow --type counsel-brief` | matter, situation | goal, urgency, counsel target, format | Audience defaults to counsel |
| 5 / `review`, `draft-review` | `review` | matter, draft path | mode, document-type, child-related, financial | Motion/opposition/response map to standard review |
| 6 / `sm`, `special-master`, `sm-triage` | `workflow --type sm-triage` | matter, situation | directive, urgency, prior OFW/attorney step, child-related, financial | Audience defaults to Special Master |
| 7 / `message`, `draft-message` | `workflow --type draft-message` | matter, audience, situation | goal, urgency, message format, child-related, financial | Audience drives tone |
| 8 / `pattern`, `pattern-review` | `workflow --type pattern-review` | matter, pattern/topic | date range, related events, goal, child-related, financial | Used to assess recurring-event patterns |
| 9 / `fact-lock` | `workflow --type fact-lock` | matter, facts/claims to verify | audience, authority/legal-claims note, child-related, financial | Focuses on source-supported facts |
| 10 / `authority`, `authority-check` | `workflow --type authority-check` | matter, legal issue/question | jurisdiction, intended use, goal | Does not claim live citator or real-time research status |
| 11 / `prep` | `workflow --type prep` | matter, prep type, situation | goal, urgency/date, child-related, financial | Covers call/hearing/SM conference/etc. |
| 12 / `journal`, `journal-entry` | `workflow --type journal-entry` | matter, what happened | date range, source anchors, category, privileged/internal-only, child-related, financial | Neutral, source-aware journal entry |
| 13 / `intake` | `intake` | matter, source files in `input/` | matter creation | Does not run if `input/` is empty |

Guided input trace storage:

- `matters/<matter>/working/guided_<timestamp>_<workflow-type>.md`

## Model Routing / Model Profiles

The original v2 `openai:*` mismatch is now mostly resolved at the parsing/resolution layer.

What changed:

- `config.default.json` now uses explicit `model_profiles` objects with `primary`, `fallbacks`, `preferred`, and `fallback` compatibility fields.
- `go/internal/config/model_routing.go` parses compact target forms such as `codex:gpt-5.4:medium` and legacy `openai:gpt-5.5`.
- Legacy `openai:*` is normalized to Codex, not silently collapsed into Claude.
- `go/internal/llm/runner.go` supports the four local CLI engines and propagates model/effort where the engine supports it.

Current `model_profiles`:

| model_profile | Primary Engine | Primary Model | Effort | Fallbacks | Used By | Notes |
|---|---|---|---|---|---|---|
| `intake_long_context` | `gemini` | `pro` | none | `claude:sonnet:medium` | Intake | Long-context ingestion profile |
| `structured_extraction` | `codex` | `gpt-5.4-mini` | `medium` | `gemini:flash`, `claude:haiku:medium` | Intake | Used for register/cards/timeline/fact extraction |
| `review_reasoning` | `codex` | `gpt-5.4` | `medium` | `claude:sonnet:medium`, `gemini:pro` | Review/workflow | Default review reasoning profile |
| `final_synthesis` | `codex` | `gpt-5.4` | `high` | `claude:opus:high`, `gemini:pro` | Review/workflow | Final synthesis profile |
| `cheap_fast` | `codex` | `gpt-5.4-mini` | `low` | `gemini:flash-lite`, `claude:haiku:medium` | Utility/compatibility | Fast lower-cost target chain |
| `strategy_reasoning` | `codex` | `gpt-5.4` | `high` | `claude:opus:high`, `gemini:pro` | Review/workflow | Strategy-heavy reasoning profile |
| `legal_research` | `codex` | `gpt-5.4` | `high` | `claude:opus:high`, `gemini:pro` | Review/workflow | Authority-focused profile |

| Model Spec Form | Supported? | Example | Resolution |
|---|---:|---|---|
| `claude` | Yes | `claude` | Uses Claude directly |
| `codex` | Yes | `codex` | Uses Codex directly |
| `gemini` | Yes | `gemini` | Uses Gemini directly |
| `copilot` | Yes | `copilot` | Uses Copilot directly |
| `codex:gpt-5.4:medium` | Yes | `codex:gpt-5.4:medium` | Engine=`codex`, model=`gpt-5.4`, effort=`medium` |
| `codex:gpt-5.5:high` | Yes | `codex:gpt-5.5:high` | Engine=`codex`, model=`gpt-5.5`, effort=`high` |
| `openai:gpt-5.5` | Yes, legacy alias | `openai:gpt-5.5` | Normalizes to `codex` + `gpt-5.5` |
| `claude:sonnet:high` | Yes | `claude:sonnet:high` | Engine=`claude`, model=`sonnet`, effort=`high` |
| `gemini:pro` | Yes | `gemini:pro` | Engine=`gemini`, model=`pro` |
| `gemini:flash` | Yes | `gemini:flash` | Engine=`gemini`, model=`flash` |
| `copilot:auto` | Yes | `copilot:auto` | Engine=`copilot`, model=`auto` |

Current model-routing reality:

- `review`, `workflow`, and `run` can all accept compact target specs.
- `--single-llm` still works and can now take a compact target spec.
- The active engine is still selected by local CLI availability on `PATH`.
- If a requested engine CLI is missing, the runner logs a warning and falls back to Claude.
- `setup` still maintains an `llms` availability map, but that is a settings/UX layer, not the route-control plane.

## LLM Runner Behavior

| Engine | Executable | Model Argument | Effort Argument | Prompt Mode | Notes |
|---|---|---|---|---|---|
| `claude` | `claude` | `--model` | `--effort` | `-p <prompt>` | Fallback target when other engines are unavailable |
| `codex` | `codex` | `--model` | `-c model_reasoning_effort=<effort>` | `codex exec ... -o <tmpfile> <prompt>` | Writes to a temp file, then reads the result |
| `gemini` | `gemini` | `--model` | none | positional prompt with `--yolo` | No effort arg passed |
| `copilot` | `copilot` | `--model` | none | `--prompt <prompt>` | No effort arg passed |

Current runner behavior:

- Persona body is injected from `agents/<persona>.md` after YAML frontmatter is stripped.
- Prompt text is passed directly to the local CLI for the chosen engine.
- Stdout is captured; stderr is either passed through or captured for logging.
- Output is written to the designated output file when one is supplied.
- Run metadata is logged to `~/.config/legal-bot/logs/legal-bot.log` unless disabled.
- `resolveLLM` falls back to Claude if `codex`, `gemini`, or `copilot` is requested but the executable is not found.
- Tests avoid real LLM calls and instead validate argument construction, persona loading, and logging behavior.

## Conditional Routing

| Trigger | Agent Added | Applies To | Implemented In | Centralized? | Duplicate Protection? | Notes |
|---|---|---|---|---:|---:|---|
| `--child-related`, child/custody keywords, child-heavy route hints | `child-best-interests-family-dynamics-reviewer` | Review and workflow | `routes_catalog.go`, `review.go`, `workflow.go` | Yes | Yes | Catch-all for child, custody, school, therapy, medical, and parent-time issues |
| `--financial`, support keywords, financial route hints | `financial-support-reviewer` | Review and workflow | `routes_catalog.go`, `review.go`, `workflow.go` | Yes | Yes | Adds financial/support review |
| Legal authority keywords or `--mode legal-research` | `legal-authority-scholar` | Review | `routes_catalog.go`, `review.go` | Yes | Yes | Review-only authority specialist |
| Strategy keywords or `--mode strategy` | `strategic-options-architect` | Review | `routes_catalog.go`, `review.go` | Yes | Yes | Review-only strategy specialist |
| Document-type-driven child/authority/strategy hints | Same as above | Review | `review.go` | Partly | Yes | Motion/declaration/custody/school/therapy/medical routes trigger heuristics |
| Workflow content heuristics | Child/financial specialists | Workflow | `workflow.go` | Partly | Yes | Workflow heuristics look at situation/goal/inputs and route metadata |

The route catalog names the rules, but `review.go` and `workflow.go` still own insertion timing. That is intentional for now and should be treated as a code-driven heuristic layer rather than a pure config layer.

## Agent Loading and Aliases

Active legal personas on disk:

- 17 root persona files in `agents/`
- shared guidance files in `agents/shared/`

The active root personas are:

- `atomic-fact-extractor`
- `attorney-prep-questioner`
- `child-best-interests-family-dynamics-reviewer`
- `chronology-clerk`
- `family-law-attorney-reviewer`
- `financial-support-reviewer`
- `legal-authority-scholar`
- `legal-writing-preservation-editor`
- `litigation-paralegal`
- `managing-partner-final-synthesizer`
- `neutral-court-reader`
- `opposing-counsel`
- `practical-resolution-reviewer`
- `relief-and-order-alignment-counsel`
- `source-document-analyst`
- `strategic-options-architect`
- `trial-fact-checker`

| Alias | Resolves To | Active Target Exists? | Notes |
|---|---|---:|---|
| `intake-mapper` | `litigation-paralegal` | Yes | Legacy legal alias |
| `document-card-generator` | `source-document-analyst` | Yes | Legacy legal alias |
| `fact-extractor` | `atomic-fact-extractor` | Yes | Legacy legal alias |
| `timeline-builder` | `chronology-clerk` | Yes | Legacy legal alias |
| `open-questions-generator` | `attorney-prep-questioner` | Yes | Legacy legal alias |
| `fact-auditor` | `trial-fact-checker` | Yes | Legacy legal alias |
| `legal-sufficiency-reviewer` | `family-law-attorney-reviewer` | Yes | Legacy legal alias |
| `court-reader` | `neutral-court-reader` | Yes | Legacy legal alias |
| `attack-surface-reviewer` | `opposing-counsel` | Yes | Legacy legal alias |
| `relief-alignment-reviewer` | `relief-and-order-alignment-counsel` | Yes | Legacy legal alias |
| `preservation-editor` | `legal-writing-preservation-editor` | Yes | Legacy legal alias |
| `bulldog-advocate` | `strategic-options-architect` | Yes | Alias only; no active file |
| `final-synthesizer` | `managing-partner-final-synthesizer` | Yes | Legacy legal alias |

Observations:

- `bulldog-advocate.md` is no longer an active persona file.
- Persona frontmatter `model_profile` is still present and used as route metadata, but active legal routing now sits above it in the Go control plane.
- `llm.Run` injects persona body text, not frontmatter, into the prompt context.

## Embedded Assets / Generation

- Embedded agent and config content still exists in `go/internal/embedded/agents_generated.go`.
- The generator source is `go/internal/embedded/gen.go`.
- Regeneration command remains `cd go && go generate ./internal/embedded/`.
- `embedded_test.go` verifies that the compiled-in assets contain the expected agent names and default config JSON keys.
- I did not independently re-run code generation in this audit; however, the embedded tests and current source layout indicate the embedded assets are part of the active build contract.

## Prompt Context Assembly

| Pipeline / Step Type | Reads Case? | Reads Knowledge? | Reads Draft? | Reads Raw Input? | Reads Prior Outputs? | Writes |
|---|---:|---:|---:|---:|---:|---|
| Intake | Yes | Yes, as it is created | No | Yes | Yes | `knowledge/` and `output/ingestion_report.md` |
| Review | Yes | Yes | Yes | No | Yes | `working/*.md` and `output/final_review_packet.md` |
| Workflow | Yes | Yes if present | Sometimes | Yes (`situation`, `goal`, `questions`, draft) | Yes | `working/*` and `output/<workflow-type>.md` |
| Guide | No separate pipeline context | Indirectly, by checking matter state | Optional | Interactive input | Plan-only metadata | `working/guided_<timestamp>_<workflow-type>.md` for workflow plans |
| Run | No | No | No | Prompt text only | No | Optional output file |
| Historian | No | No | No | Target directory contents | Prompt file if present | `input-history.md` |
| Routes | No | No | No | No | No | Standard output only |

Important distinction:

- `case/` is prompt context, not a mechanical parser input.
- `matters/<matter>/knowledge/` is the intake-created knowledge layer.
- `matters/<matter>/drafts/` is the review input layer.

## Input and Output Storage Map

| Path | Purpose | Overwrite Behavior | Notes |
|---|---|---|---|
| `matters/<matter>/input/` | Source files for intake and workflow context | Source files are user-managed | Requires `.md` or `.txt` files for intake |
| `matters/<matter>/drafts/` | Drafts for review | Latest draft may be chosen automatically | Review needs an explicit or discoverable draft |
| `matters/<matter>/knowledge/` | Intake outputs | Static names are reused | Contains `document_register.csv`, `document_cards/`, `case_timeline.md`, `fact_candidates.md`, `open_questions_for_user.md` |
| `matters/<matter>/working/` | Step artifacts and intermediate outputs | Step files are numbered/static per run; guided traces are timestamped | Workflow step files and guided request notes land here |
| `matters/<matter>/output/` | Final pipeline outputs | Static filenames are reused per route | Review writes `final_review_packet.md`; workflow writes `<workflow-type>.md` |
| `case/` | Global prompt context | File edits affect every run | Root `.md`/`.txt` files are auto-loaded |
| `scripts/` | Windows wrappers | Script contents are static | Active wrappers are kept minimal |
| `bin/` | CLI wrappers and legacy stubs | Some wrappers are intentionally stubs | Discovery/council/oracle/generate no longer look like normal commands |

Current overwrite behavior:

- Intake writes stable knowledge filenames.
- Review writes numbered `working/` step outputs and a stable `output/final_review_packet.md`.
- Workflow writes numbered `working/<type>-<step>-<slug>.md` files and a stable `output/<workflow-type>.md`.
- Guide writes timestamped `working/guided_<timestamp>_<workflow-type>.md` files for workflow plans.

## Case-Wide vs Matter-Specific Knowledge

Current `case/` inventory:

- `source_hierarchy.md`
- `people_and_roles.md`
- `preferred_framing.md`
- `language_calibration.md`
- `relief_library.md`

Behavior verified from `case/README.md`:

- Every root `.md` and `.txt` file in `case/` is loaded into prompt context.
- Subfolders are not the loading contract.
- `case/` is global prompt guidance, not a record archive.
- Matter-specific source files belong in `matters/<matter>/input/`.

Files that do not currently exist in `case/`:

- `legal_authority_guidance.md`
- `strategy_principles.md`
- `avoid_language.md`

## Active vs Legacy Pipeline Code

| Path | Purpose | Active Legal Routing? | Utility/Partial? | Legacy/Compatibility? | Status | Notes |
|---|---|---:|---:|---:|---|---|
| `go/cmd/legal-bot/intake.go` | Intake command | Yes | No | No | active | Current intake control path |
| `go/cmd/legal-bot/review.go` | Review command | Yes | No | No | active | Current review control path |
| `go/cmd/legal-bot/workflow.go` | Workflow command | Yes | No | No | active | Current workflow control path |
| `go/cmd/legal-bot/guide.go` | Guided front door | Yes | No | No | active | Deterministic operator UX layer |
| `go/cmd/legal-bot/routes_catalog.go` | Route catalog | Yes | No | No | active | Centralized route/control-plane metadata |
| `go/cmd/legal-bot/routes_command.go` | Route diagnostics | Yes | Yes | No | active | `legal-bot routes` and `--json` |
| `go/internal/config` | Config schema/loaders | Yes | Yes | No | active | Holds pipelines, profiles, modes |
| `go/internal/llm` | LLM subprocess runner | Yes | Yes | No | active | Engine/model/effort dispatch and logging |
| `go/internal/persona` | Persona loader/aliases | Yes | Yes | No | active | Legacy alias compatibility remains |
| `go/internal/embedded` | Embedded config/personas | Yes | Yes | No | active | Binary fallback source of truth |
| `agents/` | Persona bodies | Yes | No | No | active | 17 active legal personas plus shared docs |
| `case/` | Global prompt context | Yes | No | No | active | Prompt context only |
| `scripts/` | Windows helpers | Yes | Yes | No | active | `run-guide.bat` added |
| `bin/` | Wrappers and stubs | Partial | Yes | Yes | partial/legacy | Discovery/council/oracle/generate are now explicit legacy stubs |
| `go/internal/pipeline` | Historian/legacy orchestration | No | Yes | Yes | partial | Still used by historian and legacy-oriented utilities |
| `go/internal/council` | VernHole council support | No | No | Yes | legacy | Not part of the active legal workflow contract |
| `go/internal/generate` | Legacy persona-generation helpers | No | No | Yes | legacy | Compatibility / historical |
| `go/internal/vts` | Vern Task Spec helpers | No | No | Yes | legacy | Historical substrate |
| `go/internal/tobeads` | VTS-to-Beads tooling | No | No | Yes | legacy | Historical substrate |
| `go/internal/tui` | Bubble Tea utility UI | No | Yes | Yes | partial | Still exposed, but not the primary workflow front door |

## Legacy Surface Inventory

| Surface | Path | Type | Calls Existing Root Command? | Status | Replacement / Current Path | Notes |
|---|---|---|---:|---|---|---|
| `legal-bot-discovery` | `bin/legal-bot-discovery` | Wrapper | No | retired | `legal-bot guide`, `legal-bot intake`, `legal-bot review`, `legal-bot workflow` | Now an explicit legacy stub |
| `legal-bot-discovery.cmd` | `bin/legal-bot-discovery.cmd` | Windows wrapper | No | retired | Same as above | Explicit legacy stub |
| `legal-bot-council` | `bin/legal-bot-council` | Wrapper | No | retired | `legal-bot tui` only if you intentionally want the legacy utility UI | Explicit legacy stub |
| `legal-bot-council.cmd` | `bin/legal-bot-council.cmd` | Windows wrapper | No | retired | Same as above | Explicit legacy stub |
| `legal-bot-oracle` | `bin/legal-bot-oracle` | Wrapper | No | retired | None in the active legal workflow | Explicit legacy stub |
| `legal-bot-generate` | `bin/legal-bot-generate` | Wrapper | No | retired | None in the active legal workflow | Explicit legacy stub |
| `legal-bot-run` | `bin/legal-bot-run` | Wrapper | Yes, to `run` | partial | `legal-bot run` | Compatibility helper |
| `legal-bot-run.cmd` | `bin/legal-bot-run.cmd` | Windows wrapper | Yes, to `run` | partial | `legal-bot run` | Compatibility helper |
| `legal-bot-historian` | `bin/legal-bot-historian` | Wrapper | Yes, to `historian` | partial | `legal-bot historian` | Compatibility helper |
| `install-legal-bot-cli` | `bin/install-legal-bot-cli` | Installer helper | No | partial | `cd go && go build -o bin/legal-bot ./cmd/legal-bot` | Not part of the legal workflow |
| `install-legal-bot-cli.cmd` | `bin/install-legal-bot-cli.cmd` | Windows installer helper | No | partial | `cd go && go build -o bin/legal-bot.exe ./cmd/legal-bot` | Not part of the legal workflow |
| `go/internal/pipeline` | `go/internal/pipeline` | Package | N/A | partial | Keep for historian/legacy utility flows | Still referenced by historian |
| `go/internal/council` | `go/internal/council` | Package | N/A | legacy | No active legal routing dependency | VernHole legacy substrate |
| `go/internal/generate` | `go/internal/generate` | Package | N/A | legacy | No active legal routing dependency | Legacy persona-generation substrate |
| `go/internal/vts` | `go/internal/vts` | Package | N/A | legacy | No active legal routing dependency | Historical VTS substrate |
| `go/internal/tobeads` | `go/internal/tobeads` | Package | N/A | legacy | No active legal routing dependency | Historical import substrate |
| `go/internal/tui` | `go/internal/tui` | Package | N/A | partial | `legal-bot tui` | Still exposed, but clearly utility/legacy-oriented |
| `docs/legacy_surfaces.md` | `docs/legacy_surfaces.md` | Doc | N/A | active documentation | Use as inventory | Canonical legacy inventory now |

Bottom line: the old Vern wrappers are no longer silent traps. They now either point to current commands, or they loudly tell the operator they are legacy.

## Windows Scripts and Operator Helpers

| Script | Active? | Command Called | Supports Advanced Flags? | Notes |
|---|---:|---|---:|---|
| `scripts/run-intake.bat` | Yes | `legal-bot intake <matter>` | No | Requires `input/` files |
| `scripts/run-review.bat` | Yes | `legal-bot review <matter>` | No | Requires knowledge layer and draft files |
| `scripts/run-guide.bat` | Yes | `legal-bot guide [matter]` | Yes, via CLI flags | Added for the guided front door |
| `scripts/open-latest-output.bat` | Yes | Opens latest output | No | Existing helper retained |

Not currently present:

- `scripts/run-workflow.bat`
- `scripts/open-latest-workflow-output.bat`

## Documentation Audit

| Doc/File | Current Claim | Accurate? | Correction Needed? | Notes |
|---|---|---:|---:|---|
| `README.md` | Active commands include `guide`, `routes`, `intake`, `review`, `workflow`, plus utility commands | Yes | No | Now reflects the control-plane split and legacy boundary |
| `docs/architecture.md` | Active contract, route catalog, model routing, and legacy surfaces are separated | Yes | No | Strongly aligned with current code |
| `docs/user_workflow.md` | Guide and routes are operator entry points; review/document-type and workflow-type distinctions are explained | Yes | No | Good operational doc |
| `docs/agent_persona_system.md` | Legacy `openai:*` strings are accepted as aliases to Codex | Yes | No | Matches code and runner behavior |
| `docs/token_strategy.md` | Focuses on engine/model/effort strategy rather than old Vern wrappers | Mostly | Minor | It is still primarily a token-strategy doc, not a control-plane doc |
| `docs/legacy_surfaces.md` | Canonical legacy inventory | Yes | No | New authoritative legacy reference |
| `docs/legacy_discovery_migration.md` | Historical migration note | Mostly | No | Still useful, but should not be mistaken for the active operator contract |
| `case/README.md` | Root `.md`/`.txt` files are prompt context, not deterministic parsing rules | Yes | No | Correctly explains `case/` behavior |

Stale claims that are now corrected in the current docs set:

- `openai:*` no longer silently collapses into Claude.
- `guide` is now a real front door rather than a concept.
- `routes` now exists as a diagnostic command.
- Legacy Vern wrappers are explicitly marked legacy rather than shown as current workflow.
- `case/` is prompt context, not a parsed archive format.

## Tests and Validation

| Test File | What It Covers | Real LLM Calls? | Notes |
|---|---|---:|---|
| `go/cmd/legal-bot/routes_catalog_test.go` | Route catalog completeness, route mappings, conditional duplication protection, JSON snapshot, root commands, legacy wrapper documentation | No | Strongest control-plane regression coverage |
| `go/cmd/legal-bot/guide_test.go` | Guided choice mapping, intent aliases, equivalent-command generation, newest draft detection | No | Confirms deterministic route selection |
| `go/cmd/legal-bot/model_routing_test.go` | Override precedence and compact model spec behavior | No | Verifies engine/model/effort override path |
| `go/cmd/legal-bot/workflow_test.go` | Workflow spec completeness and prompt context assembly | No | Confirms active workflow catalog |
| `go/internal/config/config_test.go` | Model target parsing, profile resolution, legacy `openai:*` normalization, runnable target selection, active agents, default routes | No | Good model/control-plane coverage |
| `go/internal/persona/loader_test.go` | Persona loading and legacy alias resolution | No | Confirms alias map and persona parsing |
| `go/internal/embedded/embedded_test.go` | Embedded agent/config asset integrity | No | Confirms generated asset contents |
| `go/internal/llm/runner_test.go` | Engine arg construction, logging, persona injection, fallback behavior, truncation | No | Validates subprocess contract without real model calls |

## Test Run Results

Attempted command:

```powershell
$env:GOCACHE='C:\GIT\legal-bot\.go-cache\gocache'; $env:GOMODCACHE='C:\GIT\legal-bot\.go-cache\gomodcache'; $env:GOTMPDIR='C:\GIT\legal-bot\.go-cache\tmp'; & 'C:\Program Files\Go\bin\go.exe' test ./...
```

Result:

- Failed in the environment before compilation completed.
- Exact error: `failed to initialize build cache at C:\GIT\legal-bot\.go-cache\gocache: mkdir C:\GIT\legal-bot\.go-cache\gocache\00: Access is denied.`

I also tried a `C:\tmp`-based cache path earlier in the inspection flow and hit a similar Go temp/cache permissions error. The failure is environmental rather than a reported code compile/test failure.

## Build Results

- Not run in this audit.

## Common Failure Modes

These are the failure modes that still apply after this inspection:

- Missing local CLI executable on `PATH` for the selected engine.
- Bad model target spec or unsupported compact form.
- `review` before intake, or review without a draft.
- `intake` with no `.md` or `.txt` source files in `input/`.
- `guide` on a new matter without confirming folder creation.
- `guide` not finding a draft when choice 5 is selected.
- `workflow` without a knowledge layer, which is allowed but weaker.
- `workflow` with a missing `--audience` on routes that require one.
- `openai:*` compatibility being mistaken for a true OpenAI API path.
- Legacy wrapper confusion if someone ignores the explicit stubs.
- Static output filenames being overwritten by design.
- Go build cache or temp directory permission problems in restricted environments.
- Stale embedded assets if `config.default.json` or `agents/*.md` change without regeneration.

## Recommended Follow-Up Tasks

### Must Fix

- Keep `guide` dry-run semantics documented as LLM-free but not fully side-effect free, or change the implementation later if a truly no-write dry run is desired.
- Preserve the current route catalog as the single source for route names; avoid reintroducing duplicated literal route strings in command code or docs.

### Should Fix

- Add a dedicated `legal-bot models` or `legal-bot doctor`-style diagnostic only if Karl actually needs model selection visibility separate from `routes`.
- Consider a workflow wrapper script if Windows operator convenience becomes important.
- Keep the docs and `config.default.json` in sync whenever pipelines or agent names change.

### Could Fix Later

- Make guided trace writes optional behind a flag if traceability and dry-run side effects need to be separated.
- Add a more formal machine-readable route inventory if future tooling wants to consume it outside the CLI.
- Consider whether `tui` should remain exposed as a root command or become a clearly legacy-only utility.

### Do Not Chase Yet

- Do not replace the config-driven pipeline order with a new routing framework.
- Do not move matter-specific evidence or guidance into `case/`.
- Do not try to rewrite the remaining Vern-era historical packages unless there is a concrete breakage.
- Do not change the legal personas or prompt substance as part of control-plane cleanup.

## Final Assessment

- Legal-Bot is safer and easier to operate than v2.
- The biggest remaining technical risk is that route names are now clearer, but pipeline order and context assembly are still mixed between declarative config and code.
- The biggest remaining UX risk is that `guide --dry-run` still has a small write-side effect, and that `tui`/legacy wrappers can still be mistaken for primary workflows by a casual operator.
- The biggest remaining repo-maintenance risk is keeping docs, config, embedded assets, and wrapper inventories synchronized as routes or personas evolve.
- Karl should next use the current route catalog and tests as the baseline, then keep tightening the separation between active legal control-plane code and historical Vern compatibility surfaces.
