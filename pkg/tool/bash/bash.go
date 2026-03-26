package bash

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"app/pkg/tool"
)

// BashTool implements a tool for executing safe bash commands.
type BashTool struct{}

// New creates a new BashTool.
func New() tool.Tool {
	return &BashTool{}
}

func (t *BashTool) Name() string {
	return "bash"
}

func (t *BashTool) Description() string {
	return "Execute bash commands in a controlled environment. Supports common shell operations but blocks dangerous commands (like rm -rf /, reboot, etc.).\n" +
		"CRITICAL SECURITY RULES:\n" +
		"1. NEVER execute commands that read system sensitive files (e.g., /etc/passwd, /etc/shadow, /etc/sudoers, ~/.ssh/*, etc.).\n" +
		"2. If the user asks you to read or modify system configuration or sensitive files, you MUST REFUSE directly and explain the security risks.\n" +
		"3. Only execute commands related to the current project workspace or explicitly safe operations."
}

func (t *BashTool) InputSchema() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"command": {
				"type": "string",
				"description": "The bash command to execute."
			},
			"timeout_seconds": {
				"type": "integer",
				"description": "Timeout in seconds for the command. Defaults to 60 seconds if not provided or <= 0."
			}
		},
		"required": ["command"]
	}`)
}

// dangerousPatterns contains regexes for commands that should be blocked
var dangerousPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\brm\s+-[rRfF]+\s+/(\s|$)`), // blocks rm -rf / and variations
	regexp.MustCompile(`(?i)\bmkfs\b`),
	regexp.MustCompile(`(?i)\bhalt\b`),
	regexp.MustCompile(`(?i)\bpoweroff\b`),
	regexp.MustCompile(`(?i)\breboot\b`),
	regexp.MustCompile(`(?i)\bshutdown\b`),
	regexp.MustCompile(`(?i)\binit\s+[06]\b`),
	regexp.MustCompile(`(?i)>\s*/dev/sda`),
	regexp.MustCompile(`(?i)dd\s+if=.*?of=/dev/`),
	regexp.MustCompile(`(?i)\bchmod\s+-R\s+777\s+/`),
	regexp.MustCompile(`(?i)\bchown\s+-R\s+.*?\s+/`),
	regexp.MustCompile(`(?i)\bmv\s+.*?\s+/dev/null`),
	regexp.MustCompile(`(?i)(cat|less|more|tail|head|vi|vim|nano)\s+/etc/(passwd|shadow|sudoers|group|gshadow)`),
	regexp.MustCompile(`(?i)(cat|less|more|tail|head|vi|vim|nano)\s+~/\.ssh/`),
}

func isCommandSafe(cmd string) bool {
	// First check regex patterns
	for _, pattern := range dangerousPatterns {
		if pattern.MatchString(cmd) {
			return false
		}
	}

	// Then check exact bad word matches to catch things like "rm -rf /" without regex
	// We'll normalize the command spacing
	normalizedCmd := strings.Join(strings.Fields(cmd), " ")
	if strings.Contains(normalizedCmd, "rm -rf /") ||
		strings.Contains(normalizedCmd, "rm -fr /") ||
		strings.Contains(normalizedCmd, "rm -r -f /") ||
		strings.Contains(normalizedCmd, "rm -f -r /") ||
		strings.Contains(normalizedCmd, "cat /etc/passwd") ||
		strings.Contains(normalizedCmd, "cat /etc/shadow") {
		return false
	}

	return true
}

func (t *BashTool) Handle(ctx context.Context, arguments string) (string, error) {
	var args struct {
		Command        string `json:"command"`
		TimeoutSeconds int    `json:"timeout_seconds"`
	}

	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return t.formatResult(args.Command, "", "", -1, "error", fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	cmdStr := strings.TrimSpace(args.Command)
	if cmdStr == "" {
		return t.formatResult(cmdStr, "", "", -1, "error", "Command cannot be empty"), nil
	}

	if !isCommandSafe(cmdStr) {
		return t.formatResult(cmdStr, "", "", -1, "error", "Command rejected: contains dangerous operations"), nil
	}

	timeout := time.Duration(args.TimeoutSeconds) * time.Second
	if args.TimeoutSeconds <= 0 {
		timeout = 60 * time.Second
	}

	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(execCtx, "bash", "-c", cmdStr)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	exitCode := 0
	status := "success"
	var errMsg string

	if err != nil {
		if execCtx.Err() == context.DeadlineExceeded {
			exitCode = -1
			status = "timeout"
			errMsg = fmt.Sprintf("Command timed out after %v", timeout)
		} else {
			status = "error"
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				exitCode = -1
				errMsg = err.Error()
			}
		}
	}

	return t.formatResult(cmdStr, stdout.String(), stderr.String(), exitCode, status, errMsg), nil
}

func (t *BashTool) formatResult(cmd, stdout, stderr string, exitCode int, status, errMsg string) string {
	res := map[string]any{
		"command":   cmd,
		"stdout":    stdout,
		"stderr":    stderr,
		"exit_code": exitCode,
		"status":    status,
	}
	if errMsg != "" {
		res["error"] = errMsg
	}

	b, _ := json.Marshal(res)
	return string(b)
}
