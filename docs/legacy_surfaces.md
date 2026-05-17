# Legacy Surfaces

Legal-Bot originated from Vern-Bot. The repo still contains several Vern-era packages, wrappers, and compatibility layers.

Those legacy surfaces are not the current Legal-Bot family-law workflow contract.

## Current Command Contract

The active operator-facing Legal-Bot commands are:

- `legal-bot guide`
- `legal-bot intake`
- `legal-bot review`
- `legal-bot workflow`
- `legal-bot routes`
- `legal-bot run`
- `legal-bot setup`
- `legal-bot historian`
- `legal-bot tui`

Use `guide`, `intake`, `review`, and `workflow` for the current family-law workflow.

Use `routes` to inspect the active route catalog and the legacy inventory.

## Legacy Surface Inventory

| Surface | Status | What It Means | Use Instead |
|---|---|---|---|
| `bin/legal-bot-discovery` | retired | Legacy Vern discovery wrapper. Current Go CLI does not expose `discovery`. | `legal-bot guide`, `legal-bot intake`, `legal-bot review`, `legal-bot workflow` |
| `bin/legal-bot-discovery.cmd` | retired | Windows legacy Vern discovery wrapper. | `legal-bot guide`, `legal-bot intake`, `legal-bot review`, `legal-bot workflow` |
| `bin/legal-bot-council` | retired | Legacy VernHole council wrapper. Current Go CLI does not expose `hole`. | `legal-bot tui` only if you explicitly want the legacy utility UI |
| `bin/legal-bot-council.cmd` | retired | Windows legacy VernHole wrapper. | `legal-bot tui` only if you explicitly want the legacy utility UI |
| `bin/legal-bot-oracle` | retired | Legacy oracle wrapper. Current legal CLI does not expose `oracle`. | none in the active legal workflow; see repo history if needed |
| `bin/legal-bot-generate` | retired | Legacy persona-generation wrapper. Current legal CLI does not expose `generate`. | none in the active legal workflow |
| `bin/legal-bot-run` | partial | Compatibility wrapper for the active low-level run command. | `legal-bot run` |
| `bin/legal-bot-run.cmd` | partial | Windows compatibility wrapper for `run`. | `legal-bot run` |
| `bin/legal-bot-historian` | partial | Compatibility wrapper for the active historian utility. | `legal-bot historian` |
| `bin/install-legal-bot-cli` | partial | Binary install helper, not part of the legal workflow. | `cd go && go build -o bin/legal-bot ./cmd/legal-bot` |
| `bin/install-legal-bot-cli.cmd` | partial | Windows binary install helper. | `cd go && go build -o bin/legal-bot.exe ./cmd/legal-bot` |
| `go/internal/pipeline` | partial | Still used by `historian` and legacy-oriented utilities. | keep as compatibility infrastructure |
| `go/internal/council` | legacy | VernHole council support, not part of the active legal workflow. | active legal routes do not depend on it |
| `go/internal/generate` | legacy | Legacy persona-generation helpers. | active legal routes do not depend on it |
| `go/internal/vts` | legacy | Vern Task Spec helpers. | active legal routes do not depend on it |
| `go/internal/tobeads` | legacy | VTS-to-Beads import tooling. | active legal routes do not depend on it |
| `go/internal/tui` | partial | Still exposed through `legal-bot tui`, but oriented around legacy discovery/VernHole utilities. | `guide`, `intake`, `review`, `workflow` for legal work |

## Practical Boundary

- `guide`, `intake`, `review`, and `workflow` are the current legal workflow control plane.
- `run`, `setup`, `historian`, and `tui` remain utility commands.
- Discovery, council, oracle, generate, VTS, and VernHole concepts are retained for compatibility or historical continuity, not as the main Legal-Bot family-law operator surface.
