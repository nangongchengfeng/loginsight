package analyzer

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/nangongchengfeng/go-cli/internal/parser"
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
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "指标\t值")
	fmt.Fprintf(w, "总请求数\t%d\n", result.Total)
	w.Flush()

	fmt.Println()
	fmt.Println("状态码分布：")
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "分类\t数量\t占比")
	for _, cat := range []string{"2xx", "3xx", "4xx", "5xx"} {
		count := result.CategoryCounts[cat]
		pct := float64(count) / float64(result.Total) * 100
		fmt.Fprintf(w, "%s\t%d\t%.1f%%\n", cat, count, pct)
	}
	w.Flush()

	fmt.Println()
	fmt.Println("Top 10 IP：")
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "IP\t请求数")
	topIPs := TopN(SortMapByValueDesc(result.IPCounts), 10)
	for _, kc := range topIPs {
		fmt.Fprintf(w, "%s\t%d\n", kc.Key, kc.Count)
	}
	w.Flush()

	fmt.Println()
	fmt.Println("Top 10 URL：")
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "URL\t请求数")
	topPaths := TopN(SortMapByValueDesc(result.PathCounts), 10)
	for _, kc := range topPaths {
		fmt.Fprintf(w, "%s\t%d\n", kc.Key, kc.Count)
	}
	w.Flush()

	fmt.Println()
	fmt.Println("Top 10 用户代理：")
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "用户代理\t请求数")
	topUAs := TopN(SortMapByValueDesc(result.UACounts), 10)
	for _, kc := range topUAs {
		fmt.Fprintf(w, "%s\t%d\n", truncate(kc.Key, 60), kc.Count)
	}
	w.Flush()
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
