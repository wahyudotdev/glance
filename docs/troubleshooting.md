# Troubleshooting

Common issues and their solutions when using Glance.

## Installation Issues

### Command Not Found After Installation

**Problem**: `glance: command not found` after installing via Homebrew

**Solution**:
```bash
# Check if installed
brew list glance

# Get installation path
which glance

# Add Homebrew bin to PATH if needed
echo 'export PATH="/opt/homebrew/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Port Already in Use

**Problem**: Glance won't start - port 15500, 15501, or 15502 already in use

**Solution**:
```bash
# Find what's using the port
lsof -i :15500

# Kill the process
kill -9 <PID>

# Or start Glance on different ports
glance --proxy-port 8080 --dashboard-port 8081 --mcp-port 8082
```

## Certificate Errors

### SSL/TLS Handshake Failed

**Problem**: Certificate verification errors when making HTTPS requests

**Solution**:

1. **Export the CA certificate**:
   ```bash
   curl http://localhost:15501/ca.crt -o glance-ca.crt
   ```

2. **Trust the certificate**:

   macOS:
   ```bash
   sudo security add-trusted-cert -d -r trustRoot \
     -k /Library/Keychains/System.keychain glance-ca.crt
   ```

   Linux:
   ```bash
   sudo cp glance-ca.crt /usr/local/share/ca-certificates/
   sudo update-ca-certificates
   ```

3. **Verify trust**:
   ```bash
   curl https://api.github.com
   # Should work without errors
   ```

### Certificate Errors in Specific Tools

**cURL**:
```bash
# Use certificate
curl --cacert glance-ca.crt https://example.com

# Or ignore (dev only)
curl -k https://example.com
```

**Git**:
```bash
# Use certificate
git config --global http.sslCAInfo /path/to/glance-ca.crt

# Or disable verification (not recommended)
git config --global http.sslVerify false
```

**Node.js**:
```bash
export NODE_EXTRA_CA_CERTS=/path/to/glance-ca.crt
```

## Traffic Not Appearing

### No Requests in Dashboard

**Problem**: Dashboard is empty even though applications are running

**Diagnosis**:
```bash
# Check if proxy environment variables are set
echo $HTTP_PROXY
echo $HTTPS_PROXY

# Test with curl
curl -x http://localhost:15500 http://httpbin.org/get

# Check Glance logs
glance --log-level debug
```

**Solutions**:

1. **Set proxy variables**:
   ```bash
   export HTTP_PROXY=http://localhost:15500
   export HTTPS_PROXY=http://localhost:15500
   ```

2. **Use one-liner setup**:
   ```bash
   eval "$(curl -s http://localhost:15501/setup)"
   ```

3. **Check application proxy settings**: Some apps ignore environment variables

### Traffic Only Shows HTTP, Not HTTPS

**Problem**: Only HTTP requests appear, HTTPS requests are missing

**Cause**: Certificate not trusted or HTTPS proxy not configured

**Solution**:
1. Trust Glance CA certificate (see Certificate Errors above)
2. Ensure `HTTPS_PROXY` is set:
   ```bash
   export HTTPS_PROXY=http://localhost:15500
   ```

## Performance Issues

### Slow Request Times

**Problem**: Requests through Glance are significantly slower

**Diagnosis**:
```bash
# Compare with and without proxy
time curl https://api.github.com

HTTP_PROXY=http://localhost:15500 time curl https://api.github.com
```

**Solutions**:

1. **Check database size**:
   ```bash
   ls -lh ~/.glance.db
   # If > 1GB, clear old traffic
   ```

2. **Clear traffic**: In dashboard, click "Clear Traffic"

3. **Reduce logging**: Start with less verbose logging:
   ```bash
   glance --log-level warn
   ```

### High Memory Usage

**Problem**: Glance consuming excessive memory

**Solutions**:

1. **Clear traffic** regularly
2. **Limit traffic storage**: Configure max entries (if supported)
3. **Restart Glance** periodically

## Dashboard Issues

### Dashboard Won't Load

**Problem**: `http://localhost:15501` doesn't respond

**Diagnosis**:
```bash
# Check if Glance is running
ps aux | grep glance

# Check if port is listening
lsof -i :15501

# Try accessing via IP
curl http://127.0.0.1:15501
```

**Solutions**:

