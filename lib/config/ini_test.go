package config

import (
	"fmt"
	"strings"
	"testing"
)

func TestINI_GetKey(t *testing.T) {
	fmt.Println(DefaultString("app_name", "gow"))
	fmt.Println(DefaultString("app_mode", "dev"))
	fmt.Println(DefaultString("http_port", "8080"))

	fmt.Println(DefaultString("gkzy-user::user", "zituocn"))

	keys := Keys("gkzy-user")
	fmt.Println(strings.Join(keys, ","))
}

// TestINI_File 指定文件
func TestINI_File(t *testing.T){
	InitLoad("conf/prod.app.conf")
	fmt.Println(GetString("template_left"))
	fmt.Println(GetString("template_right"))
}