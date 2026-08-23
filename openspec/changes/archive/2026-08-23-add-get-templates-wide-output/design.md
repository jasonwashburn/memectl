## Context

`memectl get templates` currently writes one tab-separated table after the Imgflip client returns complete template records. The public response contains a direct image URL, but the local template model does not retain it. See `proposal.md` for motivation and the template-listing delta spec for behavior.

## Goals / Non-Goals

**Goals:**
- Preserve byte-for-byte default table output when no output option is provided.
- Expose the direct Imgflip image URL only in the opt-in wide view.
- Establish a stable output-selector flag surface for future formats.

**Non-Goals:**
- Supporting JSON or any format besides `wide` in this change.
- Changing template retrieval, filtering, ordering, or default columns.
- Introducing a generic output-rendering framework before more than two formats require one.

## Decisions

### Use one string output selector with long and short flag forms

The templates command will define a single Cobra output option named `output` with the `o` shorthand. An absent value selects the existing table; `wide` selects the additional URL column. Any other value is rejected before template retrieval or rendering, preventing confusing network work and ensuring errors produce no rows.

Alternative considered: add a `--wide` boolean. Rejected because it would require a second, incompatible flag when a future JSON format is added.

Alternative considered: accept future format names now. Rejected because accepting values without their promised representation would make the CLI contract misleading.

### Extend the template model with Imgflip's URL field

The Imgflip response model will retain the endpoint's `url` field as part of a template record. The command will use that field only when rendering wide output, leaving the client retrieval flow unchanged.

Alternative considered: derive the image URL from the template ID. Rejected because it duplicates provider URL conventions and can diverge from the canonical URL returned by Imgflip.

### Render default and wide tables from the same completed result set

The command will continue to retrieve and validate the full template list before producing output. It will select the header and row shape from the requested format, appending `URL` only for `wide`; the default header and rows stay unchanged.

Alternative considered: create separate command paths for default and wide output. Rejected because retrieval and error behavior would be duplicated for a single additional column.

## Risks / Trade-offs

- [Long image URLs make wide rows exceed terminal width] -> Wide output is explicitly opt-in; default output remains compact.
- [Imgflip omits an image URL] -> Render the provider's returned value without inventing a replacement URL; client response validation remains limited to the fields required by each selected view.
- [Future formats require more rendering structure] -> Keep the selector and format branching localized so a later change can introduce a renderer only when justified.

## Migration Plan

This is an additive flag and output mode. Existing invocations with no output option retain their current output. Rollback consists of removing the optional selector and URL field in a later release; no persisted data or migration is involved.
