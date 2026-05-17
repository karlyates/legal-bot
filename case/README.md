# Case Knowledge Layer

This directory contains global case-control files that are loaded into prompts for intake, workflow, and review runs.

Legal-Bot produces AI-generated legal analysis, strategy recommendations, issue spotting, risk review, evidence analysis, and drafting support for an informed user. It may recommend actions, escalation paths, attorney questions, legal research issues, draft language, and litigation strategy options. Legal-Bot is not licensed counsel, may be wrong or incomplete, and outputs should be verified before being filed, served, sent externally, or relied on in a legal matter.

## Important Loader Behavior

Every root `.md` and `.txt` file in `case/` is auto-loaded into prompt context.

That means:

- `case/` is global prompt context
- changes here affect every matter
- large files here bloat every run
- raw evidence files do not belong here

Use `matters/<matter>/input/` for matter-specific source files.

## Files

| File | Purpose |
|---|---|
| `source_hierarchy.md` | source strength, evidence vs authority, and source discipline |
| `people_and_roles.md` | stable people, roles, professionals, and identifiers |
| `preferred_framing.md` | narrative posture, strategic framing, and issue emphasis |
| `language_calibration.md` | audience-sensitive language strength, not automatic softening |
| `relief_library.md` | relief tracking, order hooks, and relief-to-evidence mapping |

## How These Files Affect Workflows

- `source_hierarchy.md` affects fact strength and safe-use handling
- `people_and_roles.md` affects naming, identifiers, and role clarity
- `preferred_framing.md` affects how the system narrates and prioritizes issues
- `language_calibration.md` affects whether strong language is preserved, narrowed, moved, or removed
- `relief_library.md` affects escalation, requested directives, and relief alignment

## Shareability Labels

Use shareability labels when relevant:

- `internal_only`
- `attorney_safe`
- `court_safe_candidate`

These help keep attorney-only or strategy-only content out of court-facing or externally shared text.

## Practical Rules

- source documents remain the source of truth
- AI summaries are not evidence
- legal authority is not factual evidence
- filed allegations are not automatically true
- strong language is allowed when it is source-supported, strategically useful, and audience-appropriate

## Good Use of `case/`

Good:

- durable narrative guidance
- people and role definitions
- source hierarchy rules
- relief mapping
- language calibration

Bad:

- giant raw evidence dumps
- matter-specific record archives
- throwaway notes that should live in `matters/<matter>/`
