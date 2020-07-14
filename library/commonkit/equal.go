package commonkit

import "github.com/spf13/cast"

// todo 比较两者是否相等. 暂时通过string
func Equal(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		return false
	}
	return cast.ToString(a) == cast.ToString(b)
}
