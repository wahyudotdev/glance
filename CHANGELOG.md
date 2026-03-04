# Changelog

All notable changes to the Glance project will be documented in this file.

## [v0.2.3] - 2026-03-05
- **Terminal Interception**: Redesigned setup script with improved compatibility and feedback.
- **Local Interception**: Default behavior now unsets `NO_PROXY` to intercept `localhost` traffic by default.
- **Enhanced Tool Support**: Added CA certificate environment variables for `curl`, `git`, `pip`, and `npm` (`CURL_CA_BUNDLE`, `GIT_SSL_CAINFO`, etc.).
- **Port Standardization**: Standardized default ports across documentation and code: `15501` for Dashboard/API and `15500` for Proxy.
- **Improved Feedback**: The terminal setup script now provides real-time feedback and success logs during configuration.
- **Testing**: Added `InitTestDB` for reliable in-memory database testing and added an integration test for terminal interception.

## [v0.2.2] - 2026-02-24
- **Mock Rules Engine**: New interactive mock editor with JSON tree view and switchable mocking rules.
- **Docker Networking**: Enhanced host resolution using `host.docker.internal` for Docker Desktop environments.
- **Docker Auto-injection**: Automatic `ExtraHosts` injection during container recreation for better connectivity.
- **Documentation**: Migrated to VitePress with new branding, dark/light themes, and detailed installation guidelines.
- **Visuals**: Added dashboard screenshots, favicons, and social assets to documentation.
- **Project Hygiene**: Integrated GitHub Sponsors and refined multi-platform release pipeline.

## [v0.2.1] - 2026-02-23
- **Improved Host Resolution**: Enhanced Host IP discovery to prioritize `host.docker.internal` for Docker Desktop users.
- **ExtraHosts Injection**: Automatically included `host.docker.internal:host-gateway` in `ExtraHosts` during container recreation to ensure reliable host resolution across all platforms (Linux/macOS/Windows).

## [v0.2.0] - 2026-02-23
- **Docker Interception**: Introduced a robust "One-Click Intercept" for Docker containers via container recreation and environment injection.
- **Java Container Support**: Added specialized logic to automatically inject and trust the Glance CA certificate into the Java Trust Store (`cacerts`) inside Docker containers using `keytool`.
- **Universal Compatibility**: Replaced the `iptables` approach with standard `HTTP_PROXY`/`HTTPS_PROXY` and `JAVA_TOOL_OPTIONS` injection, ensuring compatibility across all container OS flavors (Alpine, Debian, RHEL) without requiring `NET_ADMIN` capabilities.
- **Docker UI**: New full-width Docker Containers card in the Integrations view with real-time status tracking, container renaming for backups, and seamless restoration.
- **API Stability**: Fixed Docker SDK compatibility issues for older daemon versions (API 1.41) and properly handled binary stream demultiplexing for command execution.
