# Installation

Glance can be installed in multiple ways depending on your platform and preferences.

## Homebrew (macOS & Linux)

The easiest way to install Glance is via Homebrew:

```bash
brew tap wahyudotdev/tap
brew install glance
```

### Updating

To update to the latest version:

```bash
brew update
brew upgrade glance
```

## Direct Download

Download the latest pre-compiled binary for your platform from the [Releases](https://github.com/wahyudotdev/glance/releases/latest) page.

### macOS

```bash
# Download for macOS (ARM64)
curl -LO https://github.com/wahyudotdev/glance/releases/latest/download/glance_Darwin_arm64.tar.gz

# Extract
tar -xzf glance_Darwin_arm64.tar.gz

# Make executable
chmod +x glance

# Move to PATH (optional)
sudo mv glance /usr/local/bin/
```

### Linux

```bash
# Download for Linux (AMD64)
curl -LO https://github.com/wahyudotdev/glance/releases/latest/download/glance_Linux_x86_64.tar.gz

# Extract
tar -xzf glance_Linux_x86_64.tar.gz

# Make executable
chmod +x glance

# Move to PATH (optional)
sudo mv glance /usr/local/bin/
```

### Windows

Download the Windows binary from the [Releases](https://github.com/wahyudotdev/glance/releases/latest) page and extract it to a directory in your PATH.

## Build from Source

If you prefer to build from source:

### Prerequisites

- Go 1.24+
- Node.js & npm (for frontend development)
- Make

### Build Steps

```bash
# Clone the repository
git clone https://github.com/wahyudotdev/glance.git
cd glance

# Build the binary (includes frontend build)
make build

# Run
./glance
```

## Verify Installation

After installation, verify that Glance is working:

```bash
glance --version
```

You should see output similar to:

```
Glance v0.1.4
```

## Next Steps

- [Quick Start Guide](quickstart.md) - Get started with Glance
- [MCP Integration](mcp/) - Set up AI agent integration
- [Client Configuration](clients.md) - Configure your applications to use Glance
