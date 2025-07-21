package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zhuiye8/claude-stats/pkg/parser"
	"github.com/zhuiye8/claude-stats/pkg/formatter"
	"github.com/zhuiye8/claude-stats/pkg/models"
)

// blocksCmd 代表blocks命令
var blocksCmd = &cobra.Command{
	Use:   "blocks [目录路径]",
	Short: "分析5小时计费窗口使用情况",
	Long: `分析Claude Code的5小时计费窗口使用情况，帮助理解使用模式和预测成本。

5小时计费窗口是Claude订阅模式的计费单位，每个窗口内有固定的消息限制。
此命令可以帮助您：
• 了解不同时间段的使用强度
• 识别活跃使用窗口
• 预测当前窗口的使用趋势
• 监控实时使用率

支持功能：
• 自动检测窗口边界（每5小时一个周期）
• 显示活跃窗口状态和剩余时间
• 计算燃烧速率和预测使用量
• 实时监控模式（--live）
• Token限制预警（--token-limit）

示例：
  claude-stats blocks                    # 显示基本窗口分析
  claude-stats blocks --live             # 实时监控当前窗口
  claude-stats blocks --live -t 500000   # 设置Token限制监控
  claude-stats blocks --active           # 只显示活跃窗口
  claude-stats blocks --recent           # 显示最近的窗口
  claude-stats blocks --json             # JSON格式输出`,
	Args: cobra.MaximumNArgs(1),
	RunE: runBlocks,
}

func init() {
	rootCmd.AddCommand(blocksCmd)

	// blocks命令特定的标志位
	blocksCmd.Flags().BoolVar(&blocksLive, "live", false, "实时监控模式")
	blocksCmd.Flags().StringVarP(&blocksTokenLimit, "token-limit", "t", "", "Token限制 (数字或'max')")
	blocksCmd.Flags().IntVar(&blocksRefreshInterval, "refresh-interval", 3, "实时模式刷新间隔(秒)")
	blocksCmd.Flags().BoolVar(&blocksActive, "active", false, "只显示活跃窗口")
	blocksCmd.Flags().BoolVar(&blocksRecent, "recent", false, "显示最近的窗口")
	
	// 继承通用标志位
	blocksCmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "输出格式 (table, json)")
	blocksCmd.Flags().StringVarP(&outputFile, "output", "o", "", "输出文件路径")
	blocksCmd.Flags().StringVar(&startDate, "since", "", "开始日期 (YYYYMMDD)")
	blocksCmd.Flags().StringVar(&endDate, "until", "", "结束日期 (YYYYMMDD)")
	blocksCmd.Flags().BoolVar(&noColor, "no-color", false, "禁用颜色输出")
	blocksCmd.Flags().BoolVarP(&offline, "offline", "O", false, "离线模式")
}

func runBlocks(cmd *cobra.Command, args []string) error {
	// 确定要分析的目录列表
	targetDirs := getTargetDirectories(args)
	
	if verbose {
		fmt.Printf("🔍 分析目录: %s\n", strings.Join(targetDirs, ", "))
	}

	// 如果是实时模式，循环执行
	if blocksLive {
		return runLiveBlocks(targetDirs)
	}

	// 执行一次性分析
	return runSingleBlocks(targetDirs)
}

// runSingleBlocks 执行一次性blocks分析
func runSingleBlocks(targetDirs []string) error {
	stats, err := parseDirectories(targetDirs)
	if err != nil {
		return err
	}

	// 创建解析器并分析blocks
	claudeParser := parser.NewClaudeParser()
	blocksReport, err := claudeParser.AnalyzeBlocks(stats)
	if err != nil {
		return fmt.Errorf("分析blocks失败: %w", err)
	}

	// 应用过滤器
	blocksReport = filterBlocks(blocksReport)

	// 输出结果
	return outputBlocksReport(blocksReport)
}

// runLiveBlocks 执行实时监控模式
func runLiveBlocks(targetDirs []string) error {
	if verbose {
		fmt.Printf("🔴 启动实时监控模式 (刷新间隔: %d秒)\n", blocksRefreshInterval)
		fmt.Printf("💡 按 Ctrl+C 退出监控\n\n")
	}

	// 解析Token限制
	var tokenLimit int
	if blocksTokenLimit != "" {
		if blocksTokenLimit == "max" {
			// TODO: 从历史数据中找到最大Token使用量
			tokenLimit = 500000 // 默认值
		} else {
			if limit, err := strconv.Atoi(blocksTokenLimit); err == nil {
				tokenLimit = limit
			}
		}
	}

	for {
		// 清屏（在支持的终端中）
		fmt.Print("\033[2J\033[H")
		
		// 显示时间戳
		fmt.Printf("🕐 监控时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		if tokenLimit > 0 {
			fmt.Printf("⚡ Token限制: %s\n", formatNumber(tokenLimit))
		}
		fmt.Println()

		// 执行分析
		stats, err := parseDirectories(targetDirs)
		if err != nil {
			fmt.Printf("❌ 解析失败: %v\n", err)
		} else {
			claudeParser := parser.NewClaudeParser()
			blocksReport, err := claudeParser.AnalyzeBlocks(stats)
			if err != nil {
				fmt.Printf("❌ 分析失败: %v\n", err)
			} else {
				// 只显示活跃窗口
				activeReport := filterActiveBlocks(blocksReport)
				
				// 显示Token限制警告
				if tokenLimit > 0 {
					checkTokenLimits(activeReport, tokenLimit)
				}
				
				// 格式化输出
				formatter := formatter.NewFormatter()
				if noColor {
					formatter.Colors.Enabled = false
				}
				
				output, err := formatter.FormatBlocks(activeReport)
				if err != nil {
					fmt.Printf("❌ 格式化失败: %v\n", err)
				} else {
					fmt.Print(output)
				}
			}
		}

		// 等待下次刷新
		time.Sleep(time.Duration(blocksRefreshInterval) * time.Second)
	}
}

