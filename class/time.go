package class

import (
	"database/sql"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/class/utils"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/timekit"
	"github.com/spf13/cast"
	"time"
)

type Time struct {
	sql.NullTime
}

func (th Time) MarshalJSON() ([]byte, error) {
	if th.Valid {
		// 兼容iphone
		return jsonkit.Marshal(th.Time.Format(timekit.TimeLayout2))
	}
	return []byte("null"), nil
}

func (th *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		th.Valid = false
		return nil
	}
	str := utils.UnquoteIfQuoted(data)
	s, err := timekit.Parse(str)
	if err != nil {
		return err
	}
	th.Valid = true
	th.Time = s
	return nil
}

// todo 默认实现对时区没控制?
//func (th *Time) Scan(value any) error {
//	if value == nil {
//		th.Time, th.Valid = time.Time{}, false
//		return nil
//	}
//	var s time.Time
//	var err error
//	if v, ok := value.(time.Time); ok {
//		// todo 默认时区是0000的
//		if v.Location().String() == "" {
//			s = time.Date(v.Year(), v.Month(), v.Day(), v.Hour(), v.Minute(), v.Second(), v.Nanosecond(), time.Local)
//		} else {
//			s = v
//		}
//	} else {
//		s, err = timekit.Parse(cast.ToString(value))
//		if err != nil {
//			return err
//		}
//	}
//	th.Valid = true
//	th.Time = s
//	return err
//}

//func (th Time) IsValid() bool {
//	return th.Valid && !th.Time.IsZero()
//}

func NewTime(val any) Time {
	th := Time{}
	if val != nil {
		th.Set(val)
	}
	return th
}
func NTime(val any) *Time {
	th := &Time{}
	if val != nil {
		th.Set(val)
	}
	return th
}

func (th *Time) Set(val any) {
	if v, ok := val.(Time); ok {
		th.Time = v.Time
		th.Valid = v.Valid
	} else if v, ok := val.(string); ok {
		t, err := timekit.Parse(v)
		if err != nil {
			panic(exception.New("class.Time set error: " + err.Error()))
		}
		th.Time = t
		th.Valid = true
	} else if v, ok := val.(int64); ok {
		th.Time = time.UnixMilli(v)
		th.Valid = true
	} else {
		t, err := cast.ToTimeE(val)
		if err != nil {
			panic(exception.New("class.Time set error: " + err.Error()))
		}
		th.Time = t
		th.Valid = true
	}
}

func (th *Time) UnixMill() int64 {
	return th.Time.UnixMilli()
}
