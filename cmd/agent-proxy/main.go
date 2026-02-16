package main

import (
	"flag"
	"log"

	"agent-proxy/internal/api"
	"agent-proxy/internal/proxy"
)

func main() {
	proxyAddr := flag.String("proxy-addr", ":8080", "proxy listen address")
	apiAddr := flag.String("api-addr", ":8081", "api/dashboard listen address")
	flag.Parse()

	p := proxy.NewProxy(*proxyAddr)

	// Start API Server
	apiServer := api.NewAPIServer(p.Store)
	apiServer.RegisterRoutes()

	go func() {
		log.Printf("API server starting on %s", *apiAddr)
		if err := apiServer.Listen(*apiAddr); err != nil {
			log.Fatalf("Failed to start API server: %v", err)
		}
	}()

	if err := p.Start(); err != nil {
		log.Fatalf("Failed to start proxy: %v", err)
	}
}
