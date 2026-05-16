package filter

import (
	"testing"

	"github.com/nangongchengfeng/go-cli/internal/parser"
)

func TestMatch(t *testing.T) {
	entry := parser.LogEntry{
		IP:     "192.168.1.100",
		Method: "GET",
		Path:   "/api/users",
		Status: 200,
	}

	tests := []struct {
		name string
		opts Options
		want bool
	}{
		{
			name: "无过滤条件",
			opts: Options{},
			want: true,
		},
		{
			name: "状态码匹配",
			opts: Options{Status: "200"},
			want: true,
		},
		{
			name: "状态码不匹配",
			opts: Options{Status: "404"},
			want: false,
		},
		{
			name: "状态码范围匹配 2xx",
			opts: Options{Status: "2xx"},
			want: true,
		},
		{
			name: "状态码范围不匹配 4xx",
			opts: Options{Status: "4xx"},
			want: false,
		},
		{
			name: "IP匹配",
			opts: Options{IP: "192.168.1.100"},
			want: true,
		},
		{
			name: "IP不匹配",
			opts: Options{IP: "1.1.1.1"},
			want: false,
		},
		{
			name: "路径匹配",
			opts: Options{Path: "/api/users"},
			want: true,
		},
		{
			name: "方法匹配",
			opts: Options{Method: "GET"},
			want: true,
		},
		{
			name: "多个条件AND匹配",
			opts: Options{Status: "200", Method: "GET", Path: "/api/users"},
			want: true,
		},
		{
			name: "多个条件AND不匹配",
			opts: Options{Status: "200", Method: "POST"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Match(entry, tt.opts); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchStatusCodeCategories(t *testing.T) {
	tests := []struct {
		name   string
		status int
		filter string
		want   bool
	}{
		{"200 matches 2xx", 200, "2xx", true},
		{"201 matches 2xx", 201, "2xx", true},
		{"301 matches 3xx", 301, "3xx", true},
		{"404 matches 4xx", 404, "4xx", true},
		{"500 matches 5xx", 500, "5xx", true},
		{"404 does not match 2xx", 404, "2xx", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := parser.LogEntry{Status: tt.status}
			opts := Options{Status: tt.filter}
			if got := Match(entry, opts); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
