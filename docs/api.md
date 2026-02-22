# API Reference

Glance provides a RESTful API for programmatic access to all features.

## Base URL

```
http://localhost:15501/api
```

## Authentication

Currently, Glance runs locally without authentication. Future versions may add API keys for remote access.

## Traffic API

### List Traffic

Get a list of captured HTTP requests.

```http
GET /api/traffic
```

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `limit` | integer | Max results to return (default: 100) |
| `offset` | integer | Offset for pagination (default: 0) |
| `method` | string | Filter by HTTP method (GET, POST, etc.) |
| `search` | string | Search in URL, headers, and body |

**Response:**

```json
{
  "traffic": [
    {
      "id": "uuid",
      "method": "GET",
      "url": "https://api.example.com/users",
      "status": 200,
      "duration": 245,
      "timestamp": "2026-02-22T10:30:00Z",
      "requestHeaders": {...},
      "responseHeaders": {...},
      "requestBody": "...",
      "responseBody": "..."
    }
  ],
  "total": 150,
  "limit": 100,
  "offset": 0
}
```

### Get Traffic Details

Get full details for a specific request.

```http
GET /api/traffic/:id
```

**Response:**

```json
{
  "id": "uuid",
  "method": "GET",
  "url": "https://api.example.com/users",
  "status": 200,
  "duration": 245,
  "timestamp": "2026-02-22T10:30:00Z",
  "requestHeaders": {
    "User-Agent": "curl/7.79.1",
    "Accept": "*/*"
  },
  "responseHeaders": {
    "Content-Type": "application/json",
    "Content-Length": "1234"
  },
  "requestBody": "",
  "responseBody": "{\"users\": [...]}"
}
```

### Clear Traffic

Delete all captured traffic.

```http
DELETE /api/traffic
```

**Response:**

```json
{
  "message": "Traffic cleared",
  "deleted": 150
}
```

## Rules API

### List Rules

Get all mocking and breakpoint rules.

```http
GET /api/rules
```

**Response:**

```json
{
  "rules": [
    {
      "id": "uuid",
      "name": "Mock Users API",
      "type": "mock",
      "enabled": true,
      "urlPattern": "https://api.example.com/users*",
      "method": "GET",
      "statusCode": 200,
      "headers": {...},
      "body": {...},
      "matchCount": 42,
      "lastMatched": "2026-02-22T10:30:00Z"
    }
  ]
}
```

### Create Mock Rule

Create a new mock rule.

```http
POST /api/rules/mocks
```

**Request Body:**

```json
{
  "name": "Mock Users API",
  "urlPattern": "https://api.example.com/users*",
  "method": "GET",
  "statusCode": 200,
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    "users": [
      {"id": 1, "name": "Alice"}
    ]
  }
}
```

**Response:**

```json
{
  "id": "uuid",
  "message": "Mock rule created"
}
```

### Create Breakpoint Rule

Create a new breakpoint rule.

```http
POST /api/rules/breakpoints
```

**Request Body:**

```json
{
  "name": "Debug Login",
  "urlPattern": "https://api.example.com/login",
  "method": "POST",
  "type": "request" // or "response" or "both"
}
```

**Response:**

```json
{
  "id": "uuid",
  "message": "Breakpoint rule created"
}
```

### Update Rule

Update an existing rule.

```http
PUT /api/rules/:id
```

**Request Body:** Same as create, with fields to update

### Delete Rule

Delete a rule.

```http
DELETE /api/rules/:id
```

**Response:**

```json
{
  "message": "Rule deleted"
}
```

### Toggle Rule

Enable or disable a rule.

```http
PATCH /api/rules/:id/toggle
```

**Response:**

```json
{
  "id": "uuid",
  "enabled": false
}
```

## Scenarios API

### List Scenarios

Get all recorded scenarios.

```http
GET /api/scenarios
```

**Response:**

```json
{
  "scenarios": [
    {
      "id": "uuid",
      "name": "User Login Flow",
      "description": "Complete login workflow",
      "steps": 5,
      "created": "2026-02-22T10:00:00Z",
      "modified": "2026-02-22T10:30:00Z"
    }
  ]
}
```

### Get Scenario

