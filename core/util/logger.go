package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(env string) *zap.Logger {
	isDevelopment := env != "production"

	cfg := zap.Config{
		Development:      isDevelopment,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		Level:            zap.NewAtomicLevel(),
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "time",
			CallerKey:      "caller",
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	return logger
}
