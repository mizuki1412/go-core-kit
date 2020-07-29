package timekit

import (
	"fmt"
	"github.com/spf13/cast"
	"regexp"
	"time"
)

func Sleep(millisecond int64) {
	time.Sleep(time.Duration(millisecond) * time.Millisecond)
}

// 毫秒时间戳解析为Time
func UnixMill(t int64) time.Time {
	return time.Unix(t/1000, t%1000*1000000)
}

// cast.StringToDate 不能设置时区
func ParseString(dtString string, layout string) (time.Time, error) {
	return time.ParseInLocation(layout, dtString, time.Local)
}

func Parse(dt string) (time.Time, error) {
	var s time.Time
	// 日期时间格式 + 毫秒形式
	if ok, _ := regexp.Match("^[\\d]{13}$", []byte(dt)); ok {
		s0, err := cast.ToInt64E(dt)
		if err != nil {
			return s, err
		}
		s = UnixMill(s0)
		return s, nil
	} else {
		for _, dateType := range []string{
			time.RFC3339,
			"2006-01-02T15:04:05", // iso8601 without timezone
			time.RFC1123Z,
			time.RFC1123,
			time.RFC822Z,
			time.RFC822,
			time.RFC850,
			time.ANSIC,
			time.UnixDate,
			time.RubyDate,
			"2006-01-02 15:04:05.999999999 -0700 MST", // Time.String()
			"2006-01-02",
			"02 Jan 2006",
			"2006-01-02T15:04:05-0700", // RFC3339 without timezone hh:mm colon
			"2006-01-02 15:04:05 -07:00",
			"2006-01-02 15:04:05 -0700",
			"2006-01-02 15:04:05Z07:00", // RFC3339 without T
			"2006-01-02 15:04:05Z0700",  // RFC3339 without T or timezone hh:mm colon
			"2006-01-02 15:04:05",
			time.Kitchen,
			time.Stamp,
			time.StampMilli,
			time.StampMicro,
			time.StampNano,
		} {
			if t, e := time.ParseInLocation(dateType, dt, time.Local); e == nil {
				return t, e
			}
		}
		return s, fmt.Errorf("unable to parse date: %s", dt)
	}
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
