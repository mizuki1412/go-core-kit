package timekit

import (
	"github.com/mizuki1412/go-core-kit/v2/cli/configkey"
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
	"time"
)

func GetLocation() *time.Location {
	loc, _ := time.LoadLocation(configkit.GetString(configkey.TimeLocation))
	return loc
}
