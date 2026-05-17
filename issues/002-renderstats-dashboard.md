## What to build

在 ui 包创建 RenderStats 函数，先只渲染仪表盘卡片。

1. 在 ui 包创建 `RenderStats(result stats.StatsResult) string` 函数
2. 先只实现仪表盘（四个小卡片）的渲染
3. 更新 main.go 使用 ui.RenderStats
4. 写测试验证 RenderStats 返回的字符串包含预期内容

## Acceptance criteria

- [ ] `go test ./internal/ui` 测试通过
- [ ] `go run . stats nginx-sample.log` 仍能正常显示仪表盘
- [ ] RenderStats 返回的字符串包含 "总请求数"、"唯一 IP"、"4xx 率"、"5xx 率"

## Blocked by

#001

## Status

completed
