# User Workflow

## Main Question

The primary user-facing question is:

`What are you trying to do?`

The default front door is now:

```powershell
legal-bot guide
legal-bot guide <matter>
```

It asks a few deterministic routing questions and then invokes the existing `intake`, `workflow`, or `review` command internally.

## Current Workflows

Implemented in this pass:

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

## Commands

### Guided Front Door

```powershell
legal-bot guide
legal-bot guide <matter>
legal-bot start <matter>
legal-bot guide --list
legal-bot guide <matter> --dry-run
legal-bot routes
legal-bot routes --json
```

What the guide does:

1. asks for the matter if needed
2. detects input, drafts, and knowledge-layer readiness
3. asks what the user is trying to do
4. gathers only the route-specific details
5. prints the planned command
6. asks for confirmation before running in interactive mode

Guided choice map:

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

Examples:

```powershell
legal-bot guide Revere-042926
legal-bot guide Revere-042926 --dry-run
legal-bot guide --list
legal-bot guide Revere-042926 --intent counsel-brief --situation "Medical billing issue and provider access confusion." --goal "Prepare a concise attorney-ready brief."
legal-bot guide TRO-MTE-GAL-SM --intent review
legal-bot routes
```

`legal-bot routes` is the quick operator check for the active control plane. It shows the current command surface, review modes, document-type mappings, workflow types, conditionals, and legacy wrapper inventory.

### Matter Intake

```powershell
legal-bot intake <matter>
```

Purpose:

- build or refresh the knowledge layer
- register source documents
- produce timeline, fact candidates, and open questions

### Draftless Workflows

```powershell
legal-bot workflow <matter> --type triage
legal-bot workflow <matter> --type evidence-packet
legal-bot workflow <matter> --type counsel-brief
legal-bot workflow <matter> --type sm-triage
legal-bot workflow <matter> --type draft-message --audience ofw
legal-bot workflow <matter> --type pattern-review
legal-bot workflow <matter> --type fact-lock
legal-bot workflow <matter> --type authority-check
legal-bot workflow <matter> --type prep
legal-bot workflow <matter> --type journal-entry
```

Workflow outputs now have explicit final-packet contracts. That means the final packet structure is more predictable by workflow type and should reliably include things like Bottom Line, recommended track, source-strength framing, evidence gaps, attorney questions when needed, and Do Not Chase analysis where it belongs.

### Draft Review

```powershell
legal-bot review <matter>
legal-bot review-draft <matter>
legal-bot review <matter> --mode quick
legal-bot review <matter> --mode standard
legal-bot review <matter> --mode deep
legal-bot review <matter> --mode legal-research
legal-bot review <matter> --mode strategy
```

Document-type notes:

- `motion`, `opposition`, and `response` currently fold into standard review
- `reply`, `declaration`, `proposed-order`, and `co-parenting-communication` use document-type-specific review routes

## Model Overrides

For `intake`, `workflow`, and `review`, `--single-llm` still accepts plain engine names:

- `codex`
- `claude`
- `gemini`
- `copilot`

It also accepts compact engine/model/effort specs:

- `codex:gpt-5.4:medium`
- `codex:gpt-5.4:high`
- `codex:gpt-5.5:medium`
- `claude:sonnet:high`
- `claude:opus:xhigh`
- `gemini:pro`
- `gemini:flash`
- `copilot:auto`

Precedence is:

1. `--single-llm`
2. step `model_profile`
3. step `llm`

During a run, each step now prints the selected route, including fallback use when applicable.

## Matter File Conventions

For draftless workflows, the current command looks for:

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

What they mean:

- `situation.md`: raw description of what happened
- `goal.md`: what the user wants to decide or accomplish
- `questions.md`: specific questions to answer

If a workflow requires a situation file and none is found, the command fails with a helpful message instead of silently falling back to a draft-centric path.

## Workflow Selection Map

### 1. Understand a new situation

Use:

```powershell
legal-bot workflow <matter> --type triage
```

Output focus:

- Bottom Line Recommendation
- Pursue / Preserve / Let Go
- Why it matters or may not matter
- Evidence available
- Evidence missing
- Risks of acting
- Risks of waiting
- Best next step
- Do Not Chase Check

### 2. Build an evidence packet

