package configkit

import (
	"github.com/spf13/viper"
)

func Exist(key string) bool {
	return viper.IsSet(key)
}

func GetString(key string, defaultVal ...string) string {
	if !Exist(key) && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return viper.GetString(key)
}

func GetInt(key string, defaultVal ...int) int {
	if !Exist(key) && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return viper.GetInt(key)
}
func GetBool(key string, defaultVal ...bool) bool {
	if !Exist(key) && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return viper.GetBool(key)
}

func GetStringList(key string) []string {
	return viper.GetStringSlice(key)
}

func GetStringMap(key string) map[string]any {
	return viper.GetStringMap(key)
}

func Set(key string, val any) {
	viper.Set(key, val)
}
