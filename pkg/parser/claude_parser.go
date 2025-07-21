package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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
		ProjectStats: make(map[string]models.ProjectStats),
		MessageTypes: make(map[string]int),
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
		ProjectStats: make(map[string]models.ProjectStats),
		MessageTypes: make(map[string]int),
	}

	scanner := bufio.NewScanner(file)
	// 增加扫描器缓冲区大小以处理长行（Claude日志可能包含大量代码）
	maxCapacity := 10 * 1024 * 1024 // 10MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
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
	// 先解析到map以处理未知字段
	var rawData map[string]interface{}
	if err := json.Unmarshal([]byte(line), &rawData); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}

	// 创建entry并设置原始数据
	entry := &models.ConversationEntry{
		RawData: rawData,
	}

	// 手动提取字段以避免结构不匹配问题
	if typeVal, ok := rawData["type"].(string); ok {
		entry.Type = typeVal
	}

	if sessionID, ok := rawData["sessionId"].(string); ok {
		entry.SessionID = sessionID
	}

	if uuid, ok := rawData["uuid"].(string); ok {
		entry.UUID = uuid
	}

	if parentUUID, ok := rawData["parentUuid"].(string); ok {
		entry.ParentUUID = parentUUID
	}

	if userType, ok := rawData["userType"].(string); ok {
		entry.UserType = userType
	}

	if cwd, ok := rawData["cwd"].(string); ok {
		entry.CWD = cwd
	}

	if version, ok := rawData["version"].(string); ok {
		entry.Version = version
	}

	if requestID, ok := rawData["requestId"].(string); ok {
		entry.RequestID = requestID
	}

	if summary, ok := rawData["summary"].(string); ok {
		entry.Summary = summary
	}

	if leafUUID, ok := rawData["leafUuid"].(string); ok {
		entry.LeafUUID = leafUUID
	}

	// 处理message字段
	if messageData, ok := rawData["message"]; ok {
		entry.Message = messageData
		entry.ParsedMessage = p.parseMessage(messageData)
	}

	// 手动处理时间戳，支持多种格式
	if timestampStr, ok := rawData["timestamp"].(string); ok {
		timestamp, err := p.parseTimestamp(timestampStr)
		if err == nil {
			entry.Timestamp = timestamp
		}
	}

	// 尝试从消息中提取token使用信息
	if entry.ParsedMessage != nil {
		entry.ExtractedUsage = p.extractTokenUsage(entry.ParsedMessage)
	}

	return entry, nil
}

// parseMessage 解析消息内容
func (p *ClaudeParser) parseMessage(messageData interface{}) *models.ParsedMessage {
	parsedMsg := &models.ParsedMessage{}

	switch msg := messageData.(type) {
	case string:
		// 如果是字符串，直接设置为content
		parsedMsg.Content = msg
		// 尝试从字符串中提取model和usage信息
		parsedMsg.Model = p.extractModelFromString(msg)
		parsedMsg.Usage = p.extractUsageFromString(msg)
		
	case map[string]interface{}:
		// 如果是对象，尝试解析各个字段
		if role, ok := msg["role"].(string); ok {
			parsedMsg.Role = role
		}
		
		if content, ok := msg["content"]; ok {
			parsedMsg.Content = content
		}
		
		if model, ok := msg["model"].(string); ok {
			parsedMsg.Model = model
		}

		// 尝试解析usage信息
		if usageData, ok := msg["usage"]; ok {
			parsedMsg.Usage = p.parseUsageFromInterface(usageData)
		}
		
		// 如果没有直接的usage字段，尝试从其他字段提取
		if parsedMsg.Usage == nil || parsedMsg.Usage.IsEmpty() {
			// 检查是否有token相关的字段
			if tokenStr := p.extractStringFromMap(msg, []string{"tokens", "token_count", "usage_info"}); tokenStr != "" {
				parsedMsg.Usage = p.extractUsageFromString(tokenStr)
			}
		}
	}

	return parsedMsg
}

// extractModelFromString 从字符串中提取模型信息
func (p *ClaudeParser) extractModelFromString(content string) string {
	// 常见的Claude模型名称模式
	modelPatterns := []string{
		`claude-3-5-sonnet-[0-9]+`,
		`claude-3-5-haiku-[0-9]+`,
		`claude-3-opus-[0-9]+`,
		`claude-3-sonnet-[0-9]+`,
		`claude-3-haiku-[0-9]+`,
		`claude-[0-9]+-[a-z]+-[0-9]+`,
	}

	for _, pattern := range modelPatterns {
		re := regexp.MustCompile(pattern)
		if match := re.FindString(content); match != "" {
			return match
		}
	}

	return ""
}

