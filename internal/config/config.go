package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	ProxyAddr  string `json:"proxy_addr"`
	APIAddr    string `json:"api_addr"`
	MCPAddr    string `json:"mcp_addr"`
	MCPEnabled bool   `json:"mcp_enabled"`
}

var (
	instance *Config
	once     sync.Once
	mu       sync.RWMutex
	path     string
)

func init() {
	home, _ := os.UserHomeDir()
	path = filepath.Join(home, ".agent-proxy-config.json")
}

func Get() *Config {
	once.Do(func() {
		instance = &Config{
			ProxyAddr:  ":8000",
			APIAddr:    ":8081",
			MCPAddr:    ":8082",
			MCPEnabled: false,
		}
		load()
	})
	mu.RLock()
	defer mu.RUnlock()
	return instance
}

func load() {
	data, err := os.ReadFile(path)
	if err == nil {
		json.Unmarshal(data, instance)
	}
}

func Save(c *Config) error {
	mu.Lock()
	defer mu.Unlock()
	instance = c
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
