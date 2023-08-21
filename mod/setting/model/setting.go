package model

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/spf13/cast"
)

type Setting struct {
	Id   int32           `json:"id,omitempty" db:"id" pk:"true" table:"more_setting"`
	Data class.MapString `json:"data,omitempty" db:"data"`
}

func (th *Setting) Scan(value any) error {
	if value == nil {
		return nil
	}
	th.Id = cast.ToInt32(value)
	return nil
}
func (th Setting) Value() (driver.Value, error) {
	return int64(th.Id), nil
}
