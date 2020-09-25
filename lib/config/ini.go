package config

import (
	ini "github.com/go-ini/ini"
	"strings"
)

var (
	cfg = ini.Empty()
)

const (
	defaultConfig     = "conf/app.conf"
	defaultDevConfig  = "conf/dev.app.conf"
	defaultProdConfig = "conf/prod.app.conf"
)



// InitLoad 读取指定的配置文件
//	  config.InitLoad("conf/my.ini")
//	  config.GetString("key")
func InitLoad(fileName string) {
	var err error
	cfg, err = ini.Load(fileName)
	if err != nil {
		panic("Failed to read configuration file：" + fileName)
	}
}

// DefaultString get default string
//	 config.DefaultString("prov","四川")
func DefaultString(key, def string) string {
	if v := GetString(key); v != "" {
		return v
	}
	return def
}

// GetString get string
//	 config.GetString("prov")
func GetString(key string) string {
	return getKey(key).String()
}

//DefaultInt get default int
//	config.DefaultInt("prov_id",1)
func DefaultInt(key string, def int) int {
	if v, err := GetInt(key); err == nil {
		return v
	}
	return def
}

// GetInt  get int
//	config.GetInt("prov_id)
func GetInt(key string) (int, error) {
	return getKey(key).Int()
}

//DefaultInt DefaultInt
func DefaultInt64(key string, def int64) int64 {
	if v, err := GetInt64(key); err == nil {
		return v
	}
	return def
}

func GetInt64(key string) (int64, error) {
	return getKey(key).Int64()
}

//DefaultInt DefaultInt
func DefaultFloat(key string, def float64) float64 {
	if v, err := GetFloat(key); err == nil {
		return v
	}
	return def
}

func GetFloat(key string) (float64, error) {
	return getKey(key).Float64()
}

//GetInt64
func GetBool(key string) (bool, error) {
	return getKey(key).Bool()
}

//DefaultBool DefaultBool
func DefaultBool(key string, def bool) bool {
	if v, err := GetBool(key); err == nil {
		return v
	}
	return def
}

// Keys 获取section下的所有keys
func Keys(section string) []string {
	return cfg.Section(section).KeyStrings()
}

//getKey getKey
func getKey(key string) *ini.Key {
	if key == "" {
		return nil
	}
	sp := strings.Split(key, "::")
	switch len(sp) {
	case 1:
		return cfg.Section("").Key(key)
	case 2:
		return cfg.Section(sp[0]).Key(sp[1])
	default:
		return nil
	}

}
