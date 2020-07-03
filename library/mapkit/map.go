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
