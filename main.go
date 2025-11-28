package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// 加载配置文件
	config, err := LoadConfigFromFile("mock.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Loaded config: %s v%s", config.Manifest.Name, config.Manifest.Version)

	// 创建MCP服务器
	s := server.NewMCPServer(config.Manifest.Name, config.Manifest.Version,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	// 创建工具注册器并注册所有组件
	registry := NewToolRegistry(config)

	// 注册工具
	if err := registry.RegisterTools(s); err != nil {
		log.Fatalf("Failed to register tools: %v", err)
	}

	// 注册提示
	if err := registry.RegisterPrompts(s); err != nil {
		log.Fatalf("Failed to register prompts: %v", err)
	}

	// 注册资源
	if err := registry.RegisterResources(s); err != nil {
		log.Fatalf("Failed to register resources: %v", err)
	}

	// 根据配置选择传输方式
	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)

	log.Printf("Starting MCP server on %s with %s transport", addr, config.Server.Transport)

	switch config.Server.Transport {
	case "sse":
		// 启动SSE服务器
		sseServer := server.NewStreamableHTTPServer(s)
		if err := sseServer.Start(addr); err != nil {
			log.Fatal(err)
		}
	case "StreamableHTTP":
		fallthrough
	default:
		// 启动StreamableHTTP服务器
		httpServer := server.NewStreamableHTTPServer(s)
		if err := httpServer.Start(addr); err != nil {
			log.Fatal(err)
		}
	}
}

// 保留一些辅助函数，可能在其他地方使用
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func generateID() string {
	// Placeholder implementation
	return fmt.Sprintf("user_%d", time.Now().UnixNano())
}
