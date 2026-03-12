## Context

Greenfield Go CLI tool. No existing code or architecture to integrate with. The tool orchestrates two external programs: `claude` (Claude Code CLI) for AI review and `nvim` for displaying results.

## Goals / Non-Goals

**Goals:**
- Single Go binary with minimal dependencies (standard library only where possible)
- Simple, predictable UX: run command → spinner → editor opens with comments
- Robust parsing of Claude's JSON output with clear error messages on failure

**Non-Goals:**
- PR mode or GitHub integration
- Configurable prompts or review styles
- Support for editors other than neovim
- Diff size management (delegated to Claude)
- Streaming or incremental output

## Decisions

### Single-file structure with separate prompt constant
All Go code lives in `main.go` with the prompt in `prompt.go`. No packages, no `cmd/` directory. This is a small tool and splitting early adds ceremony without value.

**Alternative**: `cobra` or `flag`-based subcommands — rejected because there's only one command with one optional flag.

### Use `--output-format json` from Claude Code
Claude Code's `--output-format json` wraps the response in a structured envelope with a `result` field. The prompt instructs Claude to return a JSON array within its text response, and the Go binary extracts it from the envelope. This is more reliable than parsing free-form text.

**Alternative**: Parse Claude's raw text output for JSON — rejected because the envelope is more predictable.

### Prompt instructs Claude to do its own diffing
Rather than computing the diff in Go and passing it to Claude, the prompt tells Claude to review the branch against `origin/main`. Claude Code has full tool access in `-p` mode and can run `git diff`, read files, and explore context as needed. This keeps the Go code simple and lets Claude decide how to handle large diffs.

**Alternative**: Compute diff in Go and pass via stdin — rejected because it limits Claude's ability to explore context and requires diff size management in Go.

### `exec` syscall to replace process with nvim
Use `syscall.Exec` to replace the llmreview process with nvim rather than spawning a subprocess. This means nvim inherits the terminal cleanly with no wrapper process hanging around.

### Temp file for quickfix output
Write quickfix lines to `/tmp/llmreview-XXXXX.qf` using `os.CreateTemp`. The file persists after nvim exits so the user can re-open it if needed.

## Risks / Trade-offs

- **Claude output format instability** → The prompt asks for a specific JSON schema but LLMs can deviate. Mitigation: validate the parsed JSON and show the raw response on parse failure so the user can see what went wrong.
- **`claude` CLI not installed** → Check for `claude` on PATH before invoking. Clear error message.
- **Large branches** → Claude may time out or produce low-quality reviews on very large diffs. Mitigation: none for v1, accepted trade-off. Can add `--max-files` later.
