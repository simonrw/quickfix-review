## ADDED Requirements

### Requirement: Headless Claude invocation
The tool SHALL invoke Claude Code with `-p` and `--output-format json` flags, passing the review prompt as the argument to `-p`.

#### Scenario: Invoking Claude
- **WHEN** the tool starts a review for branch `foo`
- **THEN** it runs `claude -p "<prompt>" --output-format json` where the prompt includes the branch name and instructions to review against `origin/main`

### Requirement: Prompt includes branch context
The prompt SHALL instruct Claude to review the specified branch against `origin/main` and return a JSON array of review comments with `file`, `line`, `col`, and `message` fields.

#### Scenario: Prompt content
- **WHEN** the tool constructs the prompt for branch `feature-x`
- **THEN** the prompt tells Claude to review `feature-x` against `origin/main` and return results as a JSON array

### Requirement: Response parsing
The tool SHALL parse the Claude Code JSON envelope to extract the text response, then parse the JSON array of review comments from that text.

#### Scenario: Valid response
- **WHEN** Claude returns a valid JSON envelope containing a JSON array of comments
- **THEN** the tool extracts each comment's file, line, col, and message fields

#### Scenario: Unparseable response
- **WHEN** Claude's response cannot be parsed as expected JSON
- **THEN** the tool prints the raw response to stderr and exits with a non-zero exit code
