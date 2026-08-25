## 1. Local Inventory Removal

- [x] 1.1 Add shared locking for atomic inventory mutations and a single-record removal operation that preserves a valid empty versioned document, distinguishes absence, pre-replacement failures, and post-replacement durability uncertainty; verify with inventory unit tests for single and final-record removal, corrupt state, pre-replacement write failures, and post-replacement directory-sync failures.
- [x] 1.2 Change the command store contract to a single-record removal operation and update affected test doubles; verify all command packages compile and existing create/list tests continue to pass.

## 2. Delete Command

- [x] 2.1 Add and register `memectl delete meme <name> [<name>...]` with at-least-one-name and per-name local-name validation; verify command tests cover missing and invalid names without mutating the inventory.
- [x] 2.2 Invoke the inventory removal operation once for each supplied name in order, report every successful deletion, and return a non-zero named-resource not-found error for every absent request; verify command tests cover deletion of one and multiple records, preservation of other records, deletion of the final records, mixed present-and-absent batches, duplicate names, and a later storage failure after earlier success.
- [x] 2.3 Ensure the delete command uses no Imgflip client or credentials and documents local-only scope in command help; verify command tests demonstrate deletion succeeds without credentials or remote calls.

## 3. Documentation And Verification

- [x] 3.1 Document `memectl delete meme <name> [<name>...]` in the README, including that it removes only local inventory metadata and leaves the hosted Imgflip image unchanged; verify the documented command and scope match the CLI contract.
- [x] 3.2 Run the Go test suite and formatting/lint checks configured by the repository; verify all checks pass.
