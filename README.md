toolkit for golang projects

# Usage

## init.go
应用了automaxprocs，校正docker环境中的cpu核数

```go
/// 在项目的main中先导入
package main
import (_ "mizuki/project/core-kit")
```

## configconst.go
关于项目的一些配置key定义。

在service等地方也有此类定义。

## class
通用的一些类的封装和定义

## library
通用的工具库

## service
通用的服务库

## pc
应用于pc端，web+go的模式，go作为基座的一些封装。

## tools
本地使用demo。

# Related projects

## Common
- [cast](https://github.com/spf13/cast)

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

## pc

### bridge
- [go-socket.io](https://github.com/googollee/go-socket.io)