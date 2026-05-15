Status: ready-for-human

## Parent

PRD: .scratch/nginx-log-analyzer/PRD.md

## What to build

在 `stats` 命令中添加 Top 10 User-Agent 统计。

## Acceptance criteria

- [x] 按 User-Agent 分组计数
- [x] 按请求数降序排序
- [x] 输出前 10 个 User-Agent（太长时截断）及其请求数

## Blocked by

- #05: stats: Top 10 URL
