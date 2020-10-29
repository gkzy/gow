package logy

import (
	"github.com/imroc/req"
	"time"
)

var (
	timeOut = time.Duration(3 * time.Second)
)

// HTTPWriter http api writer
type HTTPWriter struct {
	url    string
	method string
	token  string
}

// WriteLog write log to http api
//
func (aw *HTTPWriter) WriteLog(t time.Time, level int, b []byte) {
	// 只上报错误及以上日志
	if level >= LevelError {
		resp, err := aw.httpRequest(aw.url, aw.method, aw.token, b)
		if err != nil {
			Errorf("[http writer] error:%v", err)
		}
		Debug(resp)
	}
}

// httpRequest http request
func (aw *HTTPWriter) httpRequest(url, method, token string, body []byte) (resp *req.Resp, err error) {
	r := req.New()
	r.SetTimeout(timeOut)
	header := req.Header{
		"token": token,
	}
	resp, err = r.Do(method, url, header, body)
	if err != nil {
		return
	}
	return
}

// NewHTTPWriter return new HTTPWriter
func NewHTTPWriter(url, method, token string) *HTTPWriter {
	return &HTTPWriter{
		url:    url,
		method: method,
		token:  token,
	}
}
