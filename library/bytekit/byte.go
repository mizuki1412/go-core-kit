package bytekit

import (
	"bytes"
	"encoding/binary"
	"github.com/spf13/cast"
	"log"
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

func num2Bytes(i interface{}) []byte {
	arr := make([]byte, 0)
	buf := bytes.NewBuffer(arr)
	// 数字转 []byte, 网络字节序为大端字节序
	err := binary.Write(buf, binary.BigEndian, i)
	if err != nil {
		log.Println(err)
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

func Bytes2Int32(bs []byte) int32 {
	buf := bytes.NewBuffer(bs)
	var i2 int32
	_ = binary.Read(buf, binary.BigEndian, &i2)
	return i2
}

func Bytes2Int64(bs []byte) int64 {
	buf := bytes.NewBuffer(bs)
	var i2 int64
	_ = binary.Read(buf, binary.BigEndian, &i2)
	return i2
}
