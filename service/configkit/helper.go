package configkit

import "time"

const ConfigKeyTimeLocation = "time.location"

func GetLocation() *time.Location {
	loc, _ := time.LoadLocation(GetString(ConfigKeyTimeLocation, "Asia/Shanghai"))
	return loc
}
