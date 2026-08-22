## Purpose

Provide a command-line catalog of public Imgflip meme templates so users can identify a template before generating a meme.

## ADDED Requirements

### Requirement: List public meme templates
The system SHALL provide `memectl get templates`, which retrieves the current public meme template list from Imgflip without requiring user credentials or local configuration.

#### Scenario: Templates are retrieved successfully
- **WHEN** a user runs `memectl get templates` and Imgflip returns a successful template list
- **THEN** the command SHALL complete successfully and display every returned template

### Requirement: Display template selection details
For each retrieved template, the command SHALL present its Imgflip identifier, name, text-box count, and image dimensions in a readable tabular form. Dimensions SHALL be displayed as width by height.

#### Scenario: Template details are displayed
- **WHEN** Imgflip returns a template with an identifier, name, text-box count, width, and height
- **THEN** the command output SHALL include those values in that template's table row

### Requirement: Report unavailable template data
The command SHALL return a non-zero result and an actionable error when it cannot retrieve or validate the template list. It SHALL NOT display partial or stale template rows in that case.

#### Scenario: Template retrieval fails
- **WHEN** the Imgflip request fails or Imgflip reports an unsuccessful response
- **THEN** the command SHALL return an error describing that templates could not be retrieved

#### Scenario: Template response is invalid
- **WHEN** Imgflip returns a successful HTTP response that cannot be interpreted as a valid template list
- **THEN** the command SHALL return an error and SHALL NOT display template rows
