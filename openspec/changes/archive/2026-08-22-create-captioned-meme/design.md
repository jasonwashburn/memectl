## Context

The current command tree has a root command and a `get templates` resource. `internal/imgflip` owns the existing unauthenticated template-list request and provides a test-injectable HTTP client. See proposal.md for the change motivation and `specs/captioned-meme-creation/spec.md` for required behavior.

## Goals / Non-Goals

**Goals:**
- Add a resource-oriented creation command without coupling it to future creation endpoints.
- Extend the Imgflip boundary with authenticated caption-image requests while keeping request formatting and response validation outside Cobra commands.
- Keep credentials out of command arguments, persistent configuration, and terminal output.

**Non-Goals:**
- Support blank template copies, local image downloads, browser launching, or machine-readable output.
- Caption GIFs, search templates, create automatic memes, or create AI memes.
- Add credential persistence, interactive authentication, watermark controls, fonts, coordinates, or other text-box styling controls.

## Decisions

### Add `create meme` as a resource command

The command tree will gain a `create` group with a `meme` resource using `memectl create meme <template-id>`. The template identifier is the required positional name, and each `--text` flag is repeatable and required at least once.

This keeps the established `memectl <verb> <resource> [name] [flags]` grammar and avoids hard-coding a top/bottom model. The `create` group remains able to add distinct resource commands such as `gif`, `auto-meme`, and `ai-meme`, each of which maps to a different Imgflip endpoint and input model.

Alternative considered: model all creation modes as `create meme` with mode flags such as `--ai`. Rejected because the endpoints have incompatible inputs and Premium constraints, which would make the initial command harder to understand and extend.

### Map text flags to Imgflip boxes

The Imgflip client will expose a caption-image operation that accepts a template ID, ordered texts, and credentials. It will send a form-encoded `POST` request to Imgflip's caption-image endpoint with `template_id`, `username`, `password`, and `boxes[<index>][text]` fields.

Using boxes preserves the caller's text casing and supports templates with more than two text positions. The CLI will validate that at least one text is supplied but will not require the supplied count to equal the template's advertised default box count; Imgflip remains the source of truth for template-specific rendering.

Alternative considered: use Imgflip's `text0` and `text1` fields. Rejected because they impose a two-text model and do not preserve the requested multi-text semantics.

### Read credentials at the command boundary

The command will read `IMGFLIP_USERNAME` and `IMGFLIP_PASSWORD` before invoking the client and reject missing or empty values without making an HTTP request. Tests will inject environment lookup or command dependencies rather than requiring process-wide environment mutation where practical.

Alternative considered: credentials as flags. Rejected because command-line arguments can be retained in shell history and observed in process metadata. Alternative considered: a configuration file or system keychain. Rejected because this first credentialed workflow does not need persisted secrets or an account-management surface.

### Validate all provider outcomes in the Imgflip client

The client will retain ownership of HTTP status validation, JSON decoding, Imgflip's success flag, and generated-image URL validation. Unsuccessful Imgflip responses will include the provider's useful error message when available, wrapped in a stable creation-context error. The command will only render output after a validated success result.

This matches template retrieval's transport boundary and prevents partial success output. Error construction must not include the submitted password or serialized request form.

Alternative considered: decode and format responses directly in the command. Rejected because it mixes terminal rendering with provider protocol behavior and makes error-path tests less focused.

### Print both generated URLs with a compact summary

On success, the command will write a short sentence identifying the source template followed by labeled direct image and Imgflip page URLs. Both URLs will be modeled and validated as part of the successful provider result.

Alternative considered: download the image locally. Rejected because hosted output is immediately usable and local file naming, overwrite, and format behavior would expand the command's scope.

## Risks / Trade-offs

- [Generated Imgflip image URLs are publicly accessible to anyone who knows the URL] -> Document this provider behavior through clear command output and avoid implying private storage.
- [Imgflip changes its caption response or rejects a template/text combination] -> Validate response success and image URL; return provider failures without emitting a success result.
- [Environment variables can be absent or exposed by a user's runtime environment] -> Validate presence before requests, never accept password flags, and never include credentials in errors.
- [Some templates have special rendering needs] -> Delegate default placement and rendering to Imgflip; defer custom coordinates and styling.

## Migration Plan

This is an additive command with no existing credential or creation compatibility surface. Release it with the normal binary workflow. If the upstream caption endpoint becomes unsuitable, a later release can remove the command; no local state or generated files require migration.
