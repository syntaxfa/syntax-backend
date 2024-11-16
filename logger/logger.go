package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/syntaxfa/syntax-backend/pkg/trace"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultFilePath        = "logs/logs.json"
	defaultUserLocalTime   = false
	defaultFileMaxSizeInMB = 10
	defaultFileAgeInDays   = 30
	defaultLogLevel        = slog.LevelInfo
)

type Config struct {
	FilePath         string `koanf:"file_path"`
	UseLocalTime     bool   `koanf:"use_local_time"`
	FileMaxSizeInMB  int    `koanf:"file_max_size_in_mb"`
	FileMaxAgeInDays int    `koanf:"file_max_age_in_days"`
	LogLevel         int    `koanf:"log_level"`
}

var l *slog.Logger

func init() {
	fileWriter := &lumberjack.Logger{
		Filename:  defaultFilePath,
		LocalTime: defaultUserLocalTime,
		MaxSize:   defaultFileMaxSizeInMB,
		MaxAge:    defaultFileAgeInDays,
	}
	l = slog.New(slog.NewJSONHandler(io.MultiWriter(fileWriter, os.Stdout), &slog.HandlerOptions{
		Level: defaultLogLevel,
	}))
}

func L() *slog.Logger {
	return l
}

func LogError(err error) {
	if err == nil {
		return
	}

	L().Error(err.Error())
}

func WithGroup(groupName string) *slog.Logger {
	t := trace.Parse()

	return l.With(slog.String("group", groupName)).With(slog.Group("trace",
		slog.String("path", t.File),
		slog.Int("line", t.Line),
		slog.String("function", t.Function),
	))
}

func New(cfg Config, opt *slog.HandlerOptions, writeInConsole bool) *slog.Logger {
	fileWriter := &lumberjack.Logger{
		Filename:  cfg.FilePath,
		LocalTime: cfg.UseLocalTime,
		MaxSize:   cfg.FileMaxSizeInMB,
		MaxAge:    cfg.FileMaxAgeInDays,
	}

	if writeInConsole {
		return slog.New(slog.NewJSONHandler(io.MultiWriter(fileWriter, os.Stdout), opt))
	}

	return slog.New(slog.NewJSONHandler(fileWriter, opt))
}
