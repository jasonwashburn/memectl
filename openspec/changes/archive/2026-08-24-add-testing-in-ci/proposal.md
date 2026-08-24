## Why

The project has unit tests and a documented `mise run test` command, but pull requests and pushes to `main` do not execute those tests in GitHub Actions. Regressions can therefore merge despite passing the existing hook, lint, and workflow-security checks.

## What Changes

- Add a dedicated GitHub Actions test workflow that runs for pull requests and pushes to `main`.
- Set up the repository's pinned toolchain with the existing mise action.
- Run only the existing `mise run test` task so CI and contributor test commands remain aligned.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None.

This is CI tooling only; it does not change memectl product requirements.

## Impact

- Adds a workflow under `.github/workflows/`.
- Uses the existing pinned `actions/checkout` and `jdx/mise-action` actions.
- Adds a required CI status that executes the existing Go test suite; no application code, dependencies, or user-facing CLI behavior changes.
