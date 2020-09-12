package gow

import (
	"encoding/json"
	"time"
)

// DataResponse data json response struct
type DataResponse struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	STime string `json:"stime"`
	Time  int    `json:"time"`
	Body  *Body  `json:"body"`
}

type Body struct {
	Pager *Pager      `json:"pager"`
	Data  interface{} `json:"data"`
}

// Pager
type Pager struct {
	Page      int64 `json:"page"`
	Limit     int64 `json:"limit"`
	Offset    int64 `json:"offset"`
	Count     int64 `json:"count"`
	PageCount int64 `json:"page_count"`
}

// DataPager middlewares
//	实现分页参数的处理
func DataPager() HandlerFunc {
	return func(c *Context) {
		pager := new(Pager)
		pager.Page, _ = c.GetInt64("page", 1)
		if pager.Page < 1 {
			pager.Page = 1
		}
		pager.Limit, _ = c.GetInt64("limit", 10)
		if pager.Limit < 1 {
			pager.Limit = 1
		}

		pager.Offset = (pager.Page - 1) * pager.Limit
		c.Pager = pager
		c.Next()
	}
}

// DataJSON DataJSON json data
//	response format json
func (c *Context) DataJSON(args ...interface{}) {
	var (
		err   error
		pager *Pager
		data  interface{}
		msg   string
		code  int
	)
	for _, v := range args {
		switch vv := v.(type) {
		case int:
			code = vv
		case string:
			msg = vv
		case error:
			err = vv
		case *Pager:
			pager = vv
		default:
			data = vv
		}
	}
	if err != nil {
		debugPrint(c.Request.URL.String(), err)
	}
	if code == 0 && msg == "" {
		msg = "success"
	}

	body := new(Body)

	if pager != nil {
		pager.PageCount = getPageCount(pager.Count, pager.Limit)
	} else {
		pager = &Pager{}
	}
	body.Pager = pager
	body.Data = data

	resp := &DataResponse{
		Code:  code,
		Msg:   msg,
		STime: time.Now().Format("2006-01-02 15:04:05"),
		Time:  int(time.Now().Unix()),
		Body:  body,
	}
	c.JSON(&resp)
	return
}

// DecodeJSONBody request body to struct or map
func (c *Context) DecodeJSONBody(v interface{}) error {
	body := c.Body()
	return json.Unmarshal(body, &v)
}

// getPageCount return pagerCount
func getPageCount(count, limit int64) (pageCount int64) {
	if count > 0 && limit > 0 {
		if count%limit == 0 {
			pageCount = count / limit
		} else {
			pageCount = (count / limit) + 1
		}
	}
	return pageCount
}
