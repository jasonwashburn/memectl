## 1. Imgflip Client

- [x] 1.1 Define the Imgflip template response and template data models, then add tests that decode a successful `get_memes` response.
- [x] 1.2 Implement a testable Imgflip client method that requests public templates and verify tests cover transport failures, unsuccessful API responses, and malformed response bodies.
- [ ] 1.3 Give the default Imgflip client a finite HTTP timeout and add coverage for that default behavior.
- [ ] 1.4 Correct client-test URL diagnostics to format URL strings consistently with the comparison.

## 2. Template Listing Command

- [x] 2.1 Add the `get` command group and `templates` subcommand to the Cobra command tree, then verify `memectl get templates --help` exposes the command.
- [x] 2.2 Connect `memectl get templates` to the Imgflip client and render each complete result as an ID, name, box-count, and `WIDTHxHEIGHT` table row; verify command tests cover successful output and no output on client failure.
- [ ] 2.3 Set explicit arguments in Cobra execution tests so they do not parse Go test process flags.

## 3. Documentation

- [x] 3.1 Update the README to describe `memectl get templates` and show how to run it, replacing the bootstrap-only project status.

## 4. Verification

- [x] 4.1 Run `mise run fmt`, `mise run vet`, `mise run test`, and `mise run build` to verify formatting, static analysis, unit tests, and the production build.
- [ ] 4.2 Re-run `mise run fmt`, `mise run vet`, `mise run test`, and `mise run build` after the review remediation tasks.
