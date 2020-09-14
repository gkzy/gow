package redis

import (
	"fmt"
	"github.com/gkzy/gow/lib/logy"
	"github.com/gomodule/redigo/redis"
)

type RDSCommon struct {
	client *redis.Pool
}

// GetRDSCommon
func GetRDSCommon() *RDSCommon {
	if redisClient == nil {
		logy.Panic("[redis]连接未初始化")
	}
	if redisClient.Get() == nil {
		logy.Panic("[redis]连接redis失败")
	}
	return &RDSCommon{
		client: redisClient,
	}
}

//==============key操作========================
// GetTTL GetTTL
func (m *RDSCommon) GetTTL(key string) (ttl int64, err error) {
	rc := m.client.Get()
	defer rc.Close()

	ttl, err = redis.Int64(rc.Do("TTL", key))
	return
}

// SetEXPIREAT 设置过期时间(以时间戳的方式)
func (m *RDSCommon) SetEXPIREAT(key string, timestamp int64) (bool, error) {
	rc := m.client.Get()
	defer rc.Close()
	resInt, err := redis.Int64(rc.Do("EXPIREAT", redis.Args{}.Add(key).Add(timestamp)...))
	if err != nil {
		return false, nil
	}
	if resInt == 0 {
		return false, nil
	}
	return true, nil
}

// SetEXPIRE 设置过期时间
func (m *RDSCommon) SetEXPIRE(key string, seconds int64) (bool, error) {
	rc := m.client.Get()
	defer rc.Close()
	resInt, err := redis.Int64(rc.Do("EXPIRE", redis.Args{}.Add(key).Add(seconds)...))
	if err != nil {
		return false, nil
	}
	if resInt == 0 {
		return false, nil
	}
	return true, nil
}

// ReName rename key:dist to newkey
func (m *RDSCommon) ReName(dist, newKey string) (err error) {
	rc := m.client.Get()
	defer rc.Close()
	_, err = redis.String(rc.Do("rename", redis.Args{}.Add(dist).Add(newKey)...))
	return
}

//DEL 删除某个Key
func (m *RDSCommon) DEL(key string) (int64, error) {
	rc := m.client.Get()
	defer rc.Close()
	return redis.Int64(rc.Do("DEL", key))
}

//Exist key是否正在
func (m *RDSCommon) Exists(key string) (bool, error) {
	rc := m.client.Get()
	defer rc.Close()
	return redis.Bool(rc.Do("EXISTS", key))
}

//=============string操作==================

//GetString 取 string
func (m *RDSCommon) GetString(key string) (v string, err error) {
	rc := m.client.Get()
	defer rc.Close()
	v, err = redis.String(rc.Do("GET", key))
	return
}

//GetInt64 取 int64
func (m *RDSCommon) GetInt64(key string) (v int64, err error) {
	rc := m.client.Get()
	defer rc.Close()
	v, err = redis.Int64(rc.Do("GET", key))
	return
}

//S etString 设置值
func (m *RDSCommon) SetString(key string, v interface{}) (ok bool, err error) {
	rc := m.client.Get()
	defer rc.Close()
	result, err := redis.String(rc.Do("SET", redis.Args{}.Add(key).Add(v)...))
	if err != nil && result != "OK" {
		return false, err
	}
	return true, err
}

//SetEx 写入string，同时设置过期时间
//		SetEx(key,"value",24*time.Hour)
func (m *RDSCommon) SetEx(key string, v interface{}, ex int64) (ok bool, err error) {
	rc := m.client.Get()
	defer rc.Close()
	result, err := redis.String(rc.Do("SETEX", redis.Args{}.Add(key).Add(ex).Add(v)...))
	if err != nil && result != "OK" {
		return false, err
	}
	return true, err
}

//===============hash操作=============

//SetHashField 设置某个field值
func (m *RDSCommon) SetHashField(key string, field, v interface{}) (int64, error) {
	rc := m.client.Get()
	defer rc.Close()
	return redis.Int64(rc.Do("HSET", redis.Args{}.Add(key).AddFlat(field).AddFlat(v)...))
}

//SetHash设置 hash值
func (m *RDSCommon) SetHash(key string, v interface{}) (err error) {
	rc := m.client.Get()
	defer rc.Close()
	result, err := redis.String(rc.Do("HMSET", redis.Args{}.Add(key).AddFlat(v)...))
	if err != nil && result != "OK" {
		return
	}
	return err
}

//GetHashInt 获取hash里某个int值
func (m *RDSCommon) GetHashInt(key string, v interface{}) (int64, error) {
	rc := m.client.Get()
	defer rc.Close()
	resInt, err := redis.Int64(rc.Do("HGET", key, v))
	if err != nil {
		return -1, err
	}
	return resInt, nil
}

//GetHashString 获取hash里某个int值
func (m *RDSCommon) GetHashString(key string, v interface{}) (string, error) {
	rc := m.client.Get()
	defer rc.Close()
	return redis.String(rc.Do("HGET", key, v))
}

//GetHash get hash all to obj
func (m *RDSCommon) GetHashAll(key string, obj interface{}) (err error) {
	rc := m.client.Get()
	defer rc.Close()
	var v []interface{}
	v, err = redis.Values(rc.Do("HGETALL", key))
	if len(v) == 0 {
		return fmt.Errorf("未查询到数据")
	}
	err = redis.ScanStruct(v, obj)
	return
}

//HashFieldExists hash某个field是否存在
func (m *RDSCommon) HashFieldExists(key, field string) (bool, error) {
	rc := m.client.Get()
	defer rc.Close()
	return redis.Bool(rc.Do("HEXISTS", redis.Args{}.Add(key).Add(field)...))
}

//================list======================

//================zset======================

//================set======================

//================bit======================
func (m *RDSCommon) GetBit(key string, offset int64) (int, error) {
	rc := m.client.Get()
	defer rc.Close()
	return redis.Int(rc.Do("GETBIT", key, offset))
}

//SetBit
func (m *RDSCommon) SetBit(key string, offset int64, v int) (err error) {
	rc := m.client.Get()
	defer rc.Close()

	if v != 0 && v != 1 {
		err = fmt.Errorf("值只能为0或1")
		return
	}

	_, err = redis.Int(rc.Do("SETBIT", key, offset, v))
	return
}
