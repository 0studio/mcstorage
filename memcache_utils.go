package storage

import (
	key "github.com/0studio/storage_key"
	jump "github.com/dgryski/go-jump"
	"github.com/dropbox/godropbox/memcache"
	"github.com/dropbox/godropbox/net2"
	"time"
)

type MemcacheConfig struct {
	AddrList             []string `json:"addr,omitempty"` // list of "ip:port"
	MaxActiveConnections int32    `json:"max_active_connections,omitempty"`
	MaxIdleConnections   uint32   `json:"max_idle_connections,omitempty"`
	ReadTimeOutMS        int      `json:"read_timeout_ms,omitempty"`
	WriteTimeOutMS       int      `json:"read_timeout_ms,omitempty"`
}

func GetClient(config MemcacheConfig, logError func(error), logInfo func(v ...interface{})) (mc memcache.Client) {
	if len(config.AddrList) == 0 {
		panic("could not load mc setting,mcAddrList len 0")
	}

	return getClientFromShardPool(config, nil, logError, logInfo)
}
func GetClient2(config MemcacheConfig, shardFunc func(key string, numShard int) (ret int), logError func(error), logInfo func(v ...interface{})) (mc memcache.Client) {
	if len(config.AddrList) == 0 {
		panic("could not load mc setting,mcAddrList len 0")
	}

	return getClientFromShardPool(config, shardFunc, logError, logInfo)
}

func getClientFromShardPool(config MemcacheConfig, shardFunc func(key string, numShard int) (ret int), logError func(error), logInfo func(v ...interface{})) (mc memcache.Client) {
	options := net2.ConnectionOptions{
		MaxActiveConnections: config.MaxActiveConnections,
		MaxIdleConnections:   config.MaxIdleConnections,
		ReadTimeout:          time.Duration(config.ReadTimeOutMS) * time.Millisecond,
		WriteTimeout:         time.Duration(config.WriteTimeOutMS) * time.Millisecond,
	}

	if shardFunc == nil {
		shardFunc = func(mcKey string, numShard int) (ret int) {
			if numShard == 0 {
				return -1
			}

			if numShard < 2 {
				return 0
			}
			// https://github.com/renstrom/go-jump-consistent-hash
			// jump 一致性hash 算法
			ret = int(jump.Hash(uint64(key.String(mcKey).ToSum()), numShard))
			// ret = int(crc32.ChecksumIEEE([]byte(key))) % len(mcAddrList)
			return
		}
	}

	manager := NewStaticShardManager(
		config.AddrList,
		logError,
		logInfo,
		shardFunc,
		options)
	mc = memcache.NewShardedClient(manager)

	return
}
func NewStaticShardManager(serverAddrs []string, logError func(error), logInfo func(v ...interface{}), shardFunc func(key string, numShard int) (shard int),
	options net2.ConnectionOptions) memcache.ShardManager {
	// 从dropbox/memcache/static_shard_manager.go copy 来
	// 将其中的log 换成zerogame.info/log

	manager := &memcache.StaticShardManager{}
	manager.Init(
		shardFunc,
		logError,
		logInfo,
		options)

	shardStates := make([]memcache.ShardState, len(serverAddrs), len(serverAddrs))
	for i, addr := range serverAddrs {
		shardStates[i].Address = addr
		shardStates[i].State = memcache.ActiveServer
	}

	manager.UpdateShardStates(shardStates)

	return manager
}
