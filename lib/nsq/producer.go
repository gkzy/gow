/*
	//init
	 pu,err:=NewProducer("192.168.0.197",4150)
	 if err!=nil{
		//error
	 }
	 b,_:=json.Marshal(obj)
  	 //发送
	 err = pu.Publish("topic",b)
	 if err!=nil{
		//error
	 }

*/
package nsq

import (
	"fmt"
	gnsq "github.com/nsqio/go-nsq"
)

//Producer Producer
type Producer struct {
	P *gnsq.Producer
}

//NewProducer init
func NewProducer(serverAddr string, serverPort int) (producer *Producer, err error) {
	if serverAddr == "" || serverPort == 0 {
		err = fmt.Errorf("[NSQ]init failed：need serverAddr and serverPort")
		return
	}
	config := gnsq.NewConfig()
	p, err := gnsq.NewProducer(fmt.Sprintf("%s:%d", serverAddr, serverPort), config)
	if err != nil {
		return
	}
	p.SetLogger(nil, 0)

	producer = &Producer{
		P: p,
	}
	return
}

//Publish publish
func (m *Producer) Publish(topic string, data []byte) (err error) {
	if m.P == nil {
		err = fmt.Errorf("[NSQ]init failed:%v", err)
	}
	err = m.P.Publish(topic, data)
	defer m.P.Stop()
	if err != nil {
		return fmt.Errorf("[NSQ] publish error:%v", err)
	}
	return
}
