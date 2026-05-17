# Language Calibration and Risk Flags

This file does not ban strong language. It helps legal-bot decide when language should be preserved, strengthened, softened, supported with evidence, moved to attorney-facing strategy, marked `internal_only`, or avoided.

Use this file to calibrate tone by audience, document type, source strength, legal authority, and strategy. Do not use it as an automatic blacklist.

## Core Rule: Flag, Do Not Automatically Sanitize

- Strong language should be evaluated in context.
- Attorney-drafted language should not be softened merely because it is forceful.
- Reviewers should ask what function the phrase serves.
- If the phrase is strong, supported, and strategically useful, preserve it or suggest a more precise version.
- If the phrase is strong but unsupported, speculative, inflammatory, or credibility-draining, flag it.
- If the phrase is really a legal standard, route it to legal-authority-scholar or counsel confirmation.
- If the phrase belongs in private strategy rather than a filing, mark it `internal_only` or `attorney_safe`.

## Language Decision Matrix

| Phrase / Pattern | Court-Facing Use | Attorney-Facing Use | Requires Source? | Requires Legal Authority? | Default Action | Notes |
|---|---|---|---|---|---|---|
| bad faith | only if standard is relevant and facts support it | yes | yes | yes | ask attorney | Legal term of art. Do not use as a casual insult. |
| knowingly | use if record supports knowledge inference | yes | yes | sometimes | preserve but narrow | Tie to notice, dates, and documents. |
| intentionally | use cautiously | yes | yes | sometimes | preserve with source citation | Avoid mind-reading. Show conduct first. |
| deliberately | use cautiously | yes | yes | no | preserve but narrow | Better when tied to sequence and repeated acts. |
| misrepresented | use if contradiction is documented | yes | yes | yes | preserve with source citation | Prefer specific inconsistency over generalized accusation. |
| concealed | use if nondisclosure is documented | yes | yes | sometimes | preserve but narrow | Clarify what was withheld, from whom, and when. |
| violated | often appropriate in enforcement context | yes | yes | yes | preserve with source citation | Cite the order paragraph or specific obligation. |
| refused | appropriate if request and response are documented | yes | yes | yes | preserve with source citation | Quote or cite the communication when possible. |
| failed to comply | appropriate in compliance framing | yes | yes | yes | preserve with source citation | Stronger than vague criticism when tied to the record. |
| unilateral | often appropriate in joint decision disputes | yes | yes | yes | preserve | Useful in school, medical, therapy, and schedule issues. |
| coercive | court-facing only with facts and careful framing | yes | yes | yes | court-safe if reframed | Describe the conduct, pressure, and impact. |
| harassing | use cautiously | yes | yes | sometimes | court-safe if reframed | Can be a legal label in some contexts. |
| abusive | use only if conduct and context support it | yes | yes | yes | court-safe if reframed | Prefer conduct description over label unless the label matters legally. |
| alienation | usually attorney-facing unless expert or authority support exists | yes | yes | yes | attorney-facing only | High-risk term. Do not state as fact without support. |
| coached | use cautiously | yes | yes | no | soften | Better as concern or source-backed observation. |
| unsafe | use if risk facts are concrete | yes | yes | yes | preserve but narrow | Identify the specific safety issue. |
| endangerment | only if supported and legally relevant | yes | yes | yes | ask attorney | Often implicates a legal standard. |
| unfit | rarely appropriate without authority and strong support | yes | yes | yes | attorney-facing only | High-risk legal conclusion. |
| narcissistic / narcissist | no | maybe internal_only | no | no | remove | Diagnosis or labeling language usually hurts credibility. |
| liar / lying | usually no | maybe | yes | no | court-safe if reframed | Replace with inconsistency language unless counsel chooses otherwise. |
| retaliatory | use cautiously | yes | yes | sometimes | preserve but narrow | Needs timing and sequence support. |
| obstruction | use if conduct fits and is documented | yes | yes | yes | preserve with source citation | May overlap with discovery or compliance terminology. |
| contempt | only if counsel is invoking the legal concept | yes | yes | yes | ask attorney | Legal term of art requiring authority and procedural fit. |
| irreparable harm | only if remedy depends on it | yes | yes | yes | needs legal authority | Legal standard, not just emphasis. |
| emergency | use when facts justify urgency | yes | yes | yes | preserve but narrow | Overuse creates credibility risk. |
| status quo | often appropriate | yes | yes | yes | preserve | Clarify what the status quo is and why it matters. |
| pattern | appropriate when multiple events are documented | yes | yes | yes | preserve with source citation | Avoid if only one event is shown. |
| weaponized | usually attorney-facing or internal-only | yes | yes | no | move | Better reframed as specific conduct and effect. |
| manipulation | usually attorney-facing unless tied to conduct | yes | yes | no | soften | Describe acts, timing, and impact instead. |

## Strong Language That May Be Appropriate

Strong language can be appropriate when it is factually supported, legally relevant, and strategically useful. It is often most defensible in:

- motion to enforce
- opposition / response
- reply
- sanctions or fees request
- protective order or injunction response
- emergency relief
- attorney strategy memo
- internal legal analysis
- sections where counsel intentionally uses forceful advocacy

Examples:

- "failed to comply with paragraph [X] of the order"
- "made unilateral decisions despite joint legal custody"
- "withheld notice regarding [school/medical/therapy issue]"
- "created a record without the other parent's participation"
- "the timing raises concern"
- "the conduct appears inconsistent with the decree"
- "the requested relief is necessary to preserve the status quo"

## Language That Usually Needs Softening

Language usually needs softening when it:

- diagnoses
- assigns motive without support
- claims certainty without evidence
- makes broad character attacks
- turns litigation theory into fact
- is not tied to a source

Examples:

