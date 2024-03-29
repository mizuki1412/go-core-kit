package class

import (
	"database/sql"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/class/utils"
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
	"github.com/spf13/cast"
)

type Float64 struct {
	sql.NullFloat64
}

// MarshalJSON 序列化为string
func (th Float64) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.Marshal(cast.ToString(th.Float64))
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

func (th Float64) IsValid() bool {
	return th.Valid
}

func NewFloat64(val ...any) Float64 {
	th := Float64{}
	if len(val) > 0 {
		th.Set(val[0])
	}
	return th
}
func NFloat64(val ...any) *Float64 {
	th := &Float64{}
	if len(val) > 0 {
		th.Set(val[0])
	}
	return th
}

func (th *Float64) Set(val any) {
	switch val.(type) {
	case Float64:
		v := val.(Float64)
		th.Float64 = v.Float64
		th.Valid = v.Valid
	case *Float64:
		v := val.(*Float64)
		th.Float64 = v.Float64
		th.Valid = v.Valid
	default:
		i, err := cast.ToFloat64E(val)
		if err == nil {
			th.Float64 = i
			th.Valid = true
		} else {
			panic(exception.New(err.Error()))
		}
	}
}
