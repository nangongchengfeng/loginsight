## What to build

在 stats 命令中新增唯一 IP 总数展示。在 Top 10 IP 表格上方显示总共有多少个不同的 IP 地址访问过。

## Acceptance criteria

- [x] `StatsResult` 结构体新增字段存储唯一 IP 数量
- [x] `Analyze` 函数计算唯一 IP 数（`len(IPCounts)`）
- [x] 在 Top 10 IP 表格上方展示唯一 IP 总数
- [x] 新增测试用例验证唯一 IP 数计算正确
- [x] 运行 `go test ./internal/analyzer` 通过

## Blocked by

None - can start immediately

## Status

✅ 已完成
