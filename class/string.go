package class

import (
	"database/sql/driver"
	"github.com/spf13/cast"
)

// 同时继承scan和value方法
type String struct {
	String string
	Valid  bool
}

func (th String) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return []byte("\"" + th.String + "\""), nil
	}
	// 返回json中的null
	return []byte("null"), nil
	//return nil,nil
}
func (th *String) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		th.Valid = false
		return nil
	}
	th.String = unquoteIfQuoted(data)
	th.Valid = true
	return nil
}
func (th *String) Scan(value interface{}) error {
	if value == nil {
		th.String, th.Valid = "", false
		return nil
	}
	th.Valid = true
	th.String = cast.ToString(value)
	return nil
}

// Value implements the driver Valuer interface.
func (th String) Value() (driver.Value, error) {
	if !th.Valid {
		return nil, nil
	}
	return th.String, nil
}

func (th *String) Set(val string) {
	th.String = val
	th.Valid = true
}
