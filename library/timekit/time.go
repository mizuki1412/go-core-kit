package timekit

import (
	"time"
)

func Sleep(millisecond int64) {
	time.Sleep(time.Duration(millisecond) * time.Millisecond)
}

// TrimDayStart 修整为当日开始时间
func TrimDayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// TrimMonthStart 修整为当月开始时间
func TrimMonthStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// TrimDayNext 修整为下一日开始时间
func TrimDayNext(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())
}

// TrimMonthNext 修整为下月开始时间
func TrimMonthNext(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
}

// CountDaysInMonth 计算一个月的最大天数
func CountDaysInMonth(year int, month time.Month) (days int) {
	if month != time.February {
		if month == time.April || month == time.June || month == time.September || month == time.November {
			return 30
		} else {
			return 31
		}
	} else {
		if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
			return 29
		} else {
			return 28
		}
	}
}

// MonthInterval 相差月份数, 不算day
func MonthInterval(t1 time.Time, t2 time.Time) int {
	yearInterval := t1.Year() - t2.Year()
	if yearInterval > 0 {
		monthInterval := int(t1.Month()) - int(t2.Month())
		return yearInterval*12 + monthInterval
	} else {
		yearInterval = 0 - yearInterval
		monthInterval := int(t2.Month()) - int(t1.Month())
		return yearInterval*12 + monthInterval
	}
}
