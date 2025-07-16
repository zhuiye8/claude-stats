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

	// 订阅限额信息 (仅订阅模式显示)
	if stats.SubscriptionQuota != nil {
		f.writeSubscriptionQuota(&output, stats)
	}

	// 详细信息
	if f.ShowDetails {
		f.writeDetailedStats(&output, stats)
	}

	return output.String(), nil
}

// writeHeader 写入美化的标题
func (f *Formatter) writeHeader(output *strings.Builder, stats *models.UsageStats) {
	// Claude Stats ASCII 艺术字
	output.WriteString("\n")
	output.WriteString(f.Colors.BrightMagenta("  ██████╗██╗      █████╗ ██╗   ██╗██████╗ ███████╗\n"))
	output.WriteString(f.Colors.BrightCyan("██╔════╝██║     ██╔══██╗██║   ██║██╔══██╗██╔════╝\n"))
	output.WriteString(f.Colors.BrightBlue("██║     ██║     ███████║██║   ██║██║  ██║█████╗  \n"))
	output.WriteString(f.Colors.BrightGreen("██║     ██║     ██╔══██║██║   ██║██║  ██║██╔══╝  \n"))
	output.WriteString(f.Colors.BrightYellow("╚██████╗███████╗██║  ██║╚██████╔╝██████╔╝███████╗\n"))
	output.WriteString(f.Colors.BrightRed(" ╚═════╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═════╝ ╚══════╝\n"))
	
	output.WriteString(f.Colors.BrightMagenta("███████╗████████╗ █████╗ ████████╗███████╗\n"))
	output.WriteString(f.Colors.BrightCyan("██╔════╝╚══██╔══╝██╔══██╗╚══██╔══╝██╔════╝\n"))
	output.WriteString(f.Colors.BrightBlue("███████╗   ██║   ███████║   ██║   ███████╗\n"))
	output.WriteString(f.Colors.BrightGreen("╚════██║   ██║   ██╔══██║   ██║   ╚════██║\n"))
	output.WriteString(f.Colors.BrightYellow("███████║   ██║   ██║  ██║   ██║   ███████║\n"))
	output.WriteString(f.Colors.BrightRed("╚══════╝   ╚═╝   ╚═╝  ╚═╝   ╚═╝   ╚══════╝\n"))
	
	// 右下角作者信息
	authorInfo := fmt.Sprintf("%s%s", 
		strings.Repeat(" ", 45), 
		f.Colors.Dim("作者: zhuiye"))
	output.WriteString(authorInfo + "\n\n")
	
	// 副标题信息
	subtitle := f.Colors.Info(fmt.Sprintf("🎯 Claude Code 使用统计分析工具  •  生成时间: %s", 
		time.Now().Format("2006-01-02 15:04:05")))
	output.WriteString(fmt.Sprintf("%s\n\n", subtitle))
	
	// 装饰性分隔线
	separator := f.Colors.BrightBlue("═") + f.Colors.BrightCyan("═") + f.Colors.BrightGreen("═") + f.Colors.BrightYellow("═")
	fullSeparator := strings.Repeat(separator, 20)
	output.WriteString(fullSeparator + "\n\n")
}

