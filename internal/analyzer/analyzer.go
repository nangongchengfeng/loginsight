package analyzer

import (
	"log-analyzer/internal/parser"
	"log-analyzer/internal/stats"
)

// Analyze 分析日志条目，返回统计结果
func Analyze(entries []parser.LogEntry) stats.StatsResult {
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

	return stats.StatsResult{
		Total:          total,
		StatusCounts:   statusCounts,
		CategoryCounts: categoryCounts,
		IPCounts:       ipCounts,
		PathCounts:     pathCounts,
		UACounts:       uaCounts,
		UniqueIPCount:  len(ipCounts),
		HourlyCounts:   hourlyCounts,
	}
}
