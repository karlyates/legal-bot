# Legal-Bot

Legal-Bot is a local-first legal analysis and legal-operations assistant for an informed user. It builds matter knowledge from local source files and runs workflow-oriented AI analysis for triage, escalation decisions, counsel briefing, draft review, evidence review, authority gaps, and external drafting.

> Legal-Bot produces AI-generated legal analysis, strategy recommendations, issue spotting, risk review, evidence analysis, and drafting support for an informed user. It may recommend actions, escalation paths, attorney questions, legal research issues, draft language, and litigation strategy options. Legal-Bot is not licensed counsel, may be wrong or incomplete, and outputs should be verified before being filed, served, sent externally, or relied on in a legal matter.

## Product Shape

Legal-Bot is no longer just a "review this draft" tool.

It is organized around practical user jobs:

1. Matter Intake
2. New Situation Triage
3. Evidence Packet Builder
4. Counsel Brief
5. Draft Review
6. Special Master Triage
7. External Message Drafting
8. Pattern Review
9. Fact Lock / Source Strength Review
10. Authority Check
11. Hearing / Call Prep
12. Post-Event Journal Entry

The internal panel still uses many personas, but the user-facing mental model is:

`What am I trying to do?`

## What It Does

1. **Matter Intake**  
   Builds the durable knowledge layer from local source files.

2. **Workflow Analysis**  
   Handles new situations before a draft exists, including triage, counsel briefing, Special Master triage, evidence-packet planning, pattern review, fact lock, authority check, prep, and journal entry creation.

3. **Draft Review**  
   Reviews a draft against the knowledge layer with quick, standard, deep, document-type, and authority-oriented paths.

## Current CLI

Implemented now:

```powershell
legal-bot guide
legal-bot guide <matter>
legal-bot guide --list
legal-bot guide <matter> --dry-run
legal-bot routes
legal-bot routes --json
legal-bot intake <matter>
legal-bot workflow <matter> --type triage
legal-bot workflow <matter> --type counsel-brief
legal-bot workflow <matter> --type sm-triage
legal-bot workflow <matter> --type draft-message --audience ofw
legal-bot workflow <matter> --type evidence-packet
legal-bot workflow <matter> --type pattern-review
legal-bot workflow <matter> --type fact-lock
legal-bot workflow <matter> --type authority-check
legal-bot workflow <matter> --type prep
legal-bot workflow <matter> --type journal-entry
legal-bot review <matter>
legal-bot review-draft <matter>
legal-bot review <matter> --document-type motion
legal-bot review <matter> --document-type opposition
legal-bot review <matter> --document-type response
```

Still available for power users:

```powershell
legal-bot review <matter> --mode quick
legal-bot review <matter> --mode standard
legal-bot review <matter> --mode deep
legal-bot review <matter> --mode legal-research
legal-bot review <matter> --mode strategy
legal-bot review <matter> --document-type declaration
legal-bot review <matter> --document-type proposed-order
legal-bot review <matter> --document-type co-parenting-communication
legal-bot review <matter> --child-related
legal-bot review <matter> --financial
```

`legal-bot start <matter>` is now available as a low-risk alias for `legal-bot guide <matter>`.

Utility commands still exposed:

```powershell
legal-bot run <llm> <prompt>
legal-bot setup
legal-bot historian <directory>
legal-bot tui
```

Legacy Vern-era wrappers still exist in `bin/`, but they are not the current family-law workflow. See [Legacy Surfaces](/c:/GIT/legal-bot/docs/legacy_surfaces.md).

## Quick Start

### 1. Set Up `case/`

Edit the case control files in [`case/README.md`](/c:/GIT/legal-bot/case/README.md).

Important:

- Every root `.md` and `.txt` file in `case/` is loaded into prompts.
- `case/` is global prompt context, not a raw evidence dump.
- Avoid dropping huge evidence files into `case/`.

### 2. Add Matter Sources

Put matter-specific source files in:

```text
matters/<matter>/input/
```

Current supported source formats for intake:

- `.md`
- `.txt`

### 3. Run Intake

```powershell
legal-bot intake <matter>
```

Intake builds:

- `knowledge/document_register.csv`
- `knowledge/document_cards/`
- `knowledge/case_timeline.md`
- `knowledge/fact_candidates.md`
- `knowledge/open_questions_for_user.md`
- `output/ingestion_report.md`

