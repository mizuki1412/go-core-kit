package storagekit

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/init/configkey"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
)

// GetFullPath 将path或filename用项目目录包裹
func GetFullPath(path string) string {
	if path[0] != '/' {
		path = "/" + path
	}
	p := configkit.GetString(configkey.ProjectDir, ".") + path
	return p
}

// SaveInHome 存入项目目录下, path是全路径
func SaveInHome(file *class.File, path string) {
	path = GetFullPath(path)
	_ = filekit.CheckFilePath(path)
	filekit.WriteClassFile(path, file)
}

func SaveBytesInHome(data []byte, path string) {
	path = GetFullPath(path)
	_ = filekit.CheckFilePath(path)
	_ = filekit.WriteFile(path, data)
}

func SaveBytesAppendInHome(data []byte, path string) {
	path = GetFullPath(path)
	_ = filekit.CheckFilePath(path)
	_ = filekit.WriteFileAppend(path, data)
}

func GetInHome(path string) []byte {
	path = GetFullPath(path)
	//file,err := os.OpenFile(path, os.O_RDONLY, 644)
	//if err!=nil{
	//	panic(exception.New("文件打开失败"))
	//}
	return filekit.ReadBytes(path)
}
