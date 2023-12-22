package timekit

import (
	"fmt"
	"github.com/spf13/cast"
	"regexp"
	"time"
)

// ParseString cast.StringToDate 不能设置时区
func ParseString(dtString string, layout string) (time.Time, error) {
	return time.ParseInLocation(layout, dtString, time.Local)
}

func Parse(dt string) (time.Time, error) {
	var s time.Time
	// 日期时间格式 + 毫秒形式
	if ok, _ := regexp.Match("^[\\d]+$", []byte(dt)); ok {
		s0, err := cast.ToInt64E(dt)
		if err != nil {
			return s, err
		}
		s = time.UnixMilli(s0)
		return s, nil
	} else {
		for _, dateType := range []string{
			time.RFC3339,
			time.RFC1123Z,
			time.RFC1123,
			time.RFC822Z,
			time.RFC822,
			time.RFC850,
			time.ANSIC,
			time.UnixDate,
			time.RubyDate,
			"2006-01-02 15:04:05.999999999 -0700 MST", // Time.String()
			"02 Jan 2006",
			"2006-01-02T15:04:05",      // iso8601 without timezone
			"2006-01-02T15:04:05-0700", // RFC3339 without timezone hh:mm colon
			"2006-01-02 15:04:05 -07:00",
			"2006-01-02 15:04:05 -0700",
			"2006-01-02 15:04:05Z07:00", // RFC3339 without T
			"2006-01-02 15:04:05Z0700",  // RFC3339 without T or timezone hh:mm colon
			"2006/01/02 15:04:05",
			time.DateTime,
			time.DateOnly,
			time.Kitchen,
			time.Stamp,
			time.StampMilli,
			time.StampMicro,
			time.StampNano,
		} {
			// ParseInLocation 在已有时区偏移时不会使用given location
			if t, e := time.ParseInLocation(dateType, dt, time.Local); e == nil {
				return t, e
			}
		}
		return s, fmt.Errorf("unable to parse date: %s", dt)
	}
}

func ParseD(dt string) time.Time {
	t, _ := Parse(dt)
	return t
}
