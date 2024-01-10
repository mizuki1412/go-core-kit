package sqlkit

import (
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/v2/class/const/sqlconst"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"strings"
)

func pgArray(arr any) (string, []any) {
	var suffix string
	var args []any
	var flags []string
	switch arr.(type) {
	case []int:
		suffix = "int[]"
		arr := arr.([]int)
		flags = make([]string, len(arr))
		args = make([]any, len(arr))
		for i := 0; i < len(flags); i++ {
			flags[i] = "?"
			args[i] = arr[i]
		}
	case []int8:
		suffix = "int[]"
		arr := arr.([]int8)
		flags = make([]string, len(arr))
		args = make([]any, len(arr))
		for i := 0; i < len(flags); i++ {
			flags[i] = "?"
			args[i] = arr[i]
		}
	case []int16:
		suffix = "int[]"
		arr := arr.([]int16)
		flags = make([]string, len(arr))
		args = make([]any, len(arr))
		for i := 0; i < len(flags); i++ {
			flags[i] = "?"
			args[i] = arr[i]
		}
	case []int32:
		suffix = "int[]"
		arr := arr.([]int32)
		flags = make([]string, len(arr))
		args = make([]any, len(arr))
		for i := 0; i < len(flags); i++ {
			flags[i] = "?"
			args[i] = arr[i]
		}
	case []int64:
		suffix = "bigint[]"
		arr := arr.([]int64)
		flags = make([]string, len(arr))
		args = make([]any, len(arr))
		for i := 0; i < len(flags); i++ {
			flags[i] = "?"
			args[i] = arr[i]
		}
	case []float32:
		suffix = "decimal[]"
		arr := arr.([]float32)
		flags = make([]string, len(arr))
		args = make([]any, len(arr))
		for i := 0; i < len(flags); i++ {
			flags[i] = "?"
			args[i] = arr[i]
		}
	case []float64:
		suffix = "decimal[]"
		arr := arr.([]float64)
		flags = make([]string, len(arr))
		args = make([]any, len(arr))
		for i := 0; i < len(flags); i++ {
			flags[i] = "?"
			args[i] = arr[i]
		}
	case []string:
		suffix = "varchar[]"
		arr := arr.([]string)
		flags = make([]string, len(arr))
		args = make([]any, len(arr))
		for i := 0; i < len(flags); i++ {
			flags[i] = "?"
			args[i] = arr[i]
		}
	default:
		panic(exception.New("pgArray params not supported"))
	}
	// 用{} 有错误：invalid input syntax for type integer
	return "ARRAY[" + strings.Join(flags, ",") + "]::" + suffix, args
}

// 占位符
func placeholder(driver string) squirrel.PlaceholderFormat {
	switch driver {
	case sqlconst.Postgres:
		return squirrel.Dollar
	case sqlconst.Oracle:
		return squirrel.Colon
	default:
		return squirrel.Question
	}
}
