## Why

The project currently implements only public template listing. Users need a credentialed, command-line path to caption a selected Imgflip image template without relying on a browser workflow.

## What Changes

- Add `memectl create meme <template-id>` for creating a captioned static image meme through Imgflip.
- Accept one or more ordered, repeatable `--text` flags so the command supports templates with any number of text boxes instead of assuming a top/bottom layout.
- Read Imgflip account credentials from `IMGFLIP_USERNAME` and `IMGFLIP_PASSWORD` environment variables.
- Print a concise creation summary plus the direct URLs for the hosted generated image and its Imgflip page.
- Document the credential setup and captioned-meme creation workflow in the README.
- Reserve `create` as the command group for future resource-specific creation modes such as GIFs, automatic memes, and AI memes; those modes are not included in this change.

## Capabilities

### New Capabilities
- `captioned-meme-creation`: Create a hosted captioned static meme from an Imgflip template using ordered text inputs and environment-provided credentials.

### Modified Capabilities
- None.

## Impact

- Affected code: Cobra command tree under `cmd` and the Imgflip HTTP client under `internal/imgflip`.
- External API: authenticated, form-encoded `POST /caption_image` requests to Imgflip.
- Documentation: README gains credential and meme-creation instructions.
- No new third-party dependencies, local credential storage, or persisted user data.
