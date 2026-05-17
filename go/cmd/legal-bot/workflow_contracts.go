package main

import (
	"fmt"
	"strings"
)

type workflowSection struct {
	Heading     string
	Description string
}

type workflowOutputContract struct {
	Type                      string
	Purpose                   string
	FinalPacketTitle          string
	DecisionType              string
	RecommendedTracks         []string
	RequiredSections          []workflowSection
	Rules                     []string
	ExternalLanguageRules     []string
	AttorneyQuestionRules     []string
	DoNotChaseRules           []string
	DraftExternalLanguageMode string
}

var workflowMetadataFields = []string{
	"Workflow:",
	"Matter:",
	"Audience:",
	"Urgency:",
	"Recommended Track:",
	"Confidence:",
	"External-Use Readiness:",
	"Source Strength:",
	"Attorney Review Needed:",
}

var workflowSourceStrengthTaxonomy = []string{
	"Locked / source-supported fact: supported by source documents, orders, filed records, preserved messages, emails, or third-party records.",
	"User recollection / user assertion: The user says it happened, but no source anchor is currently available in the matter record.",
	"Inference: a reasonable conclusion from source-supported facts, but not itself directly proven.",
	"Disputed allegation: alleged by one side or in a filing/message/record, but not accepted as true without independent support.",
	"Unsupported claim: strategically relevant, but lacking a present source.",
	"Legal proposition needing authority: a legal rule, standard, procedural claim, or remedy theory that still needs statute, rule, case, or order support.",
	"Attorney-only / internal strategy: useful for private legal strategy, but not safe for OFW, provider, school, Special Master, opposing counsel, or court-facing use without attorney review.",
	"External-safe fact: a fact that is both source-supported and phrased safely for the intended external audience.",
}

var workflowGlobalRules = []string{
	"Begin with a short Bottom Line recommendation and choose a clear Recommended Track.",
	"Answer the user's actual job instead of merely summarizing prior agents.",
	"Separate source-supported facts from user recollection, inference, disputed allegations, unsupported claims, legal propositions needing authority, and attorney-only strategy.",
	"Identify the strongest evidence, weakest evidence, evidence gaps, and the exact source that would materially change the recommendation.",
	"Identify proportionality, optics, escalation, and credibility risk.",
	"Identify legal-authority gaps and route unresolved legal questions to attorney questions instead of pretending they are resolved.",
	"Separate internal strategy from external-safe language.",
	"Do not convert internal strategy into accusations, do not overclaim, and do not use generic filler warnings instead of making a recommendation.",
	"If external-safe language is not warranted, say so plainly.",
}

