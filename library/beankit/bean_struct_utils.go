package beankit

import (
	"encoding/json"
	"reflect"
)

// Struct2Map 用标签中的json来指定map的key值, 保持原有对象中的数据类型
// bean 需要是指针形式
func Struct2Map(bean interface{}) map[string]interface{} {
	rt := reflect.TypeOf(bean).Elem()
	rv := reflect.ValueOf(bean).Elem()
	ret := map[string]interface{}{}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldV := rv.Field(i)
		key := field.Tag.Get("json")
		// 按json解析的规则，如果是null就pass
		v, ok := fieldV.Interface().(json.Marshaler)
		if ok {
			b, _ := v.MarshalJSON()
			if string(b) == "null" {
				continue
			}
		}
		ret[key] = fieldV.Interface()
	}
	return ret
}
