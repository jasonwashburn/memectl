## Purpose

Enable users to create hosted, captioned static Imgflip memes from the command line after identifying a suitable template.

## ADDED Requirements

### Requirement: Create a captioned image meme
The system SHALL provide `memectl create meme <template-id>` to create a captioned static image meme from an Imgflip template. The command SHALL require exactly one template identifier and one or more repeatable `--text <text>` flags.

#### Scenario: Captioned meme is created
- **WHEN** a user supplies a template identifier, at least one `--text` value, and valid Imgflip credentials
- **THEN** the system SHALL create a captioned image from that template and report success

#### Scenario: No caption text is supplied
- **WHEN** a user runs `memectl create meme <template-id>` without a `--text` flag
- **THEN** the system SHALL return a non-zero result that explains that at least one caption text value is required

#### Scenario: Command arguments are invalid
- **WHEN** a user supplies no template identifier or more than one template identifier
- **THEN** the system SHALL return a non-zero result and SHALL NOT create a meme

### Requirement: Preserve caption text order
The system SHALL send each supplied `--text` value as an ordered text box for the selected template. The first supplied value SHALL be the first text box, and subsequent values SHALL retain their relative order.

#### Scenario: Multiple caption texts are supplied
- **WHEN** a user supplies multiple `--text` flags
- **THEN** the created meme request SHALL contain the text values in the same order as the flags

### Requirement: Obtain creation credentials from the environment
The system SHALL obtain the Imgflip account username from `IMGFLIP_USERNAME` and password from `IMGFLIP_PASSWORD` for captioned-meme creation. It SHALL NOT require credentials for template listing.

#### Scenario: Credentials are available
- **WHEN** both required environment variables contain values
- **THEN** the system SHALL use them to authenticate the captioned-meme creation request

#### Scenario: A credential is missing
- **WHEN** either required environment variable is unset or empty
- **THEN** the command SHALL return a non-zero actionable error and SHALL NOT send a meme-creation request

### Requirement: Report the generated hosted meme
After Imgflip successfully creates a captioned meme, the command SHALL print a concise summary identifying the created meme's template and SHALL print the direct hosted image URL and the Imgflip page URL returned by Imgflip.

#### Scenario: Imgflip returns generated URLs
- **WHEN** Imgflip reports successful creation with a direct image URL and an Imgflip page URL
- **THEN** the command SHALL complete successfully and display the summary and both URLs on standard output

### Requirement: Report failed meme creation safely
The system SHALL return a non-zero actionable error when Imgflip cannot be reached, returns a non-success HTTP response, returns an invalid response, or reports unsuccessful caption creation. It SHALL NOT print a creation summary or either URL in those cases, and SHALL NOT expose the configured password in output.

#### Scenario: Imgflip rejects the request
- **WHEN** Imgflip returns an unsuccessful caption-creation response
- **THEN** the command SHALL return an error that indicates creation failed without printing a success summary or either URL

#### Scenario: Imgflip response cannot be used
- **WHEN** the caption-creation response is malformed or lacks a generated image URL
- **THEN** the command SHALL return an error without printing a success summary or either URL

### Requirement: Document captioned-meme creation
The README SHALL document the required environment variables and show how to run `memectl create meme <template-id>` with repeatable `--text` flags.

#### Scenario: Users discover creation setup
- **WHEN** a user reads the README
- **THEN** they can identify the credentials to configure and invoke the captioned-meme creation command
