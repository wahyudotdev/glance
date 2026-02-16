package main

import (
	"flag"
	"log"

	"agent-proxy/internal/api"
	"agent-proxy/internal/config"
	"agent-proxy/internal/mcp"
	"agent-proxy/internal/proxy"
)

func main() {
	cfg := config.Get()

	proxyAddr := flag.String("proxy-addr", cfg.ProxyAddr, "proxy listen address")
	apiAddr := flag.String("api-addr", cfg.APIAddr, "api/dashboard listen address")
	mcpAddr := flag.String("mcp-addr", cfg.MCPAddr, "mcp server listen address (SSE)")
	mcpMode := flag.Bool("mcp", cfg.MCPEnabled, "run as MCP server")
	flag.Parse()

	// Update config with flags if they were provided (flags override saved config)
	if *proxyAddr != cfg.ProxyAddr || *apiAddr != cfg.APIAddr || *mcpAddr != cfg.MCPAddr || *mcpMode != cfg.MCPEnabled {
		cfg.ProxyAddr = *proxyAddr
		cfg.APIAddr = *apiAddr
		cfg.MCPAddr = *mcpAddr
		cfg.MCPEnabled = *mcpMode
		config.Save(cfg)
	}

	// Check for Java Agent injection mode (used internally)
	if len(flag.Args()) > 0 && flag.Args()[0] == "inject-agent" {
		// This path is just for building/attaching, usually handled by library functions
		// but we can expose it if needed. For now, we rely on the API.
		return
	}

	p := proxy.NewProxy(*proxyAddr)

	actualProxyAddr, err := p.Start()
	if err != nil {
		log.Fatalf("Failed to start proxy: %v", err)
	}

	// Start API Server
	apiServer := api.NewAPIServer(p.Store, actualProxyAddr)
	apiServer.RegisterRoutes()

	log.Printf("API server starting on %s", *apiAddr)
	log.Printf("Dashboard available at http://localhost%s", *apiAddr)

	// Start MCP Server if requested
	if *mcpMode {
		go func() {
			mcpServer := mcp.NewMCPServer(p.Store, p.Engine, actualProxyAddr)
			log.Printf("MCP server (SSE) starting on %s", *mcpAddr)
			if err := mcpServer.ServeSSE(*mcpAddr); err != nil {
				log.Fatalf("MCP Server failed: %v", err)
			}
		}()
	}

	if err := apiServer.Listen(*apiAddr); err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}
}
