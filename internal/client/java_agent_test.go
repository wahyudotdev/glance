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

func TestBuildAndAttachAgent_Failures(t *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()

	t.Run("Javac Failure", func(t *testing.T) {
		execCommand = func(command string, args ...string) *exec.Cmd {
			cs := []string{"-test.run=TestJavaAgentHelperProcess", "--", command}
			cs = append(cs, args...)
			cmd := exec.Command(os.Args[0], cs...) // #nosec G204 G702
			cmd.Env = []string{"GO_WANT_JAVA_AGENT_HELPER_PROCESS=1", "FAIL_JAVAC=1"}
			return cmd
		}
		ca.SetupCA()
		err := BuildAndAttachAgent("1234", ":8000")
		if err == nil {
			t.Error("Expected error on javac failure")
		}
	})

	t.Run("Jar Failure", func(t *testing.T) {
		execCommand = func(command string, args ...string) *exec.Cmd {
			cs := []string{"-test.run=TestJavaAgentHelperProcess", "--", command}
			cs = append(cs, args...)
			cmd := exec.Command(os.Args[0], cs...) // #nosec G204 G702
			cmd.Env = []string{"GO_WANT_JAVA_AGENT_HELPER_PROCESS=1", "FAIL_JAR=1"}
			return cmd
		}
		ca.SetupCA()
		err := BuildAndAttachAgent("1234", ":8000")
		if err == nil {
			t.Error("Expected error on jar failure")
		}
	})
}

func TestFindToolsJar_Proc(t *testing.T) {
	oldRL := readlink
	oldST := stat
	defer func() {
		readlink = oldRL
		stat = oldST
	}()

	t.Run("Immediate lib", func(t *testing.T) {
		readlink = func(_ string) (string, error) { return "/jvm/bin/java", nil }
		stat = func(path string) (os.FileInfo, error) {
			if path == "/jvm/lib/tools.jar" {
				return nil, nil
			}
			return nil, os.ErrNotExist
		}
		if p := findToolsJar("1"); p != "/jvm/lib/tools.jar" {
			t.Errorf("Got %s", p)
		}
	})

	t.Run("One level up", func(t *testing.T) {
		readlink = func(_ string) (string, error) { return "/jvm/jre/bin/java", nil }
		stat = func(path string) (os.FileInfo, error) {
			if path == "/jvm/lib/tools.jar" {
				return nil, nil
			}
			return nil, os.ErrNotExist
		}
		if p := findToolsJar("1"); p != "/jvm/lib/tools.jar" {
			t.Errorf("Got %s", p)
		}
	})

	t.Run("Common paths", func(t *testing.T) {
		readlink = func(_ string) (string, error) { return "", os.ErrNotExist }
		stat = func(path string) (os.FileInfo, error) {
			if path == "/usr/lib/jvm/default-java/lib/tools.jar" {
				return nil, nil
			}
			return nil, os.ErrNotExist
		}
		if p := findToolsJar("1"); p != "/usr/lib/jvm/default-java/lib/tools.jar" {
			t.Errorf("Got %s", p)
		}
	})
}

func TestJavaAgentHelperProcess(_ *testing.T) {
	if os.Getenv("GO_WANT_JAVA_AGENT_HELPER_PROCESS") != "1" {
		return
	}

	if os.Getenv("FAIL_JAVAC") == "1" && os.Args[len(os.Args)-1] != "jps" {
		// Only fail if it's javac (we don't want to fail everything)
		for _, arg := range os.Args {
			if arg == "javac" {
				_, _ = fmt.Fprintln(os.Stderr, "Javac failed")
				os.Exit(1)
			}
		}
	}

	if os.Getenv("FAIL_JAR") == "1" {
		for _, arg := range os.Args {
			if arg == "jar" {
				_, _ = fmt.Fprintln(os.Stderr, "Jar failed")
				os.Exit(1)
			}
		}
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
