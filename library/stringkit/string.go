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
