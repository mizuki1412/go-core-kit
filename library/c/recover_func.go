package c

import (
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/service/logkit"
	"github.com/spf13/cast"
)

// RecoverFuncWrapper 将函数套在recover内，实现exception catch。如果捕获异常
func RecoverFuncWrapper(fun func()) (re *exception.Exception) {
	defer func() {
		if err := recover(); err != nil {
			var msg string
			if e, ok := err.(exception.Exception); ok {
				msg = e.Msg
				// 带代码位置信息
				logkit.ErrorException(e)
				re = &e
			} else {
				msg = cast.ToString(err)
				exp := exception.New(msg, 3)
				logkit.ErrorException(exp)
				re = &exp
			}
		}
	}()
	fun()
	return nil
}

func RecoverGoFuncWrapper(fun func()) {
	go func() {
		_ = RecoverFuncWrapper(fun)
	}()
}
