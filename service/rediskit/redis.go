package rediskit

import (
	"context"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/cli/configkey"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"sync"
	"time"
)

var client *redis.Client
var once sync.Once

func Instance() *redis.Client {
	if client == nil {
		once.Do(func() {
			client = redis.NewClient(&redis.Options{
				Addr:     configkit.GetStringD(configkey.RedisHost) + ":" + configkit.GetString(configkey.RedisPort, "6379"),
				Password: configkit.GetStringD(configkey.RedisPwd), // no password set
				DB:       cast.ToInt(configkit.GetString(configkey.RedisDB, "0")),
			})
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

// Set val将会以json形式存储，如果不是string的话
func Set(ctx context.Context, key string, val any, expire time.Duration) {
	client = Instance()
	if _, ok := val.(string); !ok {
		val = jsonkit.ToString(val)
	}
	_, err := client.Set(ctx, key, val, expire).Result()
	if err != nil {
		panic(exception.New("redis save failed: " + err.Error()))
	}
}

func Del(ctx context.Context, keys ...string) {
	client = Instance()
	_, err := client.Del(ctx, keys...).Result()
	if err != nil {
		logkit.Error("redis delete error: " + jsonkit.ToString(keys))
	}
}

func Expire(ctx context.Context, key string, expire time.Duration) {
	client = Instance()
	client.Expire(ctx, key, expire)
}

func GetKeyWithPrefix(key string) string {
	p := configkit.GetStringD(configkey.RedisPrefix)
	if p == "" {
		return key
	} else {
		return p + ":" + key
	}
}

// GetDecoratedKey 分层路径按:分隔，会自动加上项目前缀(ConfigKeyRedisPrefix)
func GetDecoratedKey(subPath string) string {
	name := configkit.GetStringD(configkey.RedisPrefix)
	if name != "" {
		name += ":"
	}
	return name + subPath
}

func LPush(ctx context.Context, key string, val any) {
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

func RPush(ctx context.Context, key string, val any) {
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

func LClear(ctx context.Context, key string) {
	client = Instance()
	_, err := client.LTrim(ctx, key, 1, 0).Result()
	if err != nil {
		panic(exception.New("清除redis队列失败"))
	}
}
