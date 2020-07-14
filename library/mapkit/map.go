package mapkit

import (
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
)

func PutIfAbsent(target map[string]interface{}, key string, val interface{}) {
	_, ok := target[key]
	if !ok {
		target[key] = val
	}
}

func PutAll(target, origin map[string]interface{}) {
	for k, v := range origin {
		target[k] = v
	}
}

func Map2Struct(input map[string]interface{}, obj interface{}) error {
	// todo 自定义的一些class出错
	//return mapstructure.Decode(input, obj)
	return jsonkit.ParseObj(jsonkit.ToString(input), obj)
}

func Merge(dest, src map[string]interface{}) {
	if dest == nil || src == nil {
		return
	}
	for key, val := range src {
		destVal, destOk := dest[key].(map[string]interface{})
		srcVal, srcOk := val.(map[string]interface{})
		if destOk && srcOk {
			Merge(destVal, srcVal)
		} else if val != nil {
			dest[key] = val
		}
	}
}
