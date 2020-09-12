/*
http client
使用github.com/imroc/req库
返回string，如果需要到struct，需要自己反序列化

*/
package util

import (
	"fmt"
	"github.com/imroc/req"
	"net/http"
	"time"
)

const (
	timeOut   = 10 //10s
	userAgent = "golang-http-client/1.1"
)

//HttpGet http get
func HttpGet(url string) (ret string, err error) {
	if url == "" {
		err = fmt.Errorf("url为空")
		return
	}
	header := make(http.Header)
	header.Set("User-Agent", userAgent)
	req.SetTimeout(timeOut * time.Second)

	resp, err := req.Get(url, header)
	if err != nil {
		return
	}
	ret,err = resp.ToString()
	if err != nil {
		return
	}

	return
}

//HttpPost http post
func HttpPost(url string, param req.Param) (ret string, err error) {
	if url == "" {
		err = fmt.Errorf("url为空")
		return
	}
	header := make(http.Header)
	header.Set("User-Agent", userAgent)
	req.SetTimeout(timeOut * time.Second)

	resp, err := req.Post(url, param, header)
	if err != nil {
		return
	}
	ret,err = resp.ToString()
	if err != nil {
		return
	}

	return
}
