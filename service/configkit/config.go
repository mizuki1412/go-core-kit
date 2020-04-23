package configkit

import (
	"github.com/spf13/viper"
	"log"
	"mizuki/project/core-kit/library/stringkit"
)

/**
viper是大小写不敏感的。
viper在cobra使用时，bind最好用在cmd.Run中，而不是init中
*/

// 注意，load比一般的init慢
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

func Exist(key string) bool {
	return viper.IsSet(key)
}
func Get(key string, defaultVal interface{}) interface{} {
	val := viper.Get(key)
	if val == nil {
		return defaultVal
	}
	return val
}
func GetString(key, defaultVal string) string {
	val := viper.GetString(key)
	if val == "" {
		return defaultVal
	}
	return val
}
func GetInt(key string, defaultVal int) int {
	return viper.GetInt(key)
}
func GetBool(key string, defaultVal bool) bool {
	return viper.GetBool(key)
}
func IsNil(key string) bool {
	// todo only for string
	return stringkit.IsNull(viper.Get(key))
}

//func lookupSub(key string) (interface{}, bool) {
//	arr := strings.Split(key, ".")
//	temp := viper.GetStringMap(arr[0])
//	for i := 1; i < len(arr)-1; i++ {
//		val1, ok1 := temp[arr[i]]
//		if !ok1 {
//			return nil, false
//		} else {
//			temp = val1.(map[string]interface{})
//		}
//	}
//	val1, ok1 := temp[arr[len(arr)-1]]
//	if ok1 {
//		keyPool[key] = val1
//	}
//	return val1, ok1
//}
