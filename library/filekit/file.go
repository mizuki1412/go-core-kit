package filekit

import (
	"io/ioutil"
	"os"
)

func AppendToFile(fileName string, content string) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err == nil {
		_, err = f.WriteString(content)
	}
	defer f.Close()
	return err
}

func ReadString(fileName string) (string,error) {
	f, err := ioutil.ReadFile(fileName)
	if err!=nil {
		return "",err
	}
	return string(f),err
}