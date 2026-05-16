package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
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

// nginxJSONLog 用户自定义的 Nginx JSON 日志格式
type nginxJSONLog struct {
	Timestamp      string `json:"timestamp"`
	RemoteAddr     string `json:"remote_addr"`
	XForwardedFor  string `json:"x_forwarded_for"`
	Host           string `json:"host"`
	RequestMethod  string `json:"request_method"`
	RequestURI     string `json:"request_uri"`
	Status         int    `json:"status"`
	BodyBytesSent  int    `json:"body_bytes_sent"`
	HTTPReferer    string `json:"http_referer"`
	HTTPUserAgent string `json:"http_user_agent"`
}

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

	// 先读取第一行来检测格式
	scanner := bufio.NewScanner(f)
	var firstLine string
	if scanner.Scan() {
		firstLine = scanner.Text()
	}
	// 检测格式
	isJSON := strings.HasPrefix(strings.TrimSpace(firstLine), "{")

	// 重置文件指针
	f.Seek(0, 0)

	var entries []LogEntry
	scanner = bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		var entry LogEntry
		var err error

		if isJSON {
			entry, err = ParseJSONLine(line)
		} else {
			entry, err = ParseLine(line)
		}

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

// ParseJSONLine 解析 JSON 格式日志行
func ParseJSONLine(line string) (LogEntry, error) {
	var j nginxJSONLog
	err := json.Unmarshal([]byte(line), &j)
	if err != nil {
		return LogEntry{}, fmt.Errorf("无法解析 JSON 日志行: %w", err)
	}

	entry := LogEntry{}
	entry.IP = j.RemoteAddr

	// 解析时间 (ISO8601 格式)
	if j.Timestamp != "" {
		t, err := time.Parse(time.RFC3339Nano, j.Timestamp)
		if err != nil {
			// 尝试解析不带时区的格式
			t, err = time.Parse("2006-01-02T15:04:05", j.Timestamp)
			if err != nil {
				return LogEntry{}, fmt.Errorf("无效的时间: %w", err)
			}
		}
		entry.Time = t
	}

	entry.Method = j.RequestMethod
	entry.Path = j.RequestURI
	entry.Status = j.Status
	entry.BodyBytes = j.BodyBytesSent
	entry.Referer = j.HTTPReferer
	entry.UserAgent = j.HTTPUserAgent

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
