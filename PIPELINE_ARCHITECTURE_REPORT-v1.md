# Pipeline Architecture Report

## Executive Summary

Active Legal-Bot pipelines are stored primarily in `config.default.json` under `legal_pipelines`, with important routing and context behavior hard-coded in `go/cmd/legal-bot/intake.go`, `go/cmd/legal-bot/review.go`, and model resolution logic in `go/internal/config/config.go`. The system is mixed: base sequences are declarative, but conditional routing, prompt assembly, file loading, and some fallback behavior live in Go code.

Config, code, and docs are mostly aligned on the 17-agent panel and review modes, but there are still meaningful gaps:
- conditional routing is more aggressive than the docs suggest
- `case/README.md` and any extra root `case/*.md` / `case/*.txt` files are auto-loaded into prompts
- legal research mode is implemented as a draft-based review mode, not a standalone authority-research integration
- embedded assets exist, but active persona prompt injection still loads agent markdown from disk

The biggest operator risk is assuming the system is more isolated and structured than it is today. In practice, prompt context is assembled from a broad file sweep of `case/` and `knowledge/`, outputs are overwritten in place, and wrapper scripts only expose the default review path. Before running Legal-Bot, a user should know: use `.txt` or `.md` inputs only, run from the repo with a working Go 1.25.7+ toolchain and at least one supported LLM CLI, and expect the `case/` folder to affect every matter in the repo.

## Active Pipeline Definition Locations

| Area | File/Path | Role | Active? | Notes |
|---|---|---|---|---|
| Pipeline config | `config.default.json` | Declares intake and review pipeline sequences, step order, context modes, `model_profile`, prompt prefixes | Yes | Core active pipeline contract |
| Intake CLI | `go/cmd/legal-bot/intake.go` | Runs intake pipeline, loads `case/`, input docs, and writes knowledge artifacts | Yes | Hard-codes file IO and prompt assembly |
| Review CLI | `go/cmd/legal-bot/review.go` | Resolves review mode, document type, conditional routing, and final outputs | Yes | Hard-codes route selection and conditional agent insertion |
| Config/model resolution | `go/internal/config/config.go` | Loads config, resolves pipeline definitions, LLM modes, and model profiles | Yes | Provider-qualified OpenAI profile values currently fall back to supported local engines |
| LLM runner | `go/internal/llm/runner.go` | Executes external LLM CLIs and injects persona markdown from disk | Yes | Active review/intake prompt assembly depends on `agents/*.md` on disk |
| Persona metadata / aliases | `go/internal/persona/loader.go` | Parses agent frontmatter and resolves old agent aliases | Partly | Alias map is active; full loader is not the main path for prompt injection |
| Embedded config/agents | `go/internal/embedded/*` | Embeds `agents/*.md` and `config.default.json` for fallback/build packaging | Partly | Generated and tested, but not the primary source for active prompt persona loading |
| Windows wrappers | `scripts/run-intake.bat`, `scripts/run-review.bat`, `scripts/open-latest-output.bat` | Convenience launchers | Yes | Default path only; no advanced routing flags exposed |
| Tests | `go/internal/config/config_test.go`, `go/internal/persona/loader_test.go`, `go/internal/embedded/embedded_test.go` | Validate route names, agent existence, aliases, embedded assets | Yes | Good reference for active/legal expectations |
| Docs | `README.md`, `docs/*.md`, `case/README.md` | Operator guidance | Yes, but descriptive | Not authoritative over code/config |
| Legacy pipeline packages | `go/internal/pipeline/*`, `go/internal/council/*`, `go/internal/tui/*`, `go/internal/generate/*` | Older discovery/council/oracle/TUI flows | No for active legal intake/review | Still callable or test-covered in parts |

## Actual Pipeline Sequences

The sequences below are the actual configured pipelines before conditional insertions in `go/cmd/legal-bot/review.go`.

### Intake Pipeline

Configured pipeline key: `intake`

1. `litigation-paralegal`
2. `source-document-analyst`
3. `chronology-clerk`
4. `atomic-fact-extractor`
5. `attorney-prep-questioner`

### Quick Review

Configured pipeline key: `review_quick`

1. `litigation-paralegal`
2. `family-law-attorney-reviewer`
3. `opposing-counsel`
4. `managing-partner-final-synthesizer`

Conditional additions may still insert:
- `legal-authority-scholar`
- `strategic-options-architect`

### Standard Review

Configured pipeline key: `review_standard`

1. `litigation-paralegal`
2. `trial-fact-checker`
3. `family-law-attorney-reviewer`
4. `opposing-counsel`
5. `neutral-court-reader`
6. `relief-and-order-alignment-counsel`
7. `legal-writing-preservation-editor`
8. `managing-partner-final-synthesizer`

Conditional additions may insert:
- `child-best-interests-family-dynamics-reviewer`
- `financial-support-reviewer`
- `legal-authority-scholar`
- `strategic-options-architect`

### Deep Review

Configured pipeline key: `review_deep`

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

Conditional additions may insert:
- `child-best-interests-family-dynamics-reviewer`
- `financial-support-reviewer`
- `legal-authority-scholar`

### Legal Research Review

Implemented.

Configured pipeline key: `review_legal_research`

1. `legal-authority-scholar`
2. `family-law-attorney-reviewer`
3. `managing-partner-final-synthesizer`

Important limitation: this is still a draft-based review path. It does not integrate a live authority research source by itself.

### Strategy Review

Implemented.

Configured pipeline key: `review_strategy`

1. `chronology-clerk`
2. `trial-fact-checker`
3. `family-law-attorney-reviewer`
4. `opposing-counsel`
5. `neutral-court-reader`
6. `relief-and-order-alignment-counsel`
7. `strategic-options-architect`
8. `practical-resolution-reviewer`
9. `managing-partner-final-synthesizer`

