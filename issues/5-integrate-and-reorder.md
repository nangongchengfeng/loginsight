## What to build

调整 stats 命令的输出顺序，集成所有新增指标，确保整体展示流畅美观。

## Acceptance criteria

- [x] 输出顺序调整为：1.总请求数 2.状态码（分类+Top10）3.错误率 4.时间分布 5.IP统计（唯一数+Top10）6.URL 7.User-Agent
- [x] 各指标之间有适当空行分隔
- [x] 整体布局保持一致的风格
- [x] 运行完整 stats 命令验证所有指标正常显示
- [x] 使用样本日志测试验证

## Blocked by

- 1-unique-ip-count
- 2-top-status-codes
- 3-error-rates
- 4-hourly-distribution

## Status

✅ 已完成（输出顺序已经在前面的提交中调整好）
