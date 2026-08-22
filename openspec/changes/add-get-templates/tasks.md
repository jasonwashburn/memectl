## 1. Imgflip Client

- [ ] 1.1 Define the Imgflip template response and template data models, then add tests that decode a successful `get_memes` response.
- [ ] 1.2 Implement a testable Imgflip client method that requests public templates and verify tests cover transport failures, unsuccessful API responses, and malformed response bodies.

## 2. Template Listing Command

- [ ] 2.1 Add the `get` command group and `templates` subcommand to the Cobra command tree, then verify `memectl get templates --help` exposes the command.
- [ ] 2.2 Connect `memectl get templates` to the Imgflip client and render each complete result as an ID, name, box-count, and `WIDTHxHEIGHT` table row; verify command tests cover successful output and no output on client failure.

## 3. Verification

- [ ] 3.1 Run `mise run fmt`, `mise run vet`, `mise run test`, and `mise run build` to verify formatting, static analysis, unit tests, and the production build.
