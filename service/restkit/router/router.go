package router

import (
	"embed"
	"github.com/kataras/iris/v12"
	context2 "github.com/kataras/iris/v12/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	swg "github.com/mizuki1412/go-core-kit/service/restkit/swagger"
	"mime"
	"net/http"
	"path"
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
func EmbedHtmlHandle(fs embed.FS, root string) func(c context2.Context) {
	return func(c context2.Context) {
		// 解析访问路径
		pathName := c.Params().Get("path")
		if pathName == "" {
			pathName = "index.html"
		}
		assetPath := path.Join(root, pathName)
		assets, err := fs.Open(assetPath)
		if err != nil {
			_, _ = c.Write([]byte(err.Error()))
			return
		}
		//f, err := pkger.Open(pkPath + pathName)
		//if err != nil {
		//	_, _ = c.Write([]byte(err.Error()))
		//	return
		//}
		data := make([]byte, 0, 1024*5)
		for true {
			temp := make([]byte, 1024)
			n, _ := assets.Read(temp)
			if n == 0 {
				break
			} else {
				data = append(data, temp[:n]...)
			}
		}
		//_ = mime.AddExtensionType(".js", "text/javascript")
		// mine
		i := strings.LastIndex(pathName, ".")
		if i > 0 {
			c.ContentType(mime.TypeByExtension(pathName[i:]))
		}
		_, _ = c.Write(data)
	}
}

// todo 需要外部工程在rest run之前，指定此处的值
var SwaggerAssets embed.FS

func (router *Router) RegisterSwagger() {
	if router.IsGroup {
		//router.ProxyGroup.Get("/swagger/{any:path}", swagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))
		router.ProxyGroup.Get("/swagger/doc", func(c context2.Context) {
			_, _ = c.Write([]byte(swg.Doc.ReadDoc()))
		})
		// swagger-ui 需要被pkger打包，第二个path表示匹配路径
		router.ProxyGroup.Get("/swagger/{path:path}", EmbedHtmlHandle(SwaggerAssets, "./swagger-ui"))
		router.ProxyGroup.Get("/swagger", EmbedHtmlHandle(SwaggerAssets, "./swagger-ui"))
	} else {
		router.Proxy.Get("/swagger/doc", func(c context2.Context) {
			_, _ = c.Write([]byte(swg.Doc.ReadDoc()))
		})
		//router.Proxy.HandleDir("/swagger", "./swagger-ui")
		router.Proxy.Get("/swagger/{path:path}", EmbedHtmlHandle(SwaggerAssets, "./swagger-ui"))
		router.Proxy.Get("/swagger", EmbedHtmlHandle(SwaggerAssets, "./swagger-ui"))
	}
	//swag.Register(swag.Name, &swg.Doc)
}