### 4. Use The Guided Front Door

```powershell
legal-bot guide <matter>
```

The guide asks what Karl is trying to do, gathers only the route-specific details, shows a clear execution plan, and then routes to the existing `intake`, `workflow`, or `review` command logic.

Guided choices:

1. Understand a new situation -> `workflow --type triage`
2. Decide whether to escalate, preserve, or let something go -> `workflow --type triage`
3. Build an evidence packet -> `workflow --type evidence-packet`
4. Write to Kaitlyn/Jessika -> `workflow --type counsel-brief`
5. Review a draft from counsel -> `review`
6. Decide whether to bring something to the Special Master -> `workflow --type sm-triage`
7. Write an OFW / SM / provider / school / counsel message -> `workflow --type draft-message`
8. Analyze whether recurring events add up to a pattern -> `workflow --type pattern-review`
9. Lock down facts before external use -> `workflow --type fact-lock`
10. Check legal authority gaps -> `workflow --type authority-check`
11. Prepare for a call/hearing/SM conference -> `workflow --type prep`
12. Create a journal entry -> `workflow --type journal-entry`
13. Run intake on new source material -> `intake`

The guide is deterministic route selection. It does not make an extra LLM call just to choose the route.

### 4a. Inspect The Active Route Catalog

```powershell
legal-bot routes
legal-bot routes --json
```

This prints the active command surface, review modes, document-type mappings, workflow types, conditional specialists, and legacy wrapper inventory.

### 5. Add Draftless Workflow Inputs Manually If You Want

Advanced `workflow` runs still support matter-level files:

```text
matters/<matter>/
  situation.md
  goal.md
  questions.md
```

Supported fallbacks:

- `matters/<matter>/triage.md`
- `matters/<matter>/input/situation.md`
- `matters/<matter>/input/goal.md`
- `matters/<matter>/input/questions.md`
- `matters/<matter>/working/user_situation.md`

The guide also stores traceability notes under:

```text
matters/<matter>/working/guided_<timestamp>_<workflow-type>.md
```

### 6. Run A Workflow Directly

```powershell
legal-bot workflow <matter> --type triage
```

Useful examples:

```powershell
legal-bot workflow TRO-MTE-GAL-SM --type triage
legal-bot workflow TRO-MTE-GAL-SM --type counsel-brief
legal-bot workflow TRO-MTE-GAL-SM --type sm-triage
legal-bot workflow TRO-MTE-GAL-SM --type draft-message --audience special-master
legal-bot workflow TRO-MTE-GAL-SM --type authority-check --draft matters\TRO-MTE-GAL-SM\drafts\motion.md
```

### 7. Run Draft Review

Put a draft in:

```text
matters/<matter>/drafts/
```

Then run:

```powershell
legal-bot review <matter>
```

Review effort levels:

- `quick` = fast attorney-readiness screen
- `standard` = normal review
- `deep` = high-stakes review

Document type changes emphasis, not the core product:

- `reply`
- `declaration`
- `proposed-order`
- `co-parenting-communication`

## Workflow Inputs and Outputs

The workflow command writes the final result to:

```text
matters/<matter>/output/<workflow-type>.md
```

Examples:

- `output/triage.md`
- `output/counsel-brief.md`
- `output/sm-triage.md`
- `output/draft-message.md`
- `output/evidence-packet.md`
- `output/pattern-review.md`
- `output/fact-lock.md`
- `output/authority-check.md`
- `output/prep.md`
- `output/journal-entry.md`

Workflow outputs now follow explicit output contracts so the final packets are more predictable: each workflow has required sections, source-strength framing, recommended-track behavior, Do Not Chase analysis, and clearer separation between internal strategy and external-safe language. See [Workflow Output Contracts](/c:/GIT/legal-bot/docs/workflow_output_contracts.md).

Intermediate files go into:

```text
matters/<matter>/working/
```

Guided input trace files also live in `working/`. They stay under the matter and are not written into `case/`.

## Guided Examples

```powershell
legal-bot guide Revere-042926
legal-bot guide Revere-042926 --dry-run
legal-bot guide --list
legal-bot guide Revere-042926 --choice 4 --situation "School portal access was cut off again." --goal "Prepare a concise attorney brief."
legal-bot guide TRO-MTE-GAL-SM --intent review --draft matters\TRO-MTE-GAL-SM\drafts\reply.md
```

