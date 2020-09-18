/*
使用方法：

1. init
func init(){
	InitRds(&RDSConfig{
		Host:"127.0.0.1",
		Port:2379,
		Password:"123456",
		MaxIdle:10,
		MaxActive:10000,
		DB:0,
	})
}

2.调用操作方法
rc:=new(RdsCommon)
val,err:=rc.Get(key)
....
...

*/
package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

//RDSConfig redis 配置结构
type RDSConfig struct {
	Host      string
	Port      int
	Password  string
	MaxIdle   int
	MaxActive int
	DB        int
}

var (
	//共用的redis.Pool
	redisClient *redis.Pool
)

//InitRDSClient init config
func InitRDSClient(rdc *RDSConfig) (err error) {
	if rdc == nil {
		err = fmt.Errorf("[RDS] 没有需要init的redis")
		return
	}
	if rdc.Host == "" || rdc.Port == 0 {
		err = fmt.Errorf("[RDS] 没有配置主机或端口")
		return
	}
	redisClient = &redis.Pool{
		MaxIdle:     rdc.MaxIdle,
		MaxActive:   rdc.MaxActive,
		IdleTimeout: 5 * time.Second,
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", fmt.Sprintf("%s:%d", rdc.Host, rdc.Port))
			if conn != nil {
				if len(rdc.Password) != 0 {
					if _, err := conn.Do("AUTH", rdc.Password); err != nil {
						conn.Close()
					}
				}
				if _, err := conn.Do("SELECT", rdc.DB); err != nil {
					conn.Close()
				}
			}
			return
		},
	}

	//ping一次，检测可用性
	rc := redisClient.Get()
	defer rc.Close()
	_, err = rc.Do("PING")
	if err != nil {
		err = fmt.Errorf("[RDS] redis 初始化失败 %v", err)
		return
	}

	return
}
