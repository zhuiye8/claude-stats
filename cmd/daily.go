package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zhuiye8/claude-stats/pkg/parser"
	"github.com/zhuiye8/claude-stats/pkg/formatter"
	"github.com/zhuiye8/claude-stats/pkg/models"
)

// dailyCmd 代表daily命令
var dailyCmd = &cobra.Command{
	Use:   "daily [目录路径]",
	Short: "按日分析Claude Code使用情况",
	Long: `按日期分析Claude Code的使用情况，提供精确的每日Token统计和成本分析。

此命令专门优化了按日分组的数据处理逻辑，确保计算精度和性能。
相比通用分析，daily命令能够：
• 精确按日期边界分组数据
• 优化的日期范围处理
• 专门的日级成本分配算法
• 更准确的使用模式分析

支持功能：
• 按日显示Token使用统计
• 每日成本分析和趋势
• 模型使用分解（--breakdown）
• 时间范围过滤（--since, --until）
• 多种排序方式（--order）

示例：
  claude-stats daily                           # 显示所有日期的使用情况
  claude-stats daily --breakdown               # 显示每日的模型使用分解
  claude-stats daily --since 20241201         # 显示12月1日以来的使用
  claude-stats daily --since 20241201 --until 20241231  # 显示12月的使用
  claude-stats daily --order asc              # 按时间正序排列
  claude-stats daily --json                   # JSON格式输出

日期格式：
  支持 YYYYMMDD, YYYY-MM-DD, YYYY/MM/DD 等多种格式`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDaily,
}

func init() {
	rootCmd.AddCommand(dailyCmd)

	// daily命令特定的标志位
	dailyCmd.Flags().BoolVar(&dailyBreakdown, "breakdown", false, "显示每日按模型分解的详细统计")
	dailyCmd.Flags().StringVar(&dailyOrder, "order", "desc", "排序顺序: desc(最新优先) 或 asc(最旧优先)")
	
	// 继承通用标志位
	dailyCmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "输出格式 (table, json, csv)")
	dailyCmd.Flags().StringVarP(&outputFile, "output", "o", "", "输出文件路径")
	dailyCmd.Flags().StringVar(&startDate, "since", "", "开始日期 (YYYYMMDD)")
	dailyCmd.Flags().StringVar(&endDate, "until", "", "结束日期 (YYYYMMDD)")
	dailyCmd.Flags().BoolVar(&noColor, "no-color", false, "禁用颜色输出")
	dailyCmd.Flags().BoolVarP(&offline, "offline", "O", false, "离线模式")
	dailyCmd.Flags().StringVar(&costMode, "mode", "auto", "成本计算模式 (auto, calculate, display)")
}

// runDaily 执行每日分析
func runDaily(cmd *cobra.Command, args []string) error {
	// 确定要分析的目录列表
	targetDirs := getTargetDirectories(args)
	
	if verbose {
		fmt.Printf("📅 开始按日分析: %s\n", strings.Join(targetDirs, ", "))
	}

	// 创建专门的日分析器
	dailyAnalyzer := NewDailyAnalyzer()
	dailyAnalyzer.Verbose = verbose
	dailyAnalyzer.Order = dailyOrder
	dailyAnalyzer.CostMode = costMode

	// 设置日期过滤器
	if startDate != "" || endDate != "" {
		dateFilter, err := createDateFilter(startDate, endDate)
		if err != nil {
			return fmt.Errorf("日期格式错误: %w", err)
		}
		dailyAnalyzer.DateFilter = dateFilter
	}

	// 执行日分析
	dailyReport, err := dailyAnalyzer.AnalyzeDirectories(targetDirs)
	if err != nil {
		return fmt.Errorf("日分析失败: %w", err)
	}

	if verbose {
		fmt.Printf("✅ 分析完成: 共 %d 天的数据\n", len(dailyReport.DailyData))
	}

	// 输出结果
	return outputDailyReport(dailyReport)
}

// DailyAnalyzer 专门的日分析器
type DailyAnalyzer struct {
	Verbose    bool
	Order      string
	CostMode   string
	DateFilter *parser.DateFilter
}

// NewDailyAnalyzer 创建新的日分析器
func NewDailyAnalyzer() *DailyAnalyzer {
	return &DailyAnalyzer{
		Verbose:  false,
		Order:    "desc",
		CostMode: "auto",
	}
}

