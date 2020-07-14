package arraykit

import "github.com/spf13/cast"

// todo sort.searchStrings()

func StringContains(arr []string, ele string) bool {
	for _, v := range arr {
		if v == ele {
			return true
		}
	}
	return false
}

func StringDelete(arr []string, ele string) []string {
	j := 0
	for _, val := range arr {
		if val != ele {
			arr[j] = val
			j++
		}
	}
	return arr[:j]
}
func StringDeleteAt(arr []string, index int) []string {
	j := 0
	for i, val := range arr {
		if i != index {
			arr[j] = val
			j++
		}
	}
	return arr[:j]
}

// 此种方法会修改arr原始值，使用场景必须是arr一次性覆盖的时候
// 同时要注意比较值的类型，json转过的一般是int
func Delete(arr []interface{}, ele interface{}) []interface{} {
	j := 0
	for _, val := range arr {
		if cast.ToString(val) != cast.ToString(ele) {
			arr[j] = val
			j++
		}
	}
	return arr[:j]
}
func DeleteAt(arr []interface{}, index int) []interface{} {
	j := 0
	for i, val := range arr {
		if i != index {
			arr[j] = val
			j++
		}
	}
	return arr[:j]
}
