package logy

import (
	"github.com/imroc/req"
	"sync"
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
	level  int
	sync.Mutex
}

// WriteLog write log to http api
//
func (aw *HTTPWriter) WriteLog(t time.Time, level int, b []byte) {

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
func NewHTTPWriter(url, method, token string, level int) *HTTPWriter {
	return &HTTPWriter{
		url:    url,
		method: method,
		token:  token,
		level:  level,
	}
}
