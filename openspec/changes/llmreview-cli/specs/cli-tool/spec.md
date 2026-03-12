## ADDED Requirements

### Requirement: Default branch detection
The tool SHALL detect the current git branch using `git branch --show-current` when no `--branch` flag is provided.

#### Scenario: No flag provided
- **WHEN** the user runs `llmreview` with no arguments
- **THEN** the tool uses the current git branch as the review target

#### Scenario: Branch flag provided
- **WHEN** the user runs `llmreview --branch foo`
- **THEN** the tool uses `foo` as the review target

### Requirement: Pre-flight checks
The tool SHALL verify that `claude` and `nvim` are available on PATH before proceeding.

#### Scenario: claude not on PATH
- **WHEN** the `claude` binary is not found on PATH
- **THEN** the tool exits with a non-zero exit code and prints an error message to stderr

#### Scenario: nvim not on PATH
- **WHEN** the `nvim` binary is not found on PATH
- **THEN** the tool exits with a non-zero exit code and prints an error message to stderr

### Requirement: Loading spinner
The tool SHALL display an indeterminate loading spinner on stderr while Claude is running.

#### Scenario: Review in progress
- **WHEN** the tool is waiting for Claude to complete
- **THEN** a spinner animation is displayed on stderr

#### Scenario: Review complete
- **WHEN** Claude returns a response
- **THEN** the spinner is stopped and cleared from stderr

### Requirement: Open neovim with results
The tool SHALL exec neovim with the quickfix file using `nvim -q <tempfile>`, replacing its own process.

#### Scenario: Successful review
- **WHEN** Claude returns valid review comments
- **THEN** the tool writes the quickfix file and execs `nvim -q <tempfile>`
