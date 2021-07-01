package tarkit

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Zip(dst, src string) (err error) {
	// 创建准备写入的文件
	fw, err := os.Create(dst)
	defer fw.Close()
	if err != nil {
		return err
	}

	// 通过 fw 来创建 zip.Write
	zw := zip.NewWriter(fw)
	defer zw.Close()

	// 去除压缩包中文件的目录前缀
	delPrefix := src[0:strings.LastIndexAny(src, string(filepath.Separator))]

	// 下面来将文件写入 zw ，因为有可能会有很多个目录及文件，所以递归处理
	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}

		// 通过文件信息，创建 zip 的文件信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		// 替换文件信息中的文件名
		fh.Name = strings.TrimPrefix(path, delPrefix)

		// 这步开始没有加，会发现解压的时候说它不是个目录
		if fi.IsDir() {
			fh.Name += "/"
		}

		// 写入文件信息，并返回一个 Write 结构
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}

		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
		// 如目录，也没有数据需要写
		if !fh.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩到的文件
		fr, err := os.Open(path)
		defer fr.Close()
		if err != nil {
			return err
		}

		// 将打开的文件 Copy 到 w
		_, err = io.Copy(w, fr)
		if err != nil {
			return err
		}
		return nil
	})
}
