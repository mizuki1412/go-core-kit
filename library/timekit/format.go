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

func FormatSecondHMS(second int64, zh bool) string {
	hh := second / 60 / 60
	mm := (second / 60) % 60
	ss := second % 60
	if zh {
		return fmt.Sprintf("%02d时%02d分%02d秒", hh, mm, ss)
	} else {
		return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
	}
}

func FormatSecondHM(second int64, zh bool) string {
	hh := second / 60 / 60
	mm := (second / 60) % 60
	if zh {
		return fmt.Sprintf("%02d时%02d分", hh, mm)
	} else {
		return fmt.Sprintf("%02d:%02d", hh, mm)
	}
}
