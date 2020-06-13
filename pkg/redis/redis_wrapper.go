package redis

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	redisConn *redis.Pool
}

func (r *Redis) Close() {
	r.redisConn.Close()
}

func (r *Redis) Set(key string, data interface{}) (err error) {
	conn := r.redisConn.Get()
	defer conn.Close()

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

	_, err = conn.Do("SET", key, dataByte)
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Get(key string) ([]byte, error) {
	conn := r.redisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (r *Redis) Exists(key string) (bool, error) {
	conn := r.redisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("EXISTS", key))
}

func (r *Redis) Delete(key string) error {
	conn := r.redisConn.Get()
	defer conn.Close()

	_, err := redis.Bool(conn.Do("DEL", key))
	return err
}

func (r *Redis) Keys(key string) ([]string, error) {
	conn := r.redisConn.Get()
	defer conn.Close()

	return redis.Strings(conn.Do("KEYS", key))
}

func (r *Redis) HSet(name, key string, data interface{}) (err error) {
	conn := r.redisConn.Get()
	defer conn.Close()

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

	_, err = conn.Do("HSET", name, key, dataByte)
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) HGet(name, key string) ([]byte, error) {
	conn := r.redisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("HGET", name, key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (r *Redis) HExists(name, key string) (bool, error) {
	conn := r.redisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("HEXISTS", name, key))
}

func (r *Redis) HDelete(name, key string) error {
	conn := r.redisConn.Get()
	defer conn.Close()

	_, err := redis.Bool(conn.Do("HDEL", name, key))
	return err
}

func (r *Redis) HKeys(name string) ([]string, error) {
	conn := r.redisConn.Get()
	defer conn.Close()

	return redis.Strings(conn.Do("HKEYS", name))
}

func (r *Redis) HVals(name string) ([][]byte, error) {
	conn := r.redisConn.Get()
	defer conn.Close()

	return redis.ByteSlices(conn.Do("HVALS", name))
}

func (r *Redis) HGetall(name string) ([][]byte, error) {
	conn := r.redisConn.Get()
	defer conn.Close()

	return redis.ByteSlices(conn.Do("HGETALL", name))
}

func (r *Redis) HGetallToMap(key string) (map[string]string, error) {
	conn := r.redisConn.Get()
	defer conn.Close()
	return redis.StringMap(conn.Do("HGETALL", key))
}
