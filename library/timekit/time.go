package timekit

import (
	"time"
)

func Sleep(millisecond int64) {
	time.Sleep(time.Duration(millisecond) * time.Millisecond)
}

// 毫秒时间戳解析为Time
func UnixMill(t int64) time.Time {
	return time.Unix(t/1000, t%1000*1000000)
}

// 用cast
func ParseString(dtString string, layout string) (time.Time, error) {
	return time.ParseInLocation(layout, dtString, time.Local)
}
