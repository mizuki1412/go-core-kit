package restkit

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"mizuki/project/core-kit/service/configkit"
	"mizuki/project/core-kit/service/logkit"
	"mizuki/project/core-kit/service/restkit/middleware"
	"net/http"
)

var engine *gin.Engine

func init() {
	DefaultEngine()
}

func DefaultEngine() *gin.Engine {
	engine = gin.New()
	engine.Use(middleware.Log())
	engine.Use(middleware.Cors())
	return engine
}

func Run() {
	port := cast.ToString(configkit.Get(ConfigKeyRestServerPort, 8080))
	logkit.Info("Listening and serving HTTP on "+port)
	err := http.ListenAndServe(":" + port, middleware.Session.LoadAndSave(engine))
	if err!=nil{
		logkit.Fatal(err.Error())
	}
}

// 导入业务模块，其中的路由和中间件
func AddMod(modInit func(engine *gin.Engine)) {
	modInit(engine)
}