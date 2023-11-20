package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kr/pretty"
	"golib/jk"
)

type logImpl interface {
	Status(message any, v ...any)

	Trace(message any, v ...any)
	Debug(message any, v ...any)
	Info(message any, v ...any)
	Warn(message any, v ...any)
	Error(message any, v ...any)

	AddError(lines []string, message any, v ...any) []string
}

type LogLevel int

const (
	LOff LogLevel = iota
	LError
	LWarn
	LInfo
	LDebug
	LTrace
)

var level = LTrace

const statusMarker = "<!STATUS>"

func SetLevel(logLevel LogLevel) {
	level = logLevel
}

func SetLevelS(logLevel string) {
	SetLevel(LevelFromString(logLevel))
}

func LevelFromString(logLevel string) LogLevel {
	switch strings.ToLower(logLevel) {
	case "off":
		return LOff
	case "error":
		return LError
	case "warn":
		return LWarn
	case "info":
		return LInfo
	case "debug":
		return LDebug
	case "trace":
		return LTrace
	}
	panic("Unexpected log level: [" + logLevel + "]")
}

var impl logImpl = newGolog(os.Stdout)

func UseGoLog(output io.Writer) {
	impl = newGolog(output)
}

func UseAntigloss(folder string, prefix string) {
	impl = newAntigloss(folder, prefix)
}

func Status(message any, v ...any) {
	impl.Status(message, v...)
}

func Json(message any) {
	j, _ := jk.Print(message)
	impl.Status(string(j))
}

func Trace(message any, v ...any) {
	if level >= LTrace {
		impl.Trace(message, v...)
	}
}

func Debug(message any, v ...any) {
	if level >= LDebug {
		impl.Debug(message, v...)
	}
}

func Info(message any, v ...any) {
	if level >= LInfo {
		impl.Info(message, v...)
	}
}

func Warn(message any, v ...any) {
	if level >= LWarn {
		impl.Warn(message, v...)
	}
}

func Error(message any, v ...any) {
	if level >= LError {
		impl.Error(message, v...)
	}
}

func AddError(lines []string, message any, v ...any) []string {
	return impl.AddError(lines, message, v...)
}

// this si repeated from main/lang to avoid circular dependency
func print(s any) string {
	switch s.(type) {
	case string:
		return s.(string)
	default:
		return fmt.Sprint(s)
	}
}

type prettyPrinter struct {
	value any
}

func (self *prettyPrinter) String() string {
	return pretty.Sprint(self.value)
}

func Pretty(v any) *prettyPrinter {
	return &prettyPrinter{v}
}
