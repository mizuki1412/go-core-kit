package class

import (
	"database/sql/driver"
	"github.com/spf13/cast"
	"mizuki/project/core-kit/library/jsonkit"
)

// 同时继承scan和value方法
type String struct {
	String string
	Valid  bool
}

func (th String) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.JSON().Marshal(th.String)
	}
	// 返回json中的null
	return jsonkit.JSON().Marshal(nil)
	//return nil,nil
}
func (th *String) UnmarshalJSON(data []byte) error {
	var s *string
	if err := jsonkit.JSON().Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		th.Valid = true
		th.String = *s
	} else {
		th.Valid = false
	}
	return nil
}
func (th *String) Scan(value interface{}) error {
	if value == nil {
		th.String, th.Valid = "", false
		return nil
	}
	th.Valid = true
	th.String = cast.ToString(value)
	return nil
}

// Value implements the driver Valuer interface.
func (th String) Value() (driver.Value, error) {
	if !th.Valid {
		return nil, nil
	}
	return th.String, nil
}
