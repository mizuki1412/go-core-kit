package timekit

import "fmt"

// 毫秒格式化
func FormatMillSecondHMS(mill int64) string {
	mill = mill / 1000
	hh := mill / 60 / 60
	mm := (mill / 60) % 60
	ss := mill % 60
	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
}