func allWorkflowOutputContracts() []workflowOutputContract {
	return []workflowOutputContract{
		{
			Type:             workflowTypeTriage,
			Purpose:          "Analyze a new or messy situation and decide what to do next.",
			FinalPacketTitle: "Triage Decision Packet",
			DecisionType:     "Choose the strongest next track and explain why weaker tracks should wait, narrow, or stop.",
			RecommendedTracks: []string{
				"pursue now",
				"preserve for pattern",
				"limited OFW/provider/school message",
				"counsel-first",
				"Special Master",
				"court/MTE/TRO/GAL",
				"evidence packet first",
				"journal only",
				"let it go",
				"do not send yet",
			},
			RequiredSections: []workflowSection{
				{Heading: "Bottom Line", Description: "Choose one Recommended Track and say why."},
				{Heading: "What Happened", Description: "Neutral summary that separates known facts from interpretation."},
				{Heading: "Source Strength", Description: "Use the source-strength taxonomy and say what is locked, weak, or missing."},
				{Heading: "Why It Matters", Description: "Explain the legal, co-parenting, or practical significance."},
				{Heading: "Options Matrix", Description: "Compare do nothing, journal only, preserve, limited message, counsel-first, and SM/court options."},
				{Heading: "Do Not Chase Check", Description: "Say whether the issue is too small, premature, credibility-negative, or stronger if preserved."},
				{Heading: "Evidence Needed", Description: "List the precise source items that would strengthen or change the recommendation."},
				{Heading: "Recommended Next Step", Description: "Give one to three concrete actions."},
				{Heading: "Counsel Questions", Description: "Include only when attorney judgment is needed."},
				{Heading: "External-Safe Language", Description: "Include only if a limited external message is actually recommended."},
			},
			Rules: []string{
				"Do Not Chase is mandatory analysis, not an optional cautionary note.",
				"The Options Matrix must include a genuine do-nothing or let-it-go path when that path is defensible.",
				"If no external message is recommended, say 'No external message recommended yet.'",
			},
			ExternalLanguageRules: []string{
				"Only provide external-safe language if the recommendation is a limited record-building message or similarly narrow communication step.",
				"External-safe language must use source-supported facts and low-accusation phrasing.",
			},
			AttorneyQuestionRules: []string{
				"Use counsel questions only where attorney judgment, authority, or forum choice genuinely matters.",
			},
			DoNotChaseRules: []string{
				"Say plainly if chasing this now would hurt credibility, invite backlash, or distract from stronger issues.",
			},
			DraftExternalLanguageMode: "conditional",
		},
		{
			Type:             workflowTypeCounselBrief,
			Purpose:          "Produce a Kaitlyn/Jessika-ready attorney brief or call agenda.",
			FinalPacketTitle: "Counsel Brief Packet",
			DecisionType:     "Choose the strongest attorney ask, issue framing, and evidence request.",
			RecommendedTracks: []string{
				"send concise brief now",
				"send call agenda",
				"send issue list",
				"send evidence request list",
				"collect evidence first",
				"hold for a stronger moment",
			},
			RequiredSections: []workflowSection{
				{Heading: "Subject / Proposed Email Title", Description: "Short and useful subject line."},
				{Heading: "Short Version", Description: "Three to six bullets max."},
				{Heading: "Why I Think This Matters", Description: "Candid but not theatrical."},
				{Heading: "Key Facts", Description: "Use source-strength labels."},
				{Heading: "Evidence Already Available", Description: "List available source anchors."},
				{Heading: "Evidence Missing", Description: "List exact missing sources or records."},
				{Heading: "Legal / Strategy Questions for Counsel", Description: "Frame as attorney questions, not directives."},
				{Heading: "Options I See", Description: "Concise options with risks and benefits."},
				{Heading: "Recommended Ask", Description: "What the user should ask counsel to decide or do."},
				{Heading: "Attorney-Ready Draft", Description: "Include a concise email or call agenda when requested."},
			},
			Rules: []string{
				"Counsel tone may be candid and strategic.",
				"Do not make the user sound like he is instructing counsel how to practice law.",
			},
			ExternalLanguageRules: []string{
				"The Attorney-Ready Draft may be more candid than OFW/provider/school language, but it must still distinguish source-supported fact from interpretation.",
			},
			AttorneyQuestionRules: []string{
				"Questions must sound like 'Would it help to...?', 'Is this worth...?', 'Do we need...?', 'Would you want...?', or 'Is there any risk that...?'",
			},
			DoNotChaseRules: []string{
				"If the issue is not worth spending counsel time on yet, say so and explain what would change that.",
			},
			DraftExternalLanguageMode: "allowed",
		},
		{
			Type:             workflowTypeSMTriage,
			Purpose:          "Decide whether an issue is ready and appropriate for the Special Master.",
			FinalPacketTitle: "Special Master Triage Packet",
			DecisionType:     "Decide ripeness, scope fit, requested directive, and whether to wait, narrow, or reroute the issue.",
			RecommendedTracks: []string{
				"SM-ready now",
				"SM likely, but one record-building step first",
				"counsel-first",
				"preserve for pattern",
				"too small / premature",
				"court/MTE/TRO issue instead",
				"not recommended",
			},
			RequiredSections: []workflowSection{
				{Heading: "Bottom Line", Description: "Choose a clear SM readiness posture."},
				{Heading: "Issue Framed for the Special Master", Description: "Neutral one-paragraph framing."},
				{Heading: "Scope / Authority Fit", Description: "Assess whether this belongs within likely SM scope."},
				{Heading: "Ripeness", Description: "Assess resolution attempts, OFW record, directive clarity, and evidence cleanliness."},
				{Heading: "Requested Directive", Description: "Draft a precise directive if ready, or explain what becomes appropriate later."},
				{Heading: "Evidence Needed", Description: "List the exact missing records or anchors."},
				{Heading: "Optics / Credibility Risk", Description: "Assess reasonableness and reactivity risk."},
				{Heading: "Do Not Chase Check", Description: "Say if this is too small, premature, or better preserved."},
				{Heading: "Recommended Next Step", Description: "Give one to three actions."},
				{Heading: "Draft SM Message", Description: "Only include if SM-ready or nearly ready."},
			},
			Rules: []string{
				"Special Master analysis must stay child-focused, operational, and directive-oriented.",
				"Do not ask the Special Master to function as therapist, judge, or general referee unless scope clearly supports it.",
			},
			ExternalLanguageRules: []string{
				"If a Draft SM Message is included, it must be concise, professional, evidence-tethered, and directive-focused.",
				"Do not use inflammatory labels, unsupported motive claims, or therapy-style framing.",
			},
			AttorneyQuestionRules: []string{
				"Use counsel questions when scope, authority, or escalation forum remains unclear.",
			},
			DoNotChaseRules: []string{
				"Say directly if this issue is too small or too messy to spend Special Master capital on right now.",
			},
			DraftExternalLanguageMode: "conditional",
		},
		{
			Type:             workflowTypeDraftMessage,
			Purpose:          "Draft or evaluate an external message for OFW/Emily, Special Master, provider, school, counsel, or court-facing use.",
			FinalPacketTitle: "External Message Packet",
			DecisionType:     "Decide whether to send, revise, hold, reroute to counsel, or not send a message at all.",
			RecommendedTracks: []string{
				"send",
				"revise",
				"do not send yet",
				"counsel-first",
				"no message recommended",
			},
			RequiredSections: []workflowSection{
				{Heading: "Bottom Line", Description: "Say send, revise, do not send yet, counsel-first, or no message recommended."},
				{Heading: "Audience and Tone", Description: "Name the audience and appropriate tone."},
				{Heading: "Message Goal", Description: "One sentence."},
				{Heading: "Risk Check", Description: "Assess DARVO, overstatement, accusation, optics, ambiguity, and escalation risks."},
				{Heading: "Facts Safe to Say", Description: "List only external-safe facts."},
				{Heading: "Facts or Claims to Avoid", Description: "Exclude internal strategy and unsupported claims."},
				{Heading: "Recommended Message", Description: "Provide the actual message only if sending is recommended."},
				{Heading: "Optional Shorter Version", Description: "Include for longer OFW/provider/school drafts."},
				{Heading: "Preservation Notes", Description: "Say what the message preserves in the record."},
			},
			Rules: []string{
				"Never write external language that relies on unsupported internal conclusions.",
				"OFW/Emily messaging should be brief, child-focused, low-DARVO, and record-building.",
			},
			ExternalLanguageRules: []string{
				"OFW/Emily: brief, calm, child-focused, and record-building.",
				"Special Master: directive-focused and professional.",
				"Provider/school: neutral, access/request/factual clarification focused.",
				"Counsel: candid and strategic.",
				"Court-facing: cautious and source-tethered.",
			},
			AttorneyQuestionRules: []string{
				"If the message is court-sensitive, authority-sensitive, or likely to escalate heavily, route unanswered issues into counsel questions or counsel-first guidance.",
			},
			DoNotChaseRules: []string{
				"If sending the message mainly vents, overstates, or creates quote risk, say 'No message recommended yet' or 'Do not send yet.'",
			},
			DraftExternalLanguageMode: "allowed",
		},
		{
			Type:             workflowTypeEvidencePacket,
			Purpose:          "Identify, structure, and prioritize evidence needed for an issue.",
			FinalPacketTitle: "Evidence Packet Plan",
			DecisionType:     "Choose the strongest collection and packet-building plan before escalation.",
			RecommendedTracks: []string{
				"attorney packet",
				"SM packet",
				"court packet",
				"journal support",
				"internal triage",
			},
			RequiredSections: []workflowSection{
				{Heading: "Packet Purpose", Description: "Attorney packet, SM packet, court packet, journal support, or internal triage."},
				{Heading: "Issue Summary", Description: "Neutral, concise summary."},
				{Heading: "Evidence Inventory", Description: "Table with source type, anchor, support value, strength, weakness, and needed-for fields."},
				{Heading: "Missing Evidence", Description: "List exact missing items and likely owner/location."},
				{Heading: "Chronology Needed", Description: "List the key dates and events that must be proven."},
				{Heading: "Fact Map", Description: "Table mapping facts to best source, current support, gap, and priority."},
				{Heading: "Exhibits / Attachments Candidate List", Description: "Include only if useful."},
				{Heading: "Risk / Weaknesses", Description: "Say what a neutral or opponent would attack."},
				{Heading: "Recommended Packet Structure", Description: "Suggested folder/file structure or packet order."},
				{Heading: "Next Actions", Description: "Concrete collection steps."},
			},
			Rules: []string{
				"Do not treat AI summaries or journal notes as evidence.",
				"Identify the best available source, not merely the most convenient source.",
			},
			ExternalLanguageRules: []string{
				"This workflow is usually internal-facing. Do not manufacture external language unless the packet purpose specifically requires it.",
			},
			AttorneyQuestionRules: []string{
				"Use counsel questions only where packet design, exhibit selection, or missing authority changes what should be collected.",
			},
			DoNotChaseRules: []string{
				"Say if the evidence problem means escalation should wait until cleaner records exist.",
			},
			DraftExternalLanguageMode: "prohibited",
		},
		{
			Type:             workflowTypePatternReview,
			Purpose:          "Determine whether recurring events add up to a meaningful pattern.",
			FinalPacketTitle: "Pattern Review Packet",
			DecisionType:     "Decide whether a real pattern exists, whether weak examples should be dropped, and what use is justified now.",
			RecommendedTracks: []string{
				"strong pattern",
				"emerging pattern",
				"weak/no pattern yet",
				"not worth pursuing",
				"preserve only",
			},
			RequiredSections: []workflowSection{
				{Heading: "Bottom Line", Description: "Choose pattern strength and recommended use."},
				{Heading: "Pattern Theory", Description: "Describe the alleged recurring pattern neutrally."},
				{Heading: "Strongest Examples", Description: "Table with date, event, source, significance, strength, and caveat."},
				{Heading: "Weak or Risky Examples", Description: "Table showing what should be preserved, dropped, or used cautiously."},
				{Heading: "Pattern vs Isolated Noise", Description: "Explain why this is or is not more than ordinary co-parenting friction."},
				{Heading: "Legal / Co-Parenting Significance", Description: "Explain the significance of the pattern if it exists."},
				{Heading: "Opposing / Neutral Read", Description: "Say how Emily, opposing counsel, or a neutral might frame it."},
				{Heading: "Do Not Chase Check", Description: "Say whether raising the pattern now helps or hurts."},
				{Heading: "Recommended Use", Description: "Choose counsel brief, SM issue, court support, journal only, preserve, or drop."},
				{Heading: "Evidence Needed", Description: "List missing anchors and sources."},
			},
			Rules: []string{
				"Pattern review should make the User more selective, not more reactive.",
				"Do not inflate weak examples just because they fit the theory.",
			},
			ExternalLanguageRules: []string{
				"This workflow is usually internal-facing. If external-safe language is included, it must be narrow and evidence-tethered.",
			},
			AttorneyQuestionRules: []string{
				"Use counsel questions when forum choice, pattern significance, or legal threshold remains uncertain.",
			},
			DoNotChaseRules: []string{
				"Say explicitly when the pattern is too thin or too messy to raise now.",
			},
			DraftExternalLanguageMode: "prohibited",
		},
		{
			Type:             workflowTypeFactLock,
			Purpose:          "Separate proven facts from recollection, inference, claims, allegations, and legal propositions.",
			FinalPacketTitle: "Fact Lock Packet",
			DecisionType:     "Decide what can be said safely now, what needs sourcing, and what should remain internal-only.",
			RecommendedTracks: []string{
				"safe to use now",
				"partially ready",
				"internal only",
				"counsel-first",
				"needs more sourcing",
			},
			RequiredSections: []workflowSection{
				{Heading: "Fact-Lock Summary", Description: "Say what can safely be said now."},
				{Heading: "Locked Facts", Description: "Table with fact, source, anchor, and safe external wording."},
				{Heading: "Partially Supported Facts", Description: "Table with gaps and how to lock them."},
				{Heading: "User Recollection / Unsourced Assertions", Description: "Table with needed sources and external-use risk."},
				{Heading: "Inferences", Description: "Table with supporting facts, confidence, and safer wording."},
				{Heading: "Disputed Allegations", Description: "Table with who says it, source, and caveat."},
				{Heading: "Legal Propositions Needing Authority", Description: "Table with needed authority and attorney question."},
				{Heading: "External-Use Readiness", Description: "Assess counsel, OFW, Special Master, provider/school, and court."},
				{Heading: "Recommended Next Step", Description: "Say how to lock the most important gaps."},
			},
			Rules: []string{
				"This workflow must be strict.",
				"Do not let a persuasive story override source weakness.",
			},
			ExternalLanguageRules: []string{
				"Only use the Safe External Wording column for facts that are both source-supported and audience-safe.",
			},
			AttorneyQuestionRules: []string{
				"Route unresolved authority issues and high-stakes external-use questions to counsel questions.",
			},
			DoNotChaseRules: []string{
				"Say when the safest next step is to wait, source more, or keep the issue internal.",
			},
			DraftExternalLanguageMode: "conditional",
		},
		{
			Type:             workflowTypeAuthorityCheck,
			Purpose:          "Identify legal authority gaps and questions needing live research or attorney confirmation.",
			FinalPacketTitle: "Authority Check Packet",
			DecisionType:     "Decide whether the authority posture is supported, unclear, risky, unsupported, or attorney-research dependent.",
			RecommendedTracks: []string{
				"likely supported but needs citation",
				"unclear",
				"risky",
				"unsupported",
				"attorney research needed",
				"order/decree-specific issue",
			},
			RequiredSections: []workflowSection{
				{Heading: "Bottom Line", Description: "Choose the current authority posture."},
				{Heading: "Legal Questions Presented", Description: "List the legal questions clearly."},
				{Heading: "Known Authority Sources in the Matter", Description: "List supplied orders, decree terms, rules, statutes, cases, or notes."},
				{Heading: "Authority Gaps", Description: "Table with claim, needed authority, importance, and verifier."},
				{Heading: "Currentness / Citator Status", Description: "Explicitly say whether currentness/citator status was verified."},
				{Heading: "Attorney Questions", Description: "List the attorney questions."},
				{Heading: "Safe Working Language", Description: "Provide cautious internal/counsel-facing wording only."},
				{Heading: "Do Not Overclaim", Description: "List claims that should not be made yet."},
			},
			Rules: []string{
				"Do not invent statutes, rules, cases, standards, or holdings.",
				"Do not imply authority is current unless currentness was actually checked.",
			},
			ExternalLanguageRules: []string{
				"Do not produce court-ready legal conclusions unless supported authority is actually present.",
			},
			AttorneyQuestionRules: []string{
				"Attorney questions are required when authority is missing, unclear, or not verified.",
			},
			DoNotChaseRules: []string{
				"Say when the user should not make an authority-based claim yet because the legal support is too weak or unverified.",
			},
			DraftExternalLanguageMode: "prohibited",
		},
		{
			Type:             workflowTypePrep,
			Purpose:          "Prepare for an attorney call, hearing, Special Master conference, mediation, provider/school meeting, or similar event.",
			FinalPacketTitle: "Prep Packet",
			DecisionType:     "Choose the clearest event objective, strongest lead facts, and best use of limited time.",
			RecommendedTracks: []string{
				"lead with this",
				"ask these questions",
				"bring these records",
				"avoid this lane",
				"follow up this way",
			},
			RequiredSections: []workflowSection{
				{Heading: "Prep Objective", Description: "Say what the User needs from the event."},
				{Heading: "Bottom Line", Description: "State the main recommendation and desired outcome."},
				{Heading: "Key Facts to Lead With", Description: "Short source-labeled list."},
				{Heading: "Decisions Needed", Description: "Say what must be decided during the event."},
				{Heading: "Questions to Ask", Description: "Group by audience."},
				{Heading: "Likely Pushback / Attacks", Description: "Say how opposition, Emily, neutrals, or providers may push back."},
				{Heading: "Best Responses", Description: "Concise response themes."},
				{Heading: "Evidence to Have Ready", Description: "Precise source list."},
				{Heading: "What Not to Say", Description: "Unsupported claims, emotional framing, internal-only strategy, and overbroad accusations."},
				{Heading: "Follow-Up Actions", Description: "List post-event actions."},
			},
			Rules: []string{
				"Prep should be practical and scannable.",
				"Attorney prep may be candid. External prep must be cautious.",
			},
			ExternalLanguageRules: []string{
				"If audience-specific wording is included, keep it narrow and usable in a live conversation.",
			},
			AttorneyQuestionRules: []string{
				"Use questions to maximize limited meeting time and surface decision points, not to repeat background noise.",
			},
			DoNotChaseRules: []string{
				"Say when a tangent, accusation, or weak issue should not be raised in the event.",
			},
			DraftExternalLanguageMode: "conditional",
		},
		{
			Type:             workflowTypeJournalEntry,
			Purpose:          "Create a neutral, source-aware Divorce & Parenting Journal entry.",
			FinalPacketTitle: "Journal Entry Packet",
			DecisionType:     "Decide how to preserve the event neutrally, with caveats and usable future indexing notes.",
			RecommendedTracks: []string{
				"journal only",
				"journal and preserve for pattern",
				"journal and flag for counsel",
				"journal and gather missing source",
			},
			RequiredSections: []workflowSection{
				{Heading: "Suggested Title", Description: "Short and dated if possible."},
				{Heading: "Date / Date Range", Description: "Exact or approximate."},
				{Heading: "Record Type", Description: "OFW, Medical, School, Therapy, Legal, Special Master, Court, Parent-time, Provider Access, Financial, or General."},
				{Heading: "Priority", Description: "High, medium, or low."},
				{Heading: "Category / Related Issue", Description: "Categorize the entry."},
				{Heading: "Children", Description: "Aya, Lydia, both, neither, or unclear."},
				{Heading: "Key People", Description: "List the key people."},
				{Heading: "Summary", Description: "Neutral two to four sentence summary."},
				{Heading: "What Happened", Description: "Factual, source-aware account."},
				{Heading: "Legal / Co-Parenting Significance", Description: "Explain why it might matter."},
				{Heading: "Neutral Evaluator Notes", Description: "Say how a neutral might see it, including weaknesses."},
				{Heading: "Evidence / Source Anchors", Description: "List available and missing anchors."},
				{Heading: "Caveats / Verification Needed", Description: "Preserve uncertainty."},
				{Heading: "Related Events", Description: "Link to similar incidents if known."},
				{Heading: "Import Notes", Description: "Include source strength, optics/risk, verification needs, and privileged/internal-only status."},
			},
			Rules: []string{
				"The journal is an index and analysis layer, not source evidence.",
				"Do not argue the case in the journal entry.",
			},
			ExternalLanguageRules: []string{
				"This workflow is not for external messaging. Do not produce external-facing draft language.",
			},
			AttorneyQuestionRules: []string{
				"Use counsel questions only if the event clearly needs later legal follow-up.",
			},
			DoNotChaseRules: []string{
				"If the right move is only to preserve the event without escalation, say so plainly.",
			},
			DraftExternalLanguageMode: "prohibited",
		},
	}
}

