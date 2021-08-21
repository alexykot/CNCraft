package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const testLoggerName = "t"

func MustGetTest() *zap.Logger {
	log, err := GetLevelledTest("DEBUG")
	if err != nil {
		panic(err)
	}
	return log
}

func MustGetTestNamed(name string) *zap.Logger {
	log, err := GetLevelledTest("DEBUG")
	if err != nil {
		panic(err)
	}
	return log.Named(name)
}

func GetLevelledTest(level string) (*zap.Logger, error) {
	switch level {
	case "DEBUG", "INFO", "WARN", "ERROR":
	default:
		return nil, fmt.Errorf("error level `%s` not recognised, must be one of `DEBUG`, `INFO`, `WARN`, `ERROR`", level)
	}

	var logConf zap.Config
	logConf = zap.NewDevelopmentConfig()
	logConf.EncoderConfig.TimeKey = ""
	logConf.EncoderConfig.CallerKey = ""
	logConf.EncoderConfig.LevelKey = "capitalColor"
	logConf.EncoderConfig.EncodeName = PaddedFullNameEncoder
	logConf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logLevel := zap.NewAtomicLevel()
	_ = logLevel.UnmarshalText([]byte(level))
	logConf.Level = logLevel

	log, err := logConf.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build dev logger: %w", err)
	}

	// Set the global logger, this should be an unnamed logger
	zap.ReplaceGlobals(log)

	return log.Named(testLoggerName), nil
}
