## Why

The tool currently launches a new neovim instance via `syscall.Exec`, which means users already in neovim must leave their editor to run a review. Running from within neovim keeps the user in their existing session with all their buffers and state.

## What Changes

- Add `--print` flag to the Go binary that writes quickfix lines to stdout instead of exec'ing nvim, with no spinner output on stderr
- Add `quickfix-review.lua` companion file providing a `:QuickfixReview [branch]` user command that runs the binary asynchronously and loads results into neovim's quickfix list

## Capabilities

### New Capabilities
- `print-mode`: `--print` flag that outputs quickfix-format lines to stdout without launching neovim
- `neovim-plugin`: Lua user command for running reviews asynchronously from within neovim

### Modified Capabilities

## Impact

- `main.go`: new flag, conditional skip of spinner and nvim exec when `--print` is set
- New file: `quickfix-review.lua` in repo root
- No breaking changes to existing CLI behavior
