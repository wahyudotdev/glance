# MCP Tools & Resources Reference

Complete reference for Glance's Model Context Protocol implementation.

## Resources

Resources provide read-only access to proxy state and traffic data.

### proxy://status

Get current proxy status and configuration.

**Response:**

```json
{
  "proxyAddress": "http://localhost:15500",
  "dashboardAddress": "http://localhost:15501",
  "mcpAddress": "http://localhost:15502",
  "status": "running",
  "activeRules": 5,
  "totalTraffic": 1234
}
```

### traffic://latest

Get the most recent 10 HTTP requests.

**Response:**

```json
{
  "traffic": [
    {
      "id": "uuid",
      "method": "GET",
      "url": "https://api.example.com/users",
      "status": 200,
      "timestamp": "2026-02-22T10:30:00Z"
    }
  ]
}
```

## Tools

Tools allow AI agents to perform actions.

### inspect_network_traffic

**PRIMARY** tool to list captured HTTP traffic. Must be used first for network debugging.

**Parameters:**

```typescript
{
  keyword?: string;   // Filter by URL, headers, or body content
  limit?: number;     // Max results (default: 20, max: 100)
}
```

**Returns:**

```json
{
  "traffic": [
    {
      "id": "uuid",
      "method": "GET",
      "url": "https://api.example.com/users",
      "status": 200,
      "duration": 245,
      "timestamp": "2026-02-22T10:30:00Z"
    }
  ],
  "total": 150
}
```

**Usage:**

```
Show me the last 20 requests to the /api/users endpoint
```

### inspect_request_details

**MANDATORY** tool to retrieve full headers and body for a specific traffic entry.

**Parameters:**

```typescript
{
  trafficId: string;  // ID from inspect_network_traffic
}
```

**Returns:**

```json
{
  "id": "uuid",
  "method": "POST",
  "url": "https://api.example.com/users",
  "requestHeaders": {
    "Content-Type": "application/json",
    "Authorization": "Bearer token123"
  },
  "requestBody": "{\"name\": \"Alice\"}",
  "responseHeaders": {
    "Content-Type": "application/json"
  },
  "responseBody": "{\"id\": 1, \"name\": \"Alice\"}",
  "status": 201,
  "duration": 342
}
```

**Usage:**

```
Show me the full details of request abc-123
```

### execute_request

Execute or replay custom HTTP requests through the proxy.

**Parameters:**

```typescript
{
  method: string;           // HTTP method
  url: string;             // Full URL
  headers?: object;        // Request headers
  body?: string | object;  // Request body
}
```

**Returns:**

```json
{
  "status": 200,
  "headers": {...},
  "body": "...",
  "duration": 234
}
```

**Usage:**

```
Make a POST request to https://api.example.com/users with body {"name": "Bob"}
```

### add_mock_rule

Create a mocking rule to return static responses.

**Parameters:**

```typescript
{
  name: string;            // Descriptive name
  urlPattern: string;      // URL pattern to match
  method: string;          // HTTP method (* for all)
  statusCode: number;      // Response status
  headers?: object;        // Response headers
  body?: string | object;  // Response body
}
```

**Returns:**

```json
{
  "id": "uuid",
  "message": "Mock rule created successfully"
}
```

**Usage:**

```
Create a mock for /api/users that returns a 200 with [{"id": 1, "name": "Alice"}]
```

### add_breakpoint_rule

Create a breakpoint rule to pause traffic for inspection.

**Parameters:**

```typescript
{
  name: string;        // Descriptive name
  urlPattern: string;  // URL pattern to match
  method: string;      // HTTP method
  type: string;        // "request", "response", or "both"
}
```

**Returns:**

```json
{
  "id": "uuid",
  "message": "Breakpoint rule created successfully"
}
```

**Usage:**

```
Set a breakpoint on POST requests to /api/login
```

### list_rules

List all active mocks and breakpoints.

**Parameters:** None

**Returns:**

```json
{
  "rules": [
    {
      "id": "uuid",
      "name": "Mock Users",
      "type": "mock",
      "enabled": true,
      "urlPattern": "https://api.example.com/users*",
      "matchCount": 42
    }
  ]
}
```

**Usage:**

```
Show me all active rules
```

### delete_rule

Remove an interception rule by ID.

**Parameters:**

```typescript
{
  ruleId: string;  // Rule ID from list_rules
}
```

**Returns:**

```json
{
  "message": "Rule deleted successfully"
}
```

**Usage:**

```
Delete the mock rule abc-123
```

### list_scenarios

List all recorded traffic scenarios.

**Parameters:** None

**Returns:**

