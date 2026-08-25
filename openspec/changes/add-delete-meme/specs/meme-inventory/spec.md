## ADDED Requirements

### Requirement: Delete managed meme records
The system SHALL provide `memectl delete meme <name> [<name>...]` to remove one or more named managed-meme records from the selected local inventory. It SHALL validate every supplied name, SHALL require at least one name, and SHALL make no Imgflip request or alter the hosted image.

#### Scenario: Managed memes are deleted
- **WHEN** a user runs `memectl delete meme <name> [<name>...]` with valid names present in the selected inventory
- **THEN** the system SHALL remove only the named local records, report each deleted meme, and preserve all other records

#### Scenario: The final managed memes are deleted
- **WHEN** a user deletes all records in the selected inventory
- **THEN** the system SHALL persist a valid empty versioned inventory document

#### Scenario: Some named memes do not exist
- **WHEN** a user runs `memectl delete meme <name> [<name>...]` with valid names that are present and valid names that are absent from the selected inventory
- **THEN** the system SHALL remove and report each present local record, SHALL report every absent name as not found, and SHALL return a non-zero result

#### Scenario: No named memes exist
- **WHEN** a user runs `memectl delete meme <name> [<name>...]` and every valid name is absent from the selected inventory
- **THEN** the system SHALL return a non-zero actionable not-found error and SHALL not modify the inventory

#### Scenario: Delete command arguments are invalid
- **WHEN** a user omits all meme names or supplies an invalid local meme name
- **THEN** the system SHALL return a non-zero actionable error and SHALL not modify the inventory

#### Scenario: Inventory state cannot be safely updated
- **WHEN** the selected inventory cannot be read, is malformed, contains invalid records, has an unsupported format version, or cannot be persisted
- **THEN** the system SHALL return a non-zero actionable error and SHALL preserve the existing inventory contents

### Requirement: Document local-only meme deletion
The README SHALL document `memectl delete meme <name> [<name>...]` and clarify that it removes only the selected local managed-meme records without contacting Imgflip or deleting the hosted image.

#### Scenario: Users discover deletion scope
- **WHEN** a user reads the README's managed-meme usage documentation
- **THEN** they can delete one or more local records and understand that the hosted Imgflip image remains unaffected