Conditional additions may insert:
- `legal-authority-scholar`
- `child-best-interests-family-dynamics-reviewer`
- `financial-support-reviewer`

### Document-Type Routes

Implemented in `go/cmd/legal-bot/review.go` via `reviewPipelineKey()`.

`motion`
- Resolves to `review_standard`
- Not a separate configured pipeline
- May also trigger conditional `strategic-options-architect` and `legal-authority-scholar`

`opposition`
- Resolves to `review_standard`
- Not a separate configured pipeline
- May also trigger conditional `legal-authority-scholar`

`response`
- Resolves to `review_standard`
- Not a separate configured pipeline
- May also trigger conditional `legal-authority-scholar`

`reply`
- Configured pipeline key: `review_document_reply`
- Sequence:
  1. `trial-fact-checker`
  2. `opposing-counsel`
  3. `neutral-court-reader`
  4. `legal-writing-preservation-editor`
  5. `managing-partner-final-synthesizer`

`declaration`
- Configured pipeline key: `review_document_declaration`
- Sequence:
  1. `litigation-paralegal`
  2. `trial-fact-checker`
  3. `opposing-counsel`
  4. `legal-writing-preservation-editor`
  5. `managing-partner-final-synthesizer`

`proposed-order`
- Configured pipeline key: `review_document_proposed_order`
- Sequence:
  1. `litigation-paralegal`
  2. `relief-and-order-alignment-counsel`
  3. `neutral-court-reader`
  4. `practical-resolution-reviewer`
  5. `managing-partner-final-synthesizer`

`co-parenting-communication`
- Configured pipeline key: `review_document_coparenting_communication`
- Sequence:
  1. `opposing-counsel`
  2. `neutral-court-reader`
  3. `legal-writing-preservation-editor`
  4. `practical-resolution-reviewer`
  5. `managing-partner-final-synthesizer`

Hidden route values that affect conditional routing but are not listed in CLI help:
- `custody`
- `school`
- `therapy`
- `medical`
- `financial`
- `support`
- `order`
- `communication`

## Agent Communication / Data Flow

Active legal agents do not directly call each other. They communicate through:

1. prompt-assembled shared context
2. generated knowledge artifacts on disk
3. accumulated prior review outputs passed inline to later agents

High-level flow:

```text
case/*.md + matter input docs
  -> intake agents
  -> knowledge artifacts in matters/<matter>/knowledge/

case/*.md + selected draft + knowledge artifacts
  -> review agents (mostly independent, same draft + knowledge context)
  -> working/*.md per reviewer
  -> accumulated review findings in-memory during the run
  -> managing-partner-final-synthesizer
  -> output/final_review_packet.md
```

Important details:
- Intake step 2 receives step 1 output directly in the prompt.
- Intake steps 3-5 read the generated knowledge layer, not raw input docs.
- Most review agents receive the same draft plus the same knowledge layer.
- `legal-writing-preservation-editor` receives the draft plus prior review outputs.
- `managing-partner-final-synthesizer` receives case context, draft filename, and all prior review findings, but not the full draft body or raw knowledge layer directly in its configured context mode.
- Review agents do not directly read `working/` or `output/` files during a run.
- There is no explicit structured interchange format between agents beyond markdown/plain-text prompt context and generated files.

## Pipeline Inputs and Outputs

| Step | Reads | Writes | Notes |
|---|---|---|---|
| Intake step 1 `litigation-paralegal` | `case/*.md`, `case/*.txt`, `matters/<matter>/input/*.txt|*.md` | `working/01-intake-mapping.md`, extracted `knowledge/document_register.csv` if parse succeeds | Register extraction is heuristic |
| Intake step 2 `source-document-analyst` | case context, previous step output, raw input docs | `working/02-document-cards.md`, split `knowledge/document_cards/*.md` | Old card files are not deleted first |
| Intake step 3 `chronology-clerk` | case context, knowledge layer | `knowledge/case_timeline.md` | Reads generated knowledge artifacts |
| Intake step 4 `atomic-fact-extractor` | case context, knowledge layer | `knowledge/fact_candidates.md` | No direct raw source context at this step |
| Intake step 5 `attorney-prep-questioner` | case context, knowledge layer | `knowledge/open_questions_for_user.md` | Final intake knowledge artifact |
| Intake report writer | matter metadata | `output/ingestion_report.md` | Written after intake completes |
| Review draft selection | `drafts/` | none | Picks latest `.txt` or `.md` by mod time unless `--draft` supplied |
| Review `draft_plus_knowledge` agents | case context, selected draft, top-level `knowledge/*`, `knowledge/document_cards/*` | `working/*.md` or final output depending on step | Most review steps use this |
| Review `draft_plus_reviews` agents | case context, selected draft, in-memory prior review outputs | `working/*.md` | Used by preservation editor |
| Review `all_reviews` agent | case context, draft filename, all prior review outputs | `output/final_review_packet.md` | Used by final synthesizer |
| Log writer | LLM run metadata | `~/.config/legal-bot/logs/legal-bot.log` | Can be disabled with `LEGAL_BOT_LOG=0` or `VERN_LOG=0` |

## Conditional Routing

Conditional routing is implemented in code, not config, in `go/cmd/legal-bot/review.go`.

