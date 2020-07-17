
# go-core-kit

toolkit for golang projects

# usage

## init.go
应用了automaxprocs，校正docker环境中的cpu核数

```go
/// 在项目的main中先导入
package main
import (_ "github.com/mizuki1412/go-core-kit/init")
```

## configconst.go
关于项目的一些配置key定义。

在service等地方也有此类定义。

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

# tool-local
本地使用的一些小工具

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