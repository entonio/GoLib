package lang

import (
	"github.com/go-errors/errors"
	"runtime"
	"strings"
)

var stackPrefix string
var stackInfix string
var stackSuffix string

func SetStackStyle(style int) {
	switch style {
	case 1:
		stackPrefix = ""
		stackInfix = "()\n\t"
		stackSuffix = ""
	case 2:
		stackPrefix = "\tat "
		stackInfix = "("
		stackSuffix = ")"
	}
}

type Frame struct {
	Function string
	Filename string
	Line     int
}

func CurrentStack(skipTop int, skipBottom int) string {
	stack := runtimeStack(false)
	frames := framesFromStack(stack, skipTop+2, skipBottom)
	return FormatFrames(frames)
}

func runtimeStack(all bool) string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return string(buf[:n])
		}
		buf = make([]byte, 2*len(buf))
	}
}

func errorStack(err error, skip int) string {
	// the go-errrors package sometimes produces inaccurate stacks
	return errors.Wrap(err, skip+1).ErrorStack()
}

func framesFromStack(stack string, skipTop int, skipBottom int) (frames []Frame) {
	lines := strings.Split(stack, "\n")
	limit := len(lines)/2 - skipBottom
	f := 0
	for i, line1 := range lines {
		// thread name
		if strings.HasSuffix(line1, ":") {
			continue
		}
		// not line1
		if !strings.HasSuffix(line1, ")") {
			continue
		}
		line2 := lines[i+1]
		f += 1
		if f < skipTop {
			continue
		}
		if f >= limit {
			break
		}
		var frame Frame
		s1 := strings.Split(line1, "(")
		s1a := strings.Split(s1[0], "/")
		frame.Function = s1a[len(s1a)-1]
		if strings.HasSuffix(frame.Function, ".") {
			frame.Function += "<" + strings.ReplaceAll(s1[1], ")", ">")
		}
		if len(frames) == 0 && frame.Function == "lang.Assert" {
			continue
		}
		s2 := strings.Split(line2, ":")
		s2a := strings.Split(s2[0], "/")
		frame.Filename = s2a[len(s2a)-1]
		if len(s2) > 1 {
			s2b := strings.Split(s2[1], " ")
			frame.Line = AsInt(s2b[0])
		}
		frames = append(frames, frame)
	}
	return
	/*
		frames := errors.Wrap(err, 0).StackFrames()
		limit := len(frames) - skipBottom
		for i, frame := range frames {
			if i < skipTop || i >= limit {
				continue
			}
			split := strings.Split(frame.Func().Name(), "/")
			name := split[len(split)-1]
			file := frame.File
			if mainPackagePrefix != nil {
				file = mainPackagePrefix.ReplaceAllString(file, "")
			}
			//fmt.Printf("%v %s -> %s\n", mainPackagePrefix != nil, frame.File, file)
			//file = strings.SplitAfterN(file, "/src/", 1)[0]
			if len(stack) > 0 {
				stack += "\n"
			}
			stack += stackPrefix + name + stackInfix + file + ":" + Print(frame.LineNumber) + stackSuffix
		}
		return
	*/
}

func FormatFrames(frames []Frame) string {
	var stack string
	for _, frame := range frames {
		if len(stack) > 0 {
			stack += "\n"
		}
		stack += stackPrefix + frame.Function + stackInfix + frame.Filename + ":" + Print(frame.Line) + stackSuffix
	}
	return stack
}
