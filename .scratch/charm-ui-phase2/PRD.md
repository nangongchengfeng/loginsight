# PRD: Charm UI 第二阶段优化

## Problem Statement

第一阶段已完成 stats 输出美化和帮助信息改进，现在需要继续优化 filter 输出，并添加进度条提示提升用户体验。

## Solution

- 美化 filter 命令输出，给匹配的部分添加高亮颜色
- 在处理大文件时显示进度条/Spinner

## User Stories

1. 作为运维工程师，filter 过滤结果中状态码/IP/路径等关键信息能高亮显示，便于快速识别
2. 作为工具使用者，处理大日志文件时能看到进度提示，知道程序还在运行
3. 作为工具使用者，filter 输出有更好的视觉层次和颜色区分

## Implementation Decisions

### 模块修改

- **filter 模块**: 美化输出，添加高亮
- **ui 模块**: 添加进度条/Spinner 组件和高亮样式

### 技术决策

- filter 输出给状态码添加颜色（与 stats 一致）
- filter 输出可以选择性高亮匹配的过滤条件部分
- 使用 Charm 的 bubble spinner 或简单的 ASCII 进度提示
- 进度条在读取和解析文件时显示

## Testing Decisions

- 测试只验证外部行为：输出包含预期数据
- 手动验证视觉效果

## Out of Scope

- 自动检测日志格式（第三阶段）
- 用 gonx 替换正则解析（第三阶段）

## Further Notes

渐进式迭代，第二阶段完成后可以考虑第三阶段的日志解析增强。

Status: ready-for-agent
