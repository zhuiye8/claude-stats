package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zhuiye8/claude-stats/pkg/parser"
	"github.com/zhuiye8/claude-stats/pkg/formatter"
	"github.com/zhuiye8/claude-stats/pkg/models"
)

var (
	outputFormat string
	outputFile   string
	startDate    string
	endDate      string
	modelFilter  string
	showDetails  bool
	noColor      bool
)

// analyzeCmd 代表analyze命令
var analyzeCmd = &cobra.Command{
	Use:   "analyze [目录路径]",
	Short: "分析Claude Code使用统计",
	Long: `分析指定目录中的Claude Code JSONL日志文件，生成详细的使用统计报告。

支持的功能：
• 自动检测Claude日志目录
• Token使用统计（输入、输出、缓存）
• 成本估算（API和订阅模式）
• 按模型、日期、会话分组统计
• 多种输出格式（表格、JSON、CSV）

示例：
  claude-stats analyze                    # 分析默认Claude目录
  claude-stats analyze ~/claude-logs     # 分析指定目录
  claude-stats analyze --json            # JSON格式输出
  claude-stats analyze --csv report.csv  # 导出CSV报告`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	// 添加命令特定的标志位
	analyzeCmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "输出格式 (table, json, csv)")
	analyzeCmd.Flags().StringVarP(&outputFile, "output", "o", "", "输出文件路径")
	analyzeCmd.Flags().StringVar(&startDate, "start", "", "开始日期 (YYYY-MM-DD)")
	analyzeCmd.Flags().StringVar(&endDate, "end", "", "结束日期 (YYYY-MM-DD)")
	analyzeCmd.Flags().StringVar(&modelFilter, "model", "", "过滤特定模型")
	analyzeCmd.Flags().BoolVarP(&showDetails, "details", "d", false, "显示详细信息")
	analyzeCmd.Flags().BoolVar(&noColor, "no-color", false, "禁用颜色输出")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	// 确定要分析的目录
	targetDir := getTargetDirectory(args)
	
	if verbose {
		fmt.Printf("🔍 分析目录: %s\n", targetDir)
	}

	// 检查目录是否存在
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return fmt.Errorf("目录不存在: %s", targetDir)
	}

	// 创建解析器
	claudeParser := parser.NewClaudeParser()
	claudeParser.Verbose = verbose
	claudeParser.SkipErrors = true

	// 设置日期过滤器
	if startDate != "" || endDate != "" {
		dateFilter, err := createDateFilter(startDate, endDate)
		if err != nil {
			return fmt.Errorf("日期格式错误: %w", err)
		}
		claudeParser.DateFilter = dateFilter
	}

	// 解析目录
	stats, err := claudeParser.ParseDirectory(targetDir)
	if err != nil {
		return fmt.Errorf("解析失败: %w", err)
	}

	// 应用模型过滤器
	if modelFilter != "" {
		stats = filterByModel(stats, modelFilter)
	}

	// 格式化并输出结果
	formatter := formatter.NewFormatter()
	formatter.ShowDetails = showDetails
	formatter.Verbose = verbose
	
	// 设置颜色选项
	if noColor {
		formatter.Colors.Enabled = false
	}

	output, err := formatter.Format(stats, outputFormat)
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

// getTargetDirectory 获取目标目录
func getTargetDirectory(args []string) string {
	if len(args) > 0 {
		return args[0]
	}

	// 自动检测Claude目录
	return getDefaultClaudeDirectory()
}

// getDefaultClaudeDirectory 获取默认的Claude目录
func getDefaultClaudeDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	// 根据操作系统确定Claude配置目录
	var claudeDir string
	switch runtime.GOOS {
	case "windows":
		claudeDir = filepath.Join(homeDir, "AppData", "Roaming", "claude", "projects")
	case "darwin": // macOS
		claudeDir = filepath.Join(homeDir, "Library", "Application Support", "claude", "projects")
	default: // Linux, WSL等
		claudeDir = filepath.Join(homeDir, ".config", "claude", "projects")
	}

	// 检查目录是否存在
	if _, err := os.Stat(claudeDir); err == nil {
		return claudeDir
	}

	// 如果标准目录不存在，尝试其他可能的位置
	alternatives := []string{
		filepath.Join(homeDir, ".claude", "projects"),
		filepath.Join(homeDir, "claude-logs"),
		".",
	}

	for _, alt := range alternatives {
		if _, err := os.Stat(alt); err == nil {
			return alt
		}
	}

	return "."
}

// createDateFilter 创建日期过滤器
func createDateFilter(start, end string) (*parser.DateFilter, error) {
	filter := &parser.DateFilter{}
	
	if start != "" {
		startTime, err := parseDate(start)
		if err != nil {
			return nil, fmt.Errorf("开始日期解析失败: %w", err)
		}
		filter.StartDate = &startTime
	}
	
	if end != "" {
		endTime, err := parseDate(end)
		if err != nil {
			return nil, fmt.Errorf("结束日期解析失败: %w", err)
		}
		filter.EndDate = &endTime
	}
	
	return filter, nil
}

// parseDate 解析日期字符串
func parseDate(dateStr string) (time.Time, error) {
	// 支持多种日期格式
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"01/02/2006",
		"2006-01-02 15:04:05",
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("无法解析日期: %s", dateStr)
}

// filterByModel 按模型过滤统计数据
func filterByModel(stats *models.UsageStats, model string) *models.UsageStats {
	filtered := &models.UsageStats{
		ModelStats:   make(map[string]models.TokenUsage),
		DailyStats:   make(map[string]models.TokenUsage),
		SessionStats: make(map[string]models.SessionInfo),
		DetectedMode: stats.DetectedMode,
	}

	// 过滤模型统计
	for modelName, usage := range stats.ModelStats {
		if strings.Contains(strings.ToLower(modelName), strings.ToLower(model)) {
			filtered.ModelStats[modelName] = usage
			filtered.TotalTokens.Add(usage)
		}
	}

	// 过滤会话统计
	for sessionID, session := range stats.SessionStats {
		if strings.Contains(strings.ToLower(session.Model), strings.ToLower(model)) {
			filtered.SessionStats[sessionID] = session
			filtered.TotalSessions++
		}
	}

	return filtered
}

// writeToFile 写入文件
func writeToFile(content, filename string) error {
	return os.WriteFile(filename, []byte(content), 0644)
} 