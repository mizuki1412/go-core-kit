package timekit

import (
	"github.com/mizuki1412/go-core-kit/cli/configkey"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"time"
)

func GetLocation() *time.Location {
	loc, _ := time.LoadLocation(configkit.GetString(configkey.TimeLocation, "Asia/Shanghai"))
	return loc
}
