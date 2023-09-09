package class

import (
	"database/sql"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/class/utils"
	"github.com/spf13/cast"
)

// Int32 同时继承scan和value方法
type Int32 struct {
	sql.NullInt32
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

func (th Int32) IsValid() bool {
	return th.Valid
}

func NewInt32(val ...any) Int32 {
	th := Int32{}
	if len(val) > 0 {
		th.Set(val[0])
	}
	return th
}
func NInt32(val ...any) *Int32 {
	th := &Int32{}
	if len(val) > 0 {
		th.Set(val[0])
	}
	return th
}

func (th *Int32) Set(val any) {
	switch val.(type) {
	case Int32:
		v := val.(Int32)
		th.Int32 = v.Int32
		th.Valid = v.Valid
	case *Int32:
		v := val.(*Int32)
		th.Int32 = v.Int32
		th.Valid = v.Valid
	default:
		i, err := cast.ToInt32E(val)
		if err == nil {
			th.Int32 = i
			th.Valid = true
		} else {
			panic(exception.New(err.Error()))
		}
	}
}
