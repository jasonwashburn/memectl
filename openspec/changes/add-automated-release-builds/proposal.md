## Why

memectl releases currently contain source and release notes but no installable binaries. Automating platform-specific release builds makes tagged releases immediately usable without requiring users to install Go and build from source.

## What Changes

- Add GoReleaser configuration that produces memectl archives and checksums for supported macOS and Linux architectures.
- Inject the release version, source commit, and build date into the existing Cobra version metadata at build time.
- Extend the main-branch release workflow to run GoReleaser only when Release Please creates a GitHub release, uploading the resulting artifacts to that release.
- Make the pinned GoReleaser tool available through mise and add a local snapshot-release task that builds release artifacts without publishing to GitHub.

## Capabilities

### New Capabilities
- `release-distribution`: Publishes verified, versioned memectl binary archives for supported platforms with build provenance metadata.

### Modified Capabilities

- None.

## Impact

- Adds a `.goreleaser.yaml` release configuration and GoReleaser to the managed development toolchain.
- Renames `.github/workflows/release-please.yml` to `.github/workflows/release.yml` and extends it to publish artifacts after a Release Please release is created.
- Uses the existing `cmd` package version variables as Go linker targets; the CLI command interface remains unchanged.
- Produces GitHub release assets for macOS and Linux on `amd64` and `arm64` only.
- Adds an ignored local GoReleaser output directory for snapshot-release artifacts.
