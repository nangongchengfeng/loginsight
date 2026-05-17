## What to build

把所有剩余卡片的渲染函数从 analyzer 移到 ui 包，并更新 RenderStats 函数完整渲染所有内容。

1. 把 renderStatusCard、renderHourlyCard、renderIPCard、renderURLCard、renderUACard 移到 ui 包
2. 把辅助函数 SortMapByValueDesc、TopN、truncate、SortStatusMapByValueDesc 也移过去（或者放在 stats 包）
3. 更新 ui.RenderStats 调用所有渲染函数
4. 写测试验证 RenderStats 返回的字符串包含关键内容

## Acceptance criteria

- [ ] `go test ./...` 所有测试通过
- [ ] `go run . stats nginx-sample.log` 输出和之前完全一样（包括所有卡片）

## Blocked by

#002

## Status

completed
