package class

import (
	"database/sql"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/timekit"
	"github.com/spf13/cast"
	"strings"
	"time"
)

// 同时继承scan和value方法
type Time struct {
	sql.NullTime
}

func (th Time) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.JSON().Marshal(th.Time.Format(timekit.TimeLayoutAll))
	}
	return []byte("null"), nil
}
func (th *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		th.Valid = false
		return nil
	}
	var s time.Time
	var err error
	// 日期时间格式 + 毫秒形式
	str := unquoteIfQuoted(data)
	if len(str) == 13 && strings.Index(str, "-") < 0 {
		s0, err := cast.ToInt64E(str)
		if err != nil {
			return err
		}
		s = timekit.UnixMill(s0)
	} else {
		s, err = cast.StringToDate(str)
		if err != nil {
			return err
		}
	}
	th.Valid = true
	th.Time = s
	return nil
}

func (th *Time) Set(val time.Time) {
	th.Time = val
	th.Valid = true
}
