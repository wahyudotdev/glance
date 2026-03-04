package client

import (
	"glance/internal/ca"
	"strings"
	"testing"
)

func TestGetTerminalSetupScript(t *testing.T) {
	// Set a mock CA path for testing
	ca.CAPath = "/tmp/glance-ca.crt"

	t.Run("Default NoProxy", func(t *testing.T) {
		// Default should now NOT have NO_PROXY
		script := GetTerminalSetupScript("localhost:8000", "")
		if !strings.Contains(script, "unset NO_PROXY") {
			t.Errorf("Script should unset NO_PROXY by default to intercept everything: %s", script)
		}
	})

	t.Run("Custom NoProxy", func(t *testing.T) {
		script := GetTerminalSetupScript("localhost:8000", "some-host.com")
		if !strings.Contains(script, "NO_PROXY=some-host.com") {
			t.Errorf("Script missing custom NO_PROXY: %s", script)
		}
	})

	t.Run("Certificate Variables", func(t *testing.T) {
		script := GetTerminalSetupScript("localhost:8000", "")
		expectedVars := []string{
			"SSL_CERT_FILE",
			"REQUESTS_CA_BUNDLE",
			"CURL_CA_BUNDLE",
			"GIT_SSL_CAINFO",
			"PIP_CERT",
			"NODE_EXTRA_CA_CERTS",
			"npm_config_cafile",
		}

		for _, v := range expectedVars {
			if !strings.Contains(script, "export "+v+"=") {
				t.Errorf("Script missing expected export %s: %s", v, script)
			}
		}
	})
}
