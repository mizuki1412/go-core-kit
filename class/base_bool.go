package class

import (
	"database/sql"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/class/utils"
	"github.com/spf13/cast"
)

type Bool struct {
	sql.NullBool
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

func (th Bool) IsValid() bool {
	return th.Valid
}

func NewBool(val ...any) Bool {
	th := Bool{}
	if len(val) > 0 {
		th.Set(val[0])
	}
	return th
}
func NBool(val ...any) *Bool {
	th := &Bool{}
	if len(val) > 0 {
		th.Set(val[0])
	}
	return th
}

func (th *Bool) Set(val any) {
	switch val.(type) {
	case Bool:
		v := val.(Bool)
		th.Bool = v.Bool
		th.Valid = v.Valid
	case *Bool:
		v := val.(*Bool)
		th.Bool = v.Bool
		th.Valid = v.Valid
	default:
		i, err := cast.ToBoolE(val)
		if err == nil {
			th.Bool = i
			th.Valid = true
		} else {
			panic(exception.New(err.Error()))
		}
	}
}
