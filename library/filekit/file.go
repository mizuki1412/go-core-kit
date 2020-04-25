package filekit

import (
	"github.com/spf13/afero"
	"io/ioutil"
	"os"
	"strings"
)

func WriteFileAppend(fileName string, data []byte) error {
	err := checkDir(fileName)
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
	err := checkDir(fileName)
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

func checkDir(fileName string) error {
	i := strings.LastIndex(fileName, "/")
	if i > 0 {
		exist, err := afero.Exists(afero.NewOsFs(), fileName[:i])
		if err != nil {
			return err
		}
		if !exist {
			err = os.MkdirAll(fileName[:i], 0755)
		}
		return err
	}
	return nil
}

func ReadString(fileName string) (string, error) {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(f), err
}
