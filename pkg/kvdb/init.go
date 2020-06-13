package kvdb

var defaultEngine KVDB

// 初始化默认数据库，全局唯一的，用包名访问的方法，都存在此数据库中。
func Init(engine KVDB) {
	defaultEngine = engine
}

// 获取默认数据库。
func GetKVDB() KVDB {
	return defaultEngine
}

// 导出默认数据库对象的方法，方便已包名直接访问

func Close() {
	defaultEngine.Close()
}

func Set(key string, data interface{}) error {
	return defaultEngine.Set(key, data)
}

func Get(key string) ([]byte, error) {
	return defaultEngine.Get(key)
}

func Exists(key string) (bool, error) {
	return defaultEngine.Exists(key)
}

func Delete(key string) error {
	return defaultEngine.Delete(key)
}

func Keys(key string) ([]string, error) {
	return defaultEngine.Keys(key)
}

func HSet(name, key string, data interface{}) error {
	return defaultEngine.HSet(name, key, data)
}

func HGet(name, key string) ([]byte, error) {
	return defaultEngine.HGet(name, key)
}

func HExists(name, key string) (bool, error) {
	return defaultEngine.HExists(name, key)
}

func HDelete(name, key string) error {
	return defaultEngine.HDelete(name, key)
}

func HKeys(name string) ([]string, error) {
	return defaultEngine.HKeys(name)
}

func HVals(name string) ([][]byte, error) {
	return defaultEngine.HVals(name)
}

func HGetall(name string) ([][]byte, error) {
	return defaultEngine.HGetall(name)
}
