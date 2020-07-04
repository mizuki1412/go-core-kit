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

// 修整为当日开始时间
func TrimDayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// 修整为下一日开始时间
func TrimDayNext(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())
}
