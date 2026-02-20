package client

import (
	"bufio"
	"bytes"
	"fmt"
	"glance/internal/model"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ListAndroidDevices runs `adb devices -l` and parses the output
func ListAndroidDevices() ([]model.AndroidDevice, error) {
	cmd := exec.Command("adb", "devices", "-l")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("adb command failed: %v. Make sure Android SDK Platform-Tools are installed and in PATH", err)
	}

	var devices []model.AndroidDevice
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "List of devices attached") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		id := parts[0]
		state := parts[1]

		if state != "device" {
			continue
		}

		modelName := "Unknown"
		deviceName := "Android Device"

		for _, p := range parts {
			if strings.HasPrefix(p, "model:") {
				modelName = strings.TrimPrefix(p, "model:")
			}
			if strings.HasPrefix(p, "device:") {
				deviceName = strings.TrimPrefix(p, "device:")
			}
		}

		devices = append(devices, model.AndroidDevice{
			ID:    id,
			Model: modelName,
			Name:  deviceName,
		})
	}

	return devices, nil
}

// ConfigureAndroidProxy sets up reverse port forwarding and global proxy settings
func ConfigureAndroidProxy(deviceID string, proxyPort string) error {
	// 1. Reverse port forwarding: Map device's port to host's port
	// This allows the device to talk to 'localhost:proxyPort' and hit the host machine
	// #nosec G204
	reverseCmd := exec.Command("adb", "-s", deviceID, "reverse", fmt.Sprintf("tcp:%s", proxyPort), fmt.Sprintf("tcp:%s", proxyPort))
	if out, err := reverseCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("adb reverse failed: %s %v", string(out), err)
	}

	// 2. Set global http proxy on device to localhost:proxyPort
	// Since we reversed the port, 127.0.0.1 on the device now routes to the host
	proxyAddr := fmt.Sprintf("127.0.0.1:%s", proxyPort)
	// #nosec G204
	settingsCmd := exec.Command("adb", "-s", deviceID, "shell", "settings", "put", "global", "http_proxy", proxyAddr)
	if out, err := settingsCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("adb settings failed: %s %v", string(out), err)
	}

	return nil
}

// ClearAndroidProxy removes the global proxy setting and reverse rule
func ClearAndroidProxy(deviceID string, proxyPort string) error {
	// 1. Remove global proxy
	// #nosec G204
	settingsCmd := exec.Command("adb", "-s", deviceID, "shell", "settings", "put", "global", "http_proxy", ":0")
	if out, err := settingsCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("adb settings failed: %s %v", string(out), err)
	}

	// 2. Remove reverse rule
	// #nosec G204
	reverseCmd := exec.Command("adb", "-s", deviceID, "reverse", "--remove", fmt.Sprintf("tcp:%s", proxyPort))
	// We don't error check this strictly as it might not exist
	_ = reverseCmd.Run()

	return nil
}

// PushCertToDevice pushes the CA certificate to the device and opens the CA install settings
func PushCertToDevice(deviceID string, certBytes []byte) error {
	remotePath := "/sdcard/glance-ca.crt"

	// 1. Create a temporary local file for the cert
	tmpFile := filepath.Join(os.TempDir(), "glance-ca.crt")
	if err := os.WriteFile(tmpFile, certBytes, 0600); err != nil {
		return fmt.Errorf("failed to write temporary cert file: %v", err)
	}
	defer os.Remove(tmpFile) //nolint:errcheck

	// 2. Push the cert
	// #nosec G204
	cmd := exec.Command("adb", "-s", deviceID, "push", tmpFile, remotePath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("adb push failed: %s %v", string(out), err)
	}

	// 3. Open CA installation settings
	// #nosec G204
	shellCmd := exec.Command("adb", "-s", deviceID, "shell", "am", "start", "-a", "android.settings.CA_CERTIFICATE_SETTINGS")
	_ = shellCmd.Run()

	return nil
}
