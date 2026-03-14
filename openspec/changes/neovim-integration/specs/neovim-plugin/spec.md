## ADDED Requirements

### Requirement: QuickfixReview user command
The Lua file SHALL register a `:QuickfixReview` user command that accepts an optional branch name argument.

#### Scenario: Command with branch argument
- **WHEN** user runs `:QuickfixReview feature-xyz`
- **THEN** the binary is invoked with `--print --branch feature-xyz`

#### Scenario: Command without argument
- **WHEN** user runs `:QuickfixReview` with no arguments
- **THEN** the binary is invoked with `--print` only (binary defaults to current branch)

### Requirement: Asynchronous execution
The command SHALL run the binary asynchronously so the editor remains responsive during the review.

#### Scenario: Editor stays responsive
- **WHEN** a review is in progress
- **THEN** the user can continue editing, navigating, and using neovim normally

### Requirement: Start notification
The command SHALL notify the user when a review starts.

#### Scenario: Review starts
- **WHEN** user runs `:QuickfixReview`
- **THEN** `vim.notify` displays a message indicating the review has started

### Requirement: Results loaded into quickfix list
On successful completion with results, the command SHALL populate the quickfix list and open it.

#### Scenario: Review finds issues
- **WHEN** the binary exits successfully with quickfix lines on stdout
- **THEN** the quickfix list is populated with the results and `:copen` is called

### Requirement: Empty results notification
On successful completion with no results, the command SHALL notify the user.

#### Scenario: Review finds no issues
- **WHEN** the binary exits successfully with no stdout output
- **THEN** `vim.notify` displays "No issues found"

### Requirement: Error notification
On non-zero exit, the command SHALL notify the user of the error.

#### Scenario: Binary fails
- **WHEN** the binary exits with non-zero status
- **THEN** `vim.notify` displays an error message with `vim.log.levels.ERROR`
