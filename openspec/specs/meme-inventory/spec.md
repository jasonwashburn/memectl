# meme-inventory Specification

## Purpose

Provide durable, named local meme resources that users can list without making an Imgflip request.

## Requirements

### Requirement: Persist managed meme records
The system SHALL persist a successfully created meme as a locally managed record containing its name, template identifier, ordered caption texts, generated image URL, Imgflip page URL, and creation timestamp. It SHALL NOT persist Imgflip credentials or template attributes that were not supplied or returned during creation.

#### Scenario: A created meme is recorded
- **WHEN** Imgflip successfully creates a meme and the local inventory can be persisted
- **THEN** the system SHALL create a managed record with the creation input and generated URLs

### Requirement: Resolve the inventory store
The system SHALL use `~/.meme/memes.json` as the default inventory store. When `MEME_STORE` is non-empty, the system SHALL use its value as the exact inventory store file path instead.

#### Scenario: A store override is supplied
- **WHEN** a user sets `MEME_STORE` to an inventory file path
- **THEN** the system SHALL read and write only that inventory store

### Requirement: Validate local meme names
The system SHALL require every managed meme name to be 1 through 63 characters, contain only lowercase ASCII letters, decimal digits, and hyphens, and begin and end with an ASCII letter or decimal digit.

#### Scenario: An invalid name is supplied
- **WHEN** a user supplies an empty, uppercase, malformed, or overlength meme name
- **THEN** the system SHALL return a non-zero actionable error and SHALL NOT contact Imgflip or alter the inventory

### Requirement: Reject duplicate managed meme names
The system SHALL reject creation when the selected inventory already contains the supplied local meme name. It SHALL NOT overwrite the existing record or contact Imgflip for a rejected duplicate.

#### Scenario: A duplicate name is supplied
- **WHEN** a user creates a meme with a name already in the selected inventory
- **THEN** the system SHALL return a non-zero error identifying the duplicate name without creating a remote meme or changing the inventory

### Requirement: List managed memes
The system SHALL provide `memectl get memes` to list records in the selected local inventory without contacting Imgflip. Default output SHALL show `NAME`, `TEMPLATE ID`, `AGE`, and `IMAGE URL`; `--output wide` and `-o wide` SHALL additionally show `PAGE URL`. Records SHALL be ordered by name.

#### Scenario: Managed memes are listed
- **WHEN** the selected inventory contains managed meme records
- **THEN** `memectl get memes` SHALL display the records in name order using the requested output width

#### Scenario: No managed memes exist
- **WHEN** the selected inventory store does not exist or contains no records
- **THEN** `memectl get memes` SHALL succeed and print `No resources found.`

### Requirement: Delete managed meme records
The system SHALL provide `memectl delete meme <name> [<name>...]` to remove one or more named managed-meme records from the selected local inventory. It SHALL validate every supplied name, SHALL require at least one name, SHALL process supplied names in order as independent local inventory operations, and SHALL make no Imgflip request or alter the hosted image.

#### Scenario: Managed memes are deleted
- **WHEN** a user runs `memectl delete meme <name> [<name>...]` with valid names present in the selected inventory
- **THEN** the system SHALL remove only the named local records, report each deleted meme, and preserve all other records

#### Scenario: The final managed memes are deleted
- **WHEN** a user deletes all records in the selected inventory
- **THEN** the system SHALL persist a valid empty versioned inventory document

#### Scenario: Some named memes do not exist
- **WHEN** a user runs `memectl delete meme <name> [<name>...]` with valid names that are present and valid names that are absent from the selected inventory
- **THEN** the system SHALL remove and report each present local record, SHALL report every absent name as not found, and SHALL return a non-zero result

#### Scenario: A name is requested more than once
- **WHEN** a user runs `memectl delete meme <name> [<name>...]` with the same valid present name more than once
- **THEN** the system SHALL delete and report the first occurrence, SHALL report each later occurrence as not found, and SHALL return a non-zero result

#### Scenario: A later deletion fails
- **WHEN** a user runs `memectl delete meme <name> [<name>...]` and an inventory operation for a later valid name fails after earlier names were deleted
- **THEN** the system SHALL report the failure for that name, SHALL retain the successful earlier deletions, and SHALL return a non-zero result

#### Scenario: No named memes exist
- **WHEN** a user runs `memectl delete meme <name> [<name>...]` and every valid name is absent from the selected inventory
- **THEN** the system SHALL return a non-zero actionable not-found error and SHALL not modify the inventory

#### Scenario: Delete command arguments are invalid
- **WHEN** a user omits all meme names or supplies an invalid local meme name
- **THEN** the system SHALL return a non-zero actionable error and SHALL not modify the inventory

#### Scenario: An inventory operation fails before replacement
- **WHEN** an individual inventory operation cannot be read, is malformed, contains invalid records, has an unsupported format version, or fails before replacing the inventory document
- **THEN** the system SHALL return a non-zero actionable error for that name and SHALL preserve the inventory contents from before that operation

#### Scenario: Replacement succeeds but durable persistence cannot be confirmed
- **WHEN** an individual inventory operation replaces the inventory document but cannot confirm directory persistence afterward
- **THEN** the system SHALL return a non-zero actionable error that states the deletion may have succeeded but durable persistence could not be confirmed

### Requirement: Document local-only meme deletion
The README SHALL document `memectl delete meme <name> [<name>...]` and clarify that it removes only the selected local managed-meme records without contacting Imgflip or deleting the hosted image.

#### Scenario: Users discover deletion scope
- **WHEN** a user reads the README's managed-meme usage documentation
- **THEN** they can delete one or more local records and understand that the hosted Imgflip image remains unaffected

### Requirement: Preserve inventory integrity
The system SHALL treat a missing inventory store as empty. It SHALL return a non-zero actionable error without modifying the store when the selected store cannot be read, is malformed, contains invalid records, or has an unsupported format version. Updates SHALL replace the inventory atomically so an interrupted write does not expose partial inventory content.

#### Scenario: Inventory state is corrupt
- **WHEN** the selected inventory contains malformed or invalid state
- **THEN** the command SHALL fail with the inventory path and SHALL preserve the existing file

#### Scenario: Inventory persistence fails after remote creation
- **WHEN** Imgflip successfully creates a meme but the system cannot persist its managed record
- **THEN** the command SHALL return a non-zero result explaining that the remote meme was not recorded locally and SHALL report both generated URLs
