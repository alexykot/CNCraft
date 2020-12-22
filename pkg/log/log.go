package log

import (
	"os"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ssh/terminal"
)

const DefaultLevel = "INFO"
const DefaultLogName = "CNCraft"

func GetLogger(level string) (*zap.Logger, error) {
	var err error
	var log *zap.Logger

	if level == "" {
		level = DefaultLevel
	}

	var logConf zap.Config
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		logConf = zap.NewDevelopmentConfig()

		logConf.EncoderConfig.TimeKey = ""
		logConf.EncoderConfig.NameKey = ""
		logConf.EncoderConfig.CallerKey = ""
		logConf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		logConf = zapdriver.NewProductionConfig()
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
	return log.Named(DefaultLogName), nil
}
