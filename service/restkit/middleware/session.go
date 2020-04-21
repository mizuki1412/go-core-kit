package middleware

//var Session *scs.SessionManager
//
//func init()  {
//	Session = scs.New()
//	Session.IdleTimeout = 3*time.Hour
//	Session.Cookie.SameSite = http.SameSiteNoneMode
//}

// https://github.com/kataras/iris/wiki/Sessions-database

/**
	session超时过期而不是生命周期
	session id
 */
//func Session() context.Handler {
//	return func(ctx *context.Context) {
//		ctx.Session = sessionManager.Start(ctx.Proxy)
//		logkit.Info(ctx.Session.ID())
//		// todo 超期处理
//		ctx.Proxy.Next()
//	}
//}