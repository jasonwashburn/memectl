## Context

The command tree currently contains only the Cobra root command, and `internal/imgflip` contains no client implementation. This change introduces both a nested `get templates` command and the first boundary around Imgflip's HTTP API. See `proposal.md` for motivation and `specs/template-listing/spec.md` for behavior.

## Goals / Non-Goals

**Goals:**
- Keep command orchestration separate from HTTP request and response handling.
- Establish a small, testable Imgflip client that later commands can extend.
- Render a stable, human-readable table from complete API results.
- Document the command in the README without introducing further documentation surfaces.

**Non-Goals:**
- Authentication, user configuration, caching, pagination, search, filtering, or machine-readable output.
- Creating memes or retrieving a single template by ID.
- Generalizing an API abstraction beyond the endpoints actually needed.

## Decisions

### Use the public `get_memes` endpoint without credentials

The command will use Imgflip's public template-list endpoint and model its success flag plus meme records. This keeps the first command usable immediately and avoids introducing configuration before a protected operation needs it.

Alternative considered: introduce credentials and shared configuration now. Rejected because this endpoint does not require them and the added setup would obstruct the first-use experience.

### Place HTTP behavior in `internal/imgflip`

An Imgflip client will own the endpoint URL, request execution, response decoding, and conversion of unsuccessful or malformed responses into errors. Cobra commands will request templates from that client and render the result.

Alternative considered: perform HTTP requests directly in the command. Rejected because it would mix transport failures, API decoding, and terminal presentation, making future API calls harder to test and reuse.

### Add `get` as a command group with `templates` beneath it

The command tree will follow the project's `memectl <verb> <resource>` convention: `memectl get templates`. The `get` group creates a natural extension point for future readable resources while keeping template listing explicitly plural.

Alternative considered: add `memectl templates`. Rejected because it departs from the documented kubectl-inspired interface.

### Render only after successful retrieval and validation

The command will construct terminal output from the returned slice only after the client accepts the response. Each row will include ID, name, box count, and `WIDTHxHEIGHT` dimensions.

Alternative considered: stream rows while decoding. Rejected because an invalid or unsuccessful response must not result in partial output.

### Document the command in the README only

The README will replace its bootstrap-only description with the available template-listing workflow and an invocation example. Additional command reference documentation is deferred until the CLI has more than one user-facing workflow.

## Risks / Trade-offs

- [Imgflip changes its response fields or endpoint behavior] -> Keep the API response model focused on required fields and cover success, API-failure, and invalid-response paths with client tests.
- [The public response is large or changes ordering] -> This first command returns the provider's complete ordering; add filters or output modes only after users need them.
- [Terminal column widths vary with long template names] -> Use Go's tabular formatting so columns remain readable without inventing truncation rules prematurely.

## Migration Plan

This is an additive command with no existing user data or command compatibility surface. Release it with the normal binary build; rollback consists of removing the new command in a later release if the upstream API becomes unsuitable.
