# log-analyzer

一个灵活的日志分析 CLI 工具，目前支持 Nginx access log，后续可扩展支持 Linux 系统日志等多种格式。

## 功能特性

- 📊 **统计分析** - 总请求数、状态码分布、Top 10 IP/URL/User-Agent、唯一 IP 数
- 🎯 **仪表盘** - 一目了然的核心指标：总请求数、唯一 IP、4xx 率、5xx 率
- 📈 **时间分布** - 24小时逐行展示流量分布，带可视化条形图
- 🎨 **卡片式 UI** - 使用 Charm Lipgloss 打造的美观 CLI 输出
- 🔍 **日志过滤** - 按状态码、IP、URL、HTTP 方法、时间范围过滤
- 📝 **样本生成** - 快速生成测试用的日志样本
- 🔄 **格式自动检测** - 自动识别并支持 Nginx combined 和自定义 JSON 格式

## 安装

```bash
git clone https://github.com/你的用户名/log-analyzer.git
cd log-analyzer
go build -o log-analyzer
```

## 使用说明

### 支持的日志格式

#### 1. Nginx combined 格式（默认）

```
$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"
```

#### 2. 自定义 JSON 格式

在 nginx.conf 中配置：

```nginx
log_format main escape=json '{'
    '"timestamp":"$time_iso8601",'
    '"remote_addr":"$remote_addr",'
    '"x_forwarded_for":"$http_x_forwarded_for",'
    '"host":"$host",'
    '"server_name":"$server_name",'
    '"request_method":"$request_method",'
    '"request_uri":"$request_uri",'
    '"status":$status,'
    '"body_bytes_sent":$body_bytes_sent,'
    '"bytes_sent":$bytes_sent,'
    '"request_time":$request_time,'
    '"upstream_addr":"$upstream_addr",'
    '"upstream_status":"$upstream_status",'
    '"upstream_response_time":"$upstream_response_time",'
    '"http_referer":"$http_referer",'
    '"http_user_agent":"$http_user_agent",'
    '"scheme":"$scheme",'
    '"connection":$connection,'
    '"connection_requests":$connection_requests,'
    '"request_length":$request_length'
'}';
```

工具会**自动检测**第一行是否为 JSON 格式，自动选择对应的解析器。

### 统计分析

```bash
# 查看统计信息
log-analyzer stats access.log
```

输出示例：
- 🎯 仪表盘卡片 - 总请求数、唯一 IP、4xx 率、5xx 率
- 📊 状态码分布卡片 - 分类统计 + Top 10 具体状态码（带颜色和条形图）
- ⏰ 时间分布卡片 - 24小时逐行展示，每条带条形图和百分比
- 📈 Top 10 IP 卡片 - 含唯一 IP 总数
- 📈 Top 10 URL 卡片
- 📈 Top 10 User-Agent 卡片

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

- [x] 自动检测日志格式
- [x] 支持自定义 JSON 格式
- [x] 唯一 IP 数统计
- [x] Top 10 具体状态码
- [x] 时间分布统计
- [x] 仪表盘展示
- [x] 卡片式 UI 优化
- [ ] 支持更多日志格式（Linux 系统日志等）
- [ ] 更多统计维度
- [ ] 配置文件支持

## License

MIT License
