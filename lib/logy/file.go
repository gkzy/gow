/*
logy 写日志文件的实现
sam
2020-09-29
*/
package logy

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileWriter 文件
type FileWriter struct {
	sync.Mutex
	prefix string
	dir    string
	date   string
	file   *os.File
}

// WriteLog 接口实现
func (fw *FileWriter) WriteLog(t time.Time, level int, b []byte) {
	fw.Lock()
	fw.writeFile(t)
	fw.file.Write(b)
	fw.Unlock()
}

func (fw *FileWriter) writeFile(t time.Time) {
	newDate := t.Format("20060102")
	if fw.date != newDate && fw.file != nil {
		fw.file.Close()
		fw.file = nil
	}
	if fw.file == nil {
		//目录
		dir := filepath.Dir(fw.dir)
		err := os.MkdirAll(dir, 755)
		if err != nil {
			panic(err)
		}

		//web-20200918.log
		fileName := fmt.Sprintf("%s-%s.log", fw.prefix, newDate)
		file, err := os.OpenFile(filepath.Join(fw.dir, fileName), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			panic(err)
		}
		fw.file = file
		fw.date = newDate
	}
}

// NewFileWriter return a new FileWriter
func NewFileWriter(prefix, dir string) *FileWriter {
	return &FileWriter{
		prefix: prefix,
		dir:    dir,
	}
}
