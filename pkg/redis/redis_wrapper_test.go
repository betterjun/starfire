package redis

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedisWrapper(t *testing.T) {
	db, err := NewRedis("47.108.94.209:6379", "Sdjx2020", 8, 30, 30, 200)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, db)
	defer db.Close()

	keyAndValues := map[string]interface{}{
		"keyInt":    1,
		"keyByte":   []byte("these are bytes"),
		"keyString": "these are string中文",
		"keyMap":    map[string]interface{}{"name": "betterjun", "age": 10},
	}

	/*
		Set(key string, data interface{}) error
			Get(key string) ([]byte, error)
			Exists(key string) (bool, error)
			Delete(key string) (error)
	*/
	for k, v := range keyAndValues {
		err := db.Set(k, v)
		assert.Equal(t, nil, err)

		v1, err := db.Get(k)
		assert.Equal(t, nil, err)

		switch v.(type) {
		case string:
			assert.Equal(t, 0, bytes.Compare([]byte(v.(string)), v1))
		case []byte:
			assert.Equal(t, 0, bytes.Compare(v.([]byte), v1))
		default:
			vByte, err := json.Marshal(v)
			if err != nil {
				assert.Error(t, err)
			}
			assert.Equal(t, 0, bytes.Compare(vByte, v1))
		}

		exists, err := db.Exists(k)
		assert.Equal(t, nil, err)
		assert.Equal(t, true, exists)

		err = db.Delete(k)
		assert.Equal(t, nil, err)

		exists, err = db.Exists(k)
		assert.Equal(t, nil, err)
		assert.Equal(t, false, exists)
	}

	/*
		HSet(name, key string, data interface{}) error
			HGet(name, key string) ([]byte, error)
			HExists(name, key string) (bool, error)
			HDelete(name, key string) (error)
	*/
	hashset := "hs"
	for k, v := range keyAndValues {
		err := db.HSet(hashset, k, v)
		assert.Equal(t, nil, err)

		v1, err := db.HGet(hashset, k)
		assert.Equal(t, nil, err)

		switch v.(type) {
		case string:
			assert.Equal(t, 0, bytes.Compare([]byte(v.(string)), v1))
		case []byte:
			assert.Equal(t, 0, bytes.Compare(v.([]byte), v1))
		default:
			vByte, err := json.Marshal(v)
			if err != nil {
				assert.Error(t, err)
			}
			assert.Equal(t, 0, bytes.Compare(vByte, v1))
		}

		exists, err := db.HExists(hashset, k)
		assert.Equal(t, nil, err)
		assert.Equal(t, true, exists)

		err = db.HDelete(hashset, k)
		assert.Equal(t, nil, err)

		exists, err = db.HExists(hashset, k)
		assert.Equal(t, nil, err)
		assert.Equal(t, false, exists)
	}

	/*
		Keys(key string) ([]interface{}, error)
		HKeys(name string) ([]interface{}, error)
	*/
	for k, v := range keyAndValues {
		err := db.Set(k, v)
		assert.Equal(t, nil, err)

		err = db.HSet(hashset, k, v)
		assert.Equal(t, nil, err)
	}

	keys, err := db.Keys("")
	for _, k := range keys {
		v1, err := db.Get(k)
		assert.Equal(t, nil, err)
		v2, err := db.HGet(hashset, k)
		assert.Equal(t, nil, err)
		assert.Equal(t, v1, v2)
	}

	hkeys, err := db.HKeys(hashset)
	assert.Equal(t, nil, err)
	//assert.Equal(t, len(keys), len(hkeys)+1)

	// HVals(name string) ([]interface{}, error)
	hvals, err := db.HVals(hashset)
	assert.Equal(t, nil, err)
	assert.Equal(t, len(hkeys), len(hvals))

	for i, k := range hkeys {
		v := keyAndValues[k]
		switch v.(type) {
		case string:
			assert.Equal(t, 0, bytes.Compare([]byte(v.(string)), hvals[i]))
		case []byte:
			assert.Equal(t, 0, bytes.Compare(v.([]byte), hvals[i]))
		default:
			vByte, err := json.Marshal(v)
			if err != nil {
				assert.Error(t, err)
			}
			assert.Equal(t, 0, bytes.Compare(vByte, hvals[i]))
		}
	}

	// HGetall(name string) ([]interface{}, error)
	hgetall, err := db.HGetall(hashset)
	assert.Equal(t, nil, err)
	assert.Equal(t, len(hkeys)*2, len(hgetall))
	for i := 0; i < len(hgetall); i += 2 {
		v := keyAndValues[string(hgetall[i])]
		switch v.(type) {
		case string:
			assert.Equal(t, 0, bytes.Compare([]byte(v.(string)), hgetall[i+1]))
		case []byte:
			assert.Equal(t, 0, bytes.Compare(v.([]byte), hgetall[i+1]))
		default:
			vByte, err := json.Marshal(v)
			if err != nil {
				assert.Error(t, err)
			}
			assert.Equal(t, 0, bytes.Compare(vByte, hgetall[i+1]))
		}
	}
}
