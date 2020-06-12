package jsonkit

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/tidwall/gjson"
)

var json jsoniter.API

func JSON() jsoniter.API {
	if json == nil {
		json = jsoniter.ConfigCompatibleWithStandardLibrary
	}
	return json
}

func ToString(obj interface{}) string {
	s, err := JSON().MarshalToString(obj)
	if err != nil {
		panic(exception.New("json parse error"))
	}
	return s
}

//  string, &p
func ParseObj(data string, p interface{}) {
	err := JSON().Unmarshal([]byte(data), p)
	if err != nil {
		panic(exception.New("json parse error"))
	}
}

func ParseMap(data string) map[string]interface{} {
	//m := map[string]interface{}{}
	//ParseObj(data,&m)
	m, ok := gjson.Parse(data).Value().(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
	}
	return m
}
