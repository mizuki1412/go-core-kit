package restkit

import (
	ctx "context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/init/configkey"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/middleware"
	router2 "github.com/mizuki1412/go-core-kit/service/restkit/router"
)

var router *router2.Router

func Engine() *router2.Router {
	if router == nil {
		defaultEngine()
	}
	return router
}

func defaultEngine() {
	router = &router2.Router{
		IsGroup: false,
		Proxy:   iris.New(),
		Path:    "",
	}
	//router.Proxy.Use(recover2.New())
	router.Use(middleware.Recover())
	router.Use(middleware.Log())
	router.Use(middleware.Cors())
	// max request size
	router.Proxy.Use(iris.LimitRequestBodySize(int64(configkit.GetInt(configkey.RestRequestBodySize, 100)) << 20))
	// 其他错误如404，
	router.OnError(middleware.Cors())
	// add base path
	base := configkit.GetStringD(configkey.RestServerBase)
	if base != "" {
		if base[0] != '/' {
			base = "/" + base
		}
		router.ProxyGroup = router.Proxy.Party(base)
		router.IsGroup = true
		//router.Path = base
	}

	/// init session
	context.InitSession()
}

func Run() error {
	if router == nil {
		defaultEngine()
	}
	port := configkit.GetString(configkey.RestServerPort, "8080")
	logkit.Info("Listening and serving HTTP on " + port)
	//err := http.ListenAndServe(":" + port, middleware.Session.LoadAndSave(router))
	if configkit.GetBool(configkey.RestPPROF, false) {
		p := pprof.New()
		if router.IsGroup {
			router.ProxyGroup.Any("/debug/pprof", p)
			router.ProxyGroup.Any("/debug/pprof/{action:path}", p)
		} else {
			router.Proxy.Any("/debug/pprof", p)
			router.Proxy.Any("/debug/pprof/{action:path}", p)
		}
	}
	router.RegisterSwagger()
	err := router.Proxy.Run(
		iris.Addr(":"+port),
		// 禁用，阻止如 /xx/ 自动重定向到 /xx，而不经过handle
		iris.WithoutPathCorrection)
	return err
}

func Shutdown() {
	if router != nil && router.Proxy != nil {
		err := router.Proxy.Shutdown(ctx.Background())
		if err != nil {
			logkit.Error(exception.New(err.Error()))
		}
	}
}

// AddActions 导入业务模块，其中的路由和中间件
func AddActions(actionInits ...func(r *router2.Router)) {
	if router == nil {
		defaultEngine()
	}
	for _, action := range actionInits {
		action(router)
	}
}

func GetRouter() *router2.Router {
	if router == nil {
		defaultEngine()
	}
	return router
}
