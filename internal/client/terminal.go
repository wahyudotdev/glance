package client

import (
	"agent-proxy/internal/ca"
	"fmt"
	"strings"
)

func GetTerminalSetupScript(proxyAddr string) string {
	var sb strings.Builder

	// Standard proxy environment variables
	proxyURL := fmt.Sprintf("http://%s", proxyAddr)
	sb.WriteString(fmt.Sprintf("export HTTP_PROXY=%s\n", proxyURL))
	sb.WriteString(fmt.Sprintf("export HTTPS_PROXY=%s\n", proxyURL))
	sb.WriteString(fmt.Sprintf("export ALL_PROXY=%s\n", proxyURL))
	sb.WriteString("export NO_PROXY=localhost,127.0.0.1,::1\n")

	// CA Certificate environment variables
	if ca.CAPath != "" {
		sb.WriteString(fmt.Sprintf("export SSL_CERT_FILE=%s\n", ca.CAPath))
		sb.WriteString(fmt.Sprintf("export REQUESTS_CA_BUNDLE=%s\n", ca.CAPath))
		sb.WriteString(fmt.Sprintf("export NODE_EXTRA_CA_CERTS=%s\n", ca.CAPath))
	}

	// Helper alias to unset proxy
	sb.WriteString("\n# Run 'unproxy' to clear these variables\n")
	sb.WriteString("alias unproxy='unset HTTP_PROXY HTTPS_PROXY ALL_PROXY NO_PROXY SSL_CERT_FILE REQUESTS_CA_BUNDLE NODE_EXTRA_CA_CERTS'\n")

	return sb.String()
}
