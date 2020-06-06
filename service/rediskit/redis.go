package rediskit

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"mizuki/framework/core-kit/service/configkit"
)

var client *redis.Client

func Instance() *redis.Client {
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     configkit.GetStringD(ConfigKeyRedisHost) + ":" + configkit.GetString(ConfigKeyRedisPort, "6379"),
			Password: configkit.GetStringD(ConfigKeyRedisPwd), // no password set
			DB:       cast.ToInt(configkit.GetString(ConfigKeyRedisDB, "0")),
		})
	}
	return client
}

func Get(ctx context.Context, key string, defaultVal string) string {
	client = Instance()
	val, err := client.Get(ctx, key).Result()
	if err != nil || val == "" {
		return defaultVal
	} else {
		return val
	}
}
