Status: ready-for-agent

## Parent

PRD: .scratch/nginx-log-analyzer/PRD.md

## What to build

在 `filter` 命令中添加按时间范围过滤。

## Acceptance criteria

- [ ] `--since` flag 指定开始时间（格式 2006-01-02T15:04:05Z07:00）
- [ ] `--until` flag 指定结束时间
- [ ] 解析日志中的时间并比较
- [ ] 输出时间范围内的日志

## Blocked by

- #08: filter: 按 IP/路径/方法过滤
