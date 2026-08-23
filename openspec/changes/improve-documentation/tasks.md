## 1. User Documentation

- [x] 1.1 Rewrite `README.md` around release-archive installation, including supported macOS/Linux amd64/arm64 targets, extraction, `PATH` placement, and `memectl --version`; verify every installation claim matches `.goreleaser.yaml` and the published Releases page.
- [x] 1.2 Replace source-based usage examples with direct `memectl get templates` and `memectl create meme` commands, retaining accurate Imgflip credential guidance; verify the commands and flags match the Cobra command help.
- [x] 1.3 Keep a concise OpenSpec/spec-driven-development experiment statement and link to `CONTRIBUTING.md`; verify the README addresses an end user without requiring mise, Go, or `go run`.

## 2. Contributor Documentation

- [x] 2.1 Add `CONTRIBUTING.md` with source prerequisites, `mise install`, hook installation, and the project quality-check commands; verify every documented task exists in `mise.toml`.
- [x] 2.2 Document source execution with `go run`, local release snapshots, and the project’s OpenSpec experiment context without prescribing its workflow; verify contributor-only instructions no longer need to appear in the README.

## 3. Documentation Review

- [ ] 3.1 Review both documents for concise, technically literal instructions with only sparse, meme-native personality; verify no Kubernetes-themed jokes, unsupported installation paths, or copy-paste ambiguity remain, and format every regular-text mention of `memectl` as inline code.
- [x] 3.2 Run `mise run check` and verify the repository validation suite passes after the documentation changes.
