package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
	"log-analyzer/internal/analyzer"
	"log-analyzer/internal/filter"
	"log-analyzer/internal/parser"
)

// main 程序入口
func main() {
	rootCmd := &cobra.Command{
		Use:   "log-analyzer",
		Short: "Nginx 日志分析工具",
		Long: `log-analyzer 是一个用于分析 Nginx access.log 的命令行工具。

支持三种主要功能：
  - stats:    显示日志统计信息（总请求数、状态码分布、Top 10 IP/URL/User-Agent）
  - filter:   按条件过滤日志条目
  - gen-sample: 生成 Nginx 访问日志样本用于测试`,
		Example: `  # 查看统计信息
  log-analyzer stats access.log

  # 过滤 4xx 错误
  log-analyzer filter --status 4xx access.log

  # 生成 200 行样本日志
  log-analyzer gen-sample --lines 200 --output sample.log`,
	}

	rootCmd.AddCommand(newGenSampleCmd())
	rootCmd.AddCommand(newStatsCmd())
	rootCmd.AddCommand(newFilterCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// newStatsCmd 创建 stats 子命令
func newStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats <log-file>",
		Short: "显示日志统计信息",
		Long: `显示 Nginx 访问日志的统计信息，包括：
  - 总请求数
  - 状态码分布（2xx/3xx/4xx/5xx）
  - Top 10 IP 地址
  - Top 10 请求 URL
  - Top 10 用户代理`,
		Example: `  # 查看 access.log 的统计信息
  log-analyzer stats access.log

  # 从标准输入读取（配合管道使用）
  cat access.log | log-analyzer stats -`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStats(args[0])
		},
	}

	return cmd
}

// newFilterCmd 创建 filter 子命令
func newFilterCmd() *cobra.Command {
	var opts filter.Options

	cmd := &cobra.Command{
		Use:   "filter <log-file>",
		Short: "过滤日志条目",
		Long: `按多种条件过滤 Nginx 访问日志条目，所有条件是 AND 关系。

支持的过滤条件：
  - 状态码（精确匹配或 2xx/3xx/4xx/5xx 分类）
  - 客户端 IP 地址
  - 请求路径（精确匹配）
  - HTTP 方法（GET、POST、PUT、DELETE 等）
  - 时间范围（RFC3339 格式）`,
		Example: `  # 过滤 404 错误
  log-analyzer filter --status 404 access.log

  # 过滤所有 4xx 错误
  log-analyzer filter --status 4xx access.log

  # 过滤来自特定 IP 的请求
  log-analyzer filter --ip 192.168.1.100 access.log

  # 过滤 POST 请求到 /api/users
  log-analyzer filter --method POST --path /api/users access.log

  # 过滤时间范围内的请求（RFC3339 格式）
  log-analyzer filter --since 2026-05-15T00:00:00+08:00 --until 2026-05-16T00:00:00+08:00 access.log

  # 组合多个条件：5xx 错误 + 来自 10.0.0.1
  log-analyzer filter --status 5xx --ip 10.0.0.1 access.log`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return filter.Run(args[0], opts)
		},
	}

	cmd.Flags().StringVar(&opts.Status, "status", "", "按状态码过滤，例如 404 或 4xx")
	cmd.Flags().StringVar(&opts.IP, "ip", "", "按客户端 IP 过滤，例如 192.168.1.100")
	cmd.Flags().StringVar(&opts.Path, "path", "", "按请求路径过滤，精确匹配，例如 /api/users")
	cmd.Flags().StringVar(&opts.Method, "method", "", "按 HTTP 方法过滤，例如 GET、POST")
	cmd.Flags().StringVar(&opts.Since, "since", "", "按开始时间过滤，RFC3339 格式，例如 2006-01-02T15:04:05+08:00")
	cmd.Flags().StringVar(&opts.Until, "until", "", "按结束时间过滤，RFC3339 格式，例如 2006-01-02T15:04:05+08:00")

	return cmd
}

// runStats 执行统计逻辑
func runStats(filename string) error {
	entries, err := parser.ParseFile(filename)
	if err != nil {
		return err
	}

	result := analyzer.Analyze(entries)
	analyzer.PrintStats(result)

	return nil
}

