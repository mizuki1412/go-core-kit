package context

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/go-core-kit/service/rediskit"
	"github.com/spf13/cast"
	"time"
)

var sessionManager *sessions.Sessions

func InitSession() {
	sessionManager = sessions.New(sessions.Config{
		Cookie:                      "session",
		AllowReclaim:                true,
		Expires:                     time.Duration(configkit.GetInt(ConfigKeySessionExpire, 12)) * time.Hour,
		DisableSubdomainPersistence: true, // 开发环境跨域时对火狐有效
	})
	// redis
	redisHost := configkit.GetStringD(rediskit.ConfigKeyRedisHost)
	redisPort := configkit.GetString(rediskit.ConfigKeyRedisPort, "6379")
	redisPwd := configkit.GetStringD(rediskit.ConfigKeyRedisPwd)
	redisDB := configkit.GetStringD(rediskit.ConfigKeyRedisDB)
	redisPrefix := configkit.GetString(rediskit.ConfigKeyRedisPrefix, "")
	if redisHost != "" {
		db := redis.New(redis.Config{
			Network:   "tcp",
			Addr:      redisHost + ":" + redisPort,
			Timeout:   time.Duration(30) * time.Second,
			MaxActive: 10,
			Password:  redisPwd,
			Database:  redisDB,
			Prefix:    redisPrefix + "-session-",
			Delim:     "-",
			Driver:    redis.Redigo(), // redis.Radix() can be used instead.
		})
		// Optionally configure the underline driver:
		// driver := redis.Redigo()
		// driver.MaxIdle = ...
		// driver.IdleTimeout = ...
		// driver.Wait = ...
		// redis.Config {Driver: driver}
		// Close connection when control+C/cmd+C
		iris.RegisterOnInterrupt(func() {
			_ = db.Close()
			logkit.Error("session in redis: error")
		})
		sessionManager.UseDatabase(db)
		logkit.Info("session use redis")
	}
}

func (ctx *Context) Session() *sessions.Session {
	return sessionManager.Start(ctx.Proxy)
}

// session每次请求时都会从redis中获取，所以在session中存储的务必是string，如果是对象，会被自动转为json，但如果其中有unicode，可能造成斜杆指数增加
func (ctx *Context) SessionSetUser(user interface{}) {
	if user == nil {
		return
	}
	if _, ok := user.(string); !ok {
		ctx.Session().Set("user", jsonkit.ToString(user))
	} else {
		ctx.Session().Set("user", user)
	}
}
func (ctx *Context) SessionSetSchema(schema string) {
	ctx.Session().Set("schema", schema)
}
func (ctx *Context) SessionSetToken(token string) {
	ctx.Session().Set("token", token)
}

var SessionGetUserFunc = func(ctx *Context) interface{} {
	// 默认处理，在自定义请覆盖
	json := ctx.Session().GetString("user")
	if json != "" {
		return json
	} else {
		return nil
	}
}

// eg *model.User
func (ctx *Context) SessionGetUser() interface{} {
	return SessionGetUserFunc(ctx)
}
func (ctx *Context) SessionGetSchema() string {
	return ctx.Session().GetStringDefault("schema", "public")
}
func (ctx *Context) SessionGetToken() string {
	return ctx.Session().GetString("token")
}
func (ctx *Context) SessionRemoveUser() {
	ctx.Session().Delete("user")
}

func (ctx *Context) RenewSession() *sessions.Session {
	sess := ctx.Session()
	if !sess.IsNew() {
		sessionManager.Destroy(ctx.Proxy)
		return ctx.Session()
	}
	return sess
}
func (ctx *Context) UpdateSessionExpire() {
	_ = sessionManager.UpdateExpiration(ctx.Proxy, time.Duration(cast.ToInt(configkit.GetString(ConfigKeySessionExpire, "12")))*time.Hour)
}
