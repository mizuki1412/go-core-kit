package context

import (
	"github.com/kataras/iris/v12"
	"mizuki/project/core-kit/service/restkit/swagger"
	"net/http"
)

/**
router的抽象
*/

type Router struct {
	Proxy      *iris.Application
	IsGroup    bool
	ProxyGroup iris.Party
	Path       string
}
type Handler func(ctx *Context)

func handlerTrans(handlers ...Handler) []iris.Handler {
	list := make([]iris.Handler, len(handlers), len(handlers))
	for i, v := range handlers {
		list[i] = func(ctx iris.Context) {
			// 实际ctx进入，转为抽象层的context todo 注意field更新
			v(&Context{
				Proxy:    ctx,
				Request:  ctx.Request(),
				Response: ctx.ResponseWriter(),
			})
		}
	}
	return list
}

func (router *Router) Group(path string, handlers ...Handler) *Router {
	var r iris.Party
	if router.IsGroup {
		r = router.ProxyGroup.Party(path, handlerTrans(handlers...)...)
	} else {
		r = router.Proxy.Party(path, handlerTrans(handlers...)...)
	}
	return &Router{
		IsGroup:    true,
		ProxyGroup: r,
		Path:       router.Path + path,
	}
}

func (router *Router) Use(handlers ...Handler) {
	if router.IsGroup {
		router.ProxyGroup.Use(handlerTrans(handlers...)...)
	} else {
		router.Proxy.Use(handlerTrans(handlers...)...)
	}
}

func (router *Router) Post(path string, handlers ...Handler) *swagger.SwaggerPath {
	if router.IsGroup {
		router.ProxyGroup.Post(path, handlerTrans(handlers...)...)
	} else {
		router.Proxy.Post(path, handlerTrans(handlers...)...)
	}
	return &swagger.SwaggerPath{Path: router.Path + path}
}
func (router *Router) Get(path string, handlers ...Handler) *swagger.SwaggerPath {
	if router.IsGroup {
		router.ProxyGroup.Get(path, handlerTrans(handlers...)...)
	} else {
		router.Proxy.Get(path, handlerTrans(handlers...)...)
	}
	return &swagger.SwaggerPath{Path: router.Path + path}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	router.Proxy.ServeHTTP(w, req)
}
