package ftpkit

import (
	"bytes"
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"io"
	"os"
	"strconv"
	"time"
)

type Conn struct {
	C *ftp.ServerConn
}

// NewClient 注意使用完后Close
func NewClient(host string, port string, user string, pwd string) *Conn {
	c, err := ftp.Dial(host+":"+port, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		panic(exception.New(err.Error()))
	}
	err = c.Login(user, pwd)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	return &Conn{C: c}
}

func (th *Conn) Close() {
	_ = th.C.Quit()
}

// DownloadSlice 下载，断点续传
func (th *Conn) DownloadSlice(src string, dst string, buffer int) {
	temp := dst + ".temp"
	var index int64
	if _, err := os.Stat(temp); err != nil {
		index = 0
	} else if err == nil {
		data, err := os.ReadFile(temp)
		if err != nil {
			panic(err)
		}
		index, _ = strconv.ParseInt(string(data), 10, 64)
	}
	reader, err := th.C.RetrFrom(src, uint64(index))
	if err != nil {
		panic(exception.New(err.Error()))
	}
	dstFile, _ := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer dstFile.Close()
	dstFile.Seek(index, 0)
	buf := make([]byte, buffer, buffer)
	var n2 int
	total := int(index)
	for {
		n1, err := reader.Read(buf)
		if err == io.EOF {
			fmt.Println("done")
			if err = os.Remove(temp); err != nil {
				panic(exception.New(err.Error()))
			}
			break
		}
		n2, _ = dstFile.Write(buf[:n1])
		total += n2
		err = os.WriteFile(temp, []byte(strconv.Itoa(total)), 0666)
		if err != nil {
			panic(exception.New(err.Error()))
		}
	}
}

// UploadSlice 上传，断点续传
func (th *Conn) UploadSlice(src string, dst string) {
	index := uint64(0)
	dsttmp := dst + ".tmp"
	if target, _ := th.C.List(dsttmp); target != nil {
		index = target[0].Size
	}
	srcFile, err := os.Open(src)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	defer srcFile.Close()
	srcFile.Seek(int64(index), 0)
	buf := make([]byte, 4096, 4096)
	for {
		n, err := srcFile.Read(buf)
		if err == io.EOF {
			fmt.Println("done")
			th.C.Rename(dsttmp, dst)
			break
		}
		data := bytes.NewReader(buf)
		th.C.StorFrom(dsttmp, data, index)
		index += uint64(n)
	}
}
