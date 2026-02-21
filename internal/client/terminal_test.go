package client

import (
	"strings"
	"testing"
)

func TestGetTerminalSetupScript(t *testing.T) {
	script := GetTerminalSetupScript(":8000")
	if !strings.Contains(script, "HTTP_PROXY=http://:8000") {
		t.Errorf("Script missing HTTP_PROXY: %s", script)
	}
}
