package client

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestListJavaProcesses_Mock(t *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()

	execCommand = func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestJavaHelperProcess", "--", command}
		cs = append(cs, args...)
		cmd :=
			exec.Command(os.Args[0], cs...) // #nosec G204 G702
		cmd.Env = []string{"GO_WANT_JAVA_HELPER_PROCESS=1"}
		return cmd
	}

	procs, err := ListJavaProcesses()
	if err != nil {
		t.Fatalf("ListJavaProcesses failed: %v", err)
	}
	// Verify that our mock output was parsed correctly
	found := false
	for _, p := range procs {
		if p.PID == "1234" && p.Name == "test.jar" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Mock process not found in results: %+v", procs)
	}
}

func TestListJavaProcesses_Error(_ *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()

	execCommand = func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestJavaHelperProcess", "--", command}
		cs = append(cs, args...)
		cmd :=
			exec.Command(os.Args[0], cs...) // #nosec G204 G702
		cmd.Env = []string{"GO_WANT_JAVA_HELPER_PROCESS=1", "FAIL_JAVA=1"}
		return cmd
	}

	_, _ = ListJavaProcesses()
	// It shouldn't return error because it has multiple fallbacks and just logs them.
	// But it will exercise the error paths.
}

func TestListJavaProcesses_Parsing(t *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()

	execCommand = func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestJavaHelperProcess", "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...) // #nosec G204 G702
		cmd.Env = []string{"GO_WANT_JAVA_HELPER_PROCESS=1", "BAD_JAVA_OUTPUT=1"}
		return cmd
	}

	procs, err := ListJavaProcesses()
	if err != nil {
		t.Fatalf("ListJavaProcesses failed: %v", err)
	}
	// Verify it handled malformed lines
	// Lines: "1234" (skip), "not-a-number name" (valid), "5678 valid.jar" (valid)
	if len(procs) != 2 {
		t.Errorf("Expected 2 valid processes, got %d", len(procs))
	}
}

func TestJavaHelperProcess(_ *testing.T) {
	if os.Getenv("GO_WANT_JAVA_HELPER_PROCESS") != "1" {
		return
	}

	if os.Getenv("FAIL_JAVA") == "1" {
		os.Exit(1)
	}

	if os.Getenv("BAD_JAVA_OUTPUT") == "1" {
		_, _ = fmt.Fprintln(os.Stdout, "1234") // missing name
		_, _ = fmt.Fprintln(os.Stdout, "not-a-number name")
		_, _ = fmt.Fprintln(os.Stdout, "5678 valid.jar")
		os.Exit(0)
	}

	args := os.Args
	for i, arg := range args {
		if arg == "--" {
			args = args[i+1:]
			break
		}
	}

	switch args[0] {
	case "jps":
		_, _ = fmt.Fprintln(os.Stdout, "1234 test.jar")
	case "ps":
		// Mock ps aux output if jps fails, but our test currently exercises jps path first
		_, _ = fmt.Fprintln(os.Stdout, "user 5678 0.0 0.0 123 456 ? S 12:00 0:00 java -jar other.jar")
	}

	os.Exit(0)
}
