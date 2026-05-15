Status: ready-for-human

## Parent

PRD: .scratch/nginx-log-analyzer/PRD.md

## What to build

实现 Nginx combined 格式日志解析，添加 `stats` 子命令框架，输出总请求数统计。

## Acceptance criteria

- [x] 解析 Nginx combined 格式日志行成结构化数据
- [x] `stats` 子命令接受日志文件路径作为参数
- [x] 读取日志文件
- [x] 统计并输出总请求数
- [x] 使用 text/tabwriter 格式化输出

## Blocked by

- #01: 初始化项目 + 样本生成器
