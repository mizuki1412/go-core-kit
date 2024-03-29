package cachekit

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/library/c"
	"github.com/mizuki1412/go-core-kit/v2/service/rediskit"
	"sync"
	"time"
)

var _cache *ristretto.Cache
var once sync.Once

func _getCache() {
	if _cache == nil {
		once.Do(func() {
			cache, err := ristretto.NewCache(&ristretto.Config{
				NumCounters: 1e7,     // number of keys to track frequency of (10M).
				MaxCost:     1 << 30, // maximum cost of cache (1GB).
				BufferItems: 64,      // number of keys per Get buffer.
			})
			if err != nil {
				panic(exception.New(err.Error()))
			}
			_cache = cache
		})
	}
}

type Param struct {
	Ttl  time.Duration
	Cost int64
	// 如果存在redis配置，将从redis操作，本地cache只是第二顺序处理
	// 控制忽略redis
	IgnoreRedis bool
}

func _handleParam(ps []*Param) *Param {
	_getCache()
	var p *Param
	if len(ps) == 0 {
		p = nil
	} else {
		p = ps[0]
	}
	if p == nil {
		p = &Param{}
	}
	return p
}

func Set(key string, value any, ps ...*Param) {
	p := _handleParam(ps)
	if rediskit.HasConfig() && !p.IgnoreRedis {
		_ = c.RecoverFuncWrapper(func() {
			rediskit.Set(context.Background(), rediskit.GetKeyWithPrefix(key), value, p.Ttl)
		})
	}
	// 同时也存入cache
	var res bool
	if p.Ttl > 0 {
		res = _cache.SetWithTTL(key, value, p.Cost, p.Ttl)
	} else {
		res = _cache.Set(key, value, p.Cost)
	}
	if !res {
		panic(exception.New("cache failed: " + key))
	}
}

func Get(key string, ps ...*Param) any {
	p := _handleParam(ps)
	var r any = nil
	if rediskit.HasConfig() && !p.IgnoreRedis {
		r0 := rediskit.Get(context.Background(), rediskit.GetKeyWithPrefix(key), "")
		if r0 != "" {
			r = r0
		}
	}
	if r == nil {
		r, _ = _cache.Get(key)
	}
	return r
}

func Del(key string, ps ...*Param) {
	p := _handleParam(ps)
	if rediskit.HasConfig() && !p.IgnoreRedis {
		rediskit.Del(context.Background(), rediskit.GetKeyWithPrefix(key))
	}
	_cache.Del(key)
}

func Renew(key string, ps ...*Param) {
	p := _handleParam(ps)
	if p.Ttl <= 0 {
		return
	}
	if rediskit.HasConfig() && !p.IgnoreRedis {
		rediskit.Expire(context.Background(), rediskit.GetKeyWithPrefix(key), p.Ttl)
	}
	v, ok := _cache.Get(key)
	if ok {
		_cache.SetWithTTL(key, v, p.Cost, p.Ttl)
	}
}
