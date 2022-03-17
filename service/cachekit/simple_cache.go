package cachekit

import (
	"github.com/mizuki1412/go-core-kit/class"
	"sync"
	"time"
)

// 内部简单cache，基于MapStringSync

var simpleCacheHandler *SimpleCacheHandler
var once sync.Once

type simpleCacheBean struct {
	Expire time.Time
	Data   any
}

type SimpleCacheHandler struct {
	Cache *class.MapStringSync
}

// 启用服务
func SimpleCache() *SimpleCacheHandler {
	if simpleCacheHandler == nil {
		once.Do(func() {
			simpleCacheHandler = &SimpleCacheHandler{Cache: &class.MapStringSync{Valid: true}}
			// todo clear
			//cronkit.AddFunc("@every 1h",func() {
			//	simpleCache.RLock()
			//	defer simpleCache.RUnlock()
			//	now := time.Now()
			//	for k,v:=range simpleCache.Map{
			//		v = v.(*simpleCacheBean)
			//		if
			//	}
			//})
		})
	}
	return simpleCacheHandler
}

func (th *SimpleCacheHandler) Put(key string, val any, expireSeconds int32) {
	th.Cache.Put(key, &simpleCacheBean{
		Expire: time.Now().Add(time.Duration(expireSeconds) * time.Second),
		Data:   val,
	})
}

func (th *SimpleCacheHandler) Get(key string) any {
	val := th.Cache.Get(key)
	if val != nil {
		if val.(*simpleCacheBean).Expire.Before(time.Now()) {
			return nil
		}
		return val.(*simpleCacheBean).Data
	}
	return nil
}
