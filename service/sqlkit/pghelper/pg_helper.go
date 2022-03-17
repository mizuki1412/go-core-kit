package pghelper

import (
	"github.com/Masterminds/squirrel"
	"strings"
)

// 生成sql中: sth in (select unnest(Array[?,?,?])) []any
// arr必须不能空
// 注意使用时 args...
func GenUnnestString(arr []string) (string, []any) {
	flags := make([]string, len(arr))
	args := make([]any, len(arr))
	for i := 0; i < len(arr); i++ {
		flags[i] = "?"
		args[i] = arr[i]
	}
	return " (select unnest(Array[" + strings.Join(flags, ", ") + "]))", args
}

func GenUnnestInt(arr []int32) (string, []any) {
	flags := make([]string, len(arr))
	args := make([]any, len(arr))
	for i := 0; i < len(arr); i++ {
		flags[i] = "?"
		args[i] = arr[i]
	}
	return " (select unnest(Array[" + strings.Join(flags, ", ") + "]::int[]))", args
}

func GenUnnestInt64(arr []int64) (string, []any) {
	flags := make([]string, len(arr))
	args := make([]any, len(arr))
	for i := 0; i < len(arr); i++ {
		flags[i] = "?"
		args[i] = arr[i]
	}
	return " (select unnest(Array[" + strings.Join(flags, ", ") + "]::int[]))", args
}

// 返回 Array[?,?,?]
func GenArrayFlagString(arr []string) (string, []any) {
	flags := make([]string, len(arr))
	args := make([]any, len(arr))
	for i := 0; i < len(arr); i++ {
		flags[i] = "?"
		args[i] = arr[i]
	}
	return " Array[" + strings.Join(flags, ", ") + "]::varchar[]", args
}

func GenArrayFlagInt(arr []int32) (string, []any) {
	flags := make([]string, len(arr))
	args := make([]any, len(arr))
	for i := 0; i < len(arr); i++ {
		flags[i] = "?"
		args[i] = arr[i]
	}
	return " Array[" + strings.Join(flags, ", ") + "]::int[]", args
}

// extend->'key' @> '[3]'::jsonb
func GenJsonArrayFlagInt(arr []int32) (string, []any) {
	flags := make([]string, len(arr))
	args := make([]any, len(arr))
	for i := 0; i < len(arr); i++ {
		flags[i] = "?"
		args[i] = arr[i]
	}
	return " @> '" + strings.Join(flags, ", ") + "'::jsonb", args
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

func WhereUnnestInt64(builder squirrel.SelectBuilder, sqlPrefix string, arr []int64) squirrel.SelectBuilder {
	flag, arg := GenUnnestInt64(arr)
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

func WhereJsonArrayInt(builder squirrel.SelectBuilder, sqlPrefix string, arr []int32) squirrel.SelectBuilder {
	flag, arg := GenJsonArrayFlagInt(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}

// update, 封装到builder
func WhereUnnestStringU(builder squirrel.UpdateBuilder, sqlPrefix string, arr []string) squirrel.UpdateBuilder {
	flag, arg := GenUnnestString(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}

func WhereUnnestIntU(builder squirrel.UpdateBuilder, sqlPrefix string, arr []int32) squirrel.UpdateBuilder {
	flag, arg := GenUnnestInt(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}

func WhereUnnestInt64U(builder squirrel.UpdateBuilder, sqlPrefix string, arr []int64) squirrel.UpdateBuilder {
	flag, arg := GenUnnestInt64(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}

func WhereArrayStringU(builder squirrel.UpdateBuilder, sqlPrefix string, arr []string) squirrel.UpdateBuilder {
	flag, arg := GenArrayFlagString(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}

func WhereArrayIntU(builder squirrel.UpdateBuilder, sqlPrefix string, arr []int32) squirrel.UpdateBuilder {
	flag, arg := GenArrayFlagInt(arr)
	return builder.Where(sqlPrefix+flag, arg...)
}
