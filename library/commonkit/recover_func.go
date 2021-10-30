package commonkit

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cast"
)

// RecoverFuncWrapper 将函数套在recover内，实现exception catch。eg：用在for-range时
func RecoverFuncWrapper(fun func(), throwAgain ...bool) {
	defer func() {
		if err := recover(); err != nil {
			var msg string
			if e, ok := err.(exception.Exception); ok {
				msg = e.Msg
				// 带代码位置信息
				logkit.Error(e.Error())
			} else {
				msg = cast.ToString(err)
				logkit.Error(exception.New(msg, 3).Error())
			}
			// 重新抛出
			if len(throwAgain) > 0 && throwAgain[0] {
				panic(err)
			}
		}
	}()
	fun()
}

func RecoverGoFuncWrapper(fun func(), throwAgain ...bool) {
	go func() {
		RecoverFuncWrapper(fun, throwAgain...)
	}()
}
