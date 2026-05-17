# Relief Library

Purpose: this file tracks relief theories, requested orders, supporting facts, authority status, evidence, risks, and order mechanics. It helps legal-bot check whether a motion, declaration, prayer, proposed order, and strategy packet are all pointing to the same practical outcome.

## How to Use This File

For each live relief item, fill in:

- `relief_id`
- `matter`
- `relief_type`
- `relief_requested`
- `status`
- `status_quo_effect`
- `legal_authority_status`
- `key_facts_needed`
- `source_ids`
- `proposed_order_terms`
- `deadline_or_mechanism`
- `enforcement_issue`
- `attack_surface`
- `attorney_questions`

Allowed `status_quo_effect` values:

- `preserves_status_quo`
- `restores_status_quo`
- `changes_status_quo`
- `unclear`

Allowed `legal_authority_status` values:

- `supplied`
- `verified_by_legal_authority_scholar`
- `attorney_confirmation_required`
- `unknown`

## Active Relief Templates

### TRO / Emergency Relief

| relief_id | matter | relief_type | relief_requested | status | status_quo_effect | legal_authority_status | key_facts_needed | source_ids | proposed_order_terms | deadline_or_mechanism | enforcement_issue | attack_surface | attorney_questions |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| R-TRO-1 | [MATTER] | TRO / emergency relief | [Specific temporary protection or restraint requested] | [planned/pending/granted/denied] | [preserves_status_quo/restores_status_quo/changes_status_quo/unclear] | [supplied/verified_by_legal_authority_scholar/attorney_confirmation_required/unknown] | [incident dates, immediate risk, practical harm, why ordinary timing is inadequate] | [DOC-001, DOC-014] | [who is restrained or required to act; what is prohibited; what happens next] | [effective immediately / until hearing / by date] | [how violation would be identified and addressed] | [overstatement, insufficient urgency, broad relief, weak chronology] | [What authority supports the remedy? What narrower version is most defensible?] |

### Motion to Enforce

| relief_id | matter | relief_type | relief_requested | status | status_quo_effect | legal_authority_status | key_facts_needed | source_ids | proposed_order_terms | deadline_or_mechanism | enforcement_issue | attack_surface | attorney_questions |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| R-MTE-1 | [MATTER] | Motion to Enforce | [Order enforcement, compliance deadline, makeup parent-time, records production, etc.] | [planned/pending/granted/denied] | [preserves_status_quo/restores_status_quo/changes_status_quo/unclear] | [supplied/verified_by_legal_authority_scholar/attorney_confirmation_required/unknown] | [operative order, paragraph number, notice, noncompliance events, resulting problem] | [DOC-002, DOC-017] | [cite paragraph; state who must do what; include compliance deadline] | [by date / within X days / exchange schedule] | [how future noncompliance is handled] | [missing order cite, no proof of notice, requested relief broader than violation] | [Is contempt language appropriate or should the motion stay at enforcement?] |

### Clarification

| relief_id | matter | relief_type | relief_requested | status | status_quo_effect | legal_authority_status | key_facts_needed | source_ids | proposed_order_terms | deadline_or_mechanism | enforcement_issue | attack_surface | attorney_questions |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| R-CLAR-1 | [MATTER] | Clarification | [Clarify ambiguous order language or implementation process] | [planned/pending/granted/denied] | [preserves_status_quo/restores_status_quo/changes_status_quo/unclear] | [supplied/verified_by_legal_authority_scholar/attorney_confirmation_required/unknown] | [ambiguous paragraph, conflicting interpretations, practical breakdown] | [DOC-003, DOC-011] | [clean implementation language] | [when clarification takes effect] | [what happens if parties still disagree] | [looks like modification instead of clarification] | [Does the request alter substance or only clarify mechanics?] |

### Modification

| relief_id | matter | relief_type | relief_requested | status | status_quo_effect | legal_authority_status | key_facts_needed | source_ids | proposed_order_terms | deadline_or_mechanism | enforcement_issue | attack_surface | attorney_questions |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| R-MOD-1 | [MATTER] | Modification | [Change to custody, parent-time, decision-making, support, or logistics] | [planned/pending/granted/denied] | [preserves_status_quo/restores_status_quo/changes_status_quo/unclear] | [supplied/verified_by_legal_authority_scholar/attorney_confirmation_required/unknown] | [what changed, why current order is not working, effect on child or administration] | [DOC-004, DOC-018] | [new schedule, decision rule, allocation, or process] | [effective date / review date] | [how disputes are escalated] | [insufficient change showing, overbroad ask, missing authority] | [Does counsel want temporary, interim, or final relief?] |

### GAL

| relief_id | matter | relief_type | relief_requested | status | status_quo_effect | legal_authority_status | key_facts_needed | source_ids | proposed_order_terms | deadline_or_mechanism | enforcement_issue | attack_surface | attorney_questions |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| R-GAL-1 | [MATTER] | GAL | [Appointment, scope, issue focus] | [planned/pending/appointed/completed] | [preserves_status_quo/restores_status_quo/changes_status_quo/unclear] | [supplied/verified_by_legal_authority_scholar/attorney_confirmation_required/unknown] | [child-function concern, conflict level, need for neutral input] | [DOC-007, DOC-021] | [scope, access to records, reporting path] | [appointment by date / report by date] | [how recommendations are used] | [premature, costly, unsupported, too broad] | [Would GAL, evaluator, or another process fit better?] |

