package kvdb

type KVDB interface {
	Close()

	Set(key string, data interface{}) error
	Get(key string) ([]byte, error)
	Exists(key string) (bool, error)
	Delete(key string) error
	Keys(key string) ([]string, error)

	HSet(name, key string, data interface{}) error
	HGet(name, key string) ([]byte, error)
	HExists(name, key string) (bool, error)
	HDelete(name, key string) error
	HKeys(name string) ([]string, error)
	HVals(name string) ([][]byte, error)
	HGetall(name string) ([][]byte, error)
}
