package influxkit

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/init/configkey"
	"github.com/mizuki1412/go-core-kit/library/httpkit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"strings"
)

func QueryDefaultDB(sql string) []map[string]interface{} {
	return queryResult(configkit.GetStringD(configkey.InfluxDBName), sql)
}

func QueryWithPrefix(prefix, sql string) []map[string]interface{} {
	queryResult(prefix+configkit.GetStringD(configkey.InfluxDBName), sql)
	return nil
}

func QueryWithDBName(dbName, sql string) []map[string]interface{} {
	return queryResult(dbName, sql)
}

func QueryMultiDefaultDB(sql []string) [][]map[string]interface{} {
	return queryMultiResult(configkit.GetStringD(configkey.InfluxDBName), sql)
}

func QueryMultiWithDBName(dbName string, sql []string) [][]map[string]interface{} {
	return queryMultiResult(dbName, sql)
}

func url(action, dbName string) string {
	url := configkit.GetStringD(configkey.InfluxURL)
	if url == "" {
		panic(exception.New("influx url is null"))
	}
	var params string
	if dbName != "" {
		params += "db=" + dbName
	}
	if configkit.GetStringD(configkey.InfluxUser) != "" {
		if params != "" {
			params += "&"
		}
		params += "u=" + configkit.GetStringD(configkey.InfluxUser) + "&p=" + configkit.GetStringD(ConfigKeyInfluxPwd)
	}
	if params == "" {
		return configkit.GetStringD(configkey.InfluxURL) + "/" + action
	} else {
		return configkit.GetStringD(configkey.InfluxURL) + "/" + action + "?" + params
	}
}

func queryResult(dbName, sql string) []map[string]interface{} {
	ret, code := httpkit.Request(httpkit.Req{
		Url: url("query", dbName),
		FormData: map[string]string{
			"epoch": "ms",
			"q":     sql,
		},
	})
	err := gjson.Get(ret, "error").String()
	if err != "" {
		panic(exception.New("influx query error: " + err))
	}
	if code > 300 {
		panic(exception.New("influx query error: " + cast.ToString(code)))
	}
	results := gjson.Get(ret, "results").Array()
	if len(results) > 0 {
		series := results[0].Map()["series"].Array()
		if len(series) > 0 {
			serie := series[0]
			columns := serie.Map()["columns"].Array()
			values := serie.Map()["values"].Array()
			if len(columns) > 0 && len(values) > 0 {
				list := make([]map[string]interface{}, len(values))
				for i, v := range values {
					e := map[string]interface{}{}
					for ii, vv := range columns {
						e[vv.String()] = v.Array()[ii].Value()
					}
					// 不用append，直接赋值
					list[i] = e
				}
				return list
			}
		}
	}
	return []map[string]interface{}{}
}

// 数组中可能出现nil
func queryMultiResult(dbName string, sql []string) [][]map[string]interface{} {
	sqls := strings.Join(sql, ";")
	ret, code := httpkit.Request(httpkit.Req{
		Url: url("query", dbName),
		FormData: map[string]string{
			"epoch": "ms",
			"q":     sqls,
		},
	})
	err := gjson.Get(ret, "error").String()
	if err != "" {
		panic(exception.New("influx query error: " + err))
	}
	if code > 300 {
		panic(exception.New("influx query error: " + cast.ToString(code)))
	}
	results := gjson.Get(ret, "results").Array()
	data := make([][]map[string]interface{}, 0, len(results))
	for _, result := range results {
		series := result.Map()["series"].Array()
		if len(series) > 0 {
			serie := series[0]
			columns := serie.Map()["columns"].Array()
			values := serie.Map()["values"].Array()
			if len(columns) > 0 && len(values) > 0 {
				list := make([]map[string]interface{}, len(values))
				for i, v := range values {
					e := map[string]interface{}{}
					for ii, vv := range columns {
						e[vv.String()] = v.Array()[ii].Value()
					}
					// 不用append，直接赋值
					list[i] = e
				}
				data = append(data, list)
			}
		}
	}
	return data
}

func CreateDB(name string) {
	ret, code := httpkit.Request(httpkit.Req{
		Url: url("query", ""),
		FormData: map[string]string{
			"q": "create database " + name,
		},
	})
	err := gjson.Get(ret, "error").String()
	if err != "" {
		panic(exception.New("influx query error: " + err))
	}
	if code > 300 {
		panic(exception.New("influx query error: " + cast.ToString(code)))
	}
}

// sql: dv_x key1=1,key2="val2" timestamp
func WriteDefaultDB(sql string) {
	writeData(configkit.GetStringD(configkey.InfluxDBName), sql)
}

func WriteWithDBName(dbName, sql string) {
	writeData(dbName, sql)
}

func writeData(dbName, sql string) {
	ret, code := httpkit.Request(httpkit.Req{
		Url:        url("write", dbName),
		BinaryData: []byte(sql),
	})
	err := gjson.Get(ret, "error").String()
	if err != "" {
		panic(exception.New("influx query error: " + err))
	}
	if code > 300 {
		panic(exception.New("influx query error: " + cast.ToString(code)))
	}
}

// 用于insert时或query时，val的装饰转换
func Decorate(val interface{}) string {
	if val == nil {
		return ""
	}
	switch val.(type) {
	case string:
		v := val.(string)
		v = strings.ReplaceAll(v, "\"", "\\\"")
		return "\"" + v + "\""
	case class.String:
		v := val.(class.String).String
		v = strings.ReplaceAll(v, "\"", "\\\"")
		return "\"" + v + "\""
	case *class.String:
		v := val.(*class.String).String
		v = strings.ReplaceAll(v, "\"", "\\\"")
		return "\"" + v + "\""
	case class.Decimal:
		return val.(class.Decimal).Decimal.String()
	case *class.Decimal:
		return val.(*class.Decimal).Decimal.String()
	case class.Int32:
		return cast.ToString(val.(class.Int32).Int32)
	case *class.Int32:
		return cast.ToString(val.(*class.Int32).Int32)
	case class.Int64:
		return cast.ToString(val.(class.Int64).Int64)
	case *class.Int64:
		return cast.ToString(val.(*class.Int64).Int64)
	default:
		return cast.ToString(val)
	}
}
