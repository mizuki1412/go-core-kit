package context

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/sessions"
	"mizuki/project/core-kit/class/exception"
	"net/http"
)

type Context struct {
	Proxy    iris.Context
	Request  *http.Request
	Response context.ResponseWriter
}

var sessionManager *sessions.Sessions

func init() {
	sessionManager = sessions.New(sessions.Config{
		Cookie:       "session",
		AllowReclaim: true,
	})
}

func (ctx *Context) Session() *sessions.Session {
	return sessionManager.Start(ctx.Proxy)
}

func (ctx *Context) RenewSession() *sessions.Session {
	sess := ctx.Session()
	if !sess.IsNew() {
		sessionManager.Destroy(ctx.Proxy)
		return ctx.Session()
	}
	return sess
}

// data: query, form, json/xml, param

func (ctx *Context) BindForm(bean interface{}) {
	//ctx.Proxy.Params().Get("demo")
	switch bean.(type) {
	case *map[string]interface{}:
		// query会和form合并 post时
		//allQuery := ctx.Proxy.URLParams()
		allForm := ctx.Proxy.FormValues()
		for k, v := range allForm {
			(*(bean.(*map[string]interface{})))[k] = v[len(v)-1]
		}
	default:
		err := ctx.Proxy.ReadForm(bean)
		if err != nil {
			panic(exception.New("form解析错误"))
		}
	}
	//_ = ctx.Proxy.ReadJSON(bean)
}
