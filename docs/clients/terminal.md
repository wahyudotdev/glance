# Terminal/CLI Configuration

Configure command-line tools and shell environments to route traffic through Glance.

## Quick Setup

### One-Liner (Recommended)

The easiest way to configure your current shell session:

```bash
eval "$(curl -s http://localhost:15501/setup)"
```

This automatically exports the necessary environment variables for your session.

### Manual Environment Variables

Set proxy variables manually:

```bash
export HTTP_PROXY=http://localhost:15500
export HTTPS_PROXY=http://localhost:15500
export NO_PROXY=localhost,127.0.0.1
```

## Persistent Configuration

To make the proxy configuration permanent:

### Bash

Add to `~/.bashrc` or `~/.bash_profile`:

```bash
# Glance Proxy
export HTTP_PROXY=http://localhost:15500
export HTTPS_PROXY=http://localhost:15500
export NO_PROXY=localhost,127.0.0.1
```

Reload:
```bash
source ~/.bashrc
```

### Zsh

Add to `~/.zshrc`:

```bash
# Glance Proxy
export HTTP_PROXY=http://localhost:15500
export HTTPS_PROXY=http://localhost:15500
export NO_PROXY=localhost,127.0.0.1
```

Reload:
```bash
source ~/.zshrc
```

### Fish

Add to `~/.config/fish/config.fish`:

```fish
# Glance Proxy
set -x HTTP_PROXY http://localhost:15500
set -x HTTPS_PROXY http://localhost:15500
set -x NO_PROXY localhost,127.0.0.1
```

Reload:
```fish
source ~/.config/fish/config.fish
```

## Tool-Specific Configuration

### cURL

```bash
# Use environment variables (automatic)
curl https://api.github.com/users

# Or specify proxy explicitly
curl -x http://localhost:15500 https://api.github.com/users

# Ignore certificate errors (development only)
curl -k -x http://localhost:15500 https://api.github.com/users
```

### wget

```bash
# Use environment variables
wget https://example.com

# Or specify in command
wget -e use_proxy=yes \
     -e http_proxy=http://localhost:15500 \
     -e https_proxy=http://localhost:15500 \
     https://example.com

# Ignore certificate errors
wget --no-check-certificate https://example.com
```

### HTTPie

```bash
# HTTPie respects HTTP_PROXY and HTTPS_PROXY automatically
http https://api.github.com/users

# Or specify explicitly
http --proxy=http:http://localhost:15500 \
     --proxy=https:http://localhost:15500 \
     https://api.github.com/users

# Ignore SSL
http --verify=no https://api.github.com/users
```

### Git

```bash
# Configure globally
git config --global http.proxy http://localhost:15500
git config --global https.proxy http://localhost:15500

# Per repository
git config http.proxy http://localhost:15500
git config https.proxy http://localhost:15500

# Unset proxy
git config --global --unset http.proxy
git config --global --unset https.proxy

# Ignore SSL (not recommended)
git config --global http.sslVerify false
```

### npm

```bash
# Set proxy
npm config set proxy http://localhost:15500
npm config set https-proxy http://localhost:15500

# Remove proxy
npm config delete proxy
npm config delete https-proxy

# Ignore SSL (development only)
npm config set strict-ssl false
```

### yarn

```bash
# Set proxy
yarn config set proxy http://localhost:15500
yarn config set https-proxy http://localhost:15500

# Remove proxy
yarn config delete proxy
yarn config delete https-proxy
```

### pip (Python)

```bash
# Use environment variables
pip install requests

# Or specify in command
pip install --proxy http://localhost:15500 requests

# Via config file ~/.config/pip/pip.conf
[global]
proxy = http://localhost:15500
```

### Docker

```bash
# Set proxy for docker daemon
# Edit /etc/docker/daemon.json
{
  "proxies": {
    "http-proxy": "http://localhost:15500",
    "https-proxy": "http://localhost:15500",
    "no-proxy": "localhost,127.0.0.1"
  }
}

# Restart Docker
sudo systemctl restart docker

# For docker build
docker build \
  --build-arg HTTP_PROXY=http://localhost:15500 \
  --build-arg HTTPS_PROXY=http://localhost:15500 \
  -t myimage .
```

## Certificate Trust

For HTTPS interception, you may need to trust Glance's CA certificate.

### System-Wide Certificate Trust

#### macOS

