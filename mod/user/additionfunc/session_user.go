package additionfunc

import (
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
)

// SessionGetUserFunc session定制配置, 转换user对象
// Deprecated
func SessionGetUserFunc() {
	context.SessionGetUserFunc = func(ctx *context.Context) interface{} {
		json := ctx.Session().Get("user")
		if j, ok := json.(string); ok {
			user := &model.User{}
			// todo 可能存在性能问题：几十毫秒
			err := jsonkit.ParseObj(j, user)
			if err != nil {
				return nil
			}
			return user
		} else if _, ok := ctx.Session().Get("user").(*model.User); ok {
			return ctx.Session().Get("user")
		} else {
			return nil
		}
	}
}