// writeBasicInfo 写入基本信息
func (f *Formatter) writeBasicInfo(output *strings.Builder, stats *models.UsageStats) {
	sectionTitle := f.Colors.IconHeader("📋", "基本信息", BrightBlue)
	output.WriteString(fmt.Sprintf("%s\n", sectionTitle))
	
	mode := getModeDisplay(stats.DetectedMode)
	modeColor := Green
	modeIcon := "🔧"
	if stats.DetectedMode == "subscription" {
		modeColor = BrightMagenta
		modeIcon = "💎"
	}
	
	output.WriteString(fmt.Sprintf("   %s 检测模式: %s\n", modeIcon, f.Colors.Colorize(mode, modeColor)))
	output.WriteString(fmt.Sprintf("   📊 总会话数: %s\n", f.Colors.BrightYellow(formatNumber(stats.TotalSessions))))
	output.WriteString(fmt.Sprintf("   💬 总消息数: %s\n", f.Colors.BrightCyan(formatNumber(stats.TotalMessages))))
	
	// Claude Code 特定信息
	if stats.ParsedMessages > 0 {
		parseRate := float64(stats.ParsedMessages) * 100 / float64(stats.TotalMessages)
		output.WriteString(fmt.Sprintf("   ✅ 解析成功: %s (%s)\n", 
			f.Colors.BrightGreen(formatNumber(stats.ParsedMessages)),
			f.Colors.Cyan(fmt.Sprintf("%.1f%%", parseRate))))
	}
	
	if stats.ExtractedTokens > 0 {
		extractRate := float64(stats.ExtractedTokens) * 100 / float64(stats.TotalMessages)
		output.WriteString(fmt.Sprintf("   🎯 Token提取: %s (%s)\n", 
			f.Colors.BrightGreen(formatNumber(stats.ExtractedTokens)),
			f.Colors.Cyan(fmt.Sprintf("%.1f%%", extractRate))))
	}
	
	if !stats.AnalysisPeriod.StartTime.IsZero() {
		timeRange := fmt.Sprintf("%s 至 %s", 
			stats.AnalysisPeriod.StartTime.Format("2006-01-02 15:04"),
			stats.AnalysisPeriod.EndTime.Format("2006-01-02 15:04"))
		output.WriteString(fmt.Sprintf("   📅 分析时段: %s\n", f.Colors.Info(timeRange)))
		output.WriteString(fmt.Sprintf("   ⏰ 持续时间: %s\n", f.Colors.Info(stats.AnalysisPeriod.Duration)))
	}
	
	// 显示消息类型分布
	if len(stats.MessageTypes) > 0 {
		output.WriteString(fmt.Sprintf("   🏷️  消息类型: %s\n", f.formatMessageTypes(stats.MessageTypes)))
	}
	
	output.WriteString("\n")
}

// formatMessageTypes 格式化消息类型统计
func (f *Formatter) formatMessageTypes(messageTypes map[string]int) string {
	var parts []string
	for msgType, count := range messageTypes {
		color := BrightBlue
		switch msgType {
		case "user":
			color = BrightGreen
		case "assistant":
			color = BrightMagenta
		case "summary":
			color = BrightYellow
		}
		parts = append(parts, f.Colors.Colorize(fmt.Sprintf("%s:%d", msgType, count), color))
	}
	return strings.Join(parts, ", ")
}

