package settingservice

import (
	"github.com/mizuki1412/go-core-kit/v2/class"
	"github.com/mizuki1412/go-core-kit/v2/mod/setting/dao/settingdao"
	"sync"
)

var cache = &class.MapStringSync{}
var _once sync.Once

func Get(schema ...string) *class.MapStringSync {
	_once.Do(func() {
		dao := settingdao.New()
		if len(schema) > 0 {
			dao.DataSource().Schema = schema[0]
		}
		_cache := dao.Get()
		cache.Set(_cache)
	})
	return cache
}

// Sync 先get后set
func Sync(schema ...string) {
	dao := settingdao.New()
	if len(schema) > 0 {
		dao.DataSource().Schema = schema[0]
	}
	dao.Set(Get(schema...).Map)
}