## Decision Style

Legal-Bot is expected to be useful and opinionated.

It should be able to say:

- pursue now
- use OFW first
- ask counsel first
- bring it to the Special Master
- preserve for pattern
- journal only
- let it go

`Do Not Chase` is a feature, not a failure.

## Source Discipline

Core rules:

- Do not invent legal authority.
- Do not treat user recollection as proven fact.
- Do not treat AI summaries as evidence.
- Do not treat filed allegations as true unless admitted, adopted, or independently supported.
- Do not convert internal strategy into court-facing accusations.
- Do not recommend escalation without considering proportionality, evidence strength, and credibility impact.
- When attorney review is needed, ask a precise attorney question.
- When evidence is missing, identify the exact source needed.

## Project Structure

```text
legal-bot/
  case/                         global case-control prompt context
  matters/
    <matter>/
      input/                   matter-specific source files
      drafts/                  draft files for review
      knowledge/               intake-generated knowledge layer
      working/                 intermediate workflow and review files
      output/                  final workflow and review packets
      situation.md             raw situation description for draftless workflows
      goal.md                  what the user wants
      questions.md             user questions to answer
  agents/                      persona definitions
  docs/                        product and architecture docs
  scripts/                     Windows wrappers
  go/                          Go CLI source
```

## Requirements

- **Go 1.25.7+**
- At least one supported LLM CLI installed:
  - `claude`
  - `gemini`
  - `codex`
  - `copilot`

## Model Routing

Legal-Bot now routes with three separate concepts:

- `engine`: the local CLI Legal-Bot invokes: `codex`, `claude`, `gemini`, or `copilot`
- `model`: the model string passed through to that engine, such as `gpt-5.4`, `gpt-5.5`, `gpt-5.4-mini`, `sonnet`, `opus`, `pro`, `flash`, or `auto`
- `effort`: optional reasoning level when the engine supports it, such as `low`, `medium`, `high`, `xhigh`, `max`, or `minimal`

`model_profile` remains the active routing control for intake, workflow, and review steps, but profiles now resolve to runnable local CLI targets instead of documentation-only provider strings.

Supported compact model specs:

- `codex:gpt-5.4:medium`
- `codex:gpt-5.4:high`
- `codex:gpt-5.5:medium`
- `claude:sonnet:high`
- `claude:opus:xhigh`
- `gemini:pro`
- `gemini:flash`
- `copilot:auto`

Legacy compatibility:

- `openai:gpt-5.5` is still accepted
- it now means `codex` + `gpt-5.5`
- it does not invoke the OpenAI API directly
- it no longer silently falls back to Claude just because `openai` is not an executable

Example `model_profiles` entry:

```json
"review_reasoning": {
  "primary": {
    "engine": "codex",
    "model": "gpt-5.4",
    "effort": "medium"
  },
  "fallbacks": [
    {
      "engine": "claude",
      "model": "sonnet",
      "effort": "medium"
    },
    {
      "engine": "gemini",
      "model": "pro"
    }
  ]
}
```

Fallback behavior:

- Legal-Bot checks the configured targets in order
- it picks the first target whose engine executable is available on `PATH`
- if it had to skip earlier targets, the step output prints the selected fallback and the reason
- if no configured engine is installed, the step fails with the list of engines it tried

`--single-llm` still accepts plain engine names such as `codex` or `claude`, and it now also accepts compact target specs such as `codex:gpt-5.4:high` and `gemini:flash`.

## Build

```powershell
cd go
go build -o bin/legal-bot.exe ./cmd/legal-bot
```

## Documentation

- [User Workflow](/c:/GIT/legal-bot/docs/user_workflow.md)
- [Legal Safety Boundaries](/c:/GIT/legal-bot/docs/legal_safety_boundaries.md)
- [Agent Persona System](/c:/GIT/legal-bot/docs/agent_persona_system.md)
- [Architecture](/c:/GIT/legal-bot/docs/architecture.md)
- [Legacy Surfaces](/c:/GIT/legal-bot/docs/legacy_surfaces.md)
- [Token Strategy](/c:/GIT/legal-bot/docs/token_strategy.md)
- [Legacy Discovery Migration](/c:/GIT/legal-bot/docs/legacy_discovery_migration.md)
