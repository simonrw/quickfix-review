## Why

LLM code review feedback is most useful when integrated directly into the editor workflow. Currently there's no way to get AI review comments loaded into vim's quickfix list, where developers can jump between them with `:cnext`/`:cprev` just like compiler errors.

## What Changes

- New Go CLI tool `llmreview` that shells out to Claude Code in headless mode (`-p` flag) to review the current branch against `origin/main`
- Claude's structured JSON response is parsed into vim quickfix format (`file:line:col:message`)
- Results are written to a temp file and neovim is exec'd with `-q` to load them
- A loading spinner is displayed on stderr while Claude is working

## Capabilities

### New Capabilities
- `cli-tool`: Go binary with `--branch` flag support, branch detection, spinner, and nvim exec
- `claude-integration`: Headless Claude Code invocation with structured JSON output parsing
- `quickfix-output`: Transformation of review comments into vim quickfix format

### Modified Capabilities
<!-- None — this is a greenfield project -->

## Impact

- New Go module and binary at the project root
- Runtime dependencies: `claude` CLI and `nvim` on PATH
- No existing code is modified