| Trigger | Agent(s) Added | Implemented How | File/Path | Notes |
|---|---|---|---|---|
| `--child-related` | `child-best-interests-family-dynamics-reviewer` | CLI flag | `go/cmd/legal-bot/review.go` | Force-adds child reviewer |
| `--financial` | `financial-support-reviewer` | CLI flag | `go/cmd/legal-bot/review.go` | Force-adds financial reviewer |
| Deep mode + child keywords in draft | `child-best-interests-family-dynamics-reviewer` | Draft text heuristic | `go/cmd/legal-bot/review.go` | Only automatic draft-text child detection in deep mode |
| Deep mode + financial keywords in draft | `financial-support-reviewer` | Draft text heuristic | `go/cmd/legal-bot/review.go` | Only automatic draft-text financial detection in deep mode |
| Route values `custody`, `school`, `therapy`, `medical` | `child-best-interests-family-dynamics-reviewer` | Hidden route-name heuristic | `go/cmd/legal-bot/review.go` | Not documented as supported `--document-type` values |
| Route values `financial`, `support` | `financial-support-reviewer` | Hidden route-name heuristic | `go/cmd/legal-bot/review.go` | Not documented as supported `--document-type` values |
| `--mode legal-research` | `legal-authority-scholar` | Base route | `config.default.json`, `go/cmd/legal-bot/review.go` | Full legal-research pipeline |
| Legal keywords in any draft | `legal-authority-scholar` | Draft text heuristic | `go/cmd/legal-bot/review.go` | Can affect quick/standard/deep/document routes |
| Route values `motion`, `opposition`, `response`, `reply`, `declaration`, `proposed-order`, `order`, `communication`, `custody`, `school`, `therapy`, `medical`, `financial`, `support` | `legal-authority-scholar` | Route-name heuristic | `go/cmd/legal-bot/review.go` | Broader than docs usually suggest |
| `--mode strategy` | `strategic-options-architect` | Base route | `config.default.json`, `go/cmd/legal-bot/review.go` | Full strategy pipeline |
| Draft mentions status quo/TRO/emergency/etc. | `strategic-options-architect` | Draft text heuristic | `go/cmd/legal-bot/review.go` | Can affect non-strategy modes too |
| Route values `motion`, `declaration`, `custody`, `school`, `therapy`, `medical` | `strategic-options-architect` | Route-name heuristic | `go/cmd/legal-bot/review.go` | Broad conditional strategy routing |
| `practical-resolution-reviewer` | none conditionally | Configured directly in route | `config.default.json` | Not conditionally inserted in code |

Insertion placement:
- child reviewer: before `legal-writing-preservation-editor`, else before final synthesizer
- financial reviewer: before `legal-writing-preservation-editor`, else before final synthesizer
- legal-authority-scholar: after `trial-fact-checker`, else before final synthesizer
- strategic-options-architect: before `practical-resolution-reviewer`, else before final synthesizer

## Model Profile Resolution

Exact resolution flow:

1. Pipeline step declares `model_profile` and fallback `llm` in `config.default.json`.
2. `go/internal/config/config.go` resolves the profile to a preferred provider/model string and fallback.
3. `runnableEngine()` strips provider-qualified values to the provider prefix.
4. Only `claude`, `gemini`, `codex`, and `copilot` are recognized as runnable engines today.
5. Provider-qualified values like `openai:gpt-5.5` therefore fall back to the configured fallback engine, currently usually `claude`.
6. `--single-llm` or `--llm-mode` overrides can still replace this.

`model_profile` used by active pipelines comes from `config.default.json` pipeline steps, not from agent frontmatter.

| model_profile | Configured Provider/Model | Fallback | Used By | Notes |
|---|---|---|---|---|
| `intake_long_context` | `gemini` | `claude` | Intake mapping | Directly runnable today |
| `structured_extraction` | `openai:gpt-5.4-mini` | `claude` | Source cards, chronology, fact extraction | Falls back to `claude` under current runner |
| `review_reasoning` | `openai:gpt-5.5` | `claude` | Most review agents | Falls back to `claude` under current runner |
| `final_synthesis` | `openai:gpt-5.5` | `claude` | Managing partner | Falls back to `claude` under current runner |
| `cheap_fast` | `openai:gpt-5.4-mini` | `claude` | Not heavily used in active legal pipelines | Present in config |
| `strategy_reasoning` | `openai:gpt-5.5` | `claude` | Strategic options architect | Falls back to `claude` under current runner |
| `legal_research` | `openai:gpt-5.5` | `claude` | Legal-authority-scholar | Falls back to `claude` under current runner |

Other relevant notes:
- Agent frontmatter `model_profile` is parsed by `go/internal/persona/loader.go` but is not the active pipeline routing source.
- Agent frontmatter `model` still exists as metadata and for legacy/helper code paths.
- Supported CLI runners in `go/internal/llm/runner.go`: `claude`, `gemini`, `codex`, `copilot`.
- `config.default.json` currently marks `codex` and `copilot` disabled under `llms`, though runner support exists.

## Agent Loading and Aliases

Active agent files live in `agents/*.md`.

Primary active prompt path:
- `go/internal/llm/runner.go`
- `loadPersonaContext(agentsDir, persona)`
- reads `agents/<persona>.md` from disk

Frontmatter parsing and alias support:
- `go/internal/persona/loader.go`
- parses `name`, `description`, `model`, `model_profile`, `color`
- resolves legacy names with `ResolveName()`

Active alias map includes:
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
- The alias loader is real and tested.
- Active intake/review prompt persona injection does not call `persona.Load()` first; it expects the current persona filename to exist on disk.
- Embedded fallback exists in loader code, but not as the primary active prompt path for intake/review.

## Embedded Assets / Generation

Embedded assets live under `go/internal/embedded/`.

Embedded files:
- `agents/*.md`
- `config.default.json`

Generator files:
- `go/internal/embedded/embedded.go`
- `go/internal/embedded/gen.go`
- generated output: `go/internal/embedded/agents_generated.go`

Regeneration command: found.

```powershell
cd go
go generate ./internal/embedded/
```

What requires regeneration:
- changing `agents/*.md`
- changing `config.default.json`

What does not require regeneration:
- changes to `case/*.md`
- changes to `docs/*.md`
- changes to matter files under `matters/`

Important limitation:
- embedded assets exist and are tested, but active persona prompt loading during legal intake/review still reads from disk `agents/*.md`.

## Active vs Legacy Pipeline Code

