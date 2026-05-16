package ui

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// SimpleSpinner 简单的 spinner 提示
type SimpleSpinner struct {
	message string
	running bool
	done    chan struct{}
	wg      sync.WaitGroup
}

// NewSimpleSpinner 创建新的 spinner
func NewSimpleSpinner(message string) *SimpleSpinner {
	return &SimpleSpinner{
		message: message,
		done:    make(chan struct{}),
	}
}

// Start 启动 spinner
func (s *SimpleSpinner) Start() {
	if s.running {
		return
	}
	s.running = true
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0

		for {
			select {
			case <-s.done:
				// 清除 spinner 输出
				fmt.Fprint(os.Stderr, "\r\033[K")
				return
			default:
				fmt.Fprintf(os.Stderr, "\r%s %s", chars[i], s.message)
				i = (i + 1) % len(chars)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

// Stop 停止 spinner
func (s *SimpleSpinner) Stop() {
	if !s.running {
		return
	}
	close(s.done)
	s.wg.Wait()
	s.running = false
}
