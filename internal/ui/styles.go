package ui

import "github.com/charmbracelet/lipgloss"

// 颜色方案定义
var (
	ColorGreen  = lipgloss.Color("#04B575") // 2xx 成功状态
	ColorYellow = lipgloss.Color("#E5C07B") // 3xx 重定向状态
	ColorOrange = lipgloss.Color("#FF9B5A") // 4xx 客户端错误
	ColorRed    = lipgloss.Color("#FF5555") // 5xx 服务器错误
	ColorBlue   = lipgloss.Color("#5AF")    // 强调色
	ColorGray   = lipgloss.Color("#888")    // 次要文本
)

// 基础样式
var (
	// 基础卡片样式
	BaseCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

	// 标题样式
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBlue).
			Padding(0, 1)

	// 表格样式
	TableStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorGray)
)

// StatusColor 返回状态码分类对应的颜色
func StatusColor(category string) lipgloss.Color {
	switch category {
	case "2xx":
		return ColorGreen
	case "3xx":
		return ColorYellow
	case "4xx":
		return ColorOrange
	case "5xx":
		return ColorRed
	default:
		return ColorGray
	}
}