```bash
# Download certificate
curl http://localhost:15501/ca.crt -o glance-ca.crt

# Add to system keychain
sudo security add-trusted-cert \
  -d -r trustRoot \
  -k /Library/Keychains/System.keychain \
  glance-ca.crt
```

#### Linux (Debian/Ubuntu)

```bash
# Download certificate
curl http://localhost:15501/ca.crt -o glance-ca.crt

# Copy to CA directory
sudo cp glance-ca.crt /usr/local/share/ca-certificates/

# Update certificates
sudo update-ca-certificates
```

#### Linux (RedHat/CentOS)

```bash
# Download certificate
curl http://localhost:15501/ca.crt -o glance-ca.crt

# Copy to CA directory
sudo cp glance-ca.crt /etc/pki/ca-trust/source/anchors/

# Update certificates
sudo update-ca-trust
```

### Tool-Specific Certificate

#### cURL

```bash
# Use certificate file
curl --cacert glance-ca.crt https://example.com

# Or ignore (development only)
curl -k https://example.com
```

#### Python Requests

```bash
# Set certificate bundle
export REQUESTS_CA_BUNDLE=/path/to/glance-ca.crt

# Or in code
import requests
requests.get('https://example.com', verify='/path/to/glance-ca.crt')
```

#### Node.js

```bash
# Set CA bundle
export NODE_EXTRA_CA_CERTS=/path/to/glance-ca.crt

# Or disable SSL (not recommended)
export NODE_TLS_REJECT_UNAUTHORIZED=0
```

## Testing

Verify your configuration:

```bash
# Make a request
curl https://api.github.com/users

# Should see:
# 1. Response from GitHub
# 2. Request in Glance dashboard
```

Test with httpbin:

```bash
# Check proxy is working
curl http://httpbin.org/get

# Response should show your proxy
```

## Shell Functions

Create convenient functions in your shell:

```bash
# Add to ~/.bashrc or ~/.zshrc

# Enable Glance proxy
glance-on() {
  export HTTP_PROXY=http://localhost:15500
  export HTTPS_PROXY=http://localhost:15500
  export NO_PROXY=localhost,127.0.0.1
  echo "✓ Glance proxy enabled"
}

# Disable Glance proxy
glance-off() {
  unset HTTP_PROXY
  unset HTTPS_PROXY
  unset NO_PROXY
  echo "✓ Glance proxy disabled"
}

# Check proxy status
glance-status() {
  if [ -n "$HTTP_PROXY" ]; then
    echo "✓ Proxy enabled: $HTTP_PROXY"
  else
    echo "✗ Proxy disabled"
  fi
}
```

Usage:

```bash
glance-on    # Enable proxy
glance-off   # Disable proxy
glance-status # Check status
```

## Environment-Specific Configuration

### Development

```bash
# ~/.glance_dev
export HTTP_PROXY=http://localhost:15500
export HTTPS_PROXY=http://localhost:15500
export NODE_TLS_REJECT_UNAUTHORIZED=0  # Only for dev!

# Load when needed
source ~/.glance_dev
```

### CI/CD

```yaml
# GitHub Actions
env:
  HTTP_PROXY: http://localhost:15500
  HTTPS_PROXY: http://localhost:15500

# GitLab CI
variables:
  HTTP_PROXY: "http://localhost:15500"
  HTTPS_PROXY: "http://localhost:15500"
```

## Troubleshooting

### Proxy Not Working

Check if variables are set:

```bash
echo $HTTP_PROXY
echo $HTTPS_PROXY
```

Test with verbose curl:

```bash
curl -v -x http://localhost:15500 https://example.com
```

### Certificate Errors

```bash
# Verify certificate exists
ls -la /usr/local/share/ca-certificates/

# Check certificate is valid
openssl x509 -in glance-ca.crt -text -noout

# Test with specific cert
curl --cacert glance-ca.crt https://example.com
```

### Some Tools Don't Use Proxy

Not all tools respect `HTTP_PROXY` environment variables:
- Check tool documentation
- Look for tool-specific proxy configuration
- Consider using system-wide proxy

## Best Practices

1. **Use Functions**: Create shell functions for easy on/off
2. **NO_PROXY**: Always set `NO_PROXY` to avoid loops
3. **Temporary**: Use environment variables for temporary sessions
4. **Security**: Never commit proxy credentials to git
5. **Documentation**: Document proxy setup in project README

## Next Steps

- [Client Configuration](/clients.md) - Other platforms
- [MCP Integration](/mcp/) - Analyze CLI traffic with AI
- [Troubleshooting](/troubleshooting.md) - Common issues
