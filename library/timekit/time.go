package timekit

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/class/exception"
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
func ParseString(dtString string, layout string) time.Time {
	t, err := time.ParseInLocation(layout, dtString, time.Local)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	return t
}

// 修整为当日开始时间
func TrimDayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// 修整为下一日开始时间
func TrimDayNext(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())
}

// 毫秒格式化
func FormatMillSecondHMS(mill int64) string {
	mill = mill / 1000
	hh := mill / 60 / 60
	mm := (mill / 60) % 60
	ss := mill % 60
	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
}
