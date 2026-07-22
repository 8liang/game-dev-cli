package cmd

import (
	"fmt"
	"os"

	"github.com/8liang/kit/viperparser"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

type Config struct {
	Proto struct {
		GoOut       string   `mapstructure:"default_go_out"`
		TsOut       string   `mapstructure:"default_ts_out"`
		GoModule    string   `mapstructure:"go_module"`
		Include     []string `mapstructure:"include_paths"`
		InjectTag   bool     `mapstructure:"inject_tag"`
		Plugins     []string `mapstructure:"plugins"`
		NoRecursive bool     `mapstructure:"no_recursive"`
	} `mapstructure:"proto"`
	Excel struct {
		JsonOut   string `mapstructure:"default_json_out"`
		GoOut     string `mapstructure:"default_go_out"`
		TsOut     string `mapstructure:"default_ts_out"`
		GoPackage string `mapstructure:"default_go_package"`
	} `mapstructure:"excel"`
}

var AppConfig *Config

var rootCmd = &cobra.Command{
	Use:   "game-dev-cli",
	Short: "game-dev-cli — AI vibe coding 工具集",
	Long: `game-dev-cli 是给 AI vibe coding 使用的命令行工具。
支持 proto 编译、Excel 数据导出等能力。`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func loadConfig() error {
	if cfgFile == "" {
		return nil
	}
	AppConfig = &Config{}
	opts := []viperparser.Option{viperparser.WithUrl(cfgFile)}
	if err := viperparser.Unmarshal(AppConfig, opts...); err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径 (支持 .yaml/.env/.json)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "输出调试日志")
}
