package pghelper

import (
	"github.com/Masterminds/squirrel"
	"strings"
)

// 生成sql中: in (select unnest(Array[?,?,?])) []interface{}
// arr必须不能空
// 注意使用时 args... todo 将会转为内部函数
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
	return "(select unnest(Array[" + strings.Join(flags, ", ") + "]::int[]))", args
}

// 返回 Array[?,?,?] todo 将改为内部函数
func GenArrayFlagString(arr []string) (string, []interface{}) {
	flags := make([]string, len(arr))
	args := make([]interface{}, len(arr))
	for i := 0; i < len(arr); i++ {
		flags[i] = "?"
		args[i] = arr[i]
	}
	return "Array[" + strings.Join(flags, ", ") + "]::varchar[]", args
}

func GenArrayFlagInt(arr []int32) (string, []interface{}) {
	flags := make([]string, len(arr))
	args := make([]interface{}, len(arr))
	for i := 0; i < len(arr); i++ {
		flags[i] = "?"
		args[i] = arr[i]
	}
	return "Array[" + strings.Join(flags, ", ") + "]::int[]", args
}

// 封装到builder
func WhereUnnestString(builder squirrel.SelectBuilder, sqlPrefix string, arr []string) squirrel.SelectBuilder {
	flag, arg := GenUnnestString(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}

func WhereUnnestInt(builder squirrel.SelectBuilder, sqlPrefix string, arr []int32) squirrel.SelectBuilder {
	flag, arg := GenUnnestInt(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}

func WhereArrayString(builder squirrel.SelectBuilder, sqlPrefix string, arr []string) squirrel.SelectBuilder {
	flag, arg := GenArrayFlagString(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}

func WhereArrayInt(builder squirrel.SelectBuilder, sqlPrefix string, arr []int32) squirrel.SelectBuilder {
	flag, arg := GenArrayFlagInt(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}