// parseDirectories 解析所有目录并聚合数据
func parseDirectories(targetDirs []string) (*models.UsageStats, error) {
	// 创建解析器
	claudeParser := parser.NewClaudeParser()
	claudeParser.Verbose = verbose
	claudeParser.SkipErrors = true

	// 设置日期过滤器
	if startDate != "" || endDate != "" {
		dateFilter, err := createDateFilter(startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("日期格式错误: %w", err)
		}
		claudeParser.DateFilter = dateFilter
	}

	// 解析所有目录并聚合数据
	aggregatedStats := &models.UsageStats{
		ModelStats:   make(map[string]models.TokenUsage),
		DailyStats:   make(map[string]models.TokenUsage),
		SessionStats: make(map[string]models.SessionInfo),
		ProjectStats: make(map[string]models.ProjectStats),
		MessageTypes: make(map[string]int),
	}

	var successfulDirs []string
	
	for _, targetDir := range targetDirs {
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			continue // 跳过不存在的目录
		}
		
		stats, err := claudeParser.ParseDirectory(targetDir)
		if err != nil {
			if verbose {
				fmt.Printf("⚠️  解析目录失败，跳过 %s: %v\n", targetDir, err)
			}
			continue
		}

		// 聚合数据
		mergeUsageStats(aggregatedStats, stats)
		successfulDirs = append(successfulDirs, targetDir)
	}

	if len(successfulDirs) == 0 {
		return nil, fmt.Errorf("没有找到有效的Claude配置目录")
	}

	// 最终处理
	claudeParser.FinalizeStats(aggregatedStats)
	
	return aggregatedStats, nil
}

// filterBlocks 过滤blocks报告
func filterBlocks(report *models.BlocksReport) *models.BlocksReport {
	if !blocksActive && !blocksRecent {
		return report // 不需要过滤
	}

	var filteredBlocks []models.BillingBlock
	
	for _, block := range report.Blocks {
		include := true
		
		if blocksActive && !block.IsActive {
			include = false
		}
		
		if blocksRecent {
			// 只包含最近24小时的窗口
			since := time.Now().Add(-24 * time.Hour)
			if block.StartTime.Before(since) {
				include = false
			}
		}
		
		if include {
			filteredBlocks = append(filteredBlocks, block)
		}
	}
	
	return &models.BlocksReport{
		Blocks:    filteredBlocks,
		Summary:   report.Summary,
		TotalCost: report.TotalCost,
	}
}

// filterActiveBlocks 只保留活跃窗口（用于实时监控）
func filterActiveBlocks(report *models.BlocksReport) *models.BlocksReport {
	var activeBlocks []models.BillingBlock
	
	for _, block := range report.Blocks {
		if block.IsActive {
			activeBlocks = append(activeBlocks, block)
		}
	}
	
	return &models.BlocksReport{
		Blocks:    activeBlocks,
		Summary:   report.Summary,
		TotalCost: report.TotalCost,
	}
}

// checkTokenLimits 检查Token限制并显示警告
func checkTokenLimits(report *models.BlocksReport, tokenLimit int) {
	for _, block := range report.Blocks {
		if !block.IsActive {
			continue
		}
		
		currentTokens := block.Tokens.GetTotalTokens()
		usage := float64(currentTokens) / float64(tokenLimit) * 100
		
		if usage > 90 {
			fmt.Printf("🚨 警告: 当前窗口Token使用率 %.1f%% (%.0f/%d)\n", usage, float64(currentTokens), tokenLimit)
		} else if usage > 75 {
			fmt.Printf("⚠️  注意: 当前窗口Token使用率 %.1f%% (%.0f/%d)\n", usage, float64(currentTokens), tokenLimit)
		} else if usage > 50 {
			fmt.Printf("💡 提示: 当前窗口Token使用率 %.1f%% (%.0f/%d)\n", usage, float64(currentTokens), tokenLimit)
		}
	}
}

// outputBlocksReport 输出blocks报告
func outputBlocksReport(report *models.BlocksReport) error {
	// 格式化并输出结果
	formatter := formatter.NewFormatter()
	formatter.Verbose = verbose
	
	// 设置颜色选项
	if noColor {
		formatter.Colors.Enabled = false
	}

	var output string
	var err error
	
	switch strings.ToLower(outputFormat) {
	case "json":
		output, err = formatter.FormatBlocksJSON(report)
	case "table", "":
		output, err = formatter.FormatBlocks(report)
	default:
		return fmt.Errorf("不支持的格式: %s", outputFormat)
	}
	
	if err != nil {
		return fmt.Errorf("格式化失败: %w", err)
	}

	// 输出结果
	if outputFile != "" {
		err = writeToFile(output, outputFile)
		if err != nil {
			return fmt.Errorf("写入文件失败: %w", err)
		}
		fmt.Printf("✅ 报告已保存到: %s\n", outputFile)
	} else {
		fmt.Print(output)
	}

	return nil
}

// formatNumber 格式化数字，添加千位分隔符
func formatNumber(num int) string {
	str := strconv.Itoa(num)
	if len(str) <= 3 {
		return str
	}

	var result strings.Builder
	for i, char := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(char)
	}
	return result.String()
} 