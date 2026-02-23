# Quick Start Guide

This guide will help you get started with Glance in just a few minutes.

## Step 1: Start Glance

Run Glance from your terminal:

```bash
glance
```

You should see output similar to:

```
   ______ _                             
  / ____/| |                          
 | |  __ | |  ____ _  _ __    ____  ___ 
 | | |_ || | / _' | || '_ \  / __|/ _ \
 | |__| || || (_| | || | | || (__|  __/
  \______|_| \__,_||_||_| |_| \___|\___|
                                        
  Let Your AI Understand Every Request at a Glance.

[✓] Proxy server running on localhost:15500
[✓] MCP server (SSE) running on localhost:15502/mcp
[✓] API server running on localhost:15501
[✓] Dashboard available at http://localhost:15501
```

By default, Glance runs on these ports:

- **Dashboard**: `http://localhost:15501` - Web UI for viewing and managing traffic
- **Proxy**: `http://localhost:15500` - The actual MITM proxy
- **MCP Server**: `http://localhost:15502` - Model Context Protocol server for AI agents

## Step 2: Configure Your Client

You need to configure your application or environment to route traffic through Glance's proxy.

### Option 1: One-liner Setup (Terminal)

The easiest way to get started is using the one-liner setup command:

```bash
eval "$(curl -s http://localhost:15501/setup)"
```

This automatically sets up your current terminal session to use Glance as a proxy.

### Option 2: Manual Environment Variables

Set the proxy environment variables manually:

```bash
export HTTP_PROXY=http://localhost:15500
export HTTPS_PROXY=http://localhost:15500
```

For a persistent setup, add these lines to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.).

### Option 3: Application-Specific Configuration

See the [Client Configuration Guide](clients.md) for specific instructions for:

- [Java/JVM applications](clients/java.md)
- [Docker containers](clients/docker.md)
- [Android devices](clients/android.md)
- [Chromium browsers](clients/chromium.md)
- [Terminal/CLI tools](clients/terminal.md)

## Step 3: View Traffic in the Dashboard

1. Open your browser and navigate to `http://localhost:15501`
2. You'll see the Glance dashboard with a live view of captured traffic
3. Make some HTTP/HTTPS requests from your configured client
4. Watch the requests appear in real-time!

## Step 4: Inspect Request Details

Click on any request in the traffic list to view:

- **Request Headers**: All HTTP headers sent with the request
- **Request Body**: Formatted JSON, XML, or raw body content
- **Response Headers**: All HTTP headers in the response
- **Response Body**: Formatted response content
- **Timing Information**: Request duration and timestamps
- **cURL Export**: Copy the request as a `curl` command

## Step 5: Try Creating a Mock

Let's create a simple mock to return a custom response:

1. In the dashboard, click the **"Rules"** tab
2. Click **"Add Mock"**
3. Configure the mock:
   - **URL Pattern**: `https://api.example.com/users`
   - **Method**: `GET`
   - **Status Code**: `200`
   - **Response Body**:
     ```json
     {
       "users": [
         {"id": 1, "name": "Alice"},
         {"id": 2, "name": "Bob"}
       ]
     }
     ```
4. Click **"Save"**
5. Now any request to `https://api.example.com/users` will return your mocked response!

## Step 6 (Optional): Test with cURL

Test your setup with a simple cURL command:

```bash
curl -x http://localhost:15500 https://api.github.com/users
```

You should see the GitHub API response, and the request will appear in the Glance dashboard.

## What's Next?

Now that you have Glance running, explore these features:

- 🔍 [Traffic Inspection](features/traffic-inspection.md) - Learn about advanced inspection features
- 🛠️ [Mocking & Breakpoints](features/mocking.md) - Master the rule engine
- 📝 [Scenario Recording](features/scenarios.md) - Record and replay request sequences
- 🤖 [MCP Integration](mcp/) - Connect with AI agents like Claude Desktop

## Common Issues

### HTTPS Certificate Errors

If you see certificate errors, you need to trust Glance's CA certificate. See the [Troubleshooting Guide](troubleshooting.md#certificate-errors) for instructions.

### Port Already in Use

If the default ports are already in use, you can start Glance on different ports:

```bash
glance --proxy-port 8080 --dashboard-port 8081 --mcp-port 8082
```

### Can't See Any Traffic

Make sure your client is properly configured to use the proxy. Check the [Client Configuration Guide](clients.md) for your specific environment.

## Need Help?

- 📖 Read the [FAQ](faq.md)
- 🐛 Check the [Troubleshooting Guide](troubleshooting.md)
- 💬 [Open an issue](https://github.com/wahyudotdev/glance/issues) on GitHub
