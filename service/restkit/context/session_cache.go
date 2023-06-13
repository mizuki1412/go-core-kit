package context

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/init/configkey"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/cachekit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
	"github.com/spf13/cast"
	"strings"
	"time"
)

// SessionSetUser session每次请求时都会从redis中获取，所以在session中存储的务必是string，如果是对象，会被自动转为json，但如果其中有unicode，可能造成指数增加
func (ctx *Context) SessionSetUser(user any) {
	if ctx.SessionToken() == "" {
		panic(exception.New("session token is nil"))
	}
	if user == nil {
		return
	}
	ctx._set("session-user-"+ctx.SessionToken(), user)
}
func (ctx *Context) SessionSetSchema(schema string) {
	if ctx.SessionToken() == "" {
		panic(exception.New("session token is nil"))
	}
	ctx._set("session-schema-"+ctx.SessionToken(), schema)
}

// todo 管理session的key
//func (ctx *Context) SessionSet(key string, val any) {
//	if ctx.SessionToken() == "" {
//		panic(exception.New("session token is nil"))
//	}
//	cachekit.Set("session-"+key+"-"+ctx.SessionToken(), val, &cachekit.Param{Redis: true, Ttl: time.Duration(configkit.GetInt(configkey.SessionExpire, 12)) * time.Hour})
//}

// SessionGetUser return eg *model.User
func (ctx *Context) SessionGetUser() *model.User {
	if ctx.SessionToken() == "" {
		return nil
	}
	json := cachekit.Get("session-user-"+ctx.SessionToken(), &cachekit.Param{Redis: true})
	if j, ok := json.(string); ok {
		user := &model.User{}
		err := jsonkit.ParseObj(j, user)
		if err != nil {
			logkit.Error(err.Error())
			return nil
		}
		return user
	} else {
		return nil
	}
}
func (ctx *Context) SessionGetUserOrigin() any {
	if ctx.SessionToken() == "" {
		return nil
	}
	return cachekit.Get("session-user-"+ctx.SessionToken(), &cachekit.Param{Redis: true})
}
func (ctx *Context) SessionGetSchema() string {
	r := cast.ToString(cachekit.Get("session-schema-"+ctx.SessionToken(), &cachekit.Param{Redis: true}))
	if r == "" {
		r = sqlkit.SchemaDefault
	}
	return r
}

func (ctx *Context) SessionClear() {
	cachekit.Del("session-user-"+ctx.SessionToken(), &cachekit.Param{Redis: true})
	cachekit.Del("session-schema-"+ctx.SessionToken(), &cachekit.Param{Redis: true})
	// todo session的其他key？
}

// 优化renew的频次
var renewCache = map[string]int64{}

func (ctx *Context) SessionRenew() {
	token := ctx.SessionToken()
	v := renewCache[token]
	if (time.Now().Unix() - v) > (cast.ToInt64(configkit.GetInt(configkey.SessionExpire, 12)) / 2 * 3600) {
		ctx._renew("session-user-" + token)
		ctx._renew("session-schema-" + token)
	}
}

func (ctx *Context) SessionToken() string {
	return cast.ToString(ctx.Get("_token"))
}

func (ctx *Context) _set(key string, val any) {
	cachekit.Set(key, val, &cachekit.Param{Redis: true, Ttl: time.Duration(configkit.GetInt(configkey.SessionExpire, 12)) * time.Hour})
	ctx._setCookie()
}
func (ctx *Context) _renew(key string) {
	v := renewCache[key]
	now := time.Now().Unix()
	if (now - v) > (cast.ToInt64(configkit.GetInt(configkey.SessionExpire, 12)) / 2 * 3600) {
		cachekit.Renew(key, &cachekit.Param{Redis: true, Ttl: time.Duration(configkit.GetInt(configkey.SessionExpire, 12)) * time.Hour})
		renewCache[key] = now
		ctx._setCookie()
	}
}
func (ctx *Context) _setCookie() {
	if ctx.Response.Header().Get("set-cookie") == "" {
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
		ctx.Proxy.SetCookie("token", ctx.SessionToken(), configkit.GetInt(configkey.SessionExpire, 12)*3600, "/", origin, false, true)
	}
}
