package filter

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/charmbracelet/lipgloss"
	"log-analyzer/internal/parser"
	"log-analyzer/internal/ui"
)

// Options 存储过滤选项
type Options struct {
	Status string // 状态码过滤
	IP     string // IP 过滤
	Path   string // 路径过滤
	Method string // 方法过滤
	Since  string // 开始时间
	Until  string // 结束时间
}

// Match 检查日志条目是否匹配所有过滤条件
func Match(entry parser.LogEntry, opts Options) bool {
	return matchesStatus(entry, opts.Status) &&
		matchesIP(entry, opts.IP) &&
		matchesPath(entry, opts.Path) &&
		matchesMethod(entry, opts.Method) &&
		matchesTimeRange(entry, opts.Since, opts.Until)
}

// Run 执行过滤并输出匹配的行
func Run(filename string, opts Options) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		entry, err := parser.ParseLine(line)
		if err != nil {
			continue
		}
		if Match(entry, opts) {
			fmt.Println(renderColoredLine(line, entry.Status))
		}
	}

	return scanner.Err()
}

// statusCodeRegex 匹配日志中的状态码
var statusCodeRegex = regexp.MustCompile(`"(\S+)\s+(\S+)\s+(\S+)"\s+(\d+)`)

// renderColoredLine 渲染带颜色的日志行
func renderColoredLine(line string, status int) string {
	// 确定状态码颜色
	var color lipgloss.Color
	switch {
	case status >= 200 && status < 300:
		color = ui.ColorGreen
	case status >= 300 && status < 400:
		color = ui.ColorYellow
	case status >= 400 && status < 500:
		color = ui.ColorOrange
	case status >= 500 && status < 600:
		color = ui.ColorRed
	default:
		return line
	}

	// 替换状态码为带颜色的版本
	return statusCodeRegex.ReplaceAllStringFunc(line, func(m string) string {
		matches := statusCodeRegex.FindStringSubmatch(m)
		if len(matches) == 5 {
			statusPart := lipgloss.NewStyle().Foreground(color).Bold(true).Render(matches[4])
			return fmt.Sprintf(`"%s %s %s" %s`, matches[1], matches[2], matches[3], statusPart)
		}
		return m
	})
}

// matchesStatus 检查日志是否匹配状态码过滤条件
func matchesStatus(entry parser.LogEntry, statusFilter string) bool {
	if statusFilter == "" {
		return true
	}
	if len(statusFilter) == 3 && statusFilter[1] == 'x' && statusFilter[2] == 'x' {
		switch statusFilter[0] {
		case '2':
			return entry.Status >= 200 && entry.Status < 300
		case '3':
			return entry.Status >= 300 && entry.Status < 400
		case '4':
			return entry.Status >= 400 && entry.Status < 500
		case '5':
			return entry.Status >= 500 && entry.Status < 600
		}
	}
	s, err := strconv.Atoi(statusFilter)
	if err != nil {
		return false
	}
	return entry.Status == s
}

// matchesIP 检查日志是否匹配 IP 过滤条件
func matchesIP(entry parser.LogEntry, ipFilter string) bool {
	if ipFilter == "" {
		return true
	}
	return entry.IP == ipFilter
}

// matchesPath 检查日志是否匹配路径过滤条件
func matchesPath(entry parser.LogEntry, pathFilter string) bool {
	if pathFilter == "" {
		return true
	}
	return entry.Path == pathFilter
}

// matchesMethod 检查日志是否匹配方法过滤条件
func matchesMethod(entry parser.LogEntry, methodFilter string) bool {
	if methodFilter == "" {
		return true
	}
	return entry.Method == methodFilter
}

// matchesTimeRange 检查日志是否匹配时间范围过滤条件
func matchesTimeRange(entry parser.LogEntry, sinceStr, untilStr string) bool {
	if sinceStr == "" && untilStr == "" {
		return true
	}
	var since, until time.Time
	var err error
	if sinceStr != "" {
		since, err = time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			return false
		}
		if entry.Time.Before(since) {
			return false
		}
	}
	if untilStr != "" {
		until, err = time.Parse(time.RFC3339, untilStr)
		if err != nil {
			return false
		}
		if entry.Time.After(until) {
			return false
		}
	}
	return true
}
