# Changelog

All notable changes to the Glance project will be documented in this file.

## [v0.2.0] - 2026-02-23
- **Docker Interception**: Introduced a robust "One-Click Intercept" for Docker containers via container recreation and environment injection.
- **Java Container Support**: Added specialized logic to automatically inject and trust the Glance CA certificate into the Java Trust Store (`cacerts`) inside Docker containers using `keytool`.
- **Universal Compatibility**: Replaced the `iptables` approach with standard `HTTP_PROXY`/`HTTPS_PROXY` and `JAVA_TOOL_OPTIONS` injection, ensuring compatibility across all container OS flavors (Alpine, Debian, RHEL) without requiring `NET_ADMIN` capabilities.
- **Docker UI**: New full-width Docker Containers card in the Integrations view with real-time status tracking, container renaming for backups, and seamless restoration.
- **Improved IP Detection**: Enhanced Host IP discovery to prioritize `host.docker.internal` for Docker Desktop users, ensuring reliable connectivity from container to host.
- **API Stability**: Fixed Docker SDK compatibility issues for older daemon versions (API 1.41) and properly handled binary stream demultiplexing for command execution.

## [v0.1.4] - 2026-02-21
- **MCP Enhancements**: Renamed tools to `inspect_network_traffic` and `inspect_request_details` with authoritative descriptions to guide AI agents.
- **MCP Config**: Added a configurable `limit` parameter to `inspect_network_traffic` (default 20, max 100).
- **Traffic Filtering**: Enhanced search bar with a professional toggleable filter menu for request methods (GET, POST, etc.).
- **JSON Tree Editor**: Introduced a new interactive JSON editor with tree view, collapse/expand support, and direct value editing. Integrated this into both Rule/Response editors and the Traffic Details panel.
- **Enhanced Editor**: Added search and full-screen editing capabilities for mock response bodies.
- **Scenario Fix**: Fixed a bug where scenario steps would show "Unknown URL" if the original traffic entry was not in the current view history. Steps now join with historical traffic data from the database.

## [v0.1.3] - 2026-02-20
- **UI/UX**: Fixed cURL command overflow in the details panel by adding proper text wrapping.
- **About Page**: Introduced a dedicated About page with project overview and open-source attributions.
- **UI Cleanup**: Streamlined the Integrations view by moving secondary documentation to the About page.
- **Documentation**: Added official Homebrew installation guide and compatibility matrix to the README and Dashboard.
- **Project Hygiene**: Established a dedicated `CHANGELOG.md` and integrated a "What's New" modal into the dashboard.

## [v0.1.2] - 2026-02-18
- **UI/UX**: Added full **Dark Mode** support and theme persistence.
- **Refactoring**: Extracted core logic into modular React hooks and standardized API patterns.
- **Stability**: Resolved database deadlocks and improved thread-safety in the WebSocket hub.

## [v0.1.1] - 2026-02-18
- **Release Pipeline**: Established automated multi-platform releases via GoReleaser and Homebrew.
- **Scenarios**: Implemented **Scenario Recording** with variable mapping and AI test generation prompts.
- **MCP**: Optimized server performance and moved to a dedicated port for better compatibility.
- **Real-time Stability**: Fixed critical WebSocket issues and race conditions in traffic streaming.
- **Project Hygiene**: Added MIT License, updated gitignore, and implemented safe URL parsing.
- **Build & Stability**: Resolved database deadlocks and fixed mcp-go library integration issues.

## [v0.1.0] - 2026-02-17
- **Initial Release**: Launched Glance as a specialized MITM proxy for AI Agents.
- **Traffic Inspection**: Real-time traffic inspection with request/response detail views and cURL export.
- **Rule Engine**: Unified management of Mocks and Breakpoints for traffic modification.
- **Java Integration**: Native JVM agent for auto-intercepting HttpsURLConnection in Java 8+ apps.
- **Android Support**: ADB-based device discovery, CA certificate installation, and proxy configuration.
- **Chromium Support**: Auto-launching Chromium instances with pre-configured proxy and certificate flags.
- **MCP Server**: Native implementation of the Model Context Protocol for AI Agent tools and resources.
- **Dashboard**: Modern React-based dashboard for real-time traffic monitoring and configuration.