// extractUsageFromString 从字符串中提取token使用信息
func (p *ClaudeParser) extractUsageFromString(content string) *models.TokenUsage {
	usage := &models.TokenUsage{}
	
	// 查找各种token模式
	patterns := map[string]*int{
		`"input_tokens":\s*(\d+)`:               &usage.InputTokens,
		`"output_tokens":\s*(\d+)`:              &usage.OutputTokens,
		`"cache_creation_input_tokens":\s*(\d+)`: &usage.CacheCreationTokens,
		`"cache_read_input_tokens":\s*(\d+)`:     &usage.CacheReadTokens,
		`input.*?(\d+).*?tokens`:                &usage.InputTokens,
		`output.*?(\d+).*?tokens`:               &usage.OutputTokens,
		`(\d+).*?input.*?tokens`:                &usage.InputTokens,
		`(\d+).*?output.*?tokens`:               &usage.OutputTokens,
	}

	for pattern, field := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(content); len(matches) > 1 {
			if num, err := strconv.Atoi(matches[1]); err == nil && *field == 0 {
				*field = num
			}
		}
	}

	if usage.IsEmpty() {
		return nil
	}

	usage.TotalTokens = usage.GetTotalTokens()
	return usage
}

// parseUsageFromInterface 从interface{}解析TokenUsage
func (p *ClaudeParser) parseUsageFromInterface(usageData interface{}) *models.TokenUsage {
	switch usage := usageData.(type) {
	case map[string]interface{}:
		tokenUsage := &models.TokenUsage{}
		
		if input, ok := usage["input_tokens"].(float64); ok {
			tokenUsage.InputTokens = int(input)
		}
		if output, ok := usage["output_tokens"].(float64); ok {
			tokenUsage.OutputTokens = int(output)
		}
		if cacheCreate, ok := usage["cache_creation_input_tokens"].(float64); ok {
			tokenUsage.CacheCreationTokens = int(cacheCreate)
		}
		if cacheRead, ok := usage["cache_read_input_tokens"].(float64); ok {
			tokenUsage.CacheReadTokens = int(cacheRead)
		}

		tokenUsage.TotalTokens = tokenUsage.GetTotalTokens()
		return tokenUsage
		
	case string:
		return p.extractUsageFromString(usage)
	}

	return nil
}

// extractStringFromMap 从map中提取字符串值
func (p *ClaudeParser) extractStringFromMap(data map[string]interface{}, keys []string) string {
	for _, key := range keys {
		if val, ok := data[key]; ok {
			if str, ok := val.(string); ok {
				return str
			}
		}
	}
	return ""
}

