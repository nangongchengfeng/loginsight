package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"log-analyzer/internal/stats"
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

// RenderStats 渲染完整的统计结果
func RenderStats(result stats.StatsResult) string {
	var output strings.Builder

	// 渲染仪表盘
	renderDashboard(&output, result)
	output.WriteString("\n")

	// 第一排：左侧状态码+UA垂直堆叠，右侧时间分布
	statusCard := renderStatusCardToString(result)
	uaCard := renderUACardToString(result)
	leftColumn := lipgloss.JoinVertical(lipgloss.Top, statusCard, uaCard)
	hourlyCard := renderHourlyCardToString(result)
	output.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, "  ", hourlyCard))
	output.WriteString("\n")

	// 第二排：IP 卡片 + URL 卡片并排
	ipCard := renderIPCardToString(result)
	urlCard := renderURLCardToString(result)
	output.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, ipCard, "  ", urlCard))

	return output.String()
}

// renderIPCardToString 渲染 IP 卡片并返回字符串
func renderIPCardToString(result stats.StatsResult) string {
	var sb strings.Builder
	renderIPCard(&sb, result)
	return sb.String()
}

// renderURLCardToString 渲染 URL 卡片并返回字符串
func renderURLCardToString(result stats.StatsResult) string {
	var sb strings.Builder
	renderURLCard(&sb, result)
	return sb.String()
}

// renderStatusCardToString 渲染状态码卡片并返回字符串
func renderStatusCardToString(result stats.StatsResult) string {
	var sb strings.Builder
	renderStatusCard(&sb, result)
	return sb.String()
}

// renderHourlyCardToString 渲染时间分布卡片并返回字符串
func renderHourlyCardToString(result stats.StatsResult) string {
	var sb strings.Builder
	renderHourlyCard(&sb, result)
	return sb.String()
}

// renderUACardToString 渲染 UA 卡片并返回字符串
func renderUACardToString(result stats.StatsResult) string {
	var sb strings.Builder
	renderUACard(&sb, result)
	return sb.String()
}

// renderDashboard 渲染仪表盘卡片
func renderDashboard(output *strings.Builder, result stats.StatsResult) {
	fourXXRate := float64(result.CategoryCounts["4xx"]) / float64(result.Total) * 100
	fiveXXRate := float64(result.CategoryCounts["5xx"]) / float64(result.Total) * 100

	card1 := RenderMetricCard("总请求数", fmt.Sprintf("%d", result.Total), ColorBlue)
	card2 := RenderMetricCard("唯一 IP", fmt.Sprintf("%d", result.UniqueIPCount), ColorGreen)
	card3 := RenderMetricCard("4xx 率", fmt.Sprintf("%.1f%%", fourXXRate), ColorOrange)
	card4 := RenderMetricCard("5xx 率", fmt.Sprintf("%.1f%%", fiveXXRate), ColorRed)

	output.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, card1, "  ", card2, "  ", card3, "  ", card4))
}

// renderStatusCard 渲染状态码分布卡片
func renderStatusCard(output *strings.Builder, result stats.StatsResult) {
	var content strings.Builder

	maxCatCount := 0
	for _, cat := range []string{"2xx", "3xx", "4xx", "5xx"} {
		if result.CategoryCounts[cat] > maxCatCount {
			maxCatCount = result.CategoryCounts[cat]
		}
	}

	for _, cat := range []string{"2xx", "3xx", "4xx", "5xx"} {
		count := result.CategoryCounts[cat]
		pct := float64(count) / float64(result.Total) * 100
		color := StatusColor(cat)

		catStr := lipgloss.NewStyle().Foreground(color).Bold(true).Width(6).Render(cat)
		countStr := lipgloss.NewStyle().Width(6).Align(lipgloss.Right).Render(fmt.Sprintf("%d", count))
		pctStr := lipgloss.NewStyle().Width(7).Align(lipgloss.Right).Render(fmt.Sprintf("%.1f%%", pct))
		bar := RenderBar(count, maxCatCount, 14, color) // 条形图从18缩到14

		content.WriteString(fmt.Sprintf("%s %s %s  %s\n", catStr, countStr, pctStr, bar))
	}

	content.WriteString("\n")
	content.WriteString(lipgloss.NewStyle().Foreground(ColorGray).Render("Top 10 具体状态码:\n"))
	content.WriteString(strings.Repeat("─", 36) + "\n")

	topStatus := stats.SortStatusMapByValueDesc(result.StatusCounts)
	if len(topStatus) > 10 {
		topStatus = topStatus[:10]
	}

	for _, s := range topStatus {
		pct := float64(s.Count) / float64(result.Total) * 100
		var color lipgloss.Color
		switch {
		case s.Key >= 200 && s.Key < 300:
			color = StatusColor("2xx")
		case s.Key >= 300 && s.Key < 400:
			color = StatusColor("3xx")
		case s.Key >= 400 && s.Key < 500:
			color = StatusColor("4xx")
		case s.Key >= 500 && s.Key < 600:
			color = StatusColor("5xx")
		default:
			color = ColorGray
		}

		statusStr := lipgloss.NewStyle().Foreground(color).Bold(true).Render(fmt.Sprintf("%d", s.Key))
		countStr := lipgloss.NewStyle().Width(6).Align(lipgloss.Right).Render(fmt.Sprintf("%d", s.Count))
		pctStr := lipgloss.NewStyle().Width(7).Align(lipgloss.Right).Render(fmt.Sprintf("%.1f%%", pct))

		content.WriteString(fmt.Sprintf("  %-12s %s %s\n", statusStr, countStr, pctStr))
	}

	card := BaseCardStyle.Copy().Width(44).Render( // 宽度从52缩到44
		lipgloss.JoinVertical(lipgloss.Left,
			HeaderStyle.Render("状态码分布"),
			content.String(),
		),
	)
	output.WriteString(card)
}

