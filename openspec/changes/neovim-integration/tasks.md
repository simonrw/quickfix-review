## 1. Print Mode

- [x] 1.1 Add `--print` boolean flag to main.go
- [x] 1.2 When `--print` is set, skip spinner and write quickfix lines to stdout instead of temp file
- [x] 1.3 When `--print` is set, skip the `syscall.Exec` into nvim and exit normally

## 2. Neovim Plugin

- [x] 2.1 Create `quickfix-review.lua` in repo root with `:QuickfixReview [branch]` user command
- [x] 2.2 Implement async execution via `vim.fn.jobstart` with stdout buffering
- [x] 2.3 Implement `on_exit` handler: load quickfix list on success, notify on empty results or error
