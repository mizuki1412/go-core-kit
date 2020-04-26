package stringkit

import (
	"mizuki/project/core-kit/class/exception"
	"regexp"
	"strings"
)

func IsNull(obj interface{}) bool {
	if obj == nil {
		return true
	}
	if strings.TrimSpace(obj.(string)) == "" {
		return true
	}
	return false
}

func Concat(strs ...string) string {
	// 内部用的Builder
	return strings.Join(strs, "")
}

// 正则切割
func Split(origin, reg string) []string {
	r, err := regexp.Compile(reg)
	if err != nil {
		panic(exception.New("regex error", 2))
	}
	return r.Split(origin, -1)
}

func CamelCase(str string) string {
	strs := strings.Split(str, "_")
	temp := ""
	for _, v := range strs {
		temp += strings.ToUpper(v[:1]) + v[1:]
	}
	return temp
}
