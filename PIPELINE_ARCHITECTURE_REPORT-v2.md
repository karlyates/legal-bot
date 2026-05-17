# Pipeline Architecture Report

## Executive Summary

Legal-Bot is currently a mixed system:

- The active legal pipelines are defined declaratively in `config.default.json`.
- The Go CLI decides which pipeline to run, assembles prompt context, injects conditional specialists, and writes outputs.
- Agent bodies are loaded from `agents/*.md` when present, or from embedded generated assets when not.
- The repo is mostly aligned for the current legal intake/review/workflow product, but the model layer still has a major mismatch: config profiles point at `openai:*` names that the runner cannot execute, so those profiles fall back to supported local CLI engines.
- The biggest operator risk is assuming the legacy Vern-era discovery/council/oracle/generate wrappers are part of the current legal workflow. They are not.

What a user should know first:

- Intake must run before review.
- Review needs a generated knowledge layer and a draft in `matters/<matter>/drafts/`.
- Workflow runs are draftless or draft-light jobs driven by `--type`.
- The current legal product is the Go CLI plus `config.default.json`, not the old Vern discovery surfaces.

## Active Pipeline Definition Locations

| Area | File/Path | Role | Active? | Notes |
|---|---|---:|---:|---|
| Legal pipeline definitions | [config.default.json](c:/GIT/legal-bot/config.default.json) | Declares `legal_pipelines`, `model_profiles`, and `llm_modes` | Yes | This is the main declarative source for intake, review, and workflow step order. |
| Legacy discovery config | [config.default.json](c:/GIT/legal-bot/config.default.json) | Declares `discovery_pipelines` and `pipeline_mode` | Legacy | Still loaded by config, but not used by the current legal intake/review/workflow commands. |
| Intake command | [go/cmd/legal-bot/intake.go](c:/GIT/legal-bot/go/cmd/legal-bot/intake.go) | Runs the intake pipeline, writes knowledge artifacts, chains steps | Yes | Uses `config.default.json` plus case and input context. |
| Review command | [go/cmd/legal-bot/review.go](c:/GIT/legal-bot/go/cmd/legal-bot/review.go) | Selects review mode/document-type routes and adds conditional specialists | Yes | This is where quick/standard/deep/legal-research/strategy routing is resolved. |
| Workflow command | [go/cmd/legal-bot/workflow.go](c:/GIT/legal-bot/go/cmd/legal-bot/workflow.go) | Selects draftless workflow routes and adds conditional specialists | Yes | This is where triage, counsel-brief, evidence-packet, authority-check, prep, and journal-entry routes are resolved. |
| Config loader | [go/internal/config/config.go](c:/GIT/legal-bot/go/internal/config/config.go) | Loads user/project/embedded config and resolves model profiles | Yes | It also preserves legacy discovery config compatibility. |
| Persona loader | [go/internal/persona/loader.go](c:/GIT/legal-bot/go/internal/persona/loader.go) | Parses agent frontmatter and legacy aliases | Yes | Used by roster loading and low-level persona handling. |
| Embedded assets generator | [go/internal/embedded/gen.go](c:/GIT/legal-bot/go/internal/embedded/gen.go) | Regenerates embedded agents and default config | Yes | Generates [go/internal/embedded/agents_generated.go](c:/GIT/legal-bot/go/internal/embedded/agents_generated.go). |
| Embedded assets | [go/internal/embedded/agents_generated.go](c:/GIT/legal-bot/go/internal/embedded/agents_generated.go) | Compiled-in agent markdown and default config | Yes | Used when on-disk agents/config are unavailable. |
| Agent markdown files | [agents/](c:/GIT/legal-bot/agents) | Persona definitions and frontmatter | Yes | The legal CLI loads these for prompt context when available. |
| Case guidance files | [case/](c:/GIT/legal-bot/case) | Global case-wide prompt context | Yes | Every root `.md`/`.txt` file is loaded automatically into intake/workflow/review prompts. |
| LLM subprocess runner | [go/internal/llm/runner.go](c:/GIT/legal-bot/go/internal/llm/runner.go) | Launches `claude`, `codex`, `gemini`, or `copilot` CLIs | Yes | It injects persona body text and logs runs to the user config directory. |
| Legacy discovery/council/oracle pipeline code | [go/internal/pipeline/](c:/GIT/legal-bot/go/internal/pipeline) | Vern-era discovery, historian, oracle, VernHole helpers | Mixed | Still used by `historian` and TUI/legacy surfaces, not by the main legal routing paths. |
| Legacy council roster and tiers | [go/internal/council/](c:/GIT/legal-bot/go/internal/council) | Vern council selection logic | Mixed | Supports old council/VernHole flows. |
| Legacy generation utility | [go/internal/generate/](c:/GIT/legal-bot/go/internal/generate) | Persona scaffold/generation helpers | Mixed | Present in code, but there is no current root CLI `generate` subcommand. |
| Windows wrappers | [scripts/](c:/GIT/legal-bot/scripts) | Intake/review/open-latest helpers | Yes | These are the main documented operator wrappers for the current legal product. |
| Legacy shell wrappers | [bin/](c:/GIT/legal-bot/bin) | Vern-era wrappers for discovery/council/oracle/generate/historian | Legacy | Several wrappers call subcommands that the current Go CLI does not expose. |

## Actual Pipeline Sequences

### Intake Pipeline

Defined in `legal_pipelines.intake` in [config.default.json](c:/GIT/legal-bot/config.default.json) and executed by [go/cmd/legal-bot/intake.go](c:/GIT/legal-bot/go/cmd/legal-bot/intake.go).

1. `litigation-paralegal`
2. `source-document-analyst`
3. `chronology-clerk`
4. `atomic-fact-extractor`
5. `attorney-prep-questioner`

### Quick Review

Defined in `legal_pipelines.review_quick`.

1. `litigation-paralegal`
2. `family-law-attorney-reviewer`
3. `opposing-counsel`
4. `managing-partner-final-synthesizer`

### Standard Review

Defined in `legal_pipelines.review_standard`.

1. `litigation-paralegal`
2. `trial-fact-checker`
3. `family-law-attorney-reviewer`
4. `opposing-counsel`
5. `neutral-court-reader`
6. `relief-and-order-alignment-counsel`
7. `legal-writing-preservation-editor`
8. `managing-partner-final-synthesizer`

