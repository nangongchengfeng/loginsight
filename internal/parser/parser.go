package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/satyrius/gonx"
	"log-analyzer/internal/ui"
)

// LogEntry 表示一条 Nginx 日志的结构化数据
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

// Nginx 日志格式定义
const (
	// NginxCombinedFormat Nginx combined 日志格式
	NginxCombinedFormat = `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`
)

// parser 全局 gonx 解析器
var parser *gonx.Parser

func init() {
	// 初始化 gonx 解析器
	parser = gonx.NewParser(NginxCombinedFormat)
}

// ParseFile 读取并解析整个日志文件
func ParseFile(filename string) ([]LogEntry, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 获取文件大小决定是否显示 spinner
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	var spinner *ui.SimpleSpinner
	// 文件大于 1MB 时显示 spinner
	if info.Size() > 1024*1024 {
		spinner = ui.NewSimpleSpinner("正在解析日志文件...")
		spinner.Start()
		defer spinner.Stop()
	}

	var entries []LogEntry
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		entry, err := ParseLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "警告: 第 %d 行: %v\n", lineNum, err)
			continue
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// ParseLine 解析单条日志行
func ParseLine(line string) (LogEntry, error) {
	// 使用 gonx 解析
	record, err := parser.ParseString(line)
	if err != nil {
		return LogEntry{}, fmt.Errorf("无法解析日志行: %w", err)
	}

	// 提取字段
	entry := LogEntry{}
	entry.IP, _ = record.Field("remote_addr")

	// 解析时间
	timeStr, _ := record.Field("time_local")
	if timeStr != "" {
		t, err := time.Parse("02/Jan/2006:15:04:05 -0700", timeStr)
		if err != nil {
			return LogEntry{}, fmt.Errorf("无效的时间: %w", err)
		}
		entry.Time = t
	}

	// 解析 request (method path protocol)
	request, _ := record.Field("request")
	if request != "" {
		parts := splitRequest(request)
		if len(parts) >= 1 {
			entry.Method = parts[0]
		}
		if len(parts) >= 2 {
			entry.Path = parts[1]
		}
		if len(parts) >= 3 {
			entry.Protocol = parts[2]
		}
	}

	// 解析状态码
	statusStr, _ := record.Field("status")
	if statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err != nil {
			return LogEntry{}, fmt.Errorf("无效的状态码: %w", err)
		}
		entry.Status = status
	}

	// 解析 body_bytes_sent
	bytesStr, _ := record.Field("body_bytes_sent")
	if bytesStr != "" {
		bodyBytes, err := strconv.Atoi(bytesStr)
		if err != nil {
			return LogEntry{}, fmt.Errorf("无效的字节数: %w", err)
		}
		entry.BodyBytes = bodyBytes
	}

	entry.Referer, _ = record.Field("http_referer")
	entry.UserAgent, _ = record.Field("http_user_agent")

	return entry, nil
}

// splitRequest 将 request 字符串拆分为 method, path, protocol
func splitRequest(request string) []string {
	var parts []string
	current := ""
	inQuotes := false

	for _, c := range request {
		if c == '"' {
			inQuotes = !inQuotes
		} else if c == ' ' && !inQuotes {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
