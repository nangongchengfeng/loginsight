Status: ready-for-human

## Parent

PRD: .scratch/nginx-log-analyzer/PRD.md

## What to build

初始化 Go 项目，设置 Cobra 脚手架，实现 `gen-sample` 子命令生成 Nginx combined 格式的样本日志。

## Acceptance criteria

- [x] `go mod init` 初始化项目
- [x] 安装 Cobra 依赖
- [x] 基础 Cobra CLI 结构
- [x] `gen-sample` 子命令，接受 `--lines` 参数指定行数
- [x] 生成 realistic-looking 的 Nginx combined 格式日志
- [x] 日志输出到 stdout 或指定文件

## Blocked by

None - can start immediately
