package parser

import (
	"testing"
	"time"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		want    LogEntry
		wantErr bool
	}{
		{
			name: "正常日志行",
			line: `1.2.3.4 - - [15/May/2026:10:30:00 +0800] "GET /api/users HTTP/1.1" 200 1234 "https://google.com" "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"`,
			want: LogEntry{
				IP:        "1.2.3.4",
				Time:      mustParseTime("15/May/2026:10:30:00 +0800"),
				Method:    "GET",
				Path:      "/api/users",
				Protocol:  "HTTP/1.1",
				Status:    200,
				BodyBytes: 1234,
				Referer:   "https://google.com",
				UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			},
			wantErr: false,
		},
		{
			name: "404 错误日志",
			line: `192.168.1.1 - - [15/May/2026:12:00:00 +0800] "GET /not-found HTTP/1.1" 404 0 "-" "curl/7.68.0"`,
			want: LogEntry{
				IP:        "192.168.1.1",
				Time:      mustParseTime("15/May/2026:12:00:00 +0800"),
				Method:    "GET",
				Path:      "/not-found",
				Protocol:  "HTTP/1.1",
				Status:    404,
				BodyBytes: 0,
				Referer:   "-",
				UserAgent: "curl/7.68.0",
			},
			wantErr: false,
		},
		{
			name: "无效格式",
			line: `这不是有效的日志`,
			want: LogEntry{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.IP != tt.want.IP {
					t.Errorf("IP = %v, want %v", got.IP, tt.want.IP)
				}
				if got.Method != tt.want.Method {
					t.Errorf("Method = %v, want %v", got.Method, tt.want.Method)
				}
				if got.Path != tt.want.Path {
					t.Errorf("Path = %v, want %v", got.Path, tt.want.Path)
				}
				if got.Status != tt.want.Status {
					t.Errorf("Status = %v, want %v", got.Status, tt.want.Status)
				}
				if !got.Time.Equal(tt.want.Time) {
					t.Errorf("Time = %v, want %v", got.Time, tt.want.Time)
				}
			}
		})
	}
}

func mustParseTime(s string) time.Time {
	t, err := time.Parse("02/Jan/2006:15:04:05 -0700", s)
	if err != nil {
		panic(err)
	}
	return t
}
