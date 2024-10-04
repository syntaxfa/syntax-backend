package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/syntaxfa/syntax-backend/pkg/trace"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	l *slog.Logger
}

const (
	defaultFilePath        = "logs/logs.json"
	defaultUserLocalTime   = false
	defaultFileMaxSizeInMB = 10
	defaultFileAgeInDays   = 30
	defaultLogLevel        = slog.LevelInfo
)

type Config struct {
	FilePath         string
	UseLocalTime     bool
	FileMaxSizeInMB  int
	FileMaxAgeInDays int
}

var L *Logger

func init() {
	fileWriter := &lumberjack.Logger{
		Filename:  defaultFilePath,
		LocalTime: defaultUserLocalTime,
		MaxSize:   defaultFileMaxSizeInMB,
		MaxAge:    defaultFileAgeInDays,
	}
	L.l = slog.New(slog.NewJSONHandler(io.MultiWriter(fileWriter, os.Stdout), &slog.HandlerOptions{
		Level: defaultLogLevel,
	}))
}

func New(cfg Config, opt *slog.HandlerOptions, writeInConsole bool) *Logger {
	fileWriter := &lumberjack.Logger{
		Filename:  cfg.FilePath,
		LocalTime: cfg.UseLocalTime,
		MaxSize:   cfg.FileMaxSizeInMB,
		MaxAge:    cfg.FileMaxAgeInDays,
	}

	if writeInConsole {
		return &Logger{
			l: slog.New(slog.NewJSONHandler(io.MultiWriter(fileWriter, os.Stdout), opt)),
		}
	}

	return &Logger{
		l: slog.New(slog.NewJSONHandler(fileWriter, opt)),
	}
}

func (logger *Logger) Debug(msg string, args ...any) {
	logger.l.Debug(msg, args...)
}

func (logger *Logger) Info(msg string, args ...any) {
	logger.l.Info(msg, args...)
}

func (logger *Logger) Warn(msg string, args ...any) {
	logger.l.Warn(msg, args...)
}

func (logger *Logger) Error(msg string, args ...any) {
	t := trace.Parse()

	logger.l.With(slog.Group("trace",
		slog.String("path", t.File),
		slog.Int("line", t.Line),
		slog.String("function", t.Function),
	)).Error(msg, args...)
}