### Deep Review

Defined in `legal_pipelines.review_deep`.

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

Defined in `legal_pipelines.review_legal_research`.

1. `legal-authority-scholar`
2. `family-law-attorney-reviewer`
3. `managing-partner-final-synthesizer`

### Strategy Review

Defined in `legal_pipelines.review_strategy`.

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

The document-type route is resolved in [go/cmd/legal-bot/review.go](c:/GIT/legal-bot/go/cmd/legal-bot/review.go). These are the active document-specific review sequences:

#### Reply

1. `trial-fact-checker`
2. `opposing-counsel`
3. `neutral-court-reader`
4. `legal-writing-preservation-editor`
5. `managing-partner-final-synthesizer`

#### Declaration

1. `litigation-paralegal`
2. `trial-fact-checker`
3. `opposing-counsel`
4. `legal-writing-preservation-editor`
5. `managing-partner-final-synthesizer`

#### Proposed Order

1. `litigation-paralegal`
2. `relief-and-order-alignment-counsel`
3. `neutral-court-reader`
4. `practical-resolution-reviewer`
5. `managing-partner-final-synthesizer`

#### Co-Parenting Communication

1. `opposing-counsel`
2. `neutral-court-reader`
3. `legal-writing-preservation-editor`
4. `practical-resolution-reviewer`
5. `managing-partner-final-synthesizer`

#### Route Aliases That Fold Into Standard Review

`motion`, `opposition`, and `response` do not have separate pipelines. They map to `review_standard`.

### Draftless Workflow Routes

Workflow routes are defined in `workflowSpecs()` in [go/cmd/legal-bot/workflow.go](c:/GIT/legal-bot/go/cmd/legal-bot/workflow.go).

#### Triage

1. `litigation-paralegal`
2. `chronology-clerk`
3. `atomic-fact-extractor`
4. `family-law-attorney-reviewer`
5. `opposing-counsel`
6. `neutral-court-reader`
7. `relief-and-order-alignment-counsel`
8. `strategic-options-architect`
9. `practical-resolution-reviewer`
10. `attorney-prep-questioner`
11. `managing-partner-final-synthesizer`

#### Counsel Brief

1. `litigation-paralegal`
2. `chronology-clerk`
3. `family-law-attorney-reviewer`
4. `strategic-options-architect`
5. `opposing-counsel`
6. `attorney-prep-questioner`
7. `legal-writing-preservation-editor`
8. `managing-partner-final-synthesizer`

#### Special Master Triage

1. `chronology-clerk`
2. `neutral-court-reader`
3. `practical-resolution-reviewer`
4. `opposing-counsel`
5. `relief-and-order-alignment-counsel`
6. `legal-writing-preservation-editor`
7. `managing-partner-final-synthesizer`

#### Draft Message

1. `opposing-counsel`
2. `neutral-court-reader`
3. `practical-resolution-reviewer`
4. `legal-writing-preservation-editor`
5. `managing-partner-final-synthesizer`

#### Evidence Packet

1. `source-document-analyst`
2. `litigation-paralegal`
3. `chronology-clerk`
4. `trial-fact-checker`
5. `relief-and-order-alignment-counsel`
6. `attorney-prep-questioner`
7. `managing-partner-final-synthesizer`

#### Pattern Review

1. `chronology-clerk`
2. `atomic-fact-extractor`
3. `trial-fact-checker`
4. `family-law-attorney-reviewer`
5. `opposing-counsel`
6. `neutral-court-reader`
7. `strategic-options-architect`
8. `managing-partner-final-synthesizer`

#### Fact Lock

1. `trial-fact-checker`
2. `source-document-analyst`
3. `legal-authority-scholar`
4. `opposing-counsel`
5. `managing-partner-final-synthesizer`

#### Authority Check

1. `legal-authority-scholar`
2. `family-law-attorney-reviewer`
3. `managing-partner-final-synthesizer`

#### Prep

1. `litigation-paralegal`
2. `chronology-clerk`
3. `family-law-attorney-reviewer`
4. `opposing-counsel`
5. `neutral-court-reader`
6. `strategic-options-architect`
7. `attorney-prep-questioner`
8. `managing-partner-final-synthesizer`

#### Journal Entry

1. `litigation-paralegal`
2. `chronology-clerk`
3. `trial-fact-checker`
4. `neutral-court-reader`
5. `managing-partner-final-synthesizer`

## Agent Communication / Data Flow

There is no direct agent-to-agent API call path. Agents talk to each other only through:

- generated files on disk
- prior agent output copied into the next prompt
- compact knowledge artifacts in `matters/<matter>/knowledge/`
- case-wide prompt context from `case/`

High-level flow:

```text
input files / draft / case guidance
        |
        v
intake or review or workflow command
        |
        v
llm.Run(...) injects persona body from agents/*.md
        |
        v
step output written to working/ or knowledge/ or output/
        |
        v
next step receives prior output text plus selected file context
        |
        v
final synthesizer writes final packet
```

Important details:

- `llm.Run` injects only the body of the persona markdown file, not its YAML frontmatter.
- Intake step 1 feeds raw source documents plus case context.
- Intake step 2 receives the previous step output and raw source documents again.
- Intake steps 3-5 receive the generated knowledge layer, not the full raw source corpus.
- Review steps 1-6 receive the draft plus the knowledge layer.
- Review step 7 receives the draft plus previous review findings.
- Review final synthesis receives the draft filename and all previous review findings, not the raw draft body.
- Workflow steps receive case context, situation/goal/questions, optional draft context, knowledge, and an input manifest. Later steps also receive prior workflow findings.

## Pipeline Inputs and Outputs

