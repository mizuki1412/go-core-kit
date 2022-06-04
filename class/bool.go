package class

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/class/utils"
	"github.com/spf13/cast"
)

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
func (th *Bool) Scan(value any) error {
	if value == nil {
		th.Bool, th.Valid = false, false
		return nil
	}
	th.Valid = true
	var err error
	th.Bool, err = cast.ToBoolE(value)
	return err
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

func NewBool(val any) *Bool {
	th := &Bool{}
	if val != nil {
		th.Set(val)
	}
	return th
}

func (th *Bool) Set(val any) *Bool {
	if v, ok := val.(Bool); ok {
		th.Bool = v.Bool
		th.Valid = v.Valid
	} else {
		i, err := cast.ToBoolE(val)
		if err == nil {
			th.Bool = i
			th.Valid = true
		} else {
			panic(exception.New(err.Error()))
		}
	}
	return th
}
