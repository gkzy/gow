package util

import (
	"fmt"
	"testing"
)

type UserToken struct {
	UID      int64  `json:"uid"`
	NickName string `json:"nick_name"`
	Mobile   string `json:"mobile"`
	GroupID  int    `json:"group_id"`
	ProvID   int64  `json:"prov_id"`
}

func TestAES_CBCTest(t *testing.T) {
	key := "cbcdEf910djlaLO1" //16位长度
	str := "18999998888"

	//编码：
	enStr := AESEncrypt(str, key)
	fmt.Println(enStr)

	//解密：
	val, err := AESDecrypt(enStr, key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(val)
}