| Path | Purpose | Active Legal Routing? | Legacy/Compatibility? | Notes |
|---|---|---|---|---|
| `go/cmd/legal-bot/intake.go` | Intake CLI | Yes | No | Active |
| `go/cmd/legal-bot/review.go` | Review CLI | Yes | No | Active |
| `config.default.json` | Legal pipeline definitions | Yes | No | Active |
| `go/internal/config/config.go` | Config and model resolution | Yes | No | Active |
| `go/internal/llm/runner.go` | LLM subprocess runner | Yes | No | Active |
| `go/internal/persona/loader.go` | Frontmatter parse and aliases | Yes, partly | Compatibility shim | Alias map still matters |
| `go/internal/embedded/*` | Embedded config and agents | Partly | Compatibility/build support | Real, but not the main active prompt source |
| `go/internal/council/*` | Council selection logic | No | Legacy but still present | Mentioned in migration docs |
| `go/internal/pipeline/legalbot_council.go` | Legal-bot council pipeline | No | Legacy but still callable | Not the active intake/review routing contract |
| `go/internal/pipeline/oracle.go` | Oracle flow | No | Legacy but still callable | Legacy terminology remains |
| `go/internal/pipeline/historian.go` | Historian pipeline | No | Legacy but still callable | Separate from active legal routing |
| `go/internal/pipeline/pipeline.go` | Generic pipeline scaffolding | No direct active role | Legacy/shared | Repository infrastructure |
| `go/cmd/legal-bot/historian.go` | Historian CLI command | No | Legacy but still callable | Separate workflow |
| `go/cmd/legal-bot/tui.go` | TUI CLI command | No | Legacy but still callable | Separate workflow |
| `go/internal/tui/*` | Discovery/VernHole/oracle/TUI UI | No | Legacy compatibility | Still in repo |
| `go/internal/generate/*` | Generic persona generation | No | Legacy/shared | Not active legal review routing |
| `go/internal/vts/*` | VTS support | No | Unknown/legacy | Not part of active legal routing |

Classification summary:
- active legal routing: `intake`, `review`, config, LLM runner, case/knowledge IO
- compatibility shim: persona alias map, embedded assets
- legacy but still callable: historian, TUI, pipeline/oracle/council flows
- historical/dead code: none proven dead without deeper execution tracing; some legacy paths remain build/test visible

## How To Run Legal-Bot Today

### Setup Requirements

- Run from the repository checkout.
- Build from `go/`.
- Go 1.25.7+ is required by `go/go.mod`.
- At least one supported LLM CLI must be installed and configured externally:
  - `claude`
  - `gemini`
  - `codex`
  - `copilot`

Practical current path:
- `claude` and `gemini` are the mature local runner targets
- provider-qualified OpenAI names in config do not execute directly today; they fall back

### Verified CLI Commands

```powershell
legal-bot intake <matter-name>
legal-bot intake <matter-name> --single-llm claude
legal-bot intake <matter-name> --llm-mode single_llm

legal-bot review <matter-name>
legal-bot review <matter-name> --draft path\to\draft.md
legal-bot review <matter-name> --mode quick
legal-bot review <matter-name> --mode standard
legal-bot review <matter-name> --mode deep
legal-bot review <matter-name> --mode legal-research
legal-bot review <matter-name> --mode strategy
legal-bot review <matter-name> --document-type motion
legal-bot review <matter-name> --document-type opposition
legal-bot review <matter-name> --document-type response
legal-bot review <matter-name> --document-type reply
legal-bot review <matter-name> --document-type declaration
legal-bot review <matter-name> --document-type proposed-order
legal-bot review <matter-name> --document-type co-parenting-communication
legal-bot review <matter-name> --child-related
legal-bot review <matter-name> --financial
legal-bot review <matter-name> --single-llm claude
legal-bot review <matter-name> --llm-mode single_llm

legal-bot run <llm> <prompt>
legal-bot setup
legal-bot historian <directory>
legal-bot tui
```

### Verified Windows Wrapper Scripts

```powershell
scripts\run-intake.bat <matter-name>
scripts\run-review.bat <matter-name>
scripts\open-latest-output.bat <matter-name>
```

Wrapper behavior:
- `run-intake.bat` and `run-review.bat` default the matter name to `TRO-MTE-GAL-SM` if none is passed
- wrappers build `go\bin\legal-bot.exe` if missing
- wrappers do not expose `--mode`, `--document-type`, `--draft`, `--child-related`, or `--financial`
- `open-latest-output.bat` opens `matters\<matter>\output\final_review_packet.md`

### Environment Variables Found

- `LEGAL_BOT_ROOT`
- `VERN_ROOT`
- `LEGAL_BOT_WORKING_DIR`
- `VERN_WORKING_DIR`
- `LEGAL_BOT_TIMEOUT`
- `VERN_TIMEOUT`
- `LEGAL_BOT_LOG`
- `VERN_LOG`

No repo-local API keys or provider auth flow was found in code; provider setup is expected to exist in the installed external CLI tools.

## Storage Map: Case-Wide vs Matter-Specific

