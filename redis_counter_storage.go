package storage

import (
	"github.com/0studio/storage_key"
)

func (this RedisStorage) Incr(key key.Key, step uint64) (newValue uint64, err error) {
	keyCache, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return 0, err
	}
	result, errcache := this.client.Incr(keyCache, step)
	return uint64(result), errcache
}

func (this RedisStorage) Decr(key key.Key, step uint64) (newValue uint64, err error) {
	keyCache, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return 0, err
	}
	result, errcache := this.client.Decr(keyCache, step)
	if result < 0 {
		return 0, err
	}
	return uint64(result), errcache
}
