package analyzer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"log-analyzer/internal/parser"
	"log-analyzer/internal/ui"
)

// KeyCount 表示键值对统计结果（string 键）
type KeyCount struct {
	Key   string // 键
	Count int    //计数
}

// StatusKeyCount 表示状态码统计结果（int 键）
type StatusKeyCount struct {
	Key   int // 状态码
	Count int // 计数
}

// StatsResult 包含所有统计结果
type StatsResult struct {
	Total         int
	StatusCounts  map[int]int
	CategoryCounts map[string]int
	IPCounts      map[string]int
	PathCounts    map[string]int
	UACounts      map[string]int
	UniqueIPCount int // 唯一 IP 数量
	HourlyCounts  map[int]int // 按小时的分布，key 是 0-23
}

// Analyze 分析日志并返回统计结果
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
		// 按本地时区的小时聚合
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
		UniqueIPCount: len(ipCounts), // 唯一 IP 数就是 map 的长度
		HourlyCounts:  hourlyCounts,
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

	// 错误率指标
	renderErrorRates(result)
	fmt.Println()

	// 按小时的时间分布
	renderHourlyDistribution(result)
	fmt.Println()

	// Top 10 IP
	fmt.Println(ui.HeaderStyle.Render("Top 10 IP"))
	// 显示唯一 IP 总数
	fmt.Printf("唯一 IP 数: %d\n\n", result.UniqueIPCount)
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

// renderErrorRates 渲染错误率指标
func renderErrorRates(result StatsResult) {
	fmt.Println(ui.HeaderStyle.Render("错误率指标"))

	fourXXCount := result.CategoryCounts["4xx"]
	fiveXXCount := result.CategoryCounts["5xx"]
	totalErrCount := fourXXCount + fiveXXCount

	fourXXRate := float64(fourXXCount) / float64(result.Total) * 100
	fiveXXRate := float64(fiveXXCount) / float64(result.Total) * 100
	totalErrRate := float64(totalErrCount) / float64(result.Total) * 100

	// 4xx 错误率
	fmt.Printf("4xx 错误率:  %.1f%% (%d 次)\n", fourXXRate, fourXXCount)
	// 5xx 错误率，用红色高亮
	fiveXXStr := lipgloss.NewStyle().Foreground(ui.ColorRed).Bold(true).Render(fmt.Sprintf("%.1f%%", fiveXXRate))
	fmt.Printf("5xx 错误率:  %s (%d 次)\n", fiveXXStr, fiveXXCount)
	// 总错误率
	fmt.Printf("总错误率:   %.1f%% (%d 次)\n", totalErrRate, totalErrCount)
}

// renderHourlyDistribution 渲染按小时的时间分布
func renderHourlyDistribution(result StatsResult) {
	fmt.Println(ui.HeaderStyle.Render("时间分布（按小时）"))

	// 找到最大值用于条形图比例
	maxCount := 0
	for h := 0; h < 24; h++ {
		if result.HourlyCounts[h] > maxCount {
			maxCount = result.HourlyCounts[h]
		}
	}

	// 渲染每小时的分布
	for h := 0; h < 24; h++ {
		count := result.HourlyCounts[h]
		hourStr := fmt.Sprintf("%02d:00", h)
		if count == 0 {
			fmt.Printf("%s       0\n", hourStr)
		} else {
			pct := float64(count) / float64(result.Total) * 100
			bar := ui.RenderBar(count, maxCount, 25)
			fmt.Printf("%s  %4d  %5.1f%%  %s\n", hourStr, count, pct, bar)
		}
	}
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

	// 渲染 Top 10 具体状态码
	renderTopStatusCodes(result)
}

// renderTopStatusCodes 渲染 Top 10 具体状态码
func renderTopStatusCodes(result StatsResult) {
	fmt.Println()
	fmt.Println(ui.HeaderStyle.Render("Top 10 具体状态码"))

	topStatus := SortStatusMapByValueDesc(result.StatusCounts)
	if len(topStatus) > 10 {
		topStatus = topStatus[:10]
	}

	// 表头
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.ColorGray).
		Render(fmt.Sprintf("%-10s %8s %10s", "状态码", "请求数", "占比"))
	fmt.Println(header)
	fmt.Println(strings.Repeat("─", 32))

	// 数据行
	for _, s := range topStatus {
		pct := float64(s.Count) / float64(result.Total) * 100
		// 根据状态码确定颜色
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

		statusStr := lipgloss.NewStyle().
			Foreground(color).
			Bold(true).
			Render(fmt.Sprintf("%d", s.Key))

		fmt.Printf("%-22s %8d %9.1f%%\n", statusStr, s.Count, pct)
	}
}

// renderTopTable 渲染 Top N 表格
func renderTopTable(col1, col2 string, counts map[string]int, n int) {
	topItems := TopN(SortMapByValueDesc(counts), n)

	// 表头
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.ColorGray).
		Render(fmt.Sprintf("%-25s %8s", col1, col2))
	fmt.Println(header)
	fmt.Println(strings.Repeat("─", 35))

	// 数据行
	for _, kc := range topItems {
		key := kc.Key
		if len(key) > 25 {
			key = key[:22] + "..."
		}
		fmt.Printf("%-25s %8d\n", key, kc.Count)
	}
}

// renderTopTableWithTruncate 渲染带截断的 Top N 表格
func renderTopTableWithTruncate(col1, col2 string, counts map[string]int, n, maxLen int) {
	topItems := TopN(SortMapByValueDesc(counts), n)

	// 表头
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.ColorGray).
		Render(fmt.Sprintf("%-60s %8s", col1, col2))
	fmt.Println(header)
	fmt.Println(strings.Repeat("─", 72))

	// 数据行
	for _, kc := range topItems {
		fmt.Printf("%-60s %8d\n", truncate(kc.Key, maxLen), kc.Count)
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

// SortStatusMapByValueDesc 将状态码 map 按值降序排序
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
