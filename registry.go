package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ToolRegistry 工具注册器
type ToolRegistry struct {
	config            *Config
	templateProcessor *TemplateProcessor
}

// NewToolRegistry 创建新的工具注册器
func NewToolRegistry(config *Config) *ToolRegistry {
	return &ToolRegistry{
		config:            config,
		templateProcessor: &TemplateProcessor{},
	}
}

// RegisterTools 根据配置注册所有工具
func (tr *ToolRegistry) RegisterTools(s *server.MCPServer) error {
	for _, toolConfig := range tr.config.Tools {
		err := tr.registerTool(s, toolConfig)
		if err != nil {
			return fmt.Errorf("failed to register tool %s: %w", toolConfig.Name, err)
		}
		log.Printf("Registered tool: %s", toolConfig.Name)
	}
	return nil
}

// registerTool 注册单个工具
func (tr *ToolRegistry) registerTool(s *server.MCPServer, toolConfig ToolConfig) error {
	// 创建工具选项切片
	opts := []mcp.ToolOption{
		mcp.WithDescription(toolConfig.Description),
	}

	// 为每个参数添加 WithString 配置
	for _, param := range toolConfig.Parameters {
		paramOpts := []mcp.PropertyOption{}
		if param.Required {
			paramOpts = append(paramOpts, mcp.Required())
		}

		opts = append(opts, mcp.WithString(param.Name, paramOpts...))
	}

	// 创建工具定义
	tool := mcp.NewTool(toolConfig.Name, opts...)

	// 注册工具处理函数
	handler := tr.createToolHandler(toolConfig)
	s.AddTool(tool, handler)

	return nil
}

// createToolHandler 创建工具处理函数
func (tr *ToolRegistry) createToolHandler(toolConfig ToolConfig) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 提取变量
		variables := ExtractToolVariables(req)

		// 处理模板
		resultText := tr.templateProcessor.Process(toolConfig.MockTemplate, variables)

		return mcp.NewToolResultText(resultText), nil
	}
}

// RegisterPrompts 注册提示
func (tr *ToolRegistry) RegisterPrompts(s *server.MCPServer) error {
	for _, promptConfig := range tr.config.Prompts {
		err := tr.registerPrompt(s, promptConfig)
		if err != nil {
			return fmt.Errorf("failed to register prompt %s: %w", promptConfig.Name, err)
		}
		log.Printf("Registered prompt: %s", promptConfig.Name)
	}
	return nil
}

// registerPrompt 注册单个提示
func (tr *ToolRegistry) registerPrompt(s *server.MCPServer, promptConfig PromptConfig) error {
	// 创建提示选项切片
	opts := []mcp.PromptOption{
		mcp.WithPromptDescription(promptConfig.Description),
	}

	// 为每个参数添加 WithArgument 配置
	for _, arg := range promptConfig.Arguments {
		argOpts := []mcp.ArgumentOption{}
		if arg.Required {
			argOpts = append(argOpts, mcp.RequiredArgument())
		}
		if arg.Description != "" {
			argOpts = append(argOpts, mcp.ArgumentDescription(arg.Description))
		}

		opts = append(opts, mcp.WithArgument(arg.Name, argOpts...))
	}

	// 创建提示定义
	prompt := mcp.NewPrompt(promptConfig.Name, opts...)

	// 注册提示处理函数
	handler := tr.createPromptHandler(promptConfig)
	s.AddPrompt(prompt, handler)

	return nil
}

// createPromptHandler 创建提示处理函数
func (tr *ToolRegistry) createPromptHandler(promptConfig PromptConfig) server.PromptHandlerFunc {
	return func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		// 提取变量
		variables := ExtractPromptVariables(req)
		// 处理模板
		processedTemplate := tr.templateProcessor.Process(promptConfig.MockTemplate, variables)

		// 创建提示结果
		result := &mcp.GetPromptResult{
			Description: promptConfig.Description,
			Messages: []mcp.PromptMessage{
				{
					Role: "user",
					Content: mcp.TextContent{
						Type: "text",
						Text: processedTemplate,
					},
				},
			},
		}

		return result, nil
	}
}

// RegisterResources 注册资源
func (tr *ToolRegistry) RegisterResources(s *server.MCPServer) error {
	for _, resourceConfig := range tr.config.Resources {
		err := tr.registerResource(s, resourceConfig)
		if err != nil {
			return fmt.Errorf("failed to register resource %s: %w", resourceConfig.Name, err)
		}
		log.Printf("Registered resource: %s", resourceConfig.Name)
	}
	return nil
}

// registerResource 注册单个资源
func (tr *ToolRegistry) registerResource(s *server.MCPServer, resourceConfig ResourceConfig) error {
	// 创建资源定义
	resource := mcp.NewResource(
		resourceConfig.URI,
		resourceConfig.Name,
		mcp.WithResourceDescription(resourceConfig.Description),
		mcp.WithMIMEType(resourceConfig.MIMEType),
	)

	// 注册资源处理函数
	handler := tr.createResourceHandler(resourceConfig)
	s.AddResource(resource, handler)

	return nil
}

// createResourceHandler 创建资源处理函数
func (tr *ToolRegistry) createResourceHandler(resourceConfig ResourceConfig) func(context.Context, mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// 处理内容
		processedContent := tr.templateProcessor.Process(resourceConfig.MockContent, nil)

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      req.Params.URI,
				MIMEType: resourceConfig.MIMEType,
				Text:     processedContent,
			},
		}, nil
	}
}
