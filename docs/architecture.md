# Architecture

## Overview

Legal-Bot is a local-first CLI that orchestrates LLM calls through configurable legal pipelines. The current product has three main layers:

1. matter intake and knowledge building
2. workflow-oriented legal analysis
3. draft review

The repo still contains legacy discovery, council, TUI, and compatibility surfaces, but those are not the primary contract for the current legal product.

## High-Level Shape

```text
+--------------------------------------------------+
| CLI                                              |
| guide | intake | workflow | review | routes      |
+--------------------------------------------------+
| Legal Routing                                    |
| intake pipelines | workflow pipelines | review   |
+--------------------------------------------------+
| Persona Runner                                   |
| agents/*.md + config.default.json + llm.Run      |
+--------------------------------------------------+
| File System                                      |
| case/ | matters/ | agents/ | docs/               |
+--------------------------------------------------+
| Legacy Compatibility Surfaces                    |
| council | pipeline | tui | generate | vts        |
+--------------------------------------------------+
```

## Case Control Layer

`case/` is global prompt context.

Current loader behavior:

- every root `.md` and `.txt` file in `case/` is auto-loaded into prompts
- this affects intake, workflow, and review runs

Implications:

- `case/` should hold durable guidance, not raw evidence dumps
- large files in `case/` bloat every run
- matter-specific sources belong in `matters/<matter>/input/`

## Matter Layer

Per matter, the current system uses:

```text
matters/<matter>/
  input/
  drafts/
  knowledge/
  working/
  output/
  situation.md
  goal.md
  questions.md
```

Use:

- `input/` for source materials
- `drafts/` for review drafts
- `knowledge/` for intake outputs
- `working/` for intermediate step outputs
- `output/` for final workflow and review packets
- `situation.md`, `goal.md`, `questions.md` for draftless workflows

## Current Commands

The active legal workflow contract is:

- `legal-bot guide`
- `legal-bot intake`
- `legal-bot review`
- `legal-bot workflow`
- `legal-bot routes`

Utility commands still exposed:

- `legal-bot run`
- `legal-bot setup`
- `legal-bot historian`
- `legal-bot tui`

Legacy Vern discovery/council/oracle/generate wrappers are not part of that current contract.

### Guide

Implemented in:

- `go/cmd/legal-bot/guide.go`

Behavior:

- presents a deterministic guided menu
- maps practical user jobs onto existing `intake`, `workflow`, and `review` routes
- asks only route-specific follow-up questions
- prints an execution plan before running
- can dry-run without invoking any LLM pipeline
- stores guided request trace files under `matters/<matter>/working/`

### Intake

Implemented in:

- `go/cmd/legal-bot/intake.go`

Behavior:

- loads `case/`
- reads `matters/<matter>/input/`
- runs the `intake` pipeline from `config.default.json`
- writes the knowledge layer to `knowledge/` and `output/`

### Review

Implemented in:

- `go/cmd/legal-bot/review.go`

Behavior:

- requires the knowledge layer
- loads the latest `.md` or `.txt` draft from `drafts/` unless `--draft` is passed
- selects a review pipeline based on `--mode` and `--document-type`
- applies conditional child, financial, legal-authority, and strategy routing

Review route clarifications:

- `motion`, `opposition`, and `response` fold into standard review
- `reply`, `declaration`, `proposed-order`, and `co-parenting-communication` have document-type-specific routes
- `--child-related` and `--financial` can force conditional specialists even when heuristics would not

### Workflow

Implemented in:

- `go/cmd/legal-bot/workflow.go`

Behavior:

- selects a workflow route from `--type`
- loads `case/`
- loads matter-level `situation.md`, `goal.md`, and `questions.md` conventions
- optionally includes a draft when relevant
- uses the knowledge layer if present
- falls back to an input manifest if knowledge is absent
- writes final output to `output/<workflow-type>.md`

Current workflow types:

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

### Routes Diagnostic

Implemented in:

- `go/cmd/legal-bot/routes_command.go`
- `go/cmd/legal-bot/routes_catalog.go`

