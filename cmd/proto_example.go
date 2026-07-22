package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/8liang/kit/protobuf"
	"github.com/spf13/cobra"
)

var protoExampleCmd = &cobra.Command{
	Use:   "example",
	Short: "编译 example/proto → go-project + ts-project",
	Long: `在 CWD/ 下新建 example/ 目录，结构如下：
  example/
    proto/               ← 从二进制提取的 .proto 源文件
      error.proto
      messages.proto
      svc/agent.svc.proto
    go-project/           ← Go 代码输出
    ts-project/           ← TypeScript 代码输出

根据 proto/ 下的 .proto 文件编译生成：
  Go: go-project/pkg/{messages,svc/agent}/*.pb.go
  TS: ts-project/*_pb.{js,d.ts}

执行指令（等价）：
  game-dev-cli proto gen example/proto/svc \\
    --go-out example/go-project --go-module github.com/8liang/game-dev-cli \\
    --plugin go-grain,out=example/go-project,module=github.com/8liang/game-dev-cli \\
    --include example/proto --inject-tag --no-recursive

  game-dev-cli proto gen example/proto \\
    --go-out example/go-project --go-module github.com/8liang/game-dev-cli \\
    --ts-out example/ts-project --inject-tag --no-recursive`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("获取 CWD 失败: %w", err)
		}

		exampleDir := filepath.Join(cwd, "example")
		protoDir := filepath.Join(exampleDir, "proto")
		goOut := filepath.Join(exampleDir, "go-project")
		tsOut := filepath.Join(exampleDir, "ts-project")
		mod := "github.com/8liang/game-dev-cli"

		// 1. 提取嵌入的 proto 文件到 CWD/example/proto/。
		if err := extractEmbedFS(ExampleProtoFS, "example_proto_static", protoDir); err != nil {
			return fmt.Errorf("提取 proto 文件失败: %w", err)
		}

		// 2. 创建输出目录。
		os.MkdirAll(goOut, 0755)
		os.MkdirAll(tsOut, 0755)

		// 3. proto setup --install（下载 import 依赖到 protoDir/.proto-cache）。
		fmt.Println("=== proto setup --install ===")
		if _, err := protobuf.Setup(protoDir, protobuf.SetupWithInstall()); err != nil {
			return err
		}

		// 4. STEP1: svc -> go-project。
		fmt.Println("=== STEP1: svc -> go-project ===")
		if err := protobuf.Compile(filepath.Join(protoDir, "svc"),
			protobuf.WithIncludePaths(protoDir),
			protobuf.WithGoModule(mod),
			protobuf.WithGetOutPath(func(_ string) string { return goOut }),
			protobuf.WithPluginConfig(protobuf.Plugin{Name: "go-grain", OutPath: goOut, Module: mod}),
			protobuf.WithInjectTag(),
			protobuf.WithNonRecursive(),
		); err != nil {
			return fmt.Errorf("STEP1 svc 编译失败: %w", err)
		}

		// 5. STEP2: messages/error -> go-project + ts-project。
		fmt.Println("=== STEP2: messages/error -> go-project + ts-project ===")
		if err := protobuf.Compile(protoDir,
			protobuf.WithGoModule(mod),
			protobuf.WithGetOutPath(func(_ string) string { return goOut }),
			protobuf.WithPluginConfig(protobuf.Plugin{Name: "es", OutPath: tsOut}),
			protobuf.WithInjectTag(),
			protobuf.WithNonRecursive(),
		); err != nil {
			return fmt.Errorf("STEP2 编译失败: %w", err)
		}

		slog.Info("编译完成", "go", goOut, "ts", tsOut)
		return nil
	},
}

// extractEmbedFS 把嵌入式 FS 中指定前缀的内容提取到目标目录。
func extractEmbedFS(fsys embed.FS, prefix, dest string) error {
	return fs.WalkDir(fsys, prefix, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(prefix, path)
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := fsys.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	})
}

func init() {
	protoCmd.AddCommand(protoExampleCmd)
}
