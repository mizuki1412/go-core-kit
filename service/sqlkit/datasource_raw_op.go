package sqlkit

import (
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/spf13/cast"
)

/**
* 用于简单场景的sql，无规定model
 */

func (ds *DataSource) QueryOneNumber(sql string, args ...any) int64 {
	rows := ds.Query(sql, args)
	defer rows.Close()
	for rows.Next() {
		ret, err := rows.SliceScan()
		if err != nil {
			panic(exception.New(err.Error()))
		}
		return cast.ToInt64(ret[0])
	}
	return 0
}

func (ds *DataSource) QueryOneMap(sql string, args ...any) map[string]any {
	rows := ds.Query(sql, args)
	defer rows.Close()
	for rows.Next() {
		m := map[string]any{}
		err := rows.MapScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		return m
	}
	return nil
}

func (ds *DataSource) QueryOneString(sql string, args ...any) string {
	rows := ds.Query(sql, args)
	defer rows.Close()
	for rows.Next() {
		ret, err := rows.SliceScan()
		if err != nil {
			panic(exception.New(err.Error()))
		}
		return cast.ToString(ret[0])
	}
	return ""
}

func (ds *DataSource) QueryListMap(sql string, args ...any) []map[string]any {
	rows := ds.Query(sql, args)
	defer rows.Close()
	list := make([]map[string]any, 0, 5)
	for rows.Next() {
		m := map[string]any{}
		err := rows.MapScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		list = append(list, m)
	}
	return list
}

func (ds *DataSource) QueryListString(sql string, args ...any) []string {
	rows := ds.Query(sql, args)
	defer rows.Close()
	list := make([]string, 0, 5)
	for rows.Next() {
		ret, err := rows.SliceScan()
		if err != nil {
			panic(exception.New(err.Error()))
		}
		list = append(list, cast.ToString(ret[0]))
	}
	return list
}
