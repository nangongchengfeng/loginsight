Status: ready-for-agent

## Parent

PRD: .scratch/nginx-log-analyzer/PRD.md

## What to build

在 `stats` 命令中添加 Top 10 URL 统计。

## Acceptance criteria

- [ ] 从请求行中提取 URL 路径
- [ ] 按 URL 路径分组计数
- [ ] 按请求数降序排序
- [ ] 输出前 10 个 URL 及其请求数

## Blocked by

- #04: stats: Top 10 IP