func workflowOutputContractForType(workflowType string) (workflowOutputContract, bool) {
	for _, contract := range allWorkflowOutputContracts() {
		if contract.Type == normalizeRouteName(workflowType) {
			return contract, true
		}
	}
	return workflowOutputContract{}, false
}

func renderWorkflowContractSummary(contract workflowOutputContract) string {
	var b strings.Builder
	b.WriteString("=== WORKFLOW PURPOSE AND OUTPUT CONTRACT SUMMARY ===\n")
	b.WriteString("workflow_type: " + contract.Type + "\n")
	b.WriteString("purpose: " + contract.Purpose + "\n")
	b.WriteString("final_packet_title: " + contract.FinalPacketTitle + "\n")
	b.WriteString("decision_type: " + contract.DecisionType + "\n")
	if len(contract.RecommendedTracks) > 0 {
		b.WriteString("recommended_tracks: " + strings.Join(contract.RecommendedTracks, " | ") + "\n")
	}
	b.WriteString("required_sections:\n")
	for _, section := range contract.RequiredSections {
		b.WriteString("- " + section.Heading + "\n")
	}
	b.WriteString("non_negotiable_rules:\n")
	for _, rule := range contract.Rules {
		b.WriteString("- " + rule + "\n")
	}
	if len(contract.DoNotChaseRules) > 0 {
		b.WriteString("do_not_chase_focus:\n")
		for _, rule := range contract.DoNotChaseRules {
			b.WriteString("- " + rule + "\n")
		}
	}
	b.WriteString("=== END WORKFLOW PURPOSE AND OUTPUT CONTRACT SUMMARY ===")
	return b.String()
}