| Step | Reads | Writes | Notes |
|---|---|---|---|
| Intake step 1 | `case/` root `.md`/`.txt`; `matters/<matter>/input/*.md|*.txt` | `matters/<matter>/working/01-intake-mapping.md`; `matters/<matter>/knowledge/document_register.csv` | The CSV is extracted from step output if present. |
| Intake step 2 | `case/`; raw input files; step 1 output | `matters/<matter>/working/02-document-cards.md`; `matters/<matter>/knowledge/document_cards/*.md` | The combined card output is split into per-document files. |
| Intake step 3 | `case/`; `knowledge/` files and document cards | `matters/<matter>/knowledge/case_timeline.md` | Rebuilds the knowledge layer before the prompt. |
| Intake step 4 | `case/`; `knowledge/` files and document cards | `matters/<matter>/knowledge/fact_candidates.md` | Same knowledge layer context as step 3. |
| Intake step 5 | `case/`; `knowledge/` files and document cards | `matters/<matter>/knowledge/open_questions_for_user.md` | Final intake question pass. |
| Intake report | Matter metadata | `matters/<matter>/output/ingestion_report.md` | Written after all intake steps. |
| Review steps 1..n-1 | `case/`; `knowledge/`; latest or specified draft; prior review findings | `matters/<matter>/working/<step>-<slug>.md` | Names are fixed and overwritten on rerun. |
| Review final step | `case/`; `knowledge/`; draft filename; all review findings | `matters/<matter>/output/final_review_packet.md` | The final packet is always the same filename. |
| Workflow steps 1..n-1 | `case/`; situation; goal; questions; optional draft; knowledge; input manifest; prior findings | `matters/<matter>/working/<type>-<step>-<slug>.md` | Uses route-specific working filenames. |
| Workflow final step | Same as above plus accumulated findings | `matters/<matter>/output/<workflow-type>.md` | Example: `triage.md`, `counsel-brief.md`, `authority-check.md`. |
| Historian command | Arbitrary target directory contents | `input-history.md` in target directory | Also appends a note to `prompt.md` if present. |
| Legacy Oracle/VernHole surfaces | Discovery/VTS files, synthesis files, council outputs | `oracle-vision.md`, `oracle-architect-breakdown.md`, `synthesis.md`, `vts/*.md` | These are legacy compatibility outputs, not the main legal pipeline. |

## Conditional Routing

| Trigger | Agent(s) Added | Implemented How | File/Path | Notes |
|---|---|---|---|---|
| `--child-related` in review | `child-best-interests-family-dynamics-reviewer` | `applyConditionalReviewRouting` | [go/cmd/legal-bot/review.go](c:/GIT/legal-bot/go/cmd/legal-bot/review.go) | Also auto-triggers on certain document types and draft keywords in deep review. |
| `--financial` in review | `financial-support-reviewer` | `applyConditionalReviewRouting` | [go/cmd/legal-bot/review.go](c:/GIT/legal-bot/go/cmd/legal-bot/review.go) | Also auto-triggers on financial keywords in deep review. |
| Legal-research mode or authority-heavy document/draft keywords in review | `legal-authority-scholar` | `applyConditionalReviewRouting` | [go/cmd/legal-bot/review.go](c:/GIT/legal-bot/go/cmd/legal-bot/review.go) | Inserted after `trial-fact-checker`. |
| Strategy mode or strategy-heavy document/draft keywords in review | `strategic-options-architect` | `applyConditionalReviewRouting` | [go/cmd/legal-bot/review.go](c:/GIT/legal-bot/go/cmd/legal-bot/review.go) | Inserted before `practical-resolution-reviewer`, or before final synthesis if needed. |
| `--child-related` in workflow | `child-best-interests-family-dynamics-reviewer` | `applyConditionalWorkflowRouting` | [go/cmd/legal-bot/workflow.go](c:/GIT/legal-bot/go/cmd/legal-bot/workflow.go) | Triggered by flag or situation/goal/questions/draft text heuristics. |
| `--financial` in workflow | `financial-support-reviewer` | `applyConditionalWorkflowRouting` | [go/cmd/legal-bot/workflow.go](c:/GIT/legal-bot/go/cmd/legal-bot/workflow.go) | Triggered by flag or heuristic keywords. |
| Document-type selection in review | Route-specific base pipeline | `reviewPipelineKey` | [go/cmd/legal-bot/review.go](c:/GIT/legal-bot/go/cmd/legal-bot/review.go) | Document type overrides mode selection. |
| Workflow type selection | Route-specific workflow pipeline | `workflowSpecs()` | [go/cmd/legal-bot/workflow.go](c:/GIT/legal-bot/go/cmd/legal-bot/workflow.go) | The `--type` flag is required. |

Key observations:

- Review routing is partly declarative and partly heuristic.
- Workflow routing is also partly heuristic, but it only supports child/financial conditional additions today.
- There is no separate CLI flag for conditional legal-authority or strategy routing in workflow mode; those are driven by the selected workflow type and content heuristics.

## Model Profile Resolution

Model routing is mixed:

- Declarative config names live in `config.default.json`.
- The Go config loader resolves those names to runnable engines.
- The runner only supports local CLI engines: `claude`, `codex`, `gemini`, and `copilot`.

| model_profile | Configured Provider/Model | Fallback | Used By | Notes |
|---|---|---|---|---|
| `intake_long_context` | `gemini` | `claude` | Intake step 1 | Runnable because `gemini` is a supported local CLI engine. |
| `structured_extraction` | `openai:gpt-5.4-mini` | `claude` | Intake steps 2-4 and some legacy review steps | Not runnable as OpenAI; the provider prefix is stripped and falls back to `claude`. |
| `review_reasoning` | `openai:gpt-5.5` | `claude` | Most review/workflow steps | Same mismatch as above. |
| `final_synthesis` | `openai:gpt-5.5` | `claude` | All final synthesizer steps | Same mismatch as above. |
| `cheap_fast` | `openai:gpt-5.4-mini` | `claude` | Present in config only | Not used by the current legal pipelines in `config.default.json`. |
| `strategy_reasoning` | `openai:gpt-5.5` | `claude` | Strategy-oriented steps | Falls back to `claude` under the current runner. |
| `legal_research` | `openai:gpt-5.5` | `claude` | Legal authority steps | Falls back to `claude` under the current runner. |

Resolution flow:

1. `config.Load()` reads user config, then project config, then embedded config, then hardcoded defaults.
2. `review.go`, `workflow.go`, and `intake.go` look up the step's `model_profile`.
3. `ResolveModelProfile()` looks up the profile's preferred model and fallback model.
4. `runnableEngine()` only accepts `claude`, `gemini`, `codex`, and `copilot`, optionally with provider prefixes.
5. If the preferred string is not one of those engines, the fallback is tried.
6. If neither resolves, the original step `llm` value is used.
7. `llm.Run()` then resolves the requested engine again by checking the local CLI availability.

