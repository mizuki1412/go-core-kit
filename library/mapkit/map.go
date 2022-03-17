package mapkit

import (
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
)

func PutIfAbsent(target map[string]any, key string, val any) {
	_, ok := target[key]
	if !ok {
		target[key] = val
	}
}

func PutAll(target, origin map[string]any) {
	for k, v := range origin {
		target[k] = v
	}
}

// obj need pointer
func Map2Struct(input map[string]any, obj any) error {
	// todo 自定义的一些class出错
	//return mapstructure.Decode(input, obj)
	return jsonkit.ParseObj(jsonkit.ToString(input), obj)
}

func Merge(dest, src map[string]any) {
	if dest == nil || src == nil {
		return
	}
	for key, val := range src {
		destVal, destOk := dest[key].(map[string]any)
		srcVal, srcOk := val.(map[string]any)
		if destOk && srcOk {
			Merge(destVal, srcVal)
		} else if val != nil {
			dest[key] = val
		}
	}
}
