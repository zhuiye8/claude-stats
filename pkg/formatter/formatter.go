package formatter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/zhuiye8/claude-stats/pkg/models"
	"github.com/zhuiye8/claude-stats/pkg/parser"
)

// Formatter 用于格式化输出
type Formatter struct {
	ShowDetails bool
	Verbose     bool
	Colors      *ColorSettings
}

// NewFormatter 创建新的格式化器
func NewFormatter() *Formatter {
	return &Formatter{
		ShowDetails: false,
		Verbose:     false,
		Colors:      NewColorSettings(),
	}
}

// Format 格式化统计数据
func (f *Formatter) Format(stats *models.UsageStats, format string) (string, error) {
	switch strings.ToLower(format) {
	case "json":
		return f.formatJSON(stats)
	case "csv":
		return f.formatCSV(stats)
	case "table", "":
		return f.formatTable(stats)
	default:
		return "", fmt.Errorf("不支持的格式: %s", format)
	}
}

// formatJSON 格式化为JSON
func (f *Formatter) formatJSON(stats *models.UsageStats) (string, error) {
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// formatCSV 格式化为CSV
func (f *Formatter) formatCSV(stats *models.UsageStats) (string, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	// 写入标题行
	headers := []string{
		"类型", "名称", "输入Token", "输出Token", "缓存创建Token", 
		"缓存读取Token", "总Token", "估算成本(USD)",
	}
	if err := writer.Write(headers); err != nil {
		return "", err
	}

	// 写入总体统计
	totalRow := []string{
		"总计", "全部",
		fmt.Sprintf("%d", stats.TotalTokens.InputTokens),
		fmt.Sprintf("%d", stats.TotalTokens.OutputTokens),
		fmt.Sprintf("%d", stats.TotalTokens.CacheCreationTokens),
		fmt.Sprintf("%d", stats.TotalTokens.CacheReadTokens),
		fmt.Sprintf("%d", stats.TotalTokens.GetTotalTokens()),
		fmt.Sprintf("%.4f", stats.EstimatedCost.TotalCost),
	}
	if err := writer.Write(totalRow); err != nil {
		return "", err
	}

	// 写入模型统计
	for model, usage := range stats.ModelStats {
		row := []string{
			"模型", model,
			fmt.Sprintf("%d", usage.InputTokens),
			fmt.Sprintf("%d", usage.OutputTokens),
			fmt.Sprintf("%d", usage.CacheCreationTokens),
			fmt.Sprintf("%d", usage.CacheReadTokens),
			fmt.Sprintf("%d", usage.GetTotalTokens()),
			fmt.Sprintf("%.4f", stats.EstimatedCost.ModelCosts[model]),
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	// 写入日期统计
	if f.ShowDetails {
		var dates []string
		for date := range stats.DailyStats {
			dates = append(dates, date)
		}
		sort.Strings(dates)

		for _, date := range dates {
			usage := stats.DailyStats[date]
			row := []string{
				"日期", date,
				fmt.Sprintf("%d", usage.InputTokens),
				fmt.Sprintf("%d", usage.OutputTokens),
				fmt.Sprintf("%d", usage.CacheCreationTokens),
				fmt.Sprintf("%d", usage.CacheReadTokens),
				fmt.Sprintf("%d", usage.GetTotalTokens()),
				"", // 日期级别不显示成本
			}
			if err := writer.Write(row); err != nil {
				return "", err
			}
		}
	}

	writer.Flush()
	return builder.String(), writer.Error()
}

// formatTable 格式化为表格
func (f *Formatter) formatTable(stats *models.UsageStats) (string, error) {
	var output strings.Builder

	// 添加美化的标题
	f.writeHeader(&output, stats)

	// 添加装饰性分隔符
	output.WriteString(f.Colors.GradientSeparator(80) + "\n\n")

	// 基本信息
	f.writeBasicInfo(&output, stats)

	// 总体统计表格
	f.writeTotalStats(&output, stats)

	// 模型统计表格
	if len(stats.ModelStats) > 0 {
		f.writeModelStats(&output, stats)
	}

	// 成本分析
	f.writeCostAnalysis(&output, stats)

	// 详细信息
	if f.ShowDetails {
		f.writeDetailedStats(&output, stats)
	}

	return output.String(), nil
}

// writeHeader 写入美化的标题
func (f *Formatter) writeHeader(output *strings.Builder, stats *models.UsageStats) {
	// 主标题
	title := f.Colors.Header("🎯 Claude Code 使用统计报告")
	output.WriteString(fmt.Sprintf("%s\n", title))
	
	// 版本和时间信息
	subtitle := f.Colors.Dim(fmt.Sprintf("生成时间: %s", time.Now().Format("2006-01-02 15:04:05")))
	output.WriteString(fmt.Sprintf("%s\n", subtitle))
}

// writeBasicInfo 写入基本信息
func (f *Formatter) writeBasicInfo(output *strings.Builder, stats *models.UsageStats) {
	sectionTitle := f.Colors.IconHeader("📊", "基本信息", BrightBlue)
	output.WriteString(fmt.Sprintf("%s\n", sectionTitle))
	
	mode := getModeDisplay(stats.DetectedMode)
	modeColor := Green
	if stats.DetectedMode == "subscription" {
		modeColor = BrightMagenta
	}
	
	output.WriteString(fmt.Sprintf("   检测模式: %s\n", f.Colors.Colorize(mode, modeColor)))
	output.WriteString(fmt.Sprintf("   总会话数: %s\n", f.Colors.BrightYellow(formatNumber(stats.TotalSessions))))
	output.WriteString(fmt.Sprintf("   总消息数: %s\n", f.Colors.BrightCyan(formatNumber(stats.TotalMessages))))
	
	if !stats.AnalysisPeriod.StartTime.IsZero() {
		timeRange := fmt.Sprintf("%s 至 %s", 
			stats.AnalysisPeriod.StartTime.Format("2006-01-02 15:04"),
			stats.AnalysisPeriod.EndTime.Format("2006-01-02 15:04"))
		output.WriteString(fmt.Sprintf("   分析时段: %s\n", f.Colors.Info(timeRange)))
		output.WriteString(fmt.Sprintf("   持续时间: %s\n", f.Colors.Info(stats.AnalysisPeriod.Duration)))
	}
	output.WriteString("\n")
}

// writeTotalStats 写入总体统计
func (f *Formatter) writeTotalStats(output *strings.Builder, stats *models.UsageStats) {
	sectionTitle := f.Colors.IconHeader("📈", "Token 使用统计", BrightGreen)
	output.WriteString(fmt.Sprintf("%s\n", sectionTitle))
	
	t := table.NewWriter()
	t.AppendHeader(table.Row{
		f.Colors.Header("类型"), 
		f.Colors.Header("数量"), 
		f.Colors.Header("百分比"),
		f.Colors.Header("可视化"),
	})

	total := stats.TotalTokens.GetTotalTokens()
	if total > 0 {
		// 输入Token
		inputPct := float64(stats.TotalTokens.InputTokens) * 100 / float64(total)
		inputBar := f.Colors.ProgressBar(stats.TotalTokens.InputTokens, total, 20)
		t.AppendRow(table.Row{
			f.Colors.BrightBlue("输入Token"), 
			f.Colors.BrightYellow(formatNumber(stats.TotalTokens.InputTokens)),
			f.Colors.Cyan(fmt.Sprintf("%.1f%%", inputPct)),
			inputBar,
		})
		
		// 输出Token
		outputPct := float64(stats.TotalTokens.OutputTokens) * 100 / float64(total)
		outputBar := f.Colors.ProgressBar(stats.TotalTokens.OutputTokens, total, 20)
		t.AppendRow(table.Row{
			f.Colors.BrightGreen("输出Token"), 
			f.Colors.BrightYellow(formatNumber(stats.TotalTokens.OutputTokens)),
			f.Colors.Cyan(fmt.Sprintf("%.1f%%", outputPct)),
			outputBar,
		})
		
		if stats.TotalTokens.CacheCreationTokens > 0 {
			cachePct := float64(stats.TotalTokens.CacheCreationTokens) * 100 / float64(total)
			cacheBar := f.Colors.ProgressBar(stats.TotalTokens.CacheCreationTokens, total, 20)
			t.AppendRow(table.Row{
				f.Colors.BrightMagenta("缓存创建Token"), 
				f.Colors.BrightYellow(formatNumber(stats.TotalTokens.CacheCreationTokens)),
				f.Colors.Cyan(fmt.Sprintf("%.1f%%", cachePct)),
				cacheBar,
			})
		}
		
		if stats.TotalTokens.CacheReadTokens > 0 {
			readPct := float64(stats.TotalTokens.CacheReadTokens) * 100 / float64(total)
			readBar := f.Colors.ProgressBar(stats.TotalTokens.CacheReadTokens, total, 20)
			t.AppendRow(table.Row{
				f.Colors.BrightCyan("缓存读取Token"), 
				f.Colors.BrightYellow(formatNumber(stats.TotalTokens.CacheReadTokens)),
				f.Colors.Cyan(fmt.Sprintf("%.1f%%", readPct)),
				readBar,
			})
		}
	}

	t.AppendFooter(table.Row{
		f.Colors.Bold("总计"), 
		f.Colors.Bold(f.Colors.BrightYellow(formatNumber(total))), 
		f.Colors.Bold("100.0%"),
		f.Colors.Bold("📊"),
	})
	
	// 使用更美观的表格样式
	t.SetStyle(table.StyleColoredBright)
	t.Style().Options.SeparateRows = true

	output.WriteString(t.Render() + "\n\n")
}

// writeModelStats 写入模型统计
func (f *Formatter) writeModelStats(output *strings.Builder, stats *models.UsageStats) {
	t := table.NewWriter()
	t.SetTitle("🤖 按模型统计")
	t.AppendHeader(table.Row{"模型", "输入", "输出", "缓存", "总计", "成本(USD)"})

	// 按总token数排序
	type modelStat struct {
		name  string
		usage models.TokenUsage
		cost  float64
	}
	
	var modelStats []modelStat
	for model, usage := range stats.ModelStats {
		cost := stats.EstimatedCost.ModelCosts[model]
		modelStats = append(modelStats, modelStat{model, usage, cost})
	}
	
	sort.Slice(modelStats, func(i, j int) bool {
		return modelStats[i].usage.GetTotalTokens() > modelStats[j].usage.GetTotalTokens()
	})

	for _, ms := range modelStats {
		cacheTotal := ms.usage.CacheCreationTokens + ms.usage.CacheReadTokens
		t.AppendRow(table.Row{
			ms.name,
			formatNumber(ms.usage.InputTokens),
			formatNumber(ms.usage.OutputTokens),
			formatNumber(cacheTotal),
			formatNumber(ms.usage.GetTotalTokens()),
			fmt.Sprintf("$%.4f", ms.cost),
		})
	}

	t.SetStyle(table.StyleColoredBright)
	output.WriteString(t.Render() + "\n\n")
}

// writeCostAnalysis 写入成本分析
func (f *Formatter) writeCostAnalysis(output *strings.Builder, stats *models.UsageStats) {
	cost := stats.EstimatedCost
	
	sectionTitle := f.Colors.IconHeader("💰", "成本分析", BrightYellow)
	output.WriteString(fmt.Sprintf("%s\n", sectionTitle))
	
	if cost.IsEstimated {
		note := f.Colors.Dim("(基于订阅模式的API等价成本估算)")
		output.WriteString(fmt.Sprintf("   %s\n", note))
	}
	
	// 使用颜色区分不同类型的成本
	output.WriteString(fmt.Sprintf("   输入成本:     %s\n", f.Colors.BrightGreen(fmt.Sprintf("$%.4f", cost.InputCost))))
	output.WriteString(fmt.Sprintf("   输出成本:     %s\n", f.Colors.BrightBlue(fmt.Sprintf("$%.4f", cost.OutputCost))))
	
	if cost.CacheCreationCost > 0 {
		output.WriteString(fmt.Sprintf("   缓存创建成本: %s\n", f.Colors.BrightMagenta(fmt.Sprintf("$%.4f", cost.CacheCreationCost))))
	}
	if cost.CacheReadCost > 0 {
		output.WriteString(fmt.Sprintf("   缓存读取成本: %s\n", f.Colors.BrightCyan(fmt.Sprintf("$%.4f", cost.CacheReadCost))))
	}
	
	// 总成本用醒目颜色
	totalCostStr := fmt.Sprintf("$%.4f", cost.TotalCost)
	if cost.TotalCost > 20.0 {
		totalCostStr = f.Colors.BrightRed(totalCostStr)
	} else if cost.TotalCost > 5.0 {
		totalCostStr = f.Colors.BrightYellow(totalCostStr)
	} else {
		totalCostStr = f.Colors.BrightGreen(totalCostStr)
	}
	output.WriteString(fmt.Sprintf("   总成本:       %s\n", totalCostStr))

	// 订阅模式建议
	if stats.DetectedMode == "subscription" {
		costCalculator := parser.NewCostCalculator()
		analysis := costCalculator.GetSubscriptionEquivalent(stats)
		
		output.WriteString("\n")
		suggestionTitle := f.Colors.IconHeader("🎯", "订阅计划建议", BrightMagenta)
		output.WriteString(fmt.Sprintf("%s\n", suggestionTitle))
		
		planColor := BrightGreen
		if strings.Contains(analysis.RecommendedPlan, "Max") {
			planColor = BrightYellow
		}
		output.WriteString(fmt.Sprintf("   建议计划: %s\n", f.Colors.Colorize(analysis.RecommendedPlan, planColor)))
		
		if analysis.MonthlySavings > 0 {
			savingsStr := f.Colors.Success(fmt.Sprintf("$%.2f/月", analysis.MonthlySavings))
			output.WriteString(fmt.Sprintf("   预估节省: %s\n", savingsStr))
		}
	}
	
	output.WriteString("\n")
}

// writeDetailedStats 写入详细统计
func (f *Formatter) writeDetailedStats(output *strings.Builder, stats *models.UsageStats) {
	// 按日期统计
	if len(stats.DailyStats) > 0 {
		f.writeDailyStats(output, stats)
	}

	// 会话统计
	if len(stats.SessionStats) > 0 && len(stats.SessionStats) <= 20 {
		f.writeSessionStats(output, stats)
	}
}

// writeDailyStats 写入每日统计
func (f *Formatter) writeDailyStats(output *strings.Builder, stats *models.UsageStats) {
	t := table.NewWriter()
	t.SetTitle("📅 每日使用统计")
	t.AppendHeader(table.Row{"日期", "输入", "输出", "缓存", "总计"})

	// 按日期排序
	var dates []string
	for date := range stats.DailyStats {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	for _, date := range dates {
		usage := stats.DailyStats[date]
		cacheTotal := usage.CacheCreationTokens + usage.CacheReadTokens
		t.AppendRow(table.Row{
			date,
			formatNumber(usage.InputTokens),
			formatNumber(usage.OutputTokens),
			formatNumber(cacheTotal),
			formatNumber(usage.GetTotalTokens()),
		})
	}

	t.SetStyle(table.StyleColoredBright)
	output.WriteString(t.Render() + "\n\n")
}

// writeSessionStats 写入会话统计
func (f *Formatter) writeSessionStats(output *strings.Builder, stats *models.UsageStats) {
	t := table.NewWriter()
	t.SetTitle("💬 会话统计 (最近20个)")
	t.AppendHeader(table.Row{"会话ID", "开始时间", "消息数", "Token数", "模型"})

	// 按开始时间排序
	type sessionInfo struct {
		id   string
		info models.SessionInfo
	}
	
	var sessions []sessionInfo
	for id, info := range stats.SessionStats {
		sessions = append(sessions, sessionInfo{id, info})
	}
	
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].info.StartTime.After(sessions[j].info.StartTime)
	})

	// 显示最近的20个会话
	maxSessions := 20
	if len(sessions) < maxSessions {
		maxSessions = len(sessions)
	}

	for i := 0; i < maxSessions; i++ {
		s := sessions[i]
		t.AppendRow(table.Row{
			s.id[:8] + "...", // 显示会话ID的前8位
			s.info.StartTime.Format("01-02 15:04"),
			s.info.MessageCount,
			formatNumber(s.info.Tokens.GetTotalTokens()),
			s.info.Model,
		})
	}

	t.SetStyle(table.StyleColoredBright)
	output.WriteString(t.Render() + "\n\n")
}

// 辅助函数

// getModeDisplay 获取模式显示文本
func getModeDisplay(mode string) string {
	switch mode {
	case "api":
		return "API模式 (按token计费)"
	case "subscription":
		return "订阅模式 (按请求限制)"
	default:
		return mode
	}
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