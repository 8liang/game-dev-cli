# game-dev-cli

AI vibe coding 工具集。基于 `.proto` 文件生成代码、从 Excel 文件导出数据。

## 构建

```bash
cd game-dev-cli && go build -o ../bin/game-dev-cli ./...
```

## proto gen — 从 .proto 文件生成代码

读取指定目录的 `.proto` 文件，调用 `protoc` 编译生成 Go 和/或 TypeScript 代码。

**要求：** `protoc` 已安装。TS 生成需安装 `protoc-gen-es`。

```bash
./bin/game-dev-cli proto gen <proto-dir> \
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
| `--inject-tag` | 编译后注入 struct tag（需 protoc-go-inject-tag） |

## excel gen — 从 Excel 文件生成数据

读取指定目录的 `.xlsx`/`.xls` 文件，生成 JSON 数据文件以及对应的 Go struct 和/或 TypeScript interface。

```bash
./bin/game-dev-cli excel gen <excel-dir> \
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
./bin/game-dev-cli --config gamecli.yaml proto gen ./protos
```

示例见 [`gamecli.yaml.example`](./gamecli.yaml.example)。

## 在 Claude Code 中使用

本仓库的 `CLAUDE.md` 注册了两个 skill：

- `proto-gen` — AI 检测到 `.proto` 文件变更时自动生成代码
- `excel-gen` — AI 检测到 Excel 文件变更时自动导出数据
