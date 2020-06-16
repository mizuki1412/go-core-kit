package class

import (
	"database/sql/driver"
	"github.com/spf13/cast"
)

// 同时继承scan和value方法
type Int64 struct {
	Int64 int64
	Valid bool
}

func (th Int64) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return []byte(cast.ToString(th.Int64)), nil
	}
	return []byte("null"), nil
}
func (th *Int64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		th.Valid = false
		return nil
	}
	s, err := cast.ToInt64E(unquoteIfQuoted(data))
	if err != nil {
		return err
	}
	th.Valid = true
	th.Int64 = s
	return nil
}
func (th *Int64) Scan(value interface{}) error {
	if value == nil {
		th.Int64, th.Valid = 0, false
		return nil
	}
	th.Valid = true
	th.Int64 = cast.ToInt64(value)
	return nil
}

// Value implements the driver Valuer interface.
func (th Int64) Value() (driver.Value, error) {
	if !th.Valid {
		return nil, nil
	}
	return th.Int64, nil
}

func (th *Int64) Set(val int64) {
	th.Int64 = val
	th.Valid = true
}
