## What to build

先解决 import cycle：创建独立的 stats 包放 StatsResult，避免 analyzer 和 ui 相互依赖。

1. 创建 `internal/stats` 包，把 `StatsResult`、`KeyCount`、`StatusKeyCount` 从 analyzer 移过去
2. 更新 analyzer 包使用 stats.StatsResult
3. 更新 main.go 使用新的包结构

## Acceptance criteria

- [ ] `go build` 编译通过
- [ ] `go test ./...` 所有测试通过
- [ ] `go run . stats nginx.log` 仍能正常工作

## Blocked by

None - can start immediately

## Status

completed
