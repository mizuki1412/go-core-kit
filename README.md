
# go-core-kit

toolkit for golang projects

`go get -d github.com/mizuki1412/go-core-kit@v1.0.0`

## 可替换的函数目录
- session.SessionGetUserFunc 获取session中user对象的转换处理函数

# init
本库使用的入口，以及配置参数信息相关的绑定函数

# class
通用的一些类的封装和定义

# library
通用的工具库

# service
通用的服务库

# service-third
针对第三方服务接口的封装

# pc
应用于pc端，web+go的模式，go作为基座的一些封装。

# mod
公用的业务模块

# tool-local
本地使用的一些小工具

# 其他

## 可替换的函数列表

- service.restkit.context.session.SessionGetUserFunc: 获取session中user对象的转换处理函数

# related projects

## Common
- [cast](https://github.com/spf13/cast)
- [decimal](https://github.com/shopspring/decimal)

## library

### jsonkit
- [jsoniter](https://github.com/json-iterator/go)
- [gjson](https://github.com/tidwall/gjson)

### inikit
- [ini](https://github.com/go-ini/ini)

## service

### cronkit
- [cron](https://github.com/robfig/cron)

### logkit
- [zap](https://github.com/uber-go/zap)
- [rollingwriter](https://github.com/arthurkiller/rollingwriter)

### configkit
- [viper](https://github.com/spf13/viper)

### restkit
- [iris](https://github.com/kataras/iris)

### sqlkit
- [sqlx](https://github.com/jmoiron/sqlx)
- [squirrel](https://github.com/Masterminds/squirrel)
- [pgsql driver](https://github.com/lib/pq)

### mqttkit
- [mqtt](https://github.com/eclipse/paho.mqtt.golang)

### rediskit
- [go-redis](https://github.com/go-redis/redis/v8)

### serialkit
- [serial](https://go.bug.st/serial)

## pc

### bridge
- [go-socket.io](https://github.com/googollee/go-socket.io)