package class

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/class/utils"
	"github.com/spf13/cast"
)

// 同时继承scan和value方法
type Float64 struct {
	Float64 float64
	Valid   bool
}

func (th Float64) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return []byte(cast.ToString(th.Float64)), nil
	}
	return []byte("null"), nil
}
func (th *Float64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		th.Valid = false
		return nil
	}
	s, err := cast.ToFloat64E(utils.UnquoteIfQuoted(data))
	if err != nil {
		return err
	}
	th.Valid = true
	th.Float64 = s
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

func (th Float64) IsValid() bool {
	return th.Valid
}

func (th *Float64) Set(val interface{}) {
	if v, ok := val.(Float64); ok {
		th.Float64 = v.Float64
		th.Valid = true
	} else {
		i, err := cast.ToFloat64E(val)
		if err != nil {
			panic(exception.New("class.Float64 set error"))
		}
		th.Float64 = i
		th.Valid = true
	}
}
