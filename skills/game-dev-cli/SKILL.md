---
name: game-dev-cli
description: 生成代码的工具——从 .proto 文件生成 Go/TypeScript 代码，或从 Excel (.xlsx) 配置文件生成 JSON、Go structs、TypeScript interfaces
---

# Game Dev CLI

Game-dev-cli 是游戏开发命令行工具，提供 proto 和 Excel 配置文件的代码生成能力。

## 使用流程

### Proto gen：从 .proto 生成 Go/TypeScript 代码

1. 确认 proto 目录下有 `.proto` 文件
2. 运行命令：
   ```bash
   game-dev-cli proto gen <proto-dir> \
     --go-out <dir> --ts-out <dir> --go-module <module> \
     [--plugin <spec>] [--include <path>] [--inject-tag] [--no-recursive]
   ```
3. 输出：指定目录下的 `.pb.go` 和 `.ts` 文件

### Excel gen：从 Excel 生成配置文件代码

1. 确认 Excel 目录下有 `.xlsx` 文件
2. 运行命令：
   ```bash
   game-dev-cli excel gen <excel-dir> \
     --json-out <dir> --go-out <dir> --go-package <pkg> --ts-out <dir>
   ```
3. 输出：JSON 数据文件、Go structs、TypeScript interfaces

## 参数参考

### Proto gen

| 参数 | 说明 |
|------|------|
| `--go-out` | Go 代码输出目录 |
| `--ts-out` | TypeScript 输出目录（等价于 `--plugin es,out=<dir>`） |
| `--go-module` | Go module 路径 |
| `--include` | protoc `-I` 附加包含路径（可重复） |
| `--plugin` | protoc 插件，可重复；格式: `name[,binary=<path>][,out=<dir>][,module=<mod>][,opt=<k=v>]` |
| `--inject-tag` | 编译后注入 struct tag（需 protoc-go-inject-tag） |
| `--no-recursive` | 只扫描顶层 .proto，不递归子目录 |

### Excel gen

| 参数 | 说明 |
|------|------|
| `--json-out` | JSON 输出目录（默认: `<excel-dir>/json`） |
| `--go-out` | Go struct 输出目录（需配合 `--go-package`） |
| `--go-package` | Go 包名 |
| `--ts-out` | TypeScript interface 输出目录 |

## 什么时候使用

- 游戏策划表（Excel）需要转成程序可读的代码和配置文件时
- 协议定义文件（.proto）需要同步生成 Go 后端和 TypeScript 前端的类型代码时

## 什么时候不用

- 只查 game-dev-cli 命令帮助时（直接运行 `game-dev-cli --help`）
- 需要生成其他语言（Python、Java 等）的代码时
