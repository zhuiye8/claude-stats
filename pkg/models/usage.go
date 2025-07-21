package models

import (
	"time"
)

// ConversationEntry 代表Claude Code JSONL文件中的一条记录
type ConversationEntry struct {
	// Claude Code 实际字段
	Type        string                 `json:"type"`
	Timestamp   time.Time              `json:"timestamp"`
	SessionID   string                 `json:"sessionId,omitempty"`
	UUID        string                 `json:"uuid,omitempty"`
	ParentUUID  string                 `json:"parentUuid,omitempty"`
	UserType    string                 `json:"userType,omitempty"`
	IsSidechain bool                   `json:"isSidechain,omitempty"`
	IsMeta      bool                   `json:"isMeta,omitempty"`
	CWD         string                 `json:"cwd,omitempty"`
	Version     string                 `json:"version,omitempty"`
	RequestID   string                 `json:"requestId,omitempty"`
	
	// 消息内容 - 注意：这里可能是字符串或复杂对象
	Message     interface{}            `json:"message,omitempty"`
	
	// 用于项目分组的路径信息
	Summary     string                 `json:"summary,omitempty"`
	LeafUUID    string                 `json:"leafUuid,omitempty"`
	
	// 解析后的信息
	ParsedMessage *ParsedMessage        `json:"-"` // 不序列化，仅内部使用
	ExtractedUsage *TokenUsage          `json:"-"` // 从消息中提取的token信息
	
	RawData     map[string]interface{} `json:"-"` // 存储原始数据以处理未知字段
}

// ParsedMessage 代表解析后的消息内容
type ParsedMessage struct {
	Role     string      `json:"role,omitempty"`
	Content  interface{} `json:"content,omitempty"` // 可能是字符串或复杂结构
	Model    string      `json:"model,omitempty"`
	Usage    *TokenUsage `json:"usage,omitempty"`
}

// TokenUsage 代表token使用情况
type TokenUsage struct {
	InputTokens             int `json:"input_tokens"`
	OutputTokens            int `json:"output_tokens"`
	CacheCreationTokens     int `json:"cache_creation_input_tokens"`
	CacheReadTokens         int `json:"cache_read_input_tokens"`
	TotalTokens             int `json:"total_tokens,omitempty"`
}

// Content 代表消息内容项（兼容性保留）
type Content struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// SubscriptionQuota 订阅限额信息
type SubscriptionQuota struct {
	Plan              string // 订阅计划: Pro, Max5x, Max20x
	MessagesPerWindow int    // 每个窗口的消息限制
	EstimatedUsed     int    // 估算已使用消息数
	Remaining         int    // 估算剩余消息数
	ModelSwitchPoint  int    // 模型切换点（20%处）
	DebugInfo         map[string]interface{} `json:"debug_info,omitempty"` // 调试信息
}

// UsageStats 代表统计结果
type UsageStats struct {
	TotalSessions       int                    `json:"total_sessions"`
	TotalMessages       int                    `json:"total_messages"`
	TotalTokens         TokenUsage             `json:"total_tokens"`
	ModelStats          map[string]TokenUsage  `json:"model_stats"`
	DailyStats          map[string]TokenUsage  `json:"daily_stats"`
	SessionStats        map[string]SessionInfo `json:"session_stats"`
	EstimatedCost       CostBreakdown          `json:"estimated_cost"`
	AnalysisPeriod      Period                 `json:"analysis_period"`
	DetectedMode        string                 `json:"detected_mode"` // "api" 或 "subscription"
	
	// 新增：Claude Code 特定统计
	ProjectStats        map[string]ProjectStats `json:"project_stats"`
	MessageTypes        map[string]int          `json:"message_types"`
	ParsedMessages      int                     `json:"parsed_messages"`
	ExtractedTokens     int                     `json:"extracted_tokens"`
	
	// 订阅限额信息
	SubscriptionQuota   *SubscriptionQuota     `json:"subscription_quota,omitempty"`
}

// SessionInfo 代表会话信息
type SessionInfo struct {
	ID           string     `json:"id"`
	StartTime    time.Time  `json:"start_time"`
	EndTime      time.Time  `json:"end_time"`
	Duration     string     `json:"duration"`
	MessageCount int        `json:"message_count"`
	Tokens       TokenUsage `json:"tokens"`
	Model        string     `json:"model"`
	ProjectPath  string     `json:"project_path,omitempty"`
}

// CostBreakdown 代表成本分解
type CostBreakdown struct {
	InputCost           float64            `json:"input_cost"`
	OutputCost          float64            `json:"output_cost"`
	CacheCreationCost   float64            `json:"cache_creation_cost"`
	CacheReadCost       float64            `json:"cache_read_cost"`
	TotalCost           float64            `json:"total_cost"`
	Currency            string             `json:"currency"`
	ModelCosts          map[string]float64 `json:"model_costs"`
	IsEstimated         bool               `json:"is_estimated"` // 是否为订阅模式的估算成本
}

