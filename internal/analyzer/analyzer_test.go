package analyzer

import (
	"testing"

	"log-analyzer/internal/parser"
)

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
