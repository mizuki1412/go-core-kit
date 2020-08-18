package ui

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/pc/bridge"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cast"
	"net"
	"net/http"
)

/// ui window params
type WinParam struct {
	// 默认0则开启随机端口
	Port       int
	Width      int
	Height     int
	FullScreen bool
	// url中的项目路径 eg：/base/xxx。如果前端是hash模式：/#/xxx
	Url string
	// true：只开启服务，不开启ui
	NoUI bool
}

// 开启websocket server
func startServer(param *WinParam) (chan error, string) {
	// 开启bridge
	bridge.Start()
	ch := make(chan error)
	var listener net.Listener
	var err error
	if param.Port == 0 {
		// 获取随机端口
		listener, err = net.Listen("tcp", ":0")
		if err != nil {
			panic(exception.New("随机端口开启失败"))
		}
		param.Port = listener.Addr().(*net.TCPAddr).Port
	}
	p := cast.ToString(param.Port)
	// 启动ui的http server
	go func() {
		logkit.Info("pc ui http server start at: " + p)
		if listener != nil {
			ch <- http.Serve(listener, nil)
		} else {
			ch <- http.ListenAndServe(":"+p, nil)
		}
	}()
	return ch, p
}
