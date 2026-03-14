vim.api.nvim_create_user_command("QuickfixReview", function(opts)
  local cmd = { "quickfix-review", "--print" }
  if opts.args ~= "" then
    table.insert(cmd, "--branch")
    table.insert(cmd, opts.args)
  end

  local lines = {}

  vim.notify("Starting review...", vim.log.levels.INFO)

  vim.fn.jobstart(cmd, {
    stdout_buffered = true,
    on_stdout = function(_, data)
      for _, line in ipairs(data) do
        if line ~= "" then
          table.insert(lines, line)
        end
      end
    end,
    on_exit = function(_, exit_code)
      if exit_code ~= 0 then
        vim.notify("Review failed (exit code " .. exit_code .. ")", vim.log.levels.ERROR)
        return
      end
      if #lines == 0 then
        vim.notify("No issues found", vim.log.levels.INFO)
        return
      end
      vim.fn.setqflist({}, " ", { lines = lines, title = "Code Review" })
      vim.cmd("copen")
      vim.notify("Review complete: " .. #lines .. " issue(s)", vim.log.levels.INFO)
    end,
  })
end, {
  nargs = "?",
  desc = "Run quickfix-review and load results into quickfix list",
})
