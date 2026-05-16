package analyzer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"log-analyzer/internal/parser"
	"log-analyzer/internal/ui"
)

type KeyCount struct {
	Key   string
	Count int
}

type StatusKeyCount struct {
	Key   int
	Count int
}

type StatsResult struct {
	Total         int
	StatusCounts  map[int]int
	CategoryCounts map[string]int
	IPCounts      map[string]int
	PathCounts    map[string]int
	UACounts      map[string]int
	UniqueIPCount int
	HourlyCounts  map[int]int
}

func Analyze(entries []parser.LogEntry) StatsResult {
	total := len(entries)
	statusCounts := make(map[int]int)
	categoryCounts := make(map[string]int)
	ipCounts := make(map[string]int)
	pathCounts := make(map[string]int)
	uaCounts := make(map[string]int)
	hourlyCounts := make(map[int]int)

	for _, e := range entries {
		statusCounts[e.Status]++
		ipCounts[e.IP]++
		pathCounts[e.Path]++
		uaCounts[e.UserAgent]++
		hourlyCounts[e.Time.Hour()]++
		switch {
		case e.Status >= 200 && e.Status < 300:
			categoryCounts["2xx"]++
		case e.Status >= 300 && e.Status < 400:
			categoryCounts["3xx"]++
		case e.Status >= 400 && e.Status < 500:
			categoryCounts["4xx"]++
		case e.Status >= 500 && e.Status < 600:
			categoryCounts["5xx"]++
		}
	}

	return StatsResult{
		Total:         total,
		StatusCounts:  statusCounts,
		CategoryCounts: categoryCounts,
		IPCounts:      ipCounts,
		PathCounts:    pathCounts,
		UACounts:      uaCounts,
		UniqueIPCount: len(ipCounts),
		HourlyCounts:  hourlyCounts,
	}
}

func PrintStats(result StatsResult) {
	renderDashboard(result)
	fmt.Println()
	renderStatusCard(result)
	fmt.Println()
	renderHourlyCard(result)
	fmt.Println()
	renderIPCard(result)
	fmt.Println()
	renderURLCard(result)
	fmt.Println()
	renderUACard(result)
}

func renderDashboard(result StatsResult) {
	fourXXRate := float64(result.CategoryCounts["4xx"]) / float64(result.Total) * 100
	fiveXXRate := float64(result.CategoryCounts["5xx"]) / float64(result.Total) * 100

	card1 := ui.RenderMetricCard("总请求数", fmt.Sprintf("%d", result.Total), ui.ColorBlue)
	card2 := ui.RenderMetricCard("唯一 IP", fmt.Sprintf("%d", result.UniqueIPCount), ui.ColorGreen)
	card3 := ui.RenderMetricCard("4xx 率", fmt.Sprintf("%.1f%%", fourXXRate), ui.ColorOrange)
	card4 := ui.RenderMetricCard("5xx 率", fmt.Sprintf("%.1f%%", fiveXXRate), ui.ColorRed)

	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, card1, "  ", card2, "  ", card3, "  ", card4))
}

func renderStatusCard(result StatsResult) {
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
		color := ui.StatusColor(cat)

		catStr := lipgloss.NewStyle().Foreground(color).Bold(true).Width(6).Render(cat)
		countStr := lipgloss.NewStyle().Width(6).Align(lipgloss.Right).Render(fmt.Sprintf("%d", count))
		pctStr := lipgloss.NewStyle().Width(7).Align(lipgloss.Right).Render(fmt.Sprintf("%.1f%%", pct))
		bar := ui.RenderBar(count, maxCatCount, 18, color)

		content.WriteString(fmt.Sprintf("%s %s %s  %s\n", catStr, countStr, pctStr, bar))
	}

	content.WriteString("\n")
	content.WriteString(lipgloss.NewStyle().Foreground(ui.ColorGray).Render("Top 10 具体状态码:\n"))
	content.WriteString(strings.Repeat("─", 36) + "\n")

	topStatus := SortStatusMapByValueDesc(result.StatusCounts)
	if len(topStatus) > 10 {
		topStatus = topStatus[:10]
	}

	for _, s := range topStatus {
		pct := float64(s.Count) / float64(result.Total) * 100
		var color lipgloss.Color
		switch {
		case s.Key >= 200 && s.Key < 300:
			color = ui.StatusColor("2xx")
		case s.Key >= 300 && s.Key < 400:
			color = ui.StatusColor("3xx")
		case s.Key >= 400 && s.Key < 500:
			color = ui.StatusColor("4xx")
		case s.Key >= 500 && s.Key < 600:
			color = ui.StatusColor("5xx")
		default:
			color = ui.ColorGray
		}

		statusStr := lipgloss.NewStyle().Foreground(color).Bold(true).Render(fmt.Sprintf("%d", s.Key))
		countStr := lipgloss.NewStyle().Width(6).Align(lipgloss.Right).Render(fmt.Sprintf("%d", s.Count))
		pctStr := lipgloss.NewStyle().Width(7).Align(lipgloss.Right).Render(fmt.Sprintf("%.1f%%", pct))

		content.WriteString(fmt.Sprintf("  %-12s %s %s\n", statusStr, countStr, pctStr))
	}

	card := ui.BaseCardStyle.Copy().Width(52).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			ui.HeaderStyle.Render("状态码分布"),
			content.String(),
		),
	)
	fmt.Print(card)
}

