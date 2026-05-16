package analyzer

import (
	"testing"
	"time"

	"log-analyzer/internal/parser"
)

func mustParseTime(t string) time.Time {
	ts, err := time.Parse("2006-01-02T15:04:05", t)
	if err != nil {
		panic(err)
	}
	return ts
}

func TestAnalyze(t *testing.T) {
	entries := []parser.LogEntry{
		{IP: "1.1.1.1", Status: 200, Path: "/api/users", UserAgent: "ua1"},
		{IP: "1.1.1.1", Status: 200, Path: "/api/users", UserAgent: "ua1"},
		{IP: "2.2.2.2", Status: 404, Path: "/not-found", UserAgent: "ua2"},
		{IP: "3.3.3.3", Status: 500, Path: "/api/error", UserAgent: "ua1"},
		{IP: "1.1.1.1", Status: 302, Path: "/", UserAgent: "ua2"},
	}

	result := Analyze(entries)

	if result.Total != 5 {
		t.Errorf("Total = %d, want 5", result.Total)
	}

	if result.CategoryCounts["2xx"] != 2 {
		t.Errorf("2xx count = %d, want 2", result.CategoryCounts["2xx"])
	}
	if result.CategoryCounts["4xx"] != 1 {
		t.Errorf("4xx count = %d, want 1", result.CategoryCounts["4xx"])
	}
	if result.CategoryCounts["5xx"] != 1 {
		t.Errorf("5xx count = %d, want 1", result.CategoryCounts["5xx"])
	}
	if result.CategoryCounts["3xx"] != 1 {
		t.Errorf("3xx count = %d, want 1", result.CategoryCounts["3xx"])
	}

	if result.IPCounts["1.1.1.1"] != 3 {
		t.Errorf("IP 1.1.1.1 count = %d, want 3", result.IPCounts["1.1.1.1"])
	}

	if result.PathCounts["/api/users"] != 2 {
		t.Errorf("Path /api/users count = %d, want 2", result.PathCounts["/api/users"])
	}

	// 测试唯一 IP 数
	if result.UniqueIPCount != 3 {
		t.Errorf("UniqueIPCount = %d, want 3", result.UniqueIPCount)
	}
}

func TestSortStatusMapByValueDesc(t *testing.T) {
	m := map[int]int{200: 5, 404: 3, 500: 1, 302: 2}
	result := SortStatusMapByValueDesc(m)

	if len(result) != 4 {
		t.Fatalf("len = %d, want 4", len(result))
	}
	if result[0].Key != 200 || result[0].Count != 5 {
		t.Errorf("result[0] = %v, want {200, 5}", result[0])
	}
	if result[1].Key != 404 || result[1].Count != 3 {
		t.Errorf("result[1] = %v, want {404, 3}", result[1])
	}
}

func TestErrorRates(t *testing.T) {
	entries := []parser.LogEntry{
		{IP: "1.1.1.1", Status: 200},
		{IP: "1.1.1.1", Status: 200},
		{IP: "2.2.2.2", Status: 404},
		{IP: "3.3.3.3", Status: 500},
		{IP: "1.1.1.1", Status: 403},
		{IP: "4.4.4.4", Status: 502},
		{IP: "5.5.5.5", Status: 200},
		{IP: "6.6.6.6", Status: 302},
	}

	result := Analyze(entries)

	// 4xx: 404, 403 -> 2 个
	fourXXRate := float64(result.CategoryCounts["4xx"]) / float64(result.Total) * 100
	if fourXXRate != 25.0 {
		t.Errorf("4xx rate = %.1f, want 25.0", fourXXRate)
	}

	// 5xx: 500, 502 -> 2 个
	fiveXXRate := float64(result.CategoryCounts["5xx"]) / float64(result.Total) * 100
	if fiveXXRate != 25.0 {
		t.Errorf("5xx rate = %.1f, want 25.0", fiveXXRate)
	}

	// 总错误率: 4 个
	totalErrRate := float64(result.CategoryCounts["4xx"]+result.CategoryCounts["5xx"]) / float64(result.Total) * 100
	if totalErrRate != 50.0 {
		t.Errorf("total error rate = %.1f, want 50.0", totalErrRate)
	}
}

func TestHourlyDistribution(t *testing.T) {
	entries := []parser.LogEntry{
		{IP: "1.1.1.1", Status: 200, Time: mustParseTime("2026-05-16T09:30:00")},
		{IP: "1.1.1.1", Status: 200, Time: mustParseTime("2026-05-16T09:45:00")},
		{IP: "2.2.2.2", Status: 404, Time: mustParseTime("2026-05-16T10:15:00")},
		{IP: "3.3.3.3", Status: 500, Time: mustParseTime("2026-05-16T10:30:00")},
		{IP: "1.1.1.1", Status: 302, Time: mustParseTime("2026-05-16T14:00:00")},
		{IP: "4.4.4.4", Status: 200, Time: mustParseTime("2026-05-16T14:30:00")},
		{IP: "5.5.5.5", Status: 200, Time: mustParseTime("2026-05-16T14:45:00")},
		{IP: "6.6.6.6", Status: 200, Time: mustParseTime("2026-05-16T14:50:00")},
	}

	result := Analyze(entries)

	if result.HourlyCounts[9] != 2 {
		t.Errorf("Hour 9 count = %d, want 2", result.HourlyCounts[9])
	}
	if result.HourlyCounts[10] != 2 {
		t.Errorf("Hour 10 count = %d, want 2", result.HourlyCounts[10])
	}
	if result.HourlyCounts[14] != 4 {
		t.Errorf("Hour 14 count = %d, want 4", result.HourlyCounts[14])
	}
}

func TestSortMapByValueDesc(t *testing.T) {
	m := map[string]int{"a": 3, "b": 1, "c": 2}
	result := SortMapByValueDesc(m)

	if len(result) != 3 {
		t.Fatalf("len = %d, want 3", len(result))
	}
	if result[0].Key != "a" || result[0].Count != 3 {
		t.Errorf("result[0] = %v, want {a, 3}", result[0])
	}
	if result[1].Key != "c" || result[1].Count != 2 {
		t.Errorf("result[1] = %v, want {c, 2}", result[1])
	}
	if result[2].Key != "b" || result[2].Count != 1 {
		t.Errorf("result[2] = %v, want {b, 1}", result[2])
	}
}

func TestTopN(t *testing.T) {
	list := []KeyCount{{"a", 5}, {"b", 4}, {"c", 3}, {"d", 2}, {"e", 1}}

	result := TopN(list, 3)
	if len(result) != 3 {
		t.Fatalf("len = %d, want 3", len(result))
	}
	if result[0].Key != "a" || result[2].Key != "c" {
		t.Errorf("TopN() incorrect")
	}

	result = TopN(list, 10)
	if len(result) != 5 {
		t.Errorf("TopN() len = %d, want 5", len(result))
	}
}
