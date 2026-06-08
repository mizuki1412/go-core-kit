package restkit

import (
	ctx "context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/cli/configkey"
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
	"github.com/mizuki1412/go-core-kit/v2/service/logkit"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/middleware"
	router2 "github.com/mizuki1412/go-core-kit/v2/service/restkit/router"
	"github.com/spf13/cast"
)

var router *router2.Router
var server *http.Server

func defaultEngine() {
	if !configkit.GetBool(configkey.ProfileDev) {
		gin.SetMode(gin.ReleaseMode)
	}
	router = &router2.Router{
		Proxy: gin.New(),
	}
	// add base path
	base := configkit.GetString(configkey.RestServerBase)
	if base != "" {
		if base[0] != '/' {
			base = "/" + base
		}
		if base[len(base)-1] == '/' {
			base = base[:len(base)-1]
		}
		router.ProxyGroup = router.Proxy.Group(base)
	} else {
		router.ProxyGroup = &router.Proxy.RouterGroup
	}
	router.Use(middleware.Log())
	router.Use(middleware.Cors())
	router.Use(middleware.Recover())
	if configkit.GetBool(configkey.RestPPROF) {
		//  p := pprof.New()
	}
	// max request size todo
	//router.Proxy.Use(iris.LimitRequestBodySize(int64(configkit.GetInt(configkey.RestRequestBodySize, 100)) << 20))
	// 其他错误如404，
	//router.OnError(middleware.Cors())
}

func Run(listeners ...net.Listener) error {
	if router == nil {
		defaultEngine()
	}
	port := configkit.GetString(configkey.RestServerPort)
	router.RegisterSwagger()
	if len(listeners) == 0 {
		server = &http.Server{
			Addr:    ":" + port,
			Handler: router,
		}
	} else {
		server = &http.Server{
			Handler: router,
		}
		port = cast.ToString(listeners[0].Addr().(*net.TCPAddr).Port)
	}
	go func() {
		logkit.Info("Listening and serving HTTP on " + port)
		if len(listeners) == 0 {
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				panic(exception.New(err.Error()))
			}
		} else {
			if err := server.Serve(listeners[0]); err != nil && !errors.Is(err, http.ErrServerClosed) {
				panic(exception.New(err.Error()))
			}
		}
	}()
	// https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-without-context/server.go
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, os.Interrupt, syscall.SIGTERM)
	<-quit
	logkit.Info("Shutting down server...")

	ctxt, cancel := ctx.WithTimeout(ctx.Background(), 5*time.Second)
	defer cancel()
	// 执行自定义关机逻辑
	CustomShutdownLogic(ctxt)
	if err := server.Shutdown(ctxt); err != nil {
		logkit.Error(err.Error())
		return err
	}
	return nil
}

func Shutdown() {
	if server != nil {
		logkit.Info("Shutting down server...")
		err := server.Shutdown(ctx.Background())
		if err != nil {
			logkit.Error(err.Error())
		}
	}
}

var CustomShutdownLogic = func(ctx ctx.Context) {}

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
