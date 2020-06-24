package ui

import (
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cast"
	"github.com/zserge/lorca"
	"net/http"
	"os"
	"os/signal"
)

func StartLorca(param WinParam) {
	ch := make(chan error)
	p := cast.ToString(param.Port)
	// 启动ui的http server
	go func() {
		logkit.Info("Serving at localhost: " + p)
		ch <- http.ListenAndServe(":"+p, nil)
	}()
	ui, _ := lorca.New("", "", param.Width, param.Height)
	// local web ui地址
	_ = ui.Load("http://127.0.0.1:" + p + param.Url)
	if param.FullScreen {
		_ = ui.SetBounds(lorca.Bounds{
			WindowState: lorca.WindowStateFullscreen,
		})
	}
	defer ui.Close()
	// 监听关闭信号
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
		logkit.Info("close main")
	case <-ui.Done():
		logkit.Info("close ui")
	case <-ch:
		logkit.Info("server down")
	}
}
