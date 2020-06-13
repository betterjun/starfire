package boltdb

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
)

var defaultBoltdb *Boltdb

// 初始化默认数据库，全局唯一的，用包名访问的方法，都存在此数据库中。
func Init(file string) error {
	bdb, err := NewBoltdb(file)
	if err != nil {
		return err
	}
	defaultBoltdb = bdb
	return nil
}

// 获取默认数据库。
func GetDefaultDB() *Boltdb {
	return defaultBoltdb
}

// 新打开一个数据库对象，如果程序中需要同时打开多个数据库，则可以用此方法。
func NewBoltdb(file string) (*Boltdb, error) {
	bdb, err := bolt.Open(file, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = bdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(DefaultBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		bdb.Close()
		return nil, err
	}

	return &Boltdb{boltdb: bdb}, nil
}

// 导出默认数据库对象的方法，方便已包名直接访问

func Close() {
	defaultBoltdb.Close()
}

func Set(key string, data interface{}) error {
	return defaultBoltdb.Set(key, data)
}

func Get(key string) ([]byte, error) {
	return defaultBoltdb.Get(key)
}

func Exists(key string) (bool, error) {
	return defaultBoltdb.Exists(key)
}

func Delete(key string) error {
	return defaultBoltdb.Delete(key)
}

func Keys(key string) ([]string, error) {
	return defaultBoltdb.Keys(key)
}

func HSet(name, key string, data interface{}) error {
	return defaultBoltdb.HSet(name, key, data)
}

func HGet(name, key string) ([]byte, error) {
	return defaultBoltdb.HGet(name, key)
}

func HExists(name, key string) (bool, error) {
	return defaultBoltdb.HExists(name, key)
}

func HDelete(name, key string) error {
	return defaultBoltdb.HDelete(name, key)
}

func HKeys(name string) ([]string, error) {
	return defaultBoltdb.HKeys(name)
}

func HVals(name string) ([][]byte, error) {
	return defaultBoltdb.HVals(name)
}

func HGetall(name string) ([][]byte, error) {
	return defaultBoltdb.HGetall(name)
}
