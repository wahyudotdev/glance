# Supported Clients

Glance supports various client environments and provides specialized configuration for each. This guide will help you configure your applications to use Glance as a proxy.

## Supported Environments

| Environment | Support Level | Configuration Method |
|-------------|---------------|---------------------|
| **Terminal/CLI** | ✅ Native | Environment variables or one-liner |
| **Java/JVM** | ✅ Native | Automatic HttpsURLConnection overrides |
| **Android** | ✅ ADB-based | Automatic CA installation and proxy setup |
| **Chromium** | Dashboard (One-click) | Auto-launch with proxy flags |
| **Python** | ✅ Native | Environment variables or library-specific config |
| **Node.js** | ✅ Native | Environment variables or axios/fetch config |
| **cURL** | ✅ Native | `-x` flag or environment variables |

## Quick Setup Methods

### Terminal One-Liner (Recommended)

The easiest way to configure your terminal:

```bash
eval "$(curl -s http://localhost:15501/setup)"
```

This automatically sets up the current session with proper proxy environment variables.

### Manual Environment Variables

Set these in your shell:

```bash
export HTTP_PROXY=http://localhost:15500
export HTTPS_PROXY=http://localhost:15500
export NO_PROXY=localhost,127.0.0.1
```

For persistence, add to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.).

## Platform-Specific Guides

For detailed configuration instructions for specific platforms:

- [Java/JVM Applications](clients/java.md) - Native support with automatic certificate handling
- [Android Devices](clients/android.md) - ADB-based setup with CA installation
- [Chromium Browsers](clients/chromium.md) - Auto-launch configuration
- [Terminal/CLI Tools](clients/terminal.md) - Advanced proxy configuration

## Testing Your Configuration

After configuring your client, verify it's working:

```bash
curl -x http://localhost:15500 https://api.github.com
```

You should see:
1. The response from GitHub
2. The request appear in the Glance dashboard at `http://localhost:15501`

## HTTPS Certificate Trust

For HTTPS interception, you may need to trust Glance's CA certificate:

### macOS

```bash
# Export the CA certificate
curl http://localhost:15501/ca.crt -o glance-ca.crt

# Add to keychain
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain glance-ca.crt
```

### Linux

```bash
# Export the CA certificate
curl http://localhost:15501/ca.crt -o glance-ca.crt

# Copy to trusted certificates
sudo cp glance-ca.crt /usr/local/share/ca-certificates/
sudo update-ca-certificates
```

### Windows

1. Download the CA certificate from `http://localhost:15501/ca.crt`
2. Double-click the certificate
3. Click "Install Certificate"
4. Select "Local Machine"
5. Choose "Place all certificates in the following store"
6. Select "Trusted Root Certification Authorities"
7. Complete the wizard

## Common Issues

### Certificate Errors

If you see SSL/TLS certificate errors:

1. Make sure you've trusted the Glance CA certificate
2. Some applications may need additional configuration to trust system certificates
3. See the [Troubleshooting Guide](troubleshooting.md#certificate-errors) for more help

### No Traffic Appearing

If you don't see traffic in the dashboard:

1. Verify proxy environment variables are set correctly
2. Check that your application respects proxy settings
3. Ensure no other proxy or VPN is interfering
4. Try the test command above to verify basic connectivity

### Application-Specific Issues

Some applications ignore system proxy settings. Check the specific guides:

- Java applications may need JVM arguments
- Some tools require explicit proxy configuration
- Docker containers need special network configuration

## Next Steps

- [MCP Integration](mcp/) - Connect with AI agents
- [Mocking & Breakpoints](features/mocking.md) - Modify traffic on the fly
- [Troubleshooting](troubleshooting.md) - Common issues and solutions
