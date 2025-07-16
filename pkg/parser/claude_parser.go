package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zhuiye8/claude-stats/pkg/models"
)

// ClaudeParser 用于解析Claude JSONL日志文件
type ClaudeParser struct {
	// 配置选项
	SkipErrors   bool
	Verbose      bool
	DateFilter   *DateFilter
}

// DateFilter 用于过滤日期范围
type DateFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
}

// NewClaudeParser 创建新的解析器实例
func NewClaudeParser() *ClaudeParser {
	return &ClaudeParser{
		SkipErrors: true,
		Verbose:    false,
	}
}

// ParseDirectory 解析目录中的所有JSONL文件
func (p *ClaudeParser) ParseDirectory(dirPath string) (*models.UsageStats, error) {
	stats := &models.UsageStats{
		ModelStats:   make(map[string]models.TokenUsage),
		DailyStats:   make(map[string]models.TokenUsage),
		SessionStats: make(map[string]models.SessionInfo),
		DetectedMode: p.detectMode(dirPath),
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(strings.ToLower(info.Name()), ".jsonl") {
			if p.Verbose {
				fmt.Printf("📂 处理文件: %s\n", path)
			}
			
			fileStats, err := p.ParseFile(path)
			if err != nil {
				if p.SkipErrors {
					fmt.Printf("⚠️  跳过文件 %s: %v\n", path, err)
					return nil
				}
				return fmt.Errorf("解析文件 %s 失败: %w", path, err)
			}

			p.mergeStats(stats, fileStats)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	p.calculatePeriod(stats)
	p.calculateCost(stats)
	
	return stats, nil
}

// ParseFile 解析单个JSONL文件
func (p *ClaudeParser) ParseFile(filePath string) (*models.UsageStats, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	stats := &models.UsageStats{
		ModelStats:   make(map[string]models.TokenUsage),
		DailyStats:   make(map[string]models.TokenUsage),
		SessionStats: make(map[string]models.SessionInfo),
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		entry, err := p.parseLine(line)
		if err != nil {
			if p.SkipErrors {
				if p.Verbose {
					fmt.Printf("⚠️  行 %d 解析错误: %v\n", lineNum, err)
				}
				continue
			}
			return nil, fmt.Errorf("行 %d 解析失败: %w", lineNum, err)
		}

		if entry != nil && p.shouldInclude(entry) {
			p.processEntry(stats, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	return stats, nil
}

// parseLine 解析单行JSONL内容
func (p *ClaudeParser) parseLine(line string) (*models.ConversationEntry, error) {
	var entry models.ConversationEntry
	
	// 先解析到map以处理未知字段
	var rawData map[string]interface{}
	if err := json.Unmarshal([]byte(line), &rawData); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}

	// 解析到结构体
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		return nil, fmt.Errorf("结构体解析失败: %w", err)
	}

	entry.RawData = rawData

	// 手动处理时间戳，支持多种格式
	if timestampStr, ok := rawData["timestamp"].(string); ok {
		timestamp, err := p.parseTimestamp(timestampStr)
		if err == nil {
			entry.Timestamp = timestamp
		}
	}

	return &entry, nil
}

// parseTimestamp 解析时间戳，支持多种格式
func (p *ClaudeParser) parseTimestamp(timestampStr string) (time.Time, error) {
	// 支持的时间格式
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timestampStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析时间戳: %s", timestampStr)
}

// shouldInclude 检查条目是否应该包含在统计中
func (p *ClaudeParser) shouldInclude(entry *models.ConversationEntry) bool {
	if p.DateFilter == nil {
		return true
	}

	if p.DateFilter.StartDate != nil && entry.Timestamp.Before(*p.DateFilter.StartDate) {
		return false
	}

	if p.DateFilter.EndDate != nil && entry.Timestamp.After(*p.DateFilter.EndDate) {
		return false
	}

	return true
}

// processEntry 处理单个条目，更新统计信息
func (p *ClaudeParser) processEntry(stats *models.UsageStats, entry *models.ConversationEntry) {
	stats.TotalMessages++

	// 处理token使用情况
	if entry.Usage != nil && !entry.Usage.IsEmpty() {
		stats.TotalTokens.Add(*entry.Usage)

		// 按模型统计
		if entry.Model != "" {
			modelUsage := stats.ModelStats[entry.Model]
			modelUsage.Add(*entry.Usage)
			stats.ModelStats[entry.Model] = modelUsage
		}

		// 按日期统计
		dateKey := entry.Timestamp.Format("2006-01-02")
		dailyUsage := stats.DailyStats[dateKey]
		dailyUsage.Add(*entry.Usage)
		stats.DailyStats[dateKey] = dailyUsage
	}

	// 处理会话信息
	if entry.SessionID != "" {
		session, exists := stats.SessionStats[entry.SessionID]
		if !exists {
			session = models.SessionInfo{
				ID:        entry.SessionID,
				StartTime: entry.Timestamp,
				EndTime:   entry.Timestamp,
				Model:     entry.Model,
			}
			stats.TotalSessions++
		}

		session.MessageCount++
		if entry.Timestamp.After(session.EndTime) {
			session.EndTime = entry.Timestamp
		}
		if entry.Timestamp.Before(session.StartTime) {
			session.StartTime = entry.Timestamp
		}

		if entry.Usage != nil {
			session.Tokens.Add(*entry.Usage)
		}

		session.Duration = session.EndTime.Sub(session.StartTime).String()
		stats.SessionStats[entry.SessionID] = session
	}
}

// detectMode 检测使用模式（API vs 订阅）
func (p *ClaudeParser) detectMode(dirPath string) string {
	// 简单启发式：检查是否存在cost相关信息
	// 在真实实现中，这里可以有更复杂的逻辑
	return "subscription" // 默认为订阅模式，因为大多数用户使用订阅
}

// mergeStats 合并统计数据
func (p *ClaudeParser) mergeStats(target, source *models.UsageStats) {
	target.TotalMessages += source.TotalMessages
	target.TotalSessions += source.TotalSessions
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

	// 合并会话统计
	for sessionID, session := range source.SessionStats {
		target.SessionStats[sessionID] = session
	}
}

// calculatePeriod 计算分析时间段
func (p *ClaudeParser) calculatePeriod(stats *models.UsageStats) {
	var startTime, endTime time.Time
	first := true

	for _, session := range stats.SessionStats {
		if first {
			startTime = session.StartTime
			endTime = session.EndTime
			first = false
		} else {
			if session.StartTime.Before(startTime) {
				startTime = session.StartTime
			}
			if session.EndTime.After(endTime) {
				endTime = session.EndTime
			}
		}
	}

	if !first {
		stats.AnalysisPeriod = models.Period{
			StartTime: startTime,
			EndTime:   endTime,
			Duration:  endTime.Sub(startTime).String(),
		}
	}
}

// calculateCost 计算成本
func (p *ClaudeParser) calculateCost(stats *models.UsageStats) {
	// 使用Claude 3.5 Sonnet的定价作为默认
	costCalculator := NewCostCalculator()
	stats.EstimatedCost = costCalculator.Calculate(&stats.TotalTokens, stats.ModelStats, stats.DetectedMode == "subscription")
} 