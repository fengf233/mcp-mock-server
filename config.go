package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"gopkg.in/yaml.v3"
)

// Config 表示整个mock服务器的配置
type Config struct {
	Server    ServerConfig     `yaml:"server"`
	Manifest  ManifestConfig   `yaml:"manifest"`
	Tools     []ToolConfig     `yaml:"tools"`
	Prompts   []PromptConfig   `yaml:"prompts"`
	Resources []ResourceConfig `yaml:"resources"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	BasePath  string `yaml:"base_path"`
	Transport string `yaml:"transport"`
}

// ManifestConfig 清单配置
type ManifestConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

// ToolConfig 工具配置
type ToolConfig struct {
	Name         string          `yaml:"name"`
	Description  string          `yaml:"description"`
	Parameters   []ToolParameter `yaml:"parameters"`
	MockTemplate string          `yaml:"mock_template"`
}

// ToolParameter 工具参数配置
type ToolParameter struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
}

// PromptConfig 提示配置
type PromptConfig struct {
	Name         string           `yaml:"name"`
	Description  string           `yaml:"description"`
	Arguments    []PromptArgument `yaml:"arguments"`
	MockTemplate string           `yaml:"mock_template"`
}

// PromptArgument 提示参数
type PromptArgument struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
}

// ResourceConfig 资源配置
type ResourceConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	MIMEType    string `yaml:"mime_type"`
	URI         string `yaml:"uri"`
	MockContent string `yaml:"mock_content"`
}

// LoadConfigFromFile 从YAML文件加载配置
func LoadConfigFromFile(filename string) (*Config, error) {
	yamlContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(yamlContent, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}
	return &config, nil
}

// LoadConfig 从YAML内容加载配置
func LoadConfig(yamlContent []byte) (*Config, error) {
	var config Config
	err := yaml.Unmarshal(yamlContent, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}
	return &config, nil
}

// TemplateProcessor 模板处理器，用于替换变量
type TemplateProcessor struct{}

// Process 处理模板，替换 {{variable}} 为实际值
func (tp *TemplateProcessor) Process(template string, variables map[string]string) string {
	result := template
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// ExtractVariables 从请求中提取变量
func ExtractToolVariables(req mcp.CallToolRequest) map[string]string {
	variables := make(map[string]string)

	// 检查参数是否为map类型
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		for key, value := range args {
			if strVal, ok := value.(string); ok {
				variables[key] = strVal
			}
		}
	}

	return variables
}

// ExtractPromptVariables 从请求中提取提示变量
func ExtractPromptVariables(req mcp.GetPromptRequest) map[string]string {
	variables := make(map[string]string)

	// 检查参数是否为map类型
	for key, value := range req.Params.Arguments {
		variables[key] = value
	}

	return variables
}
