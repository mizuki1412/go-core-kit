package class

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/class/utils"
	"github.com/spf13/cast"
)

// Int32 同时继承scan和value方法
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
	s, err := cast.ToInt32E(utils.UnquoteIfQuoted(data))
	if err != nil {
		return err
	}
	th.Valid = true
	th.Int32 = s
	return nil
}
func (th *Int32) Scan(value any) error {
	if value == nil {
		th.Int32, th.Valid = 0, false
		return nil
	}
	th.Valid = true
	var err error
	th.Int32, err = cast.ToInt32E(value)
	return err
}

// Value implements the driver Valuer interface.
func (th Int32) Value() (driver.Value, error) {
	if !th.Valid {
		return nil, nil
	}
	return int64(th.Int32), nil
}

func (th Int32) IsValid() bool {
	return th.Valid
}

func NewInt32(val any) *Int32 {
	th := &Int32{}
	if val != nil {
		th.Set(val)
	}
	return th
}

func (th *Int32) Set(val any) *Int32 {
	if v, ok := val.(Int32); ok {
		th.Int32 = v.Int32
		th.Valid = v.Valid
	} else {
		i, err := cast.ToInt32E(val)
		if err == nil {
			th.Int32 = i
			th.Valid = true
		} else {
			panic(exception.New(err.Error()))
		}
	}
	return th
}
