# CLAUDE.md — game-dev-cli

## Skills

### proto-gen
Generate Go/TypeScript code from .proto files.

用法:
```bash
game-dev-cli proto gen <proto-dir> \
  --go-out <dir> --ts-out <dir> --go-module <module> \
  [--plugin <spec>] [--include <path>] [--inject-tag] [--no-recursive]
```

参数:
- `--go-out` Go 代码输出目录
- `--ts-out` TypeScript 输出目录（等价于 `--plugin es,out=<dir>`）
- `--go-module` Go module 路径
- `--include` protoc `-I` 附加包含路径（可重复）
- `--plugin` protoc 插件，可重复；格式: `name[,binary=<path>][,out=<dir>][,module=<mod>][,opt=<k=v>]`
- `--inject-tag` 编译后注入 struct tag（需 protoc-go-inject-tag）
- `--no-recursive` 只扫描顶层 .proto，不递归子目录

### excel-gen
Generate JSON + Go structs + TypeScript interfaces from Excel files.

用法:
```bash
game-dev-cli excel gen <excel-dir> \
  --json-out <dir> --go-out <dir> --go-package <pkg> --ts-out <dir>
```

参数:
- `--json-out` JSON 输出目录（默认: `<excel-dir>/json`）
- `--go-out` Go struct 输出目录（需配合 `--go-package`）
- `--go-package` Go 包名
- `--ts-out` TypeScript interface 输出目录

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
