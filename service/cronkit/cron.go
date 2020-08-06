package cronkit

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cast"
)

var scheduler *cron.Cron

//
var pool map[string]*cron.Cron

func Scheduler() *cron.Cron {
	if scheduler == nil {
		scheduler = cron.New(cron.WithSeconds(), cron.WithLocation(configkit.GetLocation()))
	}
	return scheduler
}

func AddPool(key string, cron *cron.Cron) {
	RemovePool(key)
	pool[key] = cron
}

func RemovePool(key string) {
	v, ok := pool[key]
	if ok {
		v.Stop()
		delete(pool, key)
	}
}

// 给默认的scheduler add func， 封装上recover
func AddFunc(spec string, fun func()) {
	_, err := Scheduler().AddFunc(spec, func() {
		defer func() {
			if err := recover(); err != nil {
				var msg string
				if e, ok := err.(exception.Exception); ok {
					//msg = e.Msg
					// 带代码位置信息
					logkit.Error(e.Error())
				} else {
					msg = cast.ToString(err)
					logkit.Error(msg)
				}
			}
		}()
		fun()
	})
	if err != nil {
		panic(exception.New(err.Error()))
	}
}
