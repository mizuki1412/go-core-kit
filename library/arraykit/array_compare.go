package arraykit

// 比较两个数组，得出新增的和删除的，src基于base
func CompareAddAndDel(src, base []map[string]interface{}, key string) ([]map[string]interface{}, []map[string]interface{}) {
	var news []map[string]interface{}
	var dels []map[string]interface{}
	cache := map[interface{}]map[string]interface{}{}
	cacheFlag := map[interface{}]bool{}
	for _, v := range base {
		cacheFlag[v[key]] = false
		cache[v[key]] = v
	}
	for _, v := range src {
		if _, ok := cache[v[key]]; ok {
			cacheFlag[v[key]] = true
		} else {
			news = append(news, v)
		}
	}
	for k, v := range cacheFlag {
		if !v {
			dels = append(dels, cache[k])
		}
	}
	return news, dels
}

func CompareAddAndDel4StrArr(src, base []string) ([]string, []string) {
	var news []string
	var dels []string
	cache := map[string]string{}
	cacheFlag := map[string]bool{}
	for _, v := range base {
		cacheFlag[v] = false
		cache[v] = v
	}
	for _, v := range src {
		if _, ok := cache[v]; ok {
			cacheFlag[v] = true
		} else {
			news = append(news, v)
		}
	}
	for k, v := range cacheFlag {
		if !v {
			dels = append(dels, k)
		}
	}
	return news, dels
}
