package redis

import (
	"fmt"
	"testing"
	"time"
)

// TestRedis_key
func TestRedis_key(t *testing.T) {

	//init
	conf := RDSConfig{
		Host:      "192.168.0.197",
		Port:      6379,
		Password:  "123456",
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
	key := "token:1"

	_, err = redis.SetString(key, "abcdef1xd0r1jdkf")
	if err != nil {
		t.Fatal(err)
	}

	token, err := redis.GetString(key)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(token)

}

// User
type User struct {
	Username string
	Sex      int
	Created  int
}

func TestRedis_Hash(t *testing.T) {

	//init
	conf := RDSConfig{
		Host:      "192.168.0.197",
		Port:      6379,
		Password:  "123456",
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

	user := &User{
		Username: "user001",
		Sex:      1,
		Created:  int(time.Now().Unix()),
	}

	err = redis.SetHash(key, user)
	if err != nil {
		t.Fatal(err)
	}

	ret := new(User)
	err = redis.GetHashAll(key, ret)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(ret)

	val, err := redis.GetHashInt(key, "Created")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(val)

}
