package main

import (
	"flag"
	"log"

	"agent-proxy/internal/api"
	"agent-proxy/internal/proxy"
)

func main() {
	proxyAddr := flag.String("proxy-addr", ":8000", "proxy listen address")
	apiAddr := flag.String("api-addr", ":8081", "api/dashboard listen address")
	flag.Parse()

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
	if err := apiServer.Listen(*apiAddr); err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}
}
