## Purpose

Publish ready-to-run memectl binaries for supported operating systems with verifiable release provenance.

## ADDED Requirements

### Requirement: Supported-platform release archives
The project SHALL publish a compressed memectl binary archive for each release on macOS and Linux for both `amd64` and `arm64` architectures. The project SHALL NOT publish Windows release archives.

#### Scenario: A release is published
- **WHEN** a new memectl GitHub release is created through the release automation
- **THEN** its assets include one archive for each of `darwin/amd64`, `darwin/arm64`, `linux/amd64`, and `linux/arm64`

### Requirement: Release artifact integrity
The project SHALL publish a checksum manifest covering every binary archive generated for a release.

#### Scenario: A user verifies a downloaded archive
- **WHEN** a user downloads a release archive and its checksum manifest
- **THEN** the manifest contains a checksum entry for that archive

### Requirement: Embedded release provenance
Release binaries SHALL report the release version, source commit, and build date through the existing `memectl --version` output. Development builds SHALL retain their existing fallback metadata when no release build metadata is supplied.

#### Scenario: A user inspects a release binary
- **WHEN** a user runs `memectl --version` from a published archive
- **THEN** the output identifies the release version, source commit, and build date used to create that binary

#### Scenario: A developer runs a local build
- **WHEN** a developer builds memectl without release linker metadata
- **THEN** `memectl --version` reports the existing development fallback values

### Requirement: Local release-build validation
The project SHALL provide a `mise` task named `release-snapshot` that runs the GoReleaser release build pipeline locally in snapshot mode. The task SHALL produce the same four supported-platform archives and checksum manifest as a release build, and it SHALL NOT create or modify GitHub releases or upload GitHub assets.

#### Scenario: A developer runs the local release build
- **WHEN** a developer runs `mise run release-snapshot`
- **THEN** the local release output contains archives for `darwin/amd64`, `darwin/arm64`, `linux/amd64`, and `linux/arm64` plus a checksum manifest

#### Scenario: A developer validates a release build locally
- **WHEN** `mise run release-snapshot` completes
- **THEN** no GitHub release is created or modified and no GitHub assets are uploaded

### Requirement: Conditional release publishing
The main-branch release automation SHALL build and upload binary assets only when Release Please reports that it created a release. The release assets SHALL be added to the GitHub release created for the corresponding version.

#### Scenario: Release Please creates a release
- **WHEN** Release Please creates a GitHub release after a push to the main branch
- **THEN** the release automation builds the supported-platform archives and uploads them with the checksum manifest to that release

#### Scenario: Release Please does not create a release
- **WHEN** Release Please completes after a push to the main branch without creating a GitHub release
- **THEN** the release automation does not run a release artifact build or upload artifacts
