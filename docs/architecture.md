# Architecture

Technical overview of Glance's architecture and design decisions.

## High-Level Overview

```
┌─────────────────────────────────────────────────────────┐
│                      Client Apps                         │
│    (Browser, cURL, Java, Android, Node.js, etc.)        │
└───────────────────────┬─────────────────────────────────┘
                        │ HTTP/HTTPS
                        ▼
┌─────────────────────────────────────────────────────────┐
│                   Glance Proxy (:15500)                  │
│  ┌──────────────────────────────────────────────────┐  │
│  │           MITM Proxy (goproxy)                   │  │
│  │  • Certificate generation                        │  │
│  │  • Request/Response interception                 │  │
│  │  • Rule matching & execution                     │  │
│  └──────────────────────────────────────────────────┘  │
│                        │                                 │
│                        ▼                                 │
│  ┌──────────────────────────────────────────────────┐  │
│  │            Storage Layer (SQLite)                │  │
│  │  • Traffic logging                               │  │
│  │  • Rule persistence                              │  │
│  │  • Scenario recording                            │  │
│  └──────────────────────────────────────────────────┘  │
└───────────┬──────────────────────────────┬──────────────┘
            │                              │
            ▼                              ▼
┌────────────────────────┐    ┌──────────────────────────┐
│  Dashboard (:15501)    │    │    MCP Server (:15502)   │
│  ┌──────────────────┐  │    │  ┌────────────────────┐  │
│  │  API (Fiber)     │  │    │  │  MCP Protocol      │  │
│  │  • REST endpoints│  │    │  │  • Resources       │  │
│  │  • WebSocket hub │  │    │  │  • Tools           │  │
│  └──────────────────┘  │    │  │  • Prompts         │  │
│  ┌──────────────────┐  │    │  └────────────────────┘  │
│  │  Frontend (React)│  │    │                          │
│  │  • TypeScript    │  │    │  Used by:                │
│  │  • Tailwind CSS  │  │    │  • Claude Desktop        │
│  │  • Real-time UI  │  │    │  • Other MCP clients     │
│  └──────────────────┘  │    │                          │
└────────────────────────┘    └──────────────────────────┘
```

## Components

### 1. MITM Proxy

