# service

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
- context BindForm: 支持在params中直接指定基本类型和class包中的类型。
- 在action中，处理bean中的field时，注意field的valid属性，class中的类可以用Set方法来作为参数设置；自定义的field struct用指针。
- iris.Context.next() 之后的代码逻辑是在response发出之后的，不能再修改response
- router.use在使用时，多拦截器放一个use。

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

标准：https://swagger.io/specification/v2/

swagger-ui可以单独部署，后端只提供doc.json

需要在实际项目中配合使用swagger-ui，访问地址为 ip:port/projectName/swagger 

### swagger-ui

用github.com/markbates/pkger打包静态资源。swagger-ui需要放置在子项目根目录下。

需要在子项目中的main.go顶行写入 ```//go:generate pkger -include /swagger-ui```，然后运行 ```go generate``` 生成pkged.go。

在router/router.go中将自动配置```<base>/swagger```为访问地址。

todo： 暂不能从core/swagger-ui导入资源。
