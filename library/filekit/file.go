package filekit

import (
	"bytes"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/afero"
	"io"
	"os"
	"strings"
)

func WriteFileAppend(fileName string, data []byte) error {
	err := CheckFilePath(fileName)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err == nil {
		defer f.Close()
		_, err = f.Write(data)
	}
	return err
}

func WriteFile(fileName string, data []byte) error {
	err := CheckFilePath(fileName)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err == nil {
		defer f.Close()
		_, err = f.Write(data)
	}
	return err
}

func DelFile(fileName string) error {
	err := CheckFilePath(fileName)
	if err != nil {
		return err
	}
	return os.Remove(fileName)
}

func CheckFilePath(fileName string) error {
	i := strings.LastIndex(fileName, "/")
	if i > 0 {
		return CheckDir(fileName[:i])
	}
	return nil
}

func CheckDir(path string) error {
	exist, err := afero.Exists(afero.NewOsFs(), path)
	if err != nil {
		return err
	}
	if !exist {
		err = os.MkdirAll(path, 0755)
		return err
	}
	return err
}

func ReadString(fileName string) (string, error) {
	f, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(f), err
}

func ReadBytes(fileName string) []byte {
	f, err := os.ReadFile(fileName)
	if err != nil {
		logkit.Error("文件读取失败")
	}
	return f
}

func WriteClassFile(filepath string, file *class.File) {
	if file.File == nil {
		panic(exception.New("文件为空"))
	}
	_ = CheckFilePath(filepath)
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
