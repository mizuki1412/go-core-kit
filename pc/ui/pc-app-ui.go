package ui

import (
	"embed"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/cli/configkey"
	"github.com/mizuki1412/go-core-kit/service/restkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"net"
)

// WinParam ui window params
type WinParam struct {
	// 默认0则开启随机端口; 这里的port将替代config中的rest port
	Port       int
	Width      int
	Height     int
	Title      string
	FullScreen bool
	Debug      bool
	Assets     embed.FS
	// true：只开启服务，不开启ui
	NoUI bool
	// todo 不开启rest server
	NoRest bool
	// 关闭ui时不关闭主线程
	KeepMain bool
}

// 开启websocket server
func startServer(param *WinParam) chan error {
	// todo 开启bridge
	//bridge.Start()
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
	// 设置 ui assets
	restkit.GetRouter().Get("/ui/*action", router.EmbedHtmlHandle(param.Assets, "./ui"))
	restkit.GetRouter().Get("/ui", router.EmbedHtmlHandle(param.Assets, "./ui"))
	// 设置restkit的port
	viper.Set(configkey.RestServerPort, cast.ToString(param.Port))
	// 启动ui的http server
	go func() {
		if listener == nil {
			ch <- restkit.Run()
		} else {
			ch <- restkit.Run(listener)
		}
	}()
	return ch
}

func Run(param *WinParam) {
	// todo
	//serverCh := startServer(param)
	//if !param.NoUI {
	//	startUI(param)
	//}
	//waitClose(param, serverCh)
}
