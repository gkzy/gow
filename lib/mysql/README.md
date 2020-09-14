# github.com/gkzy/gow/lib/mysql


### 使用方法：

*单数据库使用时*

```go
// 初始化链接
func init(){
    config:=&DBConfig{go
    mysql.InitDefaultDB(config)
}

// 按gorm方式使用
mysql.GetORM()
....
```

*多数据库时*

```go
// 初始化链接
func init() {
    // 配置多个DBConfig
    configs:=make([]*DBConfig,0)
    mysql.InitDB(configs)
}

//按gorm方式使用
db:=mysql.GetORMByName("user")
....
```

### 一个demo

```go
type User struct {
    UID      int64  `gorm:"column:uid"`
    Nickname string `gorm:"column:nickname"`
}

func (*User) TableName() string {
    return "tbl_user"
}

func init(){
    config := &mysql.DBConfig{
        Name:     "gkzy",
        User:     "root",
        Password: "123456",
        Host:     "127.0.0.1",
        Port:     6606,
    }

    err := mysql.InitDefaultDB(config)
    if err != nil {
       fmt.Println(err)
    }

}

//MysqlDB gorm风格的数据库操作

func MysqlDB() {
    user := make([]*User, 0)
    db := GetORM()
    err = db.Model(user).Find(&user).Error
    if err != nil {
       fmt.Println(err)
    }
    .....
    ...

}


```