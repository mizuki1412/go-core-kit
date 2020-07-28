package pghelper

import "strings"

// 生成sql中: in (select unnest(Array[?,?,?])) []interface{}
// arr必须不能空
// 注意使用时 args...
func GenUnnestString(arr []string) (string, []interface{}) {
	flags := make([]string, len(arr))
	args := make([]interface{}, len(arr))
	for i := 0; i < len(arr); i++ {
		flags[i] = "?"
		args[i] = arr[i]
	}
	return "(select unnest(Array[" + strings.Join(flags, ", ") + "]))", args
}

func GenUnnestInt(arr []int32) (string, []interface{}) {
	flags := make([]string, len(arr))
	args := make([]interface{}, len(arr))
	for i := 0; i < len(arr); i++ {
		flags[i] = "?"
		args[i] = arr[i]
	}
	return "(select unnest(Array[" + strings.Join(flags, ", ") + "]))", args
}

// 返回 Array[?,?,?]
func GenArrayFlagString(arr []string) (string, []interface{}) {
	flags := make([]string, len(arr))
	args := make([]interface{}, len(arr))
	for i := 0; i < len(arr); i++ {
		flags[i] = "?"
		args[i] = arr[i]
	}
	return "Array[" + strings.Join(flags, ", ") + "]", args
}

func GenArrayFlagInt(arr []int32) (string, []interface{}) {
	flags := make([]string, len(arr))
	args := make([]interface{}, len(arr))
	for i := 0; i < len(arr); i++ {
		flags[i] = "?"
		args[i] = arr[i]
	}
	return "Array[" + strings.Join(flags, ", ") + "]", args
}
