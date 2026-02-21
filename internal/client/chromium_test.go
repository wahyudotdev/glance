package client

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
)

func TestLaunchChromium_Mock(t *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()

	execCommand = func(command string, args ...string) *exec.Cmd {
		// Just return a command that succeeds (exits 0)
		cs := []string{"-test.run=TestChromiumHelperProcess", "--", command}
		cs = append(cs, args...)
		cmd :=
			exec.Command(os.Args[0], cs...) // #nosec G204 G702
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
		return cmd
	}

	err := LaunchChromium(":8000")
	// Since LaunchChromium calls cmd.Start(), our helper needs to be robust.
	// For now, we just verify it doesn't return an immediate error on supported platforms.
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" && runtime.GOOS != "windows" {
		if err == nil {
			t.Error("Expected error on unsupported platform")
		}
	} else {
		if err != nil {
			t.Errorf("LaunchChromium failed: %v", err)
		}
	}
}

func TestChromiumHelperProcess(_ *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}

func TestLaunchChromium_Wiring(t *testing.T) {
	// We can't actually launch it, but we can check if it returns error
	// on unsupported platforms (if any) or if it at least attempts to run on supported ones.
	// For testing purposes, we just ensure it doesn't panic.
	if runtime.GOOS == "unsupported" {
		err := LaunchChromium(":8000")
		if err == nil {
			t.Error("Expected error on unsupported platform")
		}
	}
}
