// Package main provides a load testing script for the Glance.
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func main() {
	proxyAddr := flag.String("proxy", "http://localhost:8000", "Proxy address")
	targetURL := flag.String("target", "https://www.google.com", "Target URL to request")
	concurrency := flag.Int("c", 10, "Number of concurrent workers")
	totalReqs := flag.Int("n", 100, "Total number of requests")
	insecure := flag.Bool("k", true, "Allow insecure server connections when using SSL")
	flag.Parse()

	proxyURL, err := url.Parse(*proxyAddr)
	if err != nil {
		fmt.Printf("Invalid proxy URL: %v\n", err)
		return
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{
			// #nosec G402
			InsecureSkipVerify: *insecure,
		},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	fmt.Printf("Starting load test: %d requests, %d concurrency, via %s\n", *totalReqs, *concurrency, *proxyAddr)

	var wg sync.WaitGroup
	reqsPerWorker := *totalReqs / *concurrency

	start := time.Now()

	successCount := 0
	errorCount := 0
	var mu sync.Mutex

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < reqsPerWorker; j++ {
				resp, err := client.Get(*targetURL)
				if err != nil {
					mu.Lock()
					errorCount++
					mu.Unlock()
					fmt.Printf("Worker %d error: %v\n", workerID, err)
					continue
				}

				// Read body to ensure full request capture
				_, _ = io.ReadAll(resp.Body)
				_ = resp.Body.Close()

				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	fmt.Printf("\nLoad Test Results:\n")
	fmt.Printf("Total Time: %v\n", duration)
	fmt.Printf("Success:    %d\n", successCount)
	fmt.Printf("Errors:     %d\n", errorCount)
	fmt.Printf("Req/sec:    %.2f\n", float64(successCount)/duration.Seconds())
}
