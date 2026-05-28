package config

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/goapt/logger"
)

func NewLogger(conf *Config, filename string) *slog.Logger {
	logLevel := conf.GetLogLevel()
	if conf.Env == "dev" {
		return logger.New(
			logger.NewJSONHandler(os.Stdout, logger.WithLevel(logLevel)),
		)
	}
	return logger.New(
		logger.NewJSONHandler(os.Stdout, logger.WithLevel(logLevel)),
		logger.NewJSONHandler(
			logger.NewFileWriter(
				filepath.Join(conf.Common.StoragePath, "logs", filename),
				logger.WithMaxFiles(3),
			),
			logger.WithLevel(logLevel),
			logger.WithSource(),
		),
	)
}