func renderWorkflowContract(contract workflowOutputContract) string {
	var b strings.Builder
	b.WriteString("=== FULL WORKFLOW OUTPUT CONTRACT ===\n")
	b.WriteString("You are producing the final workflow packet. Follow this contract exactly. Do not merely summarize prior agents. Choose, prioritize, recommend, and name what not to chase.\n\n")
	b.WriteString("Workflow Type: " + contract.Type + "\n")
	b.WriteString("Purpose: " + contract.Purpose + "\n")
	b.WriteString("Final Packet Title: " + contract.FinalPacketTitle + "\n")
	b.WriteString("Decision Type: " + contract.DecisionType + "\n\n")

	b.WriteString("Required top metadata near the beginning of the final packet:\n")
	for _, field := range workflowMetadataFields {
		b.WriteString("- " + field + "\n")
	}

	b.WriteString("\nRecommended Track labels to choose from when relevant:\n")
	for _, track := range contract.RecommendedTracks {
		b.WriteString("- " + track + "\n")
	}

	b.WriteString("\nGlobal final packet rules:\n")
	for _, rule := range workflowGlobalRules {
		b.WriteString("- " + rule + "\n")
	}

	b.WriteString("\nSource-strength taxonomy:\n")
	for idx, item := range workflowSourceStrengthTaxonomy {
		b.WriteString(fmt.Sprintf("%d. %s\n", idx+1, item))
	}

	b.WriteString("\nWorkflow-specific required sections:\n")
	for _, section := range contract.RequiredSections {
		b.WriteString("## " + section.Heading + "\n")
		b.WriteString(section.Description + "\n\n")
	}

	if len(contract.Rules) > 0 {
		b.WriteString("Workflow-specific rules:\n")
		for _, rule := range contract.Rules {
			b.WriteString("- " + rule + "\n")
		}
		b.WriteString("\n")
	}

	if len(contract.ExternalLanguageRules) > 0 {
		b.WriteString("External-safe language rules:\n")
		for _, rule := range contract.ExternalLanguageRules {
			b.WriteString("- " + rule + "\n")
		}
		b.WriteString("\n")
	}

	if len(contract.AttorneyQuestionRules) > 0 {
		b.WriteString("Attorney-question rules:\n")
		for _, rule := range contract.AttorneyQuestionRules {
			b.WriteString("- " + rule + "\n")
		}
		b.WriteString("\n")
	}

	if len(contract.DoNotChaseRules) > 0 {
		b.WriteString("Do Not Chase rules:\n")
		for _, rule := range contract.DoNotChaseRules {
			b.WriteString("- " + rule + "\n")
		}
		b.WriteString("\n")
	}

	switch contract.DraftExternalLanguageMode {
	case "allowed":
		b.WriteString("External-safe draft language may be included when the facts and recommendation support sending it.\n")
	case "conditional":
		b.WriteString("External-safe draft language is conditional. If not recommended, say 'No external message recommended yet.'\n")
	case "prohibited":
		b.WriteString("Do not provide external-facing draft language in this workflow unless the contract explicitly calls for narrow internal/counsel-facing working language.\n")
	}

	b.WriteString("=== END FULL WORKFLOW OUTPUT CONTRACT ===")
	return b.String()
}
