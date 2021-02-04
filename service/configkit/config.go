package configkit

import (
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/spf13/viper"
)

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
	if !viper.IsSet(key) || viper.GetString(key) == "" {
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

func GetBoolD(key string) bool {
	return viper.GetBool(key)
}

func GetStringListD(key string) []string {
	str := GetStringD(key)
	if str == "" {
		return nil
	}
	var arr []string
	_ = jsonkit.ParseObj(key, &arr)
	return arr
}
