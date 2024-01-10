package filekit

import (
	"github.com/mizuki1412/go-core-kit/v2/service/logkit"
	"io/fs"
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
