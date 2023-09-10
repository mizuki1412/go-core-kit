//go:build windows && webview2

package ui

import (
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/go-core-kit/service/restkit"
	"github.com/mizuki1412/go-webview2"
	"github.com/spf13/cast"
)

var w webview2.WebView

// 单独启动ui时，用于重开ui
func startUI(param *WinParam) {
	port := cast.ToString(param.Port)
	w = webview2.New(param.Debug)
	if w == nil {
		logkit.Error("Failed to load webview2.")
		return
	}
	defer w.Destroy()
	w.SetTitle(param.Title)
	w.SetSize(param.Width, param.Height, webview2.HintFixed)
	w.Navigate("http://127.0.0.1:" + port + "/ui")
	w.Run()
	restkit.Shutdown()
	//if param.FullScreen {
	//	_ = Self.SetBounds(lorca.Bounds{
	//		WindowState: lorca.WindowStateFullscreen,
	//	})
	//}
}

func waitClose(param *WinParam, serverCh chan error) {

}
