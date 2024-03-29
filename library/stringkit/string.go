package stringkit

import (
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/spf13/cast"
	"regexp"
	"strings"
)

func IsNull(obj any) bool {
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

func ConcatWith(arr []string, sep, decorate string) string {
	if len(arr) == 0 {
		return ""
	}
	fin := ""
	for _, v := range arr {
		fin += decorate + v + decorate + sep
	}
	return fin[:strings.LastIndex(fin, sep)]
}

func ConcatIntWith(arr []int32, sep string) string {
	if len(arr) == 0 {
		return ""
	}
	fin := ""
	for _, v := range arr {
		fin += cast.ToString(v) + sep
	}
	return fin[:strings.LastIndex(fin, sep)]
}

// Split 正则切割
func Split(origin, reg string) []string {
	r, err := regexp.Compile(reg)
	if err != nil {
		panic(exception.New("regex error", 2))
	}
	return r.Split(origin, -1)
}

func CamelCase(str string) string {
	if str == "" {
		return str
	}
	strs := strings.Split(str, "_")
	temp := ""
	for _, v := range strs {
		temp += UpperFirst(v)
	}
	return temp
}

// UpperFirst 首字母大写
func UpperFirst(str string) string {
	if str == "" {
		return str
	}
	bytes := []byte(str)
	if bytes[0] >= 'a' {
		bytes[0] -= 32
	}
	return string(bytes)
}

// LowerFirst 首字母小写
func LowerFirst(str string) string {
	if str == "" {
		return str
	}
	bytes := []byte(str)
	if bytes[0] < 'a' {
		bytes[0] += 32
	}
	return string(bytes)
}

func MatchReg(reg string, target string) bool {
	r, err := regexp.Compile(reg)
	if err == nil {
		return r.Match([]byte(target))
	}
	return false
}

// ClearFilePath 去掉文件路径末尾的/或\
func ClearFilePath(path string) string {
	if len(path) > 1 {
		if path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
		if path[len(path)-1] == '\\' {
			path = path[:len(path)-1]
		}
	}
	return path
}
