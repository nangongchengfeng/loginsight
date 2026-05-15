Status: ready-for-agent

## Parent

PRD: .scratch/nginx-log-analyzer/PRD.md

## What to build

重构：将统计分析逻辑拆分到 internal/analyzer 模块。

## Acceptance criteria

- [ ] 创建 internal/analyzer 包
- [ ] 提取统计逻辑到 analyzer 包
- [ ] 保持功能不变
- [ ] 代码可编译运行

## Blocked by

- #10: 重构: 拆分 Parser 模块
