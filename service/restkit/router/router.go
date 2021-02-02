package router

import (
	"github.com/kataras/iris/v12"
	context2 "github.com/kataras/iris/v12/context"
	"github.com/markbates/pkger"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	swg "github.com/mizuki1412/go-core-kit/service/restkit/swagger"
	"mime"
	"net/http"
	"strings"
)

/**
router的抽象
*/

type Router struct {
	Proxy      *iris.Application
	IsGroup    bool
	ProxyGroup iris.Party // 存在项目前缀时，base path
	Path       string
	Swagger    *swg.SwaggerPath
}
type Handler func(ctx *context.Context)

func handlerTrans(handlers ...Handler) []iris.Handler {
	list := make([]iris.Handler, len(handlers), len(handlers))
	for i, v := range handlers {
		list[i] = func(ctx iris.Context) {
			// 实际ctx进入，转为抽象层的context todo 注意field更新
			v(&context.Context{
				Proxy:    ctx,
				Request:  ctx.Request(),
				Response: ctx.ResponseWriter(),
			})
		}
	}
	return list
}
func handlerTransOne(handler Handler) iris.Handler {
	return func(ctx iris.Context) {
		// 实际ctx进入，转为抽象层的context todo 注意field更新
		handler(&context.Context{
			Proxy:    ctx,
			Request:  ctx.Request(),
			Response: ctx.ResponseWriter(),
		})
	}
}

func (router *Router) Group(path string, handlers ...Handler) *Router {
	var r iris.Party
	if router.IsGroup {
		r = router.ProxyGroup.Party(path)
	} else {
		r = router.Proxy.Party(path)
	}
	r0 := &Router{
		IsGroup:    true,
		ProxyGroup: r,
		Path:       router.Path + path,
	}
	if len(handlers) > 0 {
		r0.Use(handlers...)
	}
	return r0
}

func (router *Router) Use(handlers ...Handler) *Router {
	// ？多参数handlers会多次执行最后一个handle？
	if router.IsGroup {
		for _, v := range handlers {
			router.ProxyGroup.Use(handlerTransOne(v))
		}
	} else {
		for _, v := range handlers {
			router.Proxy.Use(handlerTransOne(v))
		}
	}
	return router
}
func (router *Router) OnError(handlers ...Handler) {
	for _, v := range handlers {
		router.Proxy.OnAnyErrorCode(handlerTransOne(v))
	}
}

// 此处handle不能当成是use
func (router *Router) Post(path string, handlers ...Handler) *Router {
	if router.IsGroup {
		router.ProxyGroup.Post(path, handlerTrans(handlers...)...)
	} else {
		router.Proxy.Post(path, handlerTrans(handlers...)...)
	}
	router.Swagger = swg.NewPath(router.Path+path, "post")
	return router
}
func (router *Router) Get(path string, handlers ...Handler) *Router {
	if router.IsGroup {
		router.ProxyGroup.Get(path, handlerTrans(handlers...)...)
	} else {
		router.Proxy.Get(path, handlerTrans(handlers...)...)
	}
	router.Swagger = swg.NewPath(router.Path+path, "get")
	return router
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	router.Proxy.ServeHTTP(w, req)
}

// 用于pkger打包资源的html访问设置
// 注意path pattern中加入{path:path}
func EmbedHtmlHandle(pkPath string) func(c context2.Context) {
	return func(c context2.Context) {
		p := c.Params().Get("path")
		if p == "" {
			p = "index.html"
		}
		f, err := pkger.Open(pkPath + p)
		if err != nil {
			_, _ = c.Write([]byte(err.Error()))
			return
		}
		data := make([]byte, 0, 1024*5)
		for true {
			temp := make([]byte, 1024)
			n, _ := f.Read(temp)
			if n == 0 {
				break
			} else {
				data = append(data, temp[:n]...)
			}
		}
		//_ = mime.AddExtensionType(".js", "text/javascript")
		// mine
		i := strings.LastIndex(p, ".")
		if i > 0 {
			c.ContentType(mime.TypeByExtension(p[i:]))
		}
		_, _ = c.Write(data)
	}
}

func (router *Router) RegisterSwagger() {
	if router.IsGroup {
		//router.ProxyGroup.Get("/swagger/{any:path}", swagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))
		router.ProxyGroup.Get("/swagger/doc", func(c context2.Context) {
			_, _ = c.Write([]byte(swg.Doc.ReadDoc()))
		})
		// swagger-ui 需要被pkger打包，第二个path表示匹配路径
		router.ProxyGroup.Get("/swagger/{path:path}", EmbedHtmlHandle("/swagger-ui/"))
		router.ProxyGroup.Get("/swagger", EmbedHtmlHandle("/swagger-ui/"))
	} else {
		router.Proxy.Get("/swagger/doc", func(c context2.Context) {
			_, _ = c.Write([]byte(swg.Doc.ReadDoc()))
		})
		//router.Proxy.HandleDir("/swagger", "./swagger-ui")
		router.Proxy.Get("/swagger/{path:path}", EmbedHtmlHandle("/swagger-ui/"))
		router.Proxy.Get("/swagger", EmbedHtmlHandle("/swagger-ui/"))
	}
	//swag.Register(swag.Name, &swg.Doc)
}