| Information Type | Where It Belongs | Loaded Automatically? | Generated or User-Edited? | Notes |
|---|---|---|---|---|
| Case-wide people/roles | `case/people_and_roles.md` | Yes | User-edited | Global to all matters in this repo |
| Case-wide source hierarchy | `case/source_hierarchy.md` | Yes | User-edited | Global to all matters |
| Preferred framing | `case/preferred_framing.md` | Yes | User-edited | Global to all matters |
| Language calibration | `case/language_calibration.md` | Yes | User-edited | Global to all matters |
| Relief library | `case/relief_library.md` | Yes | User-edited | Global to all matters |
| Legal authority guidance | Optional `case/*.md` file if added | Yes, if root `.md`/`.txt` | User-edited | No dedicated loader; any root case markdown gets loaded |
| Strategy principles | Optional `case/*.md` file if added | Yes, if root `.md`/`.txt` | User-edited | Same generic loader behavior |
| Source documents | `matters/<matter>/input/` | Intake only | User-edited | `.txt` and `.md` only |
| Motion-specific drafts | `matters/<matter>/drafts/` | Review only | User-edited | Latest `.txt` / `.md` chosen unless `--draft` specified |
| Proposed orders | `matters/<matter>/drafts/` | Review only | User-edited | Use filename conventions; no separate folder |
| Declarations | `matters/<matter>/drafts/` | Review only | User-edited | Use document-type routing |
| Exhibits | Usually `matters/<matter>/input/` after text conversion | Intake only | User-edited | No separate exhibit folder convention in code |
| Attorney notes | No dedicated supported folder | No | User-edited if stored manually | Could be added to `input/` or `drafts/`, but no formal separation |
| Matter-specific legal research | No dedicated supported folder | Only if placed in loaded dirs | User-edited/manual | No explicit route for motion-specific research artifacts |
| Generated document cards | `matters/<matter>/knowledge/document_cards/` | Yes | Generated | Reused by review |
| Generated fact candidates | `matters/<matter>/knowledge/fact_candidates.md` | Yes | Generated | Reused by review |
| Generated timeline | `matters/<matter>/knowledge/case_timeline.md` | Yes | Generated | Reused by review |
| Review findings | `matters/<matter>/working/*.md` | No, not automatically reread from disk | Generated | Prior findings are passed in-memory during the active run |
| Final packet | `matters/<matter>/output/final_review_packet.md` | No | Generated | Final operator deliverable |

## Matter Folder Lifecycle

Current matter lifecycle:

```text
matters/<matter>/
  input/     -> user source docs for intake
  drafts/    -> user drafts for review
  knowledge/ -> generated durable knowledge layer
  working/   -> generated intermediate step outputs
  output/    -> generated final packet and intake report
```

Meaning of each folder:
- `input/`: source documents for the case or sub-project; user-edited
- `drafts/`: draft pleadings, declarations, orders, or communications; user-edited
- `knowledge/`: generated reusable artifacts used by later review runs; partly durable
- `working/`: per-step intermediate markdown outputs; generated
- `output/`: user-facing outputs like final review packet; generated

What the repo does not support clearly today:
- per-motion subfolders under a matter with separate routing metadata
- per-matter override copies of `case/*.md`
- matter-specific config files
- automatic separation between global case knowledge and one motion's private strategy or research packet

## Where To Store Case-Wide Knowledge

Current practical guidance:

- durable source documents: `matters/<matter>/input/`, not `case/`
- case-level guidance used across all matters in this repo: `case/*.md`
- people/roles: `case/people_and_roles.md`
- source hierarchy: `case/source_hierarchy.md`
- general litigation posture and framing: `case/preferred_framing.md`
- language strength/risk guidance: `case/language_calibration.md`
- global relief theories and order mechanics: `case/relief_library.md`
- legal authority guidance: no dedicated file today, but any root `case/*.md` would be loaded automatically
- strategy principles: same as above
- user notes/recollections: not clearly modeled; if factual source material, current repo expects them in `matters/<matter>/input/`
- AI-generated knowledge artifacts: `matters/<matter>/knowledge/`

What should never go in `case/`:
- raw source documents
- privileged documents that should not affect every matter in the repo
- matter-specific drafts
- large text dumps that you do not want injected into every intake/review prompt

Reality check:
- `case/` is global to the whole repository, not copied per matter
- `case/` is actually loaded into prompts, not just documentation
- `case/README.md` is also loaded into prompts because the loader reads every root `case/*.md`

## Where To Store Motion-Specific Knowledge

Current repo expectations are loose. The safest interpretation today is:

- motion-specific source documents: `matters/<motion-matter>/input/` or the broader matter's `input/`
- motion drafts: `matters/<matter>/drafts/`
- proposed orders: `matters/<matter>/drafts/`
- declarations: `matters/<matter>/drafts/`
- exhibits: `matters/<matter>/input/` after converting to `.txt` or `.md`
- attorney notes: no dedicated supported location; keep outside auto-loaded folders unless intentionally part of intake or review
- motion-specific legal research: no dedicated supported location; if stored in `knowledge/`, it may be auto-loaded broadly
- motion-specific strategy notes: no dedicated supported location; avoid dropping them casually into `knowledge/` unless you want them fed to review prompts
- generated review outputs: `matters/<matter>/working/` and `matters/<matter>/output/`
- final attorney handoff: `matters/<matter>/output/final_review_packet.md`

If one case has multiple large motion tracks, the current repo most naturally expects:
- one matter per motion/project, or
- disciplined file naming inside one matter

There is no built-in motion subfolder pattern in code.

## Output Locations

Intake outputs:
- `matters/<matter>/working/01-intake-mapping.md`
- `matters/<matter>/working/02-document-cards.md`
- `matters/<matter>/knowledge/document_register.csv`
- `matters/<matter>/knowledge/document_cards/*.md`
- `matters/<matter>/knowledge/case_timeline.md`
- `matters/<matter>/knowledge/fact_candidates.md`
- `matters/<matter>/knowledge/open_questions_for_user.md`
- `matters/<matter>/output/ingestion_report.md`

Review outputs:
- intermediate reviewer outputs in `matters/<matter>/working/*.md`
- final packet at `matters/<matter>/output/final_review_packet.md`

Logs:
- `~/.config/legal-bot/logs/legal-bot.log`

Generated embedded file:
- `go/internal/embedded/agents_generated.go`

Overwrite behavior:
- output filenames are fixed
- reruns overwrite prior files with the same name
- there is no timestamped archival behavior in active intake/review code

Finding the latest output:
- `scripts/open-latest-output.bat` opens the fixed file `matters\<matter>\output\final_review_packet.md`
- despite its name, it does not search a timestamped history