Get full scenario details including steps and variable mappings.

```http
GET /api/scenarios/:id
```

**Response:**

```json
{
  "id": "uuid",
  "name": "User Login Flow",
  "description": "Complete login workflow",
  "steps": [
    {
      "order": 1,
      "trafficId": "uuid",
      "method": "GET",
      "url": "https://example.com/login",
      "status": 200
    }
  ],
  "variableMappings": [
    {
      "from": "step_1.response.body.token",
      "to": "step_2.request.headers.Authorization"
    }
  ],
  "created": "2026-02-22T10:00:00Z",
  "modified": "2026-02-22T10:30:00Z"
}
```

### Create Scenario

Create a new scenario.

```http
POST /api/scenarios
```

**Request Body:**

```json
{
  "name": "User Login Flow",
  "description": "Complete login workflow",
  "trafficIds": ["uuid1", "uuid2", "uuid3"]
}
```

### Update Scenario

Update scenario details or steps.

```http
PUT /api/scenarios/:id
```

### Delete Scenario

Delete a scenario.

```http
DELETE /api/scenarios/:id
```



## Proxy API

### Get Proxy Status

Get current proxy status and configuration.

```http
GET /api/proxy/status
```

**Response:**

```json
{
  "proxyAddress": "http://localhost:15500",
  "dashboardAddress": "http://localhost:15501",
  "mcpAddress": "http://localhost:15502",
  "status": "running",
  "activeRules": 5,
  "totalTraffic": 1234,
  "uptime": 3600
}
```

### Get CA Certificate

Download the CA certificate.

```http
GET /api/proxy/ca.crt
```

Returns the PEM-encoded CA certificate.

## WebSocket API

### Traffic Stream

Real-time traffic updates via WebSocket.

```javascript
const ws = new WebSocket('ws://localhost:15501/ws/traffic');

ws.onmessage = (event) => {
  const traffic = JSON.parse(event.data);
  console.log('New request:', traffic);
};
```

**Message Format:**

```json
{
  "type": "traffic",
  "data": {
    "id": "uuid",
    "method": "GET",
    "url": "https://api.example.com/users",
    "status": 200,
    "timestamp": "2026-02-22T10:30:00Z"
  }
}
```

## Error Responses

All endpoints return consistent error format:

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {}
}
```

**Common HTTP Status Codes:**

| Code | Meaning |
|------|---------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid parameters |
| 404 | Not Found - Resource doesn't exist |
| 409 | Conflict - Duplicate resource |
| 500 | Internal Server Error |

## Rate Limiting

Currently, no rate limiting is enforced for local API access.

## Examples

### cURL Examples

```bash
# List traffic
curl http://localhost:15501/api/traffic

# Get specific request
curl http://localhost:15501/api/traffic/uuid

# Create mock
curl -X POST http://localhost:15501/api/rules/mocks \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Test Mock",
    "urlPattern": "https://api.example.com/*",
    "method": "GET",
    "statusCode": 200,
    "body": {"message": "mocked"}
  }'

# Clear traffic
curl -X DELETE http://localhost:15501/api/traffic
```

### JavaScript/Fetch

```javascript
// List traffic
const response = await fetch('http://localhost:15501/api/traffic');
const data = await response.json();

// Create mock
await fetch('http://localhost:15501/api/rules/mocks', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    name: 'Test Mock',
    urlPattern: 'https://api.example.com/*',
    method: 'GET',
    statusCode: 200,
    body: { message: 'mocked' }
  })
});
```

### Python

```python
import requests

# List traffic
response = requests.get('http://localhost:15501/api/traffic')
traffic = response.json()

# Create mock
mock = {
    'name': 'Test Mock',
    'urlPattern': 'https://api.example.com/*',
    'method': 'GET',
    'statusCode': 200,
    'body': {'message': 'mocked'}
}
response = requests.post(
    'http://localhost:15501/api/rules/mocks',
    json=mock
)
```

## Next Steps

- [MCP Reference](mcp/reference.md) - MCP tools and resources
- [Development Guide](development.md) - Contribute to the API
- [Examples](https://github.com/wahyudotdev/glance/tree/main/examples) - More code samples
