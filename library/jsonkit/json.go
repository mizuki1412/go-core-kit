package jsonkit

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

func ToString(obj interface{}) string {
	s, _ := json.Marshal(obj)
	return string(s)
}

//Deprecated
func ParseString(obj interface{}) string {
	return ToString(obj)
}

//  string, &p
func ParseObj(data string, p interface{}) {
	// err not handle
	_ = json.Unmarshal([]byte(data), p)
}

func ParseMap(data string) map[string]interface{} {
	m, ok := gjson.Parse(data).Value().(map[string]interface{})
	if !ok {
		return make(map[string]interface{})
	}
	return m
}