- "narcissist"
- "evil"
- "crazy"
- "liar"
- "always"
- "never"
- "obviously"
- "everyone knows"
- "she/he is trying to destroy me"

## Legal Terms of Art

Some strong phrases are not merely tone issues. They may be legal terms that require authority and attorney review.

- bad faith
- contempt
- irreparable harm
- best interests
- material change
- endangerment
- coercive control
- alienation
- harassment
- stalking
- domestic violence
- unfit parent
- emergency
- status quo
- sanctions
- attorney fees

When these terms appear, legal-bot should ask:

1. Is this term being used as a legal standard?
2. Is the legal standard supplied?
3. Is the factual support supplied?
4. Should legal-authority-scholar verify it?
5. Should counsel confirm it?

## Audience Calibration

| Audience | Tone | Strong Language Allowed? | Notes |
|---|---|---|---|
| Court-facing motion | precise, purposeful, source-backed | yes, when supported and useful | Tie language to relief, evidence, and procedural posture. |
| Declaration | factual, careful, chronological | limited | Motive language is risky. Specific conduct is better. |
| Proposed order | precise, enforceable, non-argumentative | rarely | Orders should state duties, deadlines, and mechanisms. |
| Attorney email | direct, concise, practical | yes | Can state risks and leverage more plainly. |
| Internal strategy memo | candid, exploratory | yes | Mark clearly as `internal_only` when appropriate. |
| OFW / co-parenting message | restrained, neutral, child-focused | very limited | This is usually the most sensitive audience. |
| Settlement proposal | practical, strategic, measured | yes, carefully | Preserve leverage without unnecessary escalation. |
| Mediation statement | persuasive but controlled | yes | Useful for pattern and practical impact framing. |
| Legal research memo | careful, citation-first | yes, if authority supports it | Legal labels should be accurate and qualified. |

Key rule: attorney-facing and internal strategy language may be more direct. Court-facing language should be precise, supported, and purposeful. OFW and co-parenting language should be the most restrained.

## Document-Type Calibration

| Document Type | Recommended Tone | Strong Language Rule | Default Reviewer Behavior |
|---|---|---|---|
| Motion to Enforce | firm, compliance-focused, specific | stronger compliance language may be appropriate if tied to an order provision and facts | preserve or strengthen if source-backed |
| TRO / emergency motion | urgent, specific, credible | urgent language may be needed, but avoid overstatement | preserve but narrow unless evidence is thin |
| Opposition / response | pointed, controlled, rebuttal-focused | direct rebuttal can be appropriate | preserve with source support or reframe precisely |
| Reply | concise, corrective, strategic | stronger language only if it clarifies the real dispute | trim rhetoric that adds heat without leverage |
| Declaration | factual, chronological, careful | motive and diagnosis language usually needs support or reframing | move unsupported argument out of factual sections |
| Proposed order | neutral, specific, administrable | avoid argumentative phrasing | remove rhetoric and sharpen mechanics |
| Attorney handoff email | direct, issue-focused | yes | allow attorney-safe candor |
| Co-parenting communication | calm, narrow, practical | almost never | soften unless firmness is necessary for logistics |
| Mediation statement | persuasive, pragmatic | yes, but strategic and solution-oriented | preserve persuasive framing that helps resolution |
| Settlement proposal | practical, leverage-aware | yes, carefully | distinguish negotiation posture from filing-safe language |

## Strategy Leakage

Examples of strategy leakage:

- references to attorney-client communications
- "our plan is"
- "we are trying to set up"
- settlement positions in filed documents
- private litigation theories in declarations
- internal AI analysis in court-facing drafts
- expert-shopping or timing strategy unless counsel approves

Shareability labels:

- `internal_only`
- `attorney_safe`
- `court_safe_candidate`

If a phrase reveals leverage planning, settlement positioning, or private tactical thinking, keep it out of court-facing drafts unless counsel intentionally chooses otherwise.

## Preferred Replacement Patterns

Instead of "she lied" use:

- "her statement is inconsistent with [source]"
- "the record reflects a different sequence"
- "the available documents do not support that claim"

Instead of "she always violates the decree" use:

- "on [dates], the following order provisions were not followed"
- "the documented pattern raises compliance concerns"

Instead of "she is alienating the children" use:

- "the child-related statements and timing raise concerns about loyalty conflict or pressure, which should be evaluated by the appropriate professional"

Instead of "she is abusive" use:

- "the specific conduct at issue was [describe conduct], and the impact was [describe impact]"

If attorney-drafted strong language is source-supported and strategically purposeful, do not automatically replace it. Instead classify it as:

- preserve
- preserve with source citation
- preserve but narrow
- ask attorney
- soften
- remove

## Agent Instructions

### legal-writing-preservation-editor
- Do not rewrite strong language just to make it milder.
- Check whether the language is source-backed, strategically useful, and placed in the right document.
- When possible, suggest narrower and more precise wording rather than weaker wording.

### opposing-counsel
- Identify language the other side will quote as overstatement, speculation, or personal attack.
- Distinguish between powerful language and self-inflicted credibility damage.

### neutral-court-reader
- Flag when a judge or commissioner may see the wording as overstated, distracting, or unsupported.
- Preserve strong language when it makes the requested relief easier to understand and the record supports it.

### strategic-options-architect
- Use strong internal framing when it helps counsel think clearly, but keep shareability explicit.
- Route authority-sensitive terms to legal-authority-scholar or counsel.

### managing-partner-final-synthesizer
- Decide whether strong language is worth preserving.
- Reject low-value softening that would only dilute good advocacy.
- Keep internal strategy, attorney-safe language, and court-facing language in separate buckets.

### legal-authority-scholar
- Flag when a phrase depends on a legal standard or term of art.
- Do not pass legal labels to other agents as reliable unless supported by supplied authority or research support.
