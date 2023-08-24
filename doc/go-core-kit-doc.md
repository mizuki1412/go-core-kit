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

通用的抽象的类定义，可用于本核心库下，弥补原生类型的不足。

## usage

- 定义一些通用类，用于值可能为null的场景。
- 实现数据库的读写转化(Value, Scan)和json的转化接口(MarshalJSON, UnmarshalJSON)。
- 一些常用的类自带函数，如Set等等

## 主要函数

- MarshalJSON, UnmarshalJSON ：用于json序列化
- Scan, Value： 用于sql，Value用于sql传参时驱动调用的。
- isValid：用于sql处理
- Set：值设置

## 类库

### 基本类型

Bool, Decimal, Float64, Int32, Int64, String, Time

### Map

提供了一些map常用的函数接口。

- MapString：对postgres的jsonb格式做了适配。
- MapStringSync：提供了线程安全的MapString

### 数组类型

- ArrInt: 针对postgres.array的int数组封装，提供ToInt32Slice方法。
- ArrString: 针对postgres.array的string分装。
- MapStringArr：针对postgres jsonb的分装，array形式的jsonb

### queue

队列

### file

http上传文件流 

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
	tag := "user:用户模块"
	router.Group("/rest/user/loginByUsername").Use(middleware.CreateSession()).Post("", loginByUsername).Swagger.Tag(tag).Summary("登录-用户名").Param(loginByUsernameParam{})
	router.Group("/rest/user/login").Use(middleware.CreateSession()).Post("", login).Swagger.Tag(tag).Summary("登录").Param(loginParam{})
	router.Group("/rest/user/info").Use(middleware.AuthUsernameAndPwd()).Post("", info).Swagger.Tag(tag).Summary("用户信息")
	r := router.Group("/rest/user", middleware.AuthUsernameAndPwd())
	{
		r.Post("/logout", logout).Swagger.Tag(tag).Summary("登出")
		r.Post("/updatePwd", updatePwd).Swagger.Tag(tag).Summary("密码修改").Param(updatePwdParam{})
		r.Post("/updateUserInfo", updateUserInfo).Swagger.Tag(tag).Summary("更新用户信息").Param(updateUserInfoParam{})
	}
	r1 := router.Group("/rest/user/admin", middleware.AuthUsernameAndPwd())
	{
		r1.Post("/list", listUsers).Swagger.Tag(tag).Summary("用户列表").Param(listUsersParams{})
		r1.Post("/listByRole", listByRole).Swagger.Tag(tag).Summary("用户列表 by role").Param(listByRoleParams{})
		r1.Post("/info", infoAdmin).Swagger.Tag(tag).Summary("用户信息").Param(infoAdminParams{})
	}
}
```

## 约定/注意

- action的params tags: `validate:"required" description:"xxx" default:"" trim:"true"`
- bean struct tags: `json:"" db:"db-field-name" pk:"true" table:"x" auto:"true"`
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
- recover：异常捕捉：打印、返回错误信息


# service-sqlkit

数据库服务

基于squirrel实现 sql 语句的生成，区别于大而全的 ORM 框架，本库借鉴 Mybatis-Plus/Mybatis-Flex 的使用体验，提供了极大的自主性和必要的功能封装。

```go
type Dao[T any] struct {
	meta T
	// 逻辑删除的字段，可替代全局的LogicDelVal
	LogicDelVal []any
	// 返回级联的类型
	ResultType byte
	// 级联实现的函数
	Cascade func(*T)
	// 数据源
	dataSource *DataSource
	// 目标表结构
	modelMeta ModelMeta
}

type Dao struct {
	sqlkit.Dao[$name$]
}

const (
	ResultDefault byte = iota
	ResultNone
)

func New(ds ...*sqlkit.DataSource) Dao {
	dao := Dao{}
	if len(ds) > 0 {
		dao.SetDataSource(ds[0])
	}
	dao.Cascade = func(obj *$name$) {
		switch dao.ResultType {
		case ResultDefault:
		case ResultNone:
		}
	}
	return dao
}
```

## 代码结构说明

- dao: dao 封装类

- dao_mapper: 提供 baseMapper 的增删改查功能

- dao_private: dao 中用到的私有函数

- datasource：和数据源有关，包括事务、schema、dialect

- modelmeta：当前 dao 所对应的 model 的 tablename 和 columnsName

- sqlbuilder：sql 语句构造

- transaction：事务相关的外部函数

## 支持的driver name

- postgres：`"github.com/lib/pq"`
- mysql： `"github.com/go-sql-driver/mysql"`
- mssql: todo
- oracle: todo
- kingbase: todo

## model 标签定义

- pk: bool

- table: 表名

- auto：bool，可用于所有 key

- db：数据库字段名，必须填写才能纳入管理

- logicDel: bool 逻辑删除字段

## 多数据源

通过在 Dao.New 中设置 DataSource 实现数据源的切换。

DataSource 有默认数据源，取自配置文件，也可通过手动创建 NewDataSource()

## 事务处理

通过 sqlkit.TxArea() 实现代码块的事务处理。

## 逻辑删除

全局的指定：sqlkit.LogicDelVal，默认[true,false]

Dao 中可以设置 LogicDelVal 实现局部的逻辑删除，[]any{删除的值，不删除的值}

## 注意

- **注意 commit: 如果事务中第一句是select语句，commit将会出错, 错误提示 parse C 等。**
- 在rows遍历时，注意close，特别是有级联查询存在时，如果不close将会占用连接。
- sqlx的`missing destination name sth in sth`，是查询出来的字段和类字段不符，在select中限定字段即可。
- update set时：`Set("extend",squirrel.Expr("'{}'::jsonb"))` or `Set("extend","{}")`
- class.mapString在插入数据库时将用jsonb格式，并且不是完全替换，而是merge的方式(```coalesce(extend, '{}'::jsonb) || '$param'::jsonb```)。如果要删除其中的key，需要设置key为null。 merge时只会merge顶层的keys。

## demo

todo

```go
func (dao Dao) UpdateConfirm(id int64){
    sql, args := dao.Builder().Update().Set("extend",squirrel.Expr("jsonb_set(extend, '{confirm}','true',true)")).Where("id=?",id).Sql()
    dao.Exec(sql, args...)
}

func (dao Dao) List(dTypes []string) []model.AlarmMsg {
  builder := dao.Builder().Select(dao.SelectColumns("msg")).FromAs("msg").Join(infodao.New(dao.DataSource()).Table("info")).Where("msg.deviceType=info.id").OrderBy("msg.deviceType").OrderBy("msg.id")
	if dTypes!=nil && len(dTypes) > 0 {
		flag, arg := pghelper.GenUnnestString(dTypes)
		builder = builder.Where("msg.deviceType in "+flag, arg)
	}
	sql, args := builder.MustSql()
	return dao.scan(sql, args)
}

func (dao Dao) ListId(dType []string) []string {
	builder := dao.Builder().Select("id").Where("off=?", false).OrderBy("id")
	if dType!=nil && len(dType)>0{
		builder = pghelper.WhereUnnestInt(builder,"id in ", dType)
	}
	sql, args := builder.MustSql()
	rows := dao.Query(sql, args...)
	list := make([]string, 0, 5)
	defer rows.Close()
	for rows.Next() {
		ret, err := rows.SliceScan()
		if err != nil {
			panic(exception.New(err.Error()))
		}
		list = append(list, ret[0].(string))
	}
	return list
}
```

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
