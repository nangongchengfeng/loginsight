# 03: 美化状态码分布（颜色 + ASCII 条形图）

## Parent

.scratch/charm-ui-phase1/PRD.md

## What to build

状态码分类表格添加颜色（2xx绿、3xx黄、4xx橙、5xx红），同时用 ASCII 条形图可视化各分类的比例。

## Acceptance criteria

- [ ] 2xx 显示绿色，3xx 黄色，4xx 橙色，5xx 红色
- [ ] 添加 ASCII 条形图显示各分类比例
- [ ] 保持现有 StatsResult 接口不变
- [ ] 运行 stats 命令可以看到美化后的状态码分布
- [ ] 百分比计算准确

## Blocked by

- .scratch/charm-ui-phase1/issues/02-beautify-total-requests.md

Status: ready-for-agent
