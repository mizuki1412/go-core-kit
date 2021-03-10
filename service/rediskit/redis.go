package rediskit

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/spf13/cast"
	"time"
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

// val将会以json形式存储，如果不是string的话
func Set(ctx context.Context, key string, val interface{}, expire time.Duration) {
	client = Instance()
	if _, ok := val.(string); !ok {
		val = jsonkit.ToString(val)
	}
	_, err := client.Set(ctx, key, val, expire).Result()
	if err != nil {
		panic(exception.New("redis 存入失败：" + err.Error()))
	}
}

func GetKeyPrefix() string {
	name := configkit.GetStringD(ConfigKeyRedisPrefix)
	if name != "" {
		name += "_"
	}
	return name
}

func LPush(ctx context.Context, key string, val interface{}) {
	client = Instance()
	if _, ok := val.(string); !ok {
		val = jsonkit.ToString(val)
		_, err := client.LPush(ctx, key, val).Result()
		if err != nil {
			panic(exception.New("redis 入队失败：" + err.Error()))
		}
	}
}

func LPop(ctx context.Context, key string, defaultVal string) string {
	client = Instance()
	val, err := client.LPop(ctx, key).Result()
	if err != nil || val == "" {
		return defaultVal
	}
	return val
}

func RPush(ctx context.Context, key string, val interface{}) {
	client = Instance()
	if _, ok := val.(string); !ok {
		val = jsonkit.ToString(val)
		_, err := client.RPush(ctx, key, val).Result()
		if err != nil {
			panic(exception.New("redis 入队失败：" + err.Error()))
		}
	}
}

func RPop(ctx context.Context, key string, defaultVal string) string {
	client = Instance()
	val, err := client.RPop(ctx, key).Result()
	if err != nil || val == "" {
		return defaultVal
	}
	return val
}

func LLen(ctx context.Context, key string) int64 {
	client = Instance()
	val, err := client.LLen(ctx, key).Result()
	if err != nil {
		panic(exception.New("获取redis队列失败"))
	}
	return val
}
