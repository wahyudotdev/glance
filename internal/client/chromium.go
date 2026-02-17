package client

import (
	"fmt"
	"os/exec"
	"runtime"
)

// LaunchChromium starts a Chromium-based browser configured to use the Glance.
func LaunchChromium(proxyAddr string) error {
	var cmd *exec.Cmd
	// Base flags for proxy and ignoring cert errors (for local MITM)
	args := []string{
		fmt.Sprintf("--proxy-server=http://%s", proxyAddr),
		"--ignore-certificate-errors",
		"--no-first-run",
		"--no-default-browser-check",
		"--user-data-dir=/tmp/glance-chrome", // Isolated session
		"https://www.google.com",             // Open a light site by default
	}

	switch runtime.GOOS {
	case "darwin":
		// On macOS, we try to find Chrome or Chromium
		path := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
		// #nosec G204
		cmd = exec.Command(path, args...)
	case "linux":
		// #nosec G204
		cmd = exec.Command("google-chrome", args...)
	case "windows":
		// #nosec G204
		cmd = exec.Command("chrome.exe", args...)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}
