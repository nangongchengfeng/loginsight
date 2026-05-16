package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderTotalRequests 渲染总请求数卡片
func RenderTotalRequests(total int) string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorBlue).
		Render("总请求数")

	count := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFF")).
		Background(ColorBlue).
		Padding(0, 2).
		Render(fmt.Sprintf("%d", total))

	card := BaseCardStyle.Copy().
		Width(30).
		Align(lipgloss.Center).
		Render(lipgloss.JoinVertical(lipgloss.Center, title, "", count))

	return card
}

// RenderBar 渲染 ASCII 条形图
func RenderBar(value, max, width int) string {
	if max == 0 {
		return strings.Repeat(" ", width)
	}
	filled := int(float64(value) / float64(max) * float64(width))
	if filled < 1 && value > 0 {
		filled = 1 // 至少显示一个方块
	}
	return strings.Repeat("█", filled) + strings.Repeat(" ", width-filled)
}
