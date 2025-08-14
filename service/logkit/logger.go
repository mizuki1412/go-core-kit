package logkit

// logger的抽象

import (
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/cli/configkey"
	"github.com/mizuki1412/go-core-kit/v2/library/stringkit"
	"github.com/mizuki1412/go-core-kit/v2/library/timekit"
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
	"gopkg.in/natefinch/lumberjack.v2"
)

var once sync.Once
var fileLogger *slog.Logger

func Init() {
	once.Do(func() {
		var level slog.Level
		switch configkit.GetString(configkey.LogLevel) {
		case "debug":
			level = slog.LevelDebug
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			level = slog.LevelInfo
		}
		option := &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					a.Value = slog.AnyValue(a.Value.Time().Format(timekit.TimeLayout))
				}
				return a
			},
			Level: level,
		}
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, option)))
		if configkit.Exist(configkey.LogPath) {
			switch configkit.GetString(configkey.LogType) {
			case "json":
				fileLogger = slog.New(slog.NewJSONHandler(getRollWriter(), option))
			default:
				fileLogger = slog.New(slog.NewTextHandler(getRollWriter(), option))
			}
		}
	})
}

func getRollWriter() io.Writer {
	filename := configkit.GetString(configkey.LogName)
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

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
	if fileLogger != nil {
		fileLogger.Debug(msg, args...)
	}
}
func Info(msg string, args ...any) {
	slog.Info(msg, args...)
	if fileLogger != nil {
		fileLogger.Info(msg, args...)
	}
}
func Error(msg string, args ...any) {
	err := exception.New(msg, 2).Error()
	ErrorOrigin(err, args...)
}

func ErrorOrigin(msg string, args ...any) {
	slog.Error(msg, args...)
	if fileLogger != nil {
		fileLogger.Error(msg, args...)
	}
}
