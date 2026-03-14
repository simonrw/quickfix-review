## Context

The binary currently runs Claude for code review, writes quickfix-format lines to a temp file, then `syscall.Exec`s into nvim with `-q`. This works from a terminal but not from within an already-running neovim session.

## Goals / Non-Goals

**Goals:**
- Allow the binary to output quickfix lines to stdout without launching nvim
- Provide a Lua file that integrates with neovim's job control and quickfix list

**Non-Goals:**
- Supporting `--debug` streaming output from within neovim
- Creating a full neovim plugin with install/setup infrastructure
- Modifying the review prompt or Claude interaction

## Decisions

**`--print` flag skips spinner and nvim exec.** When `--print` is set, the binary writes quickfix lines directly to stdout and exits. No spinner on stderr, no temp file, no `syscall.Exec`. This keeps stdout clean for programmatic consumption. Alternative considered: writing to a named file and printing the path — adds unnecessary coordination between binary and Lua.

**Async execution via `vim.fn.jobstart`.** The Lua command uses neovim's async job control rather than synchronous `vim.fn.system()`. Reviews take 30+ seconds; blocking the editor is unacceptable. The tradeoff is slightly more complex Lua (buffering stdout, handling exit callback) but this is straightforward with neovim's API.

**Raw quickfix lines on stdout, not JSON.** The binary already formats quickfix lines internally. Outputting them directly means the Lua side can use `vim.fn.setqflist({}, ' ', {lines = ...})` with no parsing. Alternative considered: JSON output for richer metadata — not needed now and easy to add later.

**Single Lua file in repo root.** Users `require` or `source` it from their nvim config. No plugin manager boilerplate. Can be moved into a proper plugin structure later if needed.

## Risks / Trade-offs

**Binary must be on PATH** → The Lua file assumes `quickfix-review` is available. Users who build but don't install will need to configure the path. Could add a `vim.g.quickfix_review_cmd` override.

**No progress indication during review** → User gets a notify at start and end but nothing in between. Acceptable for now; could stream stderr for progress later.
