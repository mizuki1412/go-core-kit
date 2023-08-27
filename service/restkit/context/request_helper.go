package context

import (
	"github.com/mizuki1412/go-core-kit/library/commonkit"
	"github.com/mizuki1412/go-core-kit/service/jwtkit"
)

var HeaderTokenKey = "Authorization"
var CookieTokenKey = "token"

func (ctx *Context) ReadToken() {
	token := ctx.Request.Header.Get(HeaderTokenKey)
	if token == "" {
		// 从cookie中获取
		token, _ = ctx.Proxy.Cookie(CookieTokenKey)
	}
	if token != "" {
		_ = commonkit.RecoverFuncWrapper(func() {
			c := jwtkit.Parse(token)
			ctx.Set("jwt", c)
			ctx.Set("jwt-token", token)
		})
	}
}
