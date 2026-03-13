## Why

When Claude takes a long time or produces unexpected output, the spinner hides what's happening. A `--debug` flag lets you see Claude's full raw output in real time, which is essential for prompt iteration and troubleshooting.

## What Changes

- Add `--debug` flag to the CLI
- When `--debug` is set, skip the spinner and stream Claude's stderr to the terminal so the user can see Claude's progress
- After completion, print the raw JSON response to stderr before parsing

## Capabilities

### New Capabilities
- `debug-mode`: A `--debug` flag that shows Claude's full output instead of the spinner

### Modified Capabilities
<!-- None -->

## Impact

- `main.go`: Add flag, conditionally skip spinner, wire up stderr passthrough
