## 1. Add debug flag

- [x] 1.1 Add `--debug` bool flag to flag parsing in `main.go`
- [x] 1.2 Conditionally skip spinner when `--debug` is set
- [x] 1.3 Set `cmd.Stderr = os.Stderr` when `--debug` is set so Claude's streaming output is visible
- [x] 1.4 Print raw Claude response to stderr before parsing when `--debug` is set
