
## config配置项总览
```go
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
```

## restkit
考虑到未来可能会更换mvc框架，面向项目的mvc框架相关的接口将做一层抽象。

```go
// 启动rest server，加入Action模块
restkit.AddActions(
	user.All()...)
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

action的params tags: `validate:"required" description:"xxx" default:""`

bean struct tags: `json:"" db:"db-field-name" pk:"true" tablename:"x" autoincrement:"true"`

context BindForm: 将会先trim，空字符串当做nil。

在action中，处理bean中的field时，注意field的valid属性，class中的类可以用Set方法来作为参数设置；自定义的field struct用指针。

## sqlkit
- **注意 commit: 如果只有select语句，commit将会出错。**