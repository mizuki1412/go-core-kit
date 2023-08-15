package logkit

// logger的抽象

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/cli/configkey"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/mizuki1412/go-core-kit/library/timekit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
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
	var core zapcore.Core
	level := zap.InfoLevel
	switch configkit.GetString(configkey.LogLevel) {
	case "debug":
		level = zap.DebugLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	}
	if !configkit.Exist(configkey.LogPath) {
		core = zapcore.NewTee(
			// console中基本展示
			zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), level),
		)
	} else {
		fmt.Println(111)
		core = zapcore.NewTee(
			// 日志中json方式输出
			zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(getWriter2()), level),
			// console中基本展示
			zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), level),
		)
	}
	Logger = zap.New(core)
	return Logger
}

func getWriter2() io.Writer {
	filename := configkit.GetString(configkey.LogName, "main")
	filepath := configkit.GetString(configkey.LogPath)
	if stringkit.IsNull(filepath) {
		filepath = configkit.GetString(configkey.ProjectDir) + "/log"
	}
	filepath = stringkit.ClearFilePath(filepath)
	config := &lumberjack.Logger{
		Filename:   filepath + "/" + filename + ".log",
		MaxSize:    configkit.GetInt(configkey.LogMaxSize),
		MaxBackups: configkit.GetInt(configkey.LogMaxBackups),
		MaxAge:     configkit.GetInt(configkey.LogMaxRemain),
		LocalTime:  true,
		Compress:   true,
	}
	return config
}

type Param struct {
	Key string
	Val any
}

func Debug(msg any, params ...Param) {
	if Logger == nil {
		Logger = Init()
	}
	if len(params) > 0 {
		fields := transfer(params, len(params))
		Logger.Debug(cast.ToString(msg), fields...)
	} else {
		Logger.Debug(cast.ToString(msg))
	}
}
func DebugConcat(msg ...string) {
	Debug(strings.Join(msg, " "))
}
func Info(msg any, params ...Param) {
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
func Error(msg any, params ...Param) {
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
func Fatal(msg any, params ...Param) {
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
