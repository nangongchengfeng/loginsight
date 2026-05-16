# log-analyzer

一个灵活的日志分析 CLI 工具，目前支持 Nginx access log，后续可扩展支持 Linux 系统日志等多种格式。

## 功能特性

- 📊 **统计分析** - 总请求数、状态码分布、Top 10 IP/URL/User-Agent
- 🔍 **日志过滤** - 按状态码、IP、URL、HTTP 方法、时间范围过滤
- 🎨 **美观输出** - 使用 Charm 全家桶美化 CLI 输出
- 📝 **样本生成** - 快速生成测试用的日志样本

## 安装

```bash
git clone https://github.com/你的用户名/log-analyzer.git
cd log-analyzer
go build -o log-analyzer
```

## 使用说明

### 统计分析

```bash
# 查看统计信息
log-analyzer stats access.log
```

输出示例：
- 总请求数卡片
- 状态码分布（带颜色和条形图）
- Top 10 IP/URL/User-Agent

### 日志过滤

```bash
# 按状态码过滤
log-analyzer filter --status 404 access.log

# 按状态码分类过滤
log-analyzer filter --status 5xx access.log

# 按 IP 过滤
log-analyzer filter --ip 192.168.1.1 access.log

# 组合多个条件
log-analyzer filter --status 4xx --method POST access.log

# 按时间范围过滤 (RFC3339 格式)
log-analyzer filter --since 2026-05-15T00:00:00+08:00 --until 2026-05-16T00:00:00+08:00 access.log
```

### 生成样本日志

```bash
# 生成 100 条样本日志
log-analyzer gen-sample

# 生成到文件
log-analyzer gen-sample --lines 500 --output sample.log
```

## 项目结构

```
log-analyzer/
├── main.go              # 入口和 CLI 命令
├── go.mod/go.sum        # Go 模块定义
├── internal/
│   ├── parser/          # 日志解析模块
│   ├── analyzer/        # 统计分析模块
│   ├── filter/          # 过滤模块
│   └── ui/              # UI 美化模块
```

## 技术栈

- **CLI 框架**: [Cobra](https://github.com/spf13/cobra)
- **UI 美化**: [Charm Lipgloss](https://github.com/charmbracelet/lipgloss)
- **语言**: Go 1.25+

## 后续计划

- [ ] 支持更多日志格式（Linux 系统日志等）
- [ ] 更多统计维度
- [ ] 图表输出可视化
- [ ] 配置文件支持
- [ ] 自动检测日志格式

## License

MIT License
