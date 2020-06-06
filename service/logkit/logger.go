package logkit

// logger的抽象

import (
	"github.com/arthurkiller/rollingwriter"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"mizuki/framework/core-kit/library/timekit"
	"mizuki/framework/core-kit/service/configkit"
	"os"
	"strings"
	"time"
)

var Logger *zap.Logger

func Init() *zap.Logger {
	config := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "Logger",
		MessageKey:    "msg",
		StacktraceKey: "trace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(t.Format(timekit.TimeLayoutAll))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewTee(
		// 日志中json方式输出
		zapcore.NewCore(zapcore.NewJSONEncoder(config), zapcore.AddSync(getWriter()), zap.InfoLevel),
		// console中基本展示
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), zap.InfoLevel),
	)
	Logger = zap.New(core)
	return Logger
}
func getWriter() io.Writer {
	filename := configkit.GetString(ConfigKeyLogName, "main")
	filepath := configkit.GetString(ConfigKeyLogPath, "./log")
	config := rollingwriter.Config{
		LogPath:       filepath,                                   //日志路径
		TimeTagFormat: "20060102",                                 //时间格式串
		FileName:      filename,                                   //日志文件名
		MaxRemain:     configkit.GetInt(ConfigKeyLogMaxRemain, 0), //配置日志最大存留数 0 取消自动清理
		// - 时间滚动: 配置策略如同 crontable, 例如,每天0:0切分, 则配置 0 0 0 * * *
		// - 大小滚动: 配置单个日志文件(未压缩)的滚动大小门限, 如1G, 500M
		RollingPolicy:      rollingwriter.TimeRolling, //配置滚动策略 norolling timerolling volumerolling
		RollingTimePattern: "0 0 0 * * *",             //配置时间滚动策略
		RollingVolumeSize:  "10M",                     //配置截断文件下限大小
		// writer 支持4种不同的 mode: 1. none 2. lock  3. async 4. buffer
		// - 无保护的 writer: 不提供并发安全保障
		// - lock 保护的 writer: 提供由 mutex 保护的并发安全保障
		// - 异步 writer: 异步 write, 并发安全. 异步开启后忽略 Lock 选项
		WriterMode: "lock",
		// BufferWriterThershould in B
		BufferWriterThershould: 8 * 1024 * 1024,
		// Compress will compress log file with gzip
		Compress: true,
	}
	// 创建一个 writer
	writer, err := rollingwriter.NewWriterFromConfig(&config)
	if err != nil {
		panic("rollingwriter.NewWriterFromConfig: " + err.Error())
	}
	return writer
}

type Param struct {
	Key string
	Val interface{}
}

func Info(msg interface{}, params ...Param) {
	if Logger == nil {
		Logger = Init()
	}
	if len(params) > 0 {
		fields := transfer(params, len(params))
		Logger.Info(cast.ToString(msg), fields...)
	} else {
		Logger.Info(cast.ToString(msg))
	}
}
func InfoConcat(msg ...string) {
	Info(strings.Join(msg, " "))
}
func Error(msg interface{}, params ...Param) {
	if Logger == nil {
		Logger = Init()
	}
	if len(params) > 0 {
		fields := transfer(params, len(params))
		Logger.Error(cast.ToString(msg), fields...)
	} else {
		Logger.Error(cast.ToString(msg))
	}
}
func ErrorConcat(msg ...string) {
	Error(strings.Join(msg, " "))
}
func Fatal(msg interface{}, params ...Param) {
	if Logger == nil {
		Logger = Init()
	}
	if len(params) > 0 {
		fields := transfer(params, len(params))
		Logger.Fatal(cast.ToString(msg), fields...)
	} else {
		Logger.Fatal(cast.ToString(msg))
	}
}

func transfer(params []Param, length int) []zapcore.Field {
	fields := make([]zapcore.Field, length)
	for i, p := range params {
		switch p.Val.(type) {
		case string:
			fields[i] = zap.String(p.Key, p.Val.(string))
		case int:
			fields[i] = zap.Int(p.Key, p.Val.(int))
		case float64:
			fields[i] = zap.Float64(p.Key, p.Val.(float64))
		case float32:
			fields[i] = zap.Float32(p.Key, p.Val.(float32))
		default:
			fields[i] = zap.String(p.Key, cast.ToString(p.Val))
		}
	}
	return fields
}