// Period 代表分析时间段
type Period struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  string    `json:"duration"`
}

// ProjectStats 代表项目级别的统计
type ProjectStats struct {
	ProjectName  string     `json:"project_name"`
	ProjectPath  string     `json:"project_path"`
	SessionCount int        `json:"session_count"`
	Tokens       TokenUsage `json:"tokens"`
	Cost         float64    `json:"cost"`
	LastActivity time.Time  `json:"last_activity"`
}

// BillingBlock 代表5小时计费窗口
type BillingBlock struct {
	ID            string     `json:"id"`               // 窗口标识符
	StartTime     time.Time  `json:"start_time"`       // 窗口开始时间
	EndTime       time.Time  `json:"end_time"`         // 窗口结束时间 
	ActualEndTime time.Time  `json:"actual_end_time"`  // 实际结束时间
	IsActive      bool       `json:"is_active"`        // 是否为活跃窗口
	TimeRemaining string     `json:"time_remaining"`   // 剩余时间
	Models        []string   `json:"models"`           // 使用的模型
	Tokens        TokenUsage `json:"tokens"`           // Token使用量
	CostUSD       float64    `json:"cost_usd"`         // 成本
	MessageCount  int        `json:"message_count"`    // 消息数量
	BurnRate      int        `json:"burn_rate"`        // 燃烧速率(tokens/min)
	ProjectedTotal int       `json:"projected_total"`  // 预测总量
	ProjectedCost  float64   `json:"projected_cost"`   // 预测成本
}

// BlocksReport 代表blocks报告
type BlocksReport struct {
	Blocks   []BillingBlock `json:"blocks"`
	Summary  TokenUsage     `json:"summary"`
	TotalCost float64       `json:"total_cost"`
}

// DailyReport 日报告结构
type DailyReport struct {
	Type      string            `json:"type"`
	DailyData []DailyDataPoint  `json:"data"`
	Summary   DailyDataPoint    `json:"summary"`
}

// DailyDataPoint 单日数据点
type DailyDataPoint struct {
	Date                    string                    `json:"date"`
	Models                  []string                  `json:"models"`
	InputTokens             int                       `json:"input_tokens"`
	OutputTokens            int                       `json:"output_tokens"`
	CacheCreationTokens     int                       `json:"cache_creation_tokens"`
	CacheReadTokens         int                       `json:"cache_read_tokens"`
	TotalTokens             int                       `json:"total_tokens"`
	CostUSD                 float64                   `json:"cost_usd"`
	MessageCount            int                       `json:"message_count"`
	SessionCount            int                       `json:"session_count"`
	Breakdown               map[string]DailyModelData `json:"breakdown,omitempty"`
}

// DailyModelData 每日模型数据
type DailyModelData struct {
	InputTokens         int     `json:"input_tokens"`
	OutputTokens        int     `json:"output_tokens"`
	CacheCreationTokens int     `json:"cache_creation_tokens"`
	CacheReadTokens     int     `json:"cache_read_tokens"`
	TotalTokens         int     `json:"total_tokens"`
	CostUSD             float64 `json:"cost_usd"`
	MessageCount        int     `json:"message_count"`
}

// GetTotalTokens 计算总token数
func (u *TokenUsage) GetTotalTokens() int {
	if u.TotalTokens > 0 {
		return u.TotalTokens
	}
	return u.InputTokens + u.OutputTokens
}

// Add 累加token使用量
func (u *TokenUsage) Add(other TokenUsage) {
	u.InputTokens += other.InputTokens
	u.OutputTokens += other.OutputTokens
	u.CacheCreationTokens += other.CacheCreationTokens
	u.CacheReadTokens += other.CacheReadTokens
	u.TotalTokens = u.GetTotalTokens()
}

// IsEmpty 检查是否为空的使用统计
func (u *TokenUsage) IsEmpty() bool {
	return u.InputTokens == 0 && u.OutputTokens == 0 && 
		   u.CacheCreationTokens == 0 && u.CacheReadTokens == 0
} 

