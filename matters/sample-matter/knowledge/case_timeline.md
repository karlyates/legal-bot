# Chronology Clerk Timeline

## Master Timeline

| event_id | event_date | document_date | filing_date | title | description | source_documents | source_tier | people | children | issue_tags | procedural_or_factual | dispute_status | confidence | notes |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| E-001 | 2026-04-29 | UNKNOWN | UNKNOWN | Reported family-medicine appointments and partial notice | Parent B’s intake summary reports that Parent A took both children to scheduled family-medicine appointments on 2026-04-29. Parent B later received notice about only one child’s appointment and reports no notice about the second child’s appointment. The same summary says Parent B later learned from a separate provider call note that a psychiatric referral may have been started. | [DOC-001](/C:/GIT/legal-bot/document_cards/DOC-001.md), [document_register.csv](/C:/GIT/legal-bot/document_register.csv) | Tier 4 for the underlying narrative; Tier 5 for this timeline synthesis | Parent A, Parent B | Child 1, Child 2 | joint_legal_custody; medical_notice; provider_appointment; psychiatric_referral; children | factual | unverified / reported by user | Medium | This is an intake summary, not a contemporaneous medical record or court filing. The date of the appointment is reported as 2026-04-29. No document date or filing date is supplied. |

## Procedural Timeline

| event_id | event_date | document_date | filing_date | title | description | source_documents | source_tier | people | children | issue_tags | procedural_or_factual | dispute_status | confidence | notes |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| P-001 | UNKNOWN | UNKNOWN | UNKNOWN | Intake issue raised for possible escalation | The source document functions as an intake and routing summary. It asks whether the medical-notice issue should be raised with counsel, preserved as a pattern, sent to a Special Master, or ignored for now. No filing, hearing, or court action is provided. | [DOC-001](/C:/GIT/legal-bot/document_cards/DOC-001.md) | Tier 4 | Parent B | Child 1, Child 2 | intake; attorney_review; special_master; medical_notice | procedural | unverified / internal intake only | Medium | This establishes issue-spotting posture only, not a case event. |

## Child-Related Timeline

| event_id | event_date | document_date | filing_date | title | description | source_documents | source_tier | people | children | issue_tags | procedural_or_factual | dispute_status | confidence | notes |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| C-001 | 2026-04-29 | UNKNOWN | UNKNOWN | Reported child medical appointments | The intake summary reports that both children attended family-medicine appointments on 2026-04-29. | [DOC-001](/C:/GIT/legal-bot/document_cards/DOC-001.md) | Tier 4 | Parent A, Parent B | Child 1, Child 2 | children; medical; provider_appointment | factual | reported | Medium | The appointment date is reported, not independently verified by a medical record in the current packet. |
| C-002 | After 2026-04-29 | UNKNOWN | UNKNOWN | Reported later notice to Parent B about only one child | Parent B later received notice about only one child’s appointment, according to the intake summary. | [DOC-001](/C:/GIT/legal-bot/document_cards/DOC-001.md) | Tier 4 | Parent B | Child 1, Child 2 | medical_notice; children; communication | factual | reported | Low | The exact notice date is not supplied. |
| C-003 | After 2026-04-29 | UNKNOWN | UNKNOWN | Reported later discovery of possible psychiatric referral | Parent B later learned from a separate provider call note that a psychiatric referral may have been started. | [DOC-001](/C:/GIT/legal-bot/document_cards/DOC-001.md) | Tier 4 | Parent B, referral provider | Child 1, Child 2 | psychiatric_referral; children; provider_record | factual | tentative / unverified | Low | The document uses tentative language: “may have been started.” No referral record is provided. |

## Issue-Specific Timelines

### Medical Notice / Joint Legal Custody

