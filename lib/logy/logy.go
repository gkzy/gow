/*
like https://gitea.com/lunny/log
but made some adjustments
*/
package logy

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	logDate      = 1 << iota         // the date: 2009/0123
	logTime                          // the time: 01:23:23
	microseconds                     // microsecond resolution: 01:23:23.123123.  assumes logTime.
	longFile                         // full file name and line number: /a/b/c/d.go:23
	shortFile                        // final file name element and line number: d.go:23. overrides longFile
	logModule                        // module name
	logLevel                         // logLevel: 0(Debug), 1(Info), 2(Warn), 3(Error), 4(Panic), 5(Fatal)
	longColor                        // color will start [info] end of line
	shortColor                       // color only include [info]
	stdFlags     = logDate | logTime // initial values for the standard logger
)

const (
	Lall = iota
)
const (
	Ldebug = iota
	Linfo
	Lwarn
	Lerror
	Lpanic
	Lfatal
	Lnone
)

const (
	ForeBlack  = iota + 30 //30
	ForeRed                //31
	ForeGreen              //32
	ForeYellow             //33
	ForeBlue               //34
	ForePurple             //35
	ForeCyan               //36
	ForeWhite              //37
)

const (
	BackBlack  = iota + 40 //40
	BackRed                //41
	BackGreen              //42
	BackYellow             //43
	BackBlue               //44
	BackPurple             //45
	BackCyan               //46
	BackWhite              //47
)

var levels = []string{
	"[D]",
	"[I]",
	"[W]",
	"[E]",
	"[P]",
	"[F]",
}

var colors = []int{
	ForeCyan,
	ForeGreen,
	ForeYellow,
	ForeRed,
	ForePurple,
	ForeBlue,
}

// logDefaultLevel
func logDefaultLevel() int {
	if runtime.GOOS == "windows" {
		return logLevel | stdFlags | shortFile
	}
	return logLevel | stdFlags | shortFile | longColor
}

// SetLevels MUST called before all logs
func SetLevels(lvs []string) {
	levels = lvs
}

// MUST called before all logs
func SetColors(cls []int) {
	colors = cls
}

// logData
type logData struct {
	App     string `json:"app"`     //应用
	Ver     string `json:"ver"`     //应用版本
	Msg     string `json:"msg"`     //日志内容
	Level   int    `json:"level"`   //日志等级
	Created int    `json:"created"` //当前时间
}

// Logger
type Logger struct {
	mu         sync.Mutex
	prefix     string
	flag       int
	Level      int
	out        io.Writer
	buf        bytes.Buffer
	levelStats [6]int64
	loc        *time.Location
}

// New return *Logger
func New(out io.Writer, prefix string, flag int) *Logger {
	log := &Logger{
		out:    out,
		prefix: prefix,
		Level:  0,
		flag:   flag,
		loc:    time.Local,
	}
	if out != os.Stdout {
		log.flag = rmColorFlags(flag)
	}
	return log
}

func (l *Logger) SetOutPut(out io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = out
	if out != os.Stdout {
		l.flag = rmColorFlags(l.flag)
	}
}

// ----------------Print-----------------

// Printf
func (l *Logger) Printf(format string, v ...interface{}) {
	l.outPut("", Linfo, 2, fmt.Sprintf(format, v...))
}

// Print
func (l *Logger) Print(v ...interface{}) {
	l.outPut("", Linfo, 2, fmt.Sprint(v...))
}

// Println
func (l *Logger) Println(v ...interface{}) {
	l.outPut("", Linfo, 2, fmt.Sprintln(v...))
}

// ----------------Debug-----------------

// Debugf
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.outPut("", Ldebug, 2, fmt.Sprintf(format, v...))
}

// Debug
func (l *Logger) Debug(v ...interface{}) {
	l.outPut("", Ldebug, 2, fmt.Sprintln(v...))
}

// ----------------Info-----------------

// Infof
func (l *Logger) Infof(format string, v ...interface{}) {
	l.outPut("", Linfo, 2, fmt.Sprintf(format, v...))
}

// Info
func (l *Logger) Info(v ...interface{}) {
	l.outPut("", Linfo, 2, fmt.Sprintln(v...))
}

// ----------------Warn-----------------

// Warnf
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.outPut("", Lwarn, 2, fmt.Sprintf(format, v...))
}

