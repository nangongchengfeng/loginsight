package ui

import (
	"strings"
	"testing"

	"log-analyzer/internal/stats"
)

func TestRenderStats_ContainsDashboard(t *testing.T) {
	// 准备测试数据
	result := stats.StatsResult{
		Total:         100,
		UniqueIPCount: 50,
		CategoryCounts: map[string]int{
			"2xx": 60,
			"3xx": 10,
			"4xx": 25,
			"5xx": 5,
		},
	}

	// 调用被测试函数
	output := RenderStats(result)

	// 验证输出包含关键内容
	if !strings.Contains(output, "总请求数") {
		t.Error("RenderStats() should contain '总请求数'")
	}
	if !strings.Contains(output, "100") {
		t.Error("RenderStats() should contain '100'")
	}
	if !strings.Contains(output, "唯一 IP") {
		t.Error("RenderStats() should contain '唯一 IP'")
	}
	if !strings.Contains(output, "50") {
		t.Error("RenderStats() should contain '50'")
	}
	if !strings.Contains(output, "4xx 率") {
		t.Error("RenderStats() should contain '4xx 率'")
	}
	if !strings.Contains(output, "25.0%") {
		t.Error("RenderStats() should contain '25.0%'")
	}
	if !strings.Contains(output, "5xx 率") {
		t.Error("RenderStats() should contain '5xx 率'")
	}
	if !strings.Contains(output, "5.0%") {
		t.Error("RenderStats() should contain '5.0%'")
	}
}