Behavior:

- prints the active route catalog
- shows review modes and document-type mappings
- shows workflow types and conditional specialists
- shows the legacy surface inventory
- supports `--json` for stable machine-readable output

## Pipeline Configuration

Active pipeline definitions live in:

- [`config.default.json`](/c:/GIT/legal-bot/config.default.json)

The config holds:

- `legal_pipelines`
- `model_profiles`
- `llm_modes`

`model_profiles` now distinguish:

- engine: the executable Legal-Bot runs
- model: the model string passed to that executable
- effort: optional reasoning level for engines that support it

Current legal pipeline groups:

- `intake`
- `review_*`
- workflow routes such as `triage_new_situation`, `brief_counsel`, `triage_special_master`, `draft_external_message`, `build_evidence_packet`, `review_pattern`, `fact_lock`, `authority_check`, `prep_hearing_or_call`, and `journal_entry`

## Context Assembly

Prompt assembly is still code-driven.

Important behavior:

- `case/` context is built in `loadCaseContext`
- knowledge context is built in `buildKnowledgeContext`
- intake raw input context is built in `buildInputContext`
- workflow uses situation, goal, questions, optional draft, knowledge, and input manifest
- workflow now also injects a route-specific output contract summary into all workflow prompts and the full output contract into the final synthesizer prompt
- review uses draft plus knowledge and accumulated prior step outputs

This means the repo is mixed:

- pipeline order is declarative
- routing is partly code-driven
- context assembly is code-driven
- file loading is code-driven
- output filenames are code-driven

Workflow output contracts live in:

- `go/cmd/legal-bot/workflow_contracts.go`

They define workflow-specific final packet structure, recommended-track behavior, source-strength framing, Do Not Chase requirements, and external-language rules without changing pipeline order.

## Conditional Routing

The active conditional control plane is split between:

- `go/cmd/legal-bot/routes_catalog.go` for the named specialist inventory and keyword lists
- `go/cmd/legal-bot/review.go` for review insertion logic
- `go/cmd/legal-bot/workflow.go` for workflow insertion logic

Conditional additions currently focus on:

- child/family dynamics
- financial/support issues
- authority-sensitive work
- strategy-sensitive work

Current conditional specialists:

- `child-best-interests-family-dynamics-reviewer`
- `financial-support-reviewer`
- `legal-authority-scholar`
- `strategic-options-architect`

## Workflow Implementation Status

Implemented now through `legal-bot workflow --type ...` and the shared route catalog:

- triage
- counsel-brief
- sm-triage
- draft-message
- evidence-packet
- pattern-review
- fact-lock
- authority-check
- prep
- journal-entry

## Embedded Assets

Embedded agent and config data live in:

- `go/internal/embedded/agents_generated.go`

After changing `agents/*.md` or `config.default.json`, regenerate with:

```powershell
cd go
go generate ./internal/embedded/
```

## Wrapper Scripts

Current wrappers in `scripts/` still expose only the basic intake/review path:

- `run-intake.bat`
- `run-guide.bat`
- `run-review.bat`
- `open-latest-output.bat`

Legacy compatibility wrappers live in `bin/`. See [Legacy Surfaces](/c:/GIT/legal-bot/docs/legacy_surfaces.md).

## Design Direction

The intended user experience is:

1. user states the job
2. Legal-Bot routes to the right workflow
3. internal personas debate and pressure-test
4. managing partner chooses the path
5. user gets a practical next-step output

The current architecture is now aligned enough to support that direction without breaking the existing intake and review commands.

## Current Model Routing

The active legal pipeline path is:

1. `config.default.json` defines `model_profiles`
2. intake, workflow, and review steps reference those profiles with `model_profile`
3. the command resolves a target chain such as primary plus fallbacks
4. Legal-Bot picks the first target whose engine executable exists on `PATH`
5. `llm.Run` passes the resolved `engine`, `model`, and `effort` to the selected local CLI

Legacy `openai:*` entries are normalized to Codex targets for backward compatibility. They are no longer treated as documentation-only values that silently collapse into Claude.
