package context

func GetTokenFromReq(ctx *Context) string {
	token := ctx.Request.Header.Get("token")
	if token == "" {
		// 从cookie中获取
		token, _ = ctx.Proxy.Cookie("token")
	}
	return token
}