func renderHourlyCard(result StatsResult) {
	var content strings.Builder

	maxCount := 0
	for h := 0; h < 24; h++ {
		if result.HourlyCounts[h] > maxCount {
			maxCount = result.HourlyCounts[h]
		}
	}

	// 6列一行，更整齐
	for row := 0; row < 4; row++ {
		for col := 0; col < 6; col++ {
			h := row*6 + col
			count := result.HourlyCounts[h]
			if count == 0 {
				content.WriteString(fmt.Sprintf("%02d -  ", h))
			} else {
				bar := ui.RenderBar(count, maxCount, 6, ui.ColorBlue)
				content.WriteString(fmt.Sprintf("%02d %2d %s  ", h, count, bar))
			}
		}
		content.WriteString("\n")
	}

	card := ui.BaseCardStyle.Copy().Width(70).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			ui.HeaderStyle.Render("时间分布（按小时）"),
			content.String(),
		),
	)
	fmt.Print(card)
}

func renderIPCard(result StatsResult) {
	var content strings.Builder

	topItems := TopN(SortMapByValueDesc(result.IPCounts), 10)

	header := lipgloss.NewStyle().Bold(true).Foreground(ui.ColorGray).Render(fmt.Sprintf("%-18s %8s", "IP 地址", "请求数"))
	content.WriteString(header + "\n")
	content.WriteString(strings.Repeat("─", 28) + "\n")

	for _, kc := range topItems {
		key := kc.Key
		if len(key) > 18 {
			key = key[:15] + "..."
		}
		content.WriteString(fmt.Sprintf("%-18s %8d\n", key, kc.Count))
	}

	card := ui.BaseCardStyle.Copy().Width(36).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			ui.HeaderStyle.Render("Top 10 IP"),
			content.String(),
		),
	)
	fmt.Print(card)
}

func renderURLCard(result StatsResult) {
	var content strings.Builder

	topItems := TopN(SortMapByValueDesc(result.PathCounts), 10)

	header := lipgloss.NewStyle().Bold(true).Foreground(ui.ColorGray).Render(fmt.Sprintf("%-30s %8s", "URL", "请求数"))
	content.WriteString(header + "\n")
	content.WriteString(strings.Repeat("─", 40) + "\n")

	for _, kc := range topItems {
		key := kc.Key
		if len(key) > 30 {
			key = key[:27] + "..."
		}
		content.WriteString(fmt.Sprintf("%-30s %8d\n", key, kc.Count))
	}

	card := ui.BaseCardStyle.Copy().Width(48).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			ui.HeaderStyle.Render("Top 10 URL"),
			content.String(),
		),
	)
	fmt.Print(card)
}

func renderUACard(result StatsResult) {
	var content strings.Builder

	topItems := TopN(SortMapByValueDesc(result.UACounts), 10)

	header := lipgloss.NewStyle().Bold(true).Foreground(ui.ColorGray).Render(fmt.Sprintf("%-55s %8s", "用户代理", "请求数"))
	content.WriteString(header + "\n")
	content.WriteString(strings.Repeat("─", 65) + "\n")

	for _, kc := range topItems {
		content.WriteString(fmt.Sprintf("%-55s %8d\n", truncate(kc.Key, 55), kc.Count))
	}

	card := ui.BaseCardStyle.Copy().Width(74).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			ui.HeaderStyle.Render("Top 10 用户代理"),
			content.String(),
		),
	)
	fmt.Print(card)
}

func SortMapByValueDesc(m map[string]int) []KeyCount {
	var list []KeyCount
	for k, v := range m {
		list = append(list, KeyCount{k, v})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Count > list[j].Count
	})
	return list
}

func TopN(list []KeyCount, n int) []KeyCount {
	if len(list) < n {
		return list
	}
	return list[:n]
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func SortStatusMapByValueDesc(m map[int]int) []StatusKeyCount {
	var list []StatusKeyCount
	for k, v := range m {
		list = append(list, StatusKeyCount{k, v})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Count > list[j].Count
	})
	return list
}
