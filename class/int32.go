package class

import (
	"database/sql/driver"
	"github.com/spf13/cast"
)

// 同时继承scan和value方法
type Int32 struct {
	Int32 int32
	Valid bool
}

func (th Int32) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return []byte(cast.ToString(th.Int32)), nil
	}
	return []byte("null"), nil
}
func (th *Int32) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		th.Valid = false
		return nil
	}
	s, err := cast.ToInt32E(unquoteIfQuoted(data))
	if err != nil {
		return err
	}
	th.Valid = true
	th.Int32 = s
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
