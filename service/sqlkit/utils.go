package sqlkit

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mizuki1412/go-core-kit/class/constraints"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"reflect"
)

func logReqSqlInfo(sql string, args []any) string {
	return fmt.Sprintf(`
==> %s
==> %s`, sql, jsonkit.ToString(args))
}

func logResSqlInfo(rows int64) string {
	return fmt.Sprintf(`
<== rows: %d`, rows)
}

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
	// 处理 arr, 只针对 struct; 设置 dbdriver
	rv := reflect.ValueOf(m).Elem()
	for i := 0; i < rv.NumField(); i++ {
		v := rv.Field(i)
		if v.Kind() == reflect.Struct {
			if vv, ok := v.Addr().Interface().(constraints.SetDBDriverInterface); ok {
				vv.SetDBDriver(driver)
			}
		}
	}
	err := rows.StructScan(m)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	return m
}
