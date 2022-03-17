package context

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/init/configkey"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/sessions"
	"github.com/mizuki1412/sessions/cookie"
	redistore "github.com/mizuki1412/sessions/redis"
	"github.com/spf13/cast"
	"net/http"
)

func InitSession() gin.HandlerFunc {
	// redis
	redisHost := configkit.GetStringD(configkey.RedisHost)
	redisPort := configkit.GetString(configkey.RedisPort, "6379")
	redisPwd := configkit.GetStringD(configkey.RedisPwd)
	redisDB := configkit.GetStringD(configkey.RedisDB)
	//redisPrefix := configkit.GetString(configkey.RedisPrefix, "")
	if redisHost != "" {
		logkit.Info("session use redis")
		redisClient := redis.NewClient(&redis.Options{
			Addr:     redisHost + ":" + redisPort,
			Password: redisPwd,
			DB:       cast.ToInt(redisDB),
		})
		store, _ := redistore.NewStore(context.Background(), redisClient)
		//store, _ := redis.NewStoreWithDB(10, "tcp", redisHost+":"+redisPort, redisPwd, redisDB, []byte("default"))
		return sessions.Sessions("sessionID", store)
	} else {
		store := cookie.NewStore([]byte("default"))
		logkit.Info("session use cache")
		return sessions.Sessions("sessionID", store)
	}
}

// SessionSetUser session每次请求时都会从redis中获取，所以在session中存储的务必是string，如果是对象，会被自动转为json，但如果其中有unicode，可能造成指数增加
func (ctx *Context) SessionSetUser(user any) {
	session := sessions.Default(ctx.Proxy)
	if user == nil {
		return
	}
	if _, ok := user.(string); !ok {
		session.Set("user", jsonkit.ToString(user))
	} else {
		session.Set("user", user)
	}

}
func (ctx *Context) SessionSetSchema(schema string) {
	session := sessions.Default(ctx.Proxy)
	session.Set("schema", schema)
}
func (ctx *Context) SessionSetToken(token string) {
	session := sessions.Default(ctx.Proxy)
	session.Set("token", token)
}
func (ctx *Context) SessionSet(key string, val any) {
	session := sessions.Default(ctx.Proxy)
	session.Set(key, val)
}

// SessionSave 刷新session 到cookie
func (ctx *Context) SessionSave() {
	session := sessions.Default(ctx.Proxy)
	session.Options(sessions.Options{
		Path:     "/",
		MaxAge:   configkit.GetInt(configkey.SessionExpire, 6) * 3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
	err := session.Save()
	if err != nil {
		logkit.Error(exception.New(err.Error()))
	}
}

// SessionGetUser return eg *model.User
func (ctx *Context) SessionGetUser() *model.User {
	session := sessions.Default(ctx.Proxy)
	json := session.Get("user")
	if j, ok := json.(string); ok {
		user := &model.User{}
		// todo 每次都要转换，可能存在性能问题
		err := jsonkit.ParseObj(j, user)
		if err != nil {
			return nil
		}
		return user
	} else {
		return nil
	}
}
func (ctx *Context) SessionGetSchema() string {
	session := sessions.Default(ctx.Proxy)
	r := cast.ToString(session.Get("schema"))
	if r == "" {
		r = "public"
	}
	return r
}
func (ctx *Context) SessionGetToken() string {
	session := sessions.Default(ctx.Proxy)
	return cast.ToString(session.Get("token"))
}
func (ctx *Context) SessionClear() {
	session := sessions.Default(ctx.Proxy)
	session.Clear()
	ctx.SessionSave()
}

func (ctx *Context) SessionID() string {
	session := sessions.Default(ctx.Proxy)
	return session.ID()
}
