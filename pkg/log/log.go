package log

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ssh/terminal"
)

var longestFullNameLength int
var longestShortNameLength int

const rootLoggerName = "CNC"

func GetRoot(level string) (*zap.Logger, error) {
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
		logConf.EncoderConfig.CallerKey = ""
		logConf.EncoderConfig.LevelKey = "capitalColor"
		logConf.EncoderConfig.EncodeName = PaddedFullNameEncoder
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
	registerMaxLength(rootLoggerName)
	return log.Named(rootLoggerName), nil
}

// Named is wrapping the zap.Logger.Named() method. It is needed to register logger names requested.
func Named(parent *zap.Logger, name string) *zap.Logger {
	if name == "" {
		return parent
	}

	child := parent.Named(name)

	// this hack is needed because there is no native way to get the name out of the logger
	checked := child.Check(zapcore.DebugLevel, "this is a dummy")
	registerMaxLength(checked.LoggerName)
	return child
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

// NamedLevelUp is a convenience method combining Named and LevelUp.
func NamedLevelUp(parent *zap.Logger, name, level string) *zap.Logger {
	return LevelUp(Named(parent, name), level)
}

func PaddedFullNameEncoder(loggerName string, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(fmt.Sprintf("%-"+strconv.Itoa(longestFullNameLength)+"s", loggerName))
}

func PaddedShortNameEncoder(loggerName string, encoder zapcore.PrimitiveArrayEncoder) {
	names := strings.Split(loggerName, ".")
	if len(names) > 0 {
		loggerName = names[len(names)-1]
	}

	encoder.AppendString(fmt.Sprintf("%-"+strconv.Itoa(longestShortNameLength)+"s", loggerName))
}

func registerMaxLength(loggerName string) {
	if len(loggerName) > longestFullNameLength {
		longestFullNameLength = len(loggerName)
	}

	names := strings.Split(loggerName, ".")
	if len(names[len(names)-1]) > longestShortNameLength {
		longestShortNameLength = len(names[len(names)-1])
	}
}
