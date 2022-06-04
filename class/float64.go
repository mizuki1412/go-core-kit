package class

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/class/utils"
	"github.com/spf13/cast"
)

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
func (th *Float64) Scan(value any) error {
	if value == nil {
		th.Float64, th.Valid = 0, false
		return nil
	}
	th.Valid = true
	var err error
	switch value.(type) {
	case []uint8:
		// 数据库中decimal的值是字符数组返回
		a := value.([]uint8)
		th.Float64, err = cast.ToFloat64E(string(a))
	default:
		th.Float64, err = cast.ToFloat64E(value)
	}
	return err
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

func NewFloat64(val any) *Float64 {
	th := &Float64{}
	if val != nil {
		th.Set(val)
	}
	return th
}

func (th *Float64) Set(val any) *Float64 {
	if v, ok := val.(Float64); ok {
		th.Float64 = v.Float64
		th.Valid = v.Valid
	} else {
		i, err := cast.ToFloat64E(val)
		if err == nil {
			th.Float64 = i
			th.Valid = true
		} else {
			panic(exception.New(err.Error()))
		}
	}
	return th
}
