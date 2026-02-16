# Requirements: Golang HTTP Debugging Proxy (Agent Proxy)

An application inspired by HTTP Toolkit for intercepting, inspecting, and mocking HTTP/HTTPS traffic.

## 1. Core Functional Requirements

### 1.1 Traffic Interception
- **Transparent Proxying:** Act as an HTTP/HTTPS proxy (CONNECT method support).
- **HTTPS Decryption (MITM):**
    - Generate and manage a root Certificate Authority (CA).
    - Dynamically generate certificates for intercepted domains.
- **Support for Protocols:**
    - HTTP/1.1
    - HTTP/2
    - WebSockets (optional/future)

### 1.2 Inspection
- **Live View:** Real-time stream of captured requests and responses.
- **Detailed Metadata:**
    - Request: Method, URL, Headers, Body (formatted JSON/XML/HTML), Query Parameters.
    - Response: Status Code, Headers, Body, Timing/Performance data.
- **Export Capabilities:**
    - **Copy as cURL:** Generate a standard `curl` command for any captured request for easy reproduction in terminal.
- **Search & Filtering:** Filter by host, method, status code, or content type.

### 1.3 Modification & Mocking
- **Rule Engine:** Define rules to intercept and modify traffic on the fly.
- **Request Rewriting:** Modify headers, query params, or body before it reaches the server.
- **Response Rewriting:** Modify headers or body before it reaches the client.
- **Mocking:** Return static responses (Status, Headers, Body) for specific URL patterns.
- **Breakpoint/Manual Edit:** Pause requests/responses for manual modification (advanced).

### 1.4 Client Integration
- **System Proxy Setup:** Automated configuration of system-wide proxy settings.
- **Environment Injection:** Helper for setting `HTTP_PROXY` and `HTTPS_PROXY` in terminal sessions.
- **CA Certificate Installation:** Guide or automation for installing the CA in system/browser stores.
- **Chromium/Chrome Integration:** 
    - Auto-launch browser instances with pre-configured proxy and certificate flags (e.g., `--proxy-server`, `--ignore-certificate-errors` for local dev).
- **Android Integration:**
    - Provide instructions and tools for ADB-based proxy configuration.
    - Guide for CA certificate installation on Android (system vs. user store).
    - Support for intercepting emulator traffic.

### 1.5 MCP Server Integration (Model Context Protocol)
- **MCP Compliance:** Implement the Model Context Protocol to function as an MCP Server.
- **Resources:**
    - `proxy://status`: Current proxy address, dashboard URL, and count of active rules.
    - `traffic://latest`: The most recent 10 HTTP requests captured in JSON format.
- **Tools:**
    - `list_traffic`: List captured HTTP traffic summaries (up to 20 entries), with optional keyword filtering.
    - `get_traffic_details`: Retrieve full headers and body for a specific traffic entry by ID.
    - `add_mock_rule`: Programmatically create mocking rules to intercept and return static responses.
    - `clear_traffic`: Reset/clear the captured traffic logs.
    - `get_proxy_status`: Get real-time proxy address and status.
- **Prompts:**
    - `analyze-traffic`: Pre-defined prompt to have an AI analyze recent traffic for errors or anomalies.
    - `generate-api-docs`: Pre-defined prompt to generate OpenAPI documentation from the captured traffic.

## 2. User Interface
- **Web Dashboard:** A modern, responsive web UI to manage the proxy and view traffic.
- **CLI Interface:** Basic status and logs in the terminal.

## 3. Non-Functional Requirements
- **Performance:** Low latency overhead for proxied traffic.
- **Portability:** Compile to a single binary (Go advantage).
- **Security:** Ensure the custom CA is generated locally and handled securely.

## 4. Technical Stack (Proposed)
- **Language:** Go 1.24+
- **Proxy Logic:** `net/http` (standard library) or libraries like `elazarl/goproxy`.
- **Frontend:** React or Vue (communicating via WebSockets/REST).
- **Storage:** SQLite or In-memory for active session data.
