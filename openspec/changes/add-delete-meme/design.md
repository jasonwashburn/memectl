## Context

The inventory store currently loads managed memes and atomically appends a new record while holding a file lock. The command layer exposes `create meme` and `get memes`, and accepts a store interface that supports loading and adding records. See proposal.md for the motivation and the meme-inventory delta spec for required behavior.

## Goals / Non-Goals

**Goals:**
- Add a kubectl-shaped `delete meme <name> [<name>...]` command that removes one or more local records safely.
- Preserve validation, locking, atomic replacement, permissions, and inventory-version handling already used for additions.
- Make successful deletion and not-found outcomes clear without contacting Imgflip.

**Non-Goals:**
- Delete, hide, modify, or otherwise contact the Imgflip-hosted meme.
- Add confirmation prompts, label selection, or remote resource identifiers.
- Change the inventory document schema or existing create and list behavior.

## Decisions

### Model each deletion as an inventory operation

Add a single-record removal operation to the inventory store and expose it through the command store interface. Add and removal will use a shared exclusive-lock mutation helper that validates and loads the inventory while holding the lock, then atomically writes through the existing write path. The delete command will invoke removal once per requested name, in order. Each invocation deletes one present record or reports that name as absent, so repeated names naturally receive a not-found result after the first successful deletion.

Using a store operation rather than filtering a command-level `Load` result avoids a load-modify-write race with concurrent creation or deletion. Each operation acquires and releases the shared lock independently, intentionally allowing other create or delete operations to interleave between requested names. This mirrors individual API requests: successful earlier deletions remain applied when a later request is absent or fails. Reusing the existing write path retains atomic replacement and empty-document representation per operation.

Alternative considered: filter records in the command and introduce a generic save operation. This would leak inventory persistence details into the command layer and make lock ownership less clear.

### Treat absence as a named-resource error

The deletion operation will distinguish an absent local name from storage errors, allowing the command to return a non-zero message that identifies every requested meme that was not found. The command will validate every supplied name before interacting with the inventory, matching creation's local-name handling.

Alternative considered: treat an absent record as success. This is simpler for automation but conflicts with the requested kubectl-style behavior and conceals misspelled resource names.

### Report post-replacement durability uncertainty accurately

The write path will distinguish failures before the atomic replacement from failures after it. A failure before replacement leaves the prior document intact. If replacement succeeds but directory synchronization fails, the operation will report a distinct actionable outcome that the deletion may have succeeded but durable persistence was not confirmed; it will not report the deletion as definitely failed or claim that the prior inventory was preserved.

Alternative considered: return the same generic persistence error for all write failures. This would incorrectly describe the state after a successful replacement and can lead callers to retry an operation that already took effect.

### Keep the scope strictly local

The delete command requires only the selected inventory store and does not receive an Imgflip client or credentials. Its help and README documentation will state that the hosted image remains unaffected.

Alternative considered: remotely delete the Imgflip artifact. Existing records do not retain a remote deletion handle, and adding remote behavior would expand the feature beyond local resource management.

## Risks / Trade-offs

- [A user may expect the public Imgflip image to disappear] -> State the local-only boundary in command help and README documentation.
- [Concurrent commands can interleave between requested names] -> Acquire the shared exclusive mutation lock for each single-name operation, preserving atomicity per operation while intentionally allowing API-like sequencing.
- [A request can name an unknown or previously deleted record] -> Validate every requested name, invoke one operation per name, report each absent request, and return a non-zero result.
- [A failure after replacement leaves durability uncertain] -> Return an outcome that distinguishes unconfirmed durability from failures that leave the prior document intact.

## Migration Plan

No migration is required: the versioned inventory document format remains unchanged. The new command works with existing inventory files, including an empty document. Rollback consists of removing the command; no persisted-data conversion is needed.
