package bytekit

import (
	"bytes"
	"encoding/binary"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cast"
)

func num2Bytes(i interface{}) []byte {
	arr := make([]byte, 0)
	buf := bytes.NewBuffer(arr)
	// 数字转 []byte, 网络字节序为大端字节序
	err := binary.Write(buf, binary.BigEndian, i)
	if err != nil {
		logkit.Error(err)
		return arr
	}
	return buf.Bytes()
}

func Int32ToBytes(i interface{}) []byte {
	return num2Bytes(cast.ToInt32(i))
}

func Int64ToBytes(i interface{}) []byte {
	return num2Bytes(cast.ToInt64(i))
}

func Float32ToBytes(i interface{}) []byte {
	return num2Bytes(cast.ToFloat32(i))
}

func Float64ToBytes(i interface{}) []byte {
	return num2Bytes(cast.ToFloat64(i))
}

// 必须是4个字节
func Bytes2Int32(bs []byte) int32 {
	length := len(bs)
	switch length {
	case 0:
		return 0
	case 1:
		return cast.ToInt32(bs[0])
	case 2:
		buf := bytes.NewBuffer(bs)
		var target int16
		_ = binary.Read(buf, binary.BigEndian, &target)
		return cast.ToInt32(target)
	case 3:
		data := []byte{0x00}
		data = append(data, bs...)
		bs = data
		fallthrough
	default:
		buf := bytes.NewBuffer(bs)
		var target int32
		_ = binary.Read(buf, binary.BigEndian, &target)
		return target
	}
}

func Bytes2Int64(bs []byte) int64 {
	length := len(bs)
	if length <= 4 {
		return cast.ToInt64(Bytes2Int32(bs))
	}
	if length < 8 {
		data := make([]byte, 0, 8)
		for k := 0; k < 8-length; k++ {
			data = append(data, 0x00)
		}
		data = append(data, bs...)
		bs = data
	}
	buf := bytes.NewBuffer(bs)
	var i2 int64
	_ = binary.Read(buf, binary.BigEndian, &i2)
	return i2
}

func Bytes2Float32(bs []byte) float32 {
	buf := bytes.NewBuffer(bs)
	var i2 float32
	_ = binary.Read(buf, binary.BigEndian, &i2)
	return i2
}

func Bytes2Float64(bs []byte) float64 {
	buf := bytes.NewBuffer(bs)
	var i2 float64
	_ = binary.Read(buf, binary.BigEndian, &i2)
	return i2
}
