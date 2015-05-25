package storage

import (
	"github.com/dropbox/godropbox/memcache"
	"github.com/dropbox/godropbox/net2"
	"hash/crc32"
	"strings"
	"time"
)

func GetClient(mcSetting string, maxActiveConnCnt int32, maxIdleConnCnt uint32, readTimeout, writeTimeout time.Duration, logError func(error), logInfo func(v ...interface{})) (mc memcache.Client) {
	mcAddrList := strings.Split(mcSetting, ",")
	if len(mcAddrList) == 1 {
		return getClientFromShardPool(mcAddrList, maxActiveConnCnt, maxIdleConnCnt, readTimeout, writeTimeout, logError, logInfo)
		// getSingleClient(mcAddrList[0])
	} else if len(mcAddrList) == 0 { // 0 ,
		panic("could not load mc setting,mcAddrList len 0")
		return
	} else { // >1
		return getClientFromShardPool(mcAddrList, maxActiveConnCnt, maxIdleConnCnt, readTimeout, writeTimeout, logError, logInfo)
	}
	return
}

func getClientFromShardPool(mcAddrList []string, maxActiveConnCnt int32, maxIdleConnCnt uint32, readTimeout, writeTimeout time.Duration, logError func(error), logInfo func(v ...interface{})) (mc memcache.Client) {
	options := net2.ConnectionOptions{
		MaxActiveConnections: maxActiveConnCnt,
		MaxIdleConnections:   maxIdleConnCnt,
		ReadTimeout:          readTimeout,
		WriteTimeout:         writeTimeout,
	}

	manager := NewStaticShardManager(
		mcAddrList,
		logError,
		logInfo,
		func(key string, numShard int) (ret int) {
			if numShard == 0 {
				return -1
			}

			if numShard < 2 {
				return 0
			}
			ret = int(crc32.ChecksumIEEE([]byte(key))) % len(mcAddrList)
			return
		},
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
