package router

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/mizuki1412/go-core-kit/init/httpconst"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	swg "github.com/mizuki1412/go-core-kit/service/restkit/swagger"
	"mime"
	"net/http"
	"path"
	"strings"
)

// Router router的抽象
type Router struct {
	Proxy      *gin.Engine
	Base       string
	ProxyGroup *gin.RouterGroup
	Swagger    *swg.SwaggerPath
}
type Handler func(ctx *context.Context)

func handlerTrans(handlers ...Handler) []gin.HandlerFunc {
	list := make([]gin.HandlerFunc, len(handlers), len(handlers))
	for i, v := range handlers {
		list[i] = func(ctx *gin.Context) {
			// 实际ctx进入，转为抽象层的context
			v(&context.Context{
				Proxy:    ctx,
				Request:  ctx.Request,
				Response: ctx.Writer,
			})
		}
	}
	return list
}
func handlerTransOne(handler Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 实际ctx进入，转为抽象层的context
		handler(&context.Context{
			Proxy:    ctx,
			Request:  ctx.Request,
			Response: ctx.Writer,
		})
	}
}

func (router *Router) Group(path string, handlers ...Handler) *Router {
	r0 := &Router{
		Proxy:      router.Proxy,
		ProxyGroup: router.Proxy.Group(router.Base + path),
		Base:       router.Base + path,
	}
	if len(handlers) > 0 {
		r0.Use(handlers...)
	}
	return r0
}

func (router *Router) Use(handlers ...Handler) *Router {
	if router.ProxyGroup != nil {
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

//func (router *Router) OnError(handlers ...Handler) {
//	for _, v := range handlers {
//		router.Proxy.(handlerTransOne(v))
//	}
//}

func (router *Router) setPath4Swagger(path string, method string) {
	router.Swagger = swg.NewPath(router.Base+path, method)
}

// Post 此处handle不能当成是use
func (router *Router) Post(path string, handlers ...Handler) *Router {
	if router.ProxyGroup != nil {
		router.ProxyGroup.POST(path, handlerTrans(handlers...)...)
	} else {
		router.Proxy.POST(router.Base+path, handlerTrans(handlers...)...)
	}
	router.setPath4Swagger(path, "post")
	return router
}
func (router *Router) Get(path string, handlers ...Handler) *Router {
	if router.ProxyGroup != nil {
		router.ProxyGroup.GET(path, handlerTrans(handlers...)...)
	} else {
		router.Proxy.GET(router.Base+path, handlerTrans(handlers...)...)
	}
	router.setPath4Swagger(path, "get")
	return router
}

// GetOrigin 不附带router.base
func (router *Router) GetOrigin(path string, handlers ...Handler) *Router {
	router.Proxy.GET(path, handlerTrans(handlers...)...)
	return router
}
func (router *Router) Any(path string, handlers ...Handler) *Router {
	if router.ProxyGroup != nil {
		router.ProxyGroup.Any(path, handlerTrans(handlers...)...)
	} else {
		router.Proxy.Any(router.Base+path, handlerTrans(handlers...)...)
	}
	// todo any swagger
	router.setPath4Swagger(path, "get")
	return router
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	router.Proxy.ServeHTTP(w, req)
}

// EmbedHtmlHandle 注意path pattern中加入{path:path}
// url中path的路径前缀需要和root一致
func EmbedHtmlHandle(fs embed.FS, root string) func(c *context.Context) {
	return func(c *context.Context) {
		// 解析访问路径
		var assetPath string
		pathName := c.Proxy.Param("action")
		if pathName == "" || pathName == "/" {
			pathName = "index.html"
		}
		assetPath = path.Join(root, pathName)
		assets, err := fs.Open(assetPath)
		if err != nil {
			c.Proxy.Status(http.StatusBadRequest)
			_, _ = c.Proxy.Writer.Write([]byte(err.Error()))
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
			c.Proxy.Render(http.StatusOK, render.Data{Data: data, ContentType: mime.TypeByExtension(pathName[i:])})
		} else {
			c.Proxy.Render(http.StatusOK, render.Data{Data: data})
		}
	}
}

func (router *Router) RegisterSwagger() {
	router.Get("/swagger/doc", func(c *context.Context) {
		c.Proxy.Render(http.StatusOK, render.Data{Data: []byte(swg.Doc.ReadDoc()), ContentType: httpconst.ContentTypeJSON})
		//c.Proxy.Status(http.StatusOK)
		//_, _ = c.Proxy.Writer.Write([]byte(swg.Doc.ReadDoc()))
	})
	// 第二个path表示匹配路径
	router.Get("/swagger/:path", EmbedHtmlHandle(swg.UiAssets, "./swagger-ui"))
	router.Get("/swagger", EmbedHtmlHandle(swg.UiAssets, "./swagger-ui"))
}
