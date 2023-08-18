package sqlkit

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"strings"
)

// ModelMeta 获取model中的tablename和db fields
type ModelMeta struct {
	TableName string
	Fields    []string
}

// InitModelMeta obj should be elem
func InitModelMeta(obj any) *ModelMeta {
	meta := &ModelMeta{}
	rt, _ := reflectModel(obj)
	for i := 0; i < rt.NumField(); i++ {
		if t, ok := rt.Field(i).Tag.Lookup("table"); ok {
			meta.TableName = t
		} else if t, ok := rt.Field(i).Tag.Lookup("tablename"); ok {
			// Deprecated
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

// todo 加引号装饰
func (th *ModelMeta) GetColumns(excludes ...string) []string {
	return th.GetColumnsWithPrefix("", excludes...)
}

func (th *ModelMeta) GetColumnsWithPrefix(prefix string, excludes ...string) []string {
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
	return arr
}

// GetTableName alias 可以包括table别名
func (th *ModelMeta) GetTableName(alias ...string) string {
	if len(alias) > 0 {
		return GetSchemaTable(schema, th.TableName) + " " + alias[0]
	} else {
		return GetSchemaTable(schema, th.TableName)
	}
}
