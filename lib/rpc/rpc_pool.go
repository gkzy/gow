/*
一个grpc连接池复用的实现
sam

*/

package rpc

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// PoolModeStrict 在实际创建连接数达上限后，池子中没有连接时不会新建连接
	PoolModeStrict = iota

	// PoolModeLoose 在实际创建连接数达上限后，池子中没有连接时会新建连接
	PoolModeLoose
)

var (
	// ErrorOption
	ErrorOption = "option error"
	// ErrorPoolInit 连接p池初始化出错
	ErrorPoolInit = "pool init error"
	//ErrorGetTimeout 获取连接超时
	ErrorGetTimeout = "getting connection client timeout from pool"
	//ErrorDialConn 创建连接时发生错误
	ErrorDialConn = "dialing connection error"

	// ErrorPoolIsClosed 连接池已关闭
	ErrorPoolIsClosed = "pool is closed"
)

// Option pool param option
type Option struct {
	// func() (*grpc.ClientConn, error)
	Factor Factor

	//Init init 连接数
	Init int32

	//Cap 连接上限
	Cap int32

	//IdleDur
	IdleDur time.Duration

	// MaxLifeDur
	MaxLifeDur time.Duration

	// Timeout pool关闭时
	Timeout time.Duration

	// 模式
	Mode int
}

// Factor factor func
type Factor func() (*grpc.ClientConn, error)

// Pool 连接池
type Pool struct {
	clients    chan *Client
	connCnt    int32
	cap        int32
	idleDur    time.Duration
	maxLifeDur time.Duration
	timeout    time.Duration
	factor     Factor
	lock       sync.RWMutex
	mode       int
}

// Client grpc client
type Client struct {
	*grpc.ClientConn
	timeUsed time.Time
	timeInit time.Time
	pool     *Pool
}

// DefaultPool return a default pool
func DefaultPool(factor Factor, init, cap int32) (*Pool, error) {
	option := &Option{
		Factor:     factor,
		Init:       init,
		Cap:        cap,
		IdleDur:    10 * time.Second,
		MaxLifeDur: 60 * time.Second,
		Timeout:    10 * time.Second,
		Mode:       PoolModeLoose,
	}

	return NewPool(option)
}

// NewPool return a new pool
//	need option
func NewPool(option *Option) (*Pool, error) {
	if option == nil {
		return nil, errors.New(ErrorOption)
	}
	if option.Factor == nil {
		return nil, errors.New(ErrorPoolInit)
	}

	if option.Init < 1 || option.Cap < 1 || option.IdleDur < 1 || option.MaxLifeDur < 1 {
		return nil, errors.New(ErrorPoolInit)
	}

	if option.Init > option.Cap {
		option.Init = option.Cap
	}

	pool := &Pool{
		clients:    make(chan *Client, option.Cap),
		cap:        option.Cap,
		idleDur:    option.IdleDur,
		maxLifeDur: option.MaxLifeDur,
		timeout:    option.Timeout,
		factor:     option.Factor,
		mode:       option.Mode,
	}

	for i := int32(0); i < option.Init; i++ {
		client, err := pool.createClient()
		if err != nil {
			return nil, errors.New(ErrorPoolInit)
		}
		pool.clients <- client
	}

	return pool, nil
}

// Get 从连接池取出一个链接
func (pool *Pool) Get(ctx context.Context) (*Client, error) {
	if pool.IsClose() {
		return nil, errors.New(ErrorPoolIsClosed)
	}
	var (
		client *Client
		err    error
		now    = time.Now()
	)
	select {
	case <-ctx.Done():
		if pool.mode == PoolModeStrict {
			pool.lock.Lock()
			defer pool.lock.Unlock()
			if pool.connCnt >= pool.cap {
				err = errors.New(ErrorGetTimeout)
			} else {
				client, err = pool.createClient()
			}
			return client, err
		}
	case client = <-pool.clients:
		if client != nil && pool.idleDur > 0 && client.timeUsed.Add(pool.idleDur).After(now) {
			client.timeUsed = now
			return client, nil
		}
	}

	// 如果连接已经是idle连接，或者是非严格模式下没有获取连接
	// 则新建一个连接同时销毁原有idle连接
	if client != nil {
		client.Destroy()
	}
	client, err = pool.createClient()
	if err != nil {
		return nil, err
	}
	return client, nil

}

func (pool *Pool) createClient() (*Client, error) {
	conn, err := pool.factor()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", ErrorPoolInit, err)
	}
	now := time.Now()
	client := &Client{
		ClientConn: conn,
		timeUsed:   now,
		timeInit:   now,
		pool:       pool,
	}

	atomic.AddInt32(&pool.connCnt, 1)
	return client, nil
}

// Size return pool client size
func (pool *Pool) Size() int {
	pool.lock.RLock()
	defer pool.lock.RUnlock()
	return len(pool.clients)
}

// ConnCnt 实际连接数
func (pool *Pool) ConnCnt() int32 {
	return pool.connCnt
}

func (pool *Pool) IsClose() bool {
	return pool == nil || pool.clients == nil
}

func (pool *Pool) Close() {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	if pool.IsClose() {
		return
	}
	clients := pool.clients
	pool.clients = nil

	go func() {
		for {
			select {
			case client := <-clients:
				if client != nil {
					client.Destroy()
				}
			case <-time.Tick(pool.timeout):
				if len(clients) <= 0 {
					close(clients)
					break
				}
			}
		}

	}()
}

// Close close conn
func (client *Client) Close() {
	go func() {
		pool := client.pool
		now := time.Now()
		if pool.IsClose() {
			client.Destroy()
			return
		}
		if pool.maxLifeDur > 0 && client.timeInit.Add(pool.maxLifeDur).Before(now) {
			client.Destroy()
			return
		}
		if client.ClientConn == nil {
			return
		}

		client.timeUsed = now
		client.pool.clients <- client
	}()
}

// Destroy destroy conn
func (client *Client) Destroy() {
	if client.ClientConn != nil {
		client.ClientConn.Close()
		if client.pool != nil {
			atomic.AddInt32(&client.pool.connCnt, -1)
		}
	}
	client.ClientConn = nil
	client.pool = nil
}

// TimeInit 连的创建时间
func (client *Client) TimeInit() time.Time {
	return client.timeInit
}

// TimeUsed 连接上一次的使用时间
func (client *Client) TimeUsed() time.Time {
	return client.timeUsed
}
