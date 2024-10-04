package logger_test

import (
	"bufio"
	"github.com/syntaxfa/syntax-backend/logger"
	"log/slog"
	"os"
	"strconv"
	"testing"
)

func TestLogger(t *testing.T) {
	cfg := logger.Config{
		FilePath:         "./logs.json",
		UseLocalTime:     false,
		FileMaxSizeInMB:  10,
		FileMaxAgeInDays: 1,
	}

	opt := slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}

	l := logger.New(cfg, &opt, true)

	tests := []struct {
		f    func()
		want string
	}{
		{
			f: func() {
				l.Info("an info log", "key", "value")
			},
			want: `{"level":"INFO","msg":"an info log","key":"value"}`,
		},
		{
			f: func() {
				l.Warn("WARN", "key", "value")
			},
			want: `{"level":"WARN","msg":"WARN","key":"value"}`,
		},
		{
			f: func() {
				l.Error("ERROR", "key", "value")
			},
			want: `{"level":"ERROR","msg":"ERROR","trace":{"path":"/home/ali/projects_golang/syntax/logger/logger_test.go","line":50,"function":"github.com/syntaxfa/syntax-backend/logger_test.TestLogger.func4"},"key":"value"}`,
		},
	}

	// first run logs
	for _, test := range tests {
		test.f()
	}

	f, err := os.Open(cfg.FilePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	defer func() {
		if err := os.Remove(cfg.FilePath); err != nil {
			panic(err)
		}
	}()

	scanner := bufio.NewScanner(f)

	var logs []string

	for scanner.Scan() {
		logs = append(logs, scanner.Text())
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if test.want != logs[i] {
				t.Fatalf("want: %+v, got: %+v", test.want, logs[i])
			}
		})
	}
}
