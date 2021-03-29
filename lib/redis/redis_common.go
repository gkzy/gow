package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type RDSCommon struct {
	client *redis.Pool
}

// Pool
func (m *RDSCommon) Pool(p *redis.Pool) {
	m.client = p
}

// GetRDSCommon return redis func
func GetRDSCommon() *RDSCommon {
	return cmd
}

/************************************/
/********     REDIS key 	 ********/
/************************************/

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

// GetSet 设置指定 key 的值，并返回 key 的旧值。
func (m *RDSCommon) GetSet(key string, v interface{}) (interface{}, error) {
	rc := m.client.Get()
	defer rc.Close()
	return rc.Do("GETSET", key, v)
}

/************************************/
/********     REDIS string 	 ********/
/************************************/

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

// SetString 设置key的值为v
func (m *RDSCommon) SetString(key string, v interface{}) (ok bool, err error) {
	if m.client==nil{
		return false,fmt.Errorf("缺少redis连接")
	}
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
func (m *RDSCommon) SetEX(key string, v interface{}, ex int64) (ok bool, err error) {
	rc := m.client.Get()
	defer rc.Close()
	result, err := redis.String(rc.Do("SETEX", redis.Args{}.Add(key).Add(ex).Add(v)...))
	if err != nil && result != "OK" {
		return false, err
	}
	return true, err
}

// SETNX key不存时，设置key的值为v
func (m *RDSCommon) SETNX(key string, v interface{}) (bool, error) {
	rc := m.client.Get()
	defer rc.Close()
	return redis.Bool(rc.Do("SETNX", redis.Args{}.Add(key).AddFlat(v)...))
}

// Incr 自增
func (m *RDSCommon) Incr(key string) (num int64, err error) {
	rc := m.client.Get()
	defer rc.Close()
	return redis.Int64(rc.Do("Incr", key))
}

/************************************/
/********     REDIS hash 	 ********/
/************************************/

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

// GetHash get hash all to map
func (m *RDSCommon) GetHashAll2Map(key string)(map[string]string,error) {
	rc := m.client.Get()
	defer rc.Close()
	return redis.StringMap(rc.Do("HGETALL", key))
}

// HMGET key field [field ...]
// 返回哈希表 key 中，一个或多个给定域的值。
// 如果给定的域不存在于哈希表，那么返回一个 nil 值。
func (m *RDSCommon) GetHashfields(key string, fields []string) ([]string, error) {
	rc := m.client.Get()
	defer rc.Close()
	args := redis.Args{}.Add(key)
	for _, value := range fields {
		args = args.Add(value)
	}
	return redis.Strings(rc.Do("HMGET",args...))
}

//HashFieldExists hash某个field是否存在
func (m *RDSCommon) HashFieldExists(key, field string) (bool, error) {
	rc := m.client.Get()
	defer rc.Close()
	return redis.Bool(rc.Do("HEXISTS", redis.Args{}.Add(key).Add(field)...))
}

/************************************/
/********     REDIS list 	 ********/
/************************************/

// SetList 添加list
func (m *RDSCommon) SetList(key string, values []string) error {
	rc := m.client.Get()
	defer rc.Close()
	args := redis.Args{}.Add(key)
	for _, value := range values {
		args = args.Add(value)
	}
	_, err := rc.Do("RPUSH", args...)
	return err
}

// GetList
func (m *RDSCommon) GetList(key string) (data []string, err error) {
	rc := m.client.Get()
	defer rc.Close()
	var values []interface{}
	values, err = redis.Values(rc.Do("LRANGE", key, 0, -1))
	if err != nil {
		return
	}
	if err = redis.ScanSlice(values, &data); err != nil {
		return
	}
	if len(data) == 0 {
		return data, fmt.Errorf("未查询到数据")
	}
	return
}

/************************************/
/********     REDIS ZSET 	 ********/
/************************************/

// SetZSet 添加一个新的ZSet
func (m *RDSCommon) SetZSet(key string, mp map[float64]interface{}) (err error) {
	rc := m.client.Get()
	defer rc.Close()
	_, err = redis.Int64(rc.Do("ZADD", redis.Args{}.Add(key).AddFlat(mp)...))
	return
}

// AddZSet 添加单个ZST
func (m *RDSCommon) AddZSet(key string, score float64, value []byte) (err error) {
	rc := m.client.Get()
	defer rc.Close()
	_, err = rc.Do("ZADD", key, score, value)
	return
}

// getZSetWithScore
func (m *RDSCommon) GetZSetWithScore(key string, offset, limit int) (values []interface{}, err error) {
	rc := m.client.Get()
	defer rc.Close()
	values, err = redis.Values(rc.Do("ZRANGE", key, offset, offset+limit-1, "WITHSCORES"))
	return
}

// GetZSetWithScoreToString ZSET to []string
func (m *RDSCommon) GetZSetWithScoreToStrings(key string, offset, limit int) (ss []string, err error) {
	values, err := m.GetZSetWithScore(key, offset, limit)
	if err != nil {
		return
	}
	if err = redis.ScanSlice(values, &ss); err != nil {
		return
	}
	return
}

// GetZSetWithScoreToInts ZSET to []int
func (m *RDSCommon) GetZSetWithScoreToInts(key string, offset, limit int) (ii []int, err error) {
	values, err := m.GetZSetWithScore(key, offset, limit)
	if err != nil {
		return
	}
	if err = redis.ScanSlice(values, &ii); err != nil {
		return
	}
	return
}

// GetZSetWithScoreToInt64s ZSET to []int64
func (m *RDSCommon) GetZSetWithScoreToInt64s(key string, offset, limit int) (ii []int64, err error) {
	values, err := m.GetZSetWithScore(key, offset, limit)
	if err != nil {
		return
	}
	if err = redis.ScanSlice(values, &ii); err != nil {
		return
	}
	return
}

// GetZSetCountByArea 获取指定区间min-max
//之间成员的数量
func (m *RDSCommon) GetZSetCountByArea(key string, min, max int64) (count int64, err error) {
	rc := m.client.Get()
	defer rc.Close()
	count, err = redis.Int64(rc.Do("ZCOUNT", key, min, max))
	if err != nil {
		return 0, err
	}
	return
}

// GetZSetCount 获取指定区间min-max之间成员的数量
func (m *RDSCommon) GetZSetCount(key string) (count int64, err error) {
	rc := m.client.Get()
	defer rc.Close()
	count, err = redis.Int64(rc.Do("ZCARD", key))
	if err != nil {
		return 0, err
	}
	return
}

// DelZSetMember 删除zset成员
func (m *RDSCommon) DelZSetMember(key string, member interface{}) error {
	rc := m.client.Get()
	defer rc.Close()
	_, err := rc.Do("ZREM", key, member)
	return err
}

/************************************/
/********     REDIS set 	 ********/
/************************************/

// AddSet
func (m *RDSCommon) AddSet(key string, values []string) error {
	rc := m.client.Get()
	defer rc.Close()
	args := redis.Args{}.Add(key)
	for _, value := range values {
		args = args.Add(value)
	}
	_, err := rc.Do("SADD", args...)
	return err
}

// GetSetData
func (m *RDSCommon) GetSetData(key string) (data []string, err error) {
	rc := m.client.Get()
	defer rc.Close()
	var values []interface{}
	values, err = redis.Values(rc.Do("SMEMBERS", key))
	if err != nil {
		return
	}
	if err = redis.ScanSlice(values, &data); err != nil {
		return
	}
	if len(data) == 0 {
		return data, fmt.Errorf("未查询到数据")
	}
	return
}

/************************************/
/********     REDIS bit 	 ********/
/************************************/

// get bit
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
