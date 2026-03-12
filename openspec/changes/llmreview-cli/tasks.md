## 1. Project Setup

- [x] 1.1 Initialize Go module (`go mod init`) and create `main.go` and `prompt.go` files
- [x] 1.2 Implement CLI flag parsing (`--branch`) and current branch detection via `git branch --show-current`
- [x] 1.3 Implement pre-flight checks for `claude` and `nvim` on PATH

## 2. Claude Integration

- [x] 2.1 Define the review prompt as a string constant in `prompt.go` with a `%s` placeholder for the branch name
- [x] 2.2 Implement Claude Code invocation: shell out to `claude -p "<prompt>" --output-format json` and capture stdout
- [x] 2.3 Implement stderr spinner (goroutine printing spinner frames while Claude runs)

## 3. Output Processing

- [x] 3.1 Parse the Claude Code JSON envelope to extract the text response
- [x] 3.2 Parse the JSON array of review comments from the text response (extract `file`, `line`, `col`, `message`)
- [x] 3.3 Transform comments into quickfix format lines, replacing any newlines in messages with spaces
- [x] 3.4 Write quickfix lines to a temp file (`os.CreateTemp` with `.qf` suffix)

## 4. Editor Launch

- [x] 4.1 Exec nvim with `-q <tempfile>` using `syscall.Exec` to replace the process
- [x] 4.2 Handle error case: if JSON parsing fails, print raw response to stderr and exit non-zero