| event_id | event_date | document_date | filing_date | title | description | source_documents | source_tier | people | children | issue_tags | procedural_or_factual | dispute_status | confidence | notes |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| I-MED-001 | 2026-04-29 | UNKNOWN | UNKNOWN | Reported appointment activity relevant to medical notice | The source summary ties the appointments to a possible joint-legal-custody notice issue because Parent B says notice was incomplete. | [DOC-001](/C:/GIT/legal-bot/document_cards/DOC-001.md) | Tier 4 | Parent A, Parent B | Child 1, Child 2 | joint_legal_custody; medical_notice | factual | alleged issue | Medium | The record does not include the operative custody order, so the notice obligation cannot be measured against actual order language. |

### Psychiatric Referral

| event_id | event_date | document_date | filing_date | title | description | source_documents | source_tier | people | children | issue_tags | procedural_or_factual | dispute_status | confidence | notes |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| I-PSY-001 | After 2026-04-29 | UNKNOWN | UNKNOWN | Possible psychiatric referral | The intake summary reports a later provider call note suggesting a psychiatric referral may have been initiated. | [DOC-001](/C:/GIT/legal-bot/document_cards/DOC-001.md) | Tier 4 | Parent B, referral provider | Child 1, Child 2 | psychiatric_referral; provider_record | factual | tentative | Low | This is the weakest factual point in the packet because the source is indirect and tentative. |

## Status Quo / Turning Point Timeline

| event_id | event_date | document_date | filing_date | title | description | why_it_may_be_a_turning_point | source_documents | confidence | notes |
|---|---|---|---|---|---|---|---|---|---|
| SQ-001 | 2026-04-29 | UNKNOWN | UNKNOWN | Reported child medical appointments with partial notice | The intake summary reports a same-day child medical event and later partial notice to Parent B. | Medical appointments can affect ordinary decision-making, provider relationships, and notice expectations under joint legal custody. If a referral was actually initiated, that could shift the medical decision posture further. | [DOC-001](/C:/GIT/legal-bot/document_cards/DOC-001.md) | Medium | The source is not strong enough to say a status quo changed as a proven fact. It is enough to flag the event as a possible turning point requiring source support. |
| SQ-002 | After 2026-04-29 | UNKNOWN | UNKNOWN | Possible psychiatric referral | The summary suggests a psychiatric referral may have been started. | A mental-health referral, if confirmed, can change the scope of medical decision-making and may create a new care pathway or escalation point. | [DOC-001](/C:/GIT/legal-bot/document_cards/DOC-001.md) | Low | This remains tentative until a provider record, portal entry, or billing record confirms it. |

## Timeline Gaps

- The operative custody order or decree is missing.
- The actual notice message for the one child is missing.
- The actual provider call note is missing.
- Full medical records for the 2026-04-29 appointments are missing.
- Any referral documentation is missing.
- Any insurance or billing record is missing.
- There is no filed pleading or docket entry showing this issue has been brought to court.
- No Special Master order or scope document is provided.

## Date Conflicts / Sequence Confusion

- No direct date conflict appears in the supplied material.
- The main sequence issue is incomplete chronology:
- The appointment date is reported as 2026-04-29.
- The later notice date is not given.
- The later provider-call-note date is not given.
- It is unclear whether the psychiatric referral was actually initiated, or only discussed.

## Questions That Would Change the Timeline

- What does the operative custody order require for medical notice or joint decision-making?
- What is the exact date of the notice Parent B received about one child?
- What does the actual provider call note say, and when was it created?
- Was a psychiatric referral actually entered, ordered, or merely discussed?
- Were both children seen at the same appointment time, or were there separate visits on the same date?
- Were any follow-up appointments, prescriptions, or billing entries generated?
- Is there any court filing, email, or OFW thread that confirms the notice sequence?

## Bottom-Line Chronology Assessment

- The only clearly dated event in the current packet is the reported medical appointment date of 2026-04-29.
- Everything else is either later in sequence but undated, or tentative.
- This record is enough for intake triage and source gathering.
- It is not enough, by itself, to establish a court-ready chronology of a notice violation or psychiatric-referral event.