// writeTotalStats 写入总体统计
func (f *Formatter) writeTotalStats(output *strings.Builder, stats *models.UsageStats) {
	sectionTitle := f.Colors.IconHeader("📈", "Token 使用统计", BrightGreen)
	output.WriteString(fmt.Sprintf("%s\n", sectionTitle))
	
	// 添加百分比计算说明
	baseTokens := stats.TotalTokens.InputTokens + stats.TotalTokens.OutputTokens
	output.WriteString(fmt.Sprintf("   💡 %s\n\n", 
		f.Colors.Dim(fmt.Sprintf("百分比基准: 基础Token(%s) = 输入Token + 输出Token", formatNumber(baseTokens)))))
	
	t := table.NewWriter()
	t.AppendHeader(table.Row{
		f.Colors.Header("类型"), 
		f.Colors.Header("数量"), 
		f.Colors.Header("百分比"),
		f.Colors.Header("说明"),
	})

	if baseTokens > 0 {
		// 输入Token - 基于基础Token计算百分比
		inputPct := float64(stats.TotalTokens.InputTokens) * 100 / float64(baseTokens)
		inputDesc := fmt.Sprintf("📥 %s", f.Colors.Info("用户提问成本"))
		t.AppendRow(table.Row{
			f.Colors.BrightBlue("输入Token"), 
			f.Colors.BrightYellow(formatNumber(stats.TotalTokens.InputTokens)),
			f.Colors.BrightGreen(fmt.Sprintf("%.1f%%", inputPct)),
			inputDesc,
		})
		
		// 输出Token - 基于基础Token计算百分比
		outputPct := float64(stats.TotalTokens.OutputTokens) * 100 / float64(baseTokens)
		outputDesc := fmt.Sprintf("📤 %s", f.Colors.Info("AI回复成本"))
		t.AppendRow(table.Row{
			f.Colors.BrightGreen("输出Token"), 
			f.Colors.BrightYellow(formatNumber(stats.TotalTokens.OutputTokens)),
			f.Colors.BrightGreen(fmt.Sprintf("%.1f%%", outputPct)),
			outputDesc,
		})
		
		// 添加分隔线
		t.AppendSeparator()
		
		// 缓存Token单独显示，不参与百分比计算
		if stats.TotalTokens.CacheCreationTokens > 0 {
			cacheDesc := fmt.Sprintf("📦 %s", f.Colors.Success("上下文缓存"))
			t.AppendRow(table.Row{
				f.Colors.BrightMagenta("缓存创建Token"), 
				f.Colors.BrightYellow(formatNumber(stats.TotalTokens.CacheCreationTokens)),
				f.Colors.Dim("不计入%"),
				cacheDesc,
			})
		}
		
		if stats.TotalTokens.CacheReadTokens > 0 {
			readDesc := fmt.Sprintf("⚡ %s", f.Colors.Success("缓存加速"))
			t.AppendRow(table.Row{
				f.Colors.BrightCyan("缓存读取Token"), 
				f.Colors.BrightYellow(formatNumber(stats.TotalTokens.CacheReadTokens)),
				f.Colors.Dim("不计入%"),
				readDesc,
			})
		}
	} else {
		// 如果没有token数据，显示提示信息
		t.AppendRow(table.Row{
			f.Colors.Warning("⚠️ 暂无Token数据"), 
			f.Colors.Dim("请检查JSONL格式"),
			f.Colors.Dim("--"),
			f.Colors.Dim("--"),
		})
	}

	// 总计行显示基础Token
	t.AppendFooter(table.Row{
		f.Colors.Bold("💰 基础Token总计"), 
		f.Colors.Bold(f.Colors.BrightYellow(formatNumber(baseTokens))), 
		f.Colors.Bold("100.0%"),
		f.Colors.Bold("💡 真实使用成本"),
	})
	
	// 使用更美观的表格样式
	t.SetStyle(table.StyleColoredBright)
	t.Style().Options.SeparateRows = true
	t.Style().Options.DrawBorder = true

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
		output.WriteString("\n")
		suggestionTitle := f.Colors.IconHeader("🎯", "订阅使用建议", BrightMagenta)
		output.WriteString(fmt.Sprintf("%s\n", suggestionTitle))
		
		output.WriteString(fmt.Sprintf("   ⚠️  %s\n", 
			f.Colors.Warning("以下建议基于估算数据，请结合实际使用情况判断")))
		output.WriteString("\n")
		
		if stats.SubscriptionQuota != nil {
			quota := stats.SubscriptionQuota
			
			// 通用使用建议
			output.WriteString(fmt.Sprintf("   💡 %s\n", 
				f.Colors.Info("效率提升技巧：")))
			output.WriteString(fmt.Sprintf("      • 使用 %s 清理上下文\n", f.Colors.BrightCyan("/compact")))
			output.WriteString(fmt.Sprintf("      • 使用 %s 重置对话\n", f.Colors.BrightCyan("/clear")))
			output.WriteString(fmt.Sprintf("      • 使用 %s 查看实时限额\n", f.Colors.BrightCyan("/status")))
			
			// 基于使用率的建议
			if quota.UsagePercentage > 80 {
				output.WriteString(fmt.Sprintf("\n   ⚠️  %s\n", 
					f.Colors.Warning("当前窗口使用率较高")))
				output.WriteString(fmt.Sprintf("      • 考虑使用 %s 减少上下文\n", f.Colors.BrightCyan("/compact")))
				output.WriteString(fmt.Sprintf("      • 避免长时间连续对话\n"))
				if quota.Plan == "Pro" {
					output.WriteString(fmt.Sprintf("      • 如经常遇到限制，可考虑升级计划\n"))
				}
			} else if quota.UsagePercentage < 20 {
				output.WriteString(fmt.Sprintf("\n   ✅ %s\n", 
					f.Colors.Success("当前窗口使用充裕")))
				output.WriteString(fmt.Sprintf("      • 当前计划适合您的使用模式\n"))
			}
			
			// 时区相关建议
			output.WriteString(fmt.Sprintf("\n   🌍 %s\n", 
				f.Colors.Info("时区注意事项：")))
			output.WriteString(fmt.Sprintf("      • Claude Code限额可能基于UTC时区\n"))
			output.WriteString(fmt.Sprintf("      • 重置时间可能与您的本地时间不同\n"))
			output.WriteString(fmt.Sprintf("      • 建议使用 %s 确认准确时间\n", f.Colors.BrightCyan("/status")))
			
			// 成本效益信息
			apiEquivalentCost := cost.TotalCost
			planCost := getPlanCostFloat(quota.Plan)
			if apiEquivalentCost > planCost {
				savings := apiEquivalentCost - planCost
				output.WriteString(fmt.Sprintf("\n   💰 %s\n", 
					f.Colors.Success(fmt.Sprintf("订阅模式相比API节省约 $%.2f/月", savings))))
			}
		} else {
			// 兜底建议
			output.WriteString(fmt.Sprintf("   💡 %s\n", f.Colors.Info("建议在Claude Code中使用 /status 查看详细限额信息")))
			output.WriteString(fmt.Sprintf("   📚 %s\n", f.Colors.Dim("参考官方文档了解订阅计划详情")))
		}
	}
	
	output.WriteString("\n")
}

