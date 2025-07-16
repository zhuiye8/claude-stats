package formatter

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// 颜色代码
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Dim    = "\033[2m"
	
	// 前景色
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	
	// 亮色
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"
	
	// 背景色
	BgRed    = "\033[41m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgBlue   = "\033[44m"
)

// ColorSettings 颜色设置
type ColorSettings struct {
	Enabled bool
}

// NewColorSettings 创建颜色设置
func NewColorSettings() *ColorSettings {
	return &ColorSettings{
		Enabled: supportsColor(),
	}
}

// supportsColor 检查终端是否支持颜色
func supportsColor() bool {
	// Windows CMD/PowerShell 现在也支持ANSI颜色
	if runtime.GOOS == "windows" {
		// 检查是否在新版本Windows Terminal中
		if os.Getenv("WT_SESSION") != "" {
			return true
		}
		// 检查PowerShell版本
		if os.Getenv("PSModulePath") != "" {
			return true
		}
	}
	
	// Unix系统通常支持颜色
	term := os.Getenv("TERM")
	return term != "" && term != "dumb"
}

// Colorize 给文本添加颜色
func (c *ColorSettings) Colorize(text, color string) string {
	if !c.Enabled {
		return text
	}
	return color + text + Reset
}

// 快捷方法
func (c *ColorSettings) Red(text string) string     { return c.Colorize(text, Red) }
func (c *ColorSettings) Green(text string) string   { return c.Colorize(text, Green) }
func (c *ColorSettings) Yellow(text string) string  { return c.Colorize(text, Yellow) }
func (c *ColorSettings) Blue(text string) string    { return c.Colorize(text, Blue) }
func (c *ColorSettings) Magenta(text string) string { return c.Colorize(text, Magenta) }
func (c *ColorSettings) Cyan(text string) string    { return c.Colorize(text, Cyan) }
func (c *ColorSettings) Bold(text string) string    { return c.Colorize(text, Bold) }
func (c *ColorSettings) Dim(text string) string     { return c.Colorize(text, Dim) }

func (c *ColorSettings) BrightRed(text string) string     { return c.Colorize(text, BrightRed) }
func (c *ColorSettings) BrightGreen(text string) string   { return c.Colorize(text, BrightGreen) }
func (c *ColorSettings) BrightYellow(text string) string  { return c.Colorize(text, BrightYellow) }
func (c *ColorSettings) BrightBlue(text string) string    { return c.Colorize(text, BrightBlue) }
func (c *ColorSettings) BrightMagenta(text string) string { return c.Colorize(text, BrightMagenta) }
func (c *ColorSettings) BrightCyan(text string) string    { return c.Colorize(text, BrightCyan) }

// Success, Warning, Error 语义化颜色
func (c *ColorSettings) Success(text string) string { return c.BrightGreen(text) }
func (c *ColorSettings) Warning(text string) string { return c.BrightYellow(text) }
func (c *ColorSettings) Error(text string) string   { return c.BrightRed(text) }
func (c *ColorSettings) Info(text string) string    { return c.BrightCyan(text) }
func (c *ColorSettings) Header(text string) string  { return c.Bold(c.BrightBlue(text)) }

// 创建彩色分隔符
func (c *ColorSettings) Separator(char string, length int, color string) string {
	line := strings.Repeat(char, length)
	return c.Colorize(line, color)
}

// 创建渐变效果的分隔符
func (c *ColorSettings) GradientSeparator(length int) string {
	if !c.Enabled {
		return strings.Repeat("=", length)
	}
	
	colors := []string{BrightBlue, BrightCyan, BrightGreen, BrightYellow, BrightMagenta}
	var result strings.Builder
	
	for i := 0; i < length; i++ {
		colorIndex := i % len(colors)
		result.WriteString(colors[colorIndex] + "█" + Reset)
	}
	
	return result.String()
}

// 创建带图标的标题
func (c *ColorSettings) IconHeader(icon, title, color string) string {
	if !c.Enabled {
		return fmt.Sprintf("%s %s", icon, title)
	}
	return c.Colorize(fmt.Sprintf("%s %s", icon, title), color)
}

// 创建进度条效果
func (c *ColorSettings) ProgressBar(current, total int, width int) string {
	if total == 0 {
		return ""
	}
	
	percentage := float64(current) / float64(total)
	filled := int(percentage * float64(width))
	empty := width - filled
	
	if !c.Enabled {
		return fmt.Sprintf("[%s%s] %.1f%%", 
			strings.Repeat("█", filled),
			strings.Repeat("░", empty),
			percentage*100)
	}
	
	return fmt.Sprintf("[%s%s%s%s] %s%.1f%%%s",
		BrightGreen, strings.Repeat("█", filled),
		Dim, strings.Repeat("░", empty),
		BrightYellow, percentage*100, Reset)
} 