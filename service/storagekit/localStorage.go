package storagekit

import (
	corekit "github.com/mizuki1412/go-core-kit"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
)

func SaveInHome(file *class.File, path string) {
	if path[0] != '/' {
		path = "/" + path
	}
	path = configkit.GetString(corekit.ConfigKeyProjectDir, ".") + path
	filekit.WriteClassFile(path, file)
}