What this means in practice:

- `model_profile` values are active routing control in config.
- Agent frontmatter `model_profile` is parsed, but the legal pipeline runner does not use it to choose the engine.
- Agent frontmatter `model` still matters for persona roster and legacy generation utilities.
- OpenAI-branded model names in config are documentation-level targets today, not directly runnable engines.

## Agent Loading and Aliases

Agent loading is split across two paths:

- `persona.Load()` and `persona.LoadFile()` parse markdown frontmatter and body.
- `llm.Run()` reads only the body of `agents/<persona>.md` and injects it into the prompt.

Key behaviors:

- If an agent file exists on disk, it is used.
- If the agent file is missing, the embedded copy is used.
- Legacy persona names are mapped to current names by `ResolveName()`.
- The current legal agent set is 17 personas.

Legacy alias map in [go/internal/persona/loader.go](c:/GIT/legal-bot/go/internal/persona/loader.go):

- `intake-mapper` -> `litigation-paralegal`
- `document-card-generator` -> `source-document-analyst`
- `fact-extractor` -> `atomic-fact-extractor`
- `timeline-builder` -> `chronology-clerk`
- `open-questions-generator` -> `attorney-prep-questioner`
- `fact-auditor` -> `trial-fact-checker`
- `legal-sufficiency-reviewer` -> `family-law-attorney-reviewer`
- `court-reader` -> `neutral-court-reader`
- `attack-surface-reviewer` -> `opposing-counsel`
- `relief-alignment-reviewer` -> `relief-and-order-alignment-counsel`
- `preservation-editor` -> `legal-writing-preservation-editor`
- `bulldog-advocate` -> `strategic-options-architect`
- `final-synthesizer` -> `managing-partner-final-synthesizer`

Important nuance:

- `bulldog-advocate` is a compatibility alias, not an active file in `agents/`.
- The current legal routes refer to `strategic-options-architect`, not `bulldog-advocate`.

## Embedded Assets / Generation

Embedded assets are generated from:

- all `agents/*.md` files
- `config.default.json`

Generation path:

- [go/internal/embedded/gen.go](c:/GIT/legal-bot/go/internal/embedded/gen.go) reads the repo files and writes [go/internal/embedded/agents_generated.go](c:/GIT/legal-bot/go/internal/embedded/agents_generated.go).
- The embedded package exposes `GetAgent()` and `GetDefaultConfig()`.
- `config.Load()` falls back to embedded config if an on-disk config file is not found.
- `persona.Load()` and `council.ScanRoster()` fall back to embedded agent content if the on-disk `agents/` directory is missing.

Regeneration command found:

```powershell
cd go
go generate ./internal/embedded/
```

What requires regeneration:

- changing any `agents/*.md`
- changing `config.default.json`

What does not require regeneration:

- changing docs only
- changing `case/*.md`
- changing matter input/draft/output files

## Active vs Legacy Pipeline Code

| Path | Purpose | Active Legal Routing? | Legacy/Compatibility? | Notes |
|---|---|---:|---:|---|
| [go/cmd/legal-bot/intake.go](c:/GIT/legal-bot/go/cmd/legal-bot/intake.go) | Legal intake pipeline | Yes | No | Main source-to-knowledge compression layer. |
| [go/cmd/legal-bot/review.go](c:/GIT/legal-bot/go/cmd/legal-bot/review.go) | Draft review pipeline | Yes | No | Handles quick/standard/deep/legal-research/strategy and document types. |
| [go/cmd/legal-bot/workflow.go](c:/GIT/legal-bot/go/cmd/legal-bot/workflow.go) | Draftless workflow routing | Yes | No | Handles triage, counsel brief, evidence packet, authority check, prep, and journal entry. |
| [go/internal/config/config.go](c:/GIT/legal-bot/go/internal/config/config.go) | Config loading and model profile resolution | Yes | No | Also preserves legacy discovery config compatibility. |
| [go/internal/llm/runner.go](c:/GIT/legal-bot/go/internal/llm/runner.go) | Launches local LLM subprocesses | Yes | No | Only supports `claude`, `codex`, `gemini`, and `copilot`. |
| [go/internal/persona/loader.go](c:/GIT/legal-bot/go/internal/persona/loader.go) | Persona parsing and aliasing | Yes | Partial | Alias map includes old legal names. |
| [go/internal/embedded/](c:/GIT/legal-bot/go/internal/embedded) | Embedded fallback data | Yes | No | Active support for binaries without on-disk assets. |
| [agents/](c:/GIT/legal-bot/agents) | Current legal personas | Yes | No | 17 active personas. |
| [case/](c:/GIT/legal-bot/case) | Case-wide prompt context | Yes | No | Loaded automatically for intake, workflow, and review. |
| [scripts/run-intake.bat](c:/GIT/legal-bot/scripts/run-intake.bat) | Windows intake wrapper | Yes | No | Builds the binary if needed and runs `intake`. |
| [scripts/run-review.bat](c:/GIT/legal-bot/scripts/run-review.bat) | Windows review wrapper | Yes | No | Builds the binary if needed and runs `review`. |
| [scripts/open-latest-output.bat](c:/GIT/legal-bot/scripts/open-latest-output.bat) | Open review packet helper | Yes | No | Opens only `final_review_packet.md`. |
| [bin/legal-bot-run.cmd](c:/GIT/legal-bot/bin/legal-bot-run.cmd) | Low-level Windows LLM runner wrapper | Partial | Legacy-ish | Wraps the `run` command. |
| [bin/legal-bot-historian](c:/GIT/legal-bot/bin/legal-bot-historian) | Historian wrapper | Partial | Legacy-ish | Wraps the `historian` command. |
| [bin/](c:/GIT/legal-bot/bin) | Discovery/council/oracle/generate legacy wrappers | No | Yes | See `bin/legal-bot-discovery`, `bin/legal-bot-council`, `bin/legal-bot-oracle`, and `bin/legal-bot-generate`; these wrappers call subcommands the current Go CLI does not expose. |
| [go/internal/pipeline/](c:/GIT/legal-bot/go/internal/pipeline) | Historian, VernHole, oracle, event helpers | Mixed | Yes | Active for `historian` and TUI, not for the legal workflows. |
| [go/internal/council/](c:/GIT/legal-bot/go/internal/council) | Vern council selection | No | Yes | Legacy/compat support. |
| [go/internal/generate/](c:/GIT/legal-bot/go/internal/generate) | Persona scaffolding | No | Yes | Code exists, but no current root CLI subcommand exposes it. |
| [go/internal/vts/](c:/GIT/legal-bot/go/internal/vts) | VTS parser/writer | No | Yes | Used by Oracle/VernHole legacy flows. |
| [go/internal/tobeads/](c:/GIT/legal-bot/go/internal/tobeads) | Beads importer | No | Yes | Legacy conversion pipeline. |
| [go/internal/tui/](c:/GIT/legal-bot/go/internal/tui) | Interactive Vern-era UI | No | Yes | Separate interactive surface, not the legal CLI contract. |

