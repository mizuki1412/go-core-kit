package restkit

import (
	"github.com/kataras/iris/v12"
	recover2 "github.com/kataras/iris/v12/middleware/recover"
	"github.com/spf13/cast"
	"mizuki/project/core-kit/service/configkit"
	"mizuki/project/core-kit/service/logkit"
	"mizuki/project/core-kit/service/restkit/context"
	"mizuki/project/core-kit/service/restkit/middleware"
)

var router *context.Router

func init() {
	defaultEngine()
}

//func defaultEngine2() *gin.Engine {
//	engine := gin.New()
//	//engine.Use(middleware.Log())
//	//engine.Use(middleware.Cors())
//	return engine
//}

func defaultEngine()  {
	router = &context.Router{
		IsGroup: false,
		Proxy: iris.New(),
	}
	router.Proxy.Use(recover2.New())
	//router.Use(middleware.Session())
	router.Use(middleware.Log())
	router.Use(middleware.Cors())
}

func Run() {
	port := cast.ToString(configkit.Get(ConfigKeyRestServerPort, 8080))
	logkit.Info("Listening and serving HTTP on "+port)
	//err := http.ListenAndServe(":" + port, middleware.Session.LoadAndSave(router))
	err := router.Proxy.Run(iris.Addr(":"+port))
	if err!=nil{
		logkit.Fatal(err.Error())
	}
}

// 导入业务模块，其中的路由和中间件
func AddMod(modInit func(r *context.Router)) {
	modInit(router)
}