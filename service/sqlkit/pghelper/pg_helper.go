package pghelper

import (
	"github.com/Masterminds/squirrel"
	"strings"
)

// 生成sql中: sth in (select unnest(Array[?,?,?])) []interface{}
// arr必须不能空
// 注意使用时 args...
func genUnnestString(arr []string) (string, []interface{}) {
	flags := make([]string, len(arr))
	args := make([]interface{}, len(arr))
	for i := 0; i < len(arr); i++ {
		flags[i] = "?"
		args[i] = arr[i]
	}
	return " (select unnest(Array[" + strings.Join(flags, ", ") + "]))", args
}

func genUnnestInt(arr []int32) (string, []interface{}) {
	flags := make([]string, len(arr))
	args := make([]interface{}, len(arr))
	for i := 0; i < len(arr); i++ {
		flags[i] = "?"
		args[i] = arr[i]
	}
	return " (select unnest(Array[" + strings.Join(flags, ", ") + "]::int[]))", args
}

func genUnnestInt64(arr []int64) (string, []interface{}) {
	flags := make([]string, len(arr))
	args := make([]interface{}, len(arr))
	for i := 0; i < len(arr); i++ {
		flags[i] = "?"
		args[i] = arr[i]
	}
	return " (select unnest(Array[" + strings.Join(flags, ", ") + "]::int[]))", args
}

// 返回 Array[?,?,?]
func genArrayFlagString(arr []string) (string, []interface{}) {
	flags := make([]string, len(arr))
	args := make([]interface{}, len(arr))
	for i := 0; i < len(arr); i++ {
		flags[i] = "?"
		args[i] = arr[i]
	}
	return " Array[" + strings.Join(flags, ", ") + "]::varchar[]", args
}

func genArrayFlagInt(arr []int32) (string, []interface{}) {
	flags := make([]string, len(arr))
	args := make([]interface{}, len(arr))
	for i := 0; i < len(arr); i++ {
		flags[i] = "?"
		args[i] = arr[i]
	}
	return " Array[" + strings.Join(flags, ", ") + "]::int[]", args
}

// 封装到builder
func WhereUnnestString(builder squirrel.SelectBuilder, sqlPrefix string, arr []string) squirrel.SelectBuilder {
	flag, arg := genUnnestString(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}

func WhereUnnestInt(builder squirrel.SelectBuilder, sqlPrefix string, arr []int32) squirrel.SelectBuilder {
	flag, arg := genUnnestInt(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}

func WhereUnnestInt64(builder squirrel.SelectBuilder, sqlPrefix string, arr []int64) squirrel.SelectBuilder {
	flag, arg := genUnnestInt64(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}

func WhereArrayString(builder squirrel.SelectBuilder, sqlPrefix string, arr []string) squirrel.SelectBuilder {
	flag, arg := genArrayFlagString(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}

func WhereArrayInt(builder squirrel.SelectBuilder, sqlPrefix string, arr []int32) squirrel.SelectBuilder {
	flag, arg := genArrayFlagInt(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}

// 封装到builder
func WhereUnnestStringU(builder squirrel.UpdateBuilder, sqlPrefix string, arr []string) squirrel.UpdateBuilder {
	flag, arg := genUnnestString(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}

func WhereUnnestIntU(builder squirrel.UpdateBuilder, sqlPrefix string, arr []int32) squirrel.UpdateBuilder {
	flag, arg := genUnnestInt(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}

func WhereUnnestInt64U(builder squirrel.UpdateBuilder, sqlPrefix string, arr []int64) squirrel.UpdateBuilder {
	flag, arg := genUnnestInt64(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}

func WhereArrayStringU(builder squirrel.UpdateBuilder, sqlPrefix string, arr []string) squirrel.UpdateBuilder {
	flag, arg := genArrayFlagString(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}

func WhereArrayIntU(builder squirrel.UpdateBuilder, sqlPrefix string, arr []int32) squirrel.UpdateBuilder {
	flag, arg := genArrayFlagInt(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}
