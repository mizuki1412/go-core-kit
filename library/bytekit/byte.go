package bytekit

import (
	"fmt"
)

func Bytes2HexArray(bytes []byte) string {
	str := "["
	for _,v := range bytes{
		str = str+fmt.Sprintf("0x%02x ", v)
	}
	str = str+"]"
	return str
}