## 1. Test Workflow

- [x] 1.1 Add a dedicated GitHub Actions workflow that runs on pull requests and pushes to `main`, checks out without persisted credentials, installs the pinned mise toolchain, and executes only `mise run test`; verify `zizmor .github/workflows` reports no findings.

## 2. Validation

- [x] 2.1 Run `mise run test` locally and verify the complete Go test suite passes.
