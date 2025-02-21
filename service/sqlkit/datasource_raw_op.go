package sqlkit

import (
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/library/timekit"
	"github.com/spf13/cast"
	"time"
)

/**
* 用于简单场景的sql，无规定model
 */

func (ds *DataSource) QueryOne(sql string, args ...any) any {
	rows := ds.Query(sql, args)
	defer rows.Close()
	for rows.Next() {
		ret, err := rows.SliceScan()
		if err != nil {
			panic(exception.New(err.Error()))
		}
		return ret[0]
	}
	return 0
}

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

// FormatRawMap 原始查出的map[string]any 转为 map[string]string 用于sql语句中
// 注意防注入
func (ds *DataSource) FormatRawMap(rows map[string]any) map[string]string {
	res := make(map[string]string)
	for key, val := range rows {
		res[key] = ds.FormatRawValue(val)
	}
	return res
}

func (ds *DataSource) FormatRawValue(val any) string {
	if val == nil {
		return "null"
	}
	// todo 其他情况
	switch val.(type) {
	case string:
		return "'" + val.(string) + "'"
	case []uint8:
		return "'" + string(val.([]uint8)) + "'"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return cast.ToString(val)
	case float32, float64:
		return cast.ToString(val)
	case bool:
		return cast.ToString(val)
	case time.Time:
		t := val.(time.Time)
		if t.IsZero() {
			return "null"
		}
		return "'" + t.In(timekit.GetLocation()).Format(timekit.TimeLayout) + "'"
	}
	return "'" + cast.ToString(val) + "'"
}
