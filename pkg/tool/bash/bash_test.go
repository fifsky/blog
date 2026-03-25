package bash

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCommandSafe(t *testing.T) {
	tests := []struct {
		name string
		cmd  string
		safe bool
	}{
		{"safe echo", "echo hello", true},
		{"safe ls", "ls -la /", true},
		{"safe rm file", "rm -f /tmp/test.txt", true},
		
		{"unsafe rm rf root", "rm -rf /", false},
		{"unsafe rm fr root", "rm -fr /", false},
		{"unsafe rm r f root", "rm -r -f /", false},
		{"unsafe rm rf root with spaces", "rm   -rf   /", false},
		{"unsafe rm rf root uppercase", "RM -RF /", false},
		
		{"unsafe mkfs", "mkfs.ext4 /dev/sda1", false},
		{"unsafe reboot", "sudo reboot", false},
		{"unsafe shutdown", "shutdown -h now", false},
		{"unsafe dd", "dd if=/dev/zero of=/dev/sda", false},
		{"unsafe chmod", "chmod -R 777 /", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.safe, isCommandSafe(tt.cmd))
		})
	}
}

func TestBashTool_Handle(t *testing.T) {
	tool := New()
	ctx := context.Background()

	// Test valid command
	args, _ := json.Marshal(map[string]any{"command": "echo 'hello world'"})
	res, err := tool.Handle(ctx, string(args))
	assert.NoError(t, err)
	
	var resObj map[string]any
	err = json.Unmarshal([]byte(res), &resObj)
	assert.NoError(t, err)
	assert.Equal(t, "success", resObj["status"])
	assert.Equal(t, float64(0), resObj["exit_code"])
	assert.True(t, strings.Contains(resObj["stdout"].(string), "hello world"))

	// Test dangerous command
	args, _ = json.Marshal(map[string]any{"command": "rm -rf /"})
	res, err = tool.Handle(ctx, string(args))
	assert.NoError(t, err)
	
	err = json.Unmarshal([]byte(res), &resObj)
	assert.NoError(t, err)
	assert.Equal(t, "error", resObj["status"])
	assert.True(t, strings.Contains(resObj["error"].(string), "dangerous"))

	// Test timeout
	args, _ = json.Marshal(map[string]any{
		"command": "sleep 2",
		"timeout_seconds": 1,
	})
	res, err = tool.Handle(ctx, string(args))
	assert.NoError(t, err)
	
	err = json.Unmarshal([]byte(res), &resObj)
	assert.NoError(t, err)
	assert.Equal(t, "timeout", resObj["status"])
	assert.Equal(t, float64(-1), resObj["exit_code"])
}
