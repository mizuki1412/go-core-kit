package ui

import (
	"github.com/mizuki1412/go-core-kit/library/processkit"
	"github.com/spf13/cast"
	"log"
)

/// todo
func startElectron(param WinParam) {
	p := cast.ToString(param.Port)
	ret, err := processkit.Cmd("./electron-app.exe", "--url=http://127.0.0.1:"+p, "--debug")
	if !ret {
		log.Fatalln(err)
	}
}
