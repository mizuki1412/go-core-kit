package bytekit

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"strconv"
)

// format: 0102030a0d
func HexString2Bytes1(src string) []byte {
	if len(src)%2 != 0 || len(src) == 0 {
		panic(exception.New("数据长度错误"))
	}
	ret := make([]byte, len(src)/2)
	for i := 0; i < len(src); i = i + 2 {
		v, err := strconv.ParseUint("0x"+src[i:i+2], 0, 0)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		ret[i/2] = byte(v)
	}
	return ret
}

// format: 0102030a0d
func Bytes2HexString1(data []byte) string {
	ret := ""
	for _, e := range data {
		v := strconv.FormatInt(int64(e), 16)
		if len(v) == 1 {
			v = "0" + v
		}
		ret += v
	}
	return ret
}
