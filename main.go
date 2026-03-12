package main

import (
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
	flag.Parse()

	// Pre-flight checks
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: claude not found on PATH")
		os.Exit(1)
	}
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: nvim not found on PATH")
		os.Exit(1)
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

	// Start spinner
	stop := make(chan struct{})
	go spinner(stop)

	// Invoke Claude
	cmd := exec.Command(claudePath, "-p", prompt, "--output-format", "json")
	out, err := cmd.Output()

	close(stop)
	fmt.Fprint(os.Stderr, "\r\033[K") // clear spinner line

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
		fmt.Fprintln(os.Stderr, "no review comments found")
		os.Exit(0)
	}

	// Build quickfix lines
	var lines []string
	for _, c := range comments {
		msg := strings.ReplaceAll(c.Message, "\n", " ")
		msg = strings.ReplaceAll(msg, "\r", " ")
		lines = append(lines, fmt.Sprintf("%s:%d:%d:%s", c.File, c.Line, c.Col, msg))
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
