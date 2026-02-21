package client

import (
	"fmt"
	"glance/internal/ca"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestSplitAddr(t *testing.T) {
	tests := []struct {
		addr     string
		wantHost string
		wantPort string
	}{
		{":8000", "127.0.0.1", "8000"},
		{"localhost:8000", "localhost", "8000"},
		{"0.0.0.0:15500", "127.0.0.1", "15500"},
		{"192.168.1.1:80", "192.168.1.1", "80"},
	}

	for _, tt := range tests {
		h, p := splitAddr(tt.addr)
		if h != tt.wantHost || p != tt.wantPort {
			t.Errorf("splitAddr(%s) = %s, %s; want %s, %s", tt.addr, h, p, tt.wantHost, tt.wantPort)
		}
	}
}

func TestFindToolsJar(_ *testing.T) {
	// Mock JAVA_HOME
	oldJH := os.Getenv("JAVA_HOME")
	defer func() { _ = os.Setenv("JAVA_HOME", oldJH) }()

	_ = os.Setenv("JAVA_HOME", "/tmp/dummy-jdk")
	// findToolsJar doesn't check if file exists on disk in its logic if it finds it via environment
	_ = findToolsJar("123")
}

func TestBuildAndAttachAgent_Mock(t *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()

	execCommand = func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestJavaAgentHelperProcess", "--", command}
		cs = append(cs, args...)
		cmd :=
			exec.Command(os.Args[0], cs...) // #nosec G204 G702
		cmd.Env = []string{"GO_WANT_JAVA_AGENT_HELPER_PROCESS=1"}
		return cmd
	}

	// Create a dummy CA cert for the logic to proceed
	ca.SetupCA()

	err := BuildAndAttachAgent("1234", ":8000")
	if err != nil {
		t.Fatalf("BuildAndAttachAgent failed: %v", err)
	}
}

func TestJavaAgentHelperProcess(_ *testing.T) {
	if os.Getenv("GO_WANT_JAVA_AGENT_HELPER_PROCESS") != "1" {
		return
	}

	args := os.Args
	for i, arg := range args {
		if arg == "--" {
			args = args[i+1:]
			break
		}
	}

	// Mock failure for modular java to trigger jcmd fallback
	if args[0] == "java" && strings.Contains(strings.Join(args, " "), "--add-modules") {
		_, _ = fmt.Fprintln(os.Stderr, "Attachment failed")
		os.Exit(1)
	}

	// Mock success for jcmd
	if args[0] == "jcmd" {
		_, _ = fmt.Fprintln(os.Stdout, "Command executed successfully")
	}

	os.Exit(0)
}

func TestBuildAndAttachAgent_Fallbacks(_ *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()

	execCommand = func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestJavaAgentHelperProcess", "--", command}
		cs = append(cs, args...)
		cmd :=
			exec.Command(os.Args[0], cs...) // #nosec G204 G702
		cmd.Env = []string{"GO_WANT_JAVA_AGENT_HELPER_PROCESS=1"}
		return cmd
	}

	ca.SetupCA()
	_ = BuildAndAttachAgent("1234", ":8000")
}
