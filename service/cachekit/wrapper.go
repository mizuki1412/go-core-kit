package cachekit

import (
	"github.com/mizuki1412/go-core-kit/v2/cli/configkey"
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
	"time"
)

type WrapParam struct {
	Key string
	Ttl time.Duration
}

const defaultKeyPrefix = "__cache_"

func Wrapper(wp WrapParam, f func() any) any {
	if wp.Ttl == 0 {
		wp.Ttl = time.Duration(configkit.GetInt(configkey.CacheWrapperTTL)) * time.Second
	}
	value := Get(defaultKeyPrefix + wp.Key)
	if value == nil {
		value = f()
		Set(defaultKeyPrefix+wp.Key, value, &Param{Ttl: wp.Ttl})
	}
	return value
}

/**
 sample:
	cachekit.Wrapper(cachekit.WrapParam{
		Key: "abc",
		Ttl: 0,
	}, func() any {
		dao := userdao.New(userdao.ResultDefault)
		return dao.List(userdao.ListParam{})
	}).([]*model.User)
*/
