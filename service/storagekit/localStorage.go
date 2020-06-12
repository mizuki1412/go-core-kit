package storagekit

import (
	corekit "github.com/mizuki1412/go-core-kit"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
)

func getFullPath(path string) string {
	if path[0] != '/' {
		path = "/" + path
	}
	return configkit.GetString(corekit.ConfigKeyProjectDir, ".") + path
}

func SaveInHome(file *class.File, path string) {
	path = getFullPath(path)
	filekit.WriteClassFile(path, file)
}

func GetInHome(path string) []byte {
	path = getFullPath(path)
	//file,err := os.OpenFile(path, os.O_RDONLY, 644)
	//if err!=nil{
	//	panic(exception.New("文件打开失败"))
	//}
	return filekit.ReadBytes(path)
}
