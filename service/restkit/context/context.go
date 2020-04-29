package context

import (
	"fmt"
	"github.com/go-playground/validator/v10"
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
		err = Validator.Struct(bean)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				panic(exception.New(err.Error()))
			}
			for _, err0 := range err.(validator.ValidationErrors) {
				// todo 格式化err
				panic(exception.New(fmt.Sprintf("%v", err0)))
			}
		}
	}
	//_ = ctx.Proxy.ReadJSON(bean)
}
