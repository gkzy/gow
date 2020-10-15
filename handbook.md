# gow 使用手册

gow 是基于gin源码的HTTP框架，在gin的基础上，做了更好的html模板封装和数据输出。可用于开发Web API和Web网站项目

> v0.1.2

## 1. 项目地址

[https://github.com/gkzy/gow](https://github.com/gkzy/gow)


## 2. 快速开始

```sh
mkdir hello
cd hello
```

```sh
go mod init
```

### 2.1 创建 main.go

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

### 2.2 运行

```sh
go run main.go
curl http://127.0.0.1:8080
```

或

```sh
go build && ./hello
```

---

## 3. 配置文件

* package

```sh
github.com/gkzy/gow/lib/config
```

### 3.1 初始化配置

* 可在代码中，使用以下代码初始化配置

```go
gow.InitConfig()
```

* gow.InitConfig方法 的实现

```go
func InitConfig() {
    fileName := ""
    runMode := os.Getenv("GOW_RUN_MODE")
    switch runMode {
    case DevMode:
        fileName = defaultDevConfig
    case ProdMode:
        fileName = defaultProdConfig
    default:
        fileName = defaultConfig
    }
    if fileName == "" {
        fileName = defaultConfig
    }

    config.InitLoad(fileName)
}
```

* 可读取 `GOW_RUN_MODE` = dev | prod 环境变量，实现配置文件的不同加载
* 对应使用 `conf/app.conf` 和 `conf/prod.app.conf` 配置文件

### 3.2 加载全局配置

* 配置文件内容

```ini
app_name = User-Service
run_mode = dev
http_addr = 8080
auto_render = false
session_on = false
```


* 读取配置文并应用

```go

func main() {
    r := gow.Default()
    // 加载配置并应用
    r.SetAppConfig(gow.GetAppConfig())
    //路由
    routers.APIRouter(r)
    //运行
    r.Run()
}
```

### 3.3 读取自定义配置

* 配置文件如下

```ini
app_name = User-Service
run_mode = dev
http_addr = 8080
auto_render = false
session_on = false

user = "root"
password = "root123"
host = "192.168.0.197"
port = 3306
```

```go

// 带默认值
runMode := config.DefaultString("run_mode", "dev")

// 不带默认值
user:= config.GetString("gkzy-user::user")

// 带默认值的int
port:= config.DefaultInt("gkzy-user::port",3306)
```
---


## 4. middleware

* 基本格式

```go
func FunName() gow.HandlerFunc {
    return func(c *gow.Context) {
        ......
        .....
        //需要执行此方法，否则不会执行后面的请求
        c.Next() 
    }
}
```

### 4.1 middleware 调用 

```go
r := gow.Default()
r.Use(...)
r.Run()
```
* 在 gow.Default() 方法内，已经默认调用了两个middleware

```go
func Default() *Engine {
    engine := New()
    engine.Use(Logger(), Recovery())
    return engine
}
```

### 4.2 gow自带的middleware

```go
//日志
Logger()

//recovery 
Recovery()

//session
Session()

//翻页相关
DataPager()
```

### 4.3 自定义一个middleware

```go
// APIAuth API通讯的基础鉴权
//  用户token可能不存在
func APIAuth() gow.HandlerFunc {
    return func(c *gow.Context) {
        // 非正式环境不鉴权
        if !c.IsProd() {
            return
        }

        auth := c.GetHeader("Authorization")
        source, _ := strconv.ParseInt(c.GetHeader("source"), 10, 64)
        token := c.GetHeader("token")
        timeStamp, _ := strconv.ParseInt(c.GetHeader("timestamp"), 10, 64)

        if source < 1 {
            c.DataJSON(403, "没有权限：缺少source")
            c.StopRun()
        }
        if auth == "" {
            c.DataJSON(403, "没有权限：缺少Authorization")
            c.StopRun()
        }
        if timeStamp < 1 {
            c.DataJSON(403, "没有权限：缺少timestamp")
            c.StopRun()
        }

        now := time.Now().Unix()
        var timeOut float64 = 60
        //过期请求验证
        if math.Abs(float64(now-timeStamp)) > timeOut {
            c.DataJSON(403, "没有权限，请求已过期")
            c.StopRun()
        }

        appId, appSecret := getAppInfo(source)
        key := fmt.Sprintf("%v@%v@%v@%v@%v", appId, appSecret, source, token, timeStamp)

        serverAuth := strings.ToUpper(util.MD5(key))
        if serverAuth != auth {
            c.DataJSON(403, "没有权限：接口权限验证失败")
            c.StopRun()
        }

        // 此方法一定不能漏掉
        c.Next()
    }
}
```
使用

```go
func main() {
    r := gow.Default()
    r.SetAppConfig(gow.GetAppConfig())
    v1:=r.Group("/v1")
    // 调用 
    v1.Use(APIAuth())
    v1.GET("/test", func(c *gow.Context) {
        c.JSON(gow.H{
            "code": 0,
            "msg":  "success",
        })
    })
    r.Run()
}
```
---

## 5. 日志 logy

* package

```
github.com/gkzy/gow/lib/logy
```

### 5.1 直接使用

* 基本方法

```go
logy.Info(...)
logy.Notice(...)
logy.Debug(...)
logy.Error(...)
logy.Warn(...)
logy.Fatal(...)
logy.Panic(...)
```
* 默认输出到 os.Stdout 
```sh
[gow-site] 2020/10/13 09:38:06 [I] gow.go:307: [gow-site] [dev] Listening and serving HTTP on http://:8080
[gow-site] 2020/10/13 09:38:09 [I]  302 |      69.788µs |       127.0.0.1 | GET      /
[gow-site] 2020/10/13 09:38:09 [I]  200 |    1.450247ms |       127.0.0.1 | GET      /docs/start
[gow-site] 2020/10/13 09:38:09 [I]  200 |    4.249205ms |       127.0.0.1 | GET      /static/js/tabler.min.js?v=1602553086
[gow-site] 2020/10/13 09:38:09 [I]  200 |    3.162056ms |       127.0.0.1 | GET      /static/css/md.css?v=1602553086
[gow-site] 2020/10/13 09:38:09 [I]  200 |    3.099501ms |       127.0.0.1 | GET      /static/css/prism.css?v=1602553086
[gow-site] 2020/10/13 09:38:09 [I]  200 |    6.519717ms |       127.0.0.1 | GET      /static/js/prism.js?v=1602553086
[gow-site] 2020/10/13 09:38:09 [I]  200 |    4.583973ms |       127.0.0.1 | GET      /static/css/gow.css?v=1602553086
```

### 5.1 同时记录到文件


* 配置方法

```go
// InitLog init logy
func InitLog() {
    runMode := config.DefaultString("run_mode", "dev")
    //正式环境到控制台和文件
    if runMode == gow.ProdMode {
        logy.SetOutput(
            logy.MultiWriter(
                logy.WithColor(logy.NewWriter(os.Stdout)),
                logy.NewFileWriter(logy.FileWriterOptions{
                    Dir:           "./logs",
                    Prefix:        "web",
                    StorageMaxDay: 7,
                }),
            ),
            "User-Service",
        )
    } else {
        //开发环境只到控制台
        logy.SetOutput(
            logy.WithColor(logy.NewWriter(os.Stdout)),
             "User-Service",
        )
    }
}
```

* 调用

```go
package main

import (
    "github.com/gkzy/gow"
)

func init() {
    //init 配置
    gow.InitConfig()

    //init 日志
    InitLog()

}

func main() {
    r := gow.Default()
    r.SetAppConfig(gow.GetAppConfig())
    routers.APIRouter(r)
    r.Run()
}
```
---

## 6. 路由

* 基础方法

```go
package main

import (
    "github.com/gkzy/gow"
)

func main() {
    r := gow.Default()

    r.GET("/someGet", getting)
    r.POST("/somePost", posting)
    r.PUT("/somePut", putting)
    r.DELETE("/someDelete", deleting)
    r.PATCH("/somePatch", patching)
    r.HEAD("/someHead", head)
    r.OPTIONS("/someOptions", options)
    r.Any("/some",handler)

    r.Run()
}
```

* 一个路由方法及调用

```go
// APIRouter handler to router path
func APIRouter(r *gow.Engine) {

    //分组
    sn := r.Group("/" + ServerName)

    //API auth middleware
    sn.Use(middleware.APIAuth())

    v1 := sn.Group("/v1")
    {
        //无权限接口
        allow := v1.Group("/allow")
        {

            allow.POST("/event", event.CreateReport)
            allow.POST("/msg/send", sms.SendSMS)
        }

        //有用户权限的接口
        auth := v1.Group("/auth")
        auth.Use(middleware.UserAuth()) //用户鉴权middleware
        {
            auth.GET("/test", login.TestAuth)

        }
    }
}
```

```go
func main() {
    r := gow.Default()
    r.SetAppConfig(gow.GetAppConfig())
    // 调用 
    routers.APIRouter(r)
    r.Run()
}
```

## 7. 获取值

### 7.1 获取路由参数 (router param)
 
```go
r.GET("/article/:id", handler)
```

```go
id:=c.Param("id")
```

### 7.2 获取请求参数(query param && form param)

```go
func GetUser(c *gow.Context){

    //获取字串
    c.GetString("key","default")

    //获取int
    c.GetInt("key",0)

    //获取bool
    c.GetBool("key",false)

    //获取int64
    c.GetInt64("key",-1)

    //获取float
    c.GetFloat("key",0)

    //获取[]string
    var ids []string
    ids = c.GetStrings("ids")  

    //其他方法
    c.GetInt32()
    c.GetInt16()
    c.GetInt8()
    ....
}
```
###  7.3 获取 request body 

* 获取 body

```go
func (c *Context) Body() []byte
```

* 获取 并JSON反序列化到 v

```go
func (c *Context) DecodeJSONBody(v interface{}) error 
```

* demo
```go
type User struct {
    Nickname string `json:"nickname"`
    QQ       string `json:"qq"`
}

func GetUser(c *Context){
    user := new(User)
    //
    err := c.DecodeJSONBody(&user)
    if err != nil {
        //handler error
    }

    c.JSON(gow.H{
        "user": user,
    })   

}
```

### 7.4 文件上传

* 当需要上传大于32MB的文件时，请使用以下配置

```go
r.MaxMultipartMemory = 1<<22 //64MB
```

* 上传单个文件

```go
func (c *Context) GetFile(key string) (multipart.File, *multipart.FileHeader, error) 
```

* 上传多个文件

```go
func (c *Context) GetFiles(key string) ([]*multipart.FileHeader, error) 
```

* 保存文件到服务器
```go
func (c *Context) SaveToFile(fromFile, toFile string) error 
```

* 上传演示

```
<form enctype="multipart/form-data" method="post">
    <input type="file" name="uploadname" />
    <input type="submit">
</form>
```

```go

package main

import (
    "github.com/gkzy/gow"
)

func main() {
    r := gow.Default()
    r.SetAppConfig(gow.GetAppConfig())
    r.POST("/upload",UploadFile)
    r.Run()
}

func UploadFile(c *gow.Context){
    f,h,err:=c.GetFile("file")
    if err!=nil{
        log.Fatal("getfile err ", err)
    }
    defer f.Close()
    c.SaveToFile("file","upload/"+h.Filename) //保存在upload下，没有目录，需要先创建
}
```
---

## 8. 输出值

* String 

```go
func GetUser(c *gow.Context){

    //default http.StatusOK
    c.String("hello gow...")

   //或者，指定 http.StatusCode
    c.ServerString(200,"hello gow...")
}
```

* JSON

```go
func GetUser(c *gow.Context){
    //default http.StatusOK
    c.JSON(gow.H{
        "nickname":"gow",
        "age":1,
    })

   //或者，指定 http.StatusCode
    c.ServerJSON(200,gow.H{
        "nickname":"gow",
        "age":1,
    })
}
```

* XML

```go
func GetUser(c *gow.Context){

    //default http.StatusOK
    c.XML(gow.H{
        "nickname":"gow",
        "age":18,
    })

   //或者，指定 http.StatusCode
    c.ServerXML(200,gow.H{
        "nickname":"gow",
        "age":18,
    })
}

```

* YAML

```go
func GetUser(c *gow.Context){

    //default http.StatusOK
    c.YAML(gow.H{
        "nickname":"gow",
        "age":18,
    })

   //或者，指定 http.StatusCode
    c.ServerYAML(200,gow.H{
        "nickname":"gow",
        "age":18,
    })
}
```

* 读取文件

```go
func GetUser(c *gow.Context){
    //读取go.md并输出
    c.File("go.mod")
}
```

* 下载

```go

// 下载指定内容
func GetUser(c *gow.Context){
    c.Download([]byte("download string"))
}
```

```go
// 读取main.go并下载到main.txt
func GetUser(c *gow.Context){
    c.FileAttachment("main.go","main.txt")
}
```
---


## 9. 做网站

### 9.1  目录结构

```
PROJECT_NAME
├──static
      ├── img
            ├──111.jpg
            ├──222.jpg
            ├──333.jpg
      ├──js
      ├──css
├──views
    ├──index.html
    ├──article
        ├──detail.html
├──main.go
```

### 9.2 基础方法

* 设置模板目录

```go
r := gow.Default()
r.SetView("views")
```

* 打开模板渲染

```go
r.AutoRender = true
```

* 设置表态资源

```go
r.Static("/static", "static")
```

* 设置favicon.ico

```go
r.StaticFile("favicon.ico","static/img/log.png") 
```

* 设置favicon.ico

```go
r.AddFuncMap(key string,fn interface(){})
```

* 设置html模板符号

```go
r.SetDelims("{{","}}")
```
* 向模板传递数据

```go
c.HTML("article/detail.html", gow.H{
    "title":    "这是一个文章标题",
})
```

```
<title>{{.title}}</title>
```
* 使用 `context.Data` 

```go
c.Data["title"] = "这是一个文章标题"
c.HTML("article/defail.html")
```

```
<title>{{.title}}</title>
```


### 9.3 演示代码

* main.go

```go

package main

import (
    "github.com/gkzy/gow"
)

func main() {
    r := gow.Default()
    r.AutoRender = true //打开html模板渲染
    r.SetView("views") //默认静态目录为views时，可不调用此方法
    r.StaticFile("favicon.ico","static/img/favicon.png")  //路由favicon.ico
    r.Static("/static", "static")

    //router
    r.Any("/", IndexHandler)
    r.Any("/article/1", ArticleDetailHandler)


    //自定义hello的模板函数
    //在模板文件可通过 {{hello .string}} 来执行
    r.AddFuncMap("hello", func(str string) string {
        return "hello:" + str
    })

    r.Run()
}

//IndexHandler 首页
func IndexHandler(c *gow.Context) {
    c.HTML("index.html", gow.H{
        "name":    "gow",
        "package": "github.com/gkzy/gow",
    })
}

//ArticleDetailHandler 文件详情页
func ArticleDetailHandler (c *gow.Context){
    d.Data["title"] = "这是一个标题"
    c.HTML("article/detail.html")
}
```
* 访问

```sh
https://127.0.0.1:8080/
https://127.0.0.1:8080/article/1
```
---

## 10. cookie&&session

### 10.1 cookie 方法

* 设置cookie

```go
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
```

* 读取cookie

```go
func (c *Context) GetCookie(name string) (string, error)
```

### 10.2 session 方法

* 开启session

```go
r := gow.Default()
r.SetSessionOn(true)
r.Run()
```

* 设置session

```go
func (c *Context) SetSession(key string, v interface{})
```

* 获取session

```go
func (c *Context) GetSession(key string) interface{}
func (c *Context) SessionString(key string) string 
func (c *Context) SessionInt(key string) int
func (c *Context) SessionInt64(key string) int64
func (c *Context) SessionBool(key string) bool
```

* 删除session

```go
func (c *Context) DeleteSession(key string)
```

### 10.3 演示代码

```go
package main

import (
    "github.com/gkzy/gow"
    "time"
)

func main() {
    r := gow.Default()

    r.SetSessionOn(true)

    r.GET("/session/set", SetUser)
    r.GET("/session/get", GetUser)
    r.GET("/session/del", DelUser)

    r.GET("/cookie/set", SetTopic)
    r.GET("/cookie/get", GetTopic)
    r.GET("/cookie/del", DelTopic)

    r.Run()
}

var (
    key = "nickname"
)

func SetUser(c *gow.Context) {
    c.SetSession(key, "TEST")
    c.String("OK")
}

func GetUser(c *gow.Context) {
    val := c.SessionString(key)
    c.String(val)

}

func DelUser(c *gow.Context) {
    c.DeleteSession(key)
}

//=======cookie=========

var (
    topicKey = "topic"
)

func SetTopic(c *gow.Context) {
    c.SetCookie(topicKey, "这是一个topic", int(10*time.Minute), "/", "", false, true)
}

func GetTopic(c *gow.Context) {
    v, _ := c.GetCookie(topicKey)
    c.String(v)
}

func DelTopic(c *gow.Context) {
    c.SetCookie(topicKey, "", -1, "/", "", false, true)
}
```
---

## 11. 数据分页

### 11.1 使用 DataPager middleware

* DataPager() 实现

```go
// DataPager middlewares
//  实现分页参数的处理
func DataPager() HandlerFunc {
    return func(c *Context) {
        pager := new(Pager)
        pager.Page, _ = c.GetInt64("page", 1)
        if pager.Page < 1 {
            pager.Page = 1
        }
        pager.Limit, _ = c.GetInt64("limit", 10)
        if pager.Limit < 1 {
            pager.Limit = 1
        }

        pager.Offset = (pager.Page - 1) * pager.Limit
        c.Pager = pager
        c.Next()
    }
}
```

* Pager struct

```go
type Pager struct {
    Page      int64 `json:"page"`
    Limit     int64 `json:"-"`
    Offset    int64 `json:"-"`
    Count     int64 `json:"count"`
    PageCount int64 `json:"pagecount"`
}
```

* 使用 DataPager()

```go
func main() {
    r := gow.Default()
    r.Use(gow.DataPager())
    r.GET("/", GetUser)
    r.Run()
}
```

### 11.2 设置 count

```go
func GetUser(c *gow.Context) {
    //设置总条数
    c.Pager.Count = 100
}
```

### 11.3 使用 offset与limit

> gorm 数据库翻页查询

```go
db:=conn.GetORM()
db.Model(xxx).Limit(c.Pager.Limit).Offset(c.Pager.Offset)....
```

### 11.3 使用 DataJSON输出

```go
func (c *Context) DataJSON(args ...interface{})
```

### 11.4 完整的例子

* main.go

```go
package main

import (
    "github.com/gkzy/gow"
)

type User struct {
    Nickname string `json:"nickname"`
    Age      int    `json:"age"`
}

func main() {
    r := gow.Default()
    r.Use(gow.DataPager())
    r.GET("/users", GetUser)
    r.Run()
}

func GetUser(c *gow.Context) {
    users := make([]*User, 0)
    //设置总条数
    c.Pager.Count = 100
    //输出[]*User和c.Pager
    c.DataJSON(&users, c.Pager)
}

```

* 访问

```
http://127.0.0.1:8080/users?page=1&limit=15
```

* 响应
```json
{
  "code": 0,
  "msg": "success",
  "time": 1602726745,
  "body": {
    "pager": {
      "page": 2,
      "count": 100,
      "pagecount": 7
    },
    "data": [
      
    ]
  }
}
```



---
## 15. 扩展库

### package 

```sh
github.com/gkzy/gow/lib/
```


* config

```
github.com/gkzy/gow/lib/config
```

* logy
```
github.com/gkzy/gow/lib/logy
```

* memory cache

```
github.com/gkzy/gow/lib/cache
```


* mysql

```
github.com/gkzy/gow/lib/mysql
```

* nsq

```
github.com/gkzy/gow/lib/nsq
```

* oauth

```go
github.com/gkzy/gow/lib/oauth/apple
github.com/gkzy/gow/lib/oauth/wechat
```

* oss

```
github.com/gkzy/gow/lib/oss
```


* pay

```go
github.com/gkzy/gow/lib/pay/alipay
github.com/gkzy/gow/lib/pay/wechat
github.com/gkzy/gow/lib/pay/apple
```

* pdf

```go
github.com/gkzy/gow/lib/pdf
```

* redis

```go
github.com/gkzy/gow/lib/redis
```

* grpc

```go
github.com/gkzy/gow/lib/rpc
```

* sms

```go
github.com/gkzy/gow/lib/sms
```




