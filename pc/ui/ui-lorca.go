//go:build lorca

package ui

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/init/configkey"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cast"
	"github.com/zserge/lorca"
	"os"
	"os/signal"
)

var self lorca.UI

// 单独启动ui时，用于重开ui
func startUI(param *WinParam) {
	var err error
	// 增加user-data-dir后，一些配置将会存入其中，包括安全策略的设置
	// windows下，每次开启可能提示未正确关闭：需要在设置的user-data-dir中的Default/Preferences的exit_type为Normal，并设置文件为只读
	pdr := configkit.GetStringD(configkey.ProjectDir)
	if pdr == "" {
		pdr = "."
	}
	port := cast.ToString(param.Port)
	self, err = lorca.New(
		"", pdr+"/lorca", param.Width, param.Height,
		"--disable-web-security", // 不遵守同源策略
		"--allow-insecure-localhost",
		"--allow-running-insecure-content",
		"--unsafely-treat-insecure-origin-as-secure=http://localhost:"+port+",http://127.0.0.1:"+port, // not work
		"--reduce-security-for-testing",
	)
	// local web ui地址。
	if err != nil {
		panic(exception.New(err.Error()))
	}
	_ = self.Load("http://localhost:" + port + "/ui")
	if param.FullScreen {
		_ = self.SetBounds(lorca.Bounds{
			WindowState: lorca.WindowStateFullscreen,
		})
	}
	//defer self.Close()
}

func waitClose(param *WinParam, serverCh chan error) {
	// 监听关闭信号
	sign := make(chan os.Signal)
	signal.Notify(sign, os.Interrupt)
	if self == nil || param.KeepMain {
		// todo self.Done()时 设置为nil
		select {
		case <-sign:
			logkit.Info("close main")
		case err := <-serverCh:
			logkit.Info("server down: " + err.Error())
		}
	} else if self != nil {
		//defer self.Close()
		select {
		case <-sign:
			logkit.Info("close main")
		case <-self.Done():
			logkit.Info("close ui")
		case err := <-serverCh:
			logkit.Info("server down: " + err.Error())
		}
	}
}
