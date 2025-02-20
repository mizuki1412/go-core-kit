package timekit

import (
	"github.com/mizuki1412/go-core-kit/v2/cli/configkey"
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
	"time"
)

func GetLocation() *time.Location {
	// os.Setenv("TZ", "UTC") 环境变量的方式只能生效一次
	// time.Local 可能受实际运行环境影响
	loc, _ := time.LoadLocation(configkit.GetString(configkey.TimeLocation))
	return loc
}
