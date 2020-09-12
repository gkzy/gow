package logy

import (
	"fmt"
	"io"
	"os"
)

var (
	Std = New(os.Stdout, "", logDefaultLevel())
	defaultDepth = 2
)

func SetOutPut(out io.Writer) {
	Std.SetOutPut(out)
}

// Print
func Print(v ...interface{}) {
	Std.outPut("", Linfo, defaultDepth, fmt.Sprint(v...))
}

// Printf
func Printf(format string, v ...interface{}) {
	Std.outPut("", Linfo, defaultDepth, fmt.Sprintf(format, v))
}

// Println
func Println(v ...interface{}) {
	Std.outPut("", Linfo, defaultDepth, fmt.Sprintln(v...))
}

func Info(v ...interface{}) {
	Std.outPut("", Linfo, defaultDepth, fmt.Sprintln(v...))
}

// Infof
func Infof(format string, v ...interface{}) {
	Std.outPut("", Linfo, defaultDepth, fmt.Sprintf(format, v))
}

// Warn
func Warn(v ...interface{}) {
	Std.outPut("", Lwarn, defaultDepth, fmt.Sprintln(v...))
}

// Warnf
func Warnf(format string, v ...interface{}) {
	Std.outPut("", Lwarn, defaultDepth, fmt.Sprintf(format, v))
}

// Debug
func Debug(v ...interface{}) {
	Std.outPut("", Ldebug, defaultDepth, fmt.Sprintln(v...))
}

// Debugf
func Debugf(format string, v ...interface{}) {
	Std.outPut("", Ldebug, defaultDepth, fmt.Sprintf(format, v))
}

// Error
func Error(v ...interface{}) {
	Std.outPut("", Lerror, defaultDepth, fmt.Sprintln(v...))
}

// Errorf
func Errorf(format string, v ...interface{}) {
	Std.outPut("", Lerror, defaultDepth, fmt.Sprintf(format, v))

}

// Fatal
func Fatal(v ...interface{}) {
	Std.outPut("", Lfatal, defaultDepth, fmt.Sprintln(v...))
	os.Exit(1)
}

// Fatalf
func Fatalf(format string, v ...interface{}) {
	Std.outPut("", Lfatal, defaultDepth, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Panic
func Panic(v ...interface{}) {
	s := fmt.Sprintln(v...)
	Std.outPut("", Lpanic, defaultDepth, s)
	panic(s)
}

// Panicf
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v)
	Std.outPut("", Lpanic, defaultDepth, s)
	panic(s)
}
