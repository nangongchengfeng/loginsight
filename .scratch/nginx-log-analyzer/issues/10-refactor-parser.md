Status: ready-for-human

## Parent

PRD: .scratch/nginx-log-analyzer/PRD.md

## What to build

重构：将日志解析逻辑从 main.go 拆分到 internal/parser 模块。

## Acceptance criteria

- [x] 创建 internal/parser 包
- [x] 提取解析函数到 parser 包
- [x] 保持功能不变
- [x] 代码可编译运行

## Blocked by

- #09: filter: 按时间范围过滤
