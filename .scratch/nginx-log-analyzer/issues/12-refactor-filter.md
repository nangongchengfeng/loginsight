Status: ready-for-human

## Parent

PRD: .scratch/nginx-log-analyzer/PRD.md

## What to build

重构：将过滤逻辑拆分到 internal/filter 模块。

## Acceptance criteria

- [x] 创建 internal/filter 包
- [x] 提取过滤逻辑到 filter 包
- [x] 保持功能不变
- [x] 代码可编译运行

## Blocked by

- #11: 重构: 拆分 Analyzer 模块
