## Context

Currently `main.go` runs `cmd.Output()` which captures stdout and discards stderr. A spinner goroutine provides visual feedback. In debug mode, we want to see Claude's stderr streaming output and the full response.

## Goals / Non-Goals

**Goals:**
- See Claude's streaming output in real time when `--debug` is passed
- See the raw JSON response before it's parsed

**Non-Goals:**
- Structured debug logging
- Verbosity levels

## Decisions

### Stream Claude's stderr in debug mode
When `--debug` is set, set `cmd.Stderr = os.Stderr` so Claude's streaming output is visible. Skip the spinner entirely — the streaming output serves as progress indication.

### Print raw response in debug mode
After Claude completes, print the full captured stdout to stderr before parsing. This helps diagnose JSON extraction issues.

### Use cmd.Output() in both modes
`cmd.Output()` captures stdout regardless of stderr wiring. The only difference is whether stderr is passed through or discarded, and whether we print the raw response.

## Risks / Trade-offs

- Debug output may be noisy — acceptable since it's opt-in.
