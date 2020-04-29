package class

import (
	"database/sql/driver"
	"github.com/spf13/cast"
	"mizuki/project/core-kit/library/jsonkit"
)

// 同时继承scan和value方法
type Bool struct {
	Bool  bool
	Valid bool
}

func (th Bool) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.JSON().Marshal(th.Bool)
	}
	return jsonkit.JSON().Marshal(nil)
}
func (th *Bool) UnmarshalJSON(data []byte) error {
	var s *bool
	if err := jsonkit.JSON().Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		th.Valid = true
		th.Bool = *s
	} else {
		th.Valid = false
	}
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