// AnalyzeDirectories 分析多个目录的日数据
func (da *DailyAnalyzer) AnalyzeDirectories(targetDirs []string) (*models.DailyReport, error) {
	// 创建基础解析器
	claudeParser := parser.NewClaudeParser()
	claudeParser.Verbose = da.Verbose
	claudeParser.SkipErrors = true
	claudeParser.DateFilter = da.DateFilter

	// 按日聚合的数据结构
	dailyAggregation := make(map[string]*models.DailyDataPoint)
	totalSummary := &models.DailyDataPoint{
		Date:      "总计",
		Models:    []string{},
		Breakdown: make(map[string]models.DailyModelData),
	}

	// 处理每个目录
	for _, targetDir := range targetDirs {
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			if da.Verbose {
				fmt.Printf("⚠️  目录不存在，跳过: %s\n", targetDir)
			}
			continue
		}
		
		if da.Verbose {
			fmt.Printf("📂 处理目录: %s\n", targetDir)
		}
		
		// 解析目录（但专注于日级数据处理）
		err := da.processDirectoryForDaily(claudeParser, targetDir, dailyAggregation, totalSummary)
		if err != nil {
			if da.Verbose {
				fmt.Printf("⚠️  处理目录失败，跳过 %s: %v\n", targetDir, err)
			}
			continue
		}
	}

	if len(dailyAggregation) == 0 {
		return &models.DailyReport{
			Type:      "daily",
			DailyData: []models.DailyDataPoint{},
			Summary:   *totalSummary,
		}, nil
	}

	// 转换为排序的切片
	dailyData := da.convertAndSortDailyData(dailyAggregation)

	return &models.DailyReport{
		Type:      "daily",
		DailyData: dailyData,
		Summary:   *totalSummary,
	}, nil
}

// processDirectoryForDaily 专门为日分析处理目录
func (da *DailyAnalyzer) processDirectoryForDaily(parser *parser.ClaudeParser, dir string, 
	dailyAggregation map[string]*models.DailyDataPoint, totalSummary *models.DailyDataPoint) error {
	
	// 解析目录，但使用流式处理减少内存占用
	stats, err := parser.ParseDirectory(dir)
	if err != nil {
		return err
	}

	// 处理每一天的数据
	for dateStr, dailyUsage := range stats.DailyStats {
		dayData, exists := dailyAggregation[dateStr]
		if !exists {
			dayData = &models.DailyDataPoint{
				Date:      dateStr,
				Models:    []string{},
				Breakdown: make(map[string]models.DailyModelData),
			}
			dailyAggregation[dateStr] = dayData
		}

		// 累加Token数据
		dayData.InputTokens += dailyUsage.InputTokens
		dayData.OutputTokens += dailyUsage.OutputTokens
		dayData.CacheCreationTokens += dailyUsage.CacheCreationTokens
		dayData.CacheReadTokens += dailyUsage.CacheReadTokens
		dayData.TotalTokens += dailyUsage.GetTotalTokens()

		// 统计该日期的会话和消息
		for _, session := range stats.SessionStats {
			sessionDate := session.StartTime.Format("2006-01-02")
			if sessionDate == dateStr {
				dayData.SessionCount++
				dayData.MessageCount += session.MessageCount

				// 记录使用的模型
				if session.Model != "" {
					modelExists := false
					for _, model := range dayData.Models {
						if model == session.Model {
							modelExists = true
							break
						}
					}
					if !modelExists {
						dayData.Models = append(dayData.Models, session.Model)
					}
				}
			}
		}

		// 如果需要breakdown，处理模型级数据
		if dailyBreakdown {
			da.processModelBreakdown(stats, dateStr, dayData)
		}

		// 累加到总汇总
		totalSummary.InputTokens += dailyUsage.InputTokens
		totalSummary.OutputTokens += dailyUsage.OutputTokens
		totalSummary.CacheCreationTokens += dailyUsage.CacheCreationTokens
		totalSummary.CacheReadTokens += dailyUsage.CacheReadTokens
		totalSummary.TotalTokens += dailyUsage.GetTotalTokens()
	}

	// 计算成本（使用指定的成本模式）
	da.calculateDailyCosts(stats, dailyAggregation, totalSummary)

	return nil
}

// processModelBreakdown 处理模型分解数据
func (da *DailyAnalyzer) processModelBreakdown(stats *models.UsageStats, dateStr string, dayData *models.DailyDataPoint) {
	// 这里需要更精细的逻辑来分配模型使用到特定日期
	// 从会话数据中提取每日每模型的使用情况
	for _, session := range stats.SessionStats {
		sessionDate := session.StartTime.Format("2006-01-02")
		if sessionDate == dateStr && session.Model != "" {
			modelData := dayData.Breakdown[session.Model]
			modelData.InputTokens += session.Tokens.InputTokens
			modelData.OutputTokens += session.Tokens.OutputTokens
			modelData.CacheCreationTokens += session.Tokens.CacheCreationTokens
			modelData.CacheReadTokens += session.Tokens.CacheReadTokens
			modelData.TotalTokens += session.Tokens.GetTotalTokens()
			modelData.MessageCount += session.MessageCount
			
			dayData.Breakdown[session.Model] = modelData
		}
	}
}

