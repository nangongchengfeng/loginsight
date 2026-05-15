Status: ready-for-agent

## Parent

PRD: .scratch/nginx-log-analyzer/PRD.md

## What to build

在 `filter` 命令中添加按 IP、路径、HTTP 方法过滤。

## Acceptance criteria

- [ ] `--ip` flag 按客户端 IP 过滤
- [ ] `--path` flag 按请求路径过滤（包含匹配）
- [ ] `--method` flag 按 HTTP 方法过滤（GET/POST/PUT 等）
- [ ] 多个条件是 AND 关系

## Blocked by

- #07: filter: 按状态码过滤
