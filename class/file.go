package class

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"io"
	"mime/multipart"
)

type File struct {
	File   multipart.File
	Header *multipart.FileHeader
}

func (th File) MarshalJSON() ([]byte, error) {
	if th.Header != nil {
		return []byte(fmt.Sprintf("filename=%s;size=%d", th.Header.Filename, th.Header.Size)), nil
	}
	// 返回json中的null
	return []byte("null"), nil
}
func (th *File) UnmarshalJSON(data []byte) error {
	return nil
}

func (th File) GetBytes() []byte {
	bytes, err := io.ReadAll(th.File)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	return bytes
}
