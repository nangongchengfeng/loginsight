package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

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

// logRegex 用于解析 Nginx combined 格式日志的正则表达式
var logRegex = regexp.MustCompile(`^(\S+) - - \[([^\]]+)\] "(\S+) (\S+) (\S+)" (\d+) (\d+) "([^"]+)" "([^"]+)"$`)

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
	matches := logRegex.FindStringSubmatch(line)
	if len(matches) != 10 {
		return LogEntry{}, fmt.Errorf("无法解析日志行")
	}

	status, err := strconv.Atoi(matches[6])
	if err != nil {
		return LogEntry{}, fmt.Errorf("无效的状态码: %v", err)
	}

	bodyBytes, err := strconv.Atoi(matches[7])
	if err != nil {
		return LogEntry{}, fmt.Errorf("无效的字节数: %v", err)
	}

	t, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[2])
	if err != nil {
		return LogEntry{}, fmt.Errorf("无效的时间: %v", err)
	}

	return LogEntry{
		IP:         matches[1],
		Time:       t,
		Method:     matches[3],
		Path:       matches[4],
		Protocol:   matches[5],
		Status:     status,
		BodyBytes:  bodyBytes,
		Referer:    matches[8],
		UserAgent:  matches[9],
	}, nil
}