// Warn
func (l *Logger) Warn(v ...interface{}) {
	l.outPut("", Lwarn, 2, fmt.Sprintln(v...))
}

// ----------------Error-----------------

// Errorf
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.outPut("", Lerror, 2, fmt.Sprintf(format, v...))
}

// Error
func (l *Logger) Error(v ...interface{}) {
	l.outPut("", Lerror, 2, fmt.Sprintln(v...))
}

// ----------------panic-----------------

// Panicf
func (l *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.outPut("", Lpanic, 2, s)
	panic(s)
}

// Panic
func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.outPut("", Lpanic, 2, s)
	panic(s)
}

//========================private func=====================

//formatHeader formatHeader
func (l *Logger) formatHeader(buf *bytes.Buffer, t time.Time, file string, line int, lvl int, reqId string) {
	if l.prefix != "" {
		buf.WriteString(l.prefix)
	}
	if l.flag&(logDate|logTime|microseconds) != 0 {
		if l.flag&logDate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			buf.WriteByte('/')
			itoa(buf, int(month), 2)
			buf.WriteByte('/')
			itoa(buf, day, 2)
			buf.WriteByte(' ')
		}
		if l.flag&(logTime|microseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			buf.WriteByte(':')
			itoa(buf, min, 2)
			buf.WriteByte(':')
			itoa(buf, sec, 2)
			if l.flag&microseconds != 0 {
				buf.WriteByte('.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			buf.WriteByte(' ')
		}
	}
	if reqId != "" {
		buf.WriteByte('[')
		buf.WriteString(reqId)
		buf.WriteByte(']')
		buf.WriteByte(' ')
	}

	if l.flag&(shortColor|longColor) != 0 {
		buf.WriteString(fmt.Sprintf("\033[1;%dm", colors[lvl]))
	}
	if l.flag&logLevel != 0 {
		buf.WriteString(levels[lvl])
		buf.WriteByte(' ')
	}
	if l.flag&shortColor != 0 {
		buf.WriteString("\033[0m")
	}

	if l.flag&logModule != 0 {
		buf.WriteByte('[')
		buf.WriteString(moduleOf(file))
		buf.WriteByte(']')
		buf.WriteByte(' ')
	}
	if l.flag&(shortFile|longFile) != 0 {
		if l.flag&shortFile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		buf.WriteByte('[')
		buf.WriteString(file)
		buf.WriteByte(':')
		itoa(buf, line, -1)
		buf.WriteByte(']')
		buf.WriteByte(' ')
	}
}

//outPut outPut
func (l *Logger) outPut(reqId string, lvl int, callDepth int, s string) error {
	if lvl < l.Level {
		return nil
	}
	now := time.Now().In(l.loc) // get this early.
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(shortFile|longFile|logModule) != 0 {
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(callDepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.levelStats[lvl]++
	l.buf.Reset()
	l.formatHeader(&l.buf, now, file, line, lvl, reqId)
	l.buf.WriteString(s)
	if l.flag&longColor != 0 {
		l.buf.WriteString("\033[0m")
	}
	if len(s) > 0 && s[len(s)-1] != '\n' {
		l.buf.WriteByte('\n')
	}
	_, err := l.out.Write(l.buf.Bytes())
	return err
}

func rmColorFlags(flag int) int {
	// for un std out, it should not show color since almost them don't support
	if flag&longColor != 0 {
		flag = flag ^ longColor
	}
	if flag&shortColor != 0 {
		flag = flag ^ shortColor
	}
	return flag
}

func moduleOf(file string) string {
	pos := strings.LastIndex(file, "/")
	if pos != -1 {
		pos1 := strings.LastIndex(file[:pos], "/src/")
		if pos1 != -1 {
			return file[pos1+5 : pos]
		}
	}
	return "UNKNOWN"
}

func itoa(buf *bytes.Buffer, i int, wid int) {
	var u uint = uint(i)
	if u == 0 && wid <= 1 {
		buf.WriteByte('0')
		return
	}

	// Assemble decimal in reverse order.
	var b [32]byte
	bp := len(b)
	for ; u > 0 || wid > 0; u /= 10 {
		bp--
		wid--
		b[bp] = byte(u%10) + '0'
	}

	// avoid slicing b to avoid an allocation.
	for bp < len(b) {
		buf.WriteByte(b[bp])
		bp++
	}
}