### Special Master / Parenting Coordinator

| relief_id | matter | relief_type | relief_requested | status | status_quo_effect | legal_authority_status | key_facts_needed | source_ids | proposed_order_terms | deadline_or_mechanism | enforcement_issue | attack_surface | attorney_questions |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| R-SM-1 | [MATTER] | Special Master / Parenting Coordinator | [Appointment and scope of authority] | [planned/pending/appointed/active] | [preserves_status_quo/restores_status_quo/changes_status_quo/unclear] | [supplied/verified_by_legal_authority_scholar/attorney_confirmation_required/unknown] | [repeated implementation disputes, need for neutral administration] | [DOC-006, DOC-022] | [scope, tie-break power if any, record access, reporting] | [decision deadlines / meeting schedule] | [what happens if a recommendation is rejected] | [too much authority requested, unclear mechanism, cost concerns] | [What powers are legally available and strategically wise?] |

### Custody Evaluation / Psychological Evaluation / Expert Evaluation

| relief_id | matter | relief_type | relief_requested | status | status_quo_effect | legal_authority_status | key_facts_needed | source_ids | proposed_order_terms | deadline_or_mechanism | enforcement_issue | attack_surface | attorney_questions |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| R-EVAL-1 | [MATTER] | Evaluation / expert | [Custody evaluation, psychological evaluation, vocational expert, financial expert, etc.] | [planned/pending/ordered/completed] | [preserves_status_quo/restores_status_quo/changes_status_quo/unclear] | [supplied/verified_by_legal_authority_scholar/attorney_confirmation_required/unknown] | [what question needs expertise and why the existing record is insufficient] | [DOC-008, DOC-019] | [scope, evaluator access, payment allocation, deadlines] | [appointment and report deadlines] | [noncooperation or scope disputes] | [premature, overreach, wrong expert, weak need showing] | [What exact gap will this expert fill?] |

### Medical / Therapy / School Decision-Making

| relief_id | matter | relief_type | relief_requested | status | status_quo_effect | legal_authority_status | key_facts_needed | source_ids | proposed_order_terms | deadline_or_mechanism | enforcement_issue | attack_surface | attorney_questions |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| R-MED-1 | [MATTER] | Medical / therapy / school decision-making | [Notice requirement, joint decision path, provider selection process, records access, etc.] | [planned/pending/granted/denied] | [preserves_status_quo/restores_status_quo/changes_status_quo/unclear] | [supplied/verified_by_legal_authority_scholar/attorney_confirmation_required/unknown] | [who decided what, what notice was or was not given, impact on child] | [DOC-009, DOC-013] | [specific approval process, records access, deadlines, tie-break path] | [response within X days / default rule] | [what happens if parties disagree] | [vague process, third-party obligations, authority concerns] | [Is the requested process workable for providers and schools?] |

### Parent-Time / Exchange Logistics

| relief_id | matter | relief_type | relief_requested | status | status_quo_effect | legal_authority_status | key_facts_needed | source_ids | proposed_order_terms | deadline_or_mechanism | enforcement_issue | attack_surface | attorney_questions |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| R-PT-1 | [MATTER] | Parent-time / exchange logistics | [Exchange location, timing, supervision, communication rule, makeup time] | [planned/pending/granted/denied] | [preserves_status_quo/restores_status_quo/changes_status_quo/unclear] | [supplied/verified_by_legal_authority_scholar/attorney_confirmation_required/unknown] | [missed exchanges, safety or logistics problems, child needs, prior practice] | [DOC-005, DOC-015] | [clear exchange terms and exceptions] | [start date / holiday exception / notice deadline] | [late exchange or cancellation mechanism] | [too much detail, too little detail, hidden modification] | [Would narrower logistics relief solve the real problem?] |

### Financial Support / Reimbursement

| relief_id | matter | relief_type | relief_requested | status | status_quo_effect | legal_authority_status | key_facts_needed | source_ids | proposed_order_terms | deadline_or_mechanism | enforcement_issue | attack_surface | attorney_questions |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| R-FIN-1 | [MATTER] | Financial support / reimbursement | [Support adjustment, reimbursement, allocation, arrears payment, etc.] | [planned/pending/granted/denied] | [preserves_status_quo/restores_status_quo/changes_status_quo/unclear] | [supplied/verified_by_legal_authority_scholar/attorney_confirmation_required/unknown] | [income or expense records, payment history, receipts, logs, declarations] | [DOC-010, DOC-020] | [payment amount, due date, documentation rule] | [monthly due date / reimbursement within X days] | [what counts as proof and when it is due] | [missing records, unclear math, entitlement questions] | [What legal authority or calculation support is still missing?] |

### Attorney Fees / Sanctions

