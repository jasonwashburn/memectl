## 1. Imgflip Caption Client

- [x] 1.1 Add typed caption-image request and result models plus an authenticated form-encoded Imgflip caption operation that sends ordered `boxes[n][text]` values; verify client unit tests assert the HTTP method, endpoint, form fields, and returned image and page URLs.
- [x] 1.2 Handle transport failures, non-success HTTP responses, unsuccessful Imgflip responses, malformed JSON, and missing generated URLs without exposing submitted credentials; verify focused client tests cover every failure path.

## 2. Captioned Meme Command

- [x] 2.1 Add the `create` command group and `memectl create meme <template-id>` resource command with repeatable, required `--text` input; verify command tests cover help, valid positional arguments, missing text, and invalid argument counts.
- [x] 2.2 Read `IMGFLIP_USERNAME` and `IMGFLIP_PASSWORD` before invoking the caption client and reject either missing or empty value without a request; verify command tests cover both credential failures and successful credential forwarding.
- [x] 2.3 Render a concise success summary naming the source template plus labeled direct hosted image and Imgflip page URLs only after validated creation; verify command tests cover exact success output and no success output for client failures.

## 3. Documentation And Verification

- [x] 3.1 Update the README with required environment variables, a repeatable `--text` creation example, and hosted-image output expectations; verify the documented command matches the CLI contract.
- [x] 3.2 Run `mise run fmt`, `mise run lint`, `mise run vet`, `mise run test`, `mise run coverage`, and `mise run build`; verify all commands complete successfully.
