# Configuration

Glance can be configured via command-line flags, environment variables, or a configuration file.

## Command-Line Flags

```bash
glance [flags]
```

### Available Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--proxy-port` | `15500` | Proxy server port |
| `--dashboard-port` | `15501` | Dashboard web UI port |
| `--mcp-port` | `15502` | MCP server port |
| `--db-path` | `~/.glance.db` | Path to SQLite database |
| `--log-level` | `info` | Log level (debug, info, warn, error) |
| `--mcp` | `false` | Run in MCP-only mode (for Claude Desktop) |
| `--android` | `false` | Enable Android device auto-configuration |

| `--help` | | Show help message |
| `--version` | | Show version |

### Examples

```bash
# Run on custom ports
glance --proxy-port 8080 --dashboard-port 8081

# Enable debug logging
glance --log-level debug

# Run MCP server only
glance --mcp

# Auto-configure Android device
glance --android



# Custom database location
glance --db-path /tmp/glance-test.db
```

## Environment Variables

Environment variables take precedence over default values but are overridden by command-line flags.

| Variable | Equivalent Flag |
|----------|----------------|
| `GLANCE_PROXY_PORT` | `--proxy-port` |
| `GLANCE_DASHBOARD_PORT` | `--dashboard-port` |
| `GLANCE_MCP_PORT` | `--mcp-port` |
| `GLANCE_DB_PATH` | `--db-path` |
| `GLANCE_LOG_LEVEL` | `--log-level` |

### Example

```bash
export GLANCE_PROXY_PORT=8080
export GLANCE_DASHBOARD_PORT=8081
export GLANCE_LOG_LEVEL=debug
glance
```

## Configuration File

> **Note**: Configuration file support is planned for future releases.

Future versions will support a `glance.yaml` or `.glance.yml` file:

```yaml
# glance.yaml (planned)
proxy:
  port: 15500
  host: localhost

dashboard:
  port: 15501

mcp:
  port: 15502
  enabled: true

database:
  path: ~/.glance.db
  maxSize: 1GB
  autoClean: true

logging:
  level: info
  format: json

rules:
  autoLoad: ./rules.json
```

## Database Configuration

### Location

By default, Glance stores data in `~/.glance.db`.

Change location:
```bash
glance --db-path /custom/path/glance.db
```

### Performance Tuning

The database uses:
- **Write-Ahead Logging (WAL)**: For better concurrency
- **Write-Behind Caching**: For high performance
- **Auto-vacuum**: To manage database size

### Size Management

Monitor database size:
```bash
ls -lh ~/.glance.db
```

Clear old data:
1. Via Dashboard: Click "Clear Traffic"
2. Via API: `DELETE /api/traffic`
3. Via MCP: Use `clear_traffic` tool
4. Manually: `rm ~/.glance.db` (when Glance is stopped)

## Logging

### Log Levels

- **debug**: Verbose logging for development
- **info**: Normal operation (default)
- **warn**: Warnings only
- **error**: Errors only

### Log Format

Default format: human-readable text

```
2026-02-22 10:30:00 INFO  Proxy listening on :15500
2026-02-22 10:30:01 DEBUG Intercepted GET https://api.example.com/users
```

### Log to File

```bash
# Redirect to file
glance > glance.log 2>&1

# Or use tee to see and save
glance 2>&1 | tee glance.log
```

## Network Configuration

### Binding Address

By default, Glance binds to `localhost` (127.0.0.1) for security.

To allow external connections (not recommended):

```bash
# Future feature
glance --bind-address 0.0.0.0
```

⚠️ **Security Warning**: Never expose Glance to untrusted networks. It has no authentication and can intercept sensitive data.

### Custom Ports

If default ports conflict:

```bash
glance \
  --proxy-port 8080 \
  --dashboard-port 8081 \
  --mcp-port 8082
```

Update clients accordingly:
```bash
export HTTP_PROXY=http://localhost:8080
export HTTPS_PROXY=http://localhost:8080
```

## MCP Configuration

### Claude Desktop Setup

Edit `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "glance": {
      "command": "/opt/homebrew/bin/glance",
      "args": ["--mcp"]
    }
  }
}
```

### MCP-Only Mode

Run just the MCP server (no proxy or dashboard):

```bash
glance --mcp
```

Useful for:
- Dedicated MCP instance for AI agents
- Reduced resource usage
- Separate MCP and proxy instances

### Custom MCP Port

```bash
glance --mcp-port 15503
```

## Performance Tuning

### Memory Usage

Limit memory by clearing traffic regularly:

```bash
# Clear traffic older than 1 hour
# (API endpoint - planned feature)
DELETE /api/traffic?older=1h
```

### Database Optimization

For high-traffic scenarios:

1. **Use faster storage**: SSD over HDD
2. **Increase cache**: More RAM for write-behind cache
3. **Regular cleanup**: Clear old traffic data
4. **Partition database**: Separate databases for different projects

### Request Throughput

Glance can handle:
- ~1000 requests/second on modern hardware
- Concurrent connections: 10,000+
- Database writes: Batched for efficiency

## Security Considerations

### Local Only

**Default**: Glance only listens on localhost (127.0.0.1)

Never expose to the internet:
- No authentication
- MITM capabilities
- Sensitive data in database

### CA Certificate

**Location**: Generated in memory, exported via API

**Security**:
- Private key never saved to disk (by default)
- Certificate valid for 1 year
- Only trust on development machines

### Database

**Encryption**: Not encrypted by default

Protect sensitive data:
```bash
# Use encrypted filesystem
# Set restrictive permissions
chmod 600 ~/.glance.db
```

## Multi-Instance Setup

Run multiple Glance instances:

```bash
# Instance 1
glance --proxy-port 15500 --db-path ~/.glance-dev.db

# Instance 2 (different terminal)
glance --proxy-port 16500 --db-path ~/.glance-test.db --dashboard-port 16501
```

Use cases:
- Separate environments (dev, staging, test)
- Different projects
- Isolated test scenarios

## Docker Configuration

> **Note**: Official Docker image planned for future release.

For now, run from source:

```dockerfile
FROM golang:1.24
WORKDIR /app
COPY . .
RUN make build
EXPOSE 15500 15501 15502
CMD ["./glance"]
```

## Next Steps

- [Troubleshooting](troubleshooting.md) - Common configuration issues
- [API Reference](api.md) - Programmatic configuration
- [Development Guide](development.md) - Build and customize