Gitignore behavior:
- `matters/*/input/` and `matters/*/drafts/` are gitignored
- `case/private/` is gitignored
- generated `knowledge/`, `working/`, and `output/` are not ignored by default

## Prompt Context Assembly

| Agent/Pipeline Step | Receives Case Guidance? | Receives Knowledge Artifacts? | Receives Draft? | Receives Prior Outputs? | Writes |
|---|---|---|---|---|---|
| Intake step 1 `litigation-paralegal` | Yes, all root `case/*.md|*.txt` | No | No | No | `working/01-intake-mapping.md`, register extraction |
| Intake step 2 `source-document-analyst` | Yes | No direct knowledge layer; receives prior step output and raw input docs | No | Yes, previous intake step output | `working/02-document-cards.md`, `knowledge/document_cards/*.md` |
| Intake step 3 `chronology-clerk` | Yes | Yes | No | No | `knowledge/case_timeline.md` |
| Intake step 4 `atomic-fact-extractor` | Yes | Yes | No | No | `knowledge/fact_candidates.md` |
| Intake step 5 `attorney-prep-questioner` | Yes | Yes | No | No | `knowledge/open_questions_for_user.md` |
| Review `draft_plus_knowledge` agents: paralegal, fact checker, family-law reviewer, opposing counsel, neutral court reader, relief alignment, child reviewer, financial reviewer, legal-authority conditional, strategic-options conditional | Yes | Yes, top-level `knowledge/*` plus `knowledge/document_cards/*` | Yes | No | `working/*.md` unless final |
| `legal-writing-preservation-editor` | Yes | No direct knowledge layer in prompt | Yes | Yes, accumulated previous review findings | `working/*.md` |
| `managing-partner-final-synthesizer` | Yes | No direct knowledge layer in prompt | No draft body; draft filename only | Yes, all prior review findings | `output/final_review_packet.md` |

Important assembly rules:
- `case/*.md` and `case/*.txt` root files are concatenated into every intake/review prompt
- `knowledge/` top-level files are concatenated into knowledge context regardless of extension
- `knowledge/document_cards/*` are concatenated into knowledge context
- raw source documents are only directly included during intake, not review
- `language_calibration.md`, `source_hierarchy.md`, `preferred_framing.md`, `people_and_roles.md`, `relief_library.md`, and `case/README.md` are all included automatically because they sit in `case/`

## Customization Map

| Customization Goal | Where To Change It | Risk Level | Requires Tests? | Requires Embedded Regeneration? | Notes |
|---|---|---|---|---|---|
| Change language guidance | `case/language_calibration.md` | Low | Recommended | No | Affects every matter and every intake/review prompt |
| Change preferred framing | `case/preferred_framing.md` | Low | Recommended | No | Global to repo |
| Change relief library | `case/relief_library.md` | Low | Recommended | No | Global to repo |
| Add a person/provider | `case/people_and_roles.md` | Low | Recommended | No | Global to repo |
| Change source hierarchy | `case/source_hierarchy.md` | Medium | Recommended | No | May justify rerunning intake |
| Add optional case-wide guidance | new root `case/*.md` or `.txt` | Medium | Recommended | No | Will be auto-loaded into all prompts |
| Change review mode sequences | `config.default.json` | Medium | Yes | Yes | Config is embedded |
| Change agent order | `config.default.json` | Medium | Yes | Yes | Review code also has insertion logic to consider |
| Change conditional routing | `go/cmd/legal-bot/review.go` | High | Yes | No | Code-driven |
| Change model provider/profile mapping | `config.default.json` | Medium | Yes | Yes | Current runner only supports certain engines |
| Rename an agent | `agents/*.md`, alias map, config, docs/tests | High | Yes | Yes | Disk filename matters for active prompts |
| Add a new agent | `agents/*.md`, `config.default.json`, maybe alias/tests/docs | High | Yes | Yes | Also review prompt contexts |
| Change final output format | agent prompt markdown for `managing-partner-final-synthesizer` and possibly docs | Medium | Recommended | Yes if agent file changes | Final packet shape is largely prompt-defined |

## Implemented vs Documented Features

| Feature | Documented? | Implemented? | File/Path | Notes |
|---|---|---|---|---|
| legal-research mode | Yes | Yes | `config.default.json`, `go/cmd/legal-bot/review.go` | Draft-based review mode, not standalone research integration |
| strategy mode | Yes | Yes | `config.default.json`, `go/cmd/legal-bot/review.go` | Draft-based strategy review |
| quick review | Yes | Yes | `config.default.json`, `go/cmd/legal-bot/review.go` | Active |
| standard review | Yes | Yes | `config.default.json`, `go/cmd/legal-bot/review.go` | Active default |
| deep review | Yes | Yes | `config.default.json`, `go/cmd/legal-bot/review.go` | Active |
| `--child-related` | Yes | Yes | `go/cmd/legal-bot/review.go` | Active CLI flag |
| `--financial` | Yes | Yes | `go/cmd/legal-bot/review.go` | Active CLI flag |
| `--document-type` | Yes | Yes | `go/cmd/legal-bot/review.go` | Active, but only certain route names are advertised |
| proposed-order route | Yes | Yes | `config.default.json`, `go/cmd/legal-bot/review.go` | Active |
| declaration route | Yes | Yes | `config.default.json`, `go/cmd/legal-bot/review.go` | Active |
| reply route | Yes | Yes | `config.default.json`, `go/cmd/legal-bot/review.go` | Active |
| co-parenting route | Yes | Yes | `config.default.json`, `go/cmd/legal-bot/review.go` | Active under `co-parenting-communication` |
| `language_calibration.md` loading | Yes | Yes | `go/cmd/legal-bot/intake.go` `loadCaseContext()` | Loaded generically because all root `case/*.md` are loaded |
| optional case files | Partly | Yes, generically | `go/cmd/legal-bot/intake.go` `loadCaseContext()` | Any root `case/*.md|*.txt` auto-loads; no dedicated whitelist |
| `model_profile` routing | Yes | Yes, partly | `config.default.json`, `go/internal/config/config.go`, CLI code | Pipeline step `model_profile` is active; agent frontmatter `model_profile` is not the active route source |
| OpenAI/Codex model support | Yes | Partial | `go/internal/llm/runner.go`, `go/internal/config/config.go` | Codex runner exists; provider-qualified OpenAI profile targets do not run directly today |
| embedded regeneration | Yes | Yes | `go/internal/embedded/*` | Command found |

