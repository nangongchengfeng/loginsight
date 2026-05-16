package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderTotalRequests 渲染总请求数卡片（保持兼容，但新增小卡片函数）
func RenderTotalRequests(total int) string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorBlue).
		Render("总请求数")

	count := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorWhite).
		Background(ColorBlue).
		Padding(0, 3).
		Render(fmt.Sprintf("%d", total))

	card := BaseCardStyle.Copy().
		Width(30).
		Align(lipgloss.Center).
		Render(lipgloss.JoinVertical(lipgloss.Center, title, "", count))

	return card
}

// RenderMetricCard 渲染通用指标小卡片
func RenderMetricCard(label string, value string, valueColor lipgloss.Color) string {
	labelRendered := lipgloss.NewStyle().
		Foreground(ColorGray).
		Render(label)

	valueRendered := lipgloss.NewStyle().
		Bold(true).
		Foreground(valueColor).
		Render(value)

	return SmallCardStyle.Copy().
		Render(lipgloss.JoinVertical(lipgloss.Center, labelRendered, "", valueRendered))
}

// RenderBar 渲染 ASCII 条形图（增强版，带颜色）
func RenderBar(value, max, width int, color lipgloss.Color) string {
	if max == 0 {
		return strings.Repeat(" ", width)
	}
	filled := int(float64(value) / float64(max) * float64(width))
	if filled < 1 && value > 0 {
		filled = 1
	}

	bar := lipgloss.NewStyle().Foreground(color).Render(strings.Repeat("█", filled))
	empty := strings.Repeat(" ", width-filled)
	return bar + empty
}

// RenderSimpleBar 无颜色版（兼容旧代码）
func RenderSimpleBar(value, max, width int) string {
	if max == 0 {
		return strings.Repeat(" ", width)
	}
	filled := int(float64(value) / float64(max) * float64(width))
	if filled < 1 && value > 0 {
		filled = 1
	}
	return strings.Repeat("█", filled) + strings.Repeat(" ", width-filled)
}
