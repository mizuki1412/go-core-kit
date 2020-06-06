package restkit

import (
	"github.com/kataras/iris/v12"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/middleware"
	router2 "github.com/mizuki1412/go-core-kit/service/restkit/router"
	"github.com/spf13/cast"
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
	// add base path
	base := configkit.GetStringD(ConfigKeyRestServerBase)
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

func Run() {
	if router == nil {
		defaultEngine()
	}
	port := cast.ToString(configkit.Get(ConfigKeyRestServerPort, 8080))
	logkit.Info("Listening and serving HTTP on " + port)
	//err := http.ListenAndServe(":" + port, middleware.Session.LoadAndSave(router))
	router.RegisterSwagger()
	err := router.Proxy.Run(iris.Addr(":" + port))
	if err != nil {
		logkit.Fatal(err.Error())
	}
}

// 导入业务模块，其中的路由和中间件
func AddActions(actionInits ...func(r *router2.Router)) {
	if router == nil {
		defaultEngine()
	}
	for _, action := range actionInits {
		action(router)
	}
}
