---
layout: home

hero:
  name: "Glance"
  text: "Let Your AI Understand Every Request at a Glance."
  tagline: "The specialized MITM proxy for AI Agents with native MCP integration."
  image:
    src: /icon.svg
    alt: Glance Logo
  actions:
    - theme: brand
      text: Get Started
      link: /quickstart
    - theme: alt
      text: View on GitHub
      link: https://github.com/wahyudotdev/glance

features:
  - title: 🔍 Real-time Traffic
    details: Stream captured requests and responses as they happen via WebSockets.
  - title: 🛠️ Mocking Engine
    details: Define rules to intercept and modify traffic on the fly with a powerful UI.
  - title: 🤖 Native MCP
    details: Expose traffic and tools to AI agents like Claude Desktop natively.
  - title: 📱 Multi-Platform
    details: Deep integration with Java/JVM, Docker, Android, and Chromium.

---

<div style="margin-top: 2rem; max-width: 1152px; margin-left: auto; margin-right: auto;">

::: code-group

```bash [macOS]
brew tap wahyudotdev/tap
brew install glance
```

```bash [Linux]
# Download
curl -LO https://github.com/wahyudotdev/glance/releases/latest/download/glance_Linux_x86_64.tar.gz

# Extract
tar -xzf glance_Linux_x86_64.tar.gz

# Move to path
sudo mv glance /usr/local/bin/
```

```powershell [Windows]
# Download the latest .exe from our Releases page:
https://github.com/wahyudotdev/glance/releases/latest
```

:::

</div>
