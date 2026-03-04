package client

import (
	"fmt"
	"glance/internal/ca"
	"strings"
)

// GetTerminalSetupScript generates a shell script to configure proxy environment variables.
func GetTerminalSetupScript(proxyAddr string, noProxy string) string {
	var sb strings.Builder

	// Colors for output
	sb.WriteString("BLUE='\\033[0;34m'\n")
	sb.WriteString("GREEN='\\033[0;32m'\n")
	sb.WriteString("NC='\\033[0m'\n")

	// Handle noProxy state:
	// If it's "__EMPTY__" or genuinely empty, we unset the variable to intercept everything.
	if noProxy == "__EMPTY__" {
		noProxy = ""
	}

	proxyURL := fmt.Sprintf("http://%s", proxyAddr)

	// Standard proxy environment variables (both upper and lower case for compatibility)
	fmt.Fprintf(&sb, "export HTTP_PROXY=%s\n", proxyURL)
	fmt.Fprintf(&sb, "export http_proxy=%s\n", proxyURL)
	fmt.Fprintf(&sb, "export HTTPS_PROXY=%s\n", proxyURL)
	fmt.Fprintf(&sb, "export https_proxy=%s\n", proxyURL)
	fmt.Fprintf(&sb, "export ALL_PROXY=%s\n", proxyURL)
	fmt.Fprintf(&sb, "export all_proxy=%s\n", proxyURL)

	if noProxy != "" {
		fmt.Fprintf(&sb, "export NO_PROXY=%s\n", noProxy)
		fmt.Fprintf(&sb, "export no_proxy=%s\n", noProxy)
	} else {
		sb.WriteString("unset NO_PROXY\n")
		sb.WriteString("unset no_proxy\n")
	}

	// CA Certificate environment variables
	if ca.CAPath != "" {
		fmt.Fprintf(&sb, "export SSL_CERT_FILE=%s\n", ca.CAPath)
		fmt.Fprintf(&sb, "export REQUESTS_CA_BUNDLE=%s\n", ca.CAPath)
		fmt.Fprintf(&sb, "export CURL_CA_BUNDLE=%s\n", ca.CAPath)
		fmt.Fprintf(&sb, "export GIT_SSL_CAINFO=%s\n", ca.CAPath)
		fmt.Fprintf(&sb, "export PIP_CERT=%s\n", ca.CAPath)
		fmt.Fprintf(&sb, "export NODE_EXTRA_CA_CERTS=%s\n", ca.CAPath)
		fmt.Fprintf(&sb, "export npm_config_cafile=%s\n", ca.CAPath)
	}

	// Helper alias to unset proxy
	sb.WriteString("\n# Run 'unproxy' to clear these variables\n")
	sb.WriteString("alias unproxy='unset HTTP_PROXY http_proxy HTTPS_PROXY https_proxy ALL_PROXY all_proxy NO_PROXY no_proxy SSL_CERT_FILE REQUESTS_CA_BUNDLE CURL_CA_BUNDLE GIT_SSL_CAINFO PIP_CERT NODE_EXTRA_CA_CERTS npm_config_cafile && unalias unproxy'\n")

	// Feedback
	sb.WriteString("\necho -e \"${GREEN}✓ Terminal configured for Glance interception${NC}\"\n")
	fmt.Fprintf(&sb, "echo -e \"${BLUE}Proxy:${NC} %s\"\n", proxyURL)
	if noProxy == "" {
		sb.WriteString("echo -e \"${BLUE}NO_PROXY:${NC} <none> (Intercepting EVERYTHING including localhost)\"\n")
	} else {
		fmt.Fprintf(&sb, "echo -e \"${BLUE}NO_PROXY:${NC} %s\"\n", noProxy)
	}
	if ca.CAPath != "" {
		fmt.Fprintf(&sb, "echo -e \"${BLUE}CA Cert:${NC} %s\"\n", ca.CAPath)
	}
	sb.WriteString("echo \"Type 'unproxy' to restore original terminal settings.\"\n")

	return sb.String()
}
