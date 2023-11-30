package sqlkit

import (
	"github.com/mizuki1412/go-core-kit/class/constraints"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"reflect"
	"strings"
)

// ModelMeta 获取model中的tablename和db fields
type ModelMeta struct {
	tableName   string
	keys        []ModelMetaKey
	logicDelKey ModelMetaKey
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
	Key string
	// 没有 escape 的 key
	OriKey  string
	RStruct reflect.StructField
	Primary bool
	Auto    bool
}

func (th ModelMetaKey) val(rv reflect.Value, driver string) any {
	var val any
	v := rv.FieldByName(th.RStruct.Name)
	if v.IsValid() {
		val = v.Interface()
	}
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return nil
	}
	// 改用 isValid() 判断结构体, 为了一致，必须值接收器
	if val != nil && (v.Kind() == reflect.Struct || v.Kind() == reflect.Pointer) {
		if valm, ok := val.(constraints.IsValidInterface); ok {
			if !valm.IsValid() {
				return nil
			}
			if v.Kind() == reflect.Struct {
				if vv, ok := v.Addr().Interface().(constraints.SetDBDriverInterface); ok {
					vv.SetDBDriver(driver)
					val = vv
				}
			} else {
				if vv, ok := v.Interface().(constraints.SetDBDriverInterface); ok {
					vv.SetDBDriver(driver)
					val = vv
				}
			}
		}
		//method := v.MethodByName("Value")
		//if !method.IsValid() {
		//	panic(exception.New("must add Value function or use value receiver: " + th.RStruct.Name))
		//}
		//if method.Call(nil)[0].Interface() == nil {
		//	val = nil
		//}
	}
	return val
}

// 用于存放model的解析数据： key：包路径+类名+驱动类型
var modelMetaCache = map[string]ModelMeta{}

// InitModelMeta obj should be elem
func (th ModelMeta) init(obj any) ModelMeta {
	if obj == nil {
		return ModelMeta{}
	}
	rt := reflect.TypeOf(obj)
	// 包路径+类名+驱动类型
	tk := rt.PkgPath() + "/" + rt.Name() + ":" + th.dateSource.Driver
	if tv, ok := modelMetaCache[tk]; ok {
		return tv
	}
	if rt.Kind() != reflect.Struct {
		panic(exception.New("dao model must struct"))
	}
	for i := 0; i < rt.NumField(); i++ {
		name := rt.Field(i).Tag.Get("db")
		if name == "" {
			continue
		}
		oriKey := name
		name = th.escapeName(name)
		key := ModelMetaKey{Key: name, OriKey: oriKey, RStruct: rt.Field(i)}
		// tableName; fetch once
		if th.tableName == "" {
			if t, ok := rt.Field(i).Tag.Lookup("table"); ok {
				th.tableName = t
			} else if t, ok := rt.Field(i).Tag.Lookup("tablename"); ok {
				// Deprecated
				th.tableName = t
			}
		}
		if t, ok := rt.Field(i).Tag.Lookup("logicDel"); ok && t == "true" {
			th.logicDelKey = key
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
		th.keys = append(th.keys, key)
	}
	if th.tableName == "" {
		panic(exception.New("model meta tableName error"))
	}
	// 处理
	for _, e := range th.keys {
		th.allSelectColumns = append(th.allSelectColumns, e.Key)
		if e.Primary {
			th.allPKs = append(th.allPKs, e)
		}
		if !e.Primary && !e.Auto {
			th.allUpdateKeys = append(th.allUpdateKeys, e)
		}
		if !e.Auto {
			th.allInsertKeys = append(th.allInsertKeys, e)
		}
	}
	if th.logicDelKey.OriKey != "" {
		th.allSelectColumns = append(th.allSelectColumns, th.logicDelKey.Key)
		th.allUpdateKeys = append(th.allUpdateKeys, th.logicDelKey)
	}
	modelMetaCache[tk] = th
	return th
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
	ret := make([]string, len(name))
	for i := 0; i < len(name); i++ {
		ret[i] = th.dateSource.escapeName(name[i])
	}
	return ret
}
func (th ModelMeta) escapeName(name string) string {
	return th.dateSource.escapeName(name)
}
