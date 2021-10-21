package ui

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	corekit "github.com/mizuki1412/go-core-kit/init"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/zserge/lorca"
	"os"
	"os/signal"
)

var Self lorca.UI
var serverCh chan error
var port string

// StartUI 单独启动ui时，用于重开ui
func StartUI(param *WinParam) {
	var err error
	// 增加user-data-dir后，一些配置将会存入其中，包括安全策略的设置
	// windows下，每次开启可能提示未正确关闭：需要在设置的user-data-dir中的Default/Preferences的exit_type为Normal，并设置文件为只读
	Self, err = lorca.New(
		"", configkit.GetStringD(corekit.ConfigKeyProjectDir), param.Width, param.Height,
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
	if param.Url != "" {
		_ = Self.Load("http://localhost:" + port + param.Url)
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
		// todo Self.Done()时 设置为nil
		select {
		case <-sign:
			logkit.Info("close main")
		case err := <-serverCh:
			logkit.Info("server down: " + err.Error())
		}
	} else if Self != nil {
		//defer Self.Close()
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
