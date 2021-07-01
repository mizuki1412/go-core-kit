package tarkit

import (
	"archive/tar"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CreateTar todo 存在问题：用windows解压时，出现PaxHeader文件、同时中文不显示。
func CreateTar(resource, target string, deleteIfExist bool) {
	//create tar file with targetfile name
	tf, err := os.Create(target)
	if err != nil {
		// if file is exist then delete file
		if err == os.ErrExist && deleteIfExist {
			_ = os.Remove(target)
			tf, err = os.Create(target)
			if err != nil {
				panic(exception.New("file create error:" + err.Error()))
			}
		} else {
			panic(exception.New("file create error:" + err.Error()))
		}
	}
	defer tf.Close()
	tarWriter := tar.NewWriter(tf)
	fileInfo, err := os.Stat(resource)
	if err != nil {
		panic(exception.New("file info err:" + err.Error()))
	}
	if !fileInfo.IsDir() {
		tarFile(target, resource, fileInfo, tarWriter)
	} else {
		tarFolder(resource, tarWriter)
	}
}
func tarFile(directory string, filesource string, sfileInfo os.FileInfo, tarwriter *tar.Writer) {
	sfile, err := os.Open(filesource)
	if err != nil {
		panic(err)
	}
	defer sfile.Close()
	header, err := tar.FileInfoHeader(sfileInfo, "")
	if err != nil {
		panic(exception.New("file info err:" + err.Error()))
	}
	header.Name = directory
	//header.Format = tar.FormatGNU
	err = tarwriter.WriteHeader(header)
	if err != nil {
		panic(exception.New("file header err:" + err.Error()))
	}
	//  can use buffer to copy the file to tar writer
	//    buf := make([]byte,15)
	//    if _, err = io.CopyBuffer(tarwriter, sfile, buf); err != nil {
	//        panic(err)
	//        return err
	//    }
	if _, err = io.Copy(tarwriter, sfile); err != nil {
		panic(exception.New("file copy err:" + err.Error()))
	}
}
func tarFolder(directory string, tarwriter *tar.Writer) {
	var baseFolder = filepath.Base(directory)
	err := filepath.Walk(directory, func(targetpath string, file os.FileInfo, err error) error {
		//read the file failure
		if file == nil {
			return err
		}
		if file.IsDir() {
			// information of file or folder
			header, err := tar.FileInfoHeader(file, "")
			if err != nil {
				return err
			}
			header.Name = filepath.Join(baseFolder, strings.TrimPrefix(targetpath, directory))
			header.Format = tar.FormatGNU
			if err = tarwriter.WriteHeader(header); err != nil {
				return err
			}
			_ = os.Mkdir(strings.TrimPrefix(baseFolder, file.Name()), os.ModeDir)
		} else {
			//baseFolder is the tar file path
			var fileFolder = filepath.Join(baseFolder, strings.TrimPrefix(targetpath, directory))
			tarFile(fileFolder, targetpath, file, tarwriter)
		}
		return nil
	})
	if err != nil {
		panic(exception.New("file walk err:" + err.Error()))
	}
}
