package context

// cookie-session版本

/**
func InitSession() gin.HandlerFunc {
	// redis
	redisHost := configkit.GetString(configkey.RedisHost)
	redisPort := configkit.GetString(configkey.RedisPort, "6379")
	redisPwd := configkit.GetString(configkey.RedisPwd)
	redisDB := configkit.GetString(configkey.RedisDB)
	if redisHost != "" {
		logkit.Info("session use redis")
		redisClient := redis.NewClient(&redis.Options{
			Addr:     redisHost + ":" + redisPort,
			Password: redisPwd,
			DB:       cast.ToInt(redisDB),
		})
		store, _ := redistore.NewStore(context.Background(), redisClient)
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
	// 关于跨域（chrome）：需要设置为samesite=none，secure=true，也就是必须在https下才能跨域
	// 如果需要在http下访问，但无跨域要求，需要设置secure=false
	secure := true
	samesite := http.SameSiteNoneMode
	if configkit.Exist(configkey.SessionSecure) && !configkit.GetBool(configkey.SessionSecure) {
		secure = false
		samesite = http.SameSiteLaxMode
	}
	session.Options(sessions.Options{
		Path:     "/",
		MaxAge:   configkit.GetInt(configkey.SessionExpire, 6) * 3600,
		HttpOnly: true,
		Secure:   secure,
		SameSite: samesite,
	})
	// todo save时是否也会存入redis，还是其他情况也会
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
		err := jsonkit.ParseObj(j, user)
		if err != nil {
			logkit.Error(err.Error())
			return nil
		}
		return user
	} else {
		return nil
	}
}
func (ctx *Context) SessionGetUserOrigin() any {
	session := sessions.Default(ctx.Proxy)
	return session.Get("user")
}
func (ctx *Context) SessionGetSchema() string {
	session := sessions.Default(ctx.Proxy)
	r := cast.ToString(session.Get("schema"))
	if r == "" {
		r = sqlkit.SchemaDefault
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

**/
