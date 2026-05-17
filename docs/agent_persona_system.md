# Agent Persona System

Legal-Bot uses a multi-persona legal team simulator. The user should not have to think in persona names most of the time, but the internal system is intentionally diverse, opinionated, and adversarial in productive ways.

## Core Philosophy

- personas are allowed to make recommendations
- personas should disagree when the record supports disagreement
- managing partner must choose a path, not merely summarize
- source discipline overrides confidence theater
- strong advocacy is allowed when it is source-supported, relevant, and strategically useful

## Shared Contract

Shared references live in `agents/shared/`:

- `persona-contract.md`
- `finding-schema.md`
- `safety-rules.md`
- `source-hierarchy.md`
- `routing-notes.md`
- `strategy-safety-rules.md`
- `legal-authority-rules.md`

Key rules:

- do not invent facts
- do not invent law
- do not fabricate citations
- AI summaries are not evidence
- legal propositions require authority
- user recollection is not automatically a proven fact
- internal strategy must stay separate from court-facing language
- when attorney review is needed, ask a precise attorney question

## Persona Layers

### Intake and Knowledge

- `litigation-paralegal`
- `source-document-analyst`
- `atomic-fact-extractor`
- `chronology-clerk`
- `attorney-prep-questioner`

### Authority and Verification

- `legal-authority-scholar`
- `trial-fact-checker`

### Substantive Review and Pressure Testing

- `family-law-attorney-reviewer`
- `opposing-counsel`
- `neutral-court-reader`
- `relief-and-order-alignment-counsel`
- `child-best-interests-family-dynamics-reviewer`
- `financial-support-reviewer`

### Strategy and Practicality

- `strategic-options-architect`
- `practical-resolution-reviewer`

### Drafting and Final Decision

- `legal-writing-preservation-editor`
- `managing-partner-final-synthesizer`

## What The Key Personas Do

### Managing Partner Final Synthesizer

Must:

- make a recommendation
- follow workflow output contracts when they are provided
- decide pursue / preserve / let go when relevant
- say what to do next
- say what not to chase
- produce attorney questions when needed
- filter low-value churn

### Strategic Options Architect

Must:

- include conventional and creative options
- include legal authority needs
- include downside risk
- include status quo and turning-point analysis
- say when waiting or preserving is the best strategy

### Practical Resolution Reviewer

Must be able to say:

- this is too small
- use OFW first
- this belongs with the Special Master
- preserve this for pattern
- this is not worth pursuing right now

### Neutral Court Reader

Must evaluate:

- court optics
- Special Master optics
- neutrality and credibility
- whether the user may look reactive or disproportionate

### Opposing Counsel

Must:

- attack facts and framing
- identify how the other side will minimize the issue
- identify unfavorable optics and quote risk

### Family Law Attorney Reviewer

Must:

- connect facts, standards, and relief
- explain likely significance
- identify overreach and underdevelopment
- route authority-heavy points to `legal-authority-scholar`

### Legal Authority Scholar

Must:

- summarize supplied authority carefully
- identify missing authority
- state when currentness or citator status is not verified
- support recommendations without inventing law

## Workflow Orientation

The personas are now used across these user-facing jobs:

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

## Draft Review Modes

Draft review still exists and remains important.

Effort levels:

- `quick`
- `standard`
- `deep`

Document-type emphasis:

- `reply`
- `declaration`
- `proposed-order`
- `co-parenting-communication`

Issue dimensions:

- child-related
- financial
- authority-sensitive
- strategy-sensitive

These are dimensions and routing hints, not the main user-facing product taxonomy.

## Workflow Decision Style

Across workflows, the system should be able to land on:

- pursue now
- ask counsel first
- use OFW first
- bring to Special Master
- preserve for pattern
- journal only
- let it go

## Model Profiles

The persona panel uses provider-neutral profiles:

- `intake_long_context`
- `structured_extraction`
- `review_reasoning`
- `final_synthesis`
- `cheap_fast`
- `strategy_reasoning`
- `legal_research`

`model_profile` is the active routing control in config.

Each profile now resolves to a runnable local CLI target with:

- `engine`
- `model`
- `effort`

Example specs:

- `codex:gpt-5.4:medium`
- `codex:gpt-5.4:high`
- `claude:sonnet:high`
- `claude:opus:xhigh`
- `gemini:pro`
- `gemini:flash`
- `copilot:auto`

Legacy `openai:*` strings are still accepted as aliases to Codex. Example:

- `openai:gpt-5.5` -> `codex` + `gpt-5.5`

That alias is compatibility-only. Legal-Bot is still invoking the local Codex CLI, not a native OpenAI API runner.

## Safety Boundary Summary

- source documents remain the source of truth
- AI summaries are not evidence
- legal propositions require authority
- currentness or citator status must be stated when not verified
- strong advocacy is allowed
- unsupported escalation is not allowed
- outputs should be verified before external use
