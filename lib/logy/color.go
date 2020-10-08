/*
color 实现
sam
2020-09-29
*/
package logy

import (
	"sync"
	"time"
)

const colorEnd = "\033[0m"

// 加粗的颜色输出
var logColor = []string{
	Ldebug:  "\033[1;36m",
	Linfo:   "\033[1;37m",
	Lnotice: "\033[1;33m",
	Lwarn:   "\033[1;35m",
	Lerror:  "\033[1;31m",
	Lpanic:  "\033[1;31m",
	Lfatal:  "\033[1;31m",
}

type colorWriter struct {
	m sync.Mutex
	b []byte
	w Writer
}

func (cw *colorWriter) WriteLog(t time.Time, level int, b []byte) {
	cw.m.Lock()
	cw.b = cw.b[:0]
	cw.b = append(cw.b, logColor[level]...)
	cw.b = append(cw.b, b...)
	cw.b = append(cw.b, colorEnd...)
	cw.w.WriteLog(t, level, cw.b)
	cw.m.Unlock()
}

// WithColor 指定某一个实现的writer使用颜色
func WithColor(w Writer) Writer {
	return &colorWriter{w: w}
}
