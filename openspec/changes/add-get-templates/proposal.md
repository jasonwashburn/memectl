## Why

memectl currently has no feature commands, so users cannot discover the Imgflip templates needed to generate a meme. Listing the available templates creates the first practical CLI workflow and establishes template IDs as input to future commands.

## What Changes

- Add a `memectl get templates` command that retrieves the public template list from Imgflip.
- Display each template's identifier, name, text-box count, and dimensions in a readable terminal table.
- Report retrieval and response failures as command errors without presenting partial or stale template data.

## Capabilities

### New Capabilities
- `template-listing`: Retrieves and displays Imgflip meme templates through the `memectl get templates` command.

### Modified Capabilities

- None.

## Impact

- Adds Cobra commands under `cmd`.
- Adds the first HTTP integration and response model in `internal/imgflip`.
- Updates the README to document the first user-facing command.
- Uses Imgflip's unauthenticated `get_memes` HTTP endpoint; no configuration, credentials, or new third-party Go dependencies are required.
