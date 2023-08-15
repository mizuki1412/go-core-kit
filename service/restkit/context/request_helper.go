package context

var HeaderTokenKey = "token"

func GetTokenFromReq(ctx *Context) string {
	token := ctx.Request.Header.Get(HeaderTokenKey)
	if token == "" {
		// 从cookie中获取
		token, _ = ctx.Proxy.Cookie(HeaderTokenKey)
	}
	return token
}
