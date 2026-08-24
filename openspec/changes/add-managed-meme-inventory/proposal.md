## Why

`memectl` currently creates memes remotely but forgets them immediately, so users cannot treat created memes as named local resources. A local inventory makes creation identifiable, listable, and ready for future resource-oriented `get` and `describe` commands.

## What Changes

- **BREAKING** Change creation syntax to `memectl create meme <name> --template <template-id> --text <text>...`; remove support for the positional template identifier.
- Validate local meme names and reject duplicate names before contacting Imgflip.
- Persist successful creations in a locally managed, versioned meme inventory that retains creation inputs and Imgflip-returned URLs, but never credentials or unrequested template attributes.
- Add `memectl get memes` with kubectl-style default and wide table output. The wide view adds `PAGE URL`.
- Store the inventory at `~/.meme/memes.json` by default and allow `MEME_STORE` to select a specific inventory file for isolated testing or usage.
- Safely handle absent, unreadable, corrupt, and unsupported inventory state; use atomic persistence and report remote-success/local-persistence-failure as a partial success.
- Update user documentation for the managed creation and listing workflows.

## Capabilities

### New Capabilities
- `meme-inventory`: Persist and list locally managed meme resources.

### Modified Capabilities
- `captioned-meme-creation`: Create captioned memes as named local resources using the required `--template` flag and persist them after Imgflip succeeds.

## Impact

- Affects Cobra commands and their tests in `cmd/`.
- Adds local inventory storage and tests, separate from the Imgflip HTTP client.
- Updates `README.md` usage and setup documentation.
- Introduces no new remote API calls, credential persistence, or external dependencies.
