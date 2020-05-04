package mapkit

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
