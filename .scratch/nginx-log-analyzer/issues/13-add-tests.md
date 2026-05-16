Status: ready-for-human

## Parent

PRD: .scratch/nginx-log-analyzer/PRD.md

## What to build

为 Parser、Analyzer、Filter 模块添加单元测试。

## Acceptance criteria

- [x] parser 包的单元测试，覆盖各种日志行格式
- [x] analyzer 包的单元测试，验证统计结果
- [x] filter 包的单元测试，验证过滤逻辑
- [x] `go test` 全部通过

## Blocked by

- #12: 重构: 拆分 Filter 模块
