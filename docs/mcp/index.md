# MCP Integration

Glance implements the **Model Context Protocol (MCP)**, allowing AI agents to interact with network traffic programmatically. This enables powerful AI-assisted debugging and testing workflows.

## What is MCP?

The [Model Context Protocol](https://modelcontextprotocol.io) is an open standard that allows AI agents to access external tools, resources, and context. Glance's MCP server exposes network traffic data and proxy controls to AI agents like Claude Desktop.

<div style="margin: 2rem 0; aspect-ratio: 16/9; width: 100%;">
  <iframe 
    style="width: 100%; height: 100%; border-radius: 12px;"
    src="https://www.youtube.com/embed/vgI1V-KKXhw" 
    title="Glance MCP Integration" 
    frameborder="0" 
    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" 
    allowfullscreen>
  </iframe>
</div>

![Glance MCP Settings](/_media/settings.png)

## Quick Setup (Claude Code CLI)

Run this one-liner in your terminal to instantly connect Glance to Claude Code:

```bash
claude mcp add --transport http glance http://localhost:15502/mcp
```

## Setting Up with Claude Desktop

### 1. Locate Your Claude Desktop Config

The MCP configuration file is located at:

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

### 2. Add Glance to MCP Servers

Edit the config file and add Glance as an HTTP/SSE server:

```json
{
  "mcpServers": {
    "glance": {
      "type": "http",
      "url": "http://localhost:15502/mcp"
    }
  }
}
```

### 3. Restart Claude Desktop

Restart Claude Desktop to load the new MCP configuration. You should now see Glance listed as an available MCP server.

## Setting Up with VS Code (Cline / Roo Code)

1. Open **MCP Settings** in the extension.
2. Add a new server with type `SSE`.
3. Paste the URL: `http://localhost:15502/mcp`

## Using Glance with Claude

Once configured, you can ask Claude to interact with your network traffic:

### Example Prompts

**Inspect recent traffic:**
```
Show me the last 10 HTTP requests captured by Glance
```

**Debug API errors:**
```
Analyze the recent traffic for any failed API calls and help me understand what went wrong
```

**Generate API documentation:**
```
Generate OpenAPI documentation from the captured traffic to the /api/users endpoint
```

**Create mocks:**
```
Create a mock for the /api/login endpoint that returns a success response
```

**Execute requests:**
```
Make a GET request to https://api.github.com/users through the Glance proxy
```

## MCP Capabilities

### Resources

Resources provide read-only access to proxy state and traffic data:

| Resource | Description |
|----------|-------------|
| `proxy://status` | Current proxy address, status, and active rules count |
| `traffic://latest` | Most recent 10 HTTP requests in JSON format |

### Tools

Tools allow AI agents to perform actions:

| Tool | Description |
|------|-------------|
| `inspect_network_traffic` | List captured HTTP traffic with filtering and configurable limit |
| `inspect_request_details` | Get full headers and body for a specific request |
| `execute_request` | Execute custom HTTP requests through the proxy |
| `add_mock_rule` | Create a mocking rule to return static responses |
| `add_breakpoint_rule` | Create a breakpoint rule to pause traffic |
| `list_rules` | List all active mocks and breakpoints |
| `delete_rule` | Remove an interception rule by ID |
| `list_scenarios` | List all recorded traffic scenarios |
| `get_scenario` | Get full sequence details and variable mappings |
| `add_scenario` | Create a new scenario with basic metadata |
| `update_scenario` | Modify an existing scenario |
| `delete_scenario` | Remove a scenario by ID |
| `clear_traffic` | Reset/clear captured traffic logs |
| `get_proxy_status` | Get real-time proxy address and status |

### Prompts

Pre-defined prompts for common workflows:

| Prompt | Description |
|--------|-------------|
| `analyze-traffic` | Analyze recent traffic for errors or anomalies |
| `generate-api-docs` | Generate OpenAPI documentation from captured traffic |
| `generate-scenario-test` | Transform a recorded scenario into an automated test script |

## Advanced Usage

### Custom MCP Configuration

You can customize Glance's MCP server behavior with additional arguments:

```json
{
  "mcpServers": {
    "glance": {
      "command": "/path/to/glance",
      "args": [
        "--mcp",
        "--mcp-port", "15502",
        "--log-level", "debug"
      ]
    }
  }
}
```

### Using with Other MCP Clients

Glance's MCP server can be used with any MCP-compatible client, not just Claude Desktop. The server runs on `http://localhost:15502` by default.

## Example Workflows

### 1. Debug a Failing API Call

1. Run your application with Glance proxy configured
2. Reproduce the failing API call
3. Ask Claude: *"Show me the last request to /api/endpoint and help me debug why it's failing"*
4. Claude will use `inspect_network_traffic` and `inspect_request_details` to analyze the issue

### 2. Generate Test Code from Scenarios

1. Record a scenario in the Glance dashboard
2. Ask Claude: *"Generate a Playwright test script from the 'user-login' scenario"*
3. Claude will use the `generate-scenario-test` prompt to create test code

### 3. Create Mocks for Testing

1. Capture the real API response in Glance
2. Ask Claude: *"Create a mock for this endpoint that returns a similar response but with test data"*
3. Claude will use `add_mock_rule` to create the mock

## Troubleshooting

### Claude Can't See Glance

1. Verify the path to the Glance binary is correct in your config
2. Check that you restarted Claude Desktop after editing the config
3. Look for error messages in Claude Desktop's MCP logs

### MCP Server Won't Start

1. Check if another service is using port 15502
2. Try specifying a different port with `--mcp-port`
3. Check Glance logs for error messages

### Tools Not Working

1. Make sure Glance is running and the proxy is active
2. Verify that traffic is being captured in the dashboard
3. Try using the tools directly in Claude to see specific error messages

## Learn More

- [MCP Tools & Resources Reference](reference.md) - Detailed API documentation
- [Model Context Protocol Documentation](https://modelcontextprotocol.io) - Official MCP docs

