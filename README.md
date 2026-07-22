# game-dev-cli

AI vibe coding 工具集。基于 `.proto` 文件生成代码、从 Excel 文件导出数据。

## 安装

```bash
# 方式一：go install（需 Go）
go install github.com/8liang/game-dev-cli@latest

# 方式二：npx（无需 Go，自动下载二进制）
npx @8liang/game-dev-cli --help

# 方式三：curl 一键安装
curl -fsSL https://github.com/8liang/game-dev-cli/releases/latest/download/install.sh | bash
```

## proto gen — 从 .proto 文件生成代码

读取指定目录的 `.proto` 文件，调用 `protoc` 编译生成 Go 和/或 TypeScript 代码。

**要求：** `protoc` 已安装。TS 生成需安装 `protoc-gen-es`。

```bash
game-dev-cli proto gen <proto-dir> \
  --go-out ./gen/go \
  --ts-out ./gen/ts \
  --go-module github.com/user/project
```

参数：

| flag | 说明 |
|------|------|
| `--go-out` | Go 代码输出目录 |
| `--ts-out` | TypeScript 代码输出目录（需 protoc-gen-es） |
| `--go-module` | Go module 路径 |
| `--include` | protoc `-I` 附加包含路径 |
| `--plugin` | protoc 插件，可重复；格式: `name[,binary=<path>][,out=<dir>][,module=<mod>][,opt=<k=v>]` |
| `--inject-tag` | 编译后注入 struct tag（需 protoc-go-inject-tag） |
| `--no-recursive` | 只扫描顶层 .proto，不递归子目录 |

--plugin 示例：

```bash
game-dev-cli proto gen ./protos \
  --plugin go-grain,binary=$(which protoc-gen-go-grain),out=./gen,module=github.com/user/project \
  --plugin es,binary=$(which protoc-gen-es),out=./gen/ts
```

## excel gen — 从 Excel 文件生成数据

读取指定目录的 `.xlsx`/`.xls` 文件，生成 JSON 数据文件以及对应的 Go struct 和/或 TypeScript interface。

```bash
game-dev-cli excel gen <excel-dir> \
  --json-out ./data \
  --go-out ./types \
  --go-package types \
  --ts-out ./types
```

参数：

| flag | 说明 |
|------|------|
| `--json-out` | JSON 输出目录（默认: `<excel-dir>/json`） |
| `--go-out` | Go struct 输出目录 |
| `--go-package` | Go 包名（与 --go-out 配合使用） |
| `--ts-out` | TypeScript interface 输出目录 |

## 配置文件

支持 YAML 配置文件，通过 `--config` 指定：

```bash
game-dev-cli --config gamecli.yaml proto gen ./protos
```

示例见 [`gamecli.yaml.example`](./gamecli.yaml.example)。

## 在 AI 工具中使用

### Claude Code

本仓库的 `CLAUDE.md` 注册了 `proto-gen` 和 `excel-gen` 两个 skill。Claude Code 在项目中检测到本仓库时自动可用。

### 通用 MCP

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
