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

### Model deletion as an inventory mutation

Add a record-removal operation to the inventory store and expose it through the command store interface. The operation will validate and load the inventory while holding the existing lock, identify every requested local name that is present, and atomically write the retained records using the existing write path. It will report the names that were absent so the command can emit one successful deletion line per removed record and a not-found error for each missing name.

Using a store operation rather than filtering a command-level `Load` result avoids a load-modify-write race with concurrent creation or deletion. Processing all requested names in one locked mutation mirrors kubectl's non-transactional named-resource behavior: present records are removed even when other names are absent. Reusing the existing write path retains the durable replacement and empty-document representation.

Alternative considered: filter records in the command and introduce a generic save operation. This would leak inventory persistence details into the command layer and make lock ownership less clear.

### Treat absence as a named-resource error

The deletion operation will distinguish an absent local name from storage errors, allowing the command to return a non-zero message that identifies the requested meme as not found. The command will validate its single name before interacting with the inventory, matching creation's local-name handling.

Alternative considered: treat an absent record as success. This is simpler for automation but conflicts with the requested kubectl-style behavior and conceals misspelled resource names.

### Keep the scope strictly local

The delete command requires only the selected inventory store and does not receive an Imgflip client or credentials. Its help and README documentation will state that the hosted image remains unaffected.

Alternative considered: remotely delete the Imgflip artifact. Existing records do not retain a remote deletion handle, and adding remote behavior would expand the feature beyond local resource management.

## Risks / Trade-offs

- [A user may expect the public Imgflip image to disappear] -> State the local-only boundary in command help and README documentation.
- [Concurrent commands could otherwise overwrite one another] -> Keep lookup and rewrite under the inventory's existing exclusive lock.
- [A batch can include an unknown name] -> Validate every requested name, remove records that are present, report each absent name, and return a non-zero result.
- [A write failure after lookup could leave intent unclear] -> Return an error and rely on atomic replacement so the prior inventory remains intact.

## Migration Plan

No migration is required: the versioned inventory document format remains unchanged. The new command works with existing inventory files, including an empty document. Rollback consists of removing the command; no persisted-data conversion is needed.