## Minimum Working Example

1. Setup requirements
   - Use a checkout of this repo.
   - Have Go 1.25.7+ available.
   - Have at least one supported LLM CLI installed and authenticated, preferably `claude` or `gemini`.

2. Create or choose a matter folder
   - Existing example matter: `matters/TRO-MTE-GAL-SM/`
   - Or create:

```text
matters/<matter-name>/
  input/
  drafts/
```

3. Add one source file
   - Put a `.txt` or `.md` file into `matters/<matter-name>/input/`

4. Run intake

```powershell
legal-bot intake <matter-name>
```

or

```powershell
scripts\run-intake.bat <matter-name>
```

5. Add one draft
   - Put a `.txt` or `.md` draft into `matters/<matter-name>/drafts/`

6. Run review

```powershell
legal-bot review <matter-name>
```

or

```powershell
scripts\run-review.bat <matter-name>
```

7. Find output
   - Final packet:

```text
matters/<matter-name>/output/final_review_packet.md
```

Expected artifacts after a successful first run:
- `knowledge/document_register.csv`
- `knowledge/document_cards/*.md`
- `knowledge/case_timeline.md`
- `knowledge/fact_candidates.md`
- `knowledge/open_questions_for_user.md`
- `output/ingestion_report.md`
- `output/final_review_packet.md`

If this does not work from a fresh clone, the most likely missing pieces are:
- Go 1.25.7+ not available
- LLM CLI not installed or not authenticated
- no `.txt` or `.md` source/draft files present

## Common Failure Modes

- missing LLM CLI
  - `claude`, `gemini`, `codex`, or `copilot` not installed or not configured
- wrong Go version
  - `go/go.mod` requires Go 1.25.7+
- wrong working directory
  - building/running from the wrong place can fail or miss repo paths
- unsupported file type
  - intake and review only read `.txt` and `.md`
- missing matter folder
  - CLI expects `matters/<matter>/...`
- no drafts found
  - review requires at least one `.txt` or `.md` draft
- missing `document_register.csv`
  - review refuses to start without it
- stale embedded agents/config
  - after editing `agents/*.md` or `config.default.json`, embedded data may be outdated even though active prompt loading still uses disk agents
- config references missing agent
  - step personas must match real disk agent files for active prompt context
- renamed file not loaded
  - case loader is generic; root `case/*.md` auto-loads, subfolders do not
- docs mention a workflow the wrappers do not expose
  - scripts only cover the default path
- output path unclear
  - outputs overwrite in place; there is no timestamp history
- old artifacts reused unexpectedly
  - `knowledge/` and `document_cards/` are reused and not cleaned fully
- source documents too large
  - large raw prompt contexts can become expensive or hard for the model to handle well
- model provider unavailable
  - provider-qualified OpenAI profile targets fall back today; direct OpenAI runner not found
- accidental prompt pollution
  - extra `.md/.txt` files in `case/` or extra files in `knowledge/` top level will be included automatically

## Documentation Mismatches

| Doc/File | Claim | Actual Behavior | Recommended Fix |
|---|---|---|---|
| `README.md` | Legal research and strategy support are presented at a high level | Implemented, but both are still draft-based review modes without live authority integration | Make the draft-based limitation explicit wherever legal research is described |
| `README.md` and `docs/user_workflow.md` | Wrapper scripts are shown as the easy path | True, but wrappers only expose default review and intake | Add a stronger note that alternate modes and flags require direct CLI usage |
| `docs/architecture.md` | Embedded agents/config are described as part of the standalone binary story | Embedded assets exist, but active persona prompt loading during intake/review still reads `agents/*.md` from disk | Document the disk dependency more explicitly |
| `docs/architecture.md` and `case/README.md` | `case/` is described as a control layer | True, but docs do not strongly emphasize that `case/README.md` and any extra root `case/*.md|*.txt` are also auto-loaded | Add an explicit warning about prompt pollution from extra case files |
| `docs/architecture.md`, `docs/user_workflow.md`, `README.md` | Document-type routing is described cleanly | True for supported values, but code also reacts to hidden route-name values like `custody`, `school`, `medical`, `financial`, `support` | Either document these hidden route triggers or stop implying `--document-type` is tightly validated |
| `README.md`, `docs/token_strategy.md`, `docs/agent_persona_system.md` | `model_profile` is documented as provider-neutral | True at config level, but current runtime support still falls back from provider-qualified OpenAI values to local supported engines | Keep calling this out clearly as partial implementation |
| `scripts/open-latest-output.bat` references in docs | Name implies latest output browsing | Script opens a fixed file path, not a timestamped latest artifact search | Either rename script later or document the fixed-path behavior explicitly |

## Risks / Unknowns

- No live authority-search integration was found for `legal-authority-scholar`; legal authority quality depends entirely on supplied prompt context and model behavior.
- Embedded assets are real and tested, but active prompt persona injection still appears disk-dependent for legal intake/review.
- Hidden route keywords influence conditional routing without clear user-facing validation.
- `knowledge/` prompt loading is broad and may accidentally include manually dropped files or stale generated files.
- No strong per-motion storage model exists; multi-motion matters can get cluttered.
- Some legacy packages remain callable, but a full execution audit was not performed for every non-legal command path.

