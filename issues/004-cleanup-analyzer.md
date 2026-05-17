## What to build

清理 analyzer 包，移除所有渲染相关代码，只保留数据分析核心功能。

1. 从 analyzer 包删除 PrintStats 和所有 renderXxx 函数
2. 移除 analyzer 包对 ui 包的依赖
3. 确保 analyzer 包只包含数据分析逻辑

## Acceptance criteria

- [ ] analyzer 包不再导入 ui 包
- [ ] 所有测试通过
- [ ] go run . stats 仍然正常工作（通过 main.go 调用 ui.RenderStats）

## Blocked by

#003

## Status

completed
