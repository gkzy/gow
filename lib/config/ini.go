package config

import (
	"fmt"
	"github.com/gkzy/gini"
	"github.com/gkzy/gow/lib/logy"
	"os"
	"strings"
)

var (
	ini = gini.New("conf/")
)

const (
	defaultConfig     = "app.conf"
	defaultDevConfig  = "dev.app.conf"
	defaultTestConfig = "test.app.conf"
	defaultProdConfig = "prod.app.conf"

	defaultMode = "dev"
	DevMode     = "dev"
	TestMode    = "test"
	ProdMode    = "prod"
)

func init() {
	initConfig()
}

// initConfig use "GOW_RUN_MODE"
func initConfig() {
	fileName := defaultConfig
	runMode := os.Getenv("GOW_RUN_MODE")
	switch runMode {
	case DevMode:
		fileName = defaultDevConfig
	case TestMode:
		fileName = defaultTestConfig
	case ProdMode:
		fileName = defaultProdConfig
	default:
		fileName = defaultConfig
	}
	InitLoad(fileName)
}

// InitLoad 读取指定的配置文件
//	  config.InitLoad("conf/my.ini")
//	  config.GetString("key")
func InitLoad(fileName string) {
	err := ini.Load(fileName)
	if err != nil {
		logy.Warn(fmt.Sprintf("failed to read configuration file：%v err:%v", fileName, err.Error()))
	}
}

// Reload
func Reload() error {
	return ini.ReLoad()
}

// WriteFile write an new file
//	need filename and content
func WriteFile(filename, content string) (n int, err error) {
	return ini.WriteFile(filename, content)
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
	return ini.SectionGet(getSplitSectionKey(key))
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
	return ini.SectionInt(getSplitSectionKey(key))
}

//DefaultInt DefaultInt
func DefaultInt64(key string, def int64) int64 {
	if v, err := GetInt64(key); err == nil {
		return v
	}
	return def
}

func GetInt64(key string) (int64, error) {
	return ini.SectionInt64(getSplitSectionKey(key))
}

//DefaultInt DefaultInt
func DefaultFloat(key string, def float64) float64 {
	if v, err := GetFloat(key); err == nil {
		return v
	}
	return def
}

func GetFloat(key string) (float64, error) {
	return ini.SectionFloat64(getSplitSectionKey(key))
}

//GetInt64
func GetBool(key string) (bool, error) {
	return ini.SectionBool(getSplitSectionKey(key)), nil
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
	keys := ini.GetKeys(section)
	ss := make([]string, 0)
	for _, item := range keys {
		ss = append(ss, item.K)
	}
	return ss
}

// getSplitSectionKey use gini lib
func getSplitSectionKey(name string) (section, key string) {
	if name == "" {
		return
	}
	sp := strings.Split(name, "::")
	switch len(sp) {
	case 1:
		return "", sp[0]
	case 2:
		return sp[0], sp[1]
	default:
		return
	}

}
