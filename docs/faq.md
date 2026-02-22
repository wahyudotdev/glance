# Frequently Asked Questions

## General

### What is Glance?

Glance is a specialized MITM (Man-in-the-Middle) proxy designed for AI Agents and developers to intercept, inspect, and mock HTTP/HTTPS traffic. It provides real-time visibility into network activity and integrates with AI workflows via the Model Context Protocol (MCP).

### How is Glance different from other proxies?

Glance is specifically designed for:
- **AI Agent Integration**: Native MCP support for AI-assisted debugging
- **Developer Experience**: Modern dashboard with real-time updates
- **Ease of Use**: One-liner setup, automatic certificate handling
- **Platform Support**: Specialized support for Java, Android, and Chromium

### Is Glance open source?

Yes! Glance is licensed under the MIT License. You can view the source code on [GitHub](https://github.com/wahyudotdev/glance).

### What platforms does Glance support?

Glance runs on:
- macOS (ARM64 and AMD64)
- Linux (AMD64, ARM64, ARM)
- Windows (AMD64)

## Installation & Setup

### How do I install Glance?

The easiest way is via Homebrew:
```bash
brew tap wahyudotdev/tap
brew install glance
```

See the [Installation Guide](installation.md) for other methods.

### Do I need to install certificates?

For HTTPS interception, yes. Glance generates a local CA certificate that needs to be trusted by your system. The [Client Configuration Guide](clients.md#https-certificate-trust) has detailed instructions.

### Can I run Glance on a different port?

Yes! Use command-line flags:
```bash
glance --proxy-port 8080 --dashboard-port 8081 --mcp-port 8082
```

### How do I configure my app to use Glance?

Set proxy environment variables:
```bash
export HTTP_PROXY=http://localhost:15500
export HTTPS_PROXY=http://localhost:15500
```

Or use the one-liner:
```bash
eval "$(curl -s http://localhost:15501/setup)"
```

See the [Client Configuration Guide](clients.md) for more options.

## Features

### Can I modify requests/responses on the fly?

Yes! Glance supports:
- **Mocks**: Return static responses for specific URL patterns
- **Breakpoints**: Pause traffic for manual modification
- **Scenario Recording**: Replay and modify request sequences

See the [Mocking & Breakpoints Guide](features/mocking.md).

### Does Glance support WebSockets?

Currently, Glance focuses on HTTP/HTTPS traffic. WebSocket support is planned for a future release.

### Can I export captured traffic?

Yes! You can:
- Export individual requests as cURL commands
- Record scenarios and export as test scripts (via MCP)
- Generate OpenAPI documentation from traffic (via MCP)

### How long is traffic stored?

Traffic is stored in a local SQLite database (`~/.glance.db`). You can clear it at any time using the "Clear Traffic" button in the dashboard or via the MCP `clear_traffic` tool.

## MCP Integration

### What is MCP?

The Model Context Protocol is an open standard that allows AI agents to access external tools and resources. Glance implements MCP to let AI agents like Claude Desktop interact with network traffic.

See the [MCP Integration Guide](mcp/) for details.

### How do I set up Claude Desktop with Glance?

1. Edit your Claude Desktop config file
2. Add Glance as an MCP server
3. Restart Claude Desktop

See [Setting Up with Claude Desktop](mcp/#setting-up-with-claude-desktop) for step-by-step instructions.

### Can I use Glance with other AI agents?

Yes! Any MCP-compatible AI agent can connect to Glance's MCP server. Claude Desktop is just one example.

### What can AI agents do with Glance?

AI agents can:
- Inspect and analyze network traffic
- Create mocks and breakpoints
- Execute HTTP requests
- Generate API documentation
- Create test scripts from scenarios

See the [MCP Reference](mcp/reference.md) for all capabilities.

## Performance & Security

### Does Glance slow down my requests?

Minimal overhead. Glance uses efficient streaming and caching, so most requests have negligible latency increase (typically < 10ms).

### Is Glance secure?

Glance is designed for **local development and testing only**. The CA certificate is generated locally and stored on your machine. Never use Glance in production or on untrusted networks.

⚠️ **Security Warning**: MITM proxies can decrypt HTTPS traffic. Only use Glance on your own traffic in trusted environments.

### Where is data stored?

All captured traffic is stored locally in `~/.glance.db`. No data is sent to external servers.

### Can I use Glance in production?

**No.** Glance is designed for development and testing only. Using a MITM proxy in production poses security risks.

## Troubleshooting

### I'm getting certificate errors

You need to trust Glance's CA certificate. See the [HTTPS Certificate Trust](clients.md#https-certificate-trust) section.

### No traffic is appearing in the dashboard

Common causes:
1. Proxy environment variables not set
2. Application not respecting proxy settings
3. Another proxy/VPN interfering

See the [Troubleshooting Guide](troubleshooting.md).

### Glance won't start - port already in use

Another service is using the default ports. Start Glance on different ports:
```bash
glance --proxy-port 8080 --dashboard-port 8081
```

### Can I capture traffic from Docker containers?

Yes, but you need to configure Docker to use the proxy. See the [Docker section](troubleshooting.md#docker-containers) in the Troubleshooting Guide.

## Development

### How do I contribute to Glance?

We welcome contributions! See the [Development Guide](development.md) and [Contributing Guidelines](contributing.md).

### How do I build Glance from source?

```bash
git clone https://github.com/wahyudotdev/glance.git
cd glance
make build
```

See the [Development Guide](development.md) for details.

### Can I add custom features?

Absolutely! Glance is open source. Fork the repository, add your features, and submit a pull request.

## Licensing

### What license is Glance under?

MIT License. You're free to use, modify, and distribute Glance.

### Can I use Glance commercially?

Yes, the MIT License allows commercial use. However, Glance is intended for development/testing, not production use.

## Still Have Questions?

- 💬 [Start a discussion](https://github.com/wahyudotdev/glance/discussions)
- 🐛 [Report an issue](https://github.com/wahyudotdev/glance/issues)
- 📖 [Read the documentation](/)
