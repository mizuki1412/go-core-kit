package sqlkit

import (
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/v2/class"
	"github.com/mizuki1412/go-core-kit/v2/class/const/sqlconst"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
)

// 生成pg的array表达式
func pgArray(arr any) (string, []any) {
	var suffix string
	var args []any
	var flags []string
	switch arr.(type) {
	case []int:
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

func normalArray(arr any) (string, []any) {
	var args []any
	var flags []string
	switch arr.(type) {
	case []int:
		arr := arr.([]int)
		flags = make([]string, len(arr))
		args = make([]any, len(arr))
		for i := 0; i < len(flags); i++ {
			flags[i] = "?"
			args[i] = arr[i]
		}
	case []int8:
		arr := arr.([]int8)
		flags = make([]string, len(arr))
		args = make([]any, len(arr))
		for i := 0; i < len(flags); i++ {
			flags[i] = "?"
			args[i] = arr[i]
		}
	case []int16:
		arr := arr.([]int16)
		flags = make([]string, len(arr))
		args = make([]any, len(arr))
		for i := 0; i < len(flags); i++ {
			flags[i] = "?"
			args[i] = arr[i]
		}
	case []int32:
		arr := arr.([]int32)
		flags = make([]string, len(arr))
		args = make([]any, len(arr))
		for i := 0; i < len(flags); i++ {
			flags[i] = "?"
			args[i] = arr[i]
		}
	case []int64:
		arr := arr.([]int64)
		flags = make([]string, len(arr))
		args = make([]any, len(arr))
		for i := 0; i < len(flags); i++ {
			flags[i] = "?"
			args[i] = arr[i]
		}
	case []float32:
		arr := arr.([]float32)
		flags = make([]string, len(arr))
		args = make([]any, len(arr))
		for i := 0; i < len(flags); i++ {
			flags[i] = "?"
			args[i] = arr[i]
		}
	case []float64:
		arr := arr.([]float64)
		flags = make([]string, len(arr))
		args = make([]any, len(arr))
		for i := 0; i < len(flags); i++ {
			flags[i] = "?"
			args[i] = arr[i]
		}
	case []string:
		arr := arr.([]string)
		flags = make([]string, len(arr))
		args = make([]any, len(arr))
		for i := 0; i < len(flags); i++ {
			flags[i] = "?"
			args[i] = arr[i]
		}
	default:
		panic(exception.New("normalArray params not supported"))
	}
	return "(" + strings.Join(flags, ",") + ")", args
}

// 占位符
func placeholder(driver string) squirrel.PlaceholderFormat {
	switch driver {
	case sqlconst.Postgres, sqlconst.Kingbase:
		return squirrel.Dollar
	case sqlconst.Oracle:
		return squirrel.Colon
	default:
		return squirrel.Question
	}
}

// args中部分值转换
func argsWrap(driver string, args []any) []any {
	// todo 其他值类型
	new_args := make([]any, 0, len(args))
	for _, e := range args {
		n := e
		switch e.(type) {
		case class.Time:
			v := e.(class.Time)
			if driver == sqlconst.Sqlite3 || sqlconst.IsTaos(driver) {
				n = v.UnixMill()
			} else {
				n = v.Time
			}
		case time.Time:
			v := e.(time.Time)
			if driver == sqlconst.Sqlite3 || sqlconst.IsTaos(driver) {
				n = v.UnixMilli()
			}
		case string, class.String:
			// taos未对其中的'字符转义, 但在insert中转义了？todo
			if sqlconst.IsTaos(driver) {
				var ee string
				if ev, ok := e.(string); ok {
					ee = ev
				} else {
					ee = e.(class.String).String
				}
				if strings.Contains(ee, "'") {
					n = strings.ReplaceAll(ee, "'", "''")
				}
			}
		}
		new_args = append(new_args, n)
	}
	return new_args
}

func handlePlaceholderInWhere(driver string, pred any, args ...any) any {
	_, ok := pred.(string)
	if ok && sqlconst.IsTaos(driver) && len(args) > 0 {
		noString := true
		for _, e := range args {
			if _, ok := e.(string); ok {
				noString = false
				break
			}
		}
		if noString {
			return pred
		}
		p := pred.(string)
		index := 0
		pa := ""
		ps := strings.Split(p, "?")
		for _, e := range ps {
			if len(args) <= index {
				break
			}
			if _, ok := args[index].(string); ok {
				pa += e + "'?'"
			} else {
				pa += e + "?"
			}
			index++
		}
		return pa
	}
	return pred
}
