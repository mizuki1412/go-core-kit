package ui

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/pc/bridge"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/go-core-kit/service/restkit"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"net"
	"net/http"
)

/// ui window params
type WinParam struct {
	// 默认0则开启随机端口; 这里的port将替代config中的默认port
	Port       int
	Width      int
	Height     int
	FullScreen bool
	// url中的项目路径 eg：/base/xxx。如果前端是hash模式：/#/xxx
	Url string
	// 完整的url
	CompleteUrl string
	// true：只开启服务，不开启ui
	NoUI bool
	// 关闭ui时不关闭主线程
	KeepMain bool
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
		//listener.Close()
	}
	p := cast.ToString(param.Port)
	// 设置restkit的port
	viper.Set(restkit.ConfigKeyRestServerPort, p)
	// 启动ui的http server
	go func() {
		logkit.Info("pc ui http server start at: " + p)
		if listener != nil {
			// todo 和restkit.Run()整合, 随机端口
			ch <- http.Serve(listener, nil)
		} else {
			ch <- restkit.Run()
			//ch <- http.ListenAndServe(":"+p, nil)
		}
	}()
	return ch, p
}
