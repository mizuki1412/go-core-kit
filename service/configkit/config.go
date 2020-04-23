package configkit

import (
	"github.com/spf13/viper"
	"log"
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
	if !viper.IsSet(key) {
		return defaultVal
	}
	return viper.GetString(key)
}
func GetStringD(key string) string {
	return viper.GetString(key)
}
func GetInt(key string, defaultVal int) int {
	if !viper.IsSet(key) {
		return defaultVal
	}
	return viper.GetInt(key)
}
func GetBool(key string, defaultVal bool) bool {
	if !viper.IsSet(key) {
		return defaultVal
	}
	return viper.GetBool(key)
}
