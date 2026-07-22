package cmd

import (
	"fmt"
	"os"

	"github.com/8liang/kit/protobuf"
	"github.com/spf13/cobra"
)

var (
	setupInstall       bool
	setupCacheDir      string
	setupProtocVersion string
	setupIgnoreImports []string
)

var protoSetupCmd = &cobra.Command{
	Use:   "setup <proto-dir>",
	Short: "检测并安装 proto 依赖",
	Long: `扫描 <proto-dir> 下的 .proto 文件，检测 protoc、插件、import 依赖。

默认只报告缺失项。加 --install 自动下载安装。

示例:
  game-dev-cli proto setup ./example/proto         # 检测模式
  game-dev-cli proto setup ./example/proto --install  # 检测并安装`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		protoDir := args[0]

		if _, err := os.Stat(protoDir); os.IsNotExist(err) {
			return fmt.Errorf("proto 目录不存在: %s", protoDir)
		}

		var opts []protobuf.SetupOption

		if setupInstall {
			opts = append(opts, protobuf.SetupWithInstall())
		}
		if setupCacheDir != "" {
			opts = append(opts, protobuf.SetupWithCacheDir(setupCacheDir))
		}
		if setupProtocVersion != "" {
			opts = append(opts, protobuf.SetupWithProtocVersion(setupProtocVersion))
		}
		if len(setupIgnoreImports) > 0 {
			opts = append(opts, protobuf.SetupWithIgnoreImport(setupIgnoreImports...))
		}
		if verbose {
			opts = append(opts, protobuf.SetupWithVerbose())
		}

		report, err := protobuf.Setup(protoDir, opts...)
		if err != nil {
			return err
		}

		fmt.Print(protobuf.FormatReport(report))

		if report.MissingCount > 0 && !setupInstall {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	protoSetupCmd.Flags().BoolVar(&setupInstall, "install", false, "下载安装缺失依赖")
	protoSetupCmd.Flags().StringVar(&setupCacheDir, "cache-dir", "", "缓存根目录（默认 <proto-dir>/.proto-cache）")
	protoSetupCmd.Flags().StringVar(&setupProtocVersion, "protoc-version", "", "protoc 目标版本（默认 latest）")
	protoSetupCmd.Flags().StringSliceVar(&setupIgnoreImports, "ignore-import", nil, "跳过指定 import 路径（可重复）")

	protoCmd.AddCommand(protoSetupCmd)
}
