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
	cmd.PersistentFlags().String(configkey.ProjectName, "app", "项目名称")
	cmd.PersistentFlags().String(configkey.ProjectSubDir4PublicDownload, "", "项目目录中用于公共下载的开放目录（一层），逗号分隔，.表示所有")
	cmd.PersistentFlags().String(configkey.ProjectSubDir4PrivateDownload, "", "项目目录中用于私有下载的开放目录（一层），逗号分隔，.表示所有")
	cmd.PersistentFlags().Bool(configkey.ProfileDev, false, "开发模式 default:false")
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

	cmd.PersistentFlags().String(configkey.LogPath, "", "日志目录；空则表示在project.dir/log下；不填不开启文件日志")
	cmd.PersistentFlags().String(configkey.LogName, "main", "日志文件名，无后缀")
	cmd.PersistentFlags().Int(configkey.LogMaxRemain, 0, "最大保留天数")
	cmd.PersistentFlags().Int(configkey.LogMaxBackups, 0, "最大保留个数")
	cmd.PersistentFlags().Int(configkey.LogMaxSize, 20, "单文件最大尺寸")
	cmd.PersistentFlags().String(configkey.LogLevel, "", "日志等级 debug/info/warn/error")

	cmd.PersistentFlags().String(configkey.RestServerBase, "", "rest base url")
	cmd.PersistentFlags().String(configkey.RestServerPort, "10000", "")
	cmd.PersistentFlags().String(configkey.RestRequestBodySize, "", "限制request最大，单位MB")
	cmd.PersistentFlags().Bool(configkey.RestPPROF, false, "开启pprof, /debug/pprof")

	cmd.PersistentFlags().Int(configkey.JwtExpire, 6, "jwt 过期时间")
	cmd.PersistentFlags().String(configkey.JwtSecretKey, "0123456789abcdef", "jwt 密钥")

	cmd.PersistentFlags().String(configkey.DBDriver, "", "postgres/mysql/mssql")
	cmd.PersistentFlags().String(configkey.DBHost, "", "")
	cmd.PersistentFlags().String(configkey.DBPort, "", "")
	cmd.PersistentFlags().String(configkey.DBName, "", "")
	cmd.PersistentFlags().String(configkey.DBUser, "", "")
	cmd.PersistentFlags().String(configkey.DBPwd, "", "")
	cmd.PersistentFlags().Int(configkey.DBMaxOpen, 25, "最大连接")
	cmd.PersistentFlags().Int(configkey.DBMaxIdle, 5, "最大空闲连接")
	cmd.PersistentFlags().Int(configkey.DBMaxLife, 20, "单位/分钟")

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
	cmd.PersistentFlags().String(configkey.MQTTClientID, "client", "")
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
