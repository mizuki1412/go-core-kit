package middleware

import (
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
)

func Cors() router.Handler {
	return func(c *context.Context) {
		//method := c.Request.Method
		c.Proxy.Header("Access-Control-Allow-Origin", c.Proxy.Request.Header.Get("Origin"))
		c.Proxy.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Accept, Token")
		c.Proxy.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Proxy.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Proxy.Header("Access-Control-Allow-Credentials", "true")
		//放行所有OPTIONS方法
		//if method == "OPTIONS" {
		//	c.Proxy.StatusCode(http.StatusNoContent)
		//	c.Proxy.StopExecution()
		//}
		//c.Proxy.SetCookie(&http.Cookie{
		//	Name:     "Set-Cookie",
		//	SameSite: http.SameSiteNoneMode,
		//})
		c.Proxy.Next()
	}
}
