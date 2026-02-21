package client

import (
	"fmt"
	"glance/internal/ca"
	"strings"
)

// GetTerminalSetupScript generates a shell script to configure proxy environment variables.
func GetTerminalSetupScript(proxyAddr string) string {
	var sb strings.Builder

	// Standard proxy environment variables
	proxyURL := fmt.Sprintf("http://%s", proxyAddr)
	fmt.Fprintf(&sb, "export HTTP_PROXY=%s\n", proxyURL)
	fmt.Fprintf(&sb, "export HTTPS_PROXY=%s\n", proxyURL)
	fmt.Fprintf(&sb, "export ALL_PROXY=%s\n", proxyURL)
	sb.WriteString("export NO_PROXY=localhost,127.0.0.1,::1\n")

	// CA Certificate environment variables
	if ca.CAPath != "" {
		fmt.Fprintf(&sb, "export SSL_CERT_FILE=%s\n", ca.CAPath)
		fmt.Fprintf(&sb, "export REQUESTS_CA_BUNDLE=%s\n", ca.CAPath)
		fmt.Fprintf(&sb, "export NODE_EXTRA_CA_CERTS=%s\n", ca.CAPath)
	}

	// Helper alias to unset proxy
	sb.WriteString("\n# Run 'unproxy' to clear these variables\n")
	sb.WriteString("alias unproxy='unset HTTP_PROXY HTTPS_PROXY ALL_PROXY NO_PROXY SSL_CERT_FILE REQUESTS_CA_BUNDLE NODE_EXTRA_CA_CERTS'\n")

	return sb.String()
}
