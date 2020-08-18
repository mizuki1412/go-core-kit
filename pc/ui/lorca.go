package ui

import (
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/zserge/lorca"
	"os"
	"os/signal"
)

func StartLorca(param *WinParam) {
	serverCh, p := startServer(param)
	var ui lorca.UI
	if !param.NoUI {
		ui, _ = lorca.New("", "", param.Width, param.Height)
		// local web ui地址。
		_ = ui.Load("http://127.0.0.1:" + p + param.Url)
		if param.FullScreen {
			_ = ui.SetBounds(lorca.Bounds{
				WindowState: lorca.WindowStateFullscreen,
			})
		}
		defer ui.Close()
	}
	// 监听关闭信号
	sign := make(chan os.Signal)
	signal.Notify(sign, os.Interrupt)
	if ui == nil {
		select {
		case <-sign:
			logkit.Info("close main")
		case err := <-serverCh:
			logkit.Info("server down: " + err.Error())
		}
	} else {
		select {
		case <-sign:
			logkit.Info("close main")
		case <-ui.Done():
			logkit.Info("close ui")
		case err := <-serverCh:
			logkit.Info("server down: " + err.Error())
		}
	}

}
