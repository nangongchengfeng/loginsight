# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

log-analyzer 是一个用 Go 编写的 Nginx 日志分析 CLI 工具。支持 Nginx combined 格式和自定义 JSON 格式，提供统计分析、日志过滤和样本生成功能。

## 常用命令

### 构建
```bash
go build -o log-analyzer
```

### 运行测试
```bash
# 运行所有测试
go test ./...

# 运行单个包测试
go test ./internal/parser
go test ./internal/analyzer
go test ./internal/filter

# 查看测试覆盖率
go test -cover ./...
```

### 运行工具
```bash
# 生成样本日志
go run . gen-sample --lines 100 --output test.log

# 查看统计信息
go run . stats nginx.log

# 过滤日志
go run . filter --status 4xx nginx.log
```

## 代码架构

### 目录结构
```
log-analyzer/
├── main.go                    # 入口点，Cobra CLI 命令定义
├── internal/
│   ├── parser/               # 日志解析模块
│   │   ├── parser.go         # 核心解析逻辑（支持 combined/JSON 两种格式）
│   │   └── parser_test.go
│   ├── analyzer/             # 统计分析模块
│   │   ├── analyzer.go       # 计算统计数据（Top 10、状态码分布、时间分布等）
│   │   └── analyzer_test.go
│   ├── filter/               # 日志过滤模块
│   │   ├── filter.go         # 过滤条件匹配与输出
│   │   └── filter_test.go
│   └── ui/                   # UI 美化模块
│       ├── styles.go         # Lipgloss 样式定义
│       ├── render.go         # 统计结果渲染（卡片、条形图）
│       └── spinner.go        # 加载动画
```

### 核心数据结构

**parser.LogEntry** - 统一的日志条目结构，两种解析格式都转换为此结构：
```go
type LogEntry struct {
    IP         string        // 客户端 IP
    Time       time.Time     // 请求时间
    Method     string        // HTTP 方法
    Path       string        // 请求路径
    Protocol   string        // HTTP 协议版本
    Status     int           // 状态码
    BodyBytes  int           // 响应体字节数
    Referer    string        // 来源页
    UserAgent  string        // 用户代理
}
```

**analyzer.StatsResult** - 统计结果数据结构：
```go
type StatsResult struct {
    Total         int           // 总请求数
    StatusCounts  map[int]int   // 具体状态码统计
    CategoryCounts map[string]int // 状态码分类统计（2xx/3xx/4xx/5xx）
    IPCounts      map[string]int // IP 统计
    PathCounts    map[string]int // URL 统计
    UACounts      map[string]int // User-Agent 统计
    UniqueIPCount int           // 唯一 IP 数
    HourlyCounts  map[int]int   // 按小时时间分布
}
```

### Stats 输出模块顺序

1. **仪表盘** - 总请求数、唯一 IP、4xx 率、5xx 率
2. **状态码卡片** - 分布条形图 + Top 10 具体状态码
3. **时间分布卡片** - 24小时逐行展示
4. **Top 10 IP 卡片** - 包含唯一 IP 总数
5. **Top 10 URL 卡片**
6. **Top 10 User-Agent 卡片**

### 关键设计点

1. **自动格式检测**：`parser.ParseFile()` 读取第一行判断是 JSON 还是 combined 格式
2. **流式 vs 全量**：`filter.Run()` 流式处理（适合大文件），`parser.ParseFile()` 全量加载
3. **颜色渲染**：根据状态码自动着色（2xx=绿，3xx=黄，4xx=橙，5xx=红）
4. **卡片式 UI**：使用 Lipgloss 实现带边框的美观卡片输出
5. **表格渲染优化**：表头避免使用 Lipgloss 渲染以防止自动换行问题

### 依赖库
- **CLI**: github.com/spf13/cobra
- **Nginx 解析**: github.com/satyrius/gonx
- **UI美化**: github.com/charmbracelet/lipgloss

## 开发要点

- 新增功能优先考虑在对应子包中添加，保持 main.go 只负责 CLI 绑定
- 新增解析格式时，扩展 parser 包并保持 LogEntry 统一接口
- 添加测试时参考现有测试文件的模式
- UI渲染相关修改应在 ui 包中进行，analyzer 包只负责数据处理
- 卡片宽度需要仔细测试以避免内容换行问题
