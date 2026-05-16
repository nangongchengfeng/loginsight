## What to build

在状态码分类分布下方，新增展示 Top 10 具体状态码。让用户能快速看到最常见的具体状态码（如 404、200、500 等）。

## Acceptance criteria

- [x] 新增辅助函数将 `map[int]int` 转换为可排序结构
- [x] 在状态码分类分布下方展示 Top 10 具体状态码
- [x] 具体状态码展示格式为表格：状态码 | 数量 | 占比
- [x] 新增测试用例验证状态码统计正确
- [x] 运行 `go test ./internal/analyzer` 通过

## Blocked by

None - can start immediately

## Status

✅ 已完成
