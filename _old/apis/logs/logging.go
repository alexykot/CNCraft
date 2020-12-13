// Package logs - why the fuck this invents an own logger?
package logs

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"

	"github.com/golangmc/minecraft-server/apis/base"
	"github.com/golangmc/minecraft-server/apis/data/chat"
)

type logLevel int

const (
	Debug logLevel = iota
	Info
	Warn
	Error
)

var BasicLevel = []logLevel{Info, Warn, Error}
var EveryLevel = []logLevel{Info, Warn, Error, Debug}

type Logging struct {
	name   string
	writer io.Writer
	show   []logLevel
}

func (log *Logging) Name() string {
	return log.name
}

func (log *Logging) Show() []logLevel {
	return log.show
}

func (log *Logging) formatPrint(level, message string) {
	_, _ = fmt.Fprint(log.writer,
		fmt.Sprintf("[%s] [%s] [%s] %s\n", color.HiGreenString(currentTimeAsText()), level,
			color.WhiteString(log.Name()), chat.TranslateConsole(message)))
}

func (log *Logging) data(message string) {
	log.formatPrint(color.HiBlackString("DATA"), message)
}

func (log *Logging) info(message string) {
	log.formatPrint(color.CyanString("INFO"), message)
}

func (log *Logging) warn(message string) {
	log.formatPrint(color.YellowString("WARN"), message)
}

func (log *Logging) error(message string) {
	log.formatPrint(color.RedString("FAIL"), message)
}

func (log *Logging) Data(message ...interface{}) {
	if !checkIfLevelShows(log, Debug) {
		return
	}

	log.data(base.ConvertToString(message...))
}

func (log *Logging) Info(message ...interface{}) {
	if !checkIfLevelShows(log, Info) {
		return
	}

	log.info(base.ConvertToString(message...))
}

func (log *Logging) Warn(message ...interface{}) {
	if !checkIfLevelShows(log, Warn) {
		return
	}

	log.warn(base.ConvertToString(message...))
}

func (log *Logging) Error(message ...interface{}) {
	if !checkIfLevelShows(log, Error) {
		return
	}

	log.error(base.ConvertToString(message...))
}

func (log *Logging) DebugF(format string, a ...interface{}) {
	if !checkIfLevelShows(log, Debug) {
		return
	}

	log.data(fmt.Sprintf(format, a...))
}

func (log *Logging) InfoF(format string, a ...interface{}) {
	if !checkIfLevelShows(log, Info) {
		return
	}

	log.info(fmt.Sprintf(format, a...))
}

func (log *Logging) WarnF(format string, a ...interface{}) {
	if !checkIfLevelShows(log, Warn) {
		return
	}

	log.warn(fmt.Sprintf(format, a...))
}

func (log *Logging) ErrorF(format string, a ...interface{}) {
	if !checkIfLevelShows(log, Error) {
		return
	}

	log.error(fmt.Sprintf(format, a...))
}

func NewLogging(name string, show ...logLevel) *Logging {
	return NewLoggingWith(name, os.Stdout, show...)
}

func NewLoggingWith(name string, writer io.Writer, show ...logLevel) *Logging {
	return &Logging{name: name, writer: writer, show: show}
}

func currentTimeAsText() string {
	h, m, s := time.Now().Clock()
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func checkIfLevelShows(log *Logging, lvl logLevel) bool {
	for _, a := range log.Show() {
		if a == lvl {
			return true
		}
	}
	return false
}
