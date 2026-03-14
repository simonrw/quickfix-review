package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type reviewComment struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Col     int    `json:"col"`
	Message string `json:"message"`
}

type claudeResponse struct {
	Result string `json:"result"`
}

func main() {
	branch := flag.String("branch", "", "branch to review (default: current branch)")
	debug := flag.Bool("debug", false, "show claude's full output instead of spinner")
	printOnly := flag.Bool("print", false, "print quickfix lines to stdout without launching nvim")
	flag.Parse()

	// Pre-flight checks
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: claude not found on PATH")
		os.Exit(1)
	}
	var nvimPath string
	if !*printOnly {
		nvimPath, err = exec.LookPath("nvim")
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: nvim not found on PATH")
			os.Exit(1)
		}
	}

	// Determine branch
	targetBranch := *branch
	if targetBranch == "" {
		out, err := exec.Command("git", "branch", "--show-current").Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: could not determine current branch:", err)
			os.Exit(1)
		}
		targetBranch = strings.TrimSpace(string(out))
		if targetBranch == "" {
			fmt.Fprintln(os.Stderr, "error: not on a branch (detached HEAD?)")
			os.Exit(1)
		}
	}

	// Build prompt
	prompt := fmt.Sprintf(reviewPromptTemplate, targetBranch, targetBranch)

	// Invoke Claude
	var out []byte
	if *debug {
		fmt.Fprintf(os.Stderr, "[prompt] %s\n\n", prompt)
		out, err = runClaudeDebug(claudePath, prompt)
	} else if *printOnly {
		out, err = runClaudeSilent(claudePath, prompt)
	} else {
		out, err = runClaudeQuiet(claudePath, prompt)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: claude exited with error: %v\n", err)
		if len(out) > 0 {
			fmt.Fprintf(os.Stderr, "output:\n%s\n", out)
		}
		os.Exit(1)
	}

	// Parse Claude's JSON envelope
	var envelope claudeResponse
	if err := json.Unmarshal(out, &envelope); err != nil {
		fmt.Fprintf(os.Stderr, "error: could not parse claude response envelope: %v\n", err)
		fmt.Fprintf(os.Stderr, "raw output:\n%s\n", out)
		os.Exit(1)
	}

	// Extract JSON array from the result text
	commentJSON := extractJSON(envelope.Result)

	var comments []reviewComment
	if err := json.Unmarshal([]byte(commentJSON), &comments); err != nil {
		fmt.Fprintf(os.Stderr, "error: could not parse review comments: %v\n", err)
		fmt.Fprintf(os.Stderr, "result text:\n%s\n", envelope.Result)
		os.Exit(1)
	}

	if len(comments) == 0 {
		if !*printOnly {
			fmt.Fprintln(os.Stderr, "no review comments found")
		}
		os.Exit(0)
	}

	// Build quickfix lines
	var lines []string
	for _, c := range comments {
		msg := strings.ReplaceAll(c.Message, "\n", " ")
		msg = strings.ReplaceAll(msg, "\r", " ")
		lines = append(lines, fmt.Sprintf("%s:%d:%d:%s", c.File, c.Line, c.Col, msg))
	}

	if *printOnly {
		fmt.Print(strings.Join(lines, "\n") + "\n")
		return
	}

	// Write temp file
	tmpFile, err := os.CreateTemp("", "llmreview-*.qf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not create temp file: %v\n", err)
		os.Exit(1)
	}
	_, err = tmpFile.WriteString(strings.Join(lines, "\n") + "\n")
	tmpFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not write temp file: %v\n", err)
		os.Exit(1)
	}

	// Exec nvim
	err = syscall.Exec(nvimPath, []string{"nvim", "-q", tmpFile.Name()}, os.Environ())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not exec nvim: %v\n", err)
		os.Exit(1)
	}
}

// runClaudeSilent runs claude with --output-format json and no spinner or stderr output.
func runClaudeSilent(claudePath, prompt string) ([]byte, error) {
	cmd := exec.Command(claudePath, "-p", prompt, "--output-format", "json")
	return cmd.Output()
}

// runClaudeQuiet runs claude with --output-format json and shows a spinner.
func runClaudeQuiet(claudePath, prompt string) ([]byte, error) {
	cmd := exec.Command(claudePath, "-p", prompt, "--output-format", "json")
	stop := make(chan struct{})
	go spinner(stop)
	out, err := cmd.Output()
	close(stop)
	fmt.Fprint(os.Stderr, "\r\033[K")
	return out, err
}

// streamLine represents a line from --output-format stream-json.
type streamLine struct {
	Type    string          `json:"type"`
	Message json.RawMessage `json:"message,omitempty"`
	Result  string          `json:"result,omitempty"`
}

