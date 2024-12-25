package filekit

import (
	"github.com/mizuki1412/go-core-kit/v2/service/logkit"
	"io/fs"
	"os"
	"path/filepath"
)

func ListFileNames(dir string) []string {
	var files []string
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if path == dir {
			return nil
		}
		if !info.IsDir() {
			files = append(files, path)
		} else {
			files = append(files, ListFileNames(path)...)
		}
		return nil
	})
	if err != nil {
		logkit.Error(err.Error())
	}
	return files
}

// ProgramWorkDir 当前程序运行目录
func ProgramWorkDir() (string, error) {
	return os.Getwd()
}

// ProgramLocationDir 当前程序所在目录
func ProgramLocationDir() (string, error) {
	return os.Executable()
}
