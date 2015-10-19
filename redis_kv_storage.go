package storage

import (
	"errors"
	"github.com/0studio/redisapi"
	"github.com/0studio/storage_key"
)

type RedisStorage struct {
	client            redisapi.Redis
	KeyPrefix         string
	DefaultExpireTime int
	encoding          Encoding
}

func NewRedisStorage(serverUrl string, keyPrefix string, defaultExpireTime int, encoding Encoding) (RedisStorage, error) {
	client, err := redisapi.InitDefaultClient(serverUrl)
	return RedisStorage{client, keyPrefix, defaultExpireTime, encoding}, err
}

func (this RedisStorage) Get(key key.Key) (interface{}, error) {
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return nil, err
	}
	data, err := this.client.Get(cacheKey)
	if err != nil || data == nil {
		return nil, err
	}
	object, err := this.encoding.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return object, nil

}

func (this RedisStorage) Set(key key.Key, object interface{}) error {
	buf, err := this.encoding.Marshal(object)
	if err != nil {
		return err
	}
	keyCache, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return err
	}
	this.client.Set(keyCache, buf)
	return nil
}

func (this RedisStorage) Add(key key.Key, object interface{}) error {
	return this.Set(key, object)
}

func (this RedisStorage) MultiGet(keys []key.Key) (map[key.Key]interface{}, error) {
	cacheKeys := make([]interface{}, len(keys))
	for index, key := range keys {
		cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
		if err != nil {
			return nil, err
		}
		cacheKeys[index] = cacheKey
	}
	values, err := this.client.MultiGet(cacheKeys)
	if err != nil {
		return nil, err
	}
	result := make(map[key.Key]interface{})
	for i, value := range values {
		if value == nil {
			continue
		}
		object, err := this.encoding.Unmarshal(value.([]byte))
		if err != nil {
			continue
		}
		result[keys[i]] = object
	}
	return result, nil
}

func (this RedisStorage) MultiSet(valueMap map[key.Key]interface{}) error {
	tempMap := make(map[string][]byte)
	for key, value := range valueMap {
		buf, err := this.encoding.Marshal(value)
		if err != nil {
			continue
		}
		cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
		if err != nil {
			continue
		}
		tempMap[cacheKey] = buf
	}
	return this.client.MultiSet(tempMap)
}

func (this RedisStorage) Delete(key key.Key) error {
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return err
	}
	return this.client.Delete(cacheKey)
}

func (this RedisStorage) GetKeyList(k key.Key) ([]key.Key, error) {
	return []key.Key{}, errors.New("unrealized")
}

func (this RedisStorage) PutKey(key key.Key, keys []key.Key) error {
	return errors.New("unrealized")
}

func (this RedisStorage) FlushAll() {
	this.client.ClearAll()
}
