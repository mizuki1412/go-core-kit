package sqlkit

import (
	"reflect"

	"github.com/jmoiron/sqlx"
	"github.com/mizuki1412/go-core-kit/v2/class"
	"github.com/mizuki1412/go-core-kit/v2/class/constraints"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/cli/tag"
	"github.com/spf13/cast"
)

func scanObjList[T any](dao SelectDao[T]) []*T {
	rows := dao.QueryRows()
	list := make([]*T, 0, 5)
	defer rows.Close()
	for rows.Next() {
		list = append(list, scanStruct[T](rows, dao.dataSource.Driver))
	}
	if dao.Cascade != nil {
		for i := range list {
			dao.Cascade(list[i])
		}
	}
	return list
}

func scanStruct[T any](rows *sqlx.Rows, driver string) *T {
	m := new(T)
	err := rows.StructScan(m)
	rv := reflect.ValueOf(m).Elem()
	rt := reflect.TypeOf(m).Elem()
	for i := 0; i < rv.NumField(); i++ {
		v := rv.Field(i)
		if v.Kind() == reflect.Struct {
			obj := v.Addr().Interface()
			// 处理 arr, 只针对 struct; 设置 dbdriver
			if vv, ok := obj.(constraints.SetDBDriverInterface); ok {
				vv.SetDBDriver(driver)
			}
			// 对decimal精度的处理
			precision := cast.ToInt32(rt.Field(i).Tag.Get(tag.DecimalPrecision.Name))
			if precision > 0 {
				if vv, ok := obj.(class.Decimal); ok {
					vv.Set(vv.Round(precision))
				}
				if vv, ok := obj.(*class.Decimal); ok {
					vv.Set(vv.Round(precision))
				}
			}
		}
	}
	if err != nil {
		panic(exception.New(err.Error()))
	}
	return m
}