**Technology**: [elazarl/goproxy](https://github.com/elazarl/goproxy)

**Responsibilities**:
- Intercept HTTP and HTTPS traffic
- Generate and sign certificates on-the-fly
- Apply rules (mocks, breakpoints)
- Forward modified traffic

**Key Features**:
- **Certificate Authority**: Generates a local CA certificate on startup
- **Dynamic Certificates**: Creates valid certificates for each intercepted domain
- **Request Pipeline**: Processes requests through rule engine before forwarding
- **Response Pipeline**: Processes responses before returning to client

**Flow**:
1. Client makes HTTPS request
2. Proxy intercepts, performs TLS handshake with dynamic certificate
3. Proxy checks rules (mocks, breakpoints)
4. If no rule matches, forwards to destination
5. Logs request/response to database
6. Returns response to client

### 2. Storage Layer

**Technology**: [modernc.org/sqlite](https://modernc.org/sqlite) (Pure Go SQLite)

**Schema**:

```sql
-- Traffic table
CREATE TABLE traffic (
    id TEXT PRIMARY KEY,
    method TEXT,
    url TEXT,
    status INTEGER,
    duration INTEGER,
    timestamp DATETIME,
    request_headers TEXT,  -- JSON
    request_body BLOB,
    response_headers TEXT, -- JSON
    response_body BLOB
);

-- Rules table
CREATE TABLE rules (
    id TEXT PRIMARY KEY,
    name TEXT,
    type TEXT,  -- 'mock' or 'breakpoint'
    enabled BOOLEAN,
    url_pattern TEXT,
    method TEXT,
    config TEXT,  -- JSON
    match_count INTEGER,
    last_matched DATETIME
);

-- Scenarios table
CREATE TABLE scenarios (
    id TEXT PRIMARY KEY,
    name TEXT,
    description TEXT,
    steps TEXT,  -- JSON array
    variable_mappings TEXT,  -- JSON
    created DATETIME,
    modified DATETIME
);
```

**Optimizations**:
- **WAL Mode**: Write-Ahead Logging for better concurrency
- **Write-Behind Cache**: Batches writes for performance
- **Indexes**: On url, method, timestamp for fast queries
- **Prepared Statements**: Reused for common queries

**Database Size Management**:
- Auto-vacuum enabled
- Configurable retention policy (planned)
- Manual clear via API/Dashboard

### 3. Dashboard API

**Technology**: [gofiber/fiber](https://github.com/gofiber/fiber) v2

**Endpoints**:

```
GET    /api/traffic              - List traffic
GET    /api/traffic/:id          - Get traffic details
DELETE /api/traffic              - Clear all traffic

GET    /api/rules                - List rules
POST   /api/rules/mocks          - Create mock
POST   /api/rules/breakpoints    - Create breakpoint
PUT    /api/rules/:id            - Update rule
DELETE /api/rules/:id            - Delete rule

GET    /api/scenarios            - List scenarios
GET    /api/scenarios/:id        - Get scenario
POST   /api/scenarios            - Create scenario
PUT    /api/scenarios/:id        - Update scenario
DELETE /api/scenarios/:id        - Delete scenario
POST   /api/scenarios/:id/replay - Replay scenario

GET    /api/proxy/status         - Proxy status
GET    /api/proxy/ca.crt         - Download CA cert

WS     /ws/traffic               - Real-time traffic stream
```

**WebSocket Hub**:
- Manages concurrent WebSocket connections
- Broadcasts new traffic to all connected clients
- Handles connection lifecycle (connect, disconnect, error)

**Architecture**:
```
Request → Middleware → Handler → Storage → Response
                ↓
         WebSocket Hub → Broadcast
```

### 4. Frontend

**Technology**: React 18, TypeScript, Tailwind CSS

**Architecture**:

```
src/
├── components/        # React components
│   ├── Traffic/       # Traffic list and details
│   ├── Rules/         # Rule management
│   ├── Scenarios/     # Scenario recording
│   └── Common/        # Shared UI components
├── hooks/             # Custom React hooks
│   ├── useTraffic     # Traffic state management
│   ├── useRules       # Rule state management
│   └── useWebSocket   # WebSocket connection
├── services/          # API clients
│   ├── api.ts         # REST API client
│   └── websocket.ts   # WebSocket client
└── types/             # TypeScript types
```

**State Management**:
- React hooks for local state
- WebSocket for real-time updates
- API calls for CRUD operations

**Key Features**:
- **Real-time Updates**: WebSocket connection for live traffic
- **Dark Mode**: Theme persistence with localStorage
- **Responsive**: Mobile-friendly design
- **Performance**: Virtualized lists for large datasets

### 5. MCP Server

**Technology**: [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk)

**Implementation**:

```go
type GlanceMCP struct {
    storage *Storage
    proxy   *Proxy
}

func (m *GlanceMCP) Resources() []Resource {
    return []Resource{
        {URI: "proxy://status", ...},
        {URI: "traffic://latest", ...},
    }
}

func (m *GlanceMCP) Tools() []Tool {
    return []Tool{
        {Name: "inspect_network_traffic", ...},
        {Name: "inspect_request_details", ...},
        {Name: "execute_request", ...},
        // ... more tools
    }
}

func (m *GlanceMCP) Prompts() []Prompt {
    return []Prompt{
        {Name: "analyze-traffic", ...},
        {Name: "generate-api-docs", ...},
        {Name: "generate-scenario-test", ...},
    }
}
```

**Communication**:
- Stdio transport for Claude Desktop
- HTTP transport for other clients (planned)
- JSON-RPC 2.0 protocol

## Data Flow

### Capturing Traffic

```
Client → Proxy → Rule Check → Forward → Server
  ↓         ↓                              ↓
  ↓         ↓                         Response
  ↓         ↓ ← Log to DB ← ───────────────┘
  ↓         ↓
  ↓    Broadcast to WebSocket
  ↓         ↓
  ↓    Dashboard Update
  ↓
Response
```

### Creating a Mock

```
User (Dashboard/MCP)
  ↓
API/MCP Server
  ↓
Storage.SaveRule()
  ↓
Proxy reloads rules
  ↓
Next matching request returns mock
```

### Scenario Recording

```
1. User starts recording
2. Traffic captured → Steps added to scenario
3. User defines variable mappings
4. Scenario saved to database
5. Can replay or generate test code
```

## Security Model

### Certificate Management

```
Startup:
  ↓
Generate CA cert + private key (in memory)
  ↓
Sign certificates dynamically for each domain
  ↓
Client trusts CA → Accepts dynamic certs
```

**Security Notes**:
- CA private key stored in memory only (by default)
- Certificates valid for 1 year
- No external CA communication
- Local-only by design

### Threat Model

**In Scope**:
- Local development machine
- Trusted user
- Controlled network

**Out of Scope**:
- Production use
- Untrusted networks
- Multi-user scenarios
- Remote access

## Performance Characteristics

### Throughput

- **Requests/Second**: ~1,000 on modern hardware
- **Concurrent Connections**: 10,000+
- **Latency Overhead**: < 10ms typical

### Memory Usage

- **Base**: ~50MB (empty database)
- **Per Request**: ~10KB (metadata only)
- **Database**: Grows with traffic (SQLite efficient)
- **Frontend**: ~20MB (React bundle)

### Scalability Limits

Current architecture suitable for:
- ✅ Individual developer use
- ✅ Small team (< 10 users)
- ✅ Moderate traffic (< 10,000 req/day)
- ❌ High-traffic production monitoring
- ❌ Multi-tenant SaaS

## Build Process

### Development

```
Frontend:        Backend:
  ↓                ↓
npm install      go mod download
  ↓                ↓
npm run dev      go run cmd/glance/main.go
```

### Production Build

```
1. Build Frontend:
   npm run build → dist/

2. Embed Frontend:
   go:embed dist → Go binary

3. Build Binary:
   go build -o glance

Result: Single self-contained binary
```

### Cross-Platform Build

Uses [GoReleaser](https://goreleaser.com/):

```yaml
builds:
  - goos: [darwin, linux, windows]
    goarch: [amd64, arm64]
    ldflags: -s -w
```

Produces:
- macOS (ARM64, AMD64)
- Linux (AMD64, ARM64, ARM)
- Windows (AMD64)

## Design Decisions

### Why Go?

- **Performance**: Fast, compiled language
- **Concurrency**: Great for proxy server
- **Single Binary**: Easy distribution
- **Cross-Platform**: Build for all platforms

### Why SQLite?

- **Serverless**: No separate database process
- **Portable**: Single file database
- **Fast**: Excellent for read-heavy workloads
- **Reliable**: Battle-tested, stable

### Why React?

- **Ecosystem**: Large community, many libraries
- **Performance**: Virtual DOM, efficient updates
- **TypeScript**: Type safety for large codebase
- **Developer Experience**: Great tooling

### Why Fiber?

- **Performance**: Faster than standard net/http
- **API**: Express-like, familiar to web developers
- **WebSocket**: Built-in support
- **Middleware**: Rich ecosystem

## Future Improvements

### Planned Features

- **GraphQL API**: More flexible querying
- **Plugins**: Extensibility for custom processors
- **Clustering**: Multi-instance coordination
- **Streaming**: Process large responses efficiently
- **gRPC Support**: Intercept gRPC traffic

### Performance Optimizations

- **Connection Pooling**: Reuse connections
- **Response Streaming**: Don't buffer large responses
- **Selective Logging**: Skip binary/large responses
- **Database Sharding**: Split by project/time

### Architecture Evolution

Current: **Monolith** (Single binary)
Future: **Modular** (Pluggable components)

```
glance-core      (proxy + storage)
glance-dashboard (web UI - optional)
glance-mcp       (MCP server - optional)
glance-plugins   (custom processors)
```

## Contributing

See [Development Guide](development.md) for:
- Setting up development environment
- Code structure and conventions
- Testing guidelines
- Pull request process

## Next Steps

- [Development Guide](development.md) - Contribute to Glance
- [API Reference](api.md) - Build integrations
- [MCP Reference](mcp/reference.md) - Extend MCP capabilities