// EstimateSubscriptionQuota 根据使用情况估算订阅限额
func (u *UsageStats) EstimateSubscriptionQuota() *SubscriptionQuota {
	if u.DetectedMode != "subscription" {
		return nil
	}
	
	// 只计算用户消息数作为真实请求数
	userMessages := 0
	if u.MessageTypes != nil {
		if count, exists := u.MessageTypes["user"]; exists {
			userMessages = count
		}
	}
	
	// 如果没有找到用户消息统计，使用估算
	if userMessages == 0 {
		// 假设用户消息占总消息的30-40%
		userMessages = int(float64(u.TotalMessages) * 0.35)
	}
	
	// 计算数据跨越的时间窗口数量
	duration := u.AnalysisPeriod.EndTime.Sub(u.AnalysisPeriod.StartTime)
	windowCount := int(duration.Hours()/5) + 1 // 每5小时一个窗口
	if windowCount < 1 {
		windowCount = 1
	}
	
	// 基于成本和用户消息数推断计划类型
	plan := "Pro"
	messagesPerWindow := 45
	modelSwitchPoint := 9
	
	// 智能计划检测：同时考虑成本和消息数模式
	avgMessagesPerWindow := float64(userMessages) / float64(windowCount)
	
	// 关键逻辑：基于平均使用量反推真实的计划限制
	// 特别注意：如果数据是工作期间保存的，可能不完整
	
	// 重要修正：如果平均使用量很高，可能说明用户是小计划但经常达到限额
	if avgMessagesPerWindow > 500 {
		// 极高使用量，肯定是Max20x
		plan = "Max20x"
		messagesPerWindow = 900
		modelSwitchPoint = 180
	} else if avgMessagesPerWindow > 250 {
		// 高使用量，可能是Max5x
		plan = "Max5x"
		messagesPerWindow = 225
		modelSwitchPoint = 45
	} else if avgMessagesPerWindow > 60 {
		// 中等使用量，可能是小计划用户经常达到限额
		// 关键假设：如果平均使用量远超Pro限制(45条)，很可能是Pro用户经常限额
		if avgMessagesPerWindow > 300 && u.EstimatedCost.TotalCost > 250 {
			// 极高使用量且高成本，可能是Max5x
			plan = "Max5x"
			messagesPerWindow = 225
			modelSwitchPoint = 45
		} else {
			// 中高使用量，更可能是Pro用户经常达到45条限制
			// 特别是如果数据跨越多个窗口（工作期间保存）
			plan = "Pro"
			messagesPerWindow = 45
			modelSwitchPoint = 9
		}
	} else if avgMessagesPerWindow > 30 {
		// 中低使用量，基于成本判断
		if u.EstimatedCost.TotalCost > 80 {
			plan = "Max5x"
			messagesPerWindow = 225
			modelSwitchPoint = 45
		} else {
			plan = "Pro"
			messagesPerWindow = 45
			modelSwitchPoint = 9
		}
	} else {
		// 低使用量，很可能是Pro用户
		plan = "Pro"
		messagesPerWindow = 45
		modelSwitchPoint = 9
	}
	
	// 智能推断当前窗口使用情况
	// 特别考虑：Pro用户平均234条/窗口说明经常达到45条限制
	
	var currentWindowUsed int
	
	// 情况分析：
	if avgMessagesPerWindow >= float64(messagesPerWindow)*5 {
		// 平均使用量是限制的5倍以上，用户肯定是重度用户，经常限额
		currentWindowUsed = messagesPerWindow // 当前窗口已满
	} else if avgMessagesPerWindow >= float64(messagesPerWindow)*3 {
		// 平均使用量是限制的3倍以上，用户经常达到限额
		currentWindowUsed = messagesPerWindow // 当前窗口已满  
	} else if avgMessagesPerWindow >= float64(messagesPerWindow)*2 {
		// 平均使用量是限制的2倍以上，用户很可能当前已限额
		currentWindowUsed = messagesPerWindow // 当前窗口已满
	} else if avgMessagesPerWindow >= float64(messagesPerWindow) {
		// 平均使用量达到或超过计划限制，用户肯定经常限额
		currentWindowUsed = messagesPerWindow // 当前窗口已满
	} else if avgMessagesPerWindow >= float64(messagesPerWindow)*0.9 {
		// 平均使用量接近限制（90%+），用户很可能当前已限额
		currentWindowUsed = messagesPerWindow // 当前窗口已满
	} else if avgMessagesPerWindow >= float64(messagesPerWindow)*0.7 {
		// 平均使用量较高（70-90%），用户当前可能接近限额
		currentWindowUsed = int(float64(messagesPerWindow) * 0.95) // 95%使用率
	} else if avgMessagesPerWindow >= float64(messagesPerWindow)*0.5 {
		// 中度使用（50-70%），稍微高估当前使用
		currentWindowUsed = int(avgMessagesPerWindow * 1.3)
	} else {
		// 轻度使用，使用平均值
		currentWindowUsed = int(avgMessagesPerWindow)
	}
	
	// 确保不超过限制
	if currentWindowUsed > messagesPerWindow {
		currentWindowUsed = messagesPerWindow
	}
	
	// 计算剩余数量
	remaining := messagesPerWindow - currentWindowUsed
	if remaining < 0 {
		remaining = 0
	}
	
	return &SubscriptionQuota{
		Plan:              plan,
		MessagesPerWindow: messagesPerWindow,
		EstimatedUsed:     currentWindowUsed,
		Remaining:         remaining,
		ModelSwitchPoint:  modelSwitchPoint,
		// 添加调试信息
		DebugInfo: map[string]interface{}{
			"total_user_messages": userMessages,
			"window_count":        windowCount,
			"avg_per_window":      avgMessagesPerWindow,
			"duration_hours":      duration.Hours(),
		},
	}
} 