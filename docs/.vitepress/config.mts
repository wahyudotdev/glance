import { defineConfig } from 'vitepress'

export default defineConfig({
  title: "Glance",
  description: "The specialized MITM proxy for AI Agents with MCP integration",
  base: '/glance/',
  head: [['link', { rel: 'icon', href: '/glance/icon.svg' }]],
  
  themeConfig: {
    logo: '/glance/icon.svg',
    
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Quick Start', link: '/quickstart' },
      { text: 'GitHub', link: 'https://github.com/wahyudotdev/glance' }
    ],

    sidebar: [
      {
        text: 'Getting Started',
        items: [
          { text: 'Installation', link: '/installation' },
          { text: 'Quick Start', link: '/quickstart' }
        ]
      },
      {
        text: 'Features',
        items: [
          { text: 'Traffic Inspection', link: '/features/traffic-inspection' },
          { text: 'Mocking & Breakpoints', link: '/features/mocking' },
          { text: 'Scenario Recording', link: '/features/scenarios' }
        ]
      },
      {
        text: 'Integration',
        items: [
          { text: 'Supported Clients', link: '/clients' },
          { text: 'Java/JVM', link: '/clients/java' },
          { text: 'Docker', link: '/clients/docker' },
          { text: 'Android', link: '/clients/android' },
          { text: 'Chromium', link: '/clients/chromium' },
          { text: 'Terminal', link: '/clients/terminal' }
        ]
      },
      {
        text: 'MCP',
        items: [
          { text: 'Overview', link: '/mcp/' },
          { text: 'Tools & Resources', link: '/mcp/reference' }
        ]
      },
      {
        text: 'Reference',
        items: [
          { text: 'API Reference', link: '/api' },
          { text: 'Configuration', link: '/configuration' }
        ]
      },
      {
        text: 'Development',
        items: [
          { text: 'Development Guide', link: '/development' },
          { text: 'Architecture', link: '/architecture' },
          { text: 'Contributing', link: '/contributing' }
        ]
      },
      {
        text: 'Help',
        items: [
          { text: 'FAQ', link: '/faq' },
          { text: 'Troubleshooting', link: '/troubleshooting' },
          { text: 'Changelog', link: '/changelog' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/wahyudotdev/glance' }
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2026 wahyudotdev'
    },

    search: {
      provider: 'local'
    }
  }
})