## How To Run Legal-Bot Today

### Setup Requirements

- Go 1.25.7+ is documented in [README.md](c:/GIT/legal-bot/README.md).
- At least one supported local CLI must be installed: `claude`, `gemini`, `codex`, or `copilot`.
- `case/` should contain case-wide guidance files.
- A matter folder should exist under `matters/<matter>/`.

### Build

Verified in the repo docs and build scripts:

```powershell
cd go
go build -o bin/legal-bot.exe ./cmd/legal-bot
```

### Intake

```powershell
legal-bot intake <matter>
```

Expected input:

- `matters/<matter>/input/` must exist.
- It must contain at least one `.md` or `.txt` file.

What it writes:

- `matters/<matter>/knowledge/document_register.csv`
- `matters/<matter>/knowledge/document_cards/`
- `matters/<matter>/knowledge/case_timeline.md`
- `matters/<matter>/knowledge/fact_candidates.md`
- `matters/<matter>/knowledge/open_questions_for_user.md`
- `matters/<matter>/output/ingestion_report.md`

### Review

```powershell
legal-bot review <matter>
```

Useful flags:

- `--mode quick`
- `--mode standard`
- `--mode deep`
- `--mode legal-research`
- `--mode strategy`
- `--document-type reply`
- `--document-type declaration`
- `--document-type proposed-order`
- `--document-type co-parenting-communication`
- `--child-related`
- `--financial`
- `--draft <path>`
- `--single-llm <llm>`
- `--llm-mode <mode>`

What it reads:

- `matters/<matter>/knowledge/document_register.csv` must exist.
- It loads the latest `.md` or `.txt` draft in `matters/<matter>/drafts/` unless `--draft` is given.

What it writes:

- `matters/<matter>/working/<step>-<slug>.md`
- `matters/<matter>/output/final_review_packet.md`

### Workflow

```powershell
legal-bot workflow <matter> --type triage
```

Supported `--type` values:

- `triage`
- `counsel-brief`
- `sm-triage`
- `draft-message`
- `evidence-packet`
- `pattern-review`
- `fact-lock`
- `authority-check`
- `prep`
- `journal-entry`

Helpful flags:

- `--audience`
- `--urgency`
- `--situation`
- `--goal`
- `--questions`
- `--draft`
- `--child-related`
- `--financial`
- `--single-llm`
- `--llm-mode`

### Mode Selection

- `review` chooses the route from `--document-type` first, then `--mode`.
- `workflow` chooses the route from `--type`.
- `review-draft` is an alias of `review`.

### Document-Type Selection

- `reply`, `declaration`, `proposed-order`, and `co-parenting-communication` are distinct routes.
- `motion`, `opposition`, and `response` are accepted aliases for the standard review pipeline.

### Child/Financial/Strategy/Legal-Research Routing

- Child routing is supported in both review and workflow.
- Financial routing is supported in both review and workflow.
- Legal-authority routing is supported in review and in the `authority-check` workflow.
- Strategy routing is supported in review and in deep/strategy-related workflows.

### Other Commands

- `legal-bot run <llm> <prompt>` is the low-level one-shot subprocess wrapper.
- `legal-bot setup` writes user config to `~/.config/legal-bot/config.json`.
- `legal-bot historian <directory>` indexes arbitrary directories into `input-history.md`.
- `legal-bot tui` launches the interactive legacy UI.

### Windows vs Mac/Linux

- The repo clearly documents and scripts the Windows PowerShell/batch path.
- The Go binary itself is cross-platform in principle.
- The shell wrappers under `bin/` are mostly Vern-era compatibility surfaces, not the current legal workflow contract.

## Storage Map: Case-Wide vs Matter-Specific

