package logy

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type ByType int

// GetFileFormat 获取格式
func (b ByType) GetFileFormat() string {
	return formats[b]
}

// SetFileFormat 设置文件格式
func SetFileFormat(t ByType, format string) {
	formats[t] = format
}

const (
	Day ByType = iota
	Hour
	Month
)

var (
	// formats 文件格式map
	formats = map[ByType]string{
		Day:   "2006-01-02",
		Hour:  "2006-01-02-15",
		Month: "2006-01",
	}
	// defaultMaxDay  日志文件默认的留存天数
	defaultMaxDay = 30
)

// FileOptions 写文件选项
type FileOptions struct {
	Dir    string         //日志目录
	ByType ByType         //按天/小时/月?
	Loc    *time.Location //时间
	MaxDay int            //最大留存天数
}

type Files struct {
	FileOptions
	file       *os.File
	lastFormat string
	mu         sync.Mutex
}

// NewFileWriter 设置写文件参数
// 	w:=logy.NewFileWriter(logy.FileOptions{
//		ByType:log.Day,
//		Dir:"./logs",
//		MaxDay:6,
//	})
//  logy.Std.SetOutPut(w)
func NewFileWriter(opts ...FileOptions) *Files {
	opt := prepareFileOption(opts)
	file := &Files{
		FileOptions: opt,
	}
	// init clear log
	file.clearLog()
	// timer clear
	go file.startTimer()
	return file
}

func (f *Files) getFile() (*os.File, error) {
	var err error
	t := time.Now().In(f.Loc)
	if f.file == nil {
		f.lastFormat = t.Format(f.ByType.GetFileFormat())
		f.file, err = os.OpenFile(filepath.Join(f.Dir, f.lastFormat+".log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		return f.file, err
	}
	if f.lastFormat != t.Format(f.ByType.GetFileFormat()) {
		f.file.Close()
		f.lastFormat = t.Format(f.ByType.GetFileFormat())
		f.file, err = os.OpenFile(filepath.Join(f.Dir, f.lastFormat+".log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		return f.file, err
	}
	return f.file, nil
}

func (f *Files) Write(bs []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	w, err := f.getFile()
	if err != nil {
		return 0, err
	}
	return w.Write(bs)
}

func (f *Files) Close() {
	if f.file != nil {
		f.file.Close()
		f.file = nil
	}
	f.lastFormat = ""
}

// prepareFileOption 预处理文件选项
func prepareFileOption(opts []FileOptions) FileOptions {
	var opt FileOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt.Dir == "" {
		opt.Dir = "./"
	}
	err := os.MkdirAll(opt.Dir, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	if opt.MaxDay == 0 {
		opt.MaxDay = defaultMaxDay
	}
	if opt.Loc == nil {
		opt.Loc = time.Local
	}
	return opt
}

//============private===========

// clearLog 清理掉过期的日志文件
// TODO:
func (f *Files) clearLog() {
	files := getDirFiles(f.Dir)
	now := time.Now()
	for _, item := range files {
		modTime := item.ModTime
		if modTime.Add(time.Hour * 24 * time.Duration(f.MaxDay-1)).Before(now) {
			os.Remove(item.Name)
		}
	}
}

// startTimer start time ticker
func (f *Files) startTimer() {
	now := time.Now()
	next := now.Add(time.Second * 3600)
	second := time.Duration(next.Sub(now).Seconds())
	f.timer(second)
}

// timer
func (f *Files) timer(seconds time.Duration) {
	timer := time.NewTicker(seconds * time.Second)
	for {
		select {
		case <-timer.C:
			{
				f.clearLog()
				nextTimer := time.NewTicker(3600 * time.Second)
				for {
					select {
					case <-nextTimer.C:
						{
							f.startTimer()
							return
						}
					}
				}
			}
		}
	}
}

// FileInfo
type FileInfo struct {
	Name    string
	ModTime time.Time
	Size    int64
}

// getDirFiles
func getDirFiles(path string) (files []*FileInfo) {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	files = make([]*FileInfo, 0)
	for _, fi := range dir {
		if !fi.IsDir() {
			files = append(files, &FileInfo{
				Name:    fi.Name(),
				ModTime: fi.ModTime(),
				Size:    fi.Size(),
			})
		}
	}

	return
}
