package bytekit

import (
	"strconv"
)

func Bytes2HexArray(bytes []byte) string {
	str := "["
	//str2 := "["
	for _,v := range bytes{
		//str2 = str2+fmt.Sprintf("0x%02x ", v)
		val:=strconv.FormatInt(int64(v),16)
		if len(val)==1{
			val = "0"+val
		}
		str = str+"0x"+val+" "
	}
	str = str[:len(str)-1]+"]"
	return str
}