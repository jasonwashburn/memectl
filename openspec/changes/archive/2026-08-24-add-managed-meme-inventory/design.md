## Context

The current command layer invokes Imgflip directly and writes the generated URLs to standard output. There is no durable local state or local resource abstraction. See proposal.md for motivation and the delta specs for externally observable behavior.

## Goals / Non-Goals

**Goals:**
- Add one durable, local inventory of named memes without adding remote reads.
- Keep Imgflip HTTP behavior isolated from local storage concerns.
- Make inventory access deterministic and independently testable.
- Preserve state safely across interrupted writes and report the unavoidable remote-success/local-write-failure outcome clearly.

**Non-Goals:**
- Implement individual `get meme`, `describe meme`, deletion, overwrite, import, automatic names, declarative files, remote availability checks, or credential persistence.
- Add a general configuration or context system.
- Reconcile remote memes that succeeded but could not be stored locally.

## Decisions

### Use a versioned single-file inventory

Store a versioned JSON document containing records with the local name, template ID, ordered text values, image URL, page URL, and UTC creation timestamp. A single inventory supports atomic whole-document replacement, deterministic listing, and a future `get` or `describe` command without remote reads.

One-file-per-meme storage was considered, but makes cross-record consistency, list ordering, and future migration more complex. Credentials and template metadata are intentionally excluded because they are not needed to describe the created resource and must not expand scope.

### Resolve one explicit store file

The default store is `~/.meme/memes.json`. `MEME_STORE` overrides it as an exact file path, rather than a directory root, to provide direct test isolation and avoid introducing a broad configuration-home contract. The parent directory is only created while persisting a successful creation.

For development activity launched through mise, configure `MEME_STORE` directly to `<project-root>/.meme/memes.json` so development commands do not touch persistent user state.

### Separate inventory access from command and Imgflip concerns

Introduce a local inventory abstraction with operations to load/list and persist a new record. Commands depend on it alongside the existing Imgflip interface; the Imgflip client remains responsible only for HTTP caption creation. Construction supplies the default filesystem implementation and tests supply fakes or temporary stores.

### Validate and check before remote creation

Validate the local name and load the selected inventory before credentials and remote creation proceed. Reject an existing name before calling Imgflip. Acquire an advisory lock around the reload, duplicate check, and atomic write to prevent concurrent callers from overwriting local inventory updates.

The lock is not held during the remote request and cannot turn Imgflip plus a filesystem write into one transaction. Concurrent invocations can still create an unrecorded duplicate remote meme when the later caller reaches persistence, which is an accepted limitation for this increment.

### Write atomically and preserve bad state

Serialize the complete next document to a restrictive temporary file in the target directory, close and synchronize it, then rename it over the store file. Synchronize the directory where the platform supports it. Missing state is empty; unreadable, malformed, invalid, or unknown-version state is an error and is never silently reset or overwritten.

### Make partial success explicit

Remote creation must precede local persistence because the record includes Imgflip URLs. If persistence fails after a successful response, return a non-zero error, print both URLs for recovery, and do not print the normal creation-success output. No automatic retry or recovery is possible without out-of-scope remote inspection or import behavior.

### Match kubectl-style list presentation

`get memes` uses a table sorted by name. Its normal form includes `NAME`, `TEMPLATE ID`, `AGE`, and `IMAGE URL`; `wide` adds `PAGE URL`. Persist exact UTC timestamps but render human-readable ages in list output. A missing or empty inventory prints `No resources found.` rather than a header-only table.

## Risks / Trade-offs

- [Remote meme succeeds but inventory write fails] → Return a non-zero partial-success error with both URLs; preserve no false local success claim.
- [Store file is manually edited or corrupted] → Validate fully on read, preserve the file, and identify its path in the error.
- [Concurrent creators use the same name] → The local advisory lock preserves inventory updates and rejects the later local record; document that no cross-process transaction with Imgflip exists.
- [Atomic-write behavior differs by filesystem] → Use a same-directory temporary file and rename, the strongest portable filesystem primitive available.

## Migration Plan

This is pre-1.0 and no local inventory exists, so no persisted-data migration is required. The positional-template command syntax is removed directly. Users invoke the new named syntax and receive an empty local inventory until their first successful managed creation.
