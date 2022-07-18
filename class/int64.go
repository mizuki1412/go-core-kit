package class

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/class/utils"
	"github.com/spf13/cast"
)

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
	s, err := cast.ToInt64E(utils.UnquoteIfQuoted(data))
	if err != nil {
		return err
	}
	th.Valid = true
	th.Int64 = s
	return nil
}
func (th *Int64) Scan(value any) error {
	if value == nil {
		th.Int64, th.Valid = 0, false
		return nil
	}
	th.Valid = true
	var err error
	switch value.(type) {
	case []uint8:
		th.Int64, err = cast.ToInt64E(string(value.([]uint8)))
	default:
		th.Int64, err = cast.ToInt64E(value)
	}
	return err
}

// Value implements the driver Valuer interface.
func (th Int64) Value() (driver.Value, error) {
	if !th.Valid {
		return nil, nil
	}
	return th.Int64, nil
}

func (th Int64) IsValid() bool {
	return th.Valid
}

func NewInt64(val any) *Int64 {
	th := &Int64{}
	if val != nil {
		th.Set(val)
	}
	return th
}

func (th *Int64) Set(val any) *Int64 {
	if v, ok := val.(Int64); ok {
		th.Int64 = v.Int64
		th.Valid = v.Valid
	} else {
		i, err := cast.ToInt64E(val)
		if err == nil {
			th.Int64 = i
			th.Valid = true
		} else {
			panic(exception.New(err.Error()))
		}
	}
	return th
}
