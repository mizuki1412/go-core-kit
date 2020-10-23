package ui

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/zserge/lorca"
	"os"
	"os/signal"
)

var Self lorca.UI
var serverCh chan error
var port string

// 单独启动ui时，用于重开ui
func StartUI(param *WinParam) {
	Self, _ = lorca.New("", "", param.Width, param.Height)
	// local web ui地址。
	if param.Url != "" {
		_ = Self.Load("http://127.0.0.1:" + port + param.Url)
	} else if param.CompleteUrl != "" {
		_ = Self.Load(param.CompleteUrl)
	} else {
		panic(exception.New("url未指定"))
	}
	if param.FullScreen {
		_ = Self.SetBounds(lorca.Bounds{
			WindowState: lorca.WindowStateFullscreen,
		})
	}
	//defer Self.Close()
}

func StartLorca(param *WinParam) {
	serverCh, port = startServer(param)
	if !param.NoUI {
		StartUI(param)
	}
	// 监听关闭信号
	sign := make(chan os.Signal)
	signal.Notify(sign, os.Interrupt)
	if Self == nil || param.KeepMain {
		select {
		case <-sign:
			logkit.Info("close main")
		case err := <-serverCh:
			logkit.Info("server down: " + err.Error())
		}
	} else if Self != nil {
		defer Self.Close()
		select {
		case <-sign:
			logkit.Info("close main")
		case <-Self.Done():
			logkit.Info("close ui")
		case err := <-serverCh:
			logkit.Info("server down: " + err.Error())
		}
	}
}
