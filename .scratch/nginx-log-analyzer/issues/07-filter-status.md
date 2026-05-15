Status: ready-for-human

## Parent

PRD: .scratch/nginx-log-analyzer/PRD.md

## What to build

添加 `filter` 子命令，支持按状态码过滤日志。

## Acceptance criteria

- [x] `filter` 子命令接受日志文件路径
- [x] `--status` flag 接受状态码（如 404）或状态码范围（如 4xx）
- [x] 输出匹配的日志行

## Blocked by

- #02: 解析日志 + stats 基础 + 总请求数
