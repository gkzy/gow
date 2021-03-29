# github.com/gkzy/gow/lib/cache 

> go mem cache 封装


### 使用的第三方库

```go
github.com/patrickmn/go-cache
```

### 引用 

```go
import "github.com/gkzy/gow/lib/cache"
```

### 使用

*cache_test.go*

```go
/*
一个实现的demo
*/

package cache

import (
    "fmt"
    "testing"
)

//================原始数据操作层========================

//prov model
type Prov struct {
    ID   int64
    Name string
}

//from mysql or other...
func (m *Prov) GetAllDataFromDB() (data []*Prov, err error) {
    data = make([]*Prov, 0)
    data = append(data, &Prov{
        ID:   51,
        Name: "四川",
    })
    data = append(data, &Prov{
        ID:   50,
        Name: "重庆",
    })
    data = append(data, &Prov{
        ID:   11,
        Name: "北京",
    })
    return
}

//===========================cache 实现=====================================

// ProvCache 省份信息缓存 一个ICache的实现
type ProvCache struct{}

//KeyName KeyName
func (m *ProvCache) KeyName() string {
    return "prov"
}

//PrimaryKey PrimaryKey
func (m *ProvCache) PrimaryKey(obj interface{}) string {
    return fmt.Sprintf("%v", obj.(*Prov).ID)
}

//GetAllData
func (m *ProvCache) GetAllData() (data interface{}, err error) {
    data, err = new(Prov).GetAllDataFromDB()
    return
}

//===================使用者===============

//所有prov数据
func TestCacheUtil_GetAll(t *testing.T) {
    nc := NewMemCache()
    data, err := nc.GetAll(new(ProvCache))
    if err != nil {
        t.Fatal(err)
    }

    fmt.Println("===test all===")

    //返回值断言
    value, ok := data.([]*Prov)
    if ok {
        fmt.Printf("value: %#v \n", value)
        for _, item := range value {
            fmt.Println(item)
        }
    }

}

//单个prov数据
func TestCacheUtil_GetItem(t *testing.T) {
    nc := NewMemCache()
    data, err := nc.GetItem(new(ProvCache), 51)
    if err != nil {
        t.Fatal(err)
    }

    fmt.Println("===test item===")
    value, ok := data.(*Prov)
    if ok {
        fmt.Println(value)
    }

}

```