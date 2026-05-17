# Architecture Delta Report v4

## Scope

This delta report captures the meaningful architecture changes after the v3 baseline. The main shifts are workflow output contracts, contract injection into workflow prompts, and a cleaner active-vs-legacy repo identity split on the user-facing surface.

## What Changed After v3

### Workflow Output Contracts

- Added `go/cmd/legal-bot/workflow_contracts.go`.
- Added `docs/workflow_output_contracts.md`.
- Added `go/cmd/legal-bot/workflow_contracts_test.go`.
- Every workflow type now has an explicit contract.
- Contracts define the final packet shape, source-strength taxonomy, recommended-track language, external-safe rules, attorney-question prompts, and Do Not Chase guidance.

### Prompt Injection

- `go/cmd/legal-bot/workflow.go` now injects a concise contract summary into every workflow step.
- The final `managing-partner-final-synthesizer` step receives the full workflow output contract.
- This keeps the whole workflow aligned on the same final-packet expectations instead of letting the final step improvise.

### Final Synthesizer

- `agents/managing-partner-final-synthesizer.md` now reads like a decision-maker brief, not a recap role.
- It explicitly requires the agent to choose, prioritize, recommend, separate source-supported facts from inference/strategy, and include Do Not Chase analysis when appropriate.
- The embedded agent bundle reflects this prompt corpus, which is why `go/internal/embedded/agents_generated.go` is modified.

### Repo Identity Cleanup

- Active user-facing surfaces were cleaned of obvious stale Vern branding in the earlier cleanup pass.
- `REPO_IDENTITY_CLEANUP_AUDIT.md` and `go/cmd/legal-bot/repo_identity_test.go` exist as continuity and regression artifacts.
- Legacy wrappers, docs, and compatibility env names remain intentionally preserved.

## Verification

- Full Go tests passed.
- Build passed.
- `legal-bot routes`, `legal-bot routes --json`, and `legal-bot guide --list` all passed as non-LLM smoke checks.

## Current Assessment

- The repo is now in a better place architecturally than v3 because the final packet is no longer a soft convention; it is an explicit contract layer.
- The biggest remaining architectural risk is not missing plumbing, but whether the contract language feels right on a real matter once Karl runs a controlled sample.

## Next Likely Follow-Up

- Validate one real guided workflow packet on a controlled non-sensitive sample matter.
- If the output is too verbose, too narrow, or too generic, tune the contracts rather than reopening the routing architecture.
- Defer module-path cleanup and broad legacy-package deletion until after the packet shape feels right in practice.
