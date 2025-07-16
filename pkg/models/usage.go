package models

import (
	"time"
)

// ConversationEntry 代表JSONL文件中的一条记录
type ConversationEntry struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Timestamp   time.Time              `json:"timestamp"`
	SessionID   string                 `json:"session_id,omitempty"`
	Model       string                 `json:"model,omitempty"`
	Usage       *TokenUsage            `json:"usage,omitempty"`
	Message     *Message               `json:"message,omitempty"`
	Content     []Content              `json:"content,omitempty"`
	RawData     map[string]interface{} `json:"-"` // 存储原始数据以处理未知字段
}

// TokenUsage 代表token使用情况
type TokenUsage struct {
	InputTokens             int `json:"input_tokens"`
	OutputTokens            int `json:"output_tokens"`
	CacheCreationTokens     int `json:"cache_creation_input_tokens"`
	CacheReadTokens         int `json:"cache_read_input_tokens"`
	TotalTokens             int `json:"total_tokens,omitempty"`
}

// Message 代表消息内容
type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

// Content 代表消息内容项
type Content struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
	// 可能包含其他类型的内容，如图片等
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