# Legacy Discovery Migration

## Current Product Contract

The active legal product contract is now:

- `legal-bot intake`
- `legal-bot workflow --type ...`
- `legal-bot review`
- `legal-bot review-draft`
- the multi-persona legal team simulator
- the case control layer in `case/`

The repo still contains older discovery, council, historian, oracle, TUI, and Vern-era compatibility surfaces. Those are not automatically bugs and should not be renamed casually.

## Active Legal Concepts

The active user-facing concepts are:

- matter intake
- new situation triage
- evidence packet builder
- counsel brief
- draft review
- Special Master triage
- external message drafting
- pattern review
- fact lock
- authority check
- hearing or call prep
- post-event journal entry

Internally, these route through the 17-persona legal panel and the managing partner final synthesis model.

## What Changed

The product is no longer framed mainly as:

- a generic draft reviewer
- a timid legal-analysis helper
- a pipeline collection the user has to memorize

It is now framed as:

- a local-first legal analysis and legal-operations assistant
- an AI legal team simulator
- a workflow-oriented tool for practical decisions

## Compatibility Surfaces Still Present

Examples still in the repo:

- `go/internal/council`
- `go/internal/pipeline`
- `go/internal/tui`
- `go/internal/generate`
- `go/internal/vts`
- older compatibility naming in tests, loaders, or historical paths

These should be treated as:

- active
- compatibility-only
- historical
- removable

only after actual inspection

## Safe Migration Guidance

1. Do not rename legacy surfaces just because the naming is old.
2. Confirm whether the path is active in the legal CLI.
3. Add aliases before removing active user-facing names.
4. Keep user-facing docs aligned with what is actually implemented.
5. Regenerate embedded assets after agent or config changes.
6. Run tests and build checks before claiming the migration is safe.

## Current Status

This pass updates the active legal product framing and workflow contract without trying to rewrite the entire legacy substrate.
