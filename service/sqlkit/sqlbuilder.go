package sqlkit

import (
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"reflect"
	"strings"
)

func Builder() squirrel.StatementBuilderType {
	connector()
	if driver == "postgres" {
		return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	} else {
		// todo 未处理oracle
		return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	}
}

// ModelMeta 获取model中的tablename和db fields
type ModelMeta struct {
	TableName string
	Fields    []string
}

// InitModelMeta dest should be elem
func InitModelMeta(obj any) *ModelMeta {
	meta := &ModelMeta{}
	rt := reflect.TypeOf(obj).Elem()
	for i := 0; i < rt.NumField(); i++ {
		if t, ok := rt.Field(i).Tag.Lookup("tablename"); ok {
			meta.TableName = t
		}
		if t, ok := rt.Field(i).Tag.Lookup("db"); ok {
			meta.Fields = append(meta.Fields, t)
		}
	}
	if meta.TableName == "" {
		panic(exception.New("model meta tableName error"))
	}
	return meta
}

func (th ModelMeta) GetColumns(excludes ...string) string {
	return th.GetColumnsWithPrefix("", excludes...)
}

func (th ModelMeta) GetColumnsWithPrefix(prefix string, excludes ...string) string {
	if prefix != "" {
		prefix += "."
	}
	var arr = th.Fields
	if len(excludes) > 0 {
		arr = make([]string, 0, len(th.Fields))
		ex := strings.Join(excludes, ";")
		ex += ";"
		for _, e := range th.Fields {
			if !strings.Contains(ex, e+";") {
				arr = append(arr, prefix+e)
			}
		}
	} else if prefix != "" {
		arr = make([]string, 0, len(th.Fields))
		for _, e := range th.Fields {
			arr = append(arr, prefix+e)
		}
	}
	if len(arr) == 0 {
		panic(exception.New("sql columns 不能为空"))
	}
	return strings.Join(arr, ",")
}

// GetTableName alias 可以包括table别名
func (th ModelMeta) GetTableName(schema string, alias ...string) string {
	if len(alias) > 0 {
		return GetSchemaTable(schema, th.TableName) + " " + alias[0]
	} else {
		return GetSchemaTable(schema, th.TableName)
	}
}

// GetSchemaTable name可能是表名，带join信息
func GetSchemaTable(schema string, name string) string {
	var schema0 string
	if schema != "" {
		schema0 = schema
	} else if driver == "postgres" {
		schema0 = SchemaDefault
	} else {
		schema0 = ""
	}
	if schema0 == "" {
		return name
	} else {
		return schema0 + "." + name
	}
}