// writeSubscriptionQuota 写入订阅限额信息
func (f *Formatter) writeSubscriptionQuota(output *strings.Builder, stats *models.UsageStats) {
	quota := stats.SubscriptionQuota
	sectionTitle := f.Colors.IconHeader("⚙️", "订阅限额状态", BrightMagenta)
	output.WriteString(fmt.Sprintf("%s\n", sectionTitle))
	
	// 数据准确性警告
	output.WriteString(fmt.Sprintf("   ⚠️  %s\n", 
		f.Colors.Warning("注意：以下数据为估算值，实际限额请使用Claude Code的 /status 命令查看")))
	output.WriteString(fmt.Sprintf("   💡 %s\n\n", 
		f.Colors.Dim("在Claude Code中运行 /status 可获取准确的当前窗口使用情况")))
	
	// 计划信息卡片
	planColor := BrightBlue
	planEmoji := "🥉"
	if quota.Plan == "Max20x" {
		planColor = BrightMagenta
		planEmoji = "🥇"
	} else if quota.Plan == "Max5x" {
		planColor = BrightYellow
		planEmoji = "🥈"
	}
	
	output.WriteString(fmt.Sprintf("   %s 推测计划: %s (%s)\n", 
		planEmoji,
		f.Colors.Colorize(quota.Plan, planColor),
		f.Colors.Dim(fmt.Sprintf("$%s/月", getPlanPrice(quota.Plan)))))
	
	// 窗口信息
	output.WriteString(fmt.Sprintf("   🕐 限额机制: %s窗口 (%s)\n", 
		f.Colors.Info(quota.WindowDuration),
		f.Colors.Dim("每天4个窗口，可能基于UTC时区")))
	
	// 使用情况进度条
	usageColor := BrightGreen
	usageEmoji := "🟢"
	if quota.UsagePercentage > 80 {
		usageColor = BrightRed
		usageEmoji = "🔴"
	} else if quota.UsagePercentage > 60 {
		usageColor = BrightYellow
		usageEmoji = "🟡"
	}
	
	// 创建简单的进度条
	barLength := 20
	filledLength := int(float64(barLength) * quota.UsagePercentage / 100.0)
	progressBar := strings.Repeat("█", filledLength) + strings.Repeat("░", barLength-filledLength)
	
	output.WriteString(fmt.Sprintf("   %s 估算使用: %s / %s 消息\n",
		usageEmoji,
		f.Colors.Colorize(formatNumber(quota.EstimatedUsed), usageColor),
		f.Colors.BrightCyan(formatNumber(quota.MessagesPerWindow))))
	
	output.WriteString(fmt.Sprintf("   📊 估算进度: [%s] %s\n",
		f.Colors.Colorize(progressBar, usageColor),
		f.Colors.Colorize(fmt.Sprintf("%.1f%%", quota.UsagePercentage), usageColor)))
	
	// 剩余消息数
	if quota.EstimatedRemaining > 0 {
		remainingColor := BrightGreen
		if quota.EstimatedRemaining < 10 {
			remainingColor = BrightRed
		} else if quota.EstimatedRemaining < 20 {
			remainingColor = BrightYellow
		}
		output.WriteString(fmt.Sprintf("   ✨ 估算剩余: %s\n", 
			f.Colors.Colorize(formatNumber(quota.EstimatedRemaining), remainingColor)))
	} else {
		output.WriteString(fmt.Sprintf("   ❌ 估算剩余: %s\n", 
			f.Colors.BrightRed("可能已用完")))
	}
	
	// 当前模型
	modelIcon := "🔥"
	modelDesc := "高性能模型"
	if quota.CurrentModel == "Claude 4 Sonnet" {
		modelIcon = "⚡"
		modelDesc = "标准模型"
	}
	output.WriteString(fmt.Sprintf("   %s 推测模型: %s (%s)\n", 
		modelIcon, 
		f.Colors.BrightBlue(quota.CurrentModel),
		f.Colors.Dim(modelDesc)))
	
	// 下次重置时间
	timeUntilReset := time.Until(quota.NextResetTime)
	if timeUntilReset > 0 {
		resetColor := BrightGreen
		if timeUntilReset < time.Hour {
			resetColor = BrightYellow
		}
		
		// 显示当前系统时区
		_, offset := time.Now().Zone()
		timezone := fmt.Sprintf("UTC%+d", offset/3600)
		
		output.WriteString(fmt.Sprintf("   ⏳ 预测重置: %s (%s) [%s]\n", 
			f.Colors.Dim(quota.NextResetTime.Format("15:04")),
			f.Colors.Colorize(formatDuration(timeUntilReset), resetColor),
			f.Colors.Dim(timezone)))
	} else {
		output.WriteString(fmt.Sprintf("   🔄 重置状态: %s\n", 
			f.Colors.BrightGreen("应该已重置")))
	}
	
	// 使用建议
	output.WriteString("\n")
	output.WriteString(fmt.Sprintf("   💡 %s\n", 
		f.Colors.Info("获取准确信息：在Claude Code中运行 /status 命令")))
	output.WriteString(fmt.Sprintf("   🔧 %s\n", 
		f.Colors.Dim("如果重置时间不准确，请反馈给开发者")))
	
	output.WriteString("\n")
}

