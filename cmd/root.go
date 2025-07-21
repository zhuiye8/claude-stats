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
	// 通用命令参数
	outputFormat string
	outputFile   string
	startDate    string
	endDate      string
	noColor      bool
	offline      bool
	costMode     string
	// daily命令特定参数
	dailyBreakdown bool
	dailyOrder     string
	// analyze命令特定参数
	modelFilter string
	showDetails bool
	configDirs  []string
	breakdown   bool
	order       string
	// blocks命令特定参数
	blocksLive            bool
	blocksTokenLimit      string
	blocksRefreshInterval int
	blocksActive          bool
	blocksRecent          bool
)

// rootCmd 代表基础命令
var rootCmd = &cobra.Command{
	Use:   "claude-stats",
	Short: "完美的Claude Code使用统计工具",
	Long: `claude-stats - 专业的Claude Code使用统计和分析工具

支持功能：
• 专门化命令架构（daily, monthly, session, blocks）
• 智能成本计算和趋势分析
• 多配置目录支持（CLAUDE_CONFIG_DIR环境变量）
• 实时监控和Token限制预警
• 跨平台支持（Windows、Mac、Linux、WSL）
• 多格式导出（JSON、CSV、表格）

基本命令：
  claude-stats daily             # 每日使用报告（默认）
  claude-stats monthly           # 月度使用报告  
  claude-stats session           # 会话分析
  claude-stats blocks            # 5小时计费窗口分析
  claude-stats blocks --live     # 实时监控模式
  claude-stats analyze           # 通用分析（兼容旧版）

快速开始：
  claude-stats                   # 显示每日使用情况
  claude-stats --breakdown       # 显示详细的模型分解
  claude-stats --json            # JSON格式输出

多配置目录：
  export CLAUDE_CONFIG_DIR="/path1,/path2"
  claude-stats daily --breakdown`,
	Version: "2.0.0",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 如果没有指定子命令，默认运行daily
		return runDaily(cmd, args)
	},
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

	// 支持默认daily命令的参数
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "输出格式 (table, json, csv)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "输出文件路径")
	rootCmd.Flags().StringVar(&startDate, "since", "", "开始日期 (YYYYMMDD)")
	rootCmd.Flags().StringVar(&endDate, "until", "", "结束日期 (YYYYMMDD)")
	rootCmd.Flags().BoolVar(&dailyBreakdown, "breakdown", false, "显示每日按模型分解的详细统计")
	rootCmd.Flags().StringVar(&dailyOrder, "order", "desc", "排序顺序: desc(最新优先) 或 asc(最旧优先)")
	rootCmd.Flags().BoolVar(&noColor, "no-color", false, "禁用颜色输出")
	rootCmd.Flags().BoolVarP(&offline, "offline", "O", false, "离线模式")
	rootCmd.Flags().StringVar(&costMode, "mode", "auto", "成本计算模式 (auto, calculate, display)")

	// Cobra也支持本地标志位，只对当前命令运行
	rootCmd.Flags().BoolP("version", "", false, "显示版本信息")

	// 确保其他文件被包含在编译中
	// 这些引用会强制Go编译器包含对应的文件
	_ = dailyCmd
	_ = blocksCmd
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
