## Why

The README currently teaches people to run memectl from source, even though supported release archives are now published. This mixes end-user guidance with contributor workflow and makes the project less approachable for people who simply want to use the CLI.

## What Changes

- Reorient the README around installing a published release archive, adding the executable to `PATH`, and using the `memectl` command directly.
- Preserve a short statement that memectl is also an experiment in OpenSpec and spec-driven development.
- Add a contributing guide that moves source setup, hooks, checks, source execution, and release-snapshot instructions out of the README.
- Establish a concise, meme-aware documentation voice without forced jokes or Kubernetes-themed language.
- Format regular-text mentions of `memectl` as inline code in `README.md` and `CONTRIBUTING.md`.

## Capabilities

### New Capabilities

None. This change updates repository documentation only and does not introduce runtime behavior.

### Modified Capabilities

None. This change does not modify CLI requirements.

## Impact

- Affected documentation: `README.md` and a new `CONTRIBUTING.md`.
- No application code, CLI API, release artifacts, dependencies, or supported-platform behavior changes.
