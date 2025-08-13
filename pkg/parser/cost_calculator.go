package parser

import (
	"github.com/zhuiye8/claude-stats/pkg/models"
)

// ModelPricing 代表模型定价信息
type ModelPricing struct {
	InputPricePerMToken      float64 // 每百万输入token的价格
	OutputPricePerMToken     float64 // 每百万输出token的价格
	CacheWritePricePerMToken float64 // 每百万缓存写入token的价格 (5分钟)
	CacheReadPricePerMToken  float64 // 每百万缓存读取/刷新token的价格
}

// CostCalculator 用于计算使用成本
type CostCalculator struct {
	// 定价来源: https://docs.anthropic.com/en/docs/about-claude/pricing (2025年8月)
	ModelPrices map[string]ModelPricing
}

// NewCostCalculator 创建新的成本计算器
func NewCostCalculator() *CostCalculator {
	return &CostCalculator{
		ModelPrices: map[string]ModelPricing{
			// Claude 4 / Opus
			"claude-opus-4": {
				InputPricePerMToken:      15.0,
				OutputPricePerMToken:     75.0,
				CacheWritePricePerMToken: 18.75, // 1.25x of input
				CacheReadPricePerMToken:  1.50,  // 0.1x of input
			},
			"claude-opus-3": {
				InputPricePerMToken:      15.0,
				OutputPricePerMToken:     75.0,
				CacheWritePricePerMToken: 18.75,
				CacheReadPricePerMToken:  1.50,
			},
			// Claude 4 / Sonnet
			"claude-sonnet-4": {
				InputPricePerMToken:      3.0,
				OutputPricePerMToken:     15.0,
				CacheWritePricePerMToken: 3.75, // 1.25x of input
				CacheReadPricePerMToken:  0.30, // 0.1x of input
			},
			// Claude 3.5 Sonnet
			"claude-3-5-sonnet": {
				InputPricePerMToken:      3.0,
				OutputPricePerMToken:     15.0,
				CacheWritePricePerMToken: 3.75,
				CacheReadPricePerMToken:  0.30,
			},
			// Claude 3 Sonnet
			"claude-3-sonnet": {
				InputPricePerMToken:      3.0,
				OutputPricePerMToken:     15.0,
				CacheWritePricePerMToken: 3.75,
				CacheReadPricePerMToken:  0.30,
			},
			// Claude 3.5 Haiku
			"claude-3-5-haiku": {
				InputPricePerMToken:      0.80,
				OutputPricePerMToken:     4.0,
				CacheWritePricePerMToken: 1.0, // 1.25x of input
				CacheReadPricePerMToken:  0.08, // 0.1x of input
			},
			// Claude 3 Haiku
			"claude-3-haiku": {
				InputPricePerMToken:      0.25,
				OutputPricePerMToken:     1.25,
				CacheWritePricePerMToken: 0.30,  // Table value, slightly different from 1.25x
				CacheReadPricePerMToken:  0.03, // Table value, slightly different from 0.1x
			},
			// Default fallback (matches Sonnet 3.5)
			"default": {
				InputPricePerMToken:      3.0,
				OutputPricePerMToken:     15.0,
				CacheWritePricePerMToken: 3.75,
				CacheReadPricePerMToken:  0.30,
			},
		},
	}
}

// Calculate 计算总成本
func (c *CostCalculator) Calculate(totalUsage *models.TokenUsage, modelStats map[string]models.TokenUsage, isSubscription bool) models.CostBreakdown {
	breakdown := models.CostBreakdown{
		Currency:    "USD",
		ModelCosts:  make(map[string]float64),
		IsEstimated: isSubscription,
	}

	// 如果有按模型的统计，分别计算
	if len(modelStats) > 0 {
		for model, usage := range modelStats {
			cost := c.calculateModelCost(model, &usage)
			breakdown.ModelCosts[model] = cost
			breakdown.TotalCost += cost
		}
	} else {
		// 如果没有模型信息，使用默认定价（Claude 3.5 Sonnet）
		cost := c.calculateModelCost("claude-3-5-sonnet", totalUsage)
		breakdown.TotalCost = cost
	}

	// 分解成本
	c.breakdownCosts(&breakdown, totalUsage)

	return breakdown
}

