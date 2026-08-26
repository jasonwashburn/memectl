## 1. Describe Command

- [ ] 1.1 Add `describe meme <name>` and its `desc` alias to the root CLI, reading the selected local inventory only; verify command help exposes both invocation forms.
- [ ] 1.2 Implement single-name validation, local record lookup, not-found handling, and labeled detail rendering with indexed caption texts and the stored UTC timestamp; verify focused command tests cover successful, empty-caption, invalid-argument, missing-record, and inventory-read-error cases without output on failure.

## 2. Documentation And Verification

- [ ] 2.1 Document local-only `memectl describe meme <name>` usage in the README; verify the documented example and scope statement are present.
- [ ] 2.2 Run `mise run check` and verify formatting, linting, vetting, tests, coverage, and build complete successfully.
