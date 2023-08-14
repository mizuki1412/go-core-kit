package configkit

import (
	"github.com/mizuki1412/go-core-kit/cli/configkey"
	"time"
)

// Deprecated
const ConfigKeyTimeLocation = "time.location"

func GetLocation() *time.Location {
	loc, _ := time.LoadLocation(GetString(configkey.TimeLocation, "Asia/Shanghai"))
	return loc
}
