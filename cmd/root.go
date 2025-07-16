package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
)

// rootCmd 代表基础命令
var rootCmd = &cobra.Command{
	Use:   "claude-stats",
	Short: "完美的Claude Code使用统计工具",
	Long: `claude-stats - 专业的Claude Code使用统计和分析工具

支持功能：
• 自动识别API模式和订阅模式
• 详细token使用统计（输入、输出、缓存）
• 智能成本计算和趋势分析
• 跨平台支持（Windows、Mac、Linux、WSL）
• 多格式导出（JSON、CSV、Markdown）
• 实时监控进行中的会话

使用示例：
  claude-stats analyze           # 分析默认Claude目录
  claude-stats analyze ~/logs    # 分析指定目录
  claude-stats monitor           # 实时监控模式
  claude-stats report --monthly  # 生成月度报告`,
	Version: "1.0.0",
}

// Execute 添加所有子命令到根命令并设置相应的标志位
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// 全局标志位
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件 (默认: $HOME/.claude-stats.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "详细输出")

	// Cobra也支持本地标志位，只对当前命令运行
	rootCmd.Flags().BoolP("version", "", false, "显示版本信息")
}

// initConfig 读取配置文件和环境变量
func initConfig() {
	if cfgFile != "" {
		// 使用命令行指定的配置文件
		viper.SetConfigFile(cfgFile)
	} else {
		// 查找home目录
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// 在home目录下查找".claude-stats"配置文件
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".claude-stats")
	}

	viper.AutomaticEnv() // 读取匹配的环境变量

	// 如果配置文件存在，则读取它
	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "使用配置文件:", viper.ConfigFileUsed())
	}
} 