```json
{
  "scenarios": [
    {
      "id": "uuid",
      "name": "User Login Flow",
      "description": "Complete login workflow",
      "steps": 5,
      "created": "2026-02-22T10:00:00Z"
    }
  ]
}
```

**Usage:**

```
Show me all scenarios
```

### get_scenario

Get full sequence details and variable mappings for a scenario.

**Parameters:**

```typescript
{
  scenarioId: string;  // Scenario ID
}
```

**Returns:**

```json
{
  "id": "uuid",
  "name": "User Login Flow",
  "steps": [...],
  "variableMappings": [...]
}
```

**Usage:**

```
Show me the details of the "User Login" scenario
```

### add_scenario

Create a new scenario with basic metadata.

**Parameters:**

```typescript
{
  name: string;           // Scenario name
  description?: string;   // Optional description
  trafficIds: string[];   // Array of traffic IDs
}
```

**Returns:**

```json
{
  "id": "uuid",
  "message": "Scenario created successfully"
}
```

**Usage:**

```
Create a scenario called "Checkout Flow" from the last 3 requests
```

### update_scenario

Modify an existing scenario.

**Parameters:**

```typescript
{
  scenarioId: string;         // Scenario ID
  name?: string;              // New name
  description?: string;       // New description
  steps?: array;              // Updated steps
  variableMappings?: array;   // Updated mappings
}
```

**Returns:**

```json
{
  "message": "Scenario updated successfully"
}
```

**Usage:**

```
Add the last request to the "Login Flow" scenario
```

### delete_scenario

Remove a scenario by ID.

**Parameters:**

```typescript
{
  scenarioId: string;  // Scenario ID
}
```

**Returns:**

```json
{
  "message": "Scenario deleted successfully"
}
```

**Usage:**

```
Delete the "Old Test" scenario
```

### clear_traffic

Reset/clear the captured traffic logs.

**Parameters:** None

**Returns:**

```json
{
  "message": "Traffic cleared",
  "deleted": 150
}
```

**Usage:**

```
Clear all captured traffic
```

### get_proxy_status

Get real-time proxy address and status.

**Parameters:** None

**Returns:**

```json
{
  "proxyAddress": "http://localhost:15500",
  "dashboardAddress": "http://localhost:15501",
  "status": "running",
  "activeRules": 5
}
```

**Usage:**

```
What's the proxy status?
```

## Prompts

Pre-defined prompts for common workflows.

### analyze-traffic

Analyze recent traffic for errors or anomalies.

**Description:** Have an AI analyze the most recent traffic entries to identify failed requests, unusual patterns, or potential issues.

**Usage:**

```
Analyze recent traffic for errors
```

**What it does:**
1. Fetches recent traffic
2. Identifies failed requests (4xx, 5xx)
3. Analyzes error patterns
4. Suggests potential fixes

### generate-api-docs

Generate OpenAPI documentation from captured traffic.

**Description:** Transform captured HTTP traffic into a structured OpenAPI 3.0 specification document.

**Usage:**

```
Generate API documentation from traffic to /api/users
```

**What it does:**
1. Filters traffic to specified endpoints
2. Extracts request/response schemas
3. Generates OpenAPI spec
4. Includes examples from actual traffic

### generate-scenario-test

Transform a recorded scenario into an automated test script.

**Description:** Convert a traffic scenario into test code for popular frameworks (Playwright, Cypress, Jest).

**Usage:**

```
Generate a Playwright test from the "User Login" scenario
```

**What it does:**
1. Fetches scenario steps
2. Generates test code with proper assertions
3. Handles variable dependencies
4. Includes error handling

## Best Practices

### Tool Ordering

Always use tools in this order:

1. **inspect_network_traffic** - Get overview of traffic
2. **inspect_request_details** - Get full details for specific requests
3. **Other tools** - Create mocks, scenarios, etc.

### Filtering Traffic

Use `keyword` parameter effectively:

```
Show me requests to /api/users with status 500
```

The keyword filters across URL, headers, and body.

### Creating Mocks

Be specific with URL patterns:

- ✅ `https://api.example.com/v1/users/*`
- ❌ `*users*` (too broad)

### Variable Mappings

Use clear, descriptive paths:

```json
{
  "from": "step_1.response.body.auth.token",
  "to": "step_2.request.headers.Authorization",
  "transform": "Bearer {{value}}"
}
```

## Error Handling

All tools return errors in consistent format:

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {}
}
```

Common error codes:

- `NOT_FOUND` - Resource doesn't exist
- `INVALID_PARAMETER` - Invalid parameter value
- `LIMIT_EXCEEDED` - Requested too many results

## Next Steps

- [MCP Integration](index.md) - Setup guide
- [API Reference](/api) - REST API documentation
- [Examples](https://github.com/wahyudotdev/glance/tree/main/examples) - Sample code
