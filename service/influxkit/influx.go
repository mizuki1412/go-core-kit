package influxkit

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/httpkit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
)

func QueryDefaultDB(sql string) []map[string]interface{} {
	return queryResult(configkit.GetStringD(ConfigKeyInfluxDBName), sql)
}

func QueryWithPrefix(prefix, sql string) []map[string]interface{} {
	queryResult(prefix+configkit.GetStringD(ConfigKeyInfluxDBName), sql)
	return nil
}

func QueryWithDBName(dbName, sql string) []map[string]interface{} {
	return queryResult(dbName, sql)
}

func QueryMultiDefaultDB(sql string) [][]map[string]interface{} {
	// todo
	return nil
}

func url(action, dbName string) string {
	url := configkit.GetStringD(ConfigKeyInfluxURL)
	if url == "" {
		panic(exception.New("influx url is null"))
	}
	var params string
	if dbName != "" {
		params += "db=" + dbName
	}
	if configkit.GetStringD(ConfigKeyInfluxUser) != "" {
		if params != "" {
			params += "&"
		}
		params += "u=" + configkit.GetStringD(ConfigKeyInfluxUser) + "&p=" + configkit.GetStringD(ConfigKeyInfluxPwd)
	}
	if params == "" {
		return configkit.GetStringD(ConfigKeyInfluxURL) + "/" + action
	} else {
		return configkit.GetStringD(ConfigKeyInfluxURL) + "/" + action + "?" + params
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

// sql: dv_x key1=val1,key2=val2 timestamp
func WriteDefaultDB(sql string) {
	writeData(configkit.GetStringD(ConfigKeyInfluxDBName), sql)
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
