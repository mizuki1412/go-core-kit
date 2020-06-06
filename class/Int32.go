package class

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/spf13/cast"
)

// 同时继承scan和value方法
type Int32 struct {
	Int32 int32
	Valid bool
}

func (th Int32) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.JSON().Marshal(th.Int32)
	}
	return jsonkit.JSON().Marshal(nil)
}
func (th *Int32) UnmarshalJSON(data []byte) error {
	var s *int32
	if err := jsonkit.JSON().Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		th.Valid = true
		th.Int32 = *s
	} else {
		th.Valid = false
	}
	return nil
}
func (th *Int32) Scan(value interface{}) error {
	if value == nil {
		th.Int32, th.Valid = 0, false
		return nil
	}
	th.Valid = true
	th.Int32 = cast.ToInt32(value)
	return nil
}

// Value implements the driver Valuer interface.
func (th Int32) Value() (driver.Value, error) {
	if !th.Valid {
		return nil, nil
	}
	return int64(th.Int32), nil
}

func (th *Int32) Set(val int32) {
	th.Int32 = val
	th.Valid = true
}
