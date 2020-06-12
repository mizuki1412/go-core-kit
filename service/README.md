
## config配置项总览
```go
package demo

const ConfigKeyTimeLocation = "time.location"

// 日志目录
const ConfigKeyLogPath = "logger.path"
// 文件名，无后缀
const ConfigKeyLogName = "logger.name"
// 最大保留天数
const ConfigKeyLogMaxRemain = "logger.max-remain"

const ConfigKeyMQTTBroker = "mqtt.broker"
const ConfigKeyMQTTClientID = "mqtt.clientId"
const ConfigKeyMQTTUsername = "mqtt.username"
const ConfigKeyMQTTPwd = "mqtt.pwd"

// rest server相关配置
const ConfigKeyRestServerPort = "rest.port"
// base path
const ConfigKeyRestServerBase = "rest.base"
// session过期时间 默认24h
const ConfigKeySessionExpire = "rest.sessionExpire"

const ConfigKeySwaggerTitle = "swagger.title"
const ConfigKeySwaggerDescription = "swagger.description"
const ConfigKeySwaggerVersion = "swagger.version"
const ConfigKeySwaggerHost = "swagger.host"
const ConfigKeySwaggerBasePath = "swagger.basePath"

const ConfigKeyDBDriver = "db.driver"
const ConfigKeyDBHost = "db.host"
const ConfigKeyDBPort = "db.port"
const ConfigKeyDBUser = "db.user"
const ConfigKeyDBPwd = "db.pwd"
const ConfigKeyDBName = "db.name"

const ConfigKeyRedisHost = "redis.host"
const ConfigKeyRedisPort = "redis.port"
const ConfigKeyRedisPwd = "redis.pwd"
const ConfigKeyRedisDB = "redis.db"
```

## restkit
考虑到未来可能会更换mvc框架，面向项目的mvc框架相关的接口将做一层抽象。

```go
// 启动rest server，加入Action模块
restkit.AddActions(user.All()...)
restkit.Run()

/// 例如action mod中需要定义All()
func All() []func(r *router.Router) {
	return []func(r *router.Router){user.Init, useradmin.Init, roleadmin.Init}
}

/// 其中action的初始定义demo，并配合使用swagger
func Init(router *router.Router) {
	tag := "系统用户模块"
	r := router.Group("/rest/user")
	r.Use(middleware.AuthUsernameAndPwd())
	{
		r.Post("/info", info).Tag(tag).Summary("用户信息")
		r.Post("/logout", logout).Tag(tag).Summary("登出")
		r.Post("/listRoles", listRoles).Tag(tag).Summary("角色列表").Param(listRolesParam{})
		r.Post("/listDepartment", listDepartment).Tag(tag).Summary("部分列表")
		r.Post("/updatePwd", updatePwd).Tag(tag).Summary("密码修改").Param(updatePwdParam{})
		r.Post("/updateUserInfo", updateUserInfo).Tag(tag).Summary("更新用户信息").Param(updateUserInfoParam{})
	}
	router.Group("/rest/user/loginByUsername").Post("", loginByUsername).Tag(tag).Summary("用户名登录").Param(loginByUsernameParam{})
}
```

### 约定/注意

- action的params tags: `validate:"required" description:"xxx" default:"" trim:"true"`
- bean struct tags: `json:"" db:"db-field-name" pk:"true" tablename:"x" autoincrement:"true"`
- context BindForm: 将会先trim，空字符串当做nil。
- 在action中，处理bean中的field时，注意field的valid属性，class中的类可以用Set方法来作为参数设置；自定义的field struct用指针。
- iris.Context.next() 之后的代码逻辑是在response发出之后的，不能再修改response

### context/validator

https://github.com/kataras/iris/wiki/Model-validation

https://github.com/go-playground/validator

### context/session

https://github.com/kataras/iris/wiki/Sessions-database

实际iris redis存储的内容有：
- (prefix)+sessionID
- (prefix)+sessionID-(session的每个key)

redis session key的expire时间，受iris session config控制，同时renew时，旧的也会删除。

### swagger

https://swagger.io/specification/

需要在实际项目中配合使用swagger-ui，访问地址为 ip:port/projectName/swagger 

## sqlkit

- **注意 commit: 如果只有select语句，commit将会出错。**
- bean struct中如果没有db标签，则不会被通用接口insert/update
- sqlx的`missing destination name sth in sth`，是查询出来的字段和类字段不符，在select中限定字段即可。
- update set时：`Set("extend",squirrel.Expr("'{}'::jsonb"))` or `Set("extend","{}")`

## configkit

**注意：请勿在init中获取configkit的参数值，那时还未加载。**

## cronkit

note: https://godoc.org/github.com/robfig/cron

```go
c := cron.New()
c.AddFunc("30 * * * *", func() { fmt.Println("Every hour on the half hour") })
c.AddFunc("30 3-6,20-23 * * *", func() { fmt.Println(".. in the range 3-6am, 8-11pm") })
c.AddFunc("CRON_TZ=Asia/Tokyo 30 04 * * *", func() { fmt.Println("Runs at 04:30 Tokyo time every day") })
c.AddFunc("@hourly",      func() { fmt.Println("Every hour, starting an hour from now") })
c.AddFunc("@every 1h30m10s", func() { fmt.Println("Every hour thirty, starting an hour thirty from now") })
c.Start()
..
// Funcs are invoked in their own goroutine, asynchronously.
...
// Funcs may also be added to a running Cron
c.AddFunc("@daily", func() { fmt.Println("Every day") })
..
// Inspect the cron job entries' next and previous run times.
inspect(c.Entries())
..
c.Stop()  // Stop the scheduler (does not stop any jobs already running).
```
```
// cron format
Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?

// job wrappers
- Recover any panics from jobs (activated by default)
- Delay a job's execution if the previous run hasn't completed yet
- Skip a job's execution if the previous run hasn't completed yet
- Log each job's invocations
```
