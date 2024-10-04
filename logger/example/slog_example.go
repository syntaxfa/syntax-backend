package main

import (
	"io"
	"log/slog"
	"os"
)

func main() {
	//LoggerWithTextHandler()
	//LoggerWithJsonHandler()
	//LoggerWithMultiHandler()
	LoggerGroup()
}

func LoggerWithTextHandler() {
	file, err := os.OpenFile("./logger/example/log.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}

	var logger *slog.Logger = slog.New(slog.NewTextHandler(file, nil))

	logger.Info("New log", "category", "startapp", "sub_category", "logger")
}

func LoggerWithJsonHandler() {
	file, err := os.OpenFile("./logger/example/log.json", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)

	if err != nil {
		panic(err)
	}

	var logger *slog.Logger = slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}))

	logger.Error("an error occurred", "category", "startapp", "sub_category", "LoggerWithJsonHandler")
}

func LoggerWithMultiHandler() {
	file, err := os.OpenFile("./logger/example/log.json", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}

	var logger *slog.Logger = slog.New(slog.NewJSONHandler(io.MultiWriter(file, os.Stdout),
		&slog.HandlerOptions{Level: slog.LevelInfo}))

	logger.Debug("an debug log", "category", "startapp", "sub_category", "LoggerWithMultiHandler")
	logger.Info("an info log", "category", "startapp", "sub_category", "LoggerWithMultiHandler")
	logger.Error("an error log", "category", "startapp", "sub_category", "LoggerWithMultiHandler")
}

func LoggerGroup() {
	file, err := os.OpenFile("./logger/example/log.json", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	var logger *slog.Logger = slog.New(slog.NewJSONHandler(io.MultiWriter(file, os.Stdout), nil))

	logger.Info("test log with group", slog.Group("user_info", "name", "alireza", "email", "alireza@gmail.com"))
}
