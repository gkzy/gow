package cache

// ICache 缓存接口
type ICache interface {
	//cache key
	KeyName() string

	//主键名称，
	PrimaryKey(model interface{}) string

	//获取所有数据
	GetAllData() (data interface{}, err error)

}
