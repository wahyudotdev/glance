# Changelog

All notable changes to the Glance project are documented here.

## [v0.1.4] - 2026-02-21

### Added
- **MCP Enhancements**: Renamed tools to `inspect_network_traffic` and `inspect_request_details` with authoritative descriptions to guide AI agents
- **MCP Config**: Added configurable `limit` parameter to `inspect_network_traffic` (default 20, max 100)
- **Traffic Filtering**: Enhanced search bar with professional toggleable filter menu for request methods (GET, POST, etc.)
- **JSON Tree Editor**: Introduced interactive JSON editor with tree view, collapse/expand support, and direct value editing
- **Enhanced Editor**: Added search and full-screen editing capabilities for mock response bodies

### Fixed
- **Scenario Fix**: Fixed bug where scenario steps would show "Unknown URL" if original traffic entry was not in current view history
- Steps now join with historical traffic data from database

## [v0.1.3] - 2026-02-20

### Added
- **About Page**: Introduced dedicated About page with project overview and open-source attributions
- **Documentation**: Added official Homebrew installation guide and compatibility matrix to README and Dashboard
- **Project Hygiene**: Established dedicated `CHANGELOG.md` and integrated "What's New" modal into dashboard

### Changed
- **UI Cleanup**: Streamlined Integrations view by moving secondary documentation to About page

### Fixed
- **UI/UX**: Fixed cURL command overflow in details panel by adding proper text wrapping

## [v0.1.2] - 2026-02-18

### Added
- **UI/UX**: Added full **Dark Mode** support with theme persistence

### Changed
- **Refactoring**: Extracted core logic into modular React hooks and standardized API patterns

### Fixed
- **Stability**: Resolved database deadlocks and improved thread-safety in WebSocket hub

## [v0.1.1] - 2026-02-18

### Added
- **Release Pipeline**: Established automated multi-platform releases via GoReleaser and Homebrew
- **Scenarios**: Implemented **Scenario Recording** with variable mapping and AI test generation prompts
- **Project Hygiene**: Added MIT License and updated gitignore

### Changed
- **MCP**: Optimized server performance and moved to dedicated port for better compatibility

### Fixed
- **Real-time Stability**: Fixed critical WebSocket issues and race conditions in traffic streaming
- **Build & Stability**: Resolved database deadlocks and fixed mcp-go library integration issues
- **Security**: Implemented safe URL parsing

## [v0.1.0] - 2026-02-17

### Added
- **Initial Release**: Launched Glance as a specialized MITM proxy for AI Agents
- **Traffic Inspection**: Real-time traffic inspection with request/response detail views and cURL export
- **Rule Engine**: Unified management of Mocks and Breakpoints for traffic modification
- **Java Integration**: Native JVM agent for auto-intercepting HttpsURLConnection in Java 8+ apps
- **Android Support**: ADB-based device discovery, CA certificate installation, and proxy configuration
- **Chromium Support**: Auto-launching Chromium instances with pre-configured proxy and certificate flags
- **MCP Server**: Native implementation of Model Context Protocol for AI Agent tools and resources
- **Dashboard**: Modern React-based dashboard for real-time traffic monitoring and configuration

---

## Version History

- **v0.1.4** (2026-02-21) - MCP enhancements, JSON editor, filtering improvements
- **v0.1.3** (2026-02-20) - About page, documentation improvements
- **v0.1.2** (2026-02-18) - Dark mode, stability fixes
- **v0.1.1** (2026-02-18) - Scenarios, release pipeline, MCP optimizations
- **v0.1.0** (2026-02-17) - Initial release

## Upcoming Features

We're planning to add:
- WebSocket interception support
- Request/response diffing tools
- Import/export for rules and scenarios
- gRPC traffic support
- Performance metrics and analytics
- Plugin system for custom processors

See our [GitHub Issues](https://github.com/wahyudotdev/glance/issues) for upcoming features and known bugs.
