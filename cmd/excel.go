package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/8liang/kit/excel"
	"github.com/spf13/cobra"
)

var (
	excelJsonOut   string
	excelGoOut     string
	excelGoPackage string
	excelTsOut     string
)

var excelCmd = &cobra.Command{
	Use:   "excel",
	Short: "excel 相关操作",
}

var excelGenCmd = &cobra.Command{
	Use:   "gen <excel-dir>",
	Short: "从 Excel 文件生成 JSON 和类型定义",
	Long: `从指定目录读取 .xlsx/.xls 文件，生成 JSON 数据文件，
以及对应的 Go struct 和/或 TypeScript interface。

用法:
  game-dev-cli excel gen ./excels \
    --json-out ./data \
    --go-out ./types \
    --go-package types \
    --ts-out ./types`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		excelDir := args[0]

		if _, err := os.Stat(excelDir); os.IsNotExist(err) {
			return fmt.Errorf("Excel 目录不存在: %s", excelDir)
		}

		// 加载配置默认值（如果提供了 --config）
		if err := loadConfig(); err != nil {
			return err
		}
		if AppConfig != nil {
			if excelJsonOut == "" {
				excelJsonOut = AppConfig.Excel.JsonOut
			}
			if excelGoOut == "" {
				excelGoOut = AppConfig.Excel.GoOut
			}
			if excelGoPackage == "" {
				excelGoPackage = AppConfig.Excel.GoPackage
			}
			if excelTsOut == "" {
				excelTsOut = AppConfig.Excel.TsOut
			}
		}

		var opts []excel.Option

		if excelJsonOut != "" {
			opts = append(opts, excel.WithJson(excelJsonOut))
		} else {
			// default: output json to a subdir next to excels
			opts = append(opts, excel.WithJson(filepath.Join(excelDir, "json")))
		}

		if excelGoOut != "" && excelGoPackage != "" {
			opts = append(opts, excel.WithGoStructExport(excelGoOut, excelGoPackage))
		} else if excelGoOut != "" {
			return fmt.Errorf("--go-out 需要 --go-package 指定 Go 包名")
		}

		if excelTsOut != "" {
			opts = append(opts, excel.WithTsInterfaceExport(excelTsOut))
		}

		if err := excel.ExportToJSON(excelDir, opts...); err != nil {
			return fmt.Errorf("Excel 处理失败: %w", err)
		}

		slog.Info("Excel 处理完成", "dir", excelDir)
		return nil
	},
}

func init() {
	excelGenCmd.Flags().StringVar(&excelJsonOut, "json-out", "", "JSON 输出目录（默认: <excel-dir>/json）")
	excelGenCmd.Flags().StringVar(&excelGoOut, "go-out", "", "Go struct 输出目录")
	excelGenCmd.Flags().StringVar(&excelGoPackage, "go-package", "", "Go 包名（与 --go-out 配合使用）")
	excelGenCmd.Flags().StringVar(&excelTsOut, "ts-out", "", "TypeScript interface 输出目录")

	excelCmd.AddCommand(excelGenCmd)
	rootCmd.AddCommand(excelCmd)
}
