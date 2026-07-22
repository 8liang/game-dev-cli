package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/8liang/kit/protobuf"
	"github.com/spf13/cobra"
)

var (
	protoGoOut       string
	protoTsOut       string
	protoGoModule    string
	protoInclude     []string
	protoInjectTag   bool
	protoPlugins     []string
	protoNoRecursive bool
)

var protoCmd = &cobra.Command{
	Use:   "proto",
	Short: "proto 相关操作",
}

var protoGenCmd = &cobra.Command{
	Use:   "gen <proto-dir>",
	Short: "从 .proto 文件生成代码",
	Long: `从指定目录读取 .proto 文件，生成 Go 和/或 TypeScript 代码。

用法:
  game-dev-cli proto gen ./protos \
    --go-out ./gen/go --go-module github.com/user/project \
    --plugin go-grain,out=./gen/go,module=github.com/user/project \
    --plugin es,binary=$(which protoc-gen-es),out=./gen/ts,module=github.com/user/project \
    --inject-tag --no-recursive

--plugin 可重复，格式: name[,binary=<path>][,out=<dir>][,module=<mod>][,opt=<k=v>]
  - name:        插件名，对应 protoc-gen-<name>
  - binary:      插件二进制完整路径（npm 装的 protoc-gen-es 需指定），省略则在 PATH / GOPATH/bin 查找
  - out:         --<name>_out 输出目录
  - module:      --<name>_opt=module=<mod>（protoc-gen-es 不认，省略）
  - opt:         额外 --<name>_opt=<k=v>，可重复，如 opt=target=ts
--ts-out 是 --plugin es,out=<dir> 的便捷别名（向后兼容）。
--no-recursive 只扫描顶层目录，便于分目录多次调用以按 proto 选择插件。

要求: protoc 编译器已安装。`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		protoDir := args[0]

		if _, err := os.Stat(protoDir); os.IsNotExist(err) {
			return fmt.Errorf("proto 目录不存在: %s", protoDir)
		}

		// 加载配置默认值（如果提供了 --config）。
		if err := loadConfig(); err != nil {
			return err
		}
		if AppConfig != nil {
			if protoGoOut == "" {
				protoGoOut = AppConfig.Proto.GoOut
			}
			if protoTsOut == "" {
				protoTsOut = AppConfig.Proto.TsOut
			}
			if protoGoModule == "" {
				protoGoModule = AppConfig.Proto.GoModule
			}
			if len(protoInclude) == 0 && len(AppConfig.Proto.Include) > 0 {
				protoInclude = AppConfig.Proto.Include
			}
			if !protoInjectTag && AppConfig.Proto.InjectTag {
				protoInjectTag = true
			}
			if len(protoPlugins) == 0 && len(AppConfig.Proto.Plugins) > 0 {
				protoPlugins = AppConfig.Proto.Plugins
			}
			if !protoNoRecursive && AppConfig.Proto.NoRecursive {
				protoNoRecursive = true
			}
		}

		if err := protobuf.CheckProtoc(); err != nil {
			return err
		}

		var opts []protobuf.Option

		if len(protoInclude) > 0 {
			opts = append(opts, protobuf.WithIncludePaths(protoInclude...))
		}

		if verbose {
			opts = append(opts, protobuf.WithDebug())
		}

		if protoInjectTag {
			opts = append(opts, protobuf.WithInjectTag())
		}

		// --go-module 同时用于 Go 输出分目录与插件 module。
		if protoGoModule != "" {
			opts = append(opts, protobuf.WithGoModule(protoGoModule))
		}

		if protoNoRecursive {
			opts = append(opts, protobuf.WithNonRecursive())
		}

		if protoGoOut != "" {
			opts = append(opts, protobuf.WithGetOutPath(func(_ string) string { return protoGoOut }))
		}

		// 解析 --plugin 规格。
		plugins, err := parsePlugins(protoPlugins)
		if err != nil {
			return err
		}
		// --ts-out 便捷别名：等价于 --plugin es,out=<ts-out>（protoc-gen-es 不使用 module opt）。
		if protoTsOut != "" {
			plugins = append(plugins, pluginSpec{name: "es", out: protoTsOut})
		}
		for _, p := range plugins {
			opts = append(opts, protobuf.WithPluginConfig(protobuf.Plugin{
				Name:      p.name,
				OutPath:   p.out,
				Module:    p.module,
				Binary:    p.binary,
				ExtraOpts: p.extraOpts,
			}))
		}

		if err := protobuf.Compile(protoDir, opts...); err != nil {
			return fmt.Errorf("proto 编译失败: %w", err)
		}

		slog.Info("proto 编译完成", "dir", protoDir)
		return nil
	},
}

type pluginSpec struct {
	name      string
	binary    string
	out       string
	module    string
	extraOpts []string
}

// parsePlugins 解析 --plugin 规格: name[,binary=<path>][,out=<dir>][,module=<mod>][,opt=<k=v>]。
func parsePlugins(specs []string) ([]pluginSpec, error) {
	var result []pluginSpec
	for _, spec := range specs {
		parts := strings.Split(spec, ",")
		var ps pluginSpec
		for i, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			// 第一段若不含 "=" 视为插件名。
			if i == 0 && !strings.Contains(part, "=") {
				ps.name = part
				continue
			}
			kv := strings.SplitN(part, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("无效的 plugin 规格 %q: 期望 key=value", part)
			}
			switch kv[0] {
			case "binary":
				ps.binary = kv[1]
			case "out":
				ps.out = kv[1]
			case "module":
				ps.module = kv[1]
			case "opt":
				ps.extraOpts = append(ps.extraOpts, kv[1])
			default:
				return nil, fmt.Errorf("未知的 plugin 字段 %q", kv[0])
			}
		}
		if ps.name == "" {
			return nil, fmt.Errorf("plugin 规格缺少 name: %q", spec)
		}
		result = append(result, ps)
	}
	return result, nil
}

func init() {
	protoGenCmd.Flags().StringVar(&protoGoOut, "go-out", "", "Go 代码输出目录")
	protoGenCmd.Flags().StringVar(&protoTsOut, "ts-out", "", "TypeScript 输出目录（等价于 --plugin es,out=<dir>，需 protoc-gen-es）")
	protoGenCmd.Flags().StringVar(&protoGoModule, "go-module", "", "Go module 路径（同时用于 --go_opt=module 与插件 module）")
	protoGenCmd.Flags().StringSliceVar(&protoInclude, "include", nil, "protoc -I 附加包含路径")
	protoGenCmd.Flags().BoolVar(&protoInjectTag, "inject-tag", false, "编译后注入 struct tag（需安装 protoc-go-inject-tag）")
	protoGenCmd.Flags().StringArrayVar(&protoPlugins, "plugin", nil, "protoc 插件，可重复；格式: name[,binary=<path>][,out=<dir>][,module=<mod>]")
	protoGenCmd.Flags().BoolVar(&protoNoRecursive, "no-recursive", false, "只扫描顶层目录的 .proto，不递归子目录")

	protoCmd.AddCommand(protoGenCmd)
	rootCmd.AddCommand(protoCmd)
}
