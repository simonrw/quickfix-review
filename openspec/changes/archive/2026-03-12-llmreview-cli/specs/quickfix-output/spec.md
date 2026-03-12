## ADDED Requirements

### Requirement: Quickfix format output
The tool SHALL transform each review comment into a vim quickfix line in the format `<filename>:<line>:<col>:<message>`.

#### Scenario: Comment transformation
- **WHEN** a review comment has file `src/main.go`, line `42`, col `1`, and message `Error not handled`
- **THEN** the quickfix line is `src/main.go:42:1:Error not handled`

### Requirement: Single-line messages
Each quickfix entry SHALL be on a single line with no embedded newlines in the message field.

#### Scenario: Message with newlines
- **WHEN** Claude returns a message containing newline characters
- **THEN** the newlines are replaced with spaces in the quickfix output

### Requirement: Temp file output
The tool SHALL write quickfix lines to a temporary file in the system temp directory with a `.qf` extension.

#### Scenario: File creation
- **WHEN** the review produces comments
- **THEN** a file is created at a path like `/tmp/llmreview-XXXXX.qf` containing one quickfix line per comment
