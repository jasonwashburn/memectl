# template-listing Specification

## Purpose

Provide a command-line catalog of public Imgflip meme templates so users can identify a template before generating a meme.

## Requirements

### Requirement: List public meme templates
The system SHALL provide `memectl get templates`, which retrieves the current public meme template list from Imgflip without requiring user credentials or local configuration.

#### Scenario: Templates are retrieved successfully
- **WHEN** a user runs `memectl get templates` and Imgflip returns a successful template list
- **THEN** the command SHALL complete successfully and display every returned template

### Requirement: Display template selection details
For each retrieved template, the system SHALL present its Imgflip identifier, name, text-box count, and image dimensions in a readable default tabular form. Dimensions SHALL be displayed as width by height.

#### Scenario: Template details are displayed
- **WHEN** Imgflip returns a template with an identifier, name, text-box count, width, and height and the user does not select wide output
- **THEN** the command output SHALL include those values in that template's default table row

### Requirement: Select template table output format
The system SHALL accept `--output wide` and `-o wide` as equivalent requests for wide template-table output. When neither form is provided, the system SHALL retain the existing default table output. The system SHALL return a non-zero error and SHALL NOT display template rows for an unsupported output format.

#### Scenario: Long output flag selects wide output
- **WHEN** a user runs `memectl get templates --output wide`
- **THEN** the command SHALL display the wide template table

#### Scenario: Short output flag selects wide output
- **WHEN** a user runs `memectl get templates -o wide`
- **THEN** the command SHALL display the same wide template table as the long flag

#### Scenario: Default output is unchanged
- **WHEN** a user runs `memectl get templates` without an output option
- **THEN** the command SHALL display the existing default table without a URL column

#### Scenario: Output format is unsupported
- **WHEN** a user requests an output format other than `wide`
- **THEN** the command SHALL return a non-zero error and SHALL NOT display template rows

### Requirement: Display template image URLs in wide output
When wide output is selected, the system SHALL display each template's direct Imgflip image URL in an additional `URL` table column alongside the default template selection details.

#### Scenario: Wide template details are displayed
- **WHEN** Imgflip returns a template with an identifier, name, text-box count, width, height, and image URL and the user selects wide output
- **THEN** the template's row SHALL include those values, with dimensions displayed as width by height and the image URL in the `URL` column

### Requirement: Document template listing
The README SHALL describe `memectl get templates` as the available template-discovery command and show how to run it.

#### Scenario: Users discover the command
- **WHEN** a user reads the README
- **THEN** they can identify and run `memectl get templates`

### Requirement: Report unavailable template data
The command SHALL return a non-zero result and an actionable error when it cannot retrieve or validate the template list. It SHALL NOT display partial or stale template rows in that case.

#### Scenario: Template retrieval fails
- **WHEN** the Imgflip request fails or Imgflip reports an unsuccessful response
- **THEN** the command SHALL return an error describing that templates could not be retrieved

#### Scenario: Template response is invalid
- **WHEN** Imgflip returns a successful HTTP response that cannot be interpreted as a valid template list
- **THEN** the command SHALL return an error and SHALL NOT display template rows

### Requirement: Bound default template requests
The default Imgflip client used by `memectl get templates` SHALL use a finite request timeout so an unresponsive endpoint cannot leave the command running indefinitely.

#### Scenario: Template request times out
- **WHEN** Imgflip does not respond before the default request timeout expires
- **THEN** the command SHALL return an error and SHALL NOT display template rows
