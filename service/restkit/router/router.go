package router

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/openapi"
	"mime"
	"net/http"
	"path"
	"strings"
)

// Router router的抽象
type Router struct {
	Proxy      *gin.Engine
	ProxyGroup *gin.RouterGroup
	Openapi    *openapi.Builder
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
		ProxyGroup: router.ProxyGroup.Group(path),
		Openapi:    router.Openapi,
	}
	if len(handlers) > 0 {
		r0.Use(handlers...)
	}
	return r0
}

func (router *Router) Use(handlers ...Handler) *Router {
	for _, v := range handlers {
		router.ProxyGroup.Use(handlerTransOne(v))
	}
	return router
}

func (router *Router) openapiBuilder(path string, method string) {
	router.Openapi = openapi.NewBuilder(router.ProxyGroup.BasePath()+path, method)
}

// Post 此处handle不能当成是use
func (router *Router) Post(path string, handlers ...Handler) *Router {
	router.ProxyGroup.POST(path, handlerTrans(handlers...)...)
	router.openapiBuilder(path, "post")
	return router
}
func (router *Router) Get(path string, handlers ...Handler) *Router {
	router.ProxyGroup.GET(path, handlerTrans(handlers...)...)
	router.openapiBuilder(path, "get")
	return router
}
func (router *Router) getIgnoreOpenapi(path string, handlers ...Handler) *Router {
	router.ProxyGroup.GET(path, handlerTrans(handlers...)...)
	return router
}

func (router *Router) GetPost(path string, handlers ...Handler) *Router {
	router.Get(path, handlers...)
	router.Post(path, handlers...)
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
		if root == "./knife-ui" {
			if pathName == "" || pathName == "/" {
				pathName = "doc.html"
			} else {
				pathName = "/webjars" + pathName
			}
		} else {
			if pathName == "" || pathName == "/" {
				pathName = "index.html"
			}
		}
		assetPath = path.Join(root, pathName)
		assets, err := fs.Open(assetPath)
		if err != nil {
			c.Proxy.Status(http.StatusBadRequest)
			_, _ = c.Proxy.Writer.Write([]byte(err.Error()))
			return
		}
		data := make([]byte, 0, 1024*5)
		for {
			temp := make([]byte, 1024)
			n, _ := assets.Read(temp)
			if n == 0 {
				break
			} else {
				data = append(data, temp[:n]...)
			}
		}
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
	router.getIgnoreOpenapi("/v3/api-docs", func(c *context.Context) {
		//c.Proxy.Render(http.StatusOK, render.Data{Data: []byte(openapi.Doc.ReadDoc()), ContentType: httpconst.ContentTypeJSON})
		c.Proxy.JSON(http.StatusOK, openapi.Doc.ReadDoc())
	})
	router.getIgnoreOpenapi("/v3/api-docs/swagger-config", func(c *context.Context) {
		c.Proxy.JSON(http.StatusOK, openapi.Doc.SwaggerConfig())
	})
	// 第二个path表示匹配路径
	router.getIgnoreOpenapi("/swagger/*action", EmbedHtmlHandle(openapi.UiAssets, "./swagger-ui"))
	router.getIgnoreOpenapi("/swagger", EmbedHtmlHandle(openapi.UiAssets, "./swagger-ui"))
	router.getIgnoreOpenapi("/doc.html", EmbedHtmlHandle(openapi.KUiAssets, "./knife-ui"))
	router.getIgnoreOpenapi("/webjars/*action", EmbedHtmlHandle(openapi.KUiAssets, "./knife-ui"))
}
