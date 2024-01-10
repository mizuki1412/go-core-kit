package beankit

import (
	"encoding/json"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/spf13/cast"
	"reflect"
	"strings"
)

// Struct2Map 用标签中的json来指定map的key值, 保持原有对象中的数据类型
// bean 需要是指针形式
func Struct2Map(bean any) map[string]any {
	rt := reflect.TypeOf(bean).Elem()
	rv := reflect.ValueOf(bean).Elem()
	ret := map[string]any{}
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

// Map2Struct 用于将map转换成结构体 注意可用的数据类型
func Map2Struct(m map[string]any, bean any) {
	rt := reflect.TypeOf(bean).Elem()
	rv := reflect.ValueOf(bean).Elem()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldV := rv.Field(i)
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		// 判断是否为忽略字段
		if name == "-" {
			continue
		}
		// 判断是否为空，若为空则使用字段本身的名称获取value值
		if name == "" {
			name = field.Name
		}
		//获取value值
		v, ok := m[name]
		if !ok {
			continue
		}
		//获取指定字段的类型
		kind := fieldV.Kind()
		// 若字段为指针类型
		if kind == reflect.Ptr {
			// 获取对应字段的kind
			kind = fieldV.Type().Elem().Kind()
		}
		// 设置对应字段的值
		switch kind {
		case reflect.Bool:
			fieldV.SetBool(cast.ToBool(v))
		case reflect.Int:
			fieldV.SetInt(cast.ToInt64(v))
		case reflect.Int8:
			fieldV.SetInt(cast.ToInt64(v))
		case reflect.Int16:
			fieldV.SetInt(cast.ToInt64(v))
		case reflect.Int32:
			fieldV.SetInt(cast.ToInt64(v))
		case reflect.Int64:
			fieldV.SetInt(cast.ToInt64(v))
		case reflect.Uint8:
			fieldV.SetUint(cast.ToUint64(v))
		case reflect.Uint16:
			fieldV.SetUint(cast.ToUint64(v))
		case reflect.Uint32:
			fieldV.SetUint(cast.ToUint64(v))
		case reflect.Uint64:
			fieldV.SetUint(cast.ToUint64(v))
		case reflect.Float32:
			fieldV.SetFloat(cast.ToFloat64(v))
		case reflect.Float64:
			fieldV.SetFloat(cast.ToFloat64(v))
		case reflect.String:
			fieldV.SetString(cast.ToString(v))
		}
	}
}

func ReflectElemType(dest any) reflect.Type {
	rt := reflect.TypeOf(dest)
	if rt.Kind() != reflect.Pointer {
		panic(exception.New("param should be pointer"))
	}
	return rt.Elem()
}