| relief_id | matter | relief_type | relief_requested | status | status_quo_effect | legal_authority_status | key_facts_needed | source_ids | proposed_order_terms | deadline_or_mechanism | enforcement_issue | attack_surface | attorney_questions |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| R-FEES-1 | [MATTER] | Attorney fees / sanctions | [Fees, costs, sanctions, or related relief] | [planned/pending/granted/denied] | [preserves_status_quo/restores_status_quo/changes_status_quo/unclear] | [supplied/verified_by_legal_authority_scholar/attorney_confirmation_required/unknown] | [conduct, avoidable costs, declarations, invoices, procedural history] | [DOC-012, DOC-023] | [amount or process for amount determination] | [fee declaration deadline / payment deadline] | [collection or nonpayment mechanism] | [punitive framing, weak basis, authority gap] | [Should counsel ask now, later, or only preserve the issue?] |

### Discovery / Records Production

| relief_id | matter | relief_type | relief_requested | status | status_quo_effect | legal_authority_status | key_facts_needed | source_ids | proposed_order_terms | deadline_or_mechanism | enforcement_issue | attack_surface | attorney_questions |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| R-DISC-1 | [MATTER] | Discovery / records production | [Records request, subpoena-related relief, production order] | [planned/pending/granted/denied] | [preserves_status_quo/restores_status_quo/changes_status_quo/unclear] | [supplied/verified_by_legal_authority_scholar/attorney_confirmation_required/unknown] | [what records matter, why they matter, prior requests, nonproduction] | [DOC-016, DOC-024] | [specific production categories and deadlines] | [produce within X days / authorization form deadline] | [privacy objections, third-party compliance, scope fights] | [too broad, privacy issues, authority-sensitive request] | [What narrower request is more likely to be granted?] |

### Protective Order / Stalking Injunction Response

| relief_id | matter | relief_type | relief_requested | status | status_quo_effect | legal_authority_status | key_facts_needed | source_ids | proposed_order_terms | deadline_or_mechanism | enforcement_issue | attack_surface | attorney_questions |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| R-PO-1 | [MATTER] | Protective order / stalking injunction response | [Dismissal, narrowing, no-contact clarification, practical carve-out] | [planned/pending/granted/denied] | [preserves_status_quo/restores_status_quo/changes_status_quo/unclear] | [supplied/verified_by_legal_authority_scholar/attorney_confirmation_required/unknown] | [alleged incidents, contradictory records, logistics impact, existing parenting orders] | [DOC-025, DOC-026] | [contact mechanism, exchange carve-out, records access, hearing logistics] | [effective pending hearing / by date] | [conflict with existing family-law orders] | [minimizing serious allegations, asking for relief broader than needed] | [What narrow operational relief protects both compliance and parenting logistics?] |

### Settlement / Stipulation / Proposed Order

| relief_id | matter | relief_type | relief_requested | status | status_quo_effect | legal_authority_status | key_facts_needed | source_ids | proposed_order_terms | deadline_or_mechanism | enforcement_issue | attack_surface | attorney_questions |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| R-STIP-1 | [MATTER] | Settlement / stipulation / proposed order | [Agreed process, temporary arrangement, document exchange, narrowing stipulation] | [planned/proposed/signed/rejected] | [preserves_status_quo/restores_status_quo/changes_status_quo/unclear] | [supplied/verified_by_legal_authority_scholar/attorney_confirmation_required/unknown] | [what is agreed, what remains disputed, what needs to be operationalized] | [DOC-027, DOC-028] | [clean agreed terms, deadlines, fallback process] | [effective date / expiration / review date] | [ambiguity, missing signatures, third-party implementation] | [agreement says one thing, order says another] | [What terms need to be made more enforceable before filing?] |

## Relief Drafting Checklist

For each requested relief item, ask:

1. What exactly should the court order?
2. Who must do what?
3. By when?
4. How is notice given?
5. What documents or records must be produced?
6. What happens if parties disagree?
7. Is the relief narrow enough to be granted?
8. Is the relief broad enough to solve the problem?
9. Does it preserve, restore, or change status quo?
10. Does legal authority need verification?
11. Does the proposed order match the motion?

## Relief-to-Evidence Mapping

Use this section to connect each relief item to the facts and sources that actually support it.

### Relief Mapping Template

#### [RELIEF_ID]
- relief requested:
- source-backed facts needed:
- source IDs:
- authority status:
- likely opposing attack:
- proposed order mechanics:
- missing source documents:
- attorney questions:

### Example Generic Mapping

#### R-MED-1
- relief requested: advance notice and joint process for nonemergency provider changes
- source-backed facts needed: prior unilateral provider change, missing notice, resulting confusion for school or provider
- source IDs: DOC-009, DOC-013
- authority status: attorney_confirmation_required
- likely opposing attack: too broad, not tied to actual problem, impossible for provider to administer
- proposed order mechanics: written notice through OFW, response deadline, default if no agreement, emergency exception
- missing source documents: provider note, OFW message thread
- attorney questions: does the requested process match the current custody language?