// calculateModelCost 计算单个模型的成本
func (c *CostCalculator) calculateModelCost(model string, usage *models.TokenUsage) float64 {
	pricing, exists := c.ModelPrices[model]
	if !exists {
		// 尝试匹配模型族
		switch {
		case strings.Contains(model, "opus"):
			pricing = c.ModelPrices["claude-opus-3"]
		case strings.Contains(model, "sonnet"):
			pricing = c.ModelPrices["claude-3-5-sonnet"]
		case strings.Contains(model, "haiku"):
			pricing = c.ModelPrices["claude-3-5-haiku"]
		default:
			// 使用默认定价
			pricing = c.ModelPrices["default"]
		}
	}

	inputCost := float64(usage.InputTokens) * pricing.InputPricePerMToken / 1_000_000
	outputCost := float64(usage.OutputTokens) * pricing.OutputPricePerMToken / 1_000_000
	cacheCreationCost := float64(usage.CacheCreationTokens) * pricing.CacheWritePricePerMToken / 1_000_000
	cacheReadCost := float64(usage.CacheReadTokens) * pricing.CacheReadPricePerMToken / 1_000_000

	return inputCost + outputCost + cacheCreationCost + cacheReadCost
}

// breakdownCosts 分解总成本
func (c *CostCalculator) breakdownCosts(breakdown *models.CostBreakdown, totalUsage *models.TokenUsage) {
	// 使用默认定价进行分解 (Sonnet 3.5)
	pricing := c.ModelPrices["default"]

	breakdown.InputCost = float64(totalUsage.InputTokens) * pricing.InputPricePerMToken / 1_000_000
	breakdown.OutputCost = float64(totalUsage.OutputTokens) * pricing.OutputPricePerMToken / 1_000_000
	breakdown.CacheCreationCost = float64(totalUsage.CacheCreationTokens) * pricing.CacheWritePricePerMToken / 1_000_000
	breakdown.CacheReadCost = float64(totalUsage.CacheReadTokens) * pricing.CacheReadPricePerMToken / 1_000_000
}

// GetSubscriptionEquivalent 获取订阅模式的等价成本信息
func (c *CostCalculator) GetSubscriptionEquivalent(stats *models.UsageStats) SubscriptionAnalysis {
	totalCost := stats.EstimatedCost.TotalCost
	
	analysis := SubscriptionAnalysis{
		EstimatedAPICost: totalCost,
		Currency:         "USD",
	}

	// 计算建议的订阅计划
	if totalCost <= 20.0 {
		analysis.RecommendedPlan = "Pro ($20)"
		analysis.MonthlySavings = 20.0 - totalCost
	} else if totalCost <= 100.0 {
		analysis.RecommendedPlan = "Max 5× ($100)"
		analysis.MonthlySavings = 100.0 - totalCost
	} else if totalCost <= 200.0 {
		analysis.RecommendedPlan = "Max 20× ($200)"
		analysis.MonthlySavings = 200.0 - totalCost
	} else {
		analysis.RecommendedPlan = "API模式更经济"
		analysis.MonthlySavings = 0
	}

	// 计算5小时窗口的使用情况
	analysis.FiveHourUsage = c.calculateFiveHourWindows(stats)

	return analysis
}

// SubscriptionAnalysis 订阅模式分析
type SubscriptionAnalysis struct {
	EstimatedAPICost  float64                   `json:"estimated_api_cost"`
	RecommendedPlan   string                    `json:"recommended_plan"`
	MonthlySavings    float64                   `json:"monthly_savings"`
	Currency          string                    `json:"currency"`
	FiveHourUsage     []FiveHourWindow          `json:"five_hour_usage"`
}

// FiveHourWindow 5小时窗口使用情况
type FiveHourWindow struct {
	StartTime     string  `json:"start_time"`
	EndTime       string  `json:"end_time"`
	MessageCount  int     `json:"message_count"`
	TokenCount    int     `json:"token_count"`
	EstimatedCost float64 `json:"estimated_cost"`
}

// calculateFiveHourWindows 计算5小时窗口的使用情况
func (c *CostCalculator) calculateFiveHourWindows(stats *models.UsageStats) []FiveHourWindow {
	// 这里可以实现更复杂的5小时窗口分析
	// 简化版本：假设均匀分布
	windows := []FiveHourWindow{}
	
	if len(stats.SessionStats) > 0 {
		// 基于会话数据创建示例窗口
		totalCost := stats.EstimatedCost.TotalCost
		sessionCount := len(stats.SessionStats)
		avgCostPerSession := totalCost / float64(sessionCount)
		
		for sessionID, session := range stats.SessionStats {
			if len(windows) >= 10 { // 限制显示数量
				break
			}
			
			window := FiveHourWindow{
				StartTime:     session.StartTime.Format("2006-01-02 15:04"),
				EndTime:       session.EndTime.Format("2006-01-02 15:04"),
				MessageCount:  session.MessageCount,
				TokenCount:    session.Tokens.GetTotalTokens(),
				EstimatedCost: avgCostPerSession,
			}
			windows = append(windows, window)
			_ = sessionID // 避免未使用变量警告
		}
	}
	
	return windows
} 