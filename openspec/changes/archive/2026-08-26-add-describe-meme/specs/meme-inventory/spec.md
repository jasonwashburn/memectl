## ADDED Requirements

### Requirement: Describe a managed meme record
The system SHALL provide `memectl describe meme <name>` and the equivalent
`memectl desc meme <name>` alias to inspect exactly one named record in the
selected local inventory. The command SHALL not contact Imgflip, require
Imgflip credentials, modify the inventory, or claim the current availability
of a hosted image.

#### Scenario: A managed meme is described
- **WHEN** a user describes a valid name present in the selected inventory
- **THEN** the command SHALL succeed and display a labeled detail view of that
  record

#### Scenario: Describe is invoked through its alias
- **WHEN** a user runs `memectl desc meme <name>` for a valid name present in
  the selected inventory
- **THEN** the command SHALL display the same detail view as `memectl describe
  meme <name>`

### Requirement: Display stored managed meme details
The detail view SHALL display the stored name, template identifier, ordered
caption text values, direct image URL, Imgflip page URL, and creation timestamp
in a labeled multiline format. Caption text values SHALL be shown by their
zero-based text-box index in stored order, including empty values. The creation
timestamp SHALL be shown as the stored UTC timestamp rather than a relative
age.

#### Scenario: A record has multiple or empty caption texts
- **WHEN** a described managed-meme record contains multiple caption text
  values, including an empty value
- **THEN** the detail view SHALL show every value in order with its text-box
  index and preserve the empty value

### Requirement: Report invalid describe requests
The describe command SHALL require exactly one local meme name and SHALL
validate it using the managed-meme name rules. It SHALL return a non-zero,
actionable error without contacting Imgflip or modifying the inventory when the
name is omitted, more than one name is supplied, or the supplied name is
invalid. It SHALL return a non-zero actionable not-found error when a valid
requested name is absent from the selected inventory, including when the store
does not exist.

#### Scenario: The requested meme is absent
- **WHEN** a user describes a valid local meme name that is not in the selected
  inventory
- **THEN** the command SHALL return a non-zero not-found error and produce no
  detail view

#### Scenario: Describe arguments are invalid
- **WHEN** a user omits the meme name, provides more than one name, or provides
  an invalid local meme name
- **THEN** the command SHALL return a non-zero actionable error and produce no
  detail view

### Requirement: Document local-only meme inspection
The README SHALL document `memectl describe meme <name>` and clarify that it
reads the selected local managed-meme record without contacting Imgflip.

#### Scenario: Users discover describe scope
- **WHEN** a user reads the README's managed-meme usage documentation
- **THEN** they can inspect a named local record and understand that describe
  does not query Imgflip