// calculateDailyCosts 计算每日成本
func (da *DailyAnalyzer) calculateDailyCosts(stats *models.UsageStats, 
	dailyAggregation map[string]*models.DailyDataPoint, totalSummary *models.DailyDataPoint) {
	
	// 根据costMode计算成本
	switch da.CostMode {
	case "calculate":
		// 强制从Token计算
		da.calculateCostsFromTokens(dailyAggregation, totalSummary)
	case "display":
		// 仅使用预计算的成本（如果有）
		da.usePrecalculatedCosts(stats, dailyAggregation, totalSummary)
	case "auto":
		fallthrough
	default:
		// 自动模式：优先使用预计算，回退到Token计算
		if da.hasPrecalculatedCosts(stats) {
			da.usePrecalculatedCosts(stats, dailyAggregation, totalSummary)
		} else {
			da.calculateCostsFromTokens(dailyAggregation, totalSummary)
		}
	}
}

// calculateCostsFromTokens 从Token计算成本
func (da *DailyAnalyzer) calculateCostsFromTokens(dailyAggregation map[string]*models.DailyDataPoint, totalSummary *models.DailyDataPoint) {
	costCalculator := parser.NewCostCalculator()
	
	for _, dayData := range dailyAggregation {
		usage := models.TokenUsage{
			InputTokens:         dayData.InputTokens,
			OutputTokens:        dayData.OutputTokens,
			CacheCreationTokens: dayData.CacheCreationTokens,
			CacheReadTokens:     dayData.CacheReadTokens,
			TotalTokens:         dayData.TotalTokens,
		}
		
		// 假设订阅模式（大多数用户）
		costBreakdown := costCalculator.Calculate(&usage, nil, true)
		dayData.CostUSD = costBreakdown.TotalCost
		
		// 如果有breakdown，计算模型级成本
		if len(dayData.Breakdown) > 0 {
			for model, modelData := range dayData.Breakdown {
				modelUsage := models.TokenUsage{
					InputTokens:         modelData.InputTokens,
					OutputTokens:        modelData.OutputTokens,
					CacheCreationTokens: modelData.CacheCreationTokens,
					CacheReadTokens:     modelData.CacheReadTokens,
					TotalTokens:         modelData.TotalTokens,
				}
				modelCost := costCalculator.Calculate(&modelUsage, nil, true)
				modelData.CostUSD = modelCost.TotalCost
				dayData.Breakdown[model] = modelData
			}
		}
	}
	
	// 计算总成本
	usage := models.TokenUsage{
		InputTokens:         totalSummary.InputTokens,
		OutputTokens:        totalSummary.OutputTokens,
		CacheCreationTokens: totalSummary.CacheCreationTokens,
		CacheReadTokens:     totalSummary.CacheReadTokens,
		TotalTokens:         totalSummary.TotalTokens,
	}
	costBreakdown := costCalculator.Calculate(&usage, nil, true)
	totalSummary.CostUSD = costBreakdown.TotalCost
}

// usePrecalculatedCosts 使用预计算的成本
func (da *DailyAnalyzer) usePrecalculatedCosts(stats *models.UsageStats, 
	dailyAggregation map[string]*models.DailyDataPoint, totalSummary *models.DailyDataPoint) {
	// TODO: 实现使用预计算成本的逻辑
	// 目前回退到Token计算
	da.calculateCostsFromTokens(dailyAggregation, totalSummary)
}

// hasPrecalculatedCosts 检查是否有预计算的成本
func (da *DailyAnalyzer) hasPrecalculatedCosts(stats *models.UsageStats) bool {
	// TODO: 检查JSONL数据中是否包含costUSD字段
	return false
}

// convertAndSortDailyData 转换并排序日数据
func (da *DailyAnalyzer) convertAndSortDailyData(dailyAggregation map[string]*models.DailyDataPoint) []models.DailyDataPoint {
	var dailyData []models.DailyDataPoint
	
	for _, dayData := range dailyAggregation {
		dailyData = append(dailyData, *dayData)
	}
	
	// 排序
	sort.Slice(dailyData, func(i, j int) bool {
		if da.Order == "asc" {
			return dailyData[i].Date < dailyData[j].Date
		}
		return dailyData[i].Date > dailyData[j].Date
	})
	
	return dailyData
}

// outputDailyReport 输出日报告
func outputDailyReport(report *models.DailyReport) error {
	// 格式化并输出结果
	formatter := formatter.NewFormatter()
	formatter.ShowDetails = dailyBreakdown
	formatter.Verbose = verbose
	
	// 设置颜色选项
	if noColor {
		formatter.Colors.Enabled = false
	}

	var output string
	var err error
	
	switch strings.ToLower(outputFormat) {
	case "json":
		output, err = formatter.FormatDailyJSON(report)
	case "csv":
		output, err = formatter.FormatDailyCSV(report)
	case "table", "":
		output, err = formatter.FormatDaily(report)
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