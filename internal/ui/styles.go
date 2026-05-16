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
	ColorDark   = lipgloss.Color("#555")    // 更深的灰色
	ColorWhite  = lipgloss.Color("#FFF")    // 白色
)

// 基础样式
var (
	// 基础卡片样式
	BaseCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			BorderForeground(ColorBlue)

	// 小卡片样式（仪表盘用）
	SmallCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 3).
			BorderForeground(ColorBlue).
			Align(lipgloss.Center)

	// 标题样式
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBlue).
			MarginBottom(1)

	// 表格样式
	TableStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorGray)

	// 表格表头样式
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorGray).
				MarginBottom(1)

	// 表格分隔线样式
	TableDividerStyle = lipgloss.NewStyle().
				Foreground(ColorDark)

	// 大数字样式
	LargeNumberStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorBlue).
				Render

	// 小标签样式
	LabelStyle = lipgloss.NewStyle().
			Foreground(ColorGray).
			Render
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