// newGenSampleCmd 创建 gen-sample 子命令
func newGenSampleCmd() *cobra.Command {
	var lines int
	var output string

	cmd := &cobra.Command{
		Use:   "gen-sample",
		Short: "生成 Nginx 访问日志样本",
		Long: `生成 Nginx combined 格式的访问日志样本，用于测试工具功能。

样本日志包含随机的：
  - IP 地址
  - 时间戳（最近 24 小时内）
  - HTTP 方法
  - 请求路径
  - 状态码（2xx 出现概率更高）
  - User-Agent`,
		Example: `  # 生成 100 行样本日志到标准输出
  log-analyzer gen-sample

  # 生成 500 行样本日志到文件
  log-analyzer gen-sample --lines 500 --output sample.log

  # 使用短选项
  log-analyzer gen-sample -n 200 -o sample.log`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateSample(lines, output)
		},
	}

	cmd.Flags().IntVarP(&lines, "lines", "n", 100, "生成的日志行数")
	cmd.Flags().StringVarP(&output, "output", "o", "", "输出文件，默认输出到标准输出")

	return cmd
}

// generateSample 生成样本日志
func generateSample(lines int, output string) error {
	var w io.Writer = os.Stdout
	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < lines; i++ {
		line := generateLogLine(r)
		fmt.Fprintln(w, line)
	}

	return nil
}

// generateLogLine 生成单条样本日志行
func generateLogLine(r *rand.Rand) string {
	ip := generateIP(r)
	ts := time.Now().Add(-time.Duration(r.Intn(86400)) * time.Second)
	method := pickMethod(r)
	path := pickPath(r)
	proto := "HTTP/1.1"
	status := pickStatus(r)
	bodyBytes := r.Intn(10000)
	referer := pickReferer(r)
	ua := pickUserAgent(r)

	return fmt.Sprintf(`%s - - [%s] "%s %s %s" %d %d "%s" "%s"`,
		ip,
		ts.Format("02/Jan/2006:15:04:05 -0700"),
		method,
		path,
		proto,
		status,
		bodyBytes,
		referer,
		ua,
	)
}

// generateIP 生成随机 IP 地址
func generateIP(r *rand.Rand) string {
	// 从固定 IP 池中选择，产生重复
	ips := []string{
		"192.168.1.1",
		"192.168.1.2",
		"10.0.0.1",
		"10.0.0.2",
		"172.16.0.1",
		"172.16.0.2",
		"8.8.8.8",
		"1.1.1.1",
		"203.0.113.1",
		"198.51.100.1",
	}
	return ips[r.Intn(len(ips))]
}

// pickMethod 随机选择 HTTP 方法
func pickMethod(r *rand.Rand) string {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	return methods[r.Intn(len(methods))]
}

// pickPath 随机选择请求路径
func pickPath(r *rand.Rand) string {
	paths := []string{
		"/",
		"/api/users",
		"/api/posts",
		"/static/css/style.css",
		"/static/js/app.js",
		"/images/logo.png",
		"/about",
		"/contact",
		"/blog/post-1",
		"/blog/post-2",
		"/api/products",
		"/cart",
		"/checkout",
		"/login",
		"/signup",
	}
	return paths[r.Intn(len(paths))]
}

// pickStatus 随机选择状态码（200 出现概率更高）
func pickStatus(r *rand.Rand) int {
	statuses := []int{200, 200, 200, 200, 200, 201, 301, 302, 400, 401, 403, 404, 500}
	return statuses[r.Intn(len(statuses))]
}

// pickReferer 随机选择来源页
func pickReferer(r *rand.Rand) string {
	if r.Float32() < 0.3 {
		return "-"
	}
	referers := []string{
		"https://google.com",
		"https://bing.com",
		"https://example.com",
		"https://example.com/blog",
		"https://example.com/about",
	}
	return referers[r.Intn(len(referers))]
}

// pickUserAgent 随机选择用户代理
func pickUserAgent(r *rand.Rand) string {
	uas := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"curl/7.68.0",
		"Go-http-client/1.1",
	}
	return uas[r.Intn(len(uas))]
}
