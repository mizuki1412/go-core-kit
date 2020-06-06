package class

import (
	"database/sql/driver"
	"github.com/spf13/cast"
	"mizuki/framework/core-kit/library/jsonkit"
)

// 同时继承scan和value方法
type Float64 struct {
	Float64 float64
	Valid   bool
}

func (th Float64) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.JSON().Marshal(th.Float64)
	}
	return jsonkit.JSON().Marshal(nil)
}
func (th *Float64) UnmarshalJSON(data []byte) error {
	var s *float64
	if err := jsonkit.JSON().Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		th.Valid = true
		th.Float64 = *s
	} else {
		th.Valid = false
	}
	return nil
}
func (th *Float64) Scan(value interface{}) error {
	if value == nil {
		th.Float64, th.Valid = 0, false
		return nil
	}
	th.Valid = true
	th.Float64 = cast.ToFloat64(value)
	return nil
}

// Value implements the driver Valuer interface.
func (th Float64) Value() (driver.Value, error) {
	if !th.Valid {
		return nil, nil
	}
	return th.Float64, nil
}

func (th *Float64) Set(val float64) {
	th.Float64 = val
	th.Valid = true
}
