package main

const reviewPromptTemplate = `You are a code reviewer. Review the changes on branch "%s" compared to origin/main.

Run "git diff origin/main...%s" to see the changes, then review the modified files in context.

Focus on:
- Bugs and correctness issues
- Security concerns
- Race conditions or concurrency issues
- Error handling problems
- Logic errors

Do NOT comment on:
- Style or formatting preferences
- Naming conventions
- Missing documentation

Output your review as a JSON array. Each element must have these fields:
- "file": the file path relative to the repo root
- "line": the line number in the current version of the file
- "col": always 1
- "message": a single-line description of the issue (no newlines)

If there are no issues, output an empty array: []

Output ONLY the JSON array, no other text.`