// messageContent is a content block in an assistant/user message.
type messageContent struct {
	Type     string          `json:"type"`
	Text     string          `json:"text,omitempty"`
	Thinking string          `json:"thinking,omitempty"` // thinking block content
	Name     string          `json:"name,omitempty"`
	Input    json.RawMessage `json:"input,omitempty"`
	Content  string          `json:"content,omitempty"` // tool_result text content
}

// toolInput holds common tool input fields we want to display.
type toolInput struct {
	Command  string `json:"command,omitempty"`  // Bash
	FilePath string `json:"file_path,omitempty"` // Read/Edit
	Pattern  string `json:"pattern,omitempty"`  // Grep/Glob
}

// messageBody is the structure of an assistant or user message.
type messageBody struct {
	Content []messageContent `json:"content"`
}

// runClaudeDebug runs claude with --output-format stream-json, streaming
// text deltas to stderr in real time. Returns a JSON envelope compatible
// with the non-debug path (i.e. {"result": "..."}).
func runClaudeDebug(claudePath, prompt string) ([]byte, error) {
	cmd := exec.Command(claudePath, "-p", prompt, "--output-format", "stream-json", "--verbose")
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start claude: %w", err)
	}

	var resultText string
	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024) // 10MB max line
	for scanner.Scan() {
		line := scanner.Bytes()
		var ev streamLine
		if err := json.Unmarshal(line, &ev); err != nil {
			continue
		}

		switch ev.Type {
		case "assistant":
			var msg messageBody
			if err := json.Unmarshal(ev.Message, &msg); err != nil {
				continue
			}
			for _, block := range msg.Content {
				switch block.Type {
				case "tool_use":
					var ti toolInput
					json.Unmarshal(block.Input, &ti)
					switch {
					case ti.Command != "":
						fmt.Fprintf(os.Stderr, "[%s] %s\n", block.Name, ti.Command)
					case ti.FilePath != "":
						fmt.Fprintf(os.Stderr, "[%s] %s\n", block.Name, ti.FilePath)
					case ti.Pattern != "":
						fmt.Fprintf(os.Stderr, "[%s] %s\n", block.Name, ti.Pattern)
					default:
						fmt.Fprintf(os.Stderr, "[%s]\n", block.Name)
					}
				case "thinking":
					if block.Thinking != "" {
						fmt.Fprintf(os.Stderr, "[thinking] %s\n", block.Thinking)
					}
				case "text":
					if block.Text != "" {
						fmt.Fprintf(os.Stderr, "%s\n", block.Text)
					}
				}
			}
		case "user":
			var msg messageBody
			if err := json.Unmarshal(ev.Message, &msg); err != nil {
				continue
			}
			for _, block := range msg.Content {
				if block.Type == "tool_result" && block.Content != "" {
					// Truncate long tool results
					content := block.Content
					if len(content) > 500 {
						content = content[:500] + "...(truncated)"
					}
					fmt.Fprintf(os.Stderr, "  → %s\n", content)
				}
			}
		case "result":
			resultText = ev.Result
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "[debug] scanner error: %v\n", err)
	}
	fmt.Fprintln(os.Stderr)

	waitErr := cmd.Wait()

	// Wrap in the same envelope format as --output-format json
	envelope, _ := json.Marshal(claudeResponse{Result: resultText})
	return envelope, waitErr
}

func spinner(stop chan struct{}) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-stop:
			return
		default:
			fmt.Fprintf(os.Stderr, "\r%s Reviewing...", frames[i%len(frames)])
			i++
			time.Sleep(80 * time.Millisecond)
		}
	}
}

// extractJSON finds the first JSON array in the text, handling cases where
// Claude may wrap the array in markdown code fences or add surrounding text.
func extractJSON(text string) string {
	// Try the whole text first
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "[") {
		return text
	}

	// Look for ```json ... ``` blocks
	if idx := strings.Index(text, "```json"); idx != -1 {
		start := idx + len("```json")
		if end := strings.Index(text[start:], "```"); end != -1 {
			return strings.TrimSpace(text[start : start+end])
		}
	}

	// Look for ``` ... ``` blocks
	if idx := strings.Index(text, "```"); idx != -1 {
		start := idx + len("```")
		if end := strings.Index(text[start:], "```"); end != -1 {
			return strings.TrimSpace(text[start : start+end])
		}
	}

	// Look for first [ ... last ]
	start := strings.Index(text, "[")
	end := strings.LastIndex(text, "]")
	if start != -1 && end != -1 && end > start {
		return text[start : end+1]
	}

	return text
}
