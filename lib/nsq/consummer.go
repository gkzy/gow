/*
使用方法：


var ch = make(chan []byte)

func init(){
	mh:=NewMessageHandler("192.168.0.197",4161)
	go mh.Registry("SMSTopic",ch)
	go dealSMS()
}

func dealSMS(){
	for {
		select {
		case s := <-ch:
			consumerSMS(s)
		}
	}
}


func consumerSMS(b []byte){
	//具体业务实现
	...

}

*/

package nsq

import (
	"fmt"
	gnsq "github.com/nsqio/go-nsq"
)

//MessageHandler MessageHandler
type MessageHandler struct {
	msgChan      chan *gnsq.Message
	stop         bool
	consumerAddr string
	consumerPort int
	Channel      string
}

//NewMessageHandler
func NewMessageHandler(consumerAddr string, consumerPort int, channel string) (mh *MessageHandler, err error) {
	if consumerAddr == "" || consumerPort == 0 {
		err = fmt.Errorf("[NSQ] need consumerAddr and consumerPort")
		return
	}
	mh = &MessageHandler{
		msgChan:      make(chan *gnsq.Message, 1024),
		stop:         false,
		consumerAddr: consumerAddr,
		consumerPort: consumerPort,
		Channel:      channel,
	}

	return
}

//Registry Registry
func (m *MessageHandler) Registry(topic string, ch chan []byte) {
	config := gnsq.NewConfig()
	consumer, err := gnsq.NewConsumer(topic, m.Channel, config)
	if err != nil {
		panic(err)
	}
	consumer.SetLogger(nil, 0)
	consumer.AddHandler(gnsq.HandlerFunc(m.handlerMessage))
	err = consumer.ConnectToNSQLookupd(fmt.Sprintf("%s:%d", m.consumerAddr, m.consumerPort))
	if err != nil {
		panic(err)
	}
	m.process(ch)

}

//process process
func (m *MessageHandler) process(ch chan<- []byte) {
	m.stop = false
	for {
		select {
		case message := <-m.msgChan:
			ch <- message.Body
			if m.stop {
				close(m.msgChan)
				return
			}
		}
	}
}

//handlerMessage handlerMessage
func (m *MessageHandler) handlerMessage(message *gnsq.Message) error {
	if !m.stop {
		m.msgChan <- message
	}
	return nil
}
