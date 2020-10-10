package logy

import "time"

type HTTPWriter struct {
	url    string
	method string
}

// WriteLog write log to http api
//	TODO:
func (aw *HTTPWriter) WriteLog(t time.Time, level int, b []byte) {

}

// NewHTTPWriter return new HTTPWriter
func NewHTTPWriter(url, method string) *HTTPWriter {
	return &HTTPWriter{
		url:    url,
		method: method,
	}
}
