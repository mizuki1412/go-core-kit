package tarkit

import (
	"bytes"
	"compress/gzip"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"io"
)

// GZipBytes 压缩
func GZipBytes(data []byte) []byte {
	var input bytes.Buffer
	g := gzip.NewWriter(&input)
	defer g.Close()
	_, err := g.Write(data)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	g.Flush()
	return input.Bytes()
}

// UnGZipBytes 解压
func UnGZipBytes(data []byte) []byte {
	var in bytes.Buffer
	in.Write(data)
	r, _ := gzip.NewReader(&in)
	defer r.Close()
	undatas, _ := io.ReadAll(r)
	return undatas
}
