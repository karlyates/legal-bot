# Token Strategy

## Core Principle

Ingest once. Reuse narrowly. Route deliberately.

The biggest cost risk is making every persona reread the whole source corpus. Legal-Bot is designed to push expensive reading into intake, then let later workflows rely on compact matter knowledge plus narrowly selected workflow inputs.

## Intake as the Compression Layer

Intake turns raw source files into reusable artifacts:

- document register
- document cards
- timeline
- fact candidates
- open questions

That knowledge layer exists so later triage, counsel brief, pattern review, fact lock, and draft review runs do not have to reread everything by default.

## Workflow Cost Strategy

Draftless workflows should usually depend on:

- `case/` guidance
- `situation.md`
- `goal.md`
- `questions.md`
- knowledge artifacts when available
- an input manifest when knowledge is missing

They should not dump the full raw matter corpus into every run unless the workflow truly needs it.

## Review Cost Strategy

Draft review costs are controlled mainly by:

- effort level
- document type
- conditional specialists

Use:

- `quick` when a fast readiness screen is enough
- `standard` for most reviews
- `deep` only when the extra strategic depth is justified

## Authority Cost Strategy

Authority work is expensive and high-risk.

Use it when:

- a recommendation depends on law, rules, or procedure
- a draft makes legal propositions
- the user wants authority gaps checked

Do not route every workflow through `legal-authority-scholar` by default.

## Model Profile Strategy

The repo uses provider-neutral `model_profile` names:

- `intake_long_context`
- `structured_extraction`
- `review_reasoning`
- `strategy_reasoning`
- `legal_research`
- `final_synthesis`
- `cheap_fast`

These profiles describe the work, not the vendor.

In the current runner, each profile resolves to:

- a local CLI engine
- an optional model string
- an optional effort level

Typical examples:

- `structured_extraction` -> `codex:gpt-5.4-mini:medium`
- `review_reasoning` -> `codex:gpt-5.4:medium`
- `final_synthesis` -> `codex:gpt-5.4:high`
- `legal_research` -> `codex:gpt-5.4:high`

Fallbacks remain important for cost and availability. If the preferred engine is not installed, Legal-Bot walks the configured fallback chain in order and prints the route it chose.

## Cheap Model Boundary

Cheap or smaller models are appropriate for:

- formatting cleanup
- mechanical docs updates
- reference validation
- schema checks
- low-risk maintenance

Cheap or smaller models are not the right final authority for:

- legal conclusions
- authority analysis
- emergency strategy
- credibility-sensitive review
- final synthesis

## Cost-Saving Rules

1. Do not make every agent read the full corpus.
2. Do not run deep review by default.
3. Do not run authority review unless law is actually in play.
4. Do not run strategic options work for simple proofreading.
5. Do not regenerate intake artifacts unless source materials changed.
6. Do not pass every sub-finding through to the final packet.
7. Prefer matter-level workflow inputs over giant prompt dumps.

## Product Benefit

This strategy is not only about lower token spend.

It also produces:

- less churn
- clearer recommendations
- smaller prompt surfaces
- cleaner managing-partner synthesis
- better day-to-day usability
