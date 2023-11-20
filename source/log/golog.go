package log

import (
	"fmt"
	"io"
	logger "log"
)

type golog struct {
	logger *logger.Logger
}

func newGolog(output io.Writer) *golog {
	return &golog{
		logger: logger.New(output, "", logger.Ldate|logger.Ltime|logger.Lshortfile),
	}
}

func (self *golog) Status(message any, v ...any) {
	self.printLine(statusMarker, message, v...)
}

func (self *golog) Trace(message any, v ...any) {
	self.printLine("TRACE", message, v...)
}

func (self *golog) Debug(message any, v ...any) {
	self.printLine("DEBUG", message, v...)
}

func (self *golog) Info(message any, v ...any) {
	self.printLine("INFO", message, v...)
}

func (self *golog) Warn(message any, v ...any) {
	self.printLine("WARN", message, v...)
}

func (self *golog) Error(message any, v ...any) {
	self.printLine("ERROR", message, v...)
}

func (self *golog) AddError(lines []string, message any, v ...any) []string {
	return self.addLine(lines, "ERROR", message, v...)
}

func (self *golog) printLine(level string, message any, v ...any) {
	line := self.formatted(message, v...)
	self.logger.Output(4, level+" "+line)
}

func (self *golog) addLine(lines []string, level string, message any, v ...any) []string {
	line := self.formatted(message, v...)
	self.logger.Output(4, level+" "+line)
	return append(lines, line)
}

func (self *golog) formatted(message any, v ...any) string {
	if len(v) == 0 {
		return print(message)
	} else {
		return fmt.Sprintf(print(message), v...)
	}
}
