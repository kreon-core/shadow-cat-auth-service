package clog

import (
	"fmt"
	"os"
	"slices"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	zapLog   *zap.Logger        //nolint:gochecknoglobals // global logger instance
	zapSugar *zap.SugaredLogger //nolint:gochecknoglobals // global sugared logger instance
)

func LoadZap() {
	var err error
	var baseLogger *zap.Logger

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:       "time",
			MessageKey:    "message",
			LevelKey:      "level",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		},
	}

	if slices.Contains([]string{"prod", "stag", "production", "staging"}, os.Getenv("ENV")) {
		cfg.Encoding = "json"
	} else {
		cfg.Development = true
		cfg.Encoding = "console"
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	baseLogger, err = cfg.Build(zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize zap logger: %v", err))
	}
	zapLog = baseLogger
	zapSugar = baseLogger.Sugar()
}

func CloseZap() {
	if zapLog != nil {
		_ = zapLog.Sync()
	}
}

func Log() *zap.Logger {
	if zapLog == nil {
		panic("Zap logger is not initialized. Call clog.InitZap first")
	}
	return zapLog
}

func Sugar() *zap.SugaredLogger {
	if zapSugar == nil {
		panic("Zap logger is not initialized. Call clog.InitZap first")
	}
	return zapSugar
}
