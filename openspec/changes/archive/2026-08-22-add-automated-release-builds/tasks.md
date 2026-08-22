## 1. Release Build Configuration

- [x] 1.1 Add a GoReleaser configuration that builds `memectl` for darwin and linux on `amd64` and `arm64`, packages the four outputs as archives, and generates a checksum manifest; verify `goreleaser check` accepts the configuration.
- [x] 1.2 Configure GoReleaser linker flags to inject its version, commit, and date into the existing `cmd` package metadata; verify a snapshot build reports populated values from `memectl --version` while the normal development build retains fallback metadata.
- [x] 1.3 Add the pinned GoReleaser version and a `release-snapshot` task to `mise.toml`, then ignore its local output directory; verify `mise run release-snapshot` builds the four archives and checksum manifest without publishing to GitHub.

## 2. Release Automation

- [x] 2.1 Rename the Release Please workflow as the release workflow, then add a named Release Please step and a commit-pinned GoReleaser action that runs only when `release_created` is true; verify the workflow YAML and action pinning checks pass.
- [x] 2.2 Configure the conditional GoReleaser invocation to reuse the Release Please release and upload the supported archives and checksum manifest with repository write permissions; verify the workflow supplies the release token only to the publishing invocation.

## 3. Verification

- [x] 3.1 Run formatting, linting, static analysis, tests, the normal development build, GoReleaser configuration validation, and `mise run release-snapshot`; verify all commands complete successfully and the snapshot output contains exactly four platform archives plus the checksum manifest.
