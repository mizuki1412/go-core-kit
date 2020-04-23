package jsonkit

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"log"
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
	// todo
	if err != nil {
		log.Println(err)
	}
	return s
}

//  string, &p
func ParseObj(data string, p interface{}) {
	// todo err not handle
	_ = JSON().Unmarshal([]byte(data), p)
}

func ParseMap(data string) map[string]interface{} {
	//m := map[string]interface{}{}
	//ParseObj(data,&m)
	m, ok := gjson.Parse(data).Value().(map[string]interface{})
	if !ok {
		return make(map[string]interface{})
	}
	return m
}
