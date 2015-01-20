package storage

import (
	"github.com/0studio/storage_key"
	"github.com/golang/groupcache"
)

type GroupCacheKvStorage struct {
	CacheGroup        *groupcache.Group
	DefaultExpireTime int
	encoding          Encoding
}

func (this GroupCacheKvStorage) Get(key key.Key) (interface{}, error) {
	var data []byte
	this.CacheGroup.Get(nil, key.ToString(), groupcache.AllocatingByteSliceSink(&data))
	object, err := this.encoding.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (this GroupCacheKvStorage) Set(key key.Key, object interface{}) error {
	return nil
}

func (this GroupCacheKvStorage) MultiGet(keys []key.Key) (map[key.Key]interface{}, error) {
	resultMap := make(map[key.Key]interface{})
	for _, key := range keys {
		value, err := this.Get(key)
		if err != nil {
			continue
		}
		resultMap[key] = value
	}
	return resultMap, nil
}

func (this GroupCacheKvStorage) MultiSet(objectMap map[key.Key]interface{}) error {
	return nil
}

func (this GroupCacheKvStorage) Delete(key key.Key) error {
	return nil
}

func (this GroupCacheKvStorage) FlushAll() {
}
