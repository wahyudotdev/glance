package client

import (
	"fmt"
	"os/exec"
	"runtime"
)

func LaunchChromium(proxyAddr string) error {
	var cmd *exec.Cmd
	// Base flags for proxy and ignoring cert errors (for local MITM)
	args := []string{
		fmt.Sprintf("--proxy-server=http://%s", proxyAddr),
		"--ignore-certificate-errors",
		"--no-first-run",
		"--no-default-browser-check",
		"--user-data-dir=/tmp/agent-proxy-chrome", // Isolated session
		"http://localhost:8081",                   // Open dashboard automatically
	}

	switch runtime.GOOS {
	case "darwin":
		// On macOS, we try to find Chrome or Chromium
		path := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
		cmd = exec.Command(path, args...)
	case "linux":
		cmd = exec.Command("google-chrome", args...)
	case "windows":
		cmd = exec.Command("chrome.exe", args...)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}