// extractTokenUsage 从ParsedMessage中提取token使用信息
func (p *ClaudeParser) extractTokenUsage(parsedMsg *models.ParsedMessage) *models.TokenUsage {
	// 优先使用直接的Usage字段
	if parsedMsg.Usage != nil && !parsedMsg.Usage.IsEmpty() {
		return parsedMsg.Usage
	}

	// 从Content中提取
	if contentStr, ok := parsedMsg.Content.(string); ok {
		if usage := p.extractUsageFromString(contentStr); usage != nil {
			return usage
		}
	}

	// 估算token使用量（基于文本长度的简单估算）
	if contentStr, ok := parsedMsg.Content.(string); ok && contentStr != "" {
		estimatedTokens := len(strings.Fields(contentStr)) / 3 * 4 // 粗略估算：4 tokens per 3 words
		if estimatedTokens > 0 {
			usage := &models.TokenUsage{}
			
			// 根据角色分配input/output
			if parsedMsg.Role == "user" {
				usage.InputTokens = estimatedTokens
			} else if parsedMsg.Role == "assistant" {
				usage.OutputTokens = estimatedTokens
			} else {
				// 如果角色不明确，分成一半一半
				usage.InputTokens = estimatedTokens / 2
				usage.OutputTokens = estimatedTokens - usage.InputTokens
			}
			
			usage.TotalTokens = usage.GetTotalTokens()
			return usage
		}
	}

	return nil
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
	
	// 统计消息类型
	if entry.Type != "" {
		stats.MessageTypes[entry.Type]++
	}

	// 调试：显示前几条记录的结构
	if stats.TotalMessages <= 3 && p.Verbose {
		fmt.Printf("🔍 调试信息 - 记录 #%d:\n", stats.TotalMessages)
		fmt.Printf("  Type: %s\n", entry.Type)
		
		model := ""
		if entry.ParsedMessage != nil {
			model = entry.ParsedMessage.Model
		}
		fmt.Printf("  Model: %s\n", model)
		
		fmt.Printf("  Usage: %+v\n", entry.ExtractedUsage)
		fmt.Printf("  RawData字段: %v\n", getMapKeys(entry.RawData))
		
		if entry.ParsedMessage != nil {
			fmt.Printf("  ParsedMessage.Role: %s\n", entry.ParsedMessage.Role)
			if contentStr, ok := entry.ParsedMessage.Content.(string); ok && len(contentStr) > 50 {
				fmt.Printf("  Content预览: %s...\n", contentStr[:50])
			}
		}
		fmt.Println("  ---")
	}

	// 统计解析成功的消息
	if entry.ParsedMessage != nil {
		stats.ParsedMessages++
	}

	// 处理token使用情况
	if entry.ExtractedUsage != nil && !entry.ExtractedUsage.IsEmpty() {
		stats.ExtractedTokens++
		stats.TotalTokens.Add(*entry.ExtractedUsage)

		// 按模型统计
		model := ""
		if entry.ParsedMessage != nil && entry.ParsedMessage.Model != "" {
			model = entry.ParsedMessage.Model
		} else {
			model = "unknown" // 默认模型
		}
		
		modelUsage := stats.ModelStats[model]
		modelUsage.Add(*entry.ExtractedUsage)
		stats.ModelStats[model] = modelUsage

		// 按日期统计
		dateKey := entry.Timestamp.Format("2006-01-02")
		dailyUsage := stats.DailyStats[dateKey]
		dailyUsage.Add(*entry.ExtractedUsage)
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
				ProjectPath: entry.CWD,
			}
			if entry.ParsedMessage != nil && entry.ParsedMessage.Model != "" {
				session.Model = entry.ParsedMessage.Model
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

		if entry.ExtractedUsage != nil {
			session.Tokens.Add(*entry.ExtractedUsage)
		}

		session.Duration = session.EndTime.Sub(session.StartTime).String()
		stats.SessionStats[entry.SessionID] = session
	}

	// 处理项目统计
	if entry.CWD != "" {
		projectKey := filepath.Base(entry.CWD)
		project, exists := stats.ProjectStats[projectKey]
		if !exists {
			project = models.ProjectStats{
				ProjectName: projectKey,
				ProjectPath: entry.CWD,
				LastActivity: entry.Timestamp,
			}
		}

		if entry.Timestamp.After(project.LastActivity) {
			project.LastActivity = entry.Timestamp
		}

		if entry.ExtractedUsage != nil {
			project.Tokens.Add(*entry.ExtractedUsage)
		}

		stats.ProjectStats[projectKey] = project
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

	// 合并会话统计
	for sessionID, session := range source.SessionStats {
		target.SessionStats[sessionID] = session
	}

	// 合并项目统计
	for projectKey, project := range source.ProjectStats {
		target.ProjectStats[projectKey] = project
	}

	// 合并消息类型统计
	for msgType, count := range source.MessageTypes {
		target.MessageTypes[msgType] += count
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

// FinalizeStats 完成统计数据的最终处理
func (p *ClaudeParser) FinalizeStats(stats *models.UsageStats) {
	p.calculatePeriod(stats)
	p.calculateCost(stats)
}

// AnalyzeBlocks 分析5小时计费窗口
func (p *ClaudeParser) AnalyzeBlocks(stats *models.UsageStats) (*models.BlocksReport, error) {
	if len(stats.SessionStats) == 0 {
		return &models.BlocksReport{
			Blocks:    []models.BillingBlock{},
			Summary:   stats.TotalTokens,
			TotalCost: stats.EstimatedCost.TotalCost,
		}, nil
	}

	// 收集所有活动时间点
	var timePoints []time.Time
	for _, session := range stats.SessionStats {
		timePoints = append(timePoints, session.StartTime)
		if !session.EndTime.IsZero() {
			timePoints = append(timePoints, session.EndTime)
		}
	}

	if len(timePoints) == 0 {
		return &models.BlocksReport{
			Blocks:    []models.BillingBlock{},
			Summary:   stats.TotalTokens,
			TotalCost: stats.EstimatedCost.TotalCost,
		}, nil
	}

	// 找到最早和最晚时间
	earliestTime := timePoints[0]
	latestTime := timePoints[0]
	for _, t := range timePoints {
		if t.Before(earliestTime) {
			earliestTime = t
		}
		if t.After(latestTime) {
			latestTime = t
		}
	}

	// 生成5小时窗口
	blocks := p.generateBillingBlocks(earliestTime, latestTime, stats)
	
	return &models.BlocksReport{
		Blocks:    blocks,
		Summary:   stats.TotalTokens,
		TotalCost: stats.EstimatedCost.TotalCost,
	}, nil
}

// generateBillingBlocks 生成5小时计费窗口
func (p *ClaudeParser) generateBillingBlocks(startTime, endTime time.Time, stats *models.UsageStats) []models.BillingBlock {
	var blocks []models.BillingBlock
	
	// 将开始时间向下取整到最近的5小时边界
	startHour := (startTime.Hour() / 5) * 5
	blockStart := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), startHour, 0, 0, 0, startTime.Location())
	
	// 如果开始时间在5小时边界之前，从前一个窗口开始
	if blockStart.After(startTime) {
		blockStart = blockStart.Add(-5 * time.Hour)
	}
	
	currentTime := time.Now()
	
	for blockStart.Before(endTime.Add(5 * time.Hour)) {
		blockEnd := blockStart.Add(5 * time.Hour)
		
		// 分析这个窗口内的活动
		block := p.analyzeBlock(blockStart, blockEnd, stats, currentTime)
		
		// 只包含有活动的窗口
		if block.MessageCount > 0 || block.IsActive {
			blocks = append(blocks, block)
		}
		
		blockStart = blockEnd
	}
	
	return blocks
}

// analyzeBlock 分析单个5小时窗口
func (p *ClaudeParser) analyzeBlock(blockStart, blockEnd time.Time, stats *models.UsageStats, currentTime time.Time) models.BillingBlock {
	block := models.BillingBlock{
		ID:        blockStart.Format("2006-01-02T15:04:05.000Z"),
		StartTime: blockStart,
		EndTime:   blockEnd,
		IsActive:  currentTime.After(blockStart) && currentTime.Before(blockEnd),
		Models:    []string{},
	}
	
	// 计算实际结束时间和剩余时间
	if block.IsActive {
		block.ActualEndTime = currentTime
		remaining := blockEnd.Sub(currentTime)
		if remaining > 0 {
			hours := int(remaining.Hours())
			minutes := int(remaining.Minutes()) % 60
			if hours > 0 {
				block.TimeRemaining = fmt.Sprintf("%dh %dm", hours, minutes)
			} else {
				block.TimeRemaining = fmt.Sprintf("%dm", minutes)
			}
		}
	} else {
		block.ActualEndTime = blockEnd
	}
	
	// 统计窗口内的会话
	modelSet := make(map[string]bool)
	var windowDuration time.Duration
	
	for _, session := range stats.SessionStats {
		// 检查会话是否在此窗口内
		if p.sessionInWindow(session, blockStart, blockEnd) {
			block.MessageCount += session.MessageCount
			block.Tokens.Add(session.Tokens)
			
			if session.Model != "" {
				modelSet[session.Model] = true
			}
			
			// 计算窗口内的实际活动时长
			sessionStart := session.StartTime
			sessionEnd := session.EndTime
			
			// 将会话时间限制在窗口内
			if sessionStart.Before(blockStart) {
				sessionStart = blockStart
			}
			if sessionEnd.After(blockEnd) {
				sessionEnd = blockEnd
			}
			if sessionEnd.IsZero() || sessionEnd.Before(sessionStart) {
				sessionEnd = sessionStart.Add(time.Minute) // 假设至少1分钟
			}
			
			windowDuration += sessionEnd.Sub(sessionStart)
		}
	}
	
	// 提取模型列表
	for model := range modelSet {
		block.Models = append(block.Models, model)
	}
	
	// 计算成本（从总成本按比例分配）
	if stats.TotalTokens.GetTotalTokens() > 0 {
		tokenRatio := float64(block.Tokens.GetTotalTokens()) / float64(stats.TotalTokens.GetTotalTokens())
		block.CostUSD = stats.EstimatedCost.TotalCost * tokenRatio
	}
	
	// 计算燃烧速率和预测
	if block.IsActive && windowDuration.Minutes() > 0 {
		tokensPerMinute := float64(block.Tokens.GetTotalTokens()) / windowDuration.Minutes()
		block.BurnRate = int(tokensPerMinute)
		
		// 预测到窗口结束的总量
		remainingMinutes := blockEnd.Sub(currentTime).Minutes()
		if remainingMinutes > 0 {
			projectedAdditional := tokensPerMinute * remainingMinutes
			block.ProjectedTotal = block.Tokens.GetTotalTokens() + int(projectedAdditional)
			
			// 预测成本
			if stats.TotalTokens.GetTotalTokens() > 0 {
				projectedRatio := float64(block.ProjectedTotal) / float64(stats.TotalTokens.GetTotalTokens())
				block.ProjectedCost = stats.EstimatedCost.TotalCost * projectedRatio
			}
		}
	}
	
	return block
}

// sessionInWindow 检查会话是否在指定窗口内
func (p *ClaudeParser) sessionInWindow(session models.SessionInfo, windowStart, windowEnd time.Time) bool {
	sessionStart := session.StartTime
	sessionEnd := session.EndTime
	
	// 如果会话没有结束时间，使用开始时间
	if sessionEnd.IsZero() {
		sessionEnd = sessionStart
	}
	
	// 检查会话时间范围是否与窗口有重叠
	return !(sessionEnd.Before(windowStart) || sessionStart.After(windowEnd))
}

// getMapKeys 获取map的所有键（调试用）
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
} 