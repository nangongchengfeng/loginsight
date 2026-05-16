## What to build

新增错误率指标展示：4xx 率（客户端错误率）、5xx 率（服务端错误率）、总错误率（4xx+5xx）。以百分比形式展示。

## Acceptance criteria

- [x] `StatsResult` 结构体可选择扩展或在渲染时计算百分比
- [x] 展示 4xx 率、5xx 率、总错误率
- [x] 百分比保留 1 位小数
- [x] 5xx 率用红色高亮显示
- [x] 新增测试用例验证错误率计算正确
- [x] 运行 `go test ./internal/analyzer` 通过

## Blocked by

None - can start immediately

## Status

✅ 已完成
