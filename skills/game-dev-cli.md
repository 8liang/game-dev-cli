---
name: game-dev-cli
description: Use when generating code from .proto files or Excel (.xlsx) configuration files for game development
---

# Game Dev CLI

Game-dev-cli 项目工具，提供代码生成能力。

## Proto Gen

从 .proto 文件生成 Go/TypeScript 代码。

```bash
game-dev-cli proto gen <proto-dir> \
  --go-out <dir> --ts-out <dir> --go-module <module> \
  [--plugin <spec>] [--include <path>] [--inject-tag] [--no-recursive]
```

参数:
- `--go-out` — Go 代码输出目录
- `--ts-out` — TypeScript 输出目录（等价于 `--plugin es,out=<dir>`）
- `--go-module` — Go module 路径
- `--include` — protoc `-I` 附加包含路径（可重复）
- `--plugin` — protoc 插件，可重复；格式: `name[,binary=<path>][,out=<dir>][,module=<mod>][,opt=<k=v>]`
- `--inject-tag` — 编译后注入 struct tag（需 protoc-go-inject-tag）
- `--no-recursive` — 只扫描顶层 .proto，不递归子目录

## Excel Gen

从 Excel 文件生成 JSON + Go structs + TypeScript interfaces。

```bash
game-dev-cli excel gen <excel-dir> \
  --json-out <dir> --go-out <dir> --go-package <pkg> --ts-out <dir>
```

参数:
- `--json-out` — JSON 输出目录（默认: `<excel-dir>/json`）
- `--go-out` — Go struct 输出目录（需配合 `--go-package`）
- `--go-package` — Go 包名
- `--ts-out` — TypeScript interface 输出目录
