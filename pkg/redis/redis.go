package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var defaultRedis *Redis

// 初始化默认数据库，全局唯一的，用包名访问的方法，都存在此数据库中。
func Init(host, auth string, db, maxActive, maxIdle, idleTimeout int) error {
	rdb, err := NewRedis(host, auth, db, maxActive, maxIdle, idleTimeout)
	if err != nil {
		return err
	}
	defaultRedis = rdb
	return nil
}

// 获取默认数据库。
func GetDefaultDB() *Redis {
	return defaultRedis
}

// 新打开一个数据库对象，如果程序中需要同时打开多个数据库，则可以用此方法。
func NewRedis(host, auth string, db, maxActive, maxIdle, idleTimeout int) (*Redis, error) {
	redisConn := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeout),
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp",
				host,
				redis.DialPassword(auth),
				redis.DialDatabase(db))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	_, err := redisConn.Dial()
	// 测试是不是连通
	if err != nil {
		return nil, err
	}
	return &Redis{redisConn: redisConn}, nil
}

// 导出默认数据库对象的方法，方便已包名直接访问

func Close() {
	defaultRedis.Close()
}

func Set(key string, data interface{}) error {
	return defaultRedis.Set(key, data)
}

func Get(key string) ([]byte, error) {
	return defaultRedis.Get(key)
}

func Exists(key string) (bool, error) {
	return defaultRedis.Exists(key)
}

func Delete(key string) error {
	return defaultRedis.Delete(key)
}

func Keys(key string) ([]string, error) {
	return defaultRedis.Keys(key)
}

func HSet(name, key string, data interface{}) error {
	return defaultRedis.HSet(name, key, data)
}

func HGet(name, key string) ([]byte, error) {
	return defaultRedis.HGet(name, key)
}

func HExists(name, key string) (bool, error) {
	return defaultRedis.HExists(name, key)
}

func HDelete(name, key string) error {
	return defaultRedis.HDelete(name, key)
}

func HKeys(name string) ([]string, error) {
	return defaultRedis.HKeys(name)
}

func HVals(name string) ([][]byte, error) {
	return defaultRedis.HVals(name)
}

func HGetall(name string) ([][]byte, error) {
	return defaultRedis.HGetall(name)
}