1. **Restart Glance**:
   ```bash
   killall glance
   glance
   ```

2. **Check firewall**: Ensure localhost connections allowed

3. **Try different browser**: Could be browser extension interfering

### WebSocket Connection Failed

**Problem**: "WebSocket connection failed" error in dashboard

**Cause**: Proxy or VPN interfering with WebSocket connection

**Solution**:

1. **Disable browser proxy** for dashboard:
   - Add `localhost` to proxy exceptions
   - Or access dashboard without proxy

2. **Check browser extensions**: Disable ad blockers temporarily

## Client-Specific Issues

### Docker Containers

**Problem**: Can't capture traffic from Docker containers

**Solution**:

```bash
# Use host network mode
docker run --network host myimage

# Or expose proxy to container
docker run -e HTTP_PROXY=http://host.docker.internal:15500 myimage

# On Linux, use host IP
docker run -e HTTP_PROXY=http://172.17.0.1:15500 myimage
```

### Android Apps

**Problem**: Traffic from Android app not visible

**Common Causes**:

1. **Certificate pinning**: Some apps (banking, security) use pinning
2. **API level 30+**: Requires system certificate or network security config
3. **VPN interference**: Disable VPN on device

**Solutions**:

- **Check if user certificate is trusted**: Settings → Security → Trusted credentials
- **Use Network Security Config** for your app (if you're the developer)
- **Install as system certificate** (requires root)

### Java Applications

**Problem**: Java app ignoring proxy settings

**Solution**:

```bash
# Use JVM arguments
java -Dhttp.proxyHost=localhost \
     -Dhttp.proxyPort=15500 \
     -Dhttps.proxyHost=localhost \
     -Dhttps.proxyPort=15500 \
     -jar app.jar

# Import certificate to Java keystore
keytool -import -trustcacerts -alias glance \
  -file glance-ca.crt \
  -keystore $JAVA_HOME/lib/security/cacerts \
  -storepass changeit
```

## MCP Integration Issues

### Claude Desktop Can't Connect

**Problem**: Claude Desktop doesn't show Glance MCP server

**Diagnosis**:

1. **Check config file** exists:
   ```bash
   # macOS
   cat ~/Library/Application\ Support/Claude/claude_desktop_config.json

   # Windows
   type %APPDATA%\Claude\claude_desktop_config.json
   ```

2. **Verify Glance path**:
   ```bash
   which glance
   # Use this full path in config
   ```

**Solution**:

1. **Correct config format**:
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

2. **Restart Claude Desktop** completely (quit, not just close window)

### MCP Server Port Conflict

**Problem**: MCP server won't start

**Solution**:
```bash
# Use different port
glance --mcp-port 15503

# Update Claude config with new port if needed
```

## Error Messages

### "bind: address already in use"

**Cause**: Another process using Glance's ports

**Solution**: See [Port Already in Use](#port-already-in-use) above

### "database is locked"

**Cause**: Another Glance instance or process accessing the database

**Solution**:
```bash
# Kill all Glance processes
killall glance

# Remove lock file if exists
rm ~/.glance.db-wal ~/.glance.db-shm

# Start Glance
glance
```

### "certificate signed by unknown authority"

**Cause**: Glance CA certificate not trusted

**Solution**: See [Certificate Errors](#certificate-errors) above

## Getting More Help

### Enable Debug Logging

```bash
glance --log-level debug
```

This shows detailed information about:
- Proxy connections
- Certificate handling
- Database operations
- MCP tool calls

### Check Logs

Glance logs to stdout/stderr. Redirect to file:

```bash
glance > glance.log 2>&1
```

### Report an Issue

If you can't solve the problem:

1. **Gather information**:
   ```bash
   # Glance version
   glance --version

   # OS version
   uname -a

   # Debug logs
   glance --log-level debug > debug.log 2>&1
   ```

2. **Create GitHub issue**: [github.com/wahyudotdev/glance/issues](https://github.com/wahyudotdev/glance/issues)

3. **Include**:
   - What you're trying to do
   - What happens instead
   - Steps to reproduce
   - Relevant logs (remove sensitive info!)

## Next Steps

- [FAQ](faq.md) - Frequently asked questions
- [Client Configuration](clients.md) - Platform-specific setup
- [Development Guide](development.md) - Build from source
