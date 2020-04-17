package stringkit

import (
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
