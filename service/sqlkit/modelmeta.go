package sqlkit

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"reflect"
	"strings"
)

// ModelMeta 获取model中的tablename和db fields
type ModelMeta struct {
	tableName   string
	keys        []ModelMetaKey
	logicDelKey string
	dateSource  *DataSource
	// 处理后的 keys array
	// 用于 select 的 全量columns
	allSelectColumns []string
	allInsertKeys    []ModelMetaKey
	allUpdateKeys    []ModelMetaKey
	allPKs           []ModelMetaKey
}

// ModelMetaKey 除 logicdelete 外的 keys
type ModelMetaKey struct {
	// escape 后的 key
	Key     string
	RStruct reflect.StructField
	Primary bool
	Auto    bool
}

func (th ModelMetaKey) val(rv reflect.Value) any {
	var val any
	v := rv.FieldByName(th.RStruct.Name)
	// 判断field是否指针
	if th.RStruct.Type.Kind() == reflect.Pointer && v.Elem().IsValid() {
		val = v.Elem().Interface()
	} else if th.RStruct.Type.Kind() != reflect.Pointer {
		val = v.Interface()
	}
	if val != nil {
		method := v.MethodByName("Value")
		if !method.IsValid() {
			panic(exception.New("must add Value function"))
		}
	}
	return val
}

// InitModelMeta obj should be elem
func (th ModelMeta) init(obj any) ModelMeta {
	meta := ModelMeta{}
	rt := reflect.TypeOf(obj)
	if rt.Kind() != reflect.Struct {
		panic(exception.New("dao model must struct"))
	}
	for i := 0; i < rt.NumField(); i++ {
		name := rt.Field(i).Tag.Get("db")
		if name == "" {
			continue
		}
		name = th.escapeName(name)
		key := ModelMetaKey{Key: name, RStruct: rt.Field(i)}
		// tableName; fetch once
		if meta.tableName == "" {
			if t, ok := rt.Field(i).Tag.Lookup("table"); ok {
				meta.tableName = t
			} else if t, ok := rt.Field(i).Tag.Lookup("tablename"); ok {
				// Deprecated
				meta.tableName = t
			}
		}
		if t, ok := rt.Field(i).Tag.Lookup("logicDel"); ok && t == "true" {
			meta.logicDelKey = name
			continue
		}
		// pk
		if t, ok := rt.Field(i).Tag.Lookup("pk"); ok && t == "true" {
			key.Primary = true
		}
		if t, ok := rt.Field(i).Tag.Lookup("auto"); ok && t == "true" {
			key.Auto = true
		}
		// Deprecated
		if t, ok := rt.Field(i).Tag.Lookup("autoincrement"); ok && t == "true" {
			key.Auto = true
		}
		meta.keys = append(meta.keys, key)
	}
	if meta.tableName == "" {
		panic(exception.New("model meta tableName error"))
	}
	// 处理
	for _, e := range meta.keys {
		meta.allSelectColumns = append(meta.allSelectColumns, e.Key)
		if e.Primary {
			meta.allPKs = append(meta.allPKs, e)
		}
		if !e.Primary && !e.Auto {
			meta.allUpdateKeys = append(meta.allUpdateKeys, e)
		}
		if !e.Auto {
			meta.allInsertKeys = append(meta.allInsertKeys, e)
		}
	}
	if meta.logicDelKey != "" {
		meta.allSelectColumns = append(meta.allSelectColumns, meta.logicDelKey)
	}
	return meta
}

func (th ModelMeta) getSelectColumns(excludes ...string) []string {
	return th.getSelectColumnsWithPrefix("", excludes...)
}

func (th ModelMeta) getSelectColumnsWithPrefix(prefix string, excludes ...string) []string {
	if prefix != "" {
		prefix += "."
	}
	arr := make([]string, 0, len(th.allSelectColumns))
	if len(excludes) > 0 {
		ex := strings.Join(excludes, ";")
		ex += ";"
		for _, e := range th.allSelectColumns {
			if !strings.Contains(ex, e+";") {
				arr = append(arr, prefix+e)
			}
		}
	} else if prefix != "" {
		for _, e := range th.allSelectColumns {
			arr = append(arr, prefix+e)
		}
	} else {
		for _, e := range th.allSelectColumns {
			arr = append(arr, e)
		}
	}
	if len(arr) == 0 {
		panic(exception.New("sql columns 不能为空"))
	}
	return arr
}

// getTable alias 可以包括table别名
func (th ModelMeta) getTable(alias ...string) string {
	if len(alias) > 0 {
		return th.dateSource.decoTableName(th.tableName) + " AS " + alias[0]
	} else {
		return th.dateSource.decoTableName(th.tableName)
	}
}

func (th ModelMeta) escapeNames(name []string) []string {
	if len(name) == 0 {
		panic(exception.New("modelmeta escapename nil"))
	}
	ret := make([]string, 0, len(name))
	for _, e := range name {
		ret = append(ret, th.dateSource.escapeName(e))
	}
	return ret
}
func (th ModelMeta) escapeName(name string) string {
	return th.dateSource.escapeName(name)
}
