package redis

import (
	"fmt"
	"testing"
	"time"
)

// TestRedis_key
//func TestRedis_key(t *testing.T) {
//
//	//init
//	conf := RDSConfig{
//		Host:      "192.168.0.197",
//		Port:      6379,
//		Password:  "love0021$%",
//		MaxIdle:   50,
//		MaxActive: 10000,
//		DB:        1,
//	}
//
//	err := InitRDSClient(&conf)
//	if err != nil {
//		fmt.Println("连接失败:", err)
//	}
//
//	//get rds common
//	redis := GetRDSCommon()
//	key := "token:1"
//
//	_, err = redis.SetString(key, "abcdef1xd0r1jdkf")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	token, err := redis.GetString(key)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	fmt.Println(token)
//
//}

// User
type User struct {
	Username string
	Sex      int
	Created  int64
}

func TestRedis_Hash(t *testing.T) {
	//init
	conf := RDSConfig{
		Host:      "192.168.0.197",
		Port:      6379,
		Password:  "love0021$%",
		MaxIdle:   50,
		MaxActive: 10000,
		DB:        1,
	}

	err := InitRDSClient(&conf)
	if err != nil {
		fmt.Println("连接失败:", err)
	}

	//get rds common
	redis := GetRDSCommon()
	key := "user:1"

	ok,err:=redis.SetEX(key, "key value",int64((time.Hour * 24).Seconds()))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ok)

}