// renderHourlyCard 渲染时间分布卡片
func renderHourlyCard(output *strings.Builder, result stats.StatsResult) {
	var content strings.Builder

	maxCount := 0
	for h := 0; h < 24; h++ {
		if result.HourlyCounts[h] > maxCount {
			maxCount = result.HourlyCounts[h]
		}
	}

	for h := 0; h < 24; h++ {
		count := result.HourlyCounts[h]
		pct := float64(count) / float64(result.Total) * 100

		if count == 0 {
			content.WriteString(fmt.Sprintf("%02d :   -  %-16s  %5.1f%%\n", h, "", 0.0)) // 空格从20缩到16
		} else {
			bar := RenderBar(count, maxCount, 16, ColorBlue) // 条形图从20缩到16
			content.WriteString(fmt.Sprintf("%02d : %3d  %-16s  %5.1f%%\n", h, count, bar, pct)) // 空格从20缩到16
		}
	}

	card := BaseCardStyle.Copy().Width(44).Render( // 宽度从50缩到44
		lipgloss.JoinVertical(lipgloss.Left,
			HeaderStyle.Render("时间分布（按小时）"),
			content.String(),
		),
	)
	output.WriteString(card)
}

// renderIPCard 渲染 Top IP 卡片
func renderIPCard(output *strings.Builder, result stats.StatsResult) {
	var content strings.Builder

	topItems := stats.TopN(stats.SortMapByValueDesc(result.IPCounts), 10)

	// 表头用普通字符串，不通过 lipgloss 渲染以避免换行
	content.WriteString(fmt.Sprintf("IP 地址                   请求数\n"))
	content.WriteString(strings.Repeat("─", 28) + "\n")

	for _, kc := range topItems {
		key := kc.Key
		if len(key) > 18 {
			key = key[:15] + "..."
		}
		content.WriteString(fmt.Sprintf("%-18s %8d\n", key, kc.Count))
	}

	card := BaseCardStyle.Copy().Width(40).Render( // 宽度从36加到40
		lipgloss.JoinVertical(lipgloss.Left,
			HeaderStyle.Render("Top 10 IP"),
			content.String(),
		),
	)
	output.WriteString(card)
}

// renderURLCard 渲染 Top URL 卡片
func renderURLCard(output *strings.Builder, result stats.StatsResult) {
	var content strings.Builder

	topItems := stats.TopN(stats.SortMapByValueDesc(result.PathCounts), 10)

	content.WriteString(fmt.Sprintf("URL                         请求数\n")) // 缩短标题
	content.WriteString(strings.Repeat("─", 36) + "\n") // 分隔线从40缩到36

	for _, kc := range topItems {
		key := kc.Key
		if len(key) > 26 { // 从30缩到26
			key = key[:23] + "..."
		}
		content.WriteString(fmt.Sprintf("%-26s %8d\n", key, kc.Count)) // 从30缩到26
	}

	card := BaseCardStyle.Copy().Width(44).Render( // 宽度从48缩到44
		lipgloss.JoinVertical(lipgloss.Left,
			HeaderStyle.Render("Top 10 URL"),
			content.String(),
		),
	)
	output.WriteString(card)
}

// renderUACard 渲染 Top UserAgent 卡片
func renderUACard(output *strings.Builder, result stats.StatsResult) {
	var content strings.Builder

	topItems := stats.TopN(stats.SortMapByValueDesc(result.UACounts), 10)

	content.WriteString(fmt.Sprintf("用户代理                    请求数\n")) // 缩短标题
	content.WriteString(strings.Repeat("─", 36) + "\n") // 分隔线从40缩到36

	for _, kc := range topItems {
		content.WriteString(fmt.Sprintf("%-26s %8d\n", truncate(kc.Key, 26), kc.Count)) // 从30缩到26
	}

	card := BaseCardStyle.Copy().Width(44).Render( // 宽度从48缩到44
		lipgloss.JoinVertical(lipgloss.Left,
			HeaderStyle.Render("Top 10 用户代理"),
			content.String(),
		),
	)
	output.WriteString(card)
}

// truncate 截断字符串到最大长度
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
