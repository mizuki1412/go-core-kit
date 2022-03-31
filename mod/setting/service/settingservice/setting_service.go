package settingservice

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/mod/setting/dao/settingdao"
	"sync"
)

var cache = &class.MapStringSync{}
var _once sync.Once

func Get(schema string) *class.MapStringSync {
	_once.Do(func() {
		_cache := settingdao.New(schema).Get()
		cache.Set(_cache)
	})
	return cache
}

// Sync 先get后set
func Sync(schema string) {
	settingdao.New(schema).Set(cache.Map)
}
