package storagekit

import (
	"github.com/mizuki1412/go-core-kit/class"
	corekit "github.com/mizuki1412/go-core-kit/init"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
)

func GetFullPath(path string) string {
	if path[0] != '/' {
		path = "/" + path
	}
	return configkit.GetString(corekit.ConfigKeyProjectDir, ".") + path
}

func SaveInHome(file *class.File, path string) {
	path = GetFullPath(path)
	filekit.WriteClassFile(path, file)
}

func GetInHome(path string) []byte {
	path = GetFullPath(path)
	//file,err := os.OpenFile(path, os.O_RDONLY, 644)
	//if err!=nil{
	//	panic(exception.New("文件打开失败"))
	//}
	return filekit.ReadBytes(path)
}
