## Context

The existing Release Please workflow runs after pushes to `main` and creates the version tag, GitHub release, and release notes. It will be renamed as the broader release workflow. The Cobra root command already exposes package variables for version, commit, and date that the Go linker can override. There is no release build configuration, artifact upload process, or locally managed GoReleaser tool. See `proposal.md` and `specs/release-distribution/spec.md` for motivation and release behavior.

## Goals / Non-Goals

**Goals:**
- Produce reproducible release archives for the four supported OS and architecture combinations.
- Attach the archives and a checksum manifest to the exact GitHub release created by Release Please.
- Preserve useful release provenance in `memectl --version` while keeping local development defaults intact.
- Make local GoReleaser configuration validation use the project's pinned toolchain.

**Non-Goals:**
- Windows binaries, package-manager distribution, installers, container images, or binary signing.
- Replacing Release Please's versioning, changelog generation, release notes, or release ownership.
- Changing the CLI version output format or adding runtime update checks.

## Decisions

### Keep Release Please and GoReleaser in one workflow

The renamed main-branch release workflow will retain the Release Please action as a named step with an ID. A subsequent GoReleaser step will execute only when that step reports `release_created` as true. GoReleaser will use the workflow `GITHUB_TOKEN` with `contents: write` permission to upload assets to the release for the tag produced by Release Please.

This avoids a separate workflow that listens for a release event. Releases created with the workflow token do not reliably trigger additional GitHub Actions workflows, whereas the conditional step runs in the same authenticated execution that created the release.

Alternative considered: trigger a standalone workflow from `release.published`. Rejected because it can be suppressed for releases created by `GITHUB_TOKEN` and splits tightly coupled release operations.

### Delegate the supported build matrix and metadata to GoReleaser

The GoReleaser configuration will define one `memectl` build with `darwin` and `linux` targets for `amd64` and `arm64`. Its linker flags will assign GoReleaser's version, commit, and date templates to `github.com/jasonwashburn/memectl/cmd.version`, `.commit`, and `.date`. GoReleaser will archive each target and generate one checksum manifest.

The configuration will not declare Windows targets. GoReleaser's release pipeline will upload its generated archives and checksum manifest to the existing tagged GitHub release.

Alternative considered: add shell build commands to GitHub Actions. Rejected because GoReleaser centralizes cross-compilation, archive naming, provenance injection, checksum generation, and artifact publication in a declarative configuration.

### Pin GoReleaser in both CI and the local toolchain

The workflow will use a commit-pinned GoReleaser action and a fixed GoReleaser major/minor version. The same GoReleaser version will be declared in `mise.toml` so developers can run configuration checks locally. CI will install the repository's pinned Go toolchain before invoking GoReleaser.

Alternative considered: install an unpinned GoReleaser release during CI. Rejected because it permits tool behavior changes without a repository change and makes local validation drift from CI.

### Provide a local snapshot-release task

A `mise run release-snapshot` task will invoke GoReleaser's snapshot release mode with a clean local artifact directory. Snapshot mode runs the configured build, archive, checksum, and metadata-injection pipeline but does not publish a GitHub release or upload assets. The generated output directory will be ignored by Git so developers can inspect artifacts without leaving build output in the worktree.

Alternative considered: maintain a separate local build script. Rejected because it could drift from the GoReleaser configuration used for production releases and fail to exercise the release archive and checksum pipeline.

## Risks / Trade-offs

- [A GoReleaser upgrade changes its configuration schema or archive defaults] -> Pin the tool version and validate the configuration locally and in CI before updating it.
- [The GoReleaser upload step attempts to recreate Release Please's release notes] -> Configure GoReleaser to reuse the existing tagged release and leave Release Please responsible for release content.
- [Cross-compiled output fails at runtime on a supported platform] -> Build all four targets during release and expand platform-specific smoke testing only when CI runners or user demand justify it.
- [A workflow run is not a release] -> Guard the GoReleaser step with the Release Please `release_created` output so ordinary main-branch pushes never publish assets.
- [A local validation command publishes by mistake] -> Use GoReleaser snapshot mode in the mise task and verify that it completes without GitHub release or asset changes.

## Migration Plan

The change is additive. After merging, the next Release Please release will cause GoReleaser to build and attach artifacts automatically. The release workflow can be rolled back by removing the conditional GoReleaser step; existing GitHub releases and their uploaded assets remain unaffected.
