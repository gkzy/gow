/*
report log to nsq server
sam
2020-12-21
*/
package logy

import (
	"fmt"
	"github.com/gkzy/gow/lib/nsq"
	"sync"
	"time"
)

type NsqWriter struct {
	sync.Mutex
	Host  string
	Topic string //nsq topic
	Level int
}

// NewNsqWriter return new NsqWriter
//	host nsq 主机地址
//	topic 一个nsq topic
//	level 记录错误等级
func NewNsqWriter(host, topic string, level int) *NsqWriter {
	nw := &NsqWriter{
		Host:  host,
		Level: level,
		Topic: topic,
	}

	return nw
}

// WriteLog  publish to nsq
func (nw *NsqWriter) WriteLog(t time.Time, level int, b []byte) {
	if level >= nw.Level {
		var err error
		client, err := nsq.NewProducer(nw.Host)
		if err != nil {
			fmt.Println("new nsq Err :: ", err)
		}
		err = client.Publish(nw.Topic, b)
		if err != nil {
			fmt.Println("publish topic error:", err)
		}
	}
}
