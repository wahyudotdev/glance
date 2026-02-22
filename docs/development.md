# Development Guide

This guide will help you set up a development environment for contributing to Glance.

## Prerequisites

- **Go**: 1.24 or higher
- **Node.js**: 20 or higher
- **npm**: Latest version
- **Make**: For build automation
- **Git**: For version control

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/wahyudotdev/glance.git
cd glance
```

### 2. Install Dependencies

#### Backend (Go)

```bash
go mod download
```

#### Frontend (React)

```bash
cd web/dashboard
npm install
cd ../..
```

### 3. Build the Project

```bash
# Build everything (backend + frontend)
make build

# Or build separately
make build-frontend
make build-backend
```

### 4. Run in Development Mode

#### Backend Only

```bash
go run cmd/glance/main.go
```

#### Frontend Only

```bash
cd web/dashboard
npm run dev
```

For frontend development, you'll need the backend running on `http://localhost:15500` and `http://localhost:15501`.

## Project Structure

```
glance/
├── cmd/
│   └── glance/          # Main application entry point
├── internal/            # Internal packages
│   ├── proxy/          # MITM proxy implementation
│   ├── dashboard/      # Dashboard API handlers
│   ├── mcp/           # MCP server implementation
│   ├── storage/       # Database and persistence
│   └── ...
├── web/
│   └── dashboard/      # React frontend
│       ├── src/
│       ├── public/
│       └── package.json
├── scripts/            # Build and utility scripts
├── Makefile           # Build automation
└── go.mod
```

## Make Commands

The project uses a `Makefile` for common tasks:

```bash
# Build the complete binary
make build

# Run tests
make test

# Run linter
make lint

# Generate coverage report
make test-coverage

# Clean build artifacts
make clean

# Build frontend only
make build-frontend

# Build backend only
make build-backend
```

## Running Tests

### Backend Tests

```bash
# Run all tests
go test ./internal/...

# Run with coverage
go test -cover ./internal/...

# Run specific package
go test ./internal/proxy

# Verbose output
go test -v ./internal/...
```

### Frontend Tests

```bash
cd web/dashboard
npm test
```

## Code Style

### Go

We use `golangci-lint` for Go code:

```bash
# Run linter
make lint

# Auto-fix issues
golangci-lint run --fix
```

Follow the [Effective Go](https://go.dev/doc/effective_go) guidelines.

### TypeScript/React

We use ESLint for TypeScript:

```bash
cd web/dashboard
npm run lint

# Auto-fix
npm run lint:fix
```

## Database

Glance uses SQLite with the following characteristics:

- **Location**: `~/.glance.db`
- **Mode**: Write-Ahead Logging (WAL)
- **Caching**: Write-Behind for performance

### Schema Migrations

Currently, schema is managed manually. When updating the database schema:

1. Update the initialization code in `internal/storage/db.go`
2. Test with a fresh database
3. Consider backward compatibility

## Adding New Features

### Backend Feature

1. Create new package in `internal/` if needed
2. Implement the feature with tests
3. Add API endpoints in `internal/dashboard/`
4. Update MCP tools if relevant

### Frontend Feature

1. Create components in `web/dashboard/src/components/`
2. Add hooks in `web/dashboard/src/hooks/` if needed
3. Update routes if adding new pages
4. Add corresponding API calls

### MCP Tool

1. Define tool in `internal/mcp/tools.go`
2. Implement handler logic
3. Add tests
4. Update documentation in `docs/mcp/reference.md`

## Debugging

### Backend Debugging

Use VS Code with the Go extension:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Glance",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/glance"
    }
  ]
}
```

Or use `delve` directly:

```bash
dlv debug cmd/glance/main.go
```

### Frontend Debugging

Use browser DevTools and the React DevTools extension.

## Performance Profiling

### Go Profiling

```bash
# CPU profile
go test -cpuprofile=cpu.prof ./internal/...
go tool pprof cpu.prof

# Memory profile
go test -memprofile=mem.prof ./internal/...
go tool pprof mem.prof
```

### Frontend Profiling

Use React DevTools Profiler tab to analyze component render performance.

## Release Process

Releases are automated via GoReleaser:

1. Tag a new version:
   ```bash
   git tag -a v0.1.5 -m "Release v0.1.5"
   git push origin v0.1.5
   ```

2. GitHub Actions will automatically:
   - Build binaries for all platforms
   - Create a GitHub Release
   - Update Homebrew tap

## CI/CD

We use GitHub Actions for CI/CD:

- `.github/workflows/ci.yml` - Tests, linting, coverage
- `.github/workflows/release.yml` - Release automation
- `.github/workflows/gh-pages.yml` - Documentation deployment

## Contributing Guidelines

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Pull Request Checklist

- [ ] Tests pass locally
- [ ] Linter passes
- [ ] Added tests for new features
- [ ] Updated documentation
- [ ] Followed code style guidelines
- [ ] No breaking changes (or clearly documented)

## Getting Help

- 📖 Read the [Architecture Guide](architecture.md)
- 💬 [Open a discussion](https://github.com/wahyudotdev/glance/discussions)
- 🐛 [Report issues](https://github.com/wahyudotdev/glance/issues)

## License

By contributing to Glance, you agree that your contributions will be licensed under the MIT License.
