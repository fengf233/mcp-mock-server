# MCP 服务器

基于 Model Context Protocol (MCP) 的模拟服务器，用于开发和测试 MCP 客户端。

## 功能

- **工具调用**：提供模拟工具，支持参数化模板
- **提示词管理**：动态生成提示词内容
- **资源访问**：提供模拟资源数据
- **多传输支持**：支持 StreamableHTTP 和 SSE 传输

## 快速开始

### 安装依赖
```bash
go mod tidy
```

### 配置服务器
编辑 `mock.yaml` 文件配置服务器参数：
```yaml
server:
  host: "0.0.0.0"
  port: 9999
  transport: "StreamableHTTP"  # 或 "sse"
```

### 运行服务器
```bash
# 编译
go build

# 运行
./mcp_server
```

服务器将在 `http://localhost:9999/mcp` 启动。

## 配置说明

### 工具配置
在 `mock.yaml` 中定义工具：
```yaml
tools:
  - name: "get_weather"
    description: "Get current weather for a city"
    parameters:
      - name: "city"
        type: string #只支持字符串
        required: true
    mock_template: "The current weather in {{city}} is sunny, 25°C."
```

### 提示词配置
```yaml
prompts:
  - name: "summarize_article"
    description: "Summarize a given article"
    arguments:
      - name: "article"
        required: true
    mock_template: "Summarize: {{article}}"
```

### 资源配置
```yaml
resources:
  - name: "user_profile"
    mime_type: "application/json"
    uri: "internal://resources/user_profile.json"
    mock_content: "{\"user_id\": \"{{user_id}}\"}"
```

## 项目结构

```
mcp_server/
├── main.go          # 服务器主程序
├── config.go        # 配置加载和解析
├── registry.go       # 工具、提示词、资源注册
├── mock.yaml        # 服务器配置文件
├── go.mod           # 模块定义
└── go.sum           # 依赖校验
```

## 依赖

- `github.com/mark3labs/mcp-go/server`：MCP 服务器库
- `github.com/mark3labs/mcp-go/mcp`：MCP 协议定义
- `gopkg.in/yaml.v3`：YAML 配置文件解析

## 使用示例

1. 启动服务器：`./mcp_server`
2. 使用 MCP 客户端连接：`mcp_client -type http -url http://localhost:9999/mcp`
3. 在客户端中调用工具、获取提示词、查看资源