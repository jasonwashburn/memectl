## 1. Local Inventory Storage

- [ ] 1.1 Define the versioned managed-meme record and inventory document, including name, template ID, ordered texts, URLs, and UTC creation time; verify serialization round-trip tests preserve every field and exclude credentials.
- [ ] 1.2 Implement store resolution for the default `~/.meme/memes.json` path and exact `MEME_STORE` override; verify tests cover default resolution and isolated override paths.
- [ ] 1.3 Implement inventory loading and validation for missing, malformed, invalid-record, unreadable, and unsupported-version state; verify tests confirm missing state is empty and invalid state is preserved with an actionable path-bearing error.
- [ ] 1.4 Implement atomic inventory updates with a same-directory temporary file and replacement rename; verify tests cover successful append, duplicate-safe recheck, and no partial inventory after a simulated write failure.

## 2. Managed Meme Creation

- [ ] 2.1 Change `create meme` to accept exactly one local name and required `--template`, removing positional template-ID support; verify command tests cover valid syntax and all invalid argument combinations.
- [ ] 2.2 Validate DNS-label-like local names and reject duplicates before the Imgflip client is called; verify command tests assert no remote call or store mutation for invalid and duplicate names.
- [ ] 2.3 Persist a complete managed record after Imgflip succeeds and report the named resource only after persistence; verify command tests assert stored fields, ordered caption text, and success output.
- [ ] 2.4 Handle Imgflip-success/local-store-failure as a non-zero partial success that reports both URLs without a normal success summary; verify command tests cover this recovery output and preserve existing state.

## 3. Local Meme Listing

- [ ] 3.1 Add `memectl get memes` backed only by the local inventory; verify tests confirm no Imgflip dependency is invoked and records are sorted by name.
- [ ] 3.2 Render default list columns `NAME`, `TEMPLATE ID`, `AGE`, and `IMAGE URL`, with `PAGE URL` added by `--output wide` and `-o wide`; verify exact table output and unsupported-format errors in command tests.
- [ ] 3.3 Render `No resources found.` for missing or empty inventory; verify tests cover both empty-state paths and corrupt-state failure behavior.

## 4. Documentation and Verification

- [ ] 4.1 Add `.meme/` to `.gitignore` and configure a fallback `MEME_STORE` in `mise.toml` resolving to the project-local `.meme/memes.json`; verify mise uses the local store when unset and preserves an explicitly supplied store path.
- [ ] 4.2 Update README examples and descriptions for named creation, `--template`, managed local listing, default store location, and `MEME_STORE`; verify documented commands match the command help.
- [ ] 4.3 Run `go test ./...` and the project lint/check commands from CONTRIBUTING; verify all checks pass.
