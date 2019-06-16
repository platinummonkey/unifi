package log

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

var logger *zap.Logger

// Get will get the current logger
func Get() *zap.Logger {
	return logger
}

// InitializeLogger will initialize the logger
func InitializeLogger(level string) (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.DisableStacktrace = true
	cfg.DisableCaller = true

	switch strings.TrimSpace(strings.ToLower(level)) {
	case "debug":
		cfg.Level.SetLevel(zap.DebugLevel)
	case "info":
		cfg.Level.SetLevel(zap.InfoLevel)
	case "warn":
		cfg.Level.SetLevel(zap.WarnLevel)
	case "error":
		cfg.Level.SetLevel(zap.ErrorLevel)
	case "fatal,none":
		cfg.Level.SetLevel(zap.FatalLevel)
	default:
		return nil, fmt.Errorf("invalid loglevel specified: %s", level)
	}
	var err error
	logger, err = cfg.Build()
	return logger, err
}
