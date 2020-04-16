package stringkit

import (
	"log"
	"strconv"
	"strings"
)

func ToString(obj interface{}) string {
	switch obj.(type) {
	case uint:
		return strconv.FormatUint(uint64(obj.(uint)), 10)
	case uint8:
		return strconv.FormatUint(uint64(obj.(uint8)), 10)
	case uint16:
		return strconv.FormatUint(uint64(obj.(uint16)), 10)
	case uint32:
		return strconv.FormatUint(uint64(obj.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(obj.(uint64), 10)
	case int:
		return strconv.FormatInt(int64(obj.(int)), 10)
	case int8:
		return strconv.FormatInt(int64(obj.(int8)), 10)
	case int16:
		return strconv.FormatInt(int64(obj.(int16)), 10)
	case int32:
		return strconv.FormatInt(int64(obj.(int32)), 10)
	case int64:
		return strconv.FormatInt(obj.(int64), 10)
	case bool:
		return strconv.FormatBool(obj.(bool))
	case float32:
		return strconv.FormatFloat(float64(obj.(float32)), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(obj.(float64), 'f', -1, 64)
	case string:
		return obj.(string)
	}
	log.Println("stringkit.ToString: not supported")
	return ""
}

func IsNull(obj interface{}) bool {
	if obj == nil {
		return true
	}
	if strings.Trim(obj.(string), " ") == "" {
		return true
	}
	return false
}
