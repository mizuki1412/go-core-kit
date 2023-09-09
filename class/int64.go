package class

import (
	"database/sql"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/class/utils"
	"github.com/spf13/cast"
)

type Int64 struct {
	sql.NullInt64
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

func (th Int64) IsValid() bool {
	return th.Valid
}

func NewInt64(val ...any) Int64 {
	th := Int64{}
	if len(val) > 0 {
		th.Set(val[0])
	}
	return th
}
func NInt64(val ...any) *Int64 {
	th := &Int64{}
	if len(val) > 0 {
		th.Set(val[0])
	}
	return th
}

func (th *Int64) Set(val any) {
	switch val.(type) {
	case Int64:
		v := val.(Int64)
		th.Int64 = v.Int64
		th.Valid = v.Valid
	case *Int64:
		v := val.(*Int64)
		th.Int64 = v.Int64
		th.Valid = v.Valid
	default:
		i, err := cast.ToInt64E(val)
		if err == nil {
			th.Int64 = i
			th.Valid = true
		} else {
			panic(exception.New(err.Error()))
		}
	}
}
