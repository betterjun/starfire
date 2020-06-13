package boltdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
)

const (
	DefaultBucket = "__default_bucket"
)

type Boltdb struct {
	boltdb *bolt.DB
}

// 关闭数据库
func (b *Boltdb) Close() {
	b.boltdb.Close()
}

func (b *Boltdb) Set(key string, data interface{}) error {
	return b.HSet(DefaultBucket, key, data)
}

func (b *Boltdb) Get(key string) ([]byte, error) {
	return b.HGet(DefaultBucket, key)
}

func (b *Boltdb) Exists(key string) (bool, error) {
	return b.HExists(DefaultBucket, key)
}

func (b *Boltdb) Delete(key string) error {
	return b.HDelete(DefaultBucket, key)
}

func (b *Boltdb) Keys(key string) (keys []string, err error) {
	err = b.boltdb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DefaultBucket))
		if bucket == nil {
			return fmt.Errorf("hashset not existed")
		}

		c := bucket.Cursor()
		prefix := []byte(key)
		for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			key := make([]byte, len(k))
			copy(key, k)
			keys = append(keys, string(key))
		}

		return nil
	})

	return keys, err
}

func (b *Boltdb) HSet(name, key string, data interface{}) (err error) {
	var dataByte []byte
	switch data.(type) {
	case string:
		dataByte = []byte(data.(string))
	case []byte:
		dataByte = data.([]byte)
	default:
		dataByte, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}

	//fmt.Println("set", name, key, string(dataByte))

	return b.boltdb.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(key), dataByte)
	})
}

func (b *Boltdb) HGet(name, key string) (reply []byte, err error) {
	err = b.boltdb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(name))
		if bucket == nil {
			return fmt.Errorf("hashset not exist")
		}
		v := bucket.Get([]byte(key))
		if len(v) == 0 {
			return fmt.Errorf("record not found")
		}
		reply = make([]byte, len(v))
		copy(reply, v)
		//fmt.Println("get", name, key, string(v))
		return nil
	})

	return reply, err
}

func (b *Boltdb) HExists(name, key string) (exists bool, err error) {
	err = b.boltdb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(name))
		if bucket == nil {
			return fmt.Errorf("hashset not existed")
		}
		v := bucket.Get([]byte(key))
		if len(v) > 0 {
			exists = true
		}
		return nil
	})

	return exists, err
}

func (b *Boltdb) HDelete(name, key string) (err error) {
	return b.boltdb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(name))
		if bucket == nil {
			return fmt.Errorf("hashset not existed")
		}
		return bucket.Delete([]byte(key))
	})
}

func (b *Boltdb) HKeys(name string) (keys []string, err error) {
	err = b.boltdb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(name))
		if bucket == nil {
			return fmt.Errorf("hashset not existed")
		}

		return bucket.ForEach(func(k, v []byte) error {
			key := make([]byte, len(k))
			copy(key, k)
			keys = append(keys, string(key))
			return nil
		})
	})

	return keys, err
}

func (b *Boltdb) HVals(name string) (vals [][]byte, err error) {
	err = b.boltdb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(name))
		if bucket == nil {
			return fmt.Errorf("hashset not existed")
		}

		return bucket.ForEach(func(k, v []byte) error {
			value := make([]byte, len(v))
			copy(value, v)
			vals = append(vals, value)
			return nil
		})
	})

	return vals, err
}

func (b *Boltdb) HGetall(name string) (keyAndValues [][]byte, err error) {
	err = b.boltdb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(name))
		if bucket == nil {
			return fmt.Errorf("hashset not existed")
		}

		return bucket.ForEach(func(k, v []byte) error {
			key := make([]byte, len(k))
			copy(key, k)
			keyAndValues = append(keyAndValues, key)
			value := make([]byte, len(v))
			copy(value, v)
			keyAndValues = append(keyAndValues, value)
			return nil
		})
	})

	return keyAndValues, err
}
