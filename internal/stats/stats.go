package stats

import "sort"

// KeyCount 表示字符串键值对的统计结果
type KeyCount struct {
	Key   string
	Count int
}

// StatusKeyCount 表示状态码键值对的统计结果
type StatusKeyCount struct {
	Key   int
	Count int
}

// StatsResult 包含所有日志统计结果
type StatsResult struct {
	Total          int            // 总请求数
	StatusCounts   map[int]int    // 具体状态码统计
	CategoryCounts map[string]int // 状态码分类统计（2xx/3xx/4xx/5xx）
	IPCounts       map[string]int // IP 地址统计
	PathCounts     map[string]int // URL 路径统计
	UACounts       map[string]int // UserAgent 统计
	UniqueIPCount  int            // 唯一 IP 数
	HourlyCounts   map[int]int    // 按小时的时间分布统计
}

// SortMapByValueDesc 将 map[string]int 按值降序排序
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

// TopN 返回列表前 N 个元素
func TopN(list []KeyCount, n int) []KeyCount {
	if len(list) < n {
		return list
	}
	return list[:n]
}

// SortStatusMapByValueDesc 将 map[int]int 按值降序排序
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
