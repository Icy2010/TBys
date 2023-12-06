package TBys

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var ch *cache.Cache

func InitCache() {
	ch = cache.New(60*time.Minute, 10*time.Minute)
	// 默认一小时的缓存
	// 10分钟清理一次
}

func Cache() *cache.Cache {
	return ch
}

func GetCache(key string) (Data interface{}, found bool) {
	Data, found = ch.Get(key)
	return
}

func HasCache(key string) bool {
	_, found := ch.Items()[key]
	return found
}
