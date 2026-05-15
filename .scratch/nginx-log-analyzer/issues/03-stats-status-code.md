Status: ready-for-human

## Parent

PRD: .scratch/nginx-log-analyzer/PRD.md

## What to build

在 `stats` 命令中添加状态码分布统计（2xx/3xx/4xx/5xx 分类）。

## Acceptance criteria

- [x] 按状态码分组计数
- [x] 按 2xx/3xx/4xx/5xx 分类汇总
- [x] 表格输出各分类的数量和占比

## Blocked by

- #02: 解析日志 + stats 基础 + 总请求数