## Recommended Follow-Up Tasks

### Docs-Only

- Explicitly document that `legal-research` and `strategy` are draft-based review modes, not standalone research consoles.
- Add a prominent note that all root `case/*.md|*.txt` files, including `case/README.md`, are auto-loaded.
- Clarify that `open-latest-output.bat` opens a fixed final packet path.
- Clarify that wrapper scripts only expose the default review path.

### Config

- Reconcile model profile documentation with current runtime support, especially for provider-qualified OpenAI targets.
- Consider tightening or documenting hidden route-name triggers.

### Code

- Decide whether active intake/review should truly support embedded persona fallback instead of requiring on-disk agent files.
- Consider whitelisting case control files instead of loading every root `case/*.md|*.txt`.
- Consider filtering `knowledge/` top-level files by an expected allowlist.
- Consider cleaning old `knowledge/document_cards/*.md` before rewriting them.
- Consider timestamped outputs or archived runs.

### Tests

- Add smoke coverage for conditional routing insertion behavior.
- Add tests covering broad `case/` auto-load and `knowledge/` auto-load behavior.
- Add tests or validation around wrapper-script assumptions versus CLI capabilities.

### Future Design

- Add a supported matter-specific override layer for framing, people, relief, and legal-authority guidance.
- Add a dedicated motion/project subfolder convention if multiple filings must coexist cleanly inside one case.
- Add structured legal research source integration if `legal-authority-scholar` is meant to be citation-strong in practice.

## Commands Run

Inspection commands used during this investigation included:

```powershell
Get-ChildItem -LiteralPath docs -Filter *.md | Sort-Object Name | Select-Object -ExpandProperty Name
Get-ChildItem -LiteralPath case -Filter *.md | Sort-Object Name | Select-Object -ExpandProperty Name
Get-ChildItem -LiteralPath agents -Filter *.md | Sort-Object Name | Select-Object -ExpandProperty Name
Get-ChildItem -LiteralPath agents\shared -Filter *.md | Sort-Object Name | Select-Object -ExpandProperty Name
Get-ChildItem -LiteralPath scripts -Force | Select-Object Name,Mode
Get-ChildItem -Recurse -LiteralPath matters -Depth 2 -Force | Select-Object FullName,Mode | Select-Object -First 100
Get-ChildItem -Recurse -LiteralPath go\internal\pipeline -File | Select-Object -ExpandProperty FullName
Get-ChildItem -Recurse -LiteralPath go\internal\council -File | Select-Object -ExpandProperty FullName
Get-ChildItem -Recurse -LiteralPath go\internal\tui -File | Select-Object -ExpandProperty FullName
Get-ChildItem -Recurse -LiteralPath go\internal\generate -File | Select-Object -ExpandProperty FullName

Get-Content -Raw README.md
Get-Content -Raw config.default.json
Get-Content -Raw scripts\run-intake.bat
Get-Content -Raw scripts\run-review.bat
Get-Content -Raw scripts\open-latest-output.bat
Get-Content -Raw go\cmd\legal-bot\main.go
Get-Content -Raw go\cmd\legal-bot\intake.go
Get-Content -Raw go\cmd\legal-bot\review.go
Get-Content -Raw go\cmd\legal-bot\run.go
Get-Content -Raw go\cmd\legal-bot\setup.go
Get-Content -First 80 go\cmd\legal-bot\historian.go
Get-Content -First 60 go\cmd\legal-bot\tui.go
Get-Content -Raw go\internal\config\config.go
Get-Content -Raw go\internal\llm\runner.go
Get-Content -Raw go\internal\llm\log.go
Get-Content -Raw go\internal\persona\loader.go
Get-Content -Raw go\internal\embedded\embedded.go
Get-Content -Raw go\internal\embedded\gen.go
Get-Content -Raw go\internal\config\config_test.go
Get-Content -Raw go\internal\persona\loader_test.go
Get-Content -Raw go\internal\embedded\embedded_test.go
Get-Content -Raw docs\architecture.md
Get-Content -Raw docs\user_workflow.md
Get-Content -Raw .gitignore

rg -n "pipeline|route|routes|mode|quick|standard|deep|legal-research|strategy|document-type|child-related|financial|final|synth" .
rg -n "litigation-paralegal|source-document-analyst|atomic-fact-extractor|chronology-clerk|attorney-prep-questioner|legal-authority-scholar|trial-fact-checker|family-law-attorney-reviewer|opposing-counsel|neutral-court-reader|relief-and-order-alignment-counsel|legal-writing-preservation-editor|child-best-interests-family-dynamics-reviewer|strategic-options-architect|practical-resolution-reviewer|financial-support-reviewer|managing-partner-final-synthesizer|bulldog-advocate|attack-surface-reviewer|fact-auditor|final-synthesizer" .
rg -n "case/|source_hierarchy.md|people_and_roles.md|preferred_framing.md|language_calibration.md|avoid_language.md|relief_library.md|legal_authority_guidance.md|strategy_principles.md" .
rg -n "matters/|input|drafts|knowledge|working|output|document_register|document_cards|case_timeline|fact_candidates|open_questions|final_review_packet" .
rg -n "model_profile|model|provider|claude|gemini|openai|codex|copilot|fallback" .
rg -n "embedded|agents_generated.go|generate|alias|aliases|Vern|VernHole|council|oracle|discovery|TUI" .
```

Build/test attempts:

```powershell
cd go
$env:GOTOOLCHAIN='local'
C:\GIT\legal-bot\.tools\go1.22.12\go\bin\go.exe build ./...
C:\GIT\legal-bot\.tools\go1.22.12\go\bin\go.exe test ./...
```

Result:
- build/test could not complete with the available local portable toolchain because `go/go.mod` requires Go `1.25.7`, while the available portable toolchain was `1.22.12`
- no application behavior was modified during this investigation
