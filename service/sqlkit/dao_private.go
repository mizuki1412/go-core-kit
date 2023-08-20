package sqlkit

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"reflect"
)

func (dao Dao[T]) getTable(rt reflect.Type) string {
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
	return dao.dataSource.decoTableName(tname)
}
