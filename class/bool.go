package class

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/class/utils"
	"github.com/spf13/cast"
)

// 同时继承scan和value方法
type Bool struct {
	Bool  bool
	Valid bool
}

func (th Bool) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return []byte(cast.ToString(th.Bool)), nil
	}
	return []byte("null"), nil
}
func (th *Bool) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		th.Valid = false
		return nil
	}
	s, err := cast.ToBoolE(utils.UnquoteIfQuoted(data))
	if err != nil {
		return err
	}
	th.Valid = true
	th.Bool = s
	return nil
}
func (th *Bool) Scan(value interface{}) error {
	if value == nil {
		th.Bool, th.Valid = false, false
		return nil
	}
	th.Valid = true
	th.Bool = cast.ToBool(value)
	return nil
}

// Value implements the driver Valuer interface.
func (th Bool) Value() (driver.Value, error) {
	if !th.Valid {
		return nil, nil
	}
	return th.Bool, nil
}

func (th Bool) IsValid() bool {
	return th.Valid
}

func (th *Bool) Set(val interface{}) {
	if v, ok := val.(Bool); ok {
		th.Bool = v.Bool
		th.Valid = true
	} else {
		i, err := cast.ToBoolE(val)
		if err != nil {
			panic(exception.New("class.Bool set error"))
		}
		th.Bool = i
		th.Valid = true
	}
}