| Information Type | Where It Belongs | Loaded Automatically? | Generated or User-Edited? | Notes |
|---|---|---:|---:|---|
| Case-wide people/roles | `case/people_and_roles.md` | Yes | User-edited | Global prompt context for stable names and roles. |
| Case-wide source hierarchy | `case/source_hierarchy.md` | Yes | User-edited | Defines source tiers and safe-use discipline. |
| Preferred framing | `case/preferred_framing.md` | Yes | User-edited | Narrative posture and emphasis. |
| Language calibration | `case/language_calibration.md` | Yes | User-edited | Loaded as prompt context; the code does not parse or enforce it mechanically. |
| Relief library | `case/relief_library.md` | Yes | User-edited | Global relief hooks and relief-to-evidence mapping. |
| Legal authority guidance | No dedicated case file today | No dedicated file | Mixed | The repo currently uses `agents/shared/legal-authority-rules.md` plus workflow outputs instead of a dedicated `case/legal_authority_guidance.md`. |
| Strategy principles | No dedicated case file today | No dedicated file | Mixed | The functional equivalent is `agents/shared/strategy-safety-rules.md` and strategy-oriented agents/outputs. |
| Durable source documents | `matters/<matter>/input/` | No | User-edited | Intake reads `.md` and `.txt` source files from here. |
| Motion-specific drafts | `matters/<matter>/drafts/` | Yes for review | User-edited | Review picks the newest `.md`/`.txt` unless `--draft` is given. |
| Proposed orders | `matters/<matter>/drafts/` or a draft-specific file passed to review/workflow | Yes when selected | User-edited | No special folder exists today. |
| Declarations | `matters/<matter>/drafts/` | Yes when selected | User-edited | Treated as a draft type, not a separate storage subsystem. |
| Exhibits | `matters/<matter>/input/` or draft-local supporting files | No direct automatic read in review | User-edited | Best practice is to keep them near the matter and reference them through the knowledge layer. |
| Attorney notes | `matters/<matter>/working/` or `matters/<matter>/input/` | Sometimes | User-edited | Draftless workflows can read `working/user_situation.md` as a fallback. |
| Matter-specific legal research | `matters/<matter>/knowledge/` and workflow outputs | Yes after generated | Generated | Legal research is generated by review/workflow steps, not by a special folder. |
| Generated document cards | `matters/<matter>/knowledge/document_cards/` | Yes after generation | Generated | Created by intake step 2. |
| Generated fact candidates | `matters/<matter>/knowledge/fact_candidates.md` | Yes after generation | Generated | Created by intake step 4. |
| Generated timeline | `matters/<matter>/knowledge/case_timeline.md` | Yes after generation | Generated | Created by intake step 3. |
| Review findings | `matters/<matter>/working/` and `matters/<matter>/output/final_review_packet.md` | Yes after review | Generated | Intermediate findings are accumulated through step outputs. |
| Final packet | `matters/<matter>/output/final_review_packet.md` | Yes | Generated | The final review artifact. |

## Matter Folder Lifecycle

The repo expects this lifecycle:

1. `input/` holds raw source documents for intake.
2. `knowledge/` holds intake-generated knowledge artifacts.
3. `drafts/` holds draft documents to review.
4. `working/` holds intermediate outputs from intake, workflow, and review runs.
5. `output/` holds final packets and reports.

Practical notes:

- `input/` and `drafts/` are gitignored.
- `knowledge/` and `output/` are not gitignored by default.
- The named filenames are static, so reruns overwrite previous artifacts.

## Where To Store Case-Wide Knowledge

Use `case/` for durable, shared guidance that should shape every matter:

- source hierarchy
- people and roles
- preferred framing
- language calibration
- relief library

Do not put these in `case/`:

- raw evidence dumps
- giant matter-specific archives
- throwaway notes that only matter to one case file

Reality check:

- `case/` is loaded automatically, but only root `.md` and `.txt` files are read.
- The system does not have a dedicated parser for these guidance files; they influence the model through prompt context.

## Where To Store Motion-Specific Knowledge

Use `matters/<matter>/` for motion- or project-specific material:

- source documents in `input/`
- drafts in `drafts/`
- generated knowledge in `knowledge/`
- intermediate scratch files in `working/`
- final handoff artifacts in `output/`

There is no dedicated motion-specific subfolder convention today beyond those standard matter folders.

If you need separate motion-specific research, the current repo expects you to keep it alongside the matter as generated knowledge or a draft-local artifact, not in `case/`.

## Where Do Outputs Go

| Output Type | Location | Notes |
|---|---|---|
| Intake outputs | `matters/<matter>/knowledge/` and `matters/<matter>/output/ingestion_report.md` | Static filenames, overwritten on rerun. |
| Review outputs | `matters/<matter>/working/` and `matters/<matter>/output/final_review_packet.md` | Intermediate step files plus one final packet. |
| Workflow outputs | `matters/<matter>/working/` and `matters/<matter>/output/<workflow-type>.md` | One final packet per workflow type. |
| Historian output | `<target>/input-history.md` | Also updates `prompt.md` when present. |
| Legacy oracle/VernHole outputs | Discovery/VTS directories | Legacy outputs only. |
| Logs | `~/.config/legal-bot/logs/legal-bot.log` | Written by `llm.Run` unless logging is disabled. |
| User config | `~/.config/legal-bot/config.json` | Written by `legal-bot setup`. |

Latest-output helper:

- [scripts/open-latest-output.bat](c:/GIT/legal-bot/scripts/open-latest-output.bat) opens only `matters/<matter>/output/final_review_packet.md`.
- It does not search for the newest workflow output.

Overwrite behavior:

- Intake, review, and workflow outputs use stable filenames.
- The repo does not timestamp the main artifacts.
- Reruns replace the previous files.

## Prompt Context Assembly

| Agent/Pipeline Step | Receives Case Guidance? | Receives Knowledge Artifacts? | Receives Draft? | Receives Prior Outputs? | Writes |
|---|---:|---:|---:|---:|---|
| Intake step 1 | Yes | No | No | No | Working intake mapping; register extraction. |
| Intake step 2 | Yes | No | No | Yes | Document cards. |
| Intake step 3 | Yes | Yes | No | No | Case timeline. |
| Intake step 4 | Yes | Yes | No | No | Fact candidates. |
| Intake step 5 | Yes | Yes | No | No | Open questions. |
| Review quick/standard/deep base steps | Yes | Yes | Yes | Only later steps | Working review files and final packet. |
| Review legal-research step | Yes | Yes | Yes | Only later steps | Legal authority findings. |
| Review strategy step | Yes | Yes | Yes | Only later steps | Strategy findings. |
| Review final synthesizer | Yes | Yes | Draft filename only in `all_reviews` mode | Yes | `final_review_packet.md`. |
| Workflow base steps | Yes | Yes when present | Yes when present | Only later steps | Working workflow files and final packet. |
| Workflow final synthesizer | Yes | Yes | Yes when supplied in the route | Yes | Final workflow output. |
| `legal-bot run` | No automatic case knowledge | No | Optional persona body only | No | Optional output file via `--output`. |

Additional prompt-assembly facts:

- Each step gets the persona body from `agents/<persona>.md`.
- Intake uses the full raw source files from `input/`.
- Review does not re-read the source corpus; it uses the knowledge layer plus the draft.
- Workflow does not dump the whole raw source corpus; it uses the situation/goal/questions files plus the knowledge layer and manifest.

## Customization Map

