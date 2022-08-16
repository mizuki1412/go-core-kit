package jsonkit

import (
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"strings"
)

var jsonAPI jsoniter.API

func JSON() jsoniter.API {
	if jsonAPI == nil {
		jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary
	}
	return jsonAPI
}

func Test(obj any) string {
	s, _ := json.Marshal(obj)
	return string(s)
}

func ToString(obj any) string {
	s, err := JSON().MarshalToString(obj)
	// todo ?
	if err != nil {
		return "{}"
	}
	return s
}

// ParseObj string, &p, 数组也必须point
func ParseObj(data string, p any) error {
	err := JSON().Unmarshal([]byte(data), p)
	return err
}

func ParseMap(data string) map[string]any {
	//m := map[string]any{}
	//ParseObj(data,&m)
	m, ok := gjson.Parse(data).Value().(map[string]any)
	if !ok {
		return map[string]any{}
	}
	return m
}

// ParseMapUseNumber : use json.Number，如需要高精度计算用decimal转换
// d,_:=decimal.NewFromString(jsonkit.ParseMapUseNumber(str)["key"].(json.Number).String())
// decimal.MarshalJSONWithoutQuotes=true
func ParseMapUseNumber(data string) map[string]any {
	para := make(map[string]any)
	// gjson存在精度问题，jsoniter出现nil错误
	decoder := json.NewDecoder(strings.NewReader(data))
	decoder.UseNumber()
	err := decoder.Decode(&para)
	if err != nil {
		return map[string]any{}
	}
	return para
}
