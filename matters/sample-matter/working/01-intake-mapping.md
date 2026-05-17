# Litigation Paralegal Intake

## Document Register
```csv
doc_id,filename,path,stable_file_key,document_type,document_subtype,date_of_document,date_of_event,filing_date,case_number,court,author_source,party_affiliation,recipients,source_tier,confidentiality,privilege_candidate,children_involved,issue_tags,status,duplicate_of,related_documents,notes
DOC-001,source-1.md,source-1.md,src_source-1_md_v1,user_note,intake_issue_summary,,2026-04-29,,,,"Parent B / user-side summary",user_side,,Tier 4,confidential,yes,yes,"joint_legal_custody;medical_notice;provider_appointment;psychiatric_referral;children;special_master_triage",needs_supporting_documents,,,"User summary describes two children's appointments on 2026-04-29, partial notice to Parent B for only one child, later phone-call note suggesting psychiatric referral, and absence of full medical records or billing confirmation."
```

## Intake Findings
1. **Missing source support for core event details**
   - category: `procedural`
   - severity: `High`
   - action: `Needs source document`
   - confidence: `High confidence`
   - shareability: `attorney_safe`
   - finding: The intake summary reports a two-child medical appointment event and a possible psychiatric referral, but the underlying records are not provided. Current support appears limited to one message and one phone-call note.

2. **Potential joint-decision / notice issue is not yet filing-ready**
   - category: `procedural`
   - severity: `Medium`
   - action: `Clarify`
   - confidence: `Medium confidence`
   - shareability: `attorney_safe`
   - finding: The summary suggests incomplete notice to a joint legal custodian, but the exact custody language, provider communications, and referral status are not in the source set. This is enough to preserve as an intake issue, not enough to present as a documented violation.

3. **Child-sensitive medical subject matter**
   - category: `procedural`
   - severity: `High`
   - action: `Ask attorney`
   - confidence: `High confidence`
   - shareability: `attorney_safe`
   - finding: The document references children’s medical appointments and a possible psychiatric referral. Any follow-on handling should preserve child-sensitive and medical confidentiality discipline.

4. **Strategy content mixed into factual intake note**
   - category: `procedural`
   - severity: `Low`
   - action: `Move`
   - confidence: `High confidence`
   - shareability: `internal_only`
   - finding: The document blends factual summary with next-step questions about counsel, pattern preservation, and Special Master routing. Those strategy questions should stay separate from any court-safe factual packet.

## Duplicate / Near-Duplicate Flags
- No duplicate or near-duplicate documents identified from the provided source set.

## Missing Companion Documents
- The actual message disclosing only one child’s appointment.
- The full message thread or platform metadata showing sender, recipient, and timestamp.
- The referral-provider phone-call note or call log.
- Full medical records for the 2026-04-29 appointments.
- Any referral order, referral request, or provider portal record confirming whether a psychiatric referral was initiated.
- Insurance EOB, billing record, or portal/billing screenshot if billing becomes relevant.
- The operative custody/order language if the notice issue may later be matched to a joint legal custody obligation.

## OCR / Readability Problems
- No OCR problem identified. The provided markdown source is readable.

## Privilege / Confidentiality Flags
- `DOC-001` should be treated as confidential.
- `DOC-001` is a privilege candidate because it includes attorney-directed next-step questions and strategy-routing thoughts.
- Child medical and mental-health-adjacent references increase sensitivity. Limit redistribution outside attorney-safe workflow unless necessary.

## Next Intake Steps
1. Preserve `DOC-001` as a Tier 4 intake summary, not as proof of the medical events.
2. Collect the underlying message and phone-call note and register them separately when available.
3. Obtain the operative order or decree sections governing joint legal custody and medical decision-making.
4. Request or collect the relevant provider and referral records before escalating the issue as a documented pattern.
5. If a Special Master process exists, confirm scope before routing; do not assume medical notice disputes fall within that process without the source order.

## Best Workflow Fit
The better next workflow is `triage`, followed by an `evidence packet` if the underlying message, call note, and medical/referral records can be gathered. On the current record, this appears too under-documented to escalate beyond issue preservation.