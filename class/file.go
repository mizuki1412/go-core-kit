package class

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"io"
	"mime/multipart"
)

type File struct {
	File   multipart.File
	Header *multipart.FileHeader
}

func (th *File) GetBytes() []byte {
	bytes, err := io.ReadAll(th.File)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	return bytes
}
