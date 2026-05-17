# People and Roles

List all people relevant to the case. Use stable identifiers across intake, timelines, fact candidates, review packets, and draft edits.

Do not put privileged strategy notes in people descriptions.

## Role Categories

Use these role categories where useful:

- `party`
- `child`
- `attorney`
- `opposing_counsel`
- `judge`
- `commissioner`
- `GAL`
- `special_master`
- `custody_evaluator`
- `family_systems_therapist`
- `child_therapist`
- `medical_provider`
- `school_provider`
- `financial_expert`
- `vocational_expert`
- `investigator`
- `witness`
- `family_member`
- `other`

## Core Person Table

| id | display_name | role_category | formal_role | party_affiliation | source_sensitivity | decision_role | notes_for_use |
|---|---|---|---|---|---|---|---|
| P1 | [PARTY_1_NAME] | party | [Petitioner / Respondent / Parent] | user_side | confidential | party | [Use consistent short label in drafts and summaries] |
| P2 | [PARTY_2_NAME] | party | [Petitioner / Respondent / Parent] | opposing_side | confidential | party | [Note if name variations appear in records] |
| C1 | [CHILD_1_NAME or INITIALS] | child | [Child] | child_related | child_sensitive | unknown | [Use initials if privacy is preferred; age or grade may be enough] |
| A1 | [ATTORNEY_1_NAME] | attorney | [Lead counsel] | user_side | privileged | advocate | [Attorney-safe communications only] |
| A2 | [ATTORNEY_2_NAME] | opposing_counsel | [Opposing counsel] | opposing_side | confidential | advocate | [Track opposing theories in substantive docs, not here] |
| J1 | [JUDGE_OR_COMMISSIONER_NAME] | judge | [Judge / Commissioner] | neutral | public_record | decision_maker | [Identify current forum and any posture-specific notes] |

Allowed `party_affiliation` values:

- `user_side`
- `opposing_side`
- `neutral`
- `child_related`
- `unknown`

Allowed `source_sensitivity` values:

- `public_record`
- `confidential`
- `privileged`
- `child_sensitive`
- `medical`
- `financial`
- `unknown`

Allowed `decision_role` values:

- `decision_maker`
- `recommender`
- `provider`
- `witness`
- `advocate`
- `party`
- `unknown`

## Children

Use initials or placeholders for children if privacy is preferred. Do not include unnecessary dates of birth unless needed. Age or grade is often enough.

| id | display_name | age_or_grade | household_context | source_sensitivity | notes_for_use |
|---|---|---|---|---|---|
| C1 | [CHILD_1_INITIALS] | [age / grade] | [shared schedule / primary residence / other neutral descriptor] | child_sensitive | [Child-sensitive details only when needed for relief or chronology] |

## Attorneys and Court Roles

| id | display_name | role_category | party_affiliation | decision_role | notes_for_use |
|---|---|---|---|---|---|
| A1 | [ATTORNEY_1_NAME] | attorney | user_side | advocate | [Primary drafting attorney / reviewer] |
| A2 | [ATTORNEY_2_NAME] | opposing_counsel | opposing_side | advocate | [Track opposing theories in substantive docs, not here] |
| J1 | [JUDGE_NAME] | judge | neutral | decision_maker | [Use if matter is judge-facing] |
| CMR1 | [COMMISSIONER_NAME] | commissioner | neutral | decision_maker | [Use if temporary or docket-facing issues are before a commissioner] |

## Child-Related Professionals and Neutral Roles

| id | display_name | role_category | serves_or_evaluates | party_affiliation | source_sensitivity | decision_role | notes_for_use |
|---|---|---|---|---|---|---|---|
| GAL1 | [GAL_NAME] | GAL | [child / family] | neutral | child_sensitive | recommender | [Do not imply support unless a source says so] |
| SM1 | [SPECIAL_MASTER_NAME] | special_master | [implementation / dispute resolution] | neutral | confidential | decision_maker | [Describe actual scope only if ordered] |
| EVAL1 | [EVALUATOR_NAME] | custody_evaluator | [custody / parent-time issues] | neutral | confidential | recommender | [Separate appointment from conclusions] |
| FST1 | [FAMILY_SYSTEMS_THERAPIST_NAME] | family_systems_therapist | [family system] | neutral | child_sensitive | provider | [Do not overstate role beyond records] |
| CT1 | [CHILD_THERAPIST_NAME] | child_therapist | [child] | child_related | child_sensitive | provider | [Handle records carefully] |

## Medical, School, Financial, and Other Professionals

| id | display_name | role_category | relationship_or_scope | party_affiliation | source_sensitivity | decision_role | notes_for_use |
|---|---|---|---|---|---|---|---|
| MED1 | [MEDICAL_PROVIDER_NAME] | medical_provider | [pediatrician / psychiatrist / therapist / clinic] | neutral | medical | provider | [Identify only what the records actually show] |
| SCH1 | [SCHOOL_CONTACT_NAME] | school_provider | [teacher / counselor / administrator] | neutral | child_sensitive | provider | [Useful for attendance, performance, and logistics records] |
| FIN1 | [FINANCIAL_EXPERT_NAME] | financial_expert | [tracing / support / valuation] | neutral | financial | recommender | [Do not imply conclusions not in the report] |
| VOC1 | [VOCATIONAL_EXPERT_NAME] | vocational_expert | [earning capacity / employability] | neutral | financial | recommender | [Use carefully in support disputes] |
| INV1 | [INVESTIGATOR_NAME] | investigator | [fact development / service / records] | neutral | confidential | witness | [Clarify scope and source basis] |

## Family Members, Witnesses, and Other People

| id | display_name | role_category | relationship_to_case | party_affiliation | source_sensitivity | decision_role | notes_for_use |
|---|---|---|---|---|---|---|---|
| W1 | [WITNESS_NAME] | witness | [neighbor / coach / friend / employer / family member] | unknown | confidential | witness | [Summarize relevance, not argument] |
| FM1 | [FAMILY_MEMBER_NAME] | family_member | [grandparent / sibling / partner / other] | user_side | confidential | witness | [Note whether source support exists] |
| O1 | [OTHER_NAME] | other | [describe role] | unknown | unknown | unknown | [Use when no better category fits] |

## Professional Role Notes

### GAL
- Neutral role, not automatically aligned with either side.
- Do not imply a GAL supports a position unless the source documents show it.

### Special Master
- Define only the actual scope of authority in the orders or stipulations.
- Do not assume tie-break powers or enforcement powers unless documented.

### Custody Evaluator
- Distinguish appointment, process, interim communications, and final opinions.
- Do not collapse an evaluator's limited observation into a global finding.

### Family Systems Therapist or Therapist
- Treat therapy records and professional statements as sensitive and role-limited.
- Do not imply the therapist is making custody or legal recommendations unless the record says so.

### Medical Provider
- Use the provider's role and actual records carefully.
- Medical records may establish treatment, recommendations, or reported history, but not every litigation conclusion a party wants to draw.

### School Personnel
- Strong for attendance, performance, logistics, and documented communications.
- Do not imply school staff support a custody theory unless they actually stated that.

### Investigator or Expert Witness
- Separate the person's role from the user's interpretation of the role.
- Note whether the professional is neutral, retained, informal, appointed, or unknown.

## Usage Guidance

- Use initials or placeholders for children if privacy is preferred.
- Do not include unnecessary dates of birth unless needed.
- Use age or grade when sufficient.
- Mark child-sensitive information carefully.
- Keep this file descriptive, not argumentative.
- Use the IDs consistently across document cards, timelines, findings, and review packets.
