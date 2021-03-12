package configkit

import (
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func Exist(key string) bool {
	if !viper.IsSet(key) {
		return false
	}
	switch viper.Get(key).(type) {
	case string:
		if cast.ToString(viper.Get(key)) == "" {
			return false
		}
	}
	return true
}
func Get(key string, defaultVal interface{}) interface{} {
	if !Exist(key) {
		return defaultVal
	}
	return viper.Get(key)
}
func GetString(key, defaultVal string) string {
	if !Exist(key) {
		return defaultVal
	}
	return viper.GetString(key)
}
func GetStringD(key string) string {
	return viper.GetString(key)
}
func GetInt(key string, defaultVal int) int {
	if !Exist(key) {
		return defaultVal
	}
	return viper.GetInt(key)
}
func GetBool(key string, defaultVal bool) bool {
	if !Exist(key) {
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
