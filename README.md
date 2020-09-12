# gow

gow 是基于gin的HTTP框架，在gin的基础上，做了更好的html模板封装和数据输出。可用于开发Web API和Web网站项目。


### 项目地址：

[https://github.com/gkzy/gow](https://github.com/gkzy/gow)



### 1. 快速开始

```sh
mkdir hello
cd hello
```

```sh
go mod init
```

```
go get github.com/gkzy/gow
```

##### 1.1 创建 main.go
```go
package main

import (
    "github.com/gkzy/gow"
)

func main() {
    r := gow.Default()

    r.GET("/", func(c *gow.Context) {
        c.JSON(gow.H{
            "code": 0,
            "msg":  "success",
        })
    })
    
    //default :8080
    r.Run()
}
```

##### 1.2 运行

```sh
go run main.go
curl http://127.0.0.1:8080
```
或

```sh
go build && ./hello
```