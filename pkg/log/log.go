package log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ssh/terminal"
)

const rootLoggerName = "CNCraft"

func GetRootLogger(level string) (*zap.Logger, error) {
	var err error
	var log *zap.Logger

	switch level {
	case "DEBUG", "INFO", "WARN", "ERROR":
	default:
		return nil, fmt.Errorf("error level `%s` not recognised, must be one of `DEBUG`, `INFO`, `WARN`, `ERROR`", level)
	}

	var logConf zap.Config
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		logConf = zap.NewDevelopmentConfig()

		logConf.EncoderConfig.TimeKey = ""
		logConf.EncoderConfig.NameKey = ""
		logConf.EncoderConfig.CallerKey = ""
		logConf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		logConf = zap.NewProductionConfig()
	}

	logLevel := zap.NewAtomicLevel()
	if err = logLevel.UnmarshalText([]byte(level)); err != nil {
		return nil, err
	}
	logConf.Level = logLevel

	if log, err = logConf.Build(); err != nil {
		return nil, err
	}

	// Set the global logger, this should be an unnamed logger
	zap.ReplaceGlobals(log)
	return log.Named(rootLoggerName), nil
}

// LevelUp creates a clone of the supplied logger with the new increased log level.
func LevelUp(logger *zap.Logger, level string) *zap.Logger {
	levelEnabler := zap.NewAtomicLevel()
	if err := levelEnabler.UnmarshalText([]byte(level)); err != nil {
		logger.Error("failed to increase level for logger: failed to parse log level", zap.Error(err))
		return logger
	}

	return logger.WithOptions(zap.IncreaseLevel(levelEnabler))
}
