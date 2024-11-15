package utils

import (
	"github.com/spf13/cast"
)

func UnquoteIfQuoted(bytes []byte) string {
	// If the amount is quoted, strip the quotes
	if len(bytes) >= 2 && bytes[0] == '"' && bytes[len(bytes)-1] == '"' {
		bytes = bytes[1 : len(bytes)-1]
	}
	return string(bytes)
}

func TransScanValue2String(value any) string {
	var val string
	switch value.(type) {
	case []byte:
		val = string(value.([]byte))
	case string:
		val = value.(string)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		val = cast.ToString(value)
	}
	return val
}
