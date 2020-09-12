package util


import (
	"fmt"
	"testing"
)

func TestAES_CBCTest(t *testing.T) {
	key := "cbcdEf910djlaLO1"  //16位长度
	str := "180123456789"

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
