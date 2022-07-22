package initkit

import (
	"github.com/mizuki1412/go-core-kit/init/configkey"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

// LoadConfig 注意，load比一般的init慢
func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	// 这里可以执行多次的 搜索多个地址
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		// 这里不可用logkit，因为环境变量没有设置完
		log.Println("miss config.json: " + err.Error())
	}
}

func DefFlags(cmd *cobra.Command) {
	// todo 默认值不起效？
	cmd.Flags().String(configkey.ProjectDir, ".", "项目目录")
	cmd.Flags().String(configkey.ProjectName, "", "项目名称")
	cmd.Flags().String(configkey.ProjectSubDir4PublicDownload, "", "项目目录中用于公共下载的开放目录（一层），逗号分隔，.表示所有")
	cmd.Flags().String(configkey.ProjectSubDir4PrivateDownload, "", "项目目录中用于私有下载的开放目录（一层），逗号分隔，.表示所有")
	cmd.Flags().String(configkey.ProfileDev, "", "开发模式 default:false")
	cmd.Flags().String(configkey.TimeLocation, "Asia/Shanghai", "项目中用到的时区")

	cmd.Flags().String(configkey.RedisPrefix, "", "redis key的前缀")
	cmd.Flags().String(configkey.RedisHost, "", "redis host")
	cmd.Flags().String(configkey.RedisPort, "", "")
	cmd.Flags().String(configkey.RedisDB, "", "redis db 数据库号")
	cmd.Flags().String(configkey.RedisPwd, "", "")

	cmd.Flags().String(configkey.InfluxURL, "", "")
	cmd.Flags().String(configkey.InfluxUser, "", "")
	cmd.Flags().String(configkey.InfluxPwd, "", "")
	cmd.Flags().String(configkey.InfluxDBName, "", "默认的数据库")

	cmd.Flags().String(configkey.LogPath, "", "日志目录；默认在project.dir下")
	cmd.Flags().String(configkey.LogName, "main", "日志文件名，无后缀")
	cmd.Flags().String(configkey.LogMaxRemain, "", "最大保留天数")
	cmd.Flags().String(configkey.LogMaxBackups, "", "最大保留个数")
	cmd.Flags().String(configkey.LogMaxSize, "", "单文件最大尺寸")
	cmd.Flags().String(configkey.LogFileOff, "", "关闭文件日志")
	cmd.Flags().String(configkey.LogLevel, "", "日志等级 debug/info/warn/error")

	cmd.Flags().String(configkey.RestServerBase, "", "rest base url")
	cmd.Flags().String(configkey.RestServerPort, "", "")
	cmd.Flags().String(configkey.RestRequestBodySize, "", "限制request最大，单位MB")
	cmd.Flags().String(configkey.RestPPROF, "", "开启pprof, /debug/pprof")
	cmd.Flags().String(configkey.SessionExpire, "", "session expire 单位小时")
	cmd.Flags().String(configkey.SessionSecure, "true", "上传cookie时是否需要https，关系到浏览器的跨域策略和具体是否用https部署服务")

	cmd.Flags().String(configkey.DBDriver, "", "")
	cmd.Flags().String(configkey.DBHost, "", "")
	cmd.Flags().String(configkey.DBPort, "", "")
	cmd.Flags().String(configkey.DBName, "", "")
	cmd.Flags().String(configkey.DBUser, "", "")
	cmd.Flags().String(configkey.DBPwd, "", "")
	cmd.Flags().String(configkey.DBMaxOpen, "", "最大连接 默认25")
	cmd.Flags().String(configkey.DBMaxIdle, "", "最大空闲连接 默认5")
	cmd.Flags().String(configkey.DBMaxLife, "", "单位分钟，默认20")

	cmd.Flags().String(configkey.SwaggerBasePath, "", "/path")
	cmd.Flags().String(configkey.SwaggerHost, "", "可选，默认按swagger-ui所在路径")
	cmd.Flags().String(configkey.SwaggerDescription, "", "")
	cmd.Flags().String(configkey.SwaggerTitle, "", "")
	cmd.Flags().String(configkey.SwaggerVersion, "1.0.0", "")

	cmd.Flags().String(configkey.AmapKey, "", "高德key")

	cmd.Flags().String(configkey.AliRegionId, "cn-hangzhou", "ali")
	cmd.Flags().String(configkey.AliAccessKey, "", "ali")
	cmd.Flags().String(configkey.AliAccessKeySecret, "", "ali")
	cmd.Flags().String(configkey.AliSMSTemplate1, "", "ali sms 模板1")
	cmd.Flags().String(configkey.AliSMSSign1, "", "ali sms 签名1")
	cmd.Flags().String(configkey.AliSTSRoleArn, "", "ali sts")
	cmd.Flags().String(configkey.AliOSSBucketName, "", "ali oss default bucket")

	// mqtt
	cmd.Flags().String(configkey.MQTTBroker, "", "eg: tcp://xx.xx.xx")
	cmd.Flags().String(configkey.MQTTClientID, "", "")
	cmd.Flags().String(configkey.MQTTUsername, "", "")
	cmd.Flags().String(configkey.MQTTPwd, "", "")
	// netkit
	cmd.Flags().String(configkey.NetPort, "", "")
}

func BindFlags(cmd *cobra.Command) {
	// 如果存在配置文件
	LoadConfig()
	// 从命令参数中导入
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		logkit.Error(err.Error())
	}
}
