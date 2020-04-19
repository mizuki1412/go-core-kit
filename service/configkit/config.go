package configkit

import (
	"github.com/spf13/viper"
	"log"
	"mizuki/project/core-kit/library/stringkit"
	"strings"
)

// config-key(eg: abc.exe.ss) : val
var keyPool map[string]interface{}

func init() {
	keyPool = make(map[string]interface{})
}

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	// 这里可以执行多次的 搜索多个地址
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("config load error:", err)
	}
}

func Get(key string, defaultVal interface{}) interface{} {
	if !strings.Contains(key, ".") {
		return viper.Get(key)
	}
	val, ok := keyPool[key]
	if ok {
		return val
	}
	val1, ok1 := lookupSub(key)
	if !ok1 {
		return defaultVal
	} else {
		return val1
	}
}
func GetString(key, defaultVal string) string {
	if !strings.Contains(key, ".") {
		return viper.GetString(key)
	}
	val, ok := keyPool[key]
	if ok {
		return val.(string)
	}
	val1, ok1 := lookupSub(key)
	if !ok1 {
		return defaultVal
	} else {
		return val1.(string)
	}
}
func GetInt(key string, defaultVal int) int {
	if !strings.Contains(key, ".") {
		return viper.GetInt(key)
	}
	val, ok := keyPool[key]
	if ok {
		return val.(int)
	}
	val1, ok1 := lookupSub(key)
	if !ok1 {
		return defaultVal
	} else {
		return val1.(int)
	}
}
func GetBool(key string, defaultVal bool) bool {
	if !strings.Contains(key, ".") {
		return viper.GetBool(key)
	}
	val, ok := keyPool[key]
	if ok {
		return val.(bool)
	}
	val1, ok1 := lookupSub(key)
	if !ok1 {
		return defaultVal
	} else {
		return val1.(bool)
	}
}
func IsNil(key string) bool {
	// todo only for string
	return stringkit.IsNull(viper.Get(key))
}

func lookupSub(key string) (interface{}, bool) {
	arr := strings.Split(key, ".")
	temp := viper.GetStringMap(arr[0])
	for i := 1; i < len(arr)-1; i++ {
		val1, ok1 := temp[arr[i]]
		if !ok1 {
			return nil, false
		} else {
			temp = val1.(map[string]interface{})
		}
	}
	val1, ok1 := temp[arr[len(arr)-1]]
	return val1, ok1
}
