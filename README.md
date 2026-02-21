# Glance

![CI](https://github.com/wahyudotdev/glance/actions/workflows/ci.yml/badge.svg)
![Coverage](https://raw.githubusercontent.com/wahyudotdev/glance/main/coverage.svg)

Glance is a specialized MITM (Man-in-the-Middle) proxy designed for **AI Agents** and developers to intercept, inspect, and mock HTTP/HTTPS traffic. It provides a real-time view of your application's network activity through a modern web dashboard and integrates deeply with AI workflows via the **Model Context Protocol (MCP)**.

## 🚀 Key Features

### 🔍 Traffic Inspection
- **Real-time Live View**: Stream captured requests and responses as they happen via WebSockets.
- **HTTPS Decryption**: Full MITM support with local CA generation and dynamic certificate management.
- **Detailed Metadata**: Inspect headers, query parameters, and formatted bodies (JSON/XML/HTML).
- **Export as cURL**: Quickly copy any request as a standard `curl` command for reproduction.

### 🛠️ Mocking & Breakpoints
- **Rule Engine**: Define rules to intercept and modify traffic on the fly.
- **Mocks**: Return static responses (status, headers, body) for specific URL patterns.
- **Breakpoints**: Pause matching requests or responses for manual modification in the dashboard before they proceed.
- **Scenario Recording**: Group related requests into sequences with one-click "Add to Scenario" or real-time recording mode.
- **Variable Mapping**: Define how data flows between requests (e.g., extracting a token from one response to use in the next header) to help AI understand dependencies.
- **CORS Support**: Mocked responses automatically handle preflight `OPTIONS` requests and include permissive headers.

### 🤖 AI Agent Integration (MCP)
- **Built-in MCP Server**: Native support for the Model Context Protocol.
- **AI-Powered Tools**: Allows AI agents (like Claude Desktop) to list traffic, execute requests, and manage rules.
- **Integration Docs**: Integrated documentation in the dashboard header with dynamic configuration snippets for Claude Desktop.

## 📋 MCP Specification

Glance implements the Model Context Protocol, allowing AI agents to interact with the proxy programmatically.

### Resources
- `proxy://status`: Returns the current proxy address, status, and active rules count.
- `traffic://latest`: Returns the most recent 10 HTTP requests in JSON format.

### Tools
- `inspect_network_traffic`: PRIMARY tool to list captured HTTP traffic summaries with optional keyword filtering and configurable limit (max limit follows system settings). MUST be used first for any network debugging.
- `inspect_request_details`: MANDATORY tool to retrieve full headers and body for a specific traffic entry by ID to diagnose root causes.
- `execute_request`: Execute or replay custom HTTP requests through the proxy.
- `add_mock_rule`: Create a mocking rule to return static responses.
- `add_breakpoint_rule`: Create a breakpoint rule to pause traffic for manual inspection.
- `list_rules`: List all active mocks and breakpoints.
- `delete_rule`: Remove an interception rule by ID.
- `list_scenarios`: List all recorded traffic scenarios.
- `get_scenario`: Get full sequence details and variable mappings for a scenario.
- `add_scenario`: Create a new scenario with basic metadata.
- `update_scenario`: Modify an existing scenario (name, description, steps, or variable mappings).
- `delete_scenario`: Remove a scenario by ID.
- `clear_traffic`: Reset/clear the captured traffic logs.
- `get_proxy_status`: Get real-time proxy address and status.

### Prompts
- `analyze-traffic`: Pre-defined prompt to have an AI analyze recent traffic for errors or anomalies.
- `generate-api-docs`: Pre-defined prompt to generate OpenAPI documentation from captured traffic.
- `generate-scenario-test`: Pre-defined prompt to transform a recorded scenario into an automated test script (e.g., Playwright).

### 📱 Client & Environment Support
- **Java/JVM**: Specialized overrides for `HttpsURLConnection` to support already-running apps (Java 8+).
- **Android**: ADB-based device discovery, automatic CA certificate installation, and proxy configuration.
- **Chromium**: Auto-launch browser instances with pre-configured proxy and certificate flags.
- **One-liner Setup**: `eval "$(curl -s http://localhost:15501/setup)"` for quick terminal proxy configuration.

## 📦 Installation

### Homebrew (macOS & Linux)
You can install Glance via our official tap:

```bash
brew tap wahyudotdev/tap
brew install glance
```

### Direct Download
Download the latest pre-compiled binary for your platform from the [Releases](https://github.com/wahyudotdev/glance/releases) page.

## ⚖️ Compatibility Matrix

| Environment | Support Level | Requirements | Features |
| :--- | :--- | :--- | :--- |
| **Java / JVM** | Native (Auto) | Java 8+ (v52.0+) | `HttpsURLConnection` overrides, Auto-trust CA |
| **Android** | ADB-based | Android 7+ (API 24+) | `adb reverse`, Auto-cert push, System/User trust guides |
| **Chromium** | Native (Launch) | Chrome / Edge / Brave | Auto-launch with proxy flags, Ignore cert errors |
| **Terminal** | One-liner | Bash / Zsh / Fish | `HTTP_PROXY`, `HTTPS_PROXY` injection |
| **MCP** | Built-in | Claude Desktop / MCP Host | Resource inspection, Tool execution, AI analysis |

## 📜 Changelog
See [CHANGELOG.md](CHANGELOG.md) for details on recent updates and milestones.

## 🏗️ Architecture

- **Backend**: Go (using Fiber and goproxy).
- **Frontend**: React, TypeScript, and Tailwind CSS.
- **Database**: Pure-Go SQLite (`~/.glance.db`) with Write-Ahead Logging (WAL) and Write-Behind caching for high performance.
- **Distribution**: Single binary containing the embedded dashboard assets.

## 🛠️ Development

### Prerequisites
- Go 1.24+
- Node.js & npm (for frontend development)

### Build Commands
The project uses a `Makefile` to manage builds:

```bash
# Build the binary for your current OS (includes frontend build)
make build

# Run all tests
make test

# Run linter
make lint

# Generate and view HTML coverage report
make test-coverage
```

### Running the Project
Once built, you can run the `glance` binary:

```bash
./glance
```
By default, the dashboard is available at `http://localhost:15501`, the proxy at `http://localhost:15500`, and the MCP server at `http://localhost:15502`.

## 📜 Attributions

Glance is built on top of amazing open-source projects:

- **[elazarl/goproxy](https://github.com/elazarl/goproxy)**: The core MITM proxy engine.
- **[gofiber/fiber](https://github.com/gofiber/fiber)**: High-performance web framework.
- **[Model Context Protocol](https://modelcontextprotocol.io)**: AI agent integration standard.
- **[Lucide Icons](https://lucide.dev)**: Beautiful & consistent iconography.
- **[Tailwind CSS](https://tailwindcss.com)**: Modern utility-first styling.
- **[SQLite](https://sqlite.org)**: Lightweight & high-performance persistence.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
