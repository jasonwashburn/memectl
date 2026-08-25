## 1. Local Inventory Removal

- [x] 1.1 Add shared locking for atomic inventory mutations and a removal operation that processes requested names in order, removes every present record once, preserves a valid empty versioned document, and reports absent requests separately from storage failures; verify with inventory unit tests for single and multiple removal, final-record removal, mixed present-and-absent batches, duplicate names, wholly absent batches, corrupt state, and write failures.
- [x] 1.2 Extend the command store contract for record removal and update affected test doubles; verify all command packages compile and existing create/list tests continue to pass.

## 2. Delete Command

- [x] 2.1 Add and register `memectl delete meme <name> [<name>...]` with at-least-one-name and per-name local-name validation; verify command tests cover missing and invalid names without mutating the inventory.
- [x] 2.2 Connect the command to the inventory removal operation, report every successful deletion, and return a non-zero named-resource not-found error for every absent request; verify command tests cover deletion of one and multiple records, preservation of other records, deletion of the final records, mixed present-and-absent batches, and duplicate names.
- [x] 2.3 Ensure the delete command uses no Imgflip client or credentials and documents local-only scope in command help; verify command tests demonstrate deletion succeeds without credentials or remote calls.

## 3. Documentation And Verification

- [x] 3.1 Document `memectl delete meme <name> [<name>...]` in the README, including that it removes only local inventory metadata and leaves the hosted Imgflip image unchanged; verify the documented command and scope match the CLI contract.
- [x] 3.2 Run the Go test suite and formatting/lint checks configured by the repository; verify all checks pass.
