package sqlkit

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"reflect"
)

func reflectModel(dest any) (reflect.Type, reflect.Value) {
	rt := reflect.TypeOf(dest)
	if rt.Kind() != reflect.Pointer {
		panic(exception.New("insert param should be pointer"))
	}
	rt = rt.Elem()
	rv := reflect.ValueOf(dest).Elem()
	return rt, rv
}

// 多或单主键
func getPKs(rt reflect.Type, rv reflect.Value) map[string]any {
	pks := map[string]any{}
	for i := 0; i < rt.NumField(); i++ {
		if t, ok := rt.Field(i).Tag.Lookup("pk"); ok {
			if t == "true" {
				name := rt.Field(i).Tag.Get("db")
				if name == "" {
					panic(exception.New("field "+rt.Field(i).Name+" no db tag", 2))
				}
				pks[name] = rv.Field(i).Interface()
			}
		}
	}
	if len(pks) == 0 {
		panic(exception.New("未设置pk", 2))
	}
	return pks
}

func (dao *Dao[T]) getTable(rt reflect.Type) string {
	var tname string
	for i := 0; i < rt.NumField(); i++ {
		if t, ok := rt.Field(i).Tag.Lookup("table"); ok {
			tname = t
			break
		} else if t, ok := rt.Field(i).Tag.Lookup("tablename"); ok {
			tname = t
			break
		}
	}
	if tname == "" {
		panic(exception.New("table name 未设置", 2))
	}
	return dao.DataSource.getDecoSchema() + tname
}
