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
• 多配置目录支持（通过CLAUDE_CONFIG_DIR环境变量）

示例：
  claude-stats analyze                    # 分析默认Claude目录
  claude-stats analyze ~/claude-logs     # 分析指定目录
  claude-stats analyze --json            # JSON格式输出
  claude-stats analyze --csv report.csv  # 导出CSV报告
  claude-stats analyze --breakdown       # 显示模型详细分解
  claude-stats analyze --since 20240101 --until 20241231  # 指定日期范围
  
多配置目录：
  export CLAUDE_CONFIG_DIR="/path1,/path2"  # 分析多个目录
  claude-stats analyze --breakdown           # 聚合分析所有配置目录`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	// 添加命令特定的标志位
	analyzeCmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "输出格式 (table, json, csv)")
	analyzeCmd.Flags().StringVarP(&outputFile, "output", "o", "", "输出文件路径")
	analyzeCmd.Flags().StringVar(&startDate, "since", "", "开始日期 (YYYYMMDD)")
	analyzeCmd.Flags().StringVar(&endDate, "until", "", "结束日期 (YYYYMMDD)")
	analyzeCmd.Flags().StringVar(&modelFilter, "model", "", "过滤特定模型")
	analyzeCmd.Flags().BoolVarP(&showDetails, "details", "d", false, "显示详细信息")
	analyzeCmd.Flags().BoolVar(&noColor, "no-color", false, "禁用颜色输出")
	
	// 新增：增强功能标志位
	analyzeCmd.Flags().StringSliceVar(&configDirs, "config-dirs", []string{}, "指定多个Claude配置目录，逗号分隔")
	analyzeCmd.Flags().BoolVarP(&offline, "offline", "O", false, "离线模式，使用缓存定价数据")
	analyzeCmd.Flags().BoolVar(&breakdown, "breakdown", false, "显示按模型分解的详细成本")
	analyzeCmd.Flags().StringVar(&order, "order", "desc", "排序顺序 (asc, desc)")
	analyzeCmd.Flags().StringVar(&costMode, "mode", "auto", "成本计算模式 (auto, calculate, display)")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	// 确定要分析的目录列表
	targetDirs := getTargetDirectories(args)
	
	if verbose {
		fmt.Printf("🔍 分析目录: %s\n", strings.Join(targetDirs, ", "))
		if len(targetDirs) > 1 {
			fmt.Printf("💡 将聚合分析 %d 个配置目录的数据\n", len(targetDirs))
		}
	}

	// 检查目录是否存在
	for _, dir := range targetDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("⚠️  目录不存在，跳过: %s\n", dir)
		}
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
		
		if verbose {
			fmt.Printf("📂 正在处理目录: %s\n", targetDir)
		}
		
		stats, err := claudeParser.ParseDirectory(targetDir)
		if err != nil {
			fmt.Printf("⚠️  解析目录失败，跳过 %s: %v\n", targetDir, err)
			continue
		}

		// 聚合数据
		mergeUsageStats(aggregatedStats, stats)
		successfulDirs = append(successfulDirs, targetDir)
		
		if verbose {
			fmt.Printf("✅ 完成目录: %s (会话:%d, 消息:%d)\n", 
				targetDir, stats.TotalSessions, stats.TotalMessages)
		}
	}

	if len(successfulDirs) == 0 {
		return fmt.Errorf("没有找到有效的Claude配置目录")
	}

	if verbose && len(successfulDirs) > 1 {
		fmt.Printf("📊 聚合完成: 共处理 %d 个目录\n", len(successfulDirs))
	}

	// 应用模型过滤器
	if modelFilter != "" {
		aggregatedStats = filterByModel(aggregatedStats, modelFilter)
	}

	// 最终处理
	claudeParser.FinalizeStats(aggregatedStats)
	
	// 估算订阅限额信息
	if aggregatedStats.DetectedMode == "subscription" {
		aggregatedStats.SubscriptionQuota = aggregatedStats.EstimateSubscriptionQuota()
	}

	// 格式化并输出结果
	formatter := formatter.NewFormatter()
	formatter.ShowDetails = showDetails
	formatter.Verbose = verbose
	
	// 设置颜色选项
	if noColor {
		formatter.Colors.Enabled = false
	}

	output, err := formatter.Format(aggregatedStats, outputFormat)
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

