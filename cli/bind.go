package cli

import (
	"github.com/mizuki1412/go-core-kit/cli/configkey"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 这里将在 run 之后执行
func loadConfig() {
	if viper.GetString("config") != "" {
		viper.SetConfigFile(viper.GetString("config"))
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}
	_ = viper.ReadInConfig()
}

func bindDefaultFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("config", "c", "", "配置文件全路径")
	cmd.PersistentFlags().String(configkey.ProjectDir, ".", "项目目录")
	cmd.PersistentFlags().String(configkey.ProjectName, "", "项目名称")
	cmd.PersistentFlags().String(configkey.ProjectSubDir4PublicDownload, "", "项目目录中用于公共下载的开放目录（一层），逗号分隔，.表示所有")
	cmd.PersistentFlags().String(configkey.ProjectSubDir4PrivateDownload, "", "项目目录中用于私有下载的开放目录（一层），逗号分隔，.表示所有")
	cmd.PersistentFlags().String(configkey.ProfileDev, "", "开发模式 default:false")
	cmd.PersistentFlags().String(configkey.TimeLocation, "Asia/Shanghai", "项目中用到的时区")

	cmd.PersistentFlags().String(configkey.RedisPrefix, "", "redis key的前缀")
	cmd.PersistentFlags().String(configkey.RedisHost, "", "redis host")
	cmd.PersistentFlags().String(configkey.RedisPort, "", "")
	cmd.PersistentFlags().String(configkey.RedisDB, "", "redis db 数据库号")
	cmd.PersistentFlags().String(configkey.RedisPwd, "", "")

	cmd.PersistentFlags().String(configkey.InfluxURL, "", "")
	cmd.PersistentFlags().String(configkey.InfluxUser, "", "")
	cmd.PersistentFlags().String(configkey.InfluxPwd, "", "")
	cmd.PersistentFlags().String(configkey.InfluxDBName, "", "默认的数据库")

	cmd.PersistentFlags().String(configkey.LogPath, "", "日志目录；默认在project.dir下")
	cmd.PersistentFlags().String(configkey.LogName, "main", "日志文件名，无后缀")
	cmd.PersistentFlags().String(configkey.LogMaxRemain, "", "最大保留天数")
	cmd.PersistentFlags().String(configkey.LogMaxBackups, "", "最大保留个数")
	cmd.PersistentFlags().String(configkey.LogMaxSize, "", "单文件最大尺寸")
	cmd.PersistentFlags().String(configkey.LogFileOff, "", "关闭文件日志")
	cmd.PersistentFlags().String(configkey.LogLevel, "", "日志等级 debug/info/warn/error")

	cmd.PersistentFlags().String(configkey.RestServerBase, "", "rest base url")
	cmd.PersistentFlags().String(configkey.RestServerPort, "", "")
	cmd.PersistentFlags().String(configkey.RestRequestBodySize, "", "限制request最大，单位MB")
	cmd.PersistentFlags().String(configkey.RestPPROF, "", "开启pprof, /debug/pprof")
	cmd.PersistentFlags().String(configkey.SessionExpire, "", "session expire 单位小时，默认12小时")
	cmd.PersistentFlags().String(configkey.SessionSecure, "true", "上传cookie时是否需要https，关系到浏览器的跨域策略和具体是否用https部署服务")

	cmd.PersistentFlags().String(configkey.DBDriver, "", "")
	cmd.PersistentFlags().String(configkey.DBHost, "", "")
	cmd.PersistentFlags().String(configkey.DBPort, "", "")
	cmd.PersistentFlags().String(configkey.DBName, "", "")
	cmd.PersistentFlags().String(configkey.DBUser, "", "")
	cmd.PersistentFlags().String(configkey.DBPwd, "", "")
	cmd.PersistentFlags().String(configkey.DBMaxOpen, "", "最大连接 默认25")
	cmd.PersistentFlags().String(configkey.DBMaxIdle, "", "最大空闲连接 默认5")
	cmd.PersistentFlags().String(configkey.DBMaxLife, "", "单位分钟，默认20")

	cmd.PersistentFlags().String(configkey.SwaggerBasePath, "", "/path")
	cmd.PersistentFlags().String(configkey.SwaggerHost, "", "可选，默认按swagger-ui所在路径")
	cmd.PersistentFlags().String(configkey.SwaggerDescription, "", "")
	cmd.PersistentFlags().String(configkey.SwaggerTitle, "", "")
	cmd.PersistentFlags().String(configkey.SwaggerVersion, "1.0.0", "")

	cmd.PersistentFlags().String(configkey.AmapKey, "", "高德key")

	cmd.PersistentFlags().String(configkey.AliRegionId, "cn-hangzhou", "ali")
	cmd.PersistentFlags().String(configkey.AliAccessKey, "", "ali")
	cmd.PersistentFlags().String(configkey.AliAccessKeySecret, "", "ali")
	cmd.PersistentFlags().String(configkey.AliSMSTemplate1, "", "ali sms 模板1")
	cmd.PersistentFlags().String(configkey.AliSMSSign1, "", "ali sms 签名1")
	cmd.PersistentFlags().String(configkey.AliSTSRoleArn, "", "ali sts")
	cmd.PersistentFlags().String(configkey.AliOSSBucketName, "", "ali oss default bucket")

	// mqtt
	cmd.PersistentFlags().String(configkey.MQTTBroker, "", "eg: tcp://xx.xx.xx")
	cmd.PersistentFlags().String(configkey.MQTTClientID, "", "")
	cmd.PersistentFlags().String(configkey.MQTTUsername, "", "")
	cmd.PersistentFlags().String(configkey.MQTTPwd, "", "")
	// netkit
	cmd.PersistentFlags().String(configkey.NetPort, "", "")

	// softether
	cmd.PersistentFlags().String(configkey.SoftEtherHost, "", "")
	cmd.PersistentFlags().String(configkey.SoftEtherPort, "", "")
	cmd.PersistentFlags().String(configkey.SoftEtherPwd, "", "")
	cmd.PersistentFlags().String(configkey.SoftEtherOpenVpnPort, "", "")

	bind(cmd)
}

func bind(cmd *cobra.Command) {
	err := viper.BindPFlags(cmd.PersistentFlags())
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlags(cmd.Flags())
	if err != nil {
		panic(err)
	}
}
