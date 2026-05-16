## What to build

新增按小时的时间分布展示。将所有请求按小时（00-23）聚合，展示每个小时的请求数，帮助识别流量高峰低谷期。

## Acceptance criteria

- [x] `StatsResult` 结构体新增字段存储按小时的分布（`map[int]int`）
- [x] `Analyze` 函数按本地时区的小时聚合请求
- [x] 时间分布以条形图形式展示（00-23 每个小时一行）
- [x] 新增测试用例验证时间聚合正确
- [x] 运行 `go test ./internal/analyzer` 通过

## Blocked by

None - can start immediately

## Status

✅ 已完成
