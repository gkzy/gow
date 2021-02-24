package config

import (
	"fmt"
	"testing"
)

func TestINI_GetKey(t *testing.T) {

	s := GetString("app_name")

	s1 := GetString("gkzy-user::user")

	i1, _ := GetInt("gkzy-user::port")

	f1, _ := GetFloat("score")

	b1, _ := GetBool("recover_panic")

	fmt.Printf("s type=%T,s=%v \n", s, s)
	fmt.Printf("s1 type=%T,s=%v \n", s1, s1)
	fmt.Printf("i1 type=%T,i1=%v \n", i1, i1)
	fmt.Printf("f1 type=%T,f1=%v \n", f1, f1)
	fmt.Printf("b1, type=%T,b1,=%v \n", b1, b1)

	InitLoad("test.app.conf")
}

//// TestINI_File 指定文件
//func TestINI_File(t *testing.T) {
//	InitLoad("prod.app.conf")
//	fmt.Println(GetString("nsq::producer"))
//}
