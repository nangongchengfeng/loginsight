package analyzer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"log-analyzer/internal/parser"
	"log-analyzer/internal/ui"
)

// KeyCount 表示键值对统计结果
type KeyCount struct {
	Key   string // 键
	Count int    //计数
}

// StatsResult 包含所有统计结果
type StatsResult struct {
	Total         int
	StatusCounts  map[int]int
	CategoryCounts map[string]int
	IPCounts      map[string]int
	PathCounts    map[string]int
	UACounts      map[string]int
}

// Analyze 分析日志并返回统计结果
func Analyze(entries []parser.LogEntry) StatsResult {
	total := len(entries)
	statusCounts := make(map[int]int)
	categoryCounts := make(map[string]int)
	ipCounts := make(map[string]int)
	pathCounts := make(map[string]int)
	uaCounts := make(map[string]int)

	for _, e := range entries {
		statusCounts[e.Status]++
		ipCounts[e.IP]++
		pathCounts[e.Path]++
		uaCounts[e.UserAgent]++
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
	}
}

// PrintStats 打印统计结果
func PrintStats(result StatsResult) {
	// 总请求数卡片
	fmt.Println(ui.RenderTotalRequests(result.Total))
	fmt.Println()

	// 状态码分布表格
	fmt.Println(ui.HeaderStyle.Render("状态码分布"))
	renderStatusDistribution(result)
	fmt.Println()

	// Top 10 IP
	fmt.Println(ui.HeaderStyle.Render("Top 10 IP"))
	renderTopTable("IP", "请求数", result.IPCounts, 10)
	fmt.Println()

	// Top 10 URL
	fmt.Println(ui.HeaderStyle.Render("Top 10 URL"))
	renderTopTable("URL", "请求数", result.PathCounts, 10)
	fmt.Println()

	// Top 10 User-Agent
	fmt.Println(ui.HeaderStyle.Render("Top 10 用户代理"))
	renderTopTableWithTruncate("用户代理", "请求数", result.UACounts, 10, 60)
}

// renderStatusDistribution 渲染状态码分布
func renderStatusDistribution(result StatsResult) {
	// 找到最大值用于条形图比例
	maxCount := 0
	for _, cat := range []string{"2xx", "3xx", "4xx", "5xx"} {
		if result.CategoryCounts[cat] > maxCount {
			maxCount = result.CategoryCounts[cat]
		}
	}

	// 渲染表格行
	var rows []string
	for _, cat := range []string{"2xx", "3xx", "4xx", "5xx"} {
		count := result.CategoryCounts[cat]
		pct := float64(count) / float64(result.Total) * 100
		color := ui.StatusColor(cat)

		cellCat := lipgloss.NewStyle().
			Foreground(color).
			Bold(true).
			Width(5).
			Render(cat)

		cellCount := lipgloss.NewStyle().
			Width(8).
			Align(lipgloss.Right).
			Render(fmt.Sprintf("%d", count))

		cellPct := lipgloss.NewStyle().
			Width(8).
			Align(lipgloss.Right).
			Render(fmt.Sprintf("%.1f%%", pct))

		cellBar := lipgloss.NewStyle().
			Foreground(color).
			Render(ui.RenderBar(count, maxCount, 20))

		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Left, cellCat, cellCount, cellPct, "  ", cellBar))
	}

	fmt.Println(strings.Join(rows, "\n"))
}

// renderTopTable 渲染 Top N 表格
func renderTopTable(col1, col2 string, counts map[string]int, n int) {
	topItems := TopN(SortMapByValueDesc(counts), n)

	// 表头
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.ColorGray).
		Render(fmt.Sprintf("%-20s %s", col1, col2))
	fmt.Println(header)
	fmt.Println(strings.Repeat("─", 30))

	// 数据行
	for _, kc := range topItems {
		key := kc.Key
		if len(key) > 20 {
			key = key[:17] + "..."
		}
		fmt.Printf("%-20s %d\n", key, kc.Count)
	}
}

// renderTopTableWithTruncate 渲染带截断的 Top N 表格
func renderTopTableWithTruncate(col1, col2 string, counts map[string]int, n, maxLen int) {
	topItems := TopN(SortMapByValueDesc(counts), n)

	// 表头
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.ColorGray).
		Render(fmt.Sprintf("%-60s %s", col1, col2))
	fmt.Println(header)
	fmt.Println(strings.Repeat("─", 70))

	// 数据行
	for _, kc := range topItems {
		fmt.Printf("%-60s %d\n", truncate(kc.Key, maxLen), kc.Count)
	}
}

// SortMapByValueDesc 将 map 按值降序排序
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

// TopN 返回列表的前 N 个元素
func TopN(list []KeyCount, n int) []KeyCount {
	if len(list) < n {
		return list
	}
	return list[:n]
}

// truncate 截断超长字符串
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