Use:

```powershell
legal-bot workflow <matter> --type evidence-packet
```

Output focus:

- evidence needed
- evidence available
- evidence missing
- weak evidence
- evidence that cuts against the user
- what should stay attorney-only

### 3. Write to Kaitlyn or Jessika

Use:

```powershell
legal-bot workflow <matter> --type counsel-brief
```

Optional:

- `--urgency urgent`
- `--urgency call-agenda`
- `--urgency evidence-update`
- `--urgency draft-review-feedback`
- `--urgency strategic-question`

Output focus:

- subject line
- short version
- why this matters
- key facts
- evidence available
- evidence missing
- recommended ask
- questions for counsel
- copy-paste email

### 4. Review a draft

Use:

```powershell
legal-bot review <matter>
```

Use effort level for how hard to review:

- `quick`
- `standard`
- `deep`

Use document type for emphasis:

- `reply`
- `declaration`
- `proposed-order`
- `co-parenting-communication`

### 5. Decide whether to involve the Special Master

Use:

```powershell
legal-bot workflow <matter> --type sm-triage
```

Output focus:

- whether it is within likely SM scope
- whether it is ripe
- whether OFW should be tried first
- whether more evidence is needed
- whether it is better preserved for pattern
- concise Special Master message if appropriate

### 6. Draft an external message

Use:

```powershell
legal-bot workflow <matter> --type draft-message --audience ofw
```

Supported audience hints:

- `ofw`
- `special-master`
- `provider`
- `school`
- `counsel`
- `court-cover`
- `private-journal`

### 7. Review whether events add up to a pattern

Use:

```powershell
legal-bot workflow <matter> --type pattern-review
```

Output focus:

- pattern identified or not
- strongest examples
- noisy examples
- what the pattern supports
- whether to raise now or preserve

### 8. Lock facts before external use

Use:

```powershell
legal-bot workflow <matter> --type fact-lock
```

Output focus:

- locked facts
- likely but needs source
- recollection
- inference
- disputed
- unsupported
- safe wording
- unsafe wording

### 9. Check authority gaps

Use:

```powershell
legal-bot workflow <matter> --type authority-check
```

Important:

- This does not overclaim full live legal research if only supplied authority is available.
- `review --mode legal-research` still exists for draft-based authority review.

### 10. Prepare for a call or hearing

Use:

```powershell
legal-bot workflow <matter> --type prep
```

Output focus:

- goal
- top points
- likely pushback
- evidence to have open
- what not to say
- desired outcome
- fallback position

### 11. Create a journal entry

Use:

```powershell
legal-bot workflow <matter> --type journal-entry
```

Output focus:

- date and confidence
- record type
- category
- source type
- summary
- significance
- caveats
- pattern value

## Intake vs Workflow vs Review

Rerun intake when:

- new source files are added
- source files change
- the matter knowledge layer is stale

Run a workflow when:

- a new situation happens
- there is no formal draft yet
- the real question is escalation, evidence, pattern, or messaging

Run review when:

- a draft exists
- the user needs attorney-readiness review
- the user wants effort-level or document-type review

## Guided Layer Status

Implemented:

- interactive guided routing via `legal-bot guide`
- `legal-bot guide --list`
- `legal-bot guide <matter> --dry-run`
- low-risk `start` alias
- non-interactive choice and intent selection

Guided input is saved under:

```text
matters/<matter>/working/
```

This keeps matter-specific text out of `case/`, which remains global prompt context.

## Active vs Legacy

The active family-law workflow is:

- `guide`
- `intake`
- `review`
- `workflow`
- `routes`

The repo still contains Vern-era discovery, council, oracle, generate, VTS, and related compatibility code. Those are not the normal Legal-Bot legal workflow surface. See [Legacy Surfaces](/c:/GIT/legal-bot/docs/legacy_surfaces.md).

For the final workflow packet shapes and source-strength taxonomy, see [Workflow Output Contracts](/c:/GIT/legal-bot/docs/workflow_output_contracts.md).

## Practical Notes

- `case/` is global prompt context
- `matters/<matter>/input/` is matter-specific source context
- source documents remain the source of truth
- AI summaries are not evidence
- outputs should be verified before being filed, served, sent externally, or relied on
