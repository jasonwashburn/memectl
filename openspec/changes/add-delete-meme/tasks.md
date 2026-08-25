## 1. Local Inventory Removal

- [ ] 1.1 Add a locked, atomic inventory operation that removes every present record among one or more named records, preserves a valid empty versioned document, and reports absent names separately from storage failures; verify with inventory unit tests for single and multiple removal, final-record removal, mixed present-and-absent batches, wholly absent batches, corrupt state, and write failures.
- [ ] 1.2 Extend the command store contract for record removal and update affected test doubles; verify all command packages compile and existing create/list tests continue to pass.

## 2. Delete Command

- [ ] 2.1 Add and register `memectl delete meme <name> [<name>...]` with at-least-one-name and per-name local-name validation; verify command tests cover missing and invalid names without mutating the inventory.
- [ ] 2.2 Connect the command to the inventory removal operation, report every successful deletion, and return a non-zero named-resource not-found error for every absent name; verify command tests cover deletion of one and multiple records, preservation of other records, deletion of the final records, and mixed present-and-absent batches.
- [ ] 2.3 Ensure the delete command uses no Imgflip client or credentials and documents local-only scope in command help; verify command tests demonstrate deletion succeeds without credentials or remote calls.

## 3. Documentation And Verification

- [ ] 3.1 Document `memectl delete meme <name> [<name>...]` in the README, including that it removes only local inventory metadata and leaves the hosted Imgflip image unchanged; verify the documented command and scope match the CLI contract.
- [ ] 3.2 Run the Go test suite and formatting/lint checks configured by the repository; verify all checks pass.
