## Why

The default template table is intentionally compact, but users who need to inspect a template visually must separately locate its image. An opt-in wide view exposes that direct image URL while preserving the existing output for current users and scripts.

## What Changes

- Add an output-format option to `memectl get templates`, available as `-o` and `--output`.
- Support the `wide` output format, which appends each template's direct Imgflip image URL to the existing table.
- Preserve the current default table columns and formatting when no output format is requested.
- Validate unsupported output-format values as command errors without displaying template rows.

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `template-listing`: Add an opt-in wide table that includes template image URLs while retaining the existing default output.

## Impact

- Updates the `get templates` Cobra command and its output-format validation.
- Extends the Imgflip template response model to retain the image URL returned by the public template endpoint.
- Adds command and client coverage for default and wide output behavior.
- Updates user-facing command documentation to describe wide output.
