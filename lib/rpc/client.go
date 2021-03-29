package rpc

import (
	"fmt"
	"google.golang.org/grpc"
)

//NewClient 返回rpc客户端
//serverAddr:服务端地址
//serverPort:服务端Port
func NewClient(serverAddr string, serverPort int) (client *grpc.ClientConn, err error) {
	server := fmt.Sprintf("%s:%d", serverAddr, serverPort)
	client, err = grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("[RPC] get client  error: %v", err))
		return
	}
	return
}
