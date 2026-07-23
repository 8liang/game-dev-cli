# CLAUDE.md — game-dev-cli

## Skills

见 `skills/index.md`。

## Installation

```bash
# go install（需 Go）
go install github.com/8liang/game-dev-cli@latest

# npx（无需 Go）
npx @8liang/game-dev-cli --help

# curl 一键安装
curl -fsSL https://github.com/8liang/game-dev-cli/releases/latest/download/install.sh | bash
```

## MCP 配置

在 AI 工具 MCP 配置中添加:

```json
{
  "mcpServers": {
    "game-dev-cli": {
      "command": "game-dev-cli",
      "args": ["mcp"]
    }
  }
}
```
