package context

import (
	"github.com/mizuki1412/go-core-kit/service/jwtkit"
	"github.com/spf13/cast"
	"strings"
	"time"
)

func (ctx *Context) SetJwtCookie(c jwtkit.Claims, token string) {
	//ctx.Proxy.SetSameSite(http.SameSiteNoneMode)
	origin := ctx.Proxy.GetHeader("origin")
	origins := strings.Split(origin, "//")
	if len(origins) > 1 {
		origin = origins[1]
	}
	// 可能域名是省略www的，但是origin有; 此时浏览器还是会当成不同的。所以尽量不省略。
	//if strings.Index(origin, "www") == 0 {
	//	origin = origin[3:]
	//}
	if c.ExpiresAt != nil {
		ctx.Proxy.SetCookie("token", token, cast.ToInt(c.ExpiresAt.Unix()-time.Now().Unix()), "/", origin, false, true)
	}
}

func (ctx *Context) GetJwt() jwtkit.Claims {
	if ctx.Get("jwt") == nil {
		ctx.ReadToken()
	}
	if c := ctx.Get("jwt"); c != nil {
		if cc, ok := c.(jwtkit.Claims); ok {
			return cc
		}
	}
	return jwtkit.Claims{}
}