// getTargetDirectories 获取目标目录列表，支持多配置目录
func getTargetDirectories(args []string) []string {
	// 优先级：命令行参数 > --config-dirs > CLAUDE_CONFIG_DIR环境变量 > 自动检测
	
	if len(args) > 0 {
		// 命令行指定了目录
		return []string{args[0]}
	}
	
	if len(configDirs) > 0 {
		// --config-dirs 参数指定了目录
		return configDirs
	}
	
	// 检查 CLAUDE_CONFIG_DIR 环境变量
	if envDirs := os.Getenv("CLAUDE_CONFIG_DIR"); envDirs != "" {
		// 支持逗号分隔的多个路径
		dirs := strings.Split(envDirs, ",")
		var validDirs []string
		for _, dir := range dirs {
			dir = strings.TrimSpace(dir)
			if dir != "" {
				// 展开路径中的 ~ 符号
				if strings.HasPrefix(dir, "~/") {
					if homeDir, err := os.UserHomeDir(); err == nil {
						dir = filepath.Join(homeDir, dir[2:])
					}
				}
				validDirs = append(validDirs, dir)
			}
		}
		if len(validDirs) > 0 {
			return validDirs
		}
	}
	
	// 自动检测默认目录
	return []string{getDefaultClaudeDirectory()}
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
		// Windows: 尝试多个可能的位置
		possibleDirs := []string{
			filepath.Join(homeDir, "AppData", "Roaming", "claude", "projects"),
			filepath.Join(homeDir, "AppData", "Local", "claude", "projects"),
			filepath.Join(homeDir, ".claude", "projects"),
		}
		for _, dir := range possibleDirs {
			if _, err := os.Stat(dir); err == nil {
				return dir
			}
		}
		claudeDir = possibleDirs[0] // 默认使用第一个
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

// mergeUsageStats 聚合多个UsageStats
func mergeUsageStats(target, source *models.UsageStats) {
	target.TotalMessages += source.TotalMessages
	target.TotalSessions += source.TotalSessions
	target.ParsedMessages += source.ParsedMessages
	target.ExtractedTokens += source.ExtractedTokens
	target.TotalTokens.Add(source.TotalTokens)

	// 合并模型统计
	for model, usage := range source.ModelStats {
		targetUsage := target.ModelStats[model]
		targetUsage.Add(usage)
		target.ModelStats[model] = targetUsage
	}

	// 合并日期统计
	for date, usage := range source.DailyStats {
		targetUsage := target.DailyStats[date]
		targetUsage.Add(usage)
		target.DailyStats[date] = targetUsage
	}

	// 合并会话统计（避免ID冲突）
	for sessionID, session := range source.SessionStats {
		if existing, exists := target.SessionStats[sessionID]; exists {
			// 如果会话ID冲突，合并会话数据
			existing.MessageCount += session.MessageCount
			existing.Tokens.Add(session.Tokens)
			if session.EndTime.After(existing.EndTime) {
				existing.EndTime = session.EndTime
			}
			if session.StartTime.Before(existing.StartTime) {
				existing.StartTime = session.StartTime
			}
			target.SessionStats[sessionID] = existing
		} else {
			target.SessionStats[sessionID] = session
		}
	}

	// 合并项目统计
	for projectKey, project := range source.ProjectStats {
		if existing, exists := target.ProjectStats[projectKey]; exists {
			existing.Tokens.Add(project.Tokens)
			if project.LastActivity.After(existing.LastActivity) {
				existing.LastActivity = project.LastActivity
			}
			target.ProjectStats[projectKey] = existing
		} else {
			target.ProjectStats[projectKey] = project
		}
	}

	// 合并消息类型统计
	for msgType, count := range source.MessageTypes {
		target.MessageTypes[msgType] += count
	}

	// 设置检测模式（优先使用subscription）
	if source.DetectedMode == "subscription" || target.DetectedMode == "" {
		target.DetectedMode = source.DetectedMode
	}

	// 合并分析时间段
	if !source.AnalysisPeriod.StartTime.IsZero() {
		if target.AnalysisPeriod.StartTime.IsZero() || source.AnalysisPeriod.StartTime.Before(target.AnalysisPeriod.StartTime) {
			target.AnalysisPeriod.StartTime = source.AnalysisPeriod.StartTime
		}
		if target.AnalysisPeriod.EndTime.IsZero() || source.AnalysisPeriod.EndTime.After(target.AnalysisPeriod.EndTime) {
			target.AnalysisPeriod.EndTime = source.AnalysisPeriod.EndTime
		}
		target.AnalysisPeriod.Duration = target.AnalysisPeriod.EndTime.Sub(target.AnalysisPeriod.StartTime).String()
	}
}

// createDateFilter 创建日期过滤器，支持多种日期格式
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

// parseDate 解析日期字符串，支持多种格式
func parseDate(dateStr string) (time.Time, error) {
	// 支持多种日期格式
	formats := []string{
		"20060102",      // YYYYMMDD (ccusage兼容格式)
		"2006-01-02",    // YYYY-MM-DD
		"2006/01/02",    // YYYY/MM/DD  
		"01/02/2006",    // MM/DD/YYYY
		"2006-01-02 15:04:05", // 包含时间
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("无法解析日期: %s，支持格式: YYYYMMDD, YYYY-MM-DD, YYYY/MM/DD", dateStr)
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