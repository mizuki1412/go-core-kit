[TOC]

# cli

框架的CMD入口配置，包含一些常量和可配置项。

基于 [cobra](https://cobra.dev/) 和 viper。

对 cobra 简单封装，能够在 main.go 中直接设置 rootCMD 和 childCMD 。

demo:

```go
cli.RootCMD(&cobra.Command{
  Use: "main",
  Run: func(cmd *cobra.Command, args []string) {
    restkit.AddActions(user.All()...)
    restkit.AddActions(download.Init)
    _ = restkit.Run()
  },
})
cli.AddChildCMD(&cobra.Command{
  Use: "test",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("123")
  },
})
cli.AddChildCMD(cmd.TCPServerCMD())
cli.Execute()

// 额外自定义参数
cmd.Flags().String("port", "", "端口")
// 设置必填
_ = cmd.MarkFlagRequired("port")
```

## 配置文件

默认的配置文件为当前目录的`./config.yaml`。

有个全局参数 `-c` 或 `--config` 可以指定配置文件的具体路径。

命令行的参数将覆盖配置文件中相同的参数。

## cmd 例子

`/cmd`下包括了一些场景下使用的工具：

- FrontDaoCMDNext：将 swagger 接口导出成前端 dao 文件（js）。
- File2LineCli: 配置文件转命令行文字
- MarkdownDocCMD：markdown 文件导出
- MQTTTestCMD：mqtt demo
- PGSqlToStructCMD：通过 sql 生成 model
- TCPServerCMD：tcp server
- WebStaticServerCMD：静态文件服务器

# class 封装类

# library

工具库

## jsonkit

封装 sonic：https://github.com/bytedance/sonic

## httpkit

http client

## cmdkit

调用系统 cmd。

参考：https://colobu.com/2020/12/27/go-with-os-exec/

## concurrentkit

异步等待。

## timekit

时间处理

## framekit

应用于数据流帧的拆包粘包处理。

## stringkit

字符串相关处理

## tarkit

压缩包处理

## templatekit

模板

## ftpkit

ftp 相关的封装

## inikit

note: https://ini.unknwon.cn/docs/intro/getting_started

```go
cfg, err := ini.Load(
    []byte("raw data"), // 原始数据
    "filename",         // 文件路径
    io.NopCloser(bytes.NewReader([]byte("some other data"))),
)

// 典型读取操作，默认分区可以使用空字符串表示
fmt.Println("App Mode:", cfg.Section("").Key("app_mode").String())
fmt.Println("Data Path:", cfg.Section("paths").Key("data").String())

// 试一试自动类型转换
fmt.Printf("Port Number: (%[1]T) %[1]d\n", cfg.Section("server").Key("http_port").MustInt(9999))
fmt.Printf("Enforce Domain: (%[1]T) %[1]v\n", cfg.Section("server").Key("enforce_domain").MustBool(false))

// 差不多了，修改某个值然后进行保存
cfg.Section("").Key("app_mode").SetValue("production")
cfg.SaveTo("my.ini.local")
```

## ipkit

ip 的处理

# service

## configkit

封装viper，获取配置参数

**注意：请勿在init中获取configkit的参数值，那时还未加载。**

## logkit

日志，包括rolling package。

## cachekit

缓存服务。包含内存和redis。

## rediskit



## cronkit

定时任务

note: https://godoc.org/github.com/robfig/cron

```go
c := cron.New()
c.AddFunc("0 30 * * * *", func() { fmt.Println("Every hour on the half hour") })
c.AddFunc("@hourly",      func() { fmt.Println("Every hour") })
c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty") })
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

cron库语法说明：
```text
cron format: 
Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?

```

## excelkit

excel表格处理

## influxkit

influx1

## mqttkit

mqtt服务

## netkit

tcp/udp server and client。

https://gnet.host/docs/quickstart/

## storagekit

本地文件存储服务

## serialkit

串口相关

## pdfkit

html转pdf

# service-restkit

web后端服务。

目前基于gin，考虑到以后可能会更换mvc，所以抽象了一层。

## 使用

```go
// 启动rest server，加入Action模块
restkit.AddActions(user.Init)
restkit.AddActions(useradmin.Init)
restkit.Run()

/// 其中action的初始定义demo，并配合使用swagger
func Init(router *router.Router) {
	tag := "系统用户模块"
	r := router.Group("/rest/user")
	r.Use(middleware.AuthUsernameAndPwd())
	{
		r.Post("/info", info).Tag(tag).Summary("用户信息")
		r.Post("/logout", logout).Tag(tag).Summary("登出")
		r.Post("/listRoles", listRoles).Tag(tag).Summary("角色列表").Param(listRolesParam{})
	}
	router.Group("/rest/user/loginByUsername").Post("", loginByUsername).Tag(tag).Summary("用户名登录").Param(loginByUsernameParam{})
}
```

## 约定/注意

- action的params tags: `validate:"required" description:"xxx" default:"" trim:"true"`
- bean struct tags: `json:"" db:"db-field-name" pk:"true" tablename:"x" autoincrement:"true"`
- context BindForm: 将会先trim，空字符串当做nil。
- context BindForm: 支持在params中直接指定基本类型和class包中的类型。
- 在action中，处理bean中的field时，注意field的valid属性，class中的类可以用Set方法来作为参数设置；自定义的field struct用指针。
- iris.Context.next() 之后的代码逻辑是在response发出之后的，不能再修改response
- router.use在使用时，多拦截器放一个use。

## 抽象层设计

### context

```go
type Context struct {
	Proxy    *gin.Context
	Request  *http.Request
	Response gin.ResponseWriter
}
```

- Get、Set：会话的临时变量
- DBTx：当前会话的数据库事务
- BindForm：对 request 的 form/query/json 的参数解析/校验/打印等，需要用 struct 定义参数
- Json/JsonSuccess/RawSuccess/Html/File/...：response 输出
- SessionToken/SessionSetUser/...：session 的支持

### context.session

- 不再支持 cookie，因为有禁用的场景；改用 header 中的token 字段。
- 在登录拦截的时候会获取 token，然后 context.Set/Get，判断后续逻辑。其他拦截如果需要可以参考`authup.go`。
- session的存储通过 cachekit，支持 redis 和内存。

### router

实现正常配置路由信息的同时，配置 swagger 的信息，在代码过程中配置，避免另外生成或写易错的标签。

```go
type Router struct {
	Proxy      *gin.Engine
	Base       string
	ProxyGroup *gin.RouterGroup
	Swagger    *swg.SwaggerPath
}

router.Group("/rest/user/loginByUsername").Use(middleware.CreateSession()).Post("", loginByUsername).Swagger.Tag(tag).Summary("登录-用户名").Param(loginByUsernameParam{})
	router.Group("/rest/user/login").Use(middleware.CreateSession()).Post("", login).Swagger.Tag(tag).Summary("登录").Param(loginParam{})
	router.Group("/rest/user/info").Use(middleware.AuthUsernameAndPwd()).Post("", info).Swagger.Tag(tag).Summary("用户信息")
	r := router.Group("/rest/user", middleware.AuthUsernameAndPwd())
	{
		r.Post("/logout", logout).Swagger.Tag(tag).Summary("登出")
		r.Post("/updatePwd", updatePwd).Swagger.Tag(tag).Summary("密码修改").Param(updatePwdParam{})
		r.Post("/updateUserInfo", updateUserInfo).Swagger.Tag(tag).Summary("更新用户信息").Param(updateUserInfoParam{})
	}
```

- Group: 路径组，附带 baseUrl
- Use：中间处理组件使用
- Post/Get/Any：附带 baseUrl
- GetOrigin: 不附带 baseUrl
- 和 swagger 绑定配置，通过 Router.Swagger 来配置 swagger 相关参数。
- RegisterSwagger：将内置的 swagger-ui 注册到路由

### swagger

标准：https://swagger.io/specification/v2/

swagger-ui可以单独部署，后端只提供doc.json

需要在实际项目中配合使用swagger-ui，访问地址为 `ip:port/projectName/swagger` 

更新swagger-ui：

```js
// 从github上下载更新的源码
// 取出源码中dist/下除.map外的文件放入本目录的swagger-ui中。

// 修改 index.html
<script> <style> href/src 加前缀./swagger

// 修改 swagger-initializer.js
url: "./swagger-doc",
```



## middleware

```go
router.Use(middleware.Log())
router.Use(middleware.Cors())
router.Use(middleware.Recover())
```

- log: 请求前请求后打印
- cors：跨域
- recover：异常捕捉：打印、回滚/提交、返回错误信息


# service-sqlkit

数据库服务



支持的driver name



# 框架内可配置函数或变量

## restkit

### context.TransferRestRet

转换自定义的 response 输出格式。

```go
// 自定义输出格式：{status:0, errmsg:xxx, message:xxx}
context.TransferRestRet = func(ret context.RestRet) any {
  r := Ret{}
  if ret.Result == context.ResultSuccess {
    r.Status = 0
  } else if ret.Result == context.ResultAuthErr {
    r.Status = 2
    r.Errmsg.Set(ret.Message)
  } else {
    r.Status = 1
    r.Errmsg.Set(ret.Message)
  }
  r.Message = ret.Data
  return r
}
```

### context.HeaderTokenKey

request header 中的 token 字段名，用于 session

# service-third

第三方服务集成

# mod

公用的业务模块

# pc

应用于pc端，web+go的模式，go作为基座的一些封装。

# iot

针对 IoT 相关的处理库。

暂不更新。
