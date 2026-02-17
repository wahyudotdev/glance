// Package main is the entry point for the Glance application.
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"glance/internal/apiserver"
	"glance/internal/config"
	"glance/internal/db"
	"glance/internal/interceptor"
	"glance/internal/mcp"
	"glance/internal/proxy"
	"glance/internal/repository"
	"glance/internal/rules"
)

const (
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
)

func printBanner() {
	banner := `
   ______ _                             
  / ____/| |                          
 | |  __ | |  ____ _  _ __    ____  ___ 
 | | |_ || | / _' | || '_ \  / __|/ _ \
 | |__| || || (_| | || | | || (__|  __/
  \______|_| \__,_||_||_| |_| \___|\___|
                                        `
	fmt.Printf("%s%s%s\n", colorCyan, banner, colorReset)
	fmt.Printf("%s%s  Let Your AI Understand Every Request at a Glance.%s\n\n", colorBold, colorBlue, colorReset)
}

func formatAddr(addr string) string {
	if strings.HasPrefix(addr, "[::]") {
		return "localhost" + strings.TrimPrefix(addr, "[::]")
	}
	if strings.HasPrefix(addr, ":") {
		return "localhost" + addr
	}
	return addr
}

func main() {

	db.Init()

	// Initialize repositories

	configRepo := repository.NewSQLiteConfigRepository(db.DB)

	trafficRepo := repository.NewSQLiteTrafficRepository(db.DB)

	ruleRepo := repository.NewSQLiteRuleRepository(db.DB)

	scenarioRepo := repository.NewSQLiteScenarioRepository(db.DB)

	config.Init(configRepo)

	cfg := config.Get()

	proxyAddr := flag.String("proxy-addr", cfg.ProxyAddr, "proxy listen address")

	apiAddr := flag.String("api-addr", cfg.APIAddr, "api/dashboard listen address")

	mcpMode := flag.Bool("mcp", cfg.MCPEnabled, "run as MCP server")

	versionFlag := flag.Bool("version", false, "display version information")

	flag.Parse()

	if *versionFlag {

		fmt.Printf("Glance version %s\n", config.Version)

		return

	}

	printBanner()

	// Update config with flags if they were provided (flags override saved config)
	if *proxyAddr != cfg.ProxyAddr || *apiAddr != cfg.APIAddr || *mcpMode != cfg.MCPEnabled {
		cfg.ProxyAddr = *proxyAddr
		cfg.APIAddr = *apiAddr
		cfg.MCPEnabled = *mcpMode
		if err := config.Save(cfg); err != nil {
			log.Printf("Warning: Failed to save updated config: %v", err)
		}
	}
	// Check for Java Agent injection mode (used internally)
	if len(flag.Args()) > 0 && flag.Args()[0] == "inject-agent" {
		return
	}

	store := interceptor.NewTrafficStore(trafficRepo)
	engine := rules.NewEngine(ruleRepo)
	p := proxy.NewProxyWithRepositories(*proxyAddr, store, engine)

	actualProxyAddr, err := p.Start()
	if err != nil {
		log.Fatalf("Failed to start proxy: %v", err)
	}
	fmt.Printf("%s[✓]%s Proxy server running on %s%s%s\n", colorGreen, colorReset, colorBold, formatAddr(actualProxyAddr), colorReset)

	// Initialize MCP Server if requested
	var mcpServer *mcp.Server
	if *mcpMode {
		mcpServer = mcp.NewServer(p.Store, p.Engine, actualProxyAddr, scenarioRepo)
	}

	// Start API Server
	apiServer := apiserver.NewServer(p.Store, p, actualProxyAddr, mcpServer, scenarioRepo)
	apiServer.RegisterRoutes()

	// Connect Proxy to WebSocket Hub
	p.OnEntry = apiServer.Hub.Broadcast
	p.OnIntercept = apiServer.BroadcastIntercept

	go func() {
		actualAPIAddr, err := apiServer.Listen(*apiAddr)
		if err != nil {
			log.Fatalf("Failed to start API server: %v", err)
		}
		_ = actualAPIAddr
	}()

	// Wait a tiny bit for the goroutine to potentially log its fallback
	time.Sleep(100 * time.Millisecond)

	select {}
}
