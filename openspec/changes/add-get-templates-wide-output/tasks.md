## 1. Template Data

- [x] 1.1 Extend the Imgflip template response model to retain each template's direct image URL, and verify client decoding tests cover the `url` field.

## 2. Template Output

- [x] 2.1 Add the `--output` / `-o` selector to `memectl get templates`, accepting only the default empty value and `wide`, and verify command tests cover both aliases plus unsupported-format errors with no rows.
- [x] 2.2 Render the `URL` header and image URL value only for wide output while preserving existing default table output exactly, and verify command tests cover wide rows and unchanged default rows.

## 3. Documentation And Verification

- [x] 3.1 Document `memectl get templates --output wide` in the README, and verify the example and flag spelling match command help.
- [x] 3.2 Run `mise run fmt`, `mise run vet`, `mise run test`, and `mise run build` to verify formatting, static analysis, unit tests, and the production build.
