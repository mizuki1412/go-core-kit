package bytekit

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"strconv"
)

func Bytes2HexArray(bytes []byte) string {
	str := "["
	//str2 := "["
	for _, v := range bytes {
		//str2 = str2+fmt.Sprintf("0x%02x ", v)
		val := strconv.FormatInt(int64(v), 16)
		if len(val) == 1 {
			val = "0" + val
		}
		str = str + "0x" + val + " "
	}
	str = str[:len(str)-1] + "]"
	return str
}

// [0x00 0x00]
func HexString2Bytes(src string) []byte {
	if len(src) < 6 {
		panic(exception.New("数据长度错误"))
	}
	if src[0] != '[' || src[len(src)-1] != ']' {
		panic(exception.New("数据长度错误"))
	}
	src = src[1 : len(src)-1]
	var ret []byte
	for _, e := range stringkit.Split(src, " ") {
		v, err := strconv.ParseUint(e, 0, 0)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		ret = append(ret, byte(v))
	}
	return ret
}

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
