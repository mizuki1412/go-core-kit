package jsonkit

import (
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/bytedance/sonic/decoder"
	"log"
)

func Marshal(obj any) ([]byte, error) {
	r, err := sonic.Marshal(obj)
	if err != nil {
		return []byte("null"), err
	}
	return r, nil
}

func Unmarshal(data []byte, p any) error {
	err := sonic.Unmarshal(data, p)
	return err
}

func ToString(obj any) string {
	s, err := sonic.MarshalString(obj)
	if err != nil {
		log.Println(err)
		return "{}"
	}
	return s
}

// ParseObj string, &p, 数组也必须point
func ParseObj(data string, p any) error {
	return Unmarshal([]byte(data), p)
}

func ParseMap(data string) map[string]any {
	m := map[string]any{}
	err := ParseObj(data, &m)
	if err != nil {
		return nil
	}
	return m
}

// ParseMapUseNumber : use json.Number，如需要高精度计算用decimal转换
// d,_:=decimal.NewFromString(jsonkit.ParseMapUseNumber(str)["key"].(json.Number).String())
// decimal.MarshalJSONWithoutQuotes=true
func ParseMapUseNumber(data string) map[string]any {
	para := make(map[string]any)
	// gjson存在精度问题，jsoniter出现nil错误
	dc := decoder.NewDecoder(data)
	dc.UseNumber()
	err := dc.Decode(&para)
	if err != nil {
		return map[string]any{}
	}
	return para
}

func Get(src string, path ...any) ast.Node {
	node, _ := sonic.GetFromString(src, path...)
	return node
}
