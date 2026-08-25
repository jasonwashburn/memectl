## Why

Managed memes can be created and listed, but users cannot remove stale or unwanted records from their local inventory. A delete command completes the kubectl-inspired resource lifecycle without implying that memectl can remove hosted Imgflip images.

## What Changes

- Add `memectl delete meme <name> [<name>...]` to remove one or more named managed-meme records from the selected local inventory.
- Process each requested name as an independent local inventory operation: delete locally present names, report absent names (including repeated names after their first deletion) as not found, and return a non-zero result when any name is absent. Successful earlier requests remain applied when a later request fails.
- Preserve the inventory's atomic-update and validation guarantees for each individual operation, including a valid empty versioned document after its final record is deleted and accurate reporting when replacement succeeds but durable persistence cannot be confirmed.
- Document that deletion is local only and does not contact Imgflip or affect the hosted image.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `meme-inventory`: Add lifecycle requirements for deleting named local meme records.

## Impact

- Affected commands: root command wiring and a new `delete` command group with a `meme` resource command.
- Affected local persistence: the inventory store gains an atomic record-removal operation.
- Affected documentation: README usage explains local-only deletion.
- No Imgflip API calls, credentials, remote deletion behavior, or new dependencies are introduced.
