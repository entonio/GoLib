package log

import (
	"fmt"

	"github.com/antigloss/go/logger"
)

type antigloss struct{}

func newAntigloss(folder string, prefix string) *antigloss {
	logger.Init(&logger.Config{
		LogDir:            folder,             // specify the directory to save the logfiles
		LogDest:           logger.LogDestFile, // specify the directory to save the logfiles
		LogFileMaxNum:     10000,              // maximum logfiles allowed under the specified log directory
		LogFileNumToDel:   10,                 // number of logfiles to delete when number of logfiles exceeds the configured limit
		LogFileMaxSize:    5,                  // maximum size of a logfile in MB
		LogLevel:          logger.LogLevelTrace,
		LogFilenamePrefix: prefix,
	})
	return &antigloss{}
}

func (self *antigloss) Status(message any, v ...any) {
	logger.Infof(statusMarker+" "+print(message), v...)
}

func (self *antigloss) Trace(message any, v ...any) {
	logger.Tracef(print(message), v...)
}

func (self *antigloss) Debug(message any, v ...any) {
	logger.Tracef(print(message), v...)
}

func (self *antigloss) Info(message any, v ...any) {
	logger.Infof(print(message), v...)
}

func (self *antigloss) Warn(message any, v ...any) {
	logger.Warnf(print(message), v...)
}

func (self *antigloss) Error(message any, v ...any) {
	logger.Errorf(print(message), v...)
}

func (self *antigloss) AddError(lines []string, message any, v ...any) []string {
	format := print(message)
	logger.Errorf(format, v...)
	return append(lines, fmt.Sprintf(format, v...))
}
