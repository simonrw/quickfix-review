## ADDED Requirements

### Requirement: Print flag outputs quickfix lines to stdout
When invoked with `--print`, the binary SHALL write quickfix-format lines to stdout and exit without launching neovim or writing a temp file.

#### Scenario: Basic print mode
- **WHEN** the binary is invoked with `--print`
- **THEN** quickfix-format lines are written to stdout and the process exits normally

#### Scenario: Print mode with branch flag
- **WHEN** the binary is invoked with `--print --branch feature-xyz`
- **THEN** the review targets the specified branch and quickfix lines are written to stdout

### Requirement: Print mode suppresses stderr output
When `--print` is set, the binary SHALL NOT write spinner output or other progress indicators to stderr.

#### Scenario: No spinner in print mode
- **WHEN** the binary is invoked with `--print`
- **THEN** stderr remains empty during execution

### Requirement: Print mode exits with appropriate status
The binary SHALL exit with status 0 on success (even if no review comments are found) and non-zero on error.

#### Scenario: Successful review with comments
- **WHEN** the binary completes a review that finds issues
- **THEN** it exits with status 0 and quickfix lines on stdout

#### Scenario: Successful review with no comments
- **WHEN** the binary completes a review that finds no issues
- **THEN** it exits with status 0 and empty stdout

#### Scenario: Claude invocation fails
- **WHEN** Claude fails to respond or returns invalid output
- **THEN** the binary exits with non-zero status
