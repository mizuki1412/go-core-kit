package filekit

import (
	"bytes"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/spf13/afero"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func WriteFileAppend(fileName string, data []byte) error {
	checkDir(fileName)
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err == nil {
		defer f.Close()
		_, err = f.Write(data)
	}
	return err
}

func WriteFile(fileName string, data []byte) error {
	checkDir(fileName)
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err == nil {
		defer f.Close()
		_, err = f.Write(data)
	}
	return err
}

func checkDir(fileName string) {
	i := strings.LastIndex(fileName, "/")
	if i > 0 {
		CheckDir(fileName[:i])
	}
}

func CheckDir(path string) {
	exist, err := afero.Exists(afero.NewOsFs(), path)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	if !exist {
		_ = os.MkdirAll(path, 0755)
	}
}

func ReadString(fileName string) (string, error) {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(f), err
}

func ReadBytes(fileName string) []byte {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(exception.New("文件读取失败"))
	}
	return f
}

func WriteClassFile(filepath string, file *class.File) {
	if file.File == nil {
		panic(exception.New("文件为空"))
	}
	checkDir(filepath)
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	defer f.Close()
	_, err = io.Copy(f, file.File)
	if err != nil {
		panic(exception.New(err.Error()))
	}
}

func ReadBytesFromClassFile(file *class.File) []byte {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file.File); err != nil {
		return []byte{}
	}
	return buf.Bytes()
}

// 取一个文件路径的路径和文件名 /分隔
func SplitFilePath(path string) (string, string) {
	i := strings.LastIndex(path, "/")
	if i < 0 {
		return "", path
	} else if i == len(path) {
		return path, ""
	}
	return path[0:i], path[i+1:]
}