| Customization Goal | Where To Change It | Risk Level | Requires Tests? | Requires Embedded Regeneration? | Notes |
|---|---|---:|---:|---:|---|
| Change language guidance | `case/language_calibration.md` | Low | No | No | Prompt-only influence. |
| Change preferred framing | `case/preferred_framing.md` | Low | No | No | Prompt-only influence. |
| Change relief library | `case/relief_library.md` | Low-Medium | No | No | Affects relief alignment and strategy. |
| Add or edit people/roles | `case/people_and_roles.md` | Low | No | No | Prompt-only influence. |
| Add matter-specific source material | `matters/<matter>/input/` | Low | No | No | Intake input only. |
| Add a draft | `matters/<matter>/drafts/` | Low | No | No | Review input only. |
| Change review mode routing | `go/cmd/legal-bot/review.go` | High | Yes | No | This changes pipeline selection and conditional routing. |
| Change workflow routing | `go/cmd/legal-bot/workflow.go` | High | Yes | No | This changes route resolution and context assembly. |
| Change agent order | `config.default.json` | Medium-High | Yes | Usually yes | Needed if the change should work in embedded binaries. |
| Add a new agent persona | `agents/*.md`, alias maps, tests, embedded regen | High | Yes | Yes | Also update `go/internal/persona/loader.go` and `go/internal/embedded/`. |
| Rename an agent | `agents/*.md`, alias maps, tests, embedded regen | High | Yes | Yes | Stale aliases are the biggest breakage risk. |
| Change model profile routing | `config.default.json` and `go/internal/config/config.go` | High | Yes | Yes if config is embedded | Keep runnable-engine limits in mind. |
| Change low-level LLM support | `go/internal/llm/runner.go` | Very High | Yes | No | This is where new provider CLIs would need to be wired. |
| Change generated fallback data | `config.default.json` or `agents/*.md` | High | Yes | Yes | Run `go generate ./internal/embedded/` afterward. |
| Change latest output helper | `scripts/open-latest-output.bat` | Low | No | No | Only affects the review-packet helper. |

Safe customization path:

- edit `case/*.md`
- add matter files in `matters/<matter>/input/`, `drafts/`, or `working/`
- use the existing CLI flags for route selection

Risky customization path:

- changing `config.default.json` route order
- renaming personas without updating aliases and embedded assets
- changing `model_profile` names without checking `runnableEngine()`

## Implemented vs Documented Features

| Feature | Documented? | Implemented? | File/Path | Notes |
|---|---|---:|---|---|
| Legal intake | Yes | Yes | [go/cmd/legal-bot/intake.go](c:/GIT/legal-bot/go/cmd/legal-bot/intake.go) | Active. |
| Quick review | Yes | Yes | [config.default.json](c:/GIT/legal-bot/config.default.json) / [go/cmd/legal-bot/review.go](c:/GIT/legal-bot/go/cmd/legal-bot/review.go) | Active. |
| Standard review | Yes | Yes | Same as above | Active. |
| Deep review | Yes | Yes | Same as above | Active. |
| Legal-research review | Yes | Yes | Same as above | Active. |
| Strategy review | Yes | Yes | Same as above | Active. |
| `--child-related` | Yes | Yes | [go/cmd/legal-bot/review.go](c:/GIT/legal-bot/go/cmd/legal-bot/review.go) | Active conditional route. |
| `--financial` | Yes | Yes | [go/cmd/legal-bot/review.go](c:/GIT/legal-bot/go/cmd/legal-bot/review.go) | Active conditional route. |
| `--document-type` | Yes | Yes | [go/cmd/legal-bot/review.go](c:/GIT/legal-bot/go/cmd/legal-bot/review.go) | Active, with aliases. |
| Proposed-order route | Yes | Yes | [go/cmd/legal-bot/review.go](c:/GIT/legal-bot/go/cmd/legal-bot/review.go) | Active. |
| Declaration route | Yes | Yes | Same as above | Active. |
| Reply route | Yes | Yes | Same as above | Active. |
| Co-parenting-communication route | Yes | Yes | Same as above | Active. |
| `language_calibration.md` loading | Yes | Yes | [case/README.md](c:/GIT/legal-bot/case/README.md) | Loaded automatically as case context; not parsed by code. |
| Optional case files | Partially | Yes | [case/README.md](c:/GIT/legal-bot/case/README.md) / [go/cmd/legal-bot/intake.go](c:/GIT/legal-bot/go/cmd/legal-bot/intake.go) | Any root case `.md`/`.txt` files are loaded. |
| `model_profile` routing | Yes | Yes | [go/internal/config/config.go](c:/GIT/legal-bot/go/internal/config/config.go) | Active, but only for runnable local engines. |
| OpenAI/Codex/Gemini/Claude support | Partially | Mixed | [go/internal/llm/runner.go](c:/GIT/legal-bot/go/internal/llm/runner.go) | `claude`, `gemini`, `codex`, and `copilot` are supported; `openai:*` config names are not directly runnable. |
| Embedded regeneration | Yes | Yes | [go/internal/embedded/gen.go](c:/GIT/legal-bot/go/internal/embedded/gen.go) | `go generate` command found. |
| Dedicated `strategy_principles.md` case file | Not in repo | No | `case/` | Functional equivalent lives in shared agent rules and strategy outputs. |
| Dedicated `legal_authority_guidance.md` case file | Not in repo | No | `case/` | Functional equivalent lives in shared agent rules and authority outputs. |
| Current legal CLI `generate` subcommand | Documented by legacy wrappers, not the current Go CLI | No | `bin/legal-bot-generate` / missing root command | Legacy surface only. |
| Current legal CLI `discovery`/`hole`/`oracle` subcommands | Documented by legacy wrappers, not the current Go CLI | No | `bin/legal-bot-*` legacy wrappers | Legacy surface only. |

## Minimum Working Example

This is the smallest current end-to-end path based on repo behavior:

1. Install at least one supported local CLI: `claude`, `gemini`, `codex`, or `copilot`.
2. Build the Go binary.

```powershell
cd go
go build -o bin/legal-bot.exe ./cmd/legal-bot
```

3. Create a matter folder with one source file.

```text
matters/SAMPLE/input/example.md
```

4. Run intake.

```powershell
legal-bot intake SAMPLE
```

5. Create a draft.

```text
matters/SAMPLE/drafts/motion.md
```

6. Run review.

```powershell
legal-bot review SAMPLE
```

