# Workflow Output Contracts

Legal-Bot workflows are not supposed to end in vague multi-agent recap text.

Each workflow now has an explicit output contract so the final packet is:

- practical
- opinionated
- source-disciplined
- honest about evidence and authority gaps
- clear about what to do next

## Why These Contracts Exist

The workflow router already decides *which* path to run.

The missing piece was *what a good final packet must look like* once the panel finishes debating.

The contracts give the final synthesizer a stable target:

- answer the user's actual job
- choose a recommendation
- separate facts from recollection, inference, allegations, and strategy
- identify evidence and authority gaps
- include Do Not Chase analysis when escalation is weak, premature, or credibility-negative
- keep internal strategy separate from external-safe language

## Global Final Packet Rules

Every workflow final packet should:

- begin with a short Bottom Line recommendation
- include near-top metadata for workflow, matter, audience, urgency, recommended track, confidence, external-use readiness, source strength, and attorney review needed
- choose a Recommended Track instead of merely listing possibilities
- identify the strongest evidence and weakest evidence
- say what source would materially change the recommendation
- separate internal strategy from external-safe wording
- identify optics, overstatement, and proportionality risk
- identify legal-authority gaps instead of pretending they are resolved
- avoid unsupported accusations and filler warnings

## Source-Strength Taxonomy

The final packet should use these labels when source strength matters:

This source-strength taxonomy is part of the workflow contract layer.

1. Locked / source-supported fact
2. User recollection / user assertion
3. Inference
4. Disputed allegation
5. Unsupported claim
6. Legal proposition needing authority
7. Attorney-only / internal strategy
8. External-safe fact

This taxonomy is not meant to shut analysis down. It is meant to say what the weakness means and what source would fix it.

## Recommended Tracks

Depending on the workflow, Legal-Bot should be willing to land on tracks such as:

- pursue now
- limited record-building message
- preserve for pattern evidence
- counsel-first
- Special Master
- court/MTE/TRO/GAL/court-order issue
- evidence packet first
- journal only
- let it go
- do not send yet

Do Not Chase is a valid outcome when the issue is too small, too weak, too premature, or too optics-negative to press now.

## Internal Strategy vs External-Safe Language

The contracts force a separation between:

- internal strategy
- attorney-facing candor
- external-safe language

External-safe wording must stay narrower than internal strategy. It should use source-supported, audience-safe facts and should not smuggle in unsupported motive claims, diagnoses, or escalation theater.

## Workflow Contracts

### `triage`

Expected sections include:

- `Bottom Line`
- `What Happened`
- `Source Strength`
- `Why It Matters`
- `Options Matrix`
- `Do Not Chase Check`
- `Evidence Needed`
- `Recommended Next Step`
- `Counsel Questions`
- `External-Safe Language`

### `counsel-brief`

Expected sections include:

- `Subject / Proposed Email Title`
- `Short Version`
- `Why I Think This Matters`
- `Key Facts`
- `Evidence Already Available`
- `Evidence Missing`
- `Legal / Strategy Questions for Counsel`
- `Options I See`
- `Recommended Ask`
- `Attorney-Ready Draft`

### `sm-triage`

Expected sections include:

- `Bottom Line`
- `Issue Framed for the Special Master`
- `Scope / Authority Fit`
- `Ripeness`
- `Requested Directive`
- `Evidence Needed`
- `Optics / Credibility Risk`
- `Do Not Chase Check`
- `Recommended Next Step`
- `Draft SM Message`

### `draft-message`

Expected sections include:

- `Bottom Line`
- `Audience and Tone`
- `Message Goal`
- `Risk Check`
- `Facts Safe to Say`
- `Facts or Claims to Avoid`
- `Recommended Message`
- `Optional Shorter Version`
- `Preservation Notes`

### `evidence-packet`

Expected sections include:

- `Packet Purpose`
- `Issue Summary`
- `Evidence Inventory`
- `Missing Evidence`
- `Chronology Needed`
- `Fact Map`
- `Exhibits / Attachments Candidate List`
- `Risk / Weaknesses`
- `Recommended Packet Structure`
- `Next Actions`

### `pattern-review`

Expected sections include:

- `Bottom Line`
- `Pattern Theory`
- `Strongest Examples`
- `Weak or Risky Examples`
- `Pattern vs Isolated Noise`
- `Legal / Co-Parenting Significance`
- `Opposing / Neutral Read`
- `Do Not Chase Check`
- `Recommended Use`
- `Evidence Needed`

### `fact-lock`

Expected sections include:

- `Fact-Lock Summary`
- `Locked Facts`
- `Partially Supported Facts`
- `User Recollection / Unsourced Assertions`
- `Inferences`
- `Disputed Allegations`
- `Legal Propositions Needing Authority`
- `External-Use Readiness`
- `Recommended Next Step`

### `authority-check`

Expected sections include:

- `Bottom Line`
- `Legal Questions Presented`
- `Known Authority Sources in the Matter`
- `Authority Gaps`
- `Currentness / Citator Status`
- `Attorney Questions`
- `Safe Working Language`
- `Do Not Overclaim`

### `prep`

Expected sections include:

- `Prep Objective`
- `Bottom Line`
- `Key Facts to Lead With`
- `Decisions Needed`
- `Questions to Ask`
- `Likely Pushback / Attacks`
- `Best Responses`
- `Evidence to Have Ready`
- `What Not to Say`
- `Follow-Up Actions`

### `journal-entry`

Expected sections include:

- `Suggested Title`
- `Date / Date Range`
- `Record Type`
- `Priority`
- `Category / Related Issue`
- `Children`
- `Key People`
- `Summary`
- `What Happened`
- `Legal / Co-Parenting Significance`
- `Neutral Evaluator Notes`
- `Evidence / Source Anchors`
- `Caveats / Verification Needed`
- `Related Events`
- `Import Notes`

## Audience Rules

- OFW/Emily: brief, child-focused, calm, record-building
- Special Master: specific, directive-focused, professional
- provider/school: neutral, access/request focused, non-accusatory
- counsel: candid, strategic, concise
- court-facing: cautious, source-tethered, non-theatrical

## Prompt Injection

Workflow contracts are injected in code:

- all workflow steps receive a concise workflow purpose and output contract summary
- the final `managing-partner-final-synthesizer` step receives the full workflow output contract

That keeps the whole panel oriented toward the same final packet while reserving the heavier contract text for the final decision-maker.

## Maintenance

The contract definitions live in:

- `go/cmd/legal-bot/workflow_contracts.go`

When adding or changing workflow types:

1. update the route catalog if needed
2. update the workflow contract definitions
3. update docs here
4. update the static tests that validate workflow coverage and contract headings