// getPlanPrice 获取计划价格（字符串）
func getPlanPrice(plan string) string {
	switch plan {
	case "Pro":
		return "20"
	case "Max5x":
		return "100"
	case "Max20x":
		return "200"
	default:
		return "?"
	}
}

// getPlanCostFloat 获取计划价格（浮点数）
func getPlanCostFloat(plan string) float64 {
	switch plan {
	case "Pro":
		return 20.0
	case "Max5x":
		return 100.0
	case "Max20x":
		return 200.0
	default:
		return 0.0
	}
}

// formatDuration 格式化时间间隔
func formatDuration(d time.Duration) string {
	if d < 0 {
		return "已过期"
	}
	
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	
	if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟后", hours, minutes)
	}
	return fmt.Sprintf("%d分钟后", minutes)
}

// writeDetailedStats 写入详细统计
func (f *Formatter) writeDetailedStats(output *strings.Builder, stats *models.UsageStats) {
	// 项目统计
	if len(stats.ProjectStats) > 0 {
		f.writeProjectStats(output, stats)
	}

	// 按日期统计
	if len(stats.DailyStats) > 0 {
		f.writeDailyStats(output, stats)
	}

	// 会话统计
	if len(stats.SessionStats) > 0 && len(stats.SessionStats) <= 20 {
		f.writeSessionStats(output, stats)
	}
}

// writeProjectStats 写入项目统计
func (f *Formatter) writeProjectStats(output *strings.Builder, stats *models.UsageStats) {
	t := table.NewWriter()
	t.SetTitle("📁 项目统计")
	t.AppendHeader(table.Row{"项目", "路径", "Token数", "最后活动"})

	// 按Token数排序
	type projectStat struct {
		name string
		info models.ProjectStats
	}
	
	var projectStats []projectStat
	for name, info := range stats.ProjectStats {
		projectStats = append(projectStats, projectStat{name, info})
	}
	
	sort.Slice(projectStats, func(i, j int) bool {
		return projectStats[i].info.Tokens.GetTotalTokens() > projectStats[j].info.Tokens.GetTotalTokens()
	})

	for _, ps := range projectStats {
		shortPath := ps.info.ProjectPath
		if len(shortPath) > 40 {
			shortPath = "..." + shortPath[len(shortPath)-37:]
		}
		
		t.AppendRow(table.Row{
			ps.name,
			shortPath,
			formatNumber(ps.info.Tokens.GetTotalTokens()),
			ps.info.LastActivity.Format("01-02 15:04"),
		})
	}

	t.SetStyle(table.StyleColoredBright)
	output.WriteString(t.Render() + "\n\n")
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