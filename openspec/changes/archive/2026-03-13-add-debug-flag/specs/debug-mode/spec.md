## ADDED Requirements

### Requirement: Debug flag
The tool SHALL accept a `--debug` flag.

#### Scenario: Flag not provided
- **WHEN** the user runs `llmreview` without `--debug`
- **THEN** the tool shows the spinner and hides Claude's output (existing behavior)

#### Scenario: Flag provided
- **WHEN** the user runs `llmreview --debug`
- **THEN** the tool skips the spinner and streams Claude's stderr to the terminal

### Requirement: Raw response output in debug mode
When `--debug` is set, the tool SHALL print the raw Claude JSON response to stderr before parsing.

#### Scenario: Debug with valid response
- **WHEN** `--debug` is set and Claude returns a valid response
- **THEN** the raw JSON response is printed to stderr, then parsing proceeds normally

#### Scenario: Debug with invalid response
- **WHEN** `--debug` is set and Claude returns an unparseable response
- **THEN** the raw JSON response is printed to stderr (already visible), then the parse error is shown
