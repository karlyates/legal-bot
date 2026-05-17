# Finding Schema

Use this exact structure for shared findings:

```text
finding_id:
agent:
location:
severity:
category:
issue:
why_it_matters:
recommended_action:
suggested_language_if_any:
confidence:
source_ids_if_any:
attorney_review_required:
shareability:
```

Allowed `shareability` values:
- internal_only
- attorney_safe
- court_safe_candidate

Allowed `severity` values:
- Critical
- High
- Medium
- Low
- Note

Allowed `recommended_action` values:
- Keep
- Clarify
- Support with source
- Soften
- Strengthen
- Move
- Remove
- Ask attorney
- Needs source document
- Needs user clarification
- Needs legal authority
- Do not chase

Allowed `confidence` values:
- High confidence
- Medium confidence
- Low confidence
- Unknown / source not provided

Allowed `category` values:
- factual_support
- chronology
- legal_structure
- legal_authority
- relief_alignment
- tone_credibility
- attack_surface
- child_best_interest
- family_dynamics
- financial_support
- procedural
- evidence
- readability
- strategy
- attorney_question
