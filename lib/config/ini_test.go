package config

import (
	"fmt"
	"testing"
)

func TestINI_GetKey(t *testing.T) {

	s := GetString("app_name")

	fmt.Printf("s type=%T,s=%v \n", s, s)

	_, err := WriteContent("123456")
	if err != nil {
		fmt.Println("err:", err)
	}


	Reload()

	s = GetString("app_name")

	fmt.Printf("s type=%T,s=%v \n", s, s)


}
