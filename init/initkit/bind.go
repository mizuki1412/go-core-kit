package initkit

import (
	corekit "github.com/mizuki1412/go-core-kit/init"
	"github.com/mizuki1412/go-core-kit/service-third/amapkit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/influxkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/go-core-kit/service/rediskit"
	"github.com/mizuki1412/go-core-kit/service/restkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/swagger"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

// 注意，load比一般的init慢
func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	// 这里可以执行多次的 搜索多个地址
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("no config load")
	}
}

func DefFlags(cmd *cobra.Command) {
	cmd.Flags().String(corekit.ConfigKeyProjectDir, ".", "项目目录")
	cmd.Flags().String(configkit.ConfigKeyTimeLocation, "Asia/Shanghai", "项目中用到的时区")
	cmd.Flags().String(rediskit.ConfigKeyRedisPrefix, "", "redis key的前缀")
	cmd.Flags().String(rediskit.ConfigKeyRedisHost, "", "redis host")
	cmd.Flags().String(rediskit.ConfigKeyRedisPort, "", "")
	cmd.Flags().String(rediskit.ConfigKeyRedisDB, "", "redis db 数据库号")
	cmd.Flags().String(rediskit.ConfigKeyRedisPwd, "", "")
	cmd.Flags().String(influxkit.ConfigKeyInfluxURL, "", "")
	cmd.Flags().String(influxkit.ConfigKeyInfluxUser, "", "")
	cmd.Flags().String(influxkit.ConfigKeyInfluxPwd, "", "")
	cmd.Flags().String(influxkit.ConfigKeyInfluxDBName, "", "默认的数据库")
	cmd.Flags().String(logkit.ConfigKeyLogPath, "", "日志目录；默认在project.dir下")
	cmd.Flags().String(logkit.ConfigKeyLogName, "main", "日志文件名，无后缀")
	cmd.Flags().String(logkit.ConfigKeyLogMaxRemain, "", "最大保留天数")
	cmd.Flags().String(restkit.ConfigKeyRestServerBase, "", "rest base url")
	cmd.Flags().String(restkit.ConfigKeyRestServerPort, "", "")
	cmd.Flags().String(restkit.ConfigKeyRestRequestBodySize, "", "限制request最大，单位MB")
	cmd.Flags().String(sqlkit.ConfigKeyDBDriver, "", "")
	cmd.Flags().String(sqlkit.ConfigKeyDBHost, "", "")
	cmd.Flags().String(sqlkit.ConfigKeyDBPort, "", "")
	cmd.Flags().String(sqlkit.ConfigKeyDBName, "", "")
	cmd.Flags().String(sqlkit.ConfigKeyDBUser, "", "")
	cmd.Flags().String(sqlkit.ConfigKeyDBPwd, "", "")
	cmd.Flags().String(context.ConfigKeySessionExpire, "", "session expire 单位小时")
	cmd.Flags().String(swagger.ConfigKeySwaggerBasePath, "", "/path")
	cmd.Flags().String(swagger.ConfigKeySwaggerHost, "", "")
	cmd.Flags().String(swagger.ConfigKeySwaggerDescription, "", "")
	cmd.Flags().String(swagger.ConfigKeySwaggerTitle, "", "")
	cmd.Flags().String(swagger.ConfigKeySwaggerVersion, "1.0.0", "")

	cmd.Flags().String(amapkit.ConfigKeyAmapKey, "", "高德key")
}

func BindFlags(cmd *cobra.Command) {
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		logkit.Error(err.Error())
	}
}
