package client

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

// Fake exec.Command for testing
func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd :=
		exec.Command(os.Args[0], cs...) // #nosec G204 G702
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestListAndroidDevices_Parsing(t *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()

	execCommand = func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...) // #nosec G204 G702
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1", "BAD_ADB_OUTPUT=1"}
		return cmd
	}

	devices, err := ListAndroidDevices()
	if err != nil {
		t.Fatalf("ListAndroidDevices failed: %v", err)
	}
	// Verify it skipped bad lines
	if len(devices) != 1 {
		t.Errorf("Expected 1 valid device, got %d", len(devices))
	}
}

func TestHelperProcess(_ *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	if os.Getenv("FAIL_ADB") == "1" {
		os.Exit(1)
	}

	if os.Getenv("BAD_ADB_OUTPUT") == "1" {
		_, _ = fmt.Fprintln(os.Stdout, "List of devices attached")
		_, _ = fmt.Fprintln(os.Stdout, "invalid-line")
		_, _ = fmt.Fprintln(os.Stdout, "offline-device offline")
		_, _ = fmt.Fprintln(os.Stdout, "dev1 device model:M1 device:D1")
		os.Exit(0)
	}

	args := os.Args
	for i, arg := range args {
		if arg == "--" {
			args = args[i+1:]
			break
		}
	}

	if len(args) > 1 && args[0] == "adb" && args[1] == "devices" {
		_, _ = fmt.Fprintln(os.Stdout, "List of devices attached")
		_, _ = fmt.Fprintln(os.Stdout, "emulator-5554          device product:sdk_gphone64_arm64 model:sdk_gphone64_arm64 device:emulator64_arm64 transport_id:1")
	}

	os.Exit(0)
}

func TestListAndroidDevices_Mock(t *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()
	execCommand = fakeExecCommand

	devices, err := ListAndroidDevices()
	if err != nil {
		t.Fatalf("ListAndroidDevices failed: %v", err)
	}
	if len(devices) != 1 || devices[0].ID != "emulator-5554" {
		t.Errorf("Expected 1 device, got %+v", devices)
	}
}

func TestListAndroidDevices_Error(t *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()

	execCommand = func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", command}
		cs = append(cs, args...)
		cmd :=
			exec.Command(os.Args[0], cs...) // #nosec G204 G702
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1", "FAIL_ADB=1"}
		return cmd
	}

	_, err := ListAndroidDevices()
	if err == nil {
		t.Error("Expected error")
	}
}

func TestConfigureAndroidProxy_Mock(t *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()
	execCommand = fakeExecCommand

	err := ConfigureAndroidProxy("dev1", "8000")
	if err != nil {
		t.Errorf("ConfigureAndroidProxy failed: %v", err)
	}
}

func TestClearAndroidProxy_Mock(t *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()
	execCommand = fakeExecCommand

	err := ClearAndroidProxy("dev1", "8000")
	if err != nil {
		t.Errorf("ClearAndroidProxy failed: %v", err)
	}
}

func TestPushCertToDevice_Mock(t *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()
	execCommand = fakeExecCommand

	err := PushCertToDevice("dev1", []byte("cert"))
	if err != nil {
		t.Errorf("PushCertToDevice failed: %v", err)
	}
}

func TestConfigureAndroidProxy_Error(t *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()
	execCommand = func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...) // #nosec G204 G702
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1", "FAIL_ADB=1"}
		return cmd
	}

	err := ConfigureAndroidProxy("dev1", "8000")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestClearAndroidProxy_Error(t *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()
	execCommand = func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...) // #nosec G204 G702
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1", "FAIL_ADB=1"}
		return cmd
	}

	err := ClearAndroidProxy("dev1", "8000")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestPushCertToDevice_Error(t *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()
	execCommand = func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...) // #nosec G204 G702
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1", "FAIL_ADB=1"}
		return cmd
	}

	err := PushCertToDevice("dev1", []byte("cert"))
	if err == nil {
		t.Error("Expected error")
	}
}
