/*
使用方法：

一、单数据库使用时：
//init:
func init(){
	config:=&DBConfig{}
	mysql.InitDefaultDB(config)
}
//use:
mysql.GetORM()
....


二、多数据库时：

//init:
func init() {
	configs:=make([]*DBConfig,0)
	mysql.InitDB(configs)
}

use:
db:=mysql.GetORMByName("user")
....

*/

package mysql

import (
	"fmt"
	"github.com/gkzy/gow/lib/logy"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var (
	//dbs
	dbs map[string]*gorm.DB

	//defaultDBName
	defaultDBName string
)

const (
	dbType = "mysql"
)

//DBConfig mysql配置文件
type DBConfig struct {
	Name            string //库名
	User            string //登录名
	Password        string //密码
	Host            string //主机
	Port            int    //port
	Debug           bool   //是否debug
	DisablePrepared bool   //是否disable prepared
}

//InitDefaultDB 单个数据库
func InitDefaultDB(db *DBConfig) (err error) {
	if db == nil {
		err = fmt.Errorf("没有需要init的DB")
		return
	}
	defaultDBName = db.Name
	dbs = make(map[string]*gorm.DB)
	newORM(db)
	return
}

//InitDB InitDB
func InitDB(list []*DBConfig) (err error) {
	if len(list) == 0 {
		err = fmt.Errorf("没有需要init的DB")
		return
	}
	dbs = make(map[string]*gorm.DB)
	for _, item := range list {
		newORM(item)
	}

	return
}

//newORM newORM
func newORM(db *DBConfig) {
	var (
		orm *gorm.DB
		err error
	)
	if db.User == "" || db.Password == "" || db.Host == "" || db.Port == 0 {
		panic(fmt.Sprintf("[DB]-[%v] 数据库配置信息获取失败", db.Name))
	}

	str := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db.User, db.Password, db.Host, db.Port, db.Name) + "?charset=utf8mb4&parseTime=true&loc=Local"
	if db.DisablePrepared {
		str = str + "&interpolateParams=true"
	}
	for orm, err = gorm.Open(dbType, str); err != nil; {
		logy.Error(fmt.Sprintf("[DB]-[%v] 连接异常:%v，正在重试: %v", db.Name, err, str))
		time.Sleep(5 * time.Second)
		orm, err = gorm.Open(dbType, str)
	}
	orm.LogMode(db.Debug)
	orm.CommonDB()
	dbs[db.Name] = orm
}

//GetORM GetORM
func GetORM() *gorm.DB {
	m, ok := dbs[defaultDBName]
	if !ok {
		logy.Panic("[DB] 未init，请参照使用说明")
	}
	return m
}

//GetORM GetORM
func GetORMByName(name string) *gorm.DB {
	m, ok := dbs[name]
	if !ok {
		logy.Panic("[DB] 未init，请参照使用说明")
	}
	return m
}
