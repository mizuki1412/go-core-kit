toolkit for golang projects

# Define

action的params tags: `validate:"required" description:"xxx" default:""`
bean struct tags: `json:"" db:"db-field-name" pk:"true" tablename:"x"`

context BindForm: 将会先trim，空字符串当做nil
处理bean中field时，注意valid，class中的类可以用Set方法；自定义field struct用指针。

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