7. Read the final packet.

```text
matters/SAMPLE/output/final_review_packet.md
```

What should exist after a successful run:

- `matters/SAMPLE/knowledge/document_register.csv`
- `matters/SAMPLE/knowledge/document_cards/`
- `matters/SAMPLE/knowledge/case_timeline.md`
- `matters/SAMPLE/knowledge/fact_candidates.md`
- `matters/SAMPLE/knowledge/open_questions_for_user.md`
- `matters/SAMPLE/output/ingestion_report.md`
- `matters/SAMPLE/output/final_review_packet.md`

## Common Failure Modes

- Missing local CLI: `claude`, `gemini`, `codex`, or `copilot` is not installed.
- Wrong working directory or unresolved repo root.
- Missing `matters/<matter>/input/` or no `.md`/`.txt` files there.
- Running review before intake.
- Missing `matters/<matter>/drafts/` or no draft file there.
- Assuming workflow reads the full raw source corpus. It does not.
- Assuming `case/private/` is loaded. It is not.
- Renamed agent file without updating the alias map and embedded assets.
- `openai:*` model profile expectation mismatch.
- Stale embedded assets in a binary that does not have on-disk `agents/` or `config.default.json`.
- `open-latest-output.bat` only opens the review packet, not workflow outputs.
- Static filenames causing old outputs to be overwritten unexpectedly.
- Large or noisy source files causing long run times or LLM timeout pressure.
- Legacy `bin/*` wrappers invoking missing root CLI subcommands.

## Documentation Mismatches

| Doc/File | Claim | Actual Behavior | Recommended Fix |
|---|---|---|---|
| [README.md](c:/GIT/legal-bot/README.md) | The current product includes several legal workflows and review modes. | Mostly true. | Add one explicit note that `openai:*` model profile targets are not directly runnable and currently fall back to local CLI engines. |
| [README.md](c:/GIT/legal-bot/README.md) | Document-type examples cover the main review types. | True but incomplete. | Document that `motion`, `opposition`, and `response` are aliases to standard review. |
| [docs/user_workflow.md](c:/GIT/legal-bot/docs/user_workflow.md) | Workflow and review commands are the main user path. | True. | Clarify that workflow prompt context is synthesized from scenario files and the knowledge layer, not raw source dumps. |
| [docs/agent_persona_system.md](c:/GIT/legal-bot/docs/agent_persona_system.md) | `model_profile` is the active routing control. | True, but limited by runner support. | Add that only `claude`, `gemini`, `codex`, and `copilot` are runnable and `openai:*` entries fall back. |
| [docs/token_strategy.md](c:/GIT/legal-bot/docs/token_strategy.md) | Provider-neutral model profiles are used. | Partly true. | State that some config entries are provider-qualified placeholders rather than directly runnable models. |
| [docs/architecture.md](c:/GIT/legal-bot/docs/architecture.md) | Legacy discovery/council/TUI surfaces remain in the repo. | True. | Add a warning that some legacy shell wrappers still point at missing root CLI subcommands. |
| [case/README.md](c:/GIT/legal-bot/case/README.md) | Case files are loaded into prompts and affect all workflows. | True. | Add a note that the system does not parse them mechanically; they influence prompts through LLM context only. |
| `bin/legal-bot-discovery`, `bin/legal-bot-council`, `bin/legal-bot-oracle`, `bin/legal-bot-generate` | They look like runnable CLI entrypoints. | Not with the current Go root command. | Mark them clearly as legacy compatibility wrappers or remove them from operator guidance. |
| `case/` layout implied by user expectations | Dedicated `strategy_principles.md` and `legal_authority_guidance.md` files exist. | They do not exist in the repo. | Either add those files or update docs and guidance to point at the shared agent rules and generated outputs instead. |

## Risks / Unknowns

- I could not complete `go test ./...` in this environment because the Go toolchain hit cache and module-cache permission issues outside the repo workspace.
- Some legacy Vern-era shell wrappers still exist, but the current legal CLI does not expose matching subcommands, so their reliability is uncertain.
- The model profile config uses OpenAI-style strings that the current runner cannot execute directly.
- The case guidance files are loaded as prompt context, not parsed into deterministic rules, so their effect depends on the model following the prompt.

## Recommended Follow-Up Tasks

### Docs-Only

- Document that `motion`, `opposition`, and `response` fold into standard review.
- Document that `openai:*` profile names are not directly runnable today.
- Mark the Vern-era wrapper scripts as legacy compatibility surfaces.
- Clarify that case files are prompt context, not a parsed rule engine.

### Config

- Consider replacing the `openai:*` profile targets with runnable engine names or a real OpenAI runner.
- Consider adding dedicated case files if `strategy_principles.md` and `legal_authority_guidance.md` are intended as first-class guidance.
- Consider whether workflow and review alias coverage should be centralized so config and code cannot drift.

### Code

- Decide whether legacy `bin/*` wrappers should be kept, renamed, or retired.
- Decide whether `model_profile` should be validated against runnable engines at load time.
- Consider a clearer abstraction for conditional routing so heuristic inserts are easier to reason about.

### Tests

- Add or strengthen tests for `openai:*` fallback behavior.
- Add tests for `motion` / `opposition` / `response` alias handling.
- Add tests for the legacy wrapper/documentation mismatch if those scripts are retained.
- Add tests that confirm case files are read only from the root of `case/`.

### Future Design

- If the goal is a cleaner operator experience, add a single guided entrypoint that selects workflow type, draft type, and conditional routing automatically.
- If OpenAI models are intended, add a real runner instead of config-only placeholders.
- If the product should support per-matter custom guidance, add a matter overlay layer instead of overloading `case/`.

## Commands Run

Inspection commands used during this investigation:

- `Get-ChildItem -Force`
- `rg --files`
- `rg -n` across `config.default.json`, `README.md`, `docs/`, `agents/`, `case/`, `go/`, `scripts/`, and `bin/`
- `Get-Content` on the key Go, config, script, and docs files

Validation attempts:

- `go test ./...` failed because `go` was not on PATH in this shell.
- `& 'C:\Program Files\Go\bin\go.exe' test ./...` failed because the environment could not create or use the Go build cache and module cache directories.
