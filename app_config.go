package gow

import (
	"github.com/gkzy/gow/lib/config"
	"os"
)

const (
	defaultConfig     = "conf/app.conf"
	defaultDevConfig  = "conf/dev.app.conf"
	DefaultProdConfig = "conf/prod.app.conf"
)

// AppConfig gku app 统一配置入口
//		可以通过AppConfig完成统一的app基础配置
type AppConfig struct {
	AppName       string //应用名称
	RunMode       string //运行模板
	HttpAddr      string //监听端口
	AutoRender    bool   //是否自动渲染html模板
	Views         string //html模板目录
	TemplateLeft  string //模板符号
	TemplateRight string //模板符号
	SessionOn     bool   //是否打开session
}

// GetAppConfig 获取配置文件中的信息
//  默认使用conf/app.conf配置文件
//  当环境变量 APP_RUN_MODE ="prod"时，使用 conf/prod.app.conf
//  当环境变量 APP_RUN_MODE ="dev"时，使用 conf/dev.app.conf
//  没有此环境变量时，使用conf/app.conf
func GetAppConfig() *AppConfig {
	var (
		fileName string
	)
	runMode := os.Getenv("GKU_RUN_MODE")

	switch runMode {
	case devMode:
		fileName = defaultDevConfig
	case prodMode:
		fileName = DefaultProdConfig
	default:
		fileName = defaultConfig
	}
	if fileName == "" {
		fileName = defaultConfig
	}

	config.InitLoad(fileName)

	return &AppConfig{
		AppName:       config.DefaultString("app_name", "gku"),
		RunMode:       config.DefaultString("run_mode", "dev"),
		HttpAddr:      config.DefaultString("http_addr", ":8080"),
		AutoRender:    config.DefaultBool("auto_render", false),
		Views:         config.DefaultString("views", "views"),
		TemplateLeft:  config.DefaultString("template_left", "{{"),
		TemplateRight: config.DefaultString("template_right", "}}"),
		SessionOn:     config.DefaultBool("session_on", false),
	}